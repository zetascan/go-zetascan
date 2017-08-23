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
	verify := flag.Bool("verify", false, "Verify authentication and query")
	format := flag.String("format", "", "Specify the query format (text, http, json, jsonx)")

	flag.Parse()

	fmt.Println("Welcome", *apiKey)

	var err error
	var mymetascan metascan.Api

	mymetascan, err = mymetascan.Init(*apiKey, *ipAuth)

	fmt.Println(err)

	if err != nil && *ipAuth == false {
		log.Fatal("Please specify an API key, or specify -ipauth to disable")
	}

	if *verify == true {

		fmt.Println("Running metascan diagnostics")

		var kvs []string

		if *format == "" {
			// If no format specified, use all
			kvs = []string{"text", "http", "json", "jsonx"}
		} else {
			// Otherwise use the specified format provided
			kvs = []string{*format}
		}

		for _, value := range kvs {
			mymetascan.ApiMethod = value
			fmt.Println("Testing ", value)
			fmt.Println(mymetascan.Verify(true))
		}

	}

	fmt.Println(mymetascan.GetConf())
}
