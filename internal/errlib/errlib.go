package errlib

import (
	"errors"
)

var (
	ErrInternal            = errors.New("unknown internal error happen")
	ErrUnknownCommand      = errors.New("unknown command")
	ErrInputFormat         = errors.New("invalid input format error")
	ErrClientAlreadyInClub = errors.New("YouShallNotPass")
	ErrClubClosed          = errors.New("NotOpenYet")
	ErrTableAlreadyTaken   = errors.New("PlaceIsBusy")
	ErrUnknownClient       = errors.New("ClientUnknown")
	ErrWaitingTimeExceeded = errors.New("ICanWaitNoLonger!")
)
