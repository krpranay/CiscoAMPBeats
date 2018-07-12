package beater

import (
	"bytes"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Ciscoamp struct {
	NextURL       string
	LastQueryTime time.Time
	Offst         string
}

const C_INTERVAL = 60000 // every 10 minutes
const C_LIMIT = 10

func (ciscoamp *Ciscoamp) GetEvents(clientid string, apikey string) string {
	client := &http.Client{}

	// ISO8601 2015-10-01T00:00:00+00:00
	//	t := time.Now()
	//	then := t.Add(time.Second * C_INTERVAL * -1)
	//start_date := then.Format(time.RFC3339)
	//url := "https://api.amp.cisco.com/v1/events?limit=" + strconv.Itoa(C_LIMIT) + "&start_date=" + start_date + "&event_type[]=1090519054"
	//url := "https://api.amp.cisco.com/v1/events?limit=" + strconv.Itoa(C_LIMIT) + "&start_date=" + start_date
	url := "https://api.amp.cisco.com/v1/events?limit=" + strconv.Itoa(C_LIMIT)
	log.Printf("API Call made to %s", url)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	req.SetBasicAuth(clientid, apikey) // this uses base 64 encoding, which doesn't currently work

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	buf := new(bytes.Buffer)
	body := resp.Body
	buf.ReadFrom(body)
	js := buf.String()

	defer resp.Body.Close()

	return js
}
