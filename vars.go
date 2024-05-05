package main

import "github.com/heroiclabs/nakama-common/runtime"

const (
	OK               = 0
	CANCELED         = 1
	INVALID_ARGUMENT = 2
)

var (
	errBadInput     = runtime.NewError("input contained invalid data", INVALID_ARGUMENT)
	errUserNotFound = runtime.NewError("user not found", CANCELED)
	errFileNotFound = runtime.NewError("file not found", INVALID_ARGUMENT)
)
