package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/scrollodex/ResourceUtils/air2hugo/pkg/airtableclient"
	"github.com/scrollodex/ResourceUtils/air2hugo/pkg/buildit"
)

var (
	jobIsRunning   bool
	JobIsrunningMu sync.Mutex
)

// Cite: https://stackoverflow.com/a/52793706/71978

func maybeStartJob(atc *airtableclient.AirClient) {
	JobIsrunningMu.Lock()
	start := !jobIsRunning
	jobIsRunning = true
	JobIsrunningMu.Unlock()
	if start {
		go func() {
			theJob("public/buildlog.txt", atc)
			JobIsrunningMu.Lock()
			jobIsRunning = false
			JobIsrunningMu.Unlock()
		}()
	} else {
		if debugFlag {
			log.Printf("Nope!\n")
		}
	}
}

func theJob(filename string, atc *airtableclient.AirClient) {

	startTime := time.Now()

	// if filename's filesize is greater than 1mb, delete it and start fresh
	fi, err := os.Stat(filename)
	if err == nil {
		if fi.Size() > 1*1024*1024 {
			err = os.Remove(filename)
			if err != nil {
				log.Printf("Error deleting log file: %s", err)
			}
		}
	}

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err) // Ok to panic.
	}
	defer f.Close()
	log.SetOutput(f)
	fmt.Fprintln(f)
	fmt.Fprintln(f)
	fmt.Fprintln(f)
	fmt.Fprintln(f, "==============================================")
	fmt.Fprintln(f, "==============================================")
	fmt.Fprintln(f, "==============================================")

	log.Println("RUNNING: ", startTime.Format("2006-01-02 3:4:5 PM"))

	// Refresh the data:

	log.Printf("AIRTABLE: Downloading data from Airtable\n")
	err = buildit.RefreshData(atc, debugFlag)
	if err == nil {
		log.Println("AIRTABLE: DONE!")
	} else {
		log.Printf("...ERROR: %s", err)
		return
	}

	log.Printf("HUGO: Generating website")
	cmd := exec.Command("hugo", "--noTimes")
	// --noTimes is required because DigitalOcean
	cmd.Stdout = f
	cmd.Stderr = f

	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Printf("hugo failed: %v", err)
	}

	endTime := time.Now()
	log.Println("DONE: ", startTime.Format("2006-01-02 3:4:5 PM"))
	log.Println("DONE: elapsed time: ", endTime.Sub(startTime).Round(time.Second/10))
	fmt.Fprintln(f, "==============================================")
	fmt.Fprintln(f, "==============================================")
	fmt.Fprintln(f, "==============================================")
	fmt.Fprintln(f)
	fmt.Fprintln(f)
	fmt.Fprintln(f)
}
