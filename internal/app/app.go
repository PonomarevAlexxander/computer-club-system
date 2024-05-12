package app

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/PonomarevAlexxander/computer-club-system/internal/command"
	"github.com/PonomarevAlexxander/computer-club-system/internal/config"
	"github.com/PonomarevAlexxander/computer-club-system/internal/errlib"
	"github.com/PonomarevAlexxander/computer-club-system/internal/scanner"
	"github.com/PonomarevAlexxander/computer-club-system/internal/service"
	"github.com/PonomarevAlexxander/computer-club-system/internal/util"
)

type ComputerClub struct {
	scanner scanner.IScanner
	out     io.Writer
	service service.IComputerClubService
}

func NewComputerClub(in scanner.IScanner, out io.Writer, service service.IComputerClubService) *ComputerClub {
	return &ComputerClub{
		scanner: in,
		out:     out,
		service: service,
	}
}

func (club *ComputerClub) Start() {
	fmt.Fprintln(club.out, club.service.GetOpeningHours().Format(config.TimeFormat))

	var (
		err      error     = nil
		prevTime time.Time = time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
	)

	for err != io.EOF {
		line, err := club.scanner.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			util.HandleInputError(club.out, line)
			return
		}

		cmd, err := command.CommandFromString(line)
		if err != nil {
			club.handleError(err, line)
			return
		}
		if cmd.Time.Before(prevTime) {
			util.HandleInputError(club.out, line)
			return
		}

		generatedCmd, err := club.service.Process(cmd)
		if err != nil {
			club.handleError(err, line)
			return
		}

		fmt.Fprintln(club.out, cmd)
		if generatedCmd != nil {
			fmt.Fprintln(club.out, generatedCmd)
		}

		prevTime = cmd.Time
	}

	for _, cmd := range club.service.CloseClub() {
		fmt.Fprintln(club.out, cmd)
	}

	fmt.Fprintln(club.out, club.service.GetClosingHours().Format(config.TimeFormat))

	for id, tableStat := range club.service.GetStatistics() {
		fmt.Fprintf(club.out, "%d %d %s\n", id+1, tableStat.Income, tableStat.TimeTaken.Format(config.TimeFormat))
	}
}

func (club *ComputerClub) handleError(err error, line string) {
	if errors.Is(err, errlib.ErrInputFormat) {
		util.HandleInputError(club.out, line)
		return
	}
	fmt.Fprintf(club.out, "Unexpected error: %s", err)
}
