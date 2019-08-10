package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/graham/rrule"
	"github.com/graham/rrule/eventexpander"
	"github.com/spf13/cobra"
	google_calendar "google.golang.org/api/calendar/v3"
)

var debugMode bool
var outputFilename string
var inputFilename string

var startDateString string
var endDateString string

var maxIterations int
var passThrough bool
var exitOnFail bool

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().BoolVar(&debugMode, "debug", false, "enable debug mode (prints to stderr)")

	RootCmd.PersistentFlags().StringVar(&outputFilename, "output", "", "output filename for results (default stdout)")
	RootCmd.PersistentFlags().StringVar(&inputFilename, "input", "", "input filename for results (default stdin)")
	RootCmd.PersistentFlags().StringVar(&startDateString, "start", "20100101", "start date format YYYYMMDD")
	RootCmd.PersistentFlags().StringVar(&endDateString, "end", "20200101", "end date format YYYYMMDD")
	RootCmd.PersistentFlags().IntVar(&maxIterations, "max", 0, "Maximum number of interations created, default unlimited")

	RootCmd.PersistentFlags().BoolVar(&passThrough, "pass", false, "re-emit events without a recurrence value (default false).")
	RootCmd.PersistentFlags().BoolVar(&exitOnFail, "fail", false, "exit(-1) on any recurrence related error")
}

func initConfig() {
	// This function is called after parsing and before the run.

	if debugMode == false {
		log.SetOutput(ioutil.Discard)
	}

	log.Printf("Starting run %s\n", time.Now())
}

var RootCmd = &cobra.Command{
	Use:   "evexp",
	Short: "expand google calendar events into their recurrent sub events.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var input_reader io.Reader
		if len(inputFilename) > 0 {
			if _, err := os.Stat(inputFilename); os.IsNotExist(err) {
				fmt.Errorf("input file %s does not exist.", inputFilename)
			}

			fh, err := os.Open(inputFilename)
			if err != nil {
				panic(err)
			}
			defer fh.Close()
			input_reader = fh
		} else {
			log.Printf("Reading from stdin...")
			input_reader = os.Stdin
		}

		var output_writer io.Writer
		var output_filename string = outputFilename
		file_mode := os.O_CREATE | os.O_WRONLY

		if len(output_filename) > 0 {
			outputFileHandle, err := os.OpenFile(output_filename, file_mode, 0644)
			if err != nil {
				panic(err)
			}
			output_writer = outputFileHandle
			defer outputFileHandle.Close()
		} else {
			output_writer = os.Stdout
		}

		var enc *json.Encoder = json.NewEncoder(output_writer)
		var dec *json.Decoder = json.NewDecoder(input_reader)

		var allowAfter time.Time
		var allowBefore time.Time
		var err error

		if len(startDateString) > 0 {
			allowAfter, err = time.Parse("20060102", startDateString)
			if err != nil {
				panic(err)
			}
		}

		if len(endDateString) > 0 {
			allowBefore, err = time.Parse("20060102", endDateString)
			if err != nil {
				panic(err)
			}
		}

		for {
			var readErr error
			var event *google_calendar.Event

			readErr = dec.Decode(&event)

			if readErr != nil && readErr != io.EOF {
				panic(readErr)
			}

			if readErr == io.EOF {
				return
			}

			if event.Start == nil {
				break
			}

			if len(event.Recurrence) == 0 {
				if passThrough == true {
					enc.Encode(event)
				}
				continue
			}

			rule, parseErr := rrule.Parse(strings.Join(event.Recurrence, "\n"))

			if parseErr != nil && exitOnFail {
				log.Panic(parseErr)
			}

			timeZone, tzErr := time.LoadLocation(event.Start.TimeZone)

			if tzErr != nil && exitOnFail {
				log.Panic(tzErr)
			}

			var startParseErr error

			rule.DtStart, startParseErr = time.Parse(time.RFC3339, event.Start.DateTime)

			if startParseErr != nil && exitOnFail {
				log.Panic(startParseErr)
			}

			rule.DtStart = rule.DtStart.In(timeZone)

			iter := rule.Iterator().
				After(allowAfter).
				Before(allowBefore)

			var start time.Time
			var index int = 0

			for iter.Step(&start) {
				emulated_event := eventexpander.CreateInstanceFromSource(event, start)
				enc.Encode(emulated_event)
				index += 1

				if maxIterations > 0 && index >= maxIterations {
					break
				}
			}
		}
	},
}
