package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"os/signal"

	"github.com/mattn/go-colorable"
	"github.com/spf13/cobra"
	"github.com/vito/bass"
)

var Stderr = colorable.NewColorableStderr()

//go:embed txt/help.txt
var helpText string

var rootCmd = &cobra.Command{
	Use:           "bass",
	Short:         "run bass code, or start a repl",
	Long:          helpText,
	Version:       bass.Version,
	SilenceUsage:  true,
	SilenceErrors: true,
	// Args:          cobra.MaximumNArgs(1),
	RunE: root,
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		fmt.Fprintln(Stderr, err)
		os.Exit(1)
	}
}

func root(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return repl(bass.New())
	}

	env := bass.New()

	var argVals []bass.Value
	for _, arg := range args[1:] {
		argVals = append(argVals, bass.String(arg))
	}

	env.Set("*args*", bass.NewList(argVals...))

	return run(env, args[0])
}
