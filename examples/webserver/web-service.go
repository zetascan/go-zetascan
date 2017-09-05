package main

import (
	"fmt"
	"net/http"

	"github.com/zetascanio/go-zetascan/zetascan"
)

func main() {
	fmt.Println("Launching a test webserver on port 8000")
	http.HandleFunc("/", hello)
	http.ListenAndServe(":8000", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {

	var err error
	var myzetascan zetascan.Api

	// Query the remote users address, are they blacklisted?
	query := r.RemoteAddr // "baddomain.org" , use for testing a blacklist hit

	apiKey := ""   // Speciy an IP key
	ipAuth := true // Auth via the IP address, which must be added via the zetascan developer portal

	// Init with our API key
	myzetascan, err = myzetascan.Init(apiKey, ipAuth)

	if err != nil {
		fmt.Println(err)
	}

	// Query via the JSON method
	myzetascan.ApiMethod = "dns"
	m, _ := myzetascan.Query(query)

	// Find the record score ( not supported via DNS, only DNS txt record)
	// The minimum score is -0.1, meaning that an item was found in White List only. Score 0 means that the item is not found in our DB, and the maximum score is 1. In general, items with score above 0.35 shall be considered as spam or fraud.
	//score := myzetascan.Score(&m)

	// If whitelist, trust
	if myzetascan.IsWhiteList(&m) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("200: OK - Whitelist hit, trusted record"))

	} else if myzetascan.IsBlackList(&m) {
		// If in a blacklist, throw a 403 error
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("403: Request denied - Blacklist hit!"))

	} else {
		// If no match, proceed as normal
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("200: OK - No blacklist/whitelist match found"))

	}

}
