package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PonomarevAlexxander/computer-club-system/internal/app"
	"github.com/PonomarevAlexxander/computer-club-system/internal/config"
	"github.com/PonomarevAlexxander/computer-club-system/internal/scanner"
	"github.com/PonomarevAlexxander/computer-club-system/internal/service"
	"github.com/PonomarevAlexxander/computer-club-system/internal/util"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Error: To run this 'Computer-Club' system you should provide file name.")
		return
	}

	fileName := os.Args[1]

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
	}()

	scanner := scanner.NewScanner(file)
	out := os.Stdout

	line, err := scanner.ReadLine()
	if err != nil {
		util.HandleInputError(out, line)
		return
	}

	availableTables, err := strconv.ParseUint(line, 10, 0)
	if err != nil || availableTables < 1 {
		util.HandleInputError(out, line)
		return
	}

	line, err = scanner.ReadLine()
	if err != nil {
		util.HandleInputError(out, line)
		return
	}

	timeStrings := strings.Split(line, " ")
	if len(timeStrings) != 2 {
		util.HandleInputError(out, line)
		return
	}

	openingHours, err := time.Parse(config.TimeFormat, timeStrings[0])
	if err != nil {
		util.HandleInputError(out, line)
		return
	}

	closingHours, err := time.Parse(config.TimeFormat, timeStrings[1])
	if err != nil {
		util.HandleInputError(out, line)
		return
	}

	line, err = scanner.ReadLine()
	if err != nil {
		util.HandleInputError(out, line)
		return
	}

	hourPrice, err := strconv.ParseUint(line, 10, 0)
	if err != nil || hourPrice < 1 {
		util.HandleInputError(out, line)
		return
	}

	clubService := service.NewComputerClubService(uint(availableTables), openingHours, closingHours, uint(hourPrice))

	clubApp := app.NewComputerClub(scanner, out, clubService)

	clubApp.Start()
}
