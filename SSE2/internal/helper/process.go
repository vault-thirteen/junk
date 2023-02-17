package helper

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"strings"

	"github.com/rs/zerolog"
	"github.com/vault-thirteen/junk/SSE2/internal/messages"
)

var (
	ErrNoProcess        = errors.New("process has not been created")
	ErrStdOutAlreadySet = errors.New("stdout is already set")
	ErrStdErrAlreadySet = errors.New("stderr is already set")
)

func ExecuteCommandAndGetOutput(
	ctx context.Context,
	logger *zerolog.Logger,
	command string,
	arguments []string,
) (processId *int, outputLines []string, err error) {
	cmd := exec.CommandContext(ctx, command, arguments...)

	if cmd.Stdout != nil {
		return nil, nil, ErrStdOutAlreadySet
	}
	if cmd.Stderr != nil {
		return nil, nil, ErrStdErrAlreadySet
	}

	var buffer bytes.Buffer
	cmd.Stdout = &buffer
	cmd.Stderr = &buffer

	err = cmd.Start()
	if err != nil {
		return nil, nil, err
	}

	if cmd.Process == nil {
		return nil, nil, ErrNoProcess
	}
	processId = NewIntPointer(cmd.Process.Pid)

	defer func() {
		outputLines = strings.Split(buffer.String(), "\n")
	}()

	logger.Debug().Msgf(messages.MsgFProcessHasBeenCreated, cmd.Process.Pid)

	err = cmd.Wait()
	if err != nil {
		return processId, outputLines, err
	}

	logger.Debug().Msgf(messages.MsgFProcessHasFinished, cmd.Process.Pid)

	return processId, outputLines, nil
}
