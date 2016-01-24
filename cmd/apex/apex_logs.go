package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/apex/log"
)

type LogsCmdLocalValues struct {
	Filter   string
	Follow   bool
	Duration string
	Start    string
	End      string

	name string
}

const logsCmdExample = `  Print logs for a function
  $ apex logs <name>

  Print logs for a function with a specified duration, e.g. 5 minutes
  $ apex logs <name> 5m

  Print logs for a function for a customized time range
  $ apex logs <name> --start "18/01/2016 10:00" --end "19/01/2016 22:00"`

var logsCmd = &cobra.Command{
	Use:     "logs <name> [<duration>] [--start <startDate>] [--end <endDate>]",
	Short:   "Output logs with optional filter pattern",
	Example: logsCmdExample,
	PreRun:  logsCmdPreRun,
	Run:     logsCmdRun,
}

var logsCmdLocalValues = LogsCmdLocalValues{}

func init() {
	lv := &logsCmdLocalValues
	f := logsCmd.Flags()

	f.StringVarP(&lv.Filter, "filter", "F", "", "Filter logs with pattern")
	f.BoolVarP(&lv.Follow, "follow", "f", false, "Tail logs")
	f.StringVar(&lv.Start, "start", "", "Start Date")
	f.StringVar(&lv.End, "end", "", "End Date")
}

func logsCmdPreRun(c *cobra.Command, args []string) {
	lv := &logsCmdLocalValues

	if len(args) < 1 {
		log.Fatal("Missing name argument")
	}
	lv.name = args[0]

	if len(args) >= 2 {
		lv.Duration = args[1]
	}
}

func logsCmdRun(c *cobra.Command, args []string) {
	lv := &logsCmdLocalValues

	l, err := pv.project.Logs(pv.session, lv.name, lv.Filter, lv.Duration, lv.Start, lv.End)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	if lv.Follow {
		for event := range l.Tail() {
			fmt.Printf("%s", *event.Message)
		}

		if err := l.Err(); err != nil {
			log.Fatalf("error: %s", err)
		}
	}

	events, err := l.Fetch()
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	for _, event := range events {
		fmt.Printf("%s", *event.Message)
	}

}
