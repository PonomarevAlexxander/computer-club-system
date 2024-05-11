package command

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PonomarevAlexxander/computer-club-system/internal/config"
	"github.com/PonomarevAlexxander/computer-club-system/internal/errlib"
)

type Command struct {
	Time   time.Time
	Id     uint
	Params []string
}

func NewCommand(time time.Time, id uint, params ...string) *Command {
	return &Command{
		Time:   time,
		Id:     id,
		Params: params,
	}
}

func CommandFromString(line string) (*Command, error) {
	params := strings.Split(line, " ")

	if len(params) < 2 {
		return nil, errlib.ErrInputFormat
	}

	t, err := time.Parse(config.TimeFormat, params[0])
	if err != nil {
		return nil, fmt.Errorf("%w, %w", errlib.ErrInputFormat, err)
	}

	commandId, err := strconv.ParseUint(params[1], 10, 0)
	if err != nil {
		return nil, errlib.ErrInputFormat
	}

	return &Command{
		Time:   t,
		Id:     uint(commandId),
		Params: params[2:],
	}, nil
}

func (c *Command) String() string {
	return fmt.Sprintf("%s %d %s", c.Time.Format(config.TimeFormat), c.Id, strings.Join(c.Params, " "))
}
