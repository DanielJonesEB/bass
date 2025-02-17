package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/vito/bass"
)

func run(ctx context.Context, env *bass.Env, filePath string) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = bass.EvalReader(ctx, env, file)
	if err != nil {
		return err
	}

	return nil
}
