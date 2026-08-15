package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"iasi-cli/internal/install"
	"iasi-cli/internal/source"
	"iasi-cli/internal/status"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return errors.New("usage: iasi install --workspace | iasi status")
	}

	switch args[0] {
	case "install":
		installFlags := flag.NewFlagSet("install", flag.ContinueOnError)
		installFlags.SetOutput(os.Stderr)
		workspace := installFlags.Bool("workspace", false, "install IASI in the current workspace")
		if err := installFlags.Parse(args[1:]); err != nil {
			return err
		}
		if !*workspace || installFlags.NArg() != 0 {
			return errors.New("usage: iasi install --workspace")
		}

		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("get current directory: %w", err)
		}
		root, err := source.New().Find()
		if err != nil {
			return err
		}
		_, err = install.Run(cwd, root)
		return err

	case "status":
		if len(args) != 1 {
			return errors.New("usage: iasi status")
		}
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("get current directory: %w", err)
		}
		result, err := status.Find(cwd)
		if err != nil {
			return err
		}
		fmt.Print(status.Format(result))
		return nil
	}

	return fmt.Errorf("unknown command %q", args[0])
}
