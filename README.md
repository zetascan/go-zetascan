# go-metascan
#### Development SDK/library for the metascan API in Go

## Introduction
 
The Metascan Query Services "MQS" was created to facilitate the real-time lookup of IP and Domain threat data into various applications and services. Currently there are dozens of various domain and IP data-feeds available to developers. Many of these feeds are available free of charge and some are paid for services when minimum query levels are exceeded. In addition, there are 2 main problems with trying to incorporate multiple data feed into a solution:

* The overlap between data feed providers in the content listed (IPs & URIs), and

* The absence of normalized meta-data related to the IPs or Domains.

Because of the above, many developers asked if we could do something to reduce the complexity related to accessing and using threat data as part of their applications - MQS is our solution. We are introducing a more elegant API for developers, with an affordable pricing model to match.

To start, [signup for a developer key](https://metascan.io/signup/?lang=en) and begin to integrate MQS into your web-apps and mobile applications.

## go-metascan 
The go-metascan provided an API interface to query the metascan service via HTTP or DNS, and provides examples on integrating your service, web-app or mobile application to prevent abuse via the metascan service.

## Examples

Build the metascan-query utility to provide simple CLI tools to query the metascan service as an example.

```
cd go-metascan
go build metascan-query.go
```

Or alternatively, run directly

```
go run metascan-query.go
```

### Usage:

```
./metascan-query -h

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

Query the metascan service using the JSON API method. View the [developer docs](http://docs.metascan.io/) for more information on the methods available.

```
./metascan-query -query baddomain.org -apikey YOURAPIKEY -format json

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

Query the metascan service using the DNS method. View available test IP and domains to query form the [developer docs](http://docs.metascan.io/#ip-addresses)

To use DNS, you must add your servers IP address to the metascan developer portal. An API key is not available.

```
./metascan-query.go -query 127.9.9.1 -ipauth -format dns

{Results:[
	{
		Item: Found:true
		...
	}
]}

```

## Developer example

See example.go

```
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
```