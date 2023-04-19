package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/xanderxaj/test-clips/internal/job"
)

func Test_OutputPath(t *testing.T) {
	data := []struct {
		name string
		args job.Flags
		want string
	}{
		{
			"mkv extension",
			job.Flags{
				Input:      "test.mkv",
				CRF:        24,
				Preset:     4,
				GOP:        24,
				Film_grain: 12,
			},
			"test.crf24-p4-g24-fg12.mkv",
		},
		{
			"non-mkv extension",
			job.Flags{
				Input:      "test.avi",
				CRF:        24,
				Preset:     4,
				GOP:        24,
				Film_grain: 12,
			},
			"test.crf24-p4-g24-fg12.mkv",
		},
		{
			"custom outputDir",
			job.Flags{
				Input:      "test.mkv",
				OutputDir:  "output",
				CRF:        24,
				Preset:     4,
				GOP:        24,
				Film_grain: 12,
			},
			"output/test.crf24-p4-g24-fg12.mkv",
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			got := d.args.OutputVideoPath()

			if diff := cmp.Diff(d.want, got); diff != "" {
				t.Error("unexpected output path (-want +got):", diff)
			}
		})
	}
}
