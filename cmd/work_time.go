package cmd

import "time"

type workTime struct {
	dirPaths []string
	workTime map[string]time.Duration
}
