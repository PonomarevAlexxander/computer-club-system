package command_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/PonomarevAlexxander/computer-club-system/internal/command"
)

const (
	validCommand1   = "08:48 1 client1"
	validCommand2   = "09:54 2 client1 1"
	validCommand3   = "10:59 2"
	invalidCommand1 = "1234 1 client1"
	invalidCommand2 = "10:59 -1 client1"
	invalidCommand3 = "10:59 str client1"
	invalidCommand4 = "10:59:12 123 client1"
)

var (
	command1 = command.Command{Time: time.Date(0, time.January, 1, 8, 48, 0, 0, time.Now().UTC().Location()), Id: 1, Params: []string{"client1"}}
	command2 = command.Command{Time: time.Date(0, time.January, 1, 9, 54, 0, 0, time.Now().UTC().Location()), Id: 2, Params: []string{"client1", "1"}}
	command3 = command.Command{Time: time.Date(0, time.January, 1, 10, 59, 0, 0, time.Now().UTC().Location()), Id: 2, Params: []string{}}
)

func TestCommandFromString(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name    string
		args    args
		want    *command.Command
		wantErr bool
	}{
		{
			"valid command1",
			args{validCommand1},
			&command1,
			false,
		},
		{
			"valid command2",
			args{validCommand2},
			&command2,
			false,
		},
		{
			"valid command3",
			args{validCommand3},
			&command3,
			false,
		},
		{
			"command with invalid time",
			args{invalidCommand1},
			nil,
			true,
		},
		{
			"command with invalid id",
			args{invalidCommand2},
			nil,
			true,
		},
		{
			"command with invalid id",
			args{invalidCommand3},
			nil,
			true,
		},
		{
			"command with invalid time",
			args{invalidCommand4},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := command.CommandFromString(tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("CommandFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CommandFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}
