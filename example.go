package main

import (
	"fmt"

	"github.com/metascanio/go-metascan/metascan"
)

func main() {

	var err error
	var mymetascan metascan.Api

	apiKey := ""   // Speciy an IP key
	ipAuth := true // Auth via the IP address, which must be added via the metascan developer portal
	query := "baddomain.org"

	// Init with our API key
	mymetascan, err = mymetascan.Init(apiKey, ipAuth)

	if err != nil {
		fmt.Println(err)
	}

	mymetascan.ApiMethod = "json"

	m, _ := mymetascan.Query(query)

	if mymetascan.IsMatch(&m) {

		// Add logic, IP/domain is not trusted
	} else {

		// Proceed as normal

	}

	fmt.Println("Raw struct")
	fmt.Printf("%+v", m)

}
