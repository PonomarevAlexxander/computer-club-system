package handler

import (
	"time"

	"github.com/PonomarevAlexxander/computer-club-system/internal/queue"
	"github.com/PonomarevAlexxander/computer-club-system/internal/stats"
)

type ClubStateDto struct {
	totalTables   uint
	hourPrice     uint
	clientToTable map[string]uint
	// tablesToTime - relation table to time when it was taken
	tableToTime     map[uint]time.Time
	clientsQueue    *queue.Queue
	openingTime     time.Time
	closingTime     time.Time
	tableStatistics []stats.TableStats
}

func NewClubStateDto(
	totalTables uint,
	hourPrice uint,
	clientToTable map[string]uint,
	tableToTime map[uint]time.Time,
	clientsQueue *queue.Queue,
	openingTime time.Time,
	closingTime time.Time,
	tableStatistics []stats.TableStats,
) *ClubStateDto {
	return &ClubStateDto{
		totalTables:     totalTables,
		hourPrice:       hourPrice,
		clientToTable:   clientToTable,
		tableToTime:     tableToTime,
		clientsQueue:    clientsQueue,
		openingTime:     openingTime,
		closingTime:     closingTime,
		tableStatistics: tableStatistics,
	}
}
