package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/xanderxaj/test-clips/internal/conversion"
	"github.com/xanderxaj/test-clips/internal/job"
	"github.com/xanderxaj/test-clips/internal/vmaf"
)

func main() {
	var args job.Flags
	flag.StringVar(&args.Input, "i", "", "Input")
	flag.StringVar(&args.OutputDir, "o", "", "Output directory")
	flag.IntVar(&args.CRF, "crf", 30, "CRF")
	flag.IntVar(&args.Preset, "p", 5, "Preset")
	flag.IntVar(&args.Film_grain, "fg", 8, "Film Grain")
	flag.IntVar(&args.GOP, "g", 240, "Number of frames in Group Of Pictures")
	flag.Parse()

	if err := conversion.ConversionPossible(args); err != nil {
		fmt.Println("Conversion not possible:", err)
		os.Exit(1)
	}

	if err := conversion.ConversionNeeded(args); err != nil {
		fmt.Println("Conversion not needed:", err)
	} else {
		if err := conversion.ConvertVideo(context.Background(), args); err != nil {
			fmt.Println("Error during conversion:", err)
			os.Exit(3)
		}
	}

	if err := vmaf.VMAFPossible(args); err != nil {
		fmt.Println("VMAF analysis not possible:", err)
		os.Exit(1)
	}

	if err := vmaf.VMAFNeeded(args); err != nil {
		fmt.Println("VMAF analysis not needed:", err)
	} else {
		if err := vmaf.PerformVMAFAnalysis(context.Background(), args); err != nil {
			fmt.Println("Error during VMAF analysis:", err)
			os.Exit(3)
		}
	}
}
