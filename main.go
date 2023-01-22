package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type flags struct {
	input      string
	outputDir  string
	crf        int
	preset     int
	gop        int
	film_grain int
}

func (f flags) outputPath() string {
	outputDir, inputFile := filepath.Split(f.input)

	// Figure out directory
	if f.outputDir != "" {
		outputDir = f.outputDir
	}

	// Figure out file name
	inputExt := filepath.Ext(inputFile)
	fileName := strings.Join([]string{
		inputFile[:len(inputFile)-len(inputExt)],
		fmt.Sprintf(".crf%v-p%v-g%v-fg%v", f.crf, f.preset, f.gop, f.film_grain),
		".mkv",
	}, "")
	return filepath.Join(
		outputDir,
		fileName,
	)
}

func (f flags) outputLogPath() string {
	return f.outputPath() + ".log"
}

func main() {
	var args flags
	flag.StringVar(&args.input, "i", "", "Input")
	flag.StringVar(&args.outputDir, "o", "", "Output directory")
	flag.IntVar(&args.crf, "crf", 30, "CRF")
	flag.IntVar(&args.preset, "p", 5, "Preset")
	flag.IntVar(&args.film_grain, "fg", 8, "Film Grain")
	flag.IntVar(&args.gop, "g", 240, "Number of frames in Group Of Pictures")
	flag.Parse()

	if err := conversion_possible(args); err != nil {
		fmt.Println("Conversion not possible:", err)
		os.Exit(1)
	}

	outputPath := args.outputPath()
	outputLogPath := args.outputLogPath()

	if _, err := os.Stat(outputLogPath); err == nil {
		fmt.Printf("%v already exists; skipping\n", outputLogPath)
		os.Exit(2)
	} else if errors.Is(err, os.ErrNotExist) {
		err = nil // Log file doesn't exist, continue
	} else if err != nil {
		log.Fatalln("Failure to create log file:", err)
	}

	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		log.Fatalln("Unable to create pipe:", err)
	}

	ffmpegCmd := exec.Command("time", "-v", "ffmpeg", "-y",
		"-i", args.input,
		"-map", "0:V", "-c:v", "libsvtav1", "-pix_fmt", "yuv420p10le",
		"-g", strconv.Itoa(args.gop),
		"-crf", strconv.Itoa(args.crf),
		"-svtav1-params", fmt.Sprintf("tune=0:film-grain=%v", args.film_grain),
		"-preset", strconv.Itoa(args.preset),
		outputPath,
	)
	ffmpegCmd.Stdout = writePipe
	ffmpegCmd.Stderr = writePipe
	err = ffmpegCmd.Start()
	if err != nil {
		log.Fatalln("Failed to start:", err)
	}
	defer ffmpegCmd.Wait()
	writePipe.Close() // Can now be closed as cmd has inherited the file descriptor

	// Push ffmpeg's output to both the terminal and the output file using tee,
	// both providing immediate feedback and a log for later
	teeCmd := exec.Command("tee", outputLogPath)
	teeCmd.Stdin = readPipe
	teeCmd.Stdout = os.Stdout
	teeCmd.Stderr = os.Stderr
	teeCmd.Run()
}
