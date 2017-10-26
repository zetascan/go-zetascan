# Developing a Zetascan library

If developing a new language library/SDK for Zetascan, developers are asked to adhere to the following structure, use the same method names, return data and unit-testing.

Examples below provided in Go, but general framework is language independent.

## Init

```go
	myzetascan, err = myzetascan.Init(apiKey, ipAuth)

	if err != nil {
		fmt.Println(err)
	}
```

Requirements: Pass an API key as an argument, optionally a flag for IP only access.

Errors must be returned and handled accordingly.

## Query Method

```go
	// Query via the JSON method
	myzetascan.ApiMethod = "json"
	m, _ := myzetascan.Query(query)
```

Allow the developer to choose which Zetascan method to use (json, jsonx, txt, http and optionally, dns) - However a query method is a single function, and the individual logic behind each method is hidden from the developer. They simply need to define which method to use, and call the `Query(arguments)` method.

## Returned data

In the Go and PHP library for Zetascan, internally each method (json, jsonx, txt, http, dns) insert returned data into a defined JSON format, and return a similar object. If a developer calls the txt or http method, the same return object is expected.

For example, reference the following Go code, which builds a struct, containing the same key/values as per the Zetascan JSONx documentation. This struct is used by any query method, regardless if all fields are used.

```go
// Format for JSON and JSONX responses

type JsonReason struct {
	Class       string `json:"class"`
	Rule        string `json:"rule"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Source      string `json:"source"`
	Port        string `json:"port"`
	SourcePort  string `json:"sourceport"`
	Destination string `json:"destination"`
}

type JsonExtended struct {
	ASNum   string     `json:"ASNum"`
	Route   string     `json:"route"`
	Country string     `json:"country"`
	Domain  string     `json:"domain"`
	State   string     `json:"state"`
	Time    string     `json:"time"`
	Reason  JsonReason `json:"reason"`
}

type JsonResults []struct {
	Item       string       `json:"item"`
	Found      bool         `json:"found"`
	Score      float64      `json:"score"`
	WebScore   float64      `json:"webscore"`
	FromSubnet bool         `json:"fromSubnet"`
	Sources    []string     `json:"sources"`
	Wl         bool         `json:"wl"`
	Wldata     string       `json:"wldata"`
	Extended   JsonExtended `json:"extended"`
}

type JsonRecord struct {
	Results       JsonResults `json:"results"`
	ExecutionTime int64       `json:"executionTime"`
	Status        string      `json:"status"`
}
```

# Methods

When developing a library for Zetascan, the following methods must be implemented.

## Query

```go
// Query a domain/IP via any method (text, html, json, jsonx, dns)
func (myapi Api) Query(query string) (m JsonRecord, err error) {
```

## getUrl

Method to set the domain end-point for Zetascan (default api.zetascan.com) - If Zetascan is used on-prem, allow the developer to set to their internal domain.

```go
// getUrl Return a URL to query zetascan
func (myapi Api) getUrl(domain string) string {
```

## parseResult

Parse the results from any query type, and return a common object. Reference (Returned data) for more information

```go
// parseResult returns a struct with the zetascan response, regardless of the query method
func (myapi Api) parseResult(resp *http.Response) (data JsonRecord, err error) {
```

## IsMatch

Return if a result matched a whitelist/blacklist

```go
func (myapi Api) IsMatch(response *JsonRecord) (status bool) {
```

## IsWhiteList

Return if a result matched a whitelist

```go
func (myapi Api) IsWhiteList(response *JsonRecord) (status bool) {
```

## IsBlackList

Return if a result matched a blacklist

```go
// IsBlackList return if a result matched a blacklist
func (myapi Api) IsBlackList(response *JsonRecord) (status bool) {
```

## Score

Return the score if a result matched a whitelist/blacklist

```go
func (myapi Api) Score(response *JsonRecord) (score float64) {
```

# Unit testing

Before submitting a library for Zetascan, simple unit tests must be provided that validate the test IPs/Domains successfully pass/fail, for each query method.

```go
// Verify a query to zetascan is returning valid data
func (myapi Api) Verify(status bool, verbose bool) (totalResults []Results, err error) {

	tests := make(map[string]bool)

	// Records that will pass (whitelist)
	tests["okdomain.org"] = false
	tests["127.9.9.4"] = false

	// Records that will fail (blacklisted)
	tests["baddomain.org"] = true
	tests["127.9.9.1"] = true
	tests["127.9.9.2"] = true
	tests["127.9.9.3"] = true

	// Continue tests
	...
```