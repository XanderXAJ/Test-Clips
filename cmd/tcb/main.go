package main

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

const CONFIG_FILE = "tcb.toml"

type ConfigJob struct {
	Name       string `toml:"name"`
	CRF        []int  `toml:"crf"`
	Film_grain []int  `toml:"film-grain"`
	GOP        []int  `toml:"gop"`
	Preset     []int  `toml:"preset"`
}

type Config struct {
	Jobs []ConfigJob `toml:"jobs"`
}

func main() {
	file, err := os.Open(CONFIG_FILE)
	if err != nil {
		panic(err)
	}

	var config Config
	err = toml.NewDecoder(file).Decode(&config)
	if err != nil {
		panic(err)
	}

	err = file.Close()
	if err != nil {
		panic(err)
	}

	for _, job := range config.Jobs {
		for _, crf := range job.CRF {
			for _, film_grain := range job.Film_grain {
				for _, gop := range job.GOP {
					for _, preset := range job.Preset {
						fmt.Println(job.Name, crf, film_grain, gop, preset)
					}
				}
			}
		}
	}
}
