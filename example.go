package main

import (
	"fmt"

	"github.com/zetascanio/go-zetascan/zetascan"
)

func main() {

	var err error
	var myzetascan zetascan.Api

	apiKey := ""   // Speciy an IP key
	ipAuth := true // Auth via the IP address, which must be added via the zetascan developer portal
	query := "baddomain.org"

	// Init with our API key
	myzetascan, err = myzetascan.Init(apiKey, ipAuth)

	if err != nil {
		fmt.Println(err)
	}

	myzetascan.ApiMethod = "json"

	m, _ := myzetascan.Query(query)

	if myzetascan.IsMatch(&m) {

		// Add logic, IP/domain is not trusted
	} else {

		// Proceed as normal

	}

	fmt.Println("Raw struct")
	fmt.Printf("%+v", m)

}
