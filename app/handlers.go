package main

import (
	"encoding/csv"
	"os"
	"sync"
)

// Handling struct conversion.
// Place 'NIL' when the value is empty. Exception for labels column as 'NO_LABEL'.
func (app *application) recordHandler() error {
	recordsLen := len(*app.rawData)

	if recordsLen == 0 {
		return ErrNoRecord
	}

	// column name length
	columnNo := len((*app.rawData)[0])
	app.infoLog.Printf("Total columns in record: %d", columnNo)

	for i := 1; i < recordsLen; i++ {
		r := &record{
			c: consumer{
				hostname: checkIfEmpty((*app.rawData)[i][app.indexOf("Consumer Hostname")]),
				ip:       (*app.rawData)[i][app.indexOf("Consumer IP")],
				ipl:      checkIfEmpty((*app.rawData)[i][app.indexOf("Consumer IPList")]),
				fqdn:     checkIfEmpty((*app.rawData)[i][app.indexOf("Consumer FQDN")]),
				appgroup: consolidateLabels(
					(*app.rawData)[i][app.indexOf("Consumer Role")],
					(*app.rawData)[i][app.indexOf("Consumer App")],
					(*app.rawData)[i][app.indexOf("Consumer Env")],
					(*app.rawData)[i][app.indexOf("Consumer Loc")],
				),
			},
			p: provider{
				hostname: checkIfEmpty((*app.rawData)[i][app.indexOf("Provider Hostname")]),
				ip:       (*app.rawData)[i][app.indexOf("Provider IP")],
				ipl:      checkIfEmpty((*app.rawData)[i][app.indexOf("Provider IPList")]),
				fqdn:     checkIfEmpty((*app.rawData)[i][app.indexOf("Provider FQDN")]),
				appgroup: consolidateLabels(
					(*app.rawData)[i][app.indexOf("Provider Role")],
					(*app.rawData)[i][app.indexOf("Provider App")],
					(*app.rawData)[i][app.indexOf("Provider Env")],
					(*app.rawData)[i][app.indexOf("Provider Loc")],
				),
			},
			i: info{
				transmission:   (*app.rawData)[i][app.indexOf("Transmission")],
				port:           (*app.rawData)[i][app.indexOf("Port")],
				protocol:       (*app.rawData)[i][app.indexOf("Protocol")],
				num_flows:      (*app.rawData)[i][app.indexOf("Num Flows")],
				conn_state:     checkIfEmpty((*app.rawData)[i][app.indexOf("Connection State")]),
				first_detected: (*app.rawData)[i][app.indexOf("First Detected")],
				last_detected:  (*app.rawData)[i][app.indexOf("Last Detected")],
			},
			r: reported{
				policy_decision:      checkIfEmpty((*app.rawData)[i][app.indexOf("Reported Policy Decision")]),
				enforcement_boundary: checkIfEmpty((*app.rawData)[i][app.indexOf("Reported Enforcement Boundary")]),
				by:                   checkIfEmpty((*app.rawData)[i][app.indexOf("Reported by")]),
			},
			d: draft{
				policy_decision:      (*app.rawData)[i][app.indexOf("Draft Policy Decision")],
				enforcement_boundary: checkIfEmpty((*app.rawData)[i][app.indexOf("Draft Enforcement Boundary")]),
			},
		}
		app.records = append(app.records, r)
	}

	return nil
}

// Search for the column index.
func (app *application) indexOf(input string) int {
	for k, v := range (*app.rawData)[0] {
		if input == v {
			return k
		}
	}

	return -1 //not found.
}

// Output as csv format.
func (app *application) saveAsCsv(data [][]string, name string) {
	csvFile, err := os.Create(name + ".csv")
	if err != nil {
		app.errorLog.Fatal(ErrCsvCreation)
	}
	defer csvFile.Close()

	csvwriter := csv.NewWriter(csvFile)

	for i := range data {
		_ = csvwriter.Write(data[i])
	}

	csvwriter.Flush()
	app.infoLog.Printf("File: %s created", name)
}

// Create a slice to contains data which the consumer is IPList
func (app *application) ConsumerIsIPL(header []string) [][]string {
	ss := [][]string{}
	ss = append(ss, header)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var count = 0

	for i := range app.records {
		wg.Add(1)
		s1 := []string{}
		go func(i int) {
			mu.Lock()
			if app.records[i].c.ipl != "NIL" {
				s1 = append(s1,
					removeLastOctet(app.records[i].c.ip),
					app.records[i].c.ipl,
					app.records[i].p.appgroup,
					app.records[i].i.transmission,
					app.records[i].i.port,
					app.records[i].i.protocol,
					app.records[i].i.num_flows,
					app.records[i].i.conn_state,
					app.records[i].r.policy_decision,
					app.records[i].r.by,
					app.records[i].d.policy_decision)
				ss = append(ss, s1)
				count++
			}
			mu.Unlock()
			wg.Done()
		}(i)
	}
	wg.Wait()

	app.infoLog.Printf("Total %d records - Consumer As IPList", count)
	return ss
}
