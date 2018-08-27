package beater

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/CiscoAMPBeats/tidwall/gjson"

	"github.com/CiscoAMPBeats/logmanager"
)

// Ciscoamp struct defination
type Ciscoamp struct {
	NextURL       string
	LastQueryTime time.Time
	Offst         string
}

var nextAPIurl = ""

// GetEvents method call
func (ciscoamp *Ciscoamp) GetEvents(clientid string, apikey string, requestURL string) (interface{}, string, int64, error) {
	client := &http.Client{}

	url = requestURL

	log.Printf("API Call made to %s", url)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	req.SetBasicAuth(clientid, apikey) // this uses base 64 encoding, which doesn't currently work

	resp, err := client.Do(req)

	if err != nil {
		logmanager.Log(logmanager.ERROR, "The credentials within the credentials file are corrupt", err)
		return nil, "", -1, err
	}

	var m interface{}

	buf := new(bytes.Buffer)
	body := resp.Body
	buf.ReadFrom(body)
	js := buf.String()

	nextURL := gjson.Get(js, "metadata.links.next")

	totalcount := gjson.Get(js, "metadata.results.total")

	currentItemcount := gjson.Get(js, "metadata.results.current_item_count")
	currentPageIndex := gjson.Get(js, "metadata.results.index")
	itemPerPage := gjson.Get(js, "metadata.results.items_per_page")

	events := gjson.Get(js, "data")

	d := gjson.Unmarshal([]byte(events.Raw), &m)

	if d != nil {
		log.Fatalf("Unable to parse response  %v", err)
	}

	// getting next url from response and setting it as active requestURL

	if nextURL.Str == "" {
		selfURL := gjson.Get(js, "metadata.links.self")
		if currentItemcount.Num == 0 {
			nextAPIurl = selfURL.Str
		} else {
			var nextURLSubString = ""
			if currentPageIndex.Num == 0 {
				nextURLSubString = "offset=" + currentItemcount.String()
			} else {
				// totalRecordFetched := currentPageIndex.Int()*itemPerPage.Int() + currentItemcount.Int()
				// nextURLSubString = "offset=" + strconv.FormatInt(totalRecordFetched, 10)
				nextURLSubString = "offset=" + strconv.FormatInt(currentPageIndex.Int()+1, 10)
			}
			nextAPIurl = replaceAtIndex(selfURL.Str, nextURLSubString, strings.Index(selfURL.Str, "offset="))
		}
	} else {
		nextAPIurl = nextURL.Str
	}

	fmt.Println("total records available at server : ", totalcount)
	fmt.Println("Item returned for current page : ", currentItemcount)
	fmt.Println("Current page index: ", currentPageIndex)
	fmt.Println("Item per page : ", itemPerPage)

	defer resp.Body.Close()

	return m, nextAPIurl, currentItemcount.Int(), err
}

func replaceAtIndex(str string, substr string, index int) string {
	if len(str) == index+8 {
		return str[:index] + substr
	} else {
		return str[:index] + substr + str[index+8:]
	}
}

// function to create request URL
//func createRequestURL(eventTypes map[string]int, file map[string]string) string {
func createRequestURL(file map[string]string) string {

	offsetPosition := file["offsetValue"]
	datetimeValue := file["dtValue"]
	const urlHeader string = "https://"
	const event string = "events"

	if offsetPosition != "" && datetimeValue != "" {
	}

	var apiURL = ""

	// check if url address exists
	if file["uriaddress"] != "" {
		apiURL = urlHeader + file["uriaddress"] + file["version"] + "/" + event
	} else {
		//apiURL = "https://api.amp.cisco.com/v1/events"
		fmt.Println("URI address is not available in INI file ", file)
	}

	// if startdate is availabe append it in URL
	startdt := strings.TrimSuffix(file["startdate"], "\n")
	startdt = strings.Replace(startdt, "\r", "", 1)

	// if eventype is availabe append it in URL
	eventtype := strings.TrimSuffix(file["eventtypes"], "\n")
	eventtype = strings.Replace(eventtype, "\r", "", 1)

	apiURL = apiURL + "?"

	// getting DateTime read value from postion file or INI file
	if file["dtValue"] != "" {
		// if position file has last read data for datetime
		apiURL = apiURL + "start_date=" + file["dtValue"] + "&"

	} else if len(startdt) > 0 {
		// if position file has no last read data for datetime
		apiURL = apiURL + "start_date=" + file["startdate"] + "&"
	}

	// if limit is availbale append it in URL
	if file["limit"] != "" {
		apiURL = apiURL + "limit=" + file["limit"] + "&"
	}

	if file["offsetValue"] != "" {
		apiURL = apiURL + "offset=" + file["offsetValue"] + "&"
	}

	if file["eventtypes"] != "All" {

		var allEvents = strings.Split(file["eventtypes"], ",")
		for i := 0; i < len(allEvents); i++ {
			keyValue := allEvents[i]
			apiURL = apiURL + "event_type=" + keyValue + "&"
		}
	}

	apiURL = strings.TrimSuffix(apiURL, "&")

	fmt.Println("MyURL = " + apiURL)

	return apiURL
}
