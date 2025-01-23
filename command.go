package ibk

import (
	"context"
	"errors"
)

var ErrConfigure = errors.New("error while configuring")

type Command interface {
	Configure(ctx context.Context, exec Executor) error
	Push(ctx context.Context, pusher Pusher) error
	Build() string
}

func ApplyCommand(ctx context.Context, c Command, t Transport) error {
	err := c.Configure(ctx, t)
	if err != nil {
		return err
	}

	err = c.Push(ctx, t)
	if err != nil {
		return err
	}

	err = t.Execute(ctx, c)
	if err != nil {
		return err
	}

	return nil
}
