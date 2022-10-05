package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	rawData  *[][]string
	records  []*record
}

func main() {
	file := flag.String("csv", "", "Import CSV file")
	//dedup := flag.Bool("dedup", true, "Dedup and combine records")
	flag.Parse()

	// Define logger
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Parse raw csv report
	rawdata, err := parseCsv(*file)
	if err != nil {
		errorLog.Fatal(err)
	}

	var records []*record

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		rawData:  rawdata,
		records:  records,
	}

	err = app.recordHandler()
	if err != nil {
		errorLog.Fatal(err)
	}

	var count = 0
	for i := range app.records {
		if app.records[i].c.ipl != "NIL" {
			fmt.Println(*app.records[i])
			count++
		}

	}

	fmt.Println(count)

	// c_ipl := app.ConsumerIsIPL(iplHeader)

	// if *dedup {
	// 	app.saveAsCsv(*consolidateRecord(c_ipl), "consumer_iplist")
	// }

}
