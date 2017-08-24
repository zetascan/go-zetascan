package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/metascanio/go-metascan/metascan"
)

func main() {

	apiKey := flag.String("apikey", "", "Specify API key")
	ipAuth := flag.Bool("ipauth", false, "Toggle to bypass API key and use IP authentication")

	// Verification steps
	verify := flag.Bool("verify", false, "Verify authentication and query")
	csv := flag.Bool("csv", false, "Toggle to output in CSV for -verify flag")
	count := flag.Int("count", 1, "Number of time to run tests, when -verify set")

	//
	format := flag.String("format", "", "Specify the query format (text, http, json, jsonx, dns)")
	verbose := flag.Bool("verbose", false, "Enable verbose debug log")

	// query
	// TODO: Add support for multiple queries seperated by a comma
	query := flag.String("query", "", "Specifiy domain or IP to query")

	flag.Parse()

	// If no query or verification specfied, show usage and exit
	if *verify == false && *query == "" {
		flag.Usage()
		os.Exit(1)
	}

	var err error
	var mymetascan metascan.Api

	// Init with our API key
	mymetascan, err = mymetascan.Init(*apiKey, *ipAuth)

	if err != nil {
		fmt.Println(err)
	}

	if err != nil && *ipAuth == false {
		log.Fatal("Please specify an API key, or specify -ipauth to disable")
	}

	// Verify the test IP's provided by metascan are accessible
	if *verify == true {

		if *verbose == true {
			fmt.Println("Running metascan diagnostics")
		}

		var kvs []string

		// Fetch the local host used for reporting
		hostname, _ := os.Hostname()

		// For benchmarking/performance, run the tests the specified number of times
		for cnt := 0; cnt < *count; cnt++ {

			if *format == "" {
				// If no format specified, use all
				kvs = []string{"text", "http", "json", "jsonx", "dns"}
			} else {
				// Otherwise use the specified format provided
				kvs = []string{*format}
			}

			// Display the CSV header
			if *csv == true {
				//0 , Bens-MacBook.local , jsonx , 127.9.9.4 , false , 0 , false , 378
				fmt.Println("Count, Local hostname, method, query, expected result, content length, match, time in ms")
			}

			// Loop through all query methods
			for _, value := range kvs {
				mymetascan.ApiMethod = value

				if *csv == false {
					fmt.Println("Testing ", value)
				}

				results, _ := mymetascan.Verify(true, *verbose)

				for i := range results {

					if *csv == true {
						fmt.Println(cnt, ",", hostname, ",", value, ",", results[i].IP, ",", results[i].Expected, ",", results[i].Length, ",", results[i].Match, ",", results[i].TimeElapsed)
					} else {
						fmt.Println(cnt, hostname, ",", results[i])
					}

				}

			}

		}

	}

	// Run a specific query, return the results to STDOUT
	if *query != "" {

		mymetascan.ApiMethod = *format

		m, _ := mymetascan.Query(*query)

		if err != nil {
			fmt.Println(err)
		}

		//fmt.Println(m)
		fmt.Printf("%+v", m)

	}

}
