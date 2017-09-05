# go-zetascan
#### Development SDK/library for the zetascan API in Go

## Introduction
 
The [Zetascan Query Services](https://zetascan.com/) "ZQS" was created to facilitate the real-time lookup of IP and Domain threat data into various applications and services. Currently there are dozens of various domain and IP data-feeds available to developers. Many of these feeds are available free of charge and some are paid for services when minimum query levels are exceeded. In addition, there are 2 main problems with trying to incorporate multiple data feed into a solution:

* The overlap between data feed providers in the content listed (IPs & URIs), and

* The absence of normalized meta-data related to the IPs or Domains.

Because of the above, many developers asked if we could do something to reduce the complexity related to accessing and using threat data as part of their applications - MQS is our solution. We are introducing a more elegant API for developers, with an affordable pricing model to match.

To start, [signup for a developer key](https://zetascan.com/signup/?lang=en) and begin to integrate MQS into your web-apps and mobile applications.

## go-zetascan 
The go-zetascan library provides an API interface to query zetascan via HTTP or DNS, and provides examples on how to integrate your web-app or mobile application to prevent abuse.

## Examples

Build the zetascan-query utility to provide a simple CLI tool to query the service.

```
cd go-zetascan
go build zetascan-query.go
```

Or alternatively, run directly

```
go run zetascan-query.go
```

### Usage:

```
./zetascan-query -h

-query:
  -apikey string
    	Specify API key
  -count int
    	Number of time to run tests, when -verify set (default 1)
  -csv
    	Toggle to output in CSV for -verify flag
  -format string
    	Specify the query format (text, http, json, jsonx)
  -ipauth
    	Toggle to enable IP authentication
  -query string
    	Specifiy domain or IP to query
  -verbose
    	Enable verbose debug log
  -verify
    	Verify authentication and query
```

### Example domain query via JSON

Query the zetascan service using the JSON API method. View the [developer docs](http://docs.zetascan.com/) for more information on the methods available.

```
./zetascan-query -query baddomain.org -apikey YOURAPIKEY -format json

{Results:[
	{
		Item:baddomain.org
		Found:true
		Score:1
		FromSubnet:false
		
		Sources:[shDBL ubGrey ubGold ubRed ubBlack]
		Wl:false
		Wldata:
		
		...
	}
	ExecutionTime:0
	Status:success
]}

```

### Example IP query via DNS

Query the zetascan service using the DNS method. View available test IP and domains to query form the [developer docs](http://docs.zetascan.com/#ip-addresses)

To use DNS, you must add your servers IP address to the zetascan developer portal. An API key is not available.

```
./zetascan-query.go -query 127.9.9.1 -ipauth -format dns

{Results:[
	{
		Item: Found:true
		...
	}
]}

```

## Developer example

See examples/cli/test-query.go

```go
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
	if(len(os.Args) > 1) {
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
	score := myzetascan.Score(&m)

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
```

### Usage

Blacklist hit:

```
$ go run ./examples/cli/test-query.go baddomain.org

Querying baddomain.org
Blacklist hit, with a high score

```

Whitelist hit:

```
$ go run ./examples/cli/test-query.go 127.9.9.4

Querying 127.9.9.4
Whitelist hit, trusted record
```


## Web-server Example

A sample go HTTP server is provided, which upon receiving a request, initiates a DNS lookup to Zetascan if the clients IP address is contained in a blacklist, and throws a HTTP 403 response code (Forbidden) if matched.

A useful example to protect user signups forms, API end-points and critical services from known bot-nets, spammers, and to prevent blacklisted IP's abusing your infrastructure.

```go
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
```

### Usage

```
$ go run ./examples/webserver/web-service.go
Launching a test webserver on port 8000
```

Next, launch `http://localhost:8000` in your browser, your remote IP address will be looked up via the Zetascan service and a 200 (OK) response returned if no match/whitelist, otherwise a 403 (Forbidden) response returned if listed in a known blacklist.

