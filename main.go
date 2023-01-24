package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
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

func (f flags) outputVideoPath() string {
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
	return f.outputVideoPath() + ".log"
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

	if err := conversionPossible(args); err != nil {
		fmt.Println("Conversion not possible:", err)
		os.Exit(1)
	}

	if err := conversionNeeded(args); err != nil {
		fmt.Println("Conversion not needed:", err)
		os.Exit(2)
	}

	if err := convertVideo(context.Background(), args); err != nil {
		fmt.Println("Error during conversion:", err)
		os.Exit(3)
	}
}
