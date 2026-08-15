package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"iasi-cli/internal/adapters/copilot"
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
		return errors.New("usage: iasi version | iasi install --workspace | iasi status")
	}

	switch args[0] {
	case "adapt":
		if len(args) != 2 || args[1] != "copilot" {
			return errors.New("usage: iasi adapt copilot")
		}
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("get current directory: %w", err)
		}
		output, err := copilot.Run(cwd)
		if err != nil {
			return err
		}
		fmt.Print(output)
		return nil

	case "version":
		if len(args) != 1 {
			return errors.New("usage: iasi version")
		}
		version, err := source.Version()
		if err != nil {
			return err
		}
		fmt.Printf("IASI %s\n", version)
		return nil

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
		version, err := source.Version()
		if err != nil {
			return err
		}
		_, err = install.Run(cwd, source.Methodology(), version)
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
		version, err := source.Version()
		if err != nil {
			return err
		}
		fmt.Print(status.Format(result, version))
		return nil
	}

	return fmt.Errorf("unknown command %q", args[0])
}
