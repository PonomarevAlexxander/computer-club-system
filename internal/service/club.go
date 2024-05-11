package service

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/PonomarevAlexxander/computer-club-system/internal/command"
	"github.com/PonomarevAlexxander/computer-club-system/internal/errlib"
	"github.com/PonomarevAlexxander/computer-club-system/internal/handler"
	"github.com/PonomarevAlexxander/computer-club-system/internal/queue"
	"github.com/PonomarevAlexxander/computer-club-system/internal/stats"
)

type IComputerClubService interface {
	Process(cmd *command.Command) (*command.Command, error)
	GetOpeningHours() time.Time
	GetClosingHours() time.Time
	CloseClub() []*command.Command
	GetStatistics() []stats.TableStats
}

type ComputerClubService struct {
	totalTables   uint
	hourPrice     uint
	clientToTable map[string]uint
	// tablesToTime - relation table to time when it was taken
	tableToTime     map[uint]time.Time
	clientsQueue    *queue.Queue
	openingTime     time.Time
	closingTime     time.Time
	tableStatistics []stats.TableStats
	handlers        []handler.CommandHandler
}

func NewComputerClubService(tables uint, openingHours time.Time, closingHours time.Time, hourPrice uint) *ComputerClubService {
	return &ComputerClubService{
		totalTables:     tables,
		openingTime:     openingHours,
		closingTime:     closingHours,
		hourPrice:       hourPrice,
		tableToTime:     make(map[uint]time.Time),
		clientToTable:   make(map[string]uint),
		clientsQueue:    queue.NewQueue(),
		tableStatistics: make([]stats.TableStats, tables),
		handlers:        handler.NewCommandHandlers(),
	}
}

func (c *ComputerClubService) Process(cmd *command.Command) (*command.Command, error) {
	currentState := handler.NewClubStateDto(
		c.totalTables,
		c.hourPrice,
		c.clientToTable,
		c.tableToTime,
		c.clientsQueue,
		c.openingTime,
		c.closingTime,
		c.tableStatistics,
	)

	for _, item := range c.handlers {
		if item.Handles(cmd) {
			return item.Process(cmd, currentState)
		}
	}
	return nil, fmt.Errorf("%w. Command: '%s'", errlib.ErrUnknownCommand, cmd)
}

func (c *ComputerClubService) CloseClub() []*command.Command {
	outCommands := make([]*command.Command, 0)

	for client, table := range c.clientToTable {
		duration := c.closingTime.Sub(c.tableToTime[table])

		if table > 0 {
			c.tableStatistics[table-1].AddIncomeUpperBound(duration, c.hourPrice)
			c.tableStatistics[table-1].AddTime(duration)
		}

		outCommands = append(outCommands, command.NewCommand(c.closingTime, handler.ClientLeftOutputCommandId, client))
		delete(c.clientToTable, client)
		delete(c.tableToTime, table)
	}

	slices.SortFunc[[]*command.Command](outCommands, func(l, r *command.Command) int {
		return cmp.Compare(strings.ToLower(l.Params[0]), strings.ToLower(r.Params[0]))
	})
	return outCommands
}

func (c *ComputerClubService) GetStatistics() []stats.TableStats {
	return c.tableStatistics
}

func (c *ComputerClubService) GetOpeningHours() time.Time {
	return c.openingTime
}

func (c *ComputerClubService) GetClosingHours() time.Time {
	return c.closingTime
}
