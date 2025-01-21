package cmd

import (
	"deep-work-tui/models"
	"deep-work-tui/utils"
	"flag"
	"fmt"
	"strings"
)

func ParseArgs() ([]models.Timeblock, error) {
	blocksArg := flag.String("blocks", "", "Time blocks in the format 'START-END TASK'")

	flag.Usage = CustomUsage

	flag.Parse()

	if *blocksArg == "" {
		return nil, nil
	}

	blockArgs := strings.Split(*blocksArg, ";")
	return ParseBlocks(blockArgs)
}

func ParseBlocks(blockArgs []string) ([]models.Timeblock, error) {
	var timeblocks []models.Timeblock

	for _, block := range blockArgs {
		parts := strings.SplitN(block, " ", 2)
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid format: %s. Expected 'START-END TASKNAME;'", block)
		}
		timeRange := parts[0]
		task := parts[1]
		timeParts := strings.Split(timeRange, "-")
		if len(timeParts) != 2 {
			return nil, fmt.Errorf("invalid time range: %s", timeRange)
		}
		start := utils.ParseTime(strings.TrimSpace(timeParts[0]))
		end := utils.ParseTime(strings.TrimSpace(timeParts[1]))
		timeblocks = append(timeblocks, models.Timeblock{
			Task:      task,
			Starttime: start,
			Endtime:   end,
		})
	}
	return timeblocks, nil
}
