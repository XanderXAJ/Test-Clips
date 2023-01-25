package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/hashicorp/go-multierror"
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
	path := f.outputVMAFPath()
	if _, err := os.Stat(path); err == nil {
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

func performVMAFAnalysis(ctx context.Context, f flags) error {
	err := executeVMAF(ctx, f)
	if err != nil {
		cleanupFailedVMAFAnalysis(f)
	}
	return err
}

func executeVMAF(ctx context.Context, f flags) error {
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		return fmt.Errorf("unable to create pipe: %w", err)
	}

	// Create output directory, otherwise tee will fail
	if err := os.MkdirAll(f.outputDir, 0750); err != nil {
		return fmt.Errorf("failed to create output directory %v: %w", f.outputDir, err)
	}

	vmafComplexFilter := fmt.Sprintf(`
	[0:v]setpts=PTS-STARTPTS[reference];
	[1:v]setpts=PTS-STARTPTS[distorted];
	[distorted][reference]libvmaf=log_fmt=json:log_path=%v:n_threads=%v:feature='name=psnr|name=float_ssim|name=float_ms_ssim|name=float_ansnr'
	`, f.outputVMAFPath(), runtime.NumCPU())

	vmafCmd := exec.Command("time", "-v", "ffmpeg",
		"-i", f.input,
		"-i", f.outputVideoPath(),
		"-lavfi", vmafComplexFilter,
		"-f", "null", "-",
	)
	vmafCmd.Stdout = writePipe
	vmafCmd.Stderr = writePipe
	if err := vmafCmd.Start(); err != nil {
		return fmt.Errorf("failed to start VMAF analysis: %w", err)
	}
	writePipe.Close() // Can now be closed as cmd has inherited the file descriptor. If we don't do this, go won't close after all tasks are complete.

	// Push ffmpeg's output to both the terminal and the output file using tee,
	// both providing immediate feedback and a log for later.
	// Explicitly do not tie tee to the context, since it will terminate when its pipes are closed.
	teeCmd := exec.Command("tee", f.outputVMAFLogPath())
	teeCmd.Stdin = readPipe
	teeCmd.Stdout = os.Stdout
	teeCmd.Stderr = os.Stderr
	if err := teeCmd.Start(); err != nil {
		return fmt.Errorf("failed to start tee: %w", err)
	}

	var result error
	if err := interruptibleWait(vmafCmd, os.Interrupt); err != nil {
		result = multierror.Append(result, fmt.Errorf("VMAF command failed: %w", err))
	}
	log.Println("VMAF ended")
	if err := teeCmd.Wait(); err != nil {
		result = multierror.Append(result, fmt.Errorf("tee command failed: %w", err))
	}
	log.Println("tee ended")
	return result
}

func cleanupFailedVMAFAnalysis(f flags) {
	log.Println("Cleaning up failed VMAF analysis")
	// Move failed VMAF analysis
	if err := os.Rename(f.outputVMAFPath(), generateFailedPath(f.outputVMAFPath())); err != nil {
		log.Println("Error during cleanup, continuing:", err)
	}
	// Move failed VMAF analysis log
	if err := os.Rename(f.outputVMAFLogPath(), generateFailedPath(f.outputVMAFLogPath())); err != nil {
		log.Println("Error during cleanup, continuing:", err)
	}
}
