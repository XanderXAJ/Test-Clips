package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_generateFailedPath(t *testing.T) {
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
			got := generateFailedPath(d.originalPath)

			if diff := cmp.Diff(d.want, got); diff != "" {
				t.Error("unexpected failed file name (-want +got):", diff)
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
