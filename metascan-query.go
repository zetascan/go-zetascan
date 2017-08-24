package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/metascanio/go-metascan/metascan"
)

func main() {

	apiKey := flag.String("apikey", "", "Specify API key")
	ipAuth := flag.Bool("ipauth", false, "Toggle to enable IP authentication")

	// Verification steps
	verify := flag.Bool("verify", false, "Verify authentication and query")
	csv := flag.Bool("csv", false, "Toggle to output in CSV for -verify flag")
	count := flag.Int("count", 1, "Number of time to run tests, when -verify set")

	//
	format := flag.String("format", "", "Specify the query format (text, http, json, jsonx)")
	verbose := flag.Bool("verbose", false, "Enable verbose debug log")

	// query
	query := flag.String("query", "", "Specifiy domain or IP to query")

	flag.Parse()

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

		// For benchmarking/performance, run the tests the specified number of times
		for cnt := 0; cnt < *count; cnt++ {

			if *format == "" {
				// If no format specified, use all
				kvs = []string{"text", "http", "json", "jsonx", "dns"}
			} else {
				// Otherwise use the specified format provided
				kvs = []string{*format}
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
						fmt.Println(cnt, ",", value, ",", results[i].IP, ",", results[i].Expected, ",", results[i].Length, ",", results[i].Match, ",", results[i].TimeElapsed)
					} else {
						fmt.Println(cnt, results[i])
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
