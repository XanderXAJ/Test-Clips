package main

import (
	"errors"
	"fmt"
	"os"
)

func conversion_possible(f flags) error {
	if f.input == "" {
		return fmt.Errorf("no input provided -- use -i <path>")
	}
	if _, err := os.Stat(f.input); err != nil {
		// There's a race condition here since we're checking for the existence of the file
		// before using it. However, since we're just trying to be helpful, and the error
		// will be caught by FFmpeg later if conversion fails, the risk is acceptable.
		// Additionally, we wouldn't normally anticipate video files suddenly being created
		// just before they're converted, although I'm sure it could happen in someone's use case eventually.
		return fmt.Errorf("error finding input file: %w", err)
	}
	return nil
}

func conversion_needed(f flags) error {
	outputLogPath := f.outputLogPath()

	if _, err := os.Stat(outputLogPath); errors.Is(err, os.ErrNotExist) { // Log file doesn't exist, conversion needed
		return nil
	} else if err != nil { // A different file error occurred
		return err
	} else { // No error, file exists
		return fmt.Errorf("log already exists: %v", outputLogPath)
	}
}

func convert_video(f flags) {

}
