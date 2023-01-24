package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
)

func conversionPossible(f flags) error {
	if f.input == "" {
		return fmt.Errorf("no input provided -- use -i <path>")
	}
	if _, err := os.Stat(f.input); err != nil {
		// There's a race condition here since we're checking for the existence of the file
		// before using it. However, since we're just trying to be helpful, and the error
		// will be caught by FFmpeg later if conversion fails, the risk is acceptable.
		// Additionally, we wouldn't normally anticipate video files suddenly being created
		// just before they're converted, although I'm sure it could happen in someone's use case eventually.
		// See: https://xkcd.com/1172/
		return fmt.Errorf("error finding input file: %w", err)
	}
	return nil
}

func conversionNeeded(f flags) error {
	outputLogPath := f.outputLogPath()

	if _, err := os.Stat(outputLogPath); errors.Is(err, os.ErrNotExist) { // Log file doesn't exist, conversion needed
		return nil
	} else if err != nil { // A different file error occurred
		return err
	} else { // No error, file exists
		return fmt.Errorf("log already exists: %v", outputLogPath)
	}
}

func convertVideo(ctx context.Context, f flags) error {
	// Allow conversion to be stopped via interrupt, e.g. Ctrl-C on CLI
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	var err error
	done := make(chan struct{}, 1)
	go func() {
		err = runVideoConversion(ctx, f)
		close(done)
	}()

	select {
	case <-done: // Video conversion completed, potentially with err
		break
	case <-ctx.Done(): // Interrupt signal received
		err = ctx.Err()
		break
	}

	stop() // Restore original signal processing
	if err != nil {
		cleanupConversion(f)
	}
	return err
}

func runVideoConversion(ctx context.Context, f flags) error {
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		return fmt.Errorf("unable to create pipe: %w", err)
	}

	// Create output directory, otherwise tee will fail
	if err := os.MkdirAll(f.outputDir, 0750); err != nil {
		return fmt.Errorf("failed to create output directory %v: %w", f.outputDir, err)
	}

	ffmpegCmd := exec.CommandContext(ctx, "time", "-v", "ffmpeg", "-y",
		"-i", f.input,
		"-map", "0:V", "-c:v", "libsvtav1", "-pix_fmt", "yuv420p10le",
		"-g", strconv.Itoa(f.gop),
		"-crf", strconv.Itoa(f.crf),
		"-svtav1-params", fmt.Sprintf("tune=0:film-grain=%v", f.film_grain),
		"-preset", strconv.Itoa(f.preset),
		f.outputVideoPath(),
	)
	ffmpegCmd.Stdout = writePipe
	ffmpegCmd.Stderr = writePipe
	if err := ffmpegCmd.Start(); err != nil {
		return fmt.Errorf("failed to start ffmpeg: %w", err)
	}
	writePipe.Close() // Can now be closed as cmd has inherited the file descriptor

	// Push ffmpeg's output to both the terminal and the output file using tee,
	// both providing immediate feedback and a log for later
	teeCmd := exec.CommandContext(ctx, "tee", f.outputLogPath())
	teeCmd.Stdin = readPipe
	teeCmd.Stdout = os.Stdout
	teeCmd.Stderr = os.Stderr

	if err := teeCmd.Run(); err != nil {
		return err
	}

	if err := ffmpegCmd.Wait(); err != nil {
		return err
	}

	return nil
}

func cleanupConversion(f flags) error {
	log.Println("cleaning up conversion")
	// Move video file
	if err := os.Rename(f.outputVideoPath(), generateFailedPath(f.outputVideoPath())); err != nil {
		log.Println("Error during cleanup, continuing:", err)
	}
	// Move log file
	if err := os.Rename(f.outputLogPath(), generateFailedPath(f.outputLogPath())); err != nil {
		log.Println("Error during cleanup, continuing:", err)
	}

	return nil
}

// generateFailedPath inserts an indicator in to originalPath to indicate failure.
// The intention is that moving the file located at originalPath to the path returned by generateFailedPath would indicate
// to a user that processing for that file has failed without needing to inspect the file or command output.
func generateFailedPath(originalPath string) string {
	return appendToFileName(originalPath, ".failed")
}

// appendtoFilename appends a string suffix to the filename, before any extension, defined by a dot.
// For example, appending suffix "-suffix" to path "filename.ext" will result in "filename-suffix.ext".
func appendToFileName(path string, suffix string) string {
	originalExt := fullExt(path)
	insertionIndex := len(path) - len(originalExt)
	return strings.Join([]string{
		path[:insertionIndex],
		suffix,
		path[insertionIndex:],
	}, "")
}

// fullExt returns the full file name extension used by the path.
// The full extension is the suffix beginning after the first dot in the final element of the path.
// An empty string is returned if there is no dot.
func fullExt(path string) string {
	extIndex := -1
	for i := len(path) - 1; i >= 0 && !os.IsPathSeparator(path[i]); i-- {
		if path[i] == '.' {
			extIndex = i
		}
	}
	if extIndex != -1 {
		return path[extIndex:]
	} else {
		return ""
	}
}
