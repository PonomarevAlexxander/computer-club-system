package handler

import (
	"regexp"
	"strconv"

	"github.com/PonomarevAlexxander/computer-club-system/internal/command"
	"github.com/PonomarevAlexxander/computer-club-system/internal/config"
	"github.com/PonomarevAlexxander/computer-club-system/internal/errlib"
)

const (
	clientHasNoTable uint = 0
)

const (
	ClientArrivedInputCommandId uint = iota + 1
	ClientTookTableInputCommandId
	ClientWaitingInputCommandId
	ClientLeftInputCommandId
)

const (
	ClientLeftOutputCommandId uint = iota + 11
	ClientTookTableOutputCommandId
	ErrorOutputCommandId
)

func NewCommandHandlers() []CommandHandler {
	return []CommandHandler{
		NewClientArrivedHandler(ClientArrivedInputCommandId),
		NewClientWaitingHandler(ClientWaitingInputCommandId),
		NewClientTookTableHandler(ClientTookTableInputCommandId),
		NewClientLeftHandler(ClientLeftInputCommandId),
	}
}

type CommandHandler interface {
	// Process returns: new command if it is needed to generate output or error if one happens
	Process(cmd *command.Command, clubStateDto *ClubStateDto) (*command.Command, error)
	Handles(cmd *command.Command) bool
}

type commonCommandHandler struct {
	commandId uint
}

func (e *commonCommandHandler) Handles(cmd *command.Command) bool {
	return e.commandId == cmd.Id
}

type ClientArrivedHandler struct {
	commonCommandHandler
}

func NewClientArrivedHandler(id uint) *ClientArrivedHandler {
	return &ClientArrivedHandler{
		commonCommandHandler{
			commandId: id,
		},
	}
}

func (e *ClientArrivedHandler) Process(cmd *command.Command, clubStateDto *ClubStateDto) (*command.Command, error) {
	if len(cmd.Params) != 1 {
		return nil, errlib.ErrInputFormat
	}

	clientName := cmd.Params[0]
	match, err := regexp.MatchString(config.ClientNameRegex, clientName)
	if !match || err != nil {
		return nil, errlib.ErrInputFormat
	}

	if cmd.Time.Before(clubStateDto.openingTime) || cmd.Time.After(clubStateDto.closingTime) {
		return command.NewCommand(cmd.Time, ErrorOutputCommandId, errlib.ErrClubClosed.Error()), nil
	}

	_, alreadyInClub := clubStateDto.clientToTable[clientName]
	if alreadyInClub {
		return command.NewCommand(cmd.Time, ErrorOutputCommandId, errlib.ErrClientAlreadyInClub.Error()), nil
	}

	clubStateDto.clientToTable[clientName] = clientHasNoTable
	return nil, nil
}

type ClientTookTableHandler struct {
	commonCommandHandler
}

func NewClientTookTableHandler(id uint) *ClientTookTableHandler {
	return &ClientTookTableHandler{
		commonCommandHandler{
			commandId: id,
		},
	}
}

func (e *ClientTookTableHandler) Process(cmd *command.Command, clubStateDto *ClubStateDto) (*command.Command, error) {
	if len(cmd.Params) != 2 {
		return nil, errlib.ErrInputFormat
	}

	clientName := cmd.Params[0]
	match, err := regexp.MatchString(config.ClientNameRegex, clientName)
	if !match || err != nil {
		return nil, errlib.ErrInputFormat
	}

	tempTableId, err := strconv.ParseUint(cmd.Params[1], 10, 0)
	if err != nil || tempTableId > uint64(clubStateDto.totalTables) || tempTableId < 1 {
		return nil, errlib.ErrInputFormat
	}
	tableId := uint(tempTableId)

	currentTable, clientInClub := clubStateDto.clientToTable[clientName]
	if !clientInClub {
		return command.NewCommand(cmd.Time, ErrorOutputCommandId, errlib.ErrUnknownClient.Error()), nil
	}

	_, tableTaken := clubStateDto.tableToTime[tableId]
	if tableTaken {
		return command.NewCommand(cmd.Time, ErrorOutputCommandId, errlib.ErrTableAlreadyTaken.Error()), nil
	}

	clubStateDto.tableToTime[tableId] = cmd.Time
	clubStateDto.clientToTable[clientName] = tableId
	if currentTable != clientHasNoTable {

		duration := cmd.Time.Sub(clubStateDto.tableToTime[currentTable])
		clubStateDto.tableStatistics[currentTable-1].AddTime(duration)
		clubStateDto.tableStatistics[currentTable-1].AddIncomeLowerBound(
			duration,
			clubStateDto.hourPrice,
		)

		delete(clubStateDto.tableToTime, currentTable)
	}
	return nil, nil
}

type ClientWaitingHandler struct {
	commonCommandHandler
}

func NewClientWaitingHandler(id uint) *ClientWaitingHandler {
	return &ClientWaitingHandler{
		commonCommandHandler{
			commandId: id,
		},
	}
}

func (e *ClientWaitingHandler) Process(cmd *command.Command, clubStateDto *ClubStateDto) (*command.Command, error) {
	if len(cmd.Params) != 1 {
		return nil, errlib.ErrInputFormat
	}

	clientName := cmd.Params[0]
	match, err := regexp.MatchString(config.ClientNameRegex, clientName)
	if !match || err != nil {
		return nil, errlib.ErrInputFormat
	}

	table, clientInClub := clubStateDto.clientToTable[clientName]
	if !clientInClub {
		return command.NewCommand(cmd.Time, ErrorOutputCommandId, errlib.ErrUnknownClient.Error()), nil
	}

	if len(clubStateDto.tableToTime) < int(clubStateDto.totalTables) {
		return command.NewCommand(cmd.Time, ErrorOutputCommandId, errlib.ErrWaitingTimeExceeded.Error()), nil
	}

	// In task specification this case is not mentioned
	if table != clientHasNoTable {
		return nil, nil
	}

	if clubStateDto.clientsQueue.Length() > int(clubStateDto.totalTables) {
		delete(clubStateDto.clientToTable, clientName)
		return command.NewCommand(cmd.Time, ClientLeftOutputCommandId, clientName), nil
	}

	clubStateDto.clientsQueue.Add(clientName)
	return nil, nil
}

type ClientLeftHandler struct {
	commonCommandHandler
}

func NewClientLeftHandler(id uint) *ClientLeftHandler {
	return &ClientLeftHandler{
		commonCommandHandler{
			commandId: id,
		},
	}
}

func (e *ClientLeftHandler) Process(cmd *command.Command, clubStateDto *ClubStateDto) (*command.Command, error) {
	if len(cmd.Params) != 1 {
		return nil, errlib.ErrInputFormat
	}

	clientName := cmd.Params[0]
	match, err := regexp.MatchString(config.ClientNameRegex, clientName)
	if !match || err != nil {
		return nil, errlib.ErrInputFormat
	}

	currentTable, clientInClub := clubStateDto.clientToTable[clientName]
	if !clientInClub {
		return command.NewCommand(cmd.Time, ErrorOutputCommandId, errlib.ErrUnknownClient.Error()), nil
	}

	if currentTable == clientHasNoTable {
		delete(clubStateDto.clientToTable, clientName)
		return nil, nil
	}

	duration := cmd.Time.Sub(clubStateDto.tableToTime[currentTable])
	clubStateDto.tableStatistics[currentTable-1].AddIncomeUpperBound(duration, clubStateDto.hourPrice)
	clubStateDto.tableStatistics[currentTable-1].AddTime(duration)

	delete(clubStateDto.clientToTable, clientName)
	delete(clubStateDto.tableToTime, currentTable)

	waitingClient := clubStateDto.clientsQueue.Poll()
	if waitingClient != nil {
		name := waitingClient.(string)
		clubStateDto.clientToTable[name] = currentTable
		clubStateDto.tableToTime[currentTable] = cmd.Time

		return command.NewCommand(
			cmd.Time,
			ClientTookTableOutputCommandId,
			name,
			strconv.FormatUint(uint64(currentTable), 10),
		), nil
	}

	return nil, nil
}
