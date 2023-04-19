package file

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_GenerateFailedPath(t *testing.T) {
	data := []struct {
		name         string
		originalPath string
		want         string
	}{
		{
			"video",
			"video.mkv",
			"video.failed.mkv",
		},
		{
			"log",
			"video.mkv.log",
			"video.failed.mkv.log",
		},
		{
			"vmaf",
			"video.mkv.vmaf.json",
			"video.failed.mkv.vmaf.json",
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			got := GenerateFailedPath(d.originalPath)

			if diff := cmp.Diff(d.want, got); diff != "" {
				t.Error("unexpected failed file name (-want +got):", diff)
			}
		})
	}
}

func Test_appendtoFileName(t *testing.T) {
	data := []struct {
		name   string
		path   string
		suffix string
		want   string
	}{
		{
			"filename suffix",
			"filename.ext",
			"-suffix",
			"filename-suffix.ext",
		},
		{
			"prepend extension",
			"filename.ext2",
			".ext1",
			"filename.ext1.ext2",
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			got := appendToFileName(d.path, d.suffix)

			if diff := cmp.Diff(d.want, got); diff != "" {
				t.Error("unexpected file name (-want +got)", diff)
			}
		})
	}
}

func Test_fullExt(t *testing.T) {
	data := []struct {
		name string
		path string
		want string
	}{
		{
			"standard file",
			"video.mkv",
			".mkv",
		},
		{
			"multi-extension file",
			"video.mkv.log",
			".mkv.log",
		},
		{
			"dotfile",
			".dotfile",
			".dotfile",
		},
		{
			"dotfile with extension",
			".dotfile.ext",
			".dotfile.ext",
		},
		{
			"only considers final path element",
			"path/with/file.extension/video.mkv.log",
			".mkv.log",
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			got := fullExt(d.path)

			if diff := cmp.Diff(d.want, got); diff != "" {
				t.Error("unexpected full extension (-want +got):", diff)
			}
		})
	}
}
