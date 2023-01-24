package main

import (
	"context"
	"errors"
	"fmt"
	"os"
)

func vmafPossible(f flags) error {
	if f.input == "" {
		return fmt.Errorf("no input provided -- use -i <path>")
	}
	if _, err := os.Stat(f.input); err != nil {
		// There's a race condition here since we're checking for the existence of the file
		// before using it. However, since we're just trying to be helpful, and the error
		// will be caught by VMAF later if conversion fails, the risk is acceptable.
		// Additionally, we wouldn't normally anticipate video files suddenly being created
		// just before they're analysed, although I'm sure it could happen in someone's use case eventually.
		// See: https://xkcd.com/1172/
		return fmt.Errorf("error finding input file: %w", err)
	}
	return nil
}

func vmafNeeded(f flags) error {
	path := f.outputVmafPath()
	if _, err := os.Stat(path); err != nil {
		// No error, file exists
		return fmt.Errorf("file already exists: %v", path)
	} else if errors.Is(err, os.ErrNotExist) {
		// File doesn't exist, we're good to go
		return nil
	} else {
		// Another file error occurred
		return err
	}
}

func performVMAFAnalysis(f flags) error {
	return nil
}

func executeVMAF(ctx context.Context, f flags) error {
	return nil
}
