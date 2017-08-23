package metascan

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Api struct for key, URL and method
type Api struct {
	apiKey      string
	apiURL      string
	ApiMethod   string
	apiVersion  string
	apiProtocol string
}

type Query struct {
	apiKey   string
	apiQuery string
}

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
	Score      float32      `json:"score"`
	FromSubnet bool         `json:"fromSubnet"`
	Sources    []string     `json:"sources"`
	Wl         bool         `json:"wl"`
	Wldata     string       `json:"wldata"`
	Extended   JsonExtended `json:"extended"`
}

type JsonRecord struct {
	Results       JsonResults `json:"results"`
	ExecutionTime int         `json:"executionTime"`
	Status        string      `json:"status"`
}

// Init specify an authentication key for authentication
func (myapi Api) Init(apiKey string, ipcheck bool) (myapi2 Api, err error) {

	if apiKey != "" {
		myapi.apiKey = apiKey
		//return myapi, errors.New("API Key must be specified")
	}

	fmt.Println("Using ", apiKey)
	myapi.apiURL = "api.metascan.io"
	myapi.apiProtocol = myapi.ToggleSSL(true) // Default to SSL
	myapi.ApiMethod = "http"
	myapi.apiVersion = "v1"

	// Check
	if myapi.apiProtocol == "http" && apiKey != "" && ipcheck == false {
		return myapi, errors.New("https required if using API key without ip check")
	}

	return myapi, nil
}

// Verify a query to metascan is returning valid data
func (myapi Api) Verify(status bool) error {

	// Good
	tests := make(map[int]string)

	tests[0] = "okdomain.org"
	tests[1] = "127.9.9.4"

	tests[2] = "baddomain.org"
	tests[3] = "127.9.9.1"
	tests[4] = "127.9.9.2"
	tests[5] = "127.9.9.3"
	tests[6] = "127.9.9.4"

	for i := 0; i < len(tests); i++ {

		fmt.Println("Testing", tests[i])

		res, _ := http.Get(myapi.getUrl(tests[i]))

		match, err := myapi.isMatch(res)

		if err != nil {
			fmt.Println(err)
		}

		if match == true {
			fmt.Println(tests[i], ": FAIL!")
		} else {
			fmt.Println(tests[i], ": OK!")
		}

		fmt.Println("Resp => ", res, "\n")

	}

	return nil

}

// getUrl Return a URL to query metascan
func (myapi Api) getUrl(domain string) string {

	// Encode the apiKey if specified
	v := url.Values{}

	if myapi.apiKey != "" {
		v.Set("key", myapi.apiKey)
	}

	str := myapi.apiProtocol + "://" + myapi.apiURL + "/" + myapi.apiVersion + "/check/" + myapi.ApiMethod + "/" + domain + "?" + v.Encode()

	fmt.Println(str)

	return str
}

//
func (myapi Api) isMatch(resp *http.Response) (status bool, err error) {

	fmt.Println("\nisMatch launch")

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return false, err
	}

	switch myapi.ApiMethod {

	case "http":
		{

			if resp.StatusCode == 204 {
				return false, nil
			} else {
				return true, nil
			}

		}

	case "text":
		{

			// Read the body and split from the specified API formatting

			bodyString := string(body)
			head := strings.Split(bodyString, ":")
			data := strings.Split(head[1], ",")

			fmt.Println("Head =>", head, bodyString)
			fmt.Println("Data =>", data)

			/*
				http://docs.metascan.io/?php#http-format
				item:bool,bool,wldata,score,source

				Where:

				the first bool is true, if found in any black list,
				the second bool is true, if found in any white list,
				wldata contains the data from the white list, and
				score is followed by the list of sources where the item was found.
			*/
			if data[0] == "true" {
				return true, nil
			} else {
				return false, nil
			}

		}

	case "json", "jsonx":
		{

			/*
				http://docs.metascan.io/?php#json-format

				Formatting of a JSON response:

				{
					"results": [{
					"item": "123.123.123.123",
					"found": true,
					"score": 0.2,
					"fromSubnet": true,
					"sources": ["shPBL"],
					"wl": false,
					"wldata": ""
					}],
					"executionTime": 2,
					"status": "success"
				}
			*/

			dec := json.NewDecoder(strings.NewReader(string(body)))
			for {

				var m JsonRecord

				if err := dec.Decode(&m); err == io.EOF {
					return false, errors.New("JSON parse error")
				} else if err != nil {
					return false, err
				}

				fmt.Println(m)

				// Return if the query matched
				if m.Results[0].Found == true {
					return true, nil
				} else {
					return false, nil
				}

			}

		}

	}

	return false, nil

}

// Toggle SSL support
func (myapi Api) ToggleSSL(ssl bool) (str string) {

	if ssl == false {
		myapi.apiProtocol = "http"
	} else {
		myapi.apiProtocol = "https"
	}

	return myapi.apiProtocol

}

func (myapi Api) GetConf() string {

	return myapi.apiKey
}

// Query an IP(s) against the metascan API
func (myapi Api) QueryJSON(ip []string) (json string) {

	return ""

}

// Preform a DNS query against the metascan API
func (myapi Api) QueryDNS(ip []string) (json string) {

	return ""

}
