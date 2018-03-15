package main

import (
	"fmt"
	"os"

	"github.com/zetascanio/go-zetascan/zetascan"
)

func main() {

	var err error
	var myzetascan zetascan.Api

	apiKey := ""   // Speciy an IP key
	ipAuth := true // Auth via the IP address, which must be added via the zetascan developer portal

	var query string

	// First argument is the domain or IP
	if len(os.Args) > 1 {
		query = os.Args[1]

	} else {
		query = "baddomain.org"
	}

	fmt.Println("Querying", query)

	// Init with our API key
	myzetascan, err = myzetascan.Init(apiKey, ipAuth)

	if err != nil {
		fmt.Println(err)
	}

	// Query via the JSON method
	myzetascan.ApiMethod = "json"
	m, _ := myzetascan.Query(query)

	// Find the record score
	// The minimum score is -0.1, meaning that an item was found in White List only. Score 0 means that the item is not found in our DB, and the maximum score is 1. In general, items with score above 0.35 shall be considered as spam or fraud.
	// Zetascan provides 2 scoring methods, the default score for MTA/SMTP use, or a WebScore used by web-apps
	//score := myzetascan.Score(&m)
	score := myzetascan.WebScore(&m)

	// If whitelist, trust
	if myzetascan.IsWhiteList(&m) {
		fmt.Println("Whitelist hit, trusted record")
	} else if myzetascan.IsBlackList(&m) && score > 0.35 {
		// If blacklist and high score
		fmt.Println("Blacklist hit, with a high score")
	} else if myzetascan.IsBlackList(&m) && score < 0.35 {
		// If blacklist low score
		fmt.Println("Blacklist hit, with a lower score")
	} else {
		// If no match
		fmt.Println("No blacklist/whitelist match found")
	}

	fmt.Println("\n\nRaw struct")
	fmt.Printf("%+v", m)

}
