package cmd

import (
	"log/slog"
	"os"
)

func assertRunWithErrorMessage(toExecute func() error, exitCode int, msg string, args ...any) {
	err := toExecute()
	if err != nil {
		args = append(args, "error", err)
		slog.Error(msg, args...)
		os.Exit(exitCode)
	}
}

func assertRunWithParamAndErrorMessage[T interface{}](toExecute func(T) error, param T, exitCode int, msg string, args ...any) {
	err := toExecute(param)
	if err != nil {
		args = append(args, "error", err)
		slog.Error(msg, args...)
		os.Exit(exitCode)
	}
}

func assertRunWithResultAndErrorMessage[T interface{}](toExecute func() (T, error), exitCode int, msg string, args ...any) T {
	result, err := toExecute()
	if err != nil {
		args = append(args, "error", err)
		slog.Error(msg, args...)
		os.Exit(exitCode)
	}
	return result
}

func assertRunWithResult[T interface{}](toExecute func() (T, error), exitCode int) T {
	return assertRunWithResultAndErrorMessage(toExecute, exitCode, "%s")
}
