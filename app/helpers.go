package main

import (
	"encoding/csv"
	"errors"
	"os"
	"strings"
)

var ErrNoRecord = errors.New("no matching record found")
var ErrCsvCreation = errors.New("failed to create csv")

var iplHeader = []string{
	"subnet", "consumer_ipl", "provider_apps",
	"transmission", "port", "protocol", "num_flows",
	"conn_state", "reported policy decision",
	"reported by", "draft policy decision",
}

// Import CSV and parse the data.
func parseCsv(file string) (*[][]string, error) {
	csvfile, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	reader.FieldsPerRecord = -1

	rawdata, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return &rawdata, nil
}

func consolidateLabels(labels ...string) string {
	var combined string

	for _, label := range labels {
		if label == "" {
			combined += "NO_LABEL"
		}

		combined += label + " | "
	}

	return strings.TrimSuffix(combined, " | ")
}

// Check if the column is empty. Assign "NIL" when empty.
func checkIfEmpty(input string) string {
	if input == "" {
		return "NIL"
	}

	return input
}

// Remove chars from the last dot in IP.
func removeLastOctet(input string) string {
	if len(input) > 0 {
		if i := strings.LastIndex(input, "."); i > 0 {
			input = input[:i]
		}
	}

	return input
}

// Deduplicate and combine record.
func consolidateRecord(data [][]string) *[][]string {
	return &data
}
