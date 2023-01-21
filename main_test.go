package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_OutputPath(t *testing.T) {
	data := []struct {
		name string
		args flags
		want string
	}{
		{
			"mkv extension",
			flags{
				input:      "test.mkv",
				crf:        24,
				preset:     4,
				gop:        24,
				film_grain: 12,
			},
			"test.crf24-p4-g24-fg12.mkv",
		},
		{
			"non-mkv extension",
			flags{
				input:      "test.avi",
				crf:        24,
				preset:     4,
				gop:        24,
				film_grain: 12,
			},
			"test.crf24-p4-g24-fg12.mkv",
		},
		{
			"custom outputDir",
			flags{
				input:      "test.mkv",
				outputDir:  "output",
				crf:        24,
				preset:     4,
				gop:        24,
				film_grain: 12,
			},
			"output/test.crf24-p4-g24-fg12.mkv",
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			got := d.args.outputPath()

			if diff := cmp.Diff(d.want, got); diff != "" {
				t.Error("unexpected output path (-want +got):", diff)
			}
		})
	}
}
