package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
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

func convertVideo(f flags) error {
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		return fmt.Errorf("unable to create pipe: %w", err)
	}

	// Create output directory, otherwise tee will fail
	if err := os.MkdirAll(f.outputDir, 0750); err != nil {
		return fmt.Errorf("failed to create output directory %v: %w", f.outputDir, err)
	}

	ffmpegCmd := exec.Command("time", "-v", "ffmpeg", "-y",
		"-i", f.input,
		"-map", "0:V", "-c:v", "libsvtav1", "-pix_fmt", "yuv420p10le",
		"-g", strconv.Itoa(f.gop),
		"-crf", strconv.Itoa(f.crf),
		"-svtav1-params", fmt.Sprintf("tune=0:film-grain=%v", f.film_grain),
		"-preset", strconv.Itoa(f.preset),
		f.outputPath(),
	)
	ffmpegCmd.Stdout = writePipe
	ffmpegCmd.Stderr = writePipe
	if err := ffmpegCmd.Start(); err != nil {
		return fmt.Errorf("failed to start ffmpeg: %w", err)
	}
	writePipe.Close() // Can now be closed as cmd has inherited the file descriptor

	// Push ffmpeg's output to both the terminal and the output file using tee,
	// both providing immediate feedback and a log for later
	teeCmd := exec.Command("tee", f.outputLogPath())
	teeCmd.Stdin = readPipe
	teeCmd.Stdout = os.Stdout
	teeCmd.Stderr = os.Stderr

	if err := teeCmd.Run(); err != nil {
		log.Println("tee command failed")
		conversionCleanup(f)
		return err
	}

	if err := ffmpegCmd.Wait(); err != nil {
		log.Println("ffmpeg command failed")
		conversionCleanup(f)
		return err
	}

	return nil
}

func conversionCleanup(f flags) error {
	log.Println("cleaning up conversion")
	// TODO move video file
	if err := os.Rename(f.outputPath(), generateFailedPath(f.outputPath())); err != nil {
		log.Println("Error during cleanup, continuing:", err)
	}
	// TODO move log file
	if err := os.Rename(f.outputLogPath(), generateFailedPath(f.outputLogPath())); err != nil {
		log.Println("Error during cleanup, continuing:", err)
	}

	return nil
}

func generateFailedPath(originalPath string) string {
	originalExt := fullExt(originalPath)
	insertionIndex := len(originalPath) - len(originalExt)
	return strings.Join([]string{
		originalPath[:insertionIndex],
		".failed",
		originalPath[insertionIndex:],
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
