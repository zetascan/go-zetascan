package metascanjson

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
)

func main() {
	const jsonStream = `
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
	`

	const jsonStream2 = `
	{"results":[{"item":"127.9.9.1","found":true,"score":0.95,"fromSubnet":false,"sources":["shXBL","shSBL"],"extended":{"ASNum":"23969","route":"1.0.200.0/24","country":"AU","domain":"veridas.net","state":"","time":"1486447729","reason":{"class":"BOT","rule":"9904","type":"sinkhole","name":"conficker","source":"104.244.14.252","port":"80","sourceport":"23915","destination":"1"},"emailslastday":"0"},"wl":false,"wldata":"","lastModified":1500972200}],"executionTime":1,"status":"success"}
	`
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

	dec := json.NewDecoder(strings.NewReader(jsonStream))
	for {
		var m JsonRecord
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s: %s\n", m)
	}

	dec = json.NewDecoder(strings.NewReader(jsonStream2))
	for {
		var m JsonRecord
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		fmt.Println(m)
		fmt.Printf("%s: %s\n", m)
	}

}
