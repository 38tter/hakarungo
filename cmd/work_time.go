package cmd

import "time"

type workTime struct {
	dirPaths []string
	workTime map[string]time.Duration
}

func TotalWorkTime(workTimes map[string]time.Duration) time.Duration {
	var total time.Duration
	for _, t := range workTimes {
		total += t
	}
	return total
}
