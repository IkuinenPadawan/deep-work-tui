package cmd

import (
	"flag"
	"fmt"
)

func CustomUsage() {
	fmt.Fprintf(flag.CommandLine.Output(), `
Usage:
  [options]

Options:
  --blocks "START-END TASK[;START-END TASK;...]" 
      Provide time blocks separated by semicolons (;) where each block has:
      - START and END in the format HH:MM
      - TASK is the description of the activity

Example:
  --blocks "07:00-09:00 Deep Work;09:00-09:30 Emails"

Help:
  Use this program to schedule and visualize timeblocks.
`)
}
