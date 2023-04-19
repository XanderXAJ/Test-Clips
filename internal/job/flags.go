package job

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Flags struct {
	Input      string
	OutputDir  string
	CRF        int
	Preset     int
	GOP        int
	Film_grain int
}

func (f Flags) OutputVideoPath() string {
	outputDir, inputFile := filepath.Split(f.Input)

	// Figure out directory
	if f.OutputDir != "" {
		outputDir = f.OutputDir
	}

	// Figure out file name
	inputExt := filepath.Ext(inputFile)
	fileName := strings.Join([]string{
		inputFile[:len(inputFile)-len(inputExt)],
		fmt.Sprintf(".crf%v-p%v-g%v-fg%v", f.CRF, f.Preset, f.GOP, f.Film_grain),
		".mkv",
	}, "")
	return filepath.Join(
		outputDir,
		fileName,
	)
}

func (f Flags) OutputVideoLogPath() string {
	return f.OutputVideoPath() + ".log"
}

func (f Flags) OutputVideoProcessStatsPath() string {
	return f.OutputVideoPath() + ".stats.json"
}

func (f Flags) OutputVMAFPath() string {
	return f.OutputVideoPath() + ".vmaf.json"
}

func (f Flags) OutputVMAFLogPath() string {
	return f.OutputVideoPath() + ".vmaf.log"
}
