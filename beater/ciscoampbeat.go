package beater

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/CiscoAMPBeats/config"
	"github.com/CiscoAMPBeats/logmanager"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
)

var nurl = ""
var url = ""
var requestCounter = 0
var lastUrl = ""
var filename = ""
var offsetPosition = ""

// Ciscoampbeat configuration.
type Ciscoampbeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
}

// New creates an instance of ciscoampbeat.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		logmanager.Log(logmanager.ERROR, "Error reading config file - ", err)
		return nil, nil // fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Ciscoampbeat{
		done:   make(chan struct{}),
		config: c,
	}
	return bt, nil
}

// Run starts ciscoampbeat.
func (bt *Ciscoampbeat) Run(b *beat.Beat) error {

	logp.Info("ciscoampbeat is running! Hit CTRL-C to stop it.")

	// Reading INI file
	iniFile, err1 := config.ReadINIFile(bt.config.INIPath)
	jsonStr := iniFile["sample"]

	personMap := make(map[string]interface{})

	json.Unmarshal([]byte(jsonStr), &personMap)

	if err1 == nil {

		// Creating position file
		filename = config.DefaultConfig.POSFILEPATH
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			//fmt.Println("File does not exist")
			var file, err = os.Create(filename)
			if err != nil {
				logmanager.Log(logmanager.ERROR, "Unable to create position file - ", err)
			}
			defer file.Close()

			// Writing offset value in position file
			writePositionFile(filename, "0")

		} else {
			//fmt.Println("File exists")
			offsetValue, dtValue := readPositionFile(filename)
			fmt.Println(offsetValue + "," + dtValue)
			iniFile["offsetValue"] = offsetValue
			iniFile["dtValue"] = dtValue

		}

		// Creating API request URL

		//url = createRequestURL(config.EventTypes, iniFile)
		url = createRequestURL(iniFile)

		var err error
		bt.client, err = b.Publisher.Connect()
		if err != nil {
			logmanager.Log(logmanager.ERROR, "Failed to call the AMP4EP API, Exiting", err)
		}

		camp := Ciscoamp{}

		ticker := time.NewTicker(bt.config.Period)
		counter := 1
		for {
			select {
			case <-bt.done:
				return nil
			case <-ticker.C:
			}

			// maintaing a temprory variable for last url
			if url != "" {
				lastUrl = url
			}

			if requestCounter != 0 && strings.Contains(url, "offset=") {

				offset := strings.Split(url, "offset=")
				fmt.Println(offset[1])
				offsetPosition = offset[1]
			}

			requestCounter++
			events, re, currentitemcount, apierror := camp.GetEvents(iniFile["clientid"], iniFile["apikey"], url)
			if apierror == nil {
				url = re
				if currentitemcount != 0 {
					event := beat.Event{
						Timestamp: time.Now(),
						Fields: common.MapStr{
							"type":     b.Info.Name,
							"response": events,
						},
					}
					bt.client.Publish(event)
					logp.Info("Event sent")
					counter++
				}
			}
		}
	}
	return err1
}

// func itemExist(v interface{}) bool {
// 	x := reflect.TypeOf(v)
// 	switch x.Kind() {
// 	case reflect.Slice:
// 	case reflect.Array:
// 		s := reflect.ValueOf(v)
// 		return s.Len() > 0
// 	default:
// 		return false
// 	}
// 	return false
// }

// Stop stops ciscoampbeat.
func (bt *Ciscoampbeat) Stop() {
	fmt.Println("Stopping CiscoAMP Beats")

	bt.client.Close()
	close(bt.done)
	writePositionFile(filename, offsetPosition)
}

// function to create and write position file
func writePositionFile(filename string, offsetValue string) {
	// open file using READ & WRITE permission
	var file, err = os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		logmanager.Log(logmanager.ERROR, "Unable to find Position file", err)
	}
	defer file.Close()

	// write some text line-by-line to file
	file.Truncate(0)
	_, err = file.WriteString(offsetValue)
	if err != nil {
		logmanager.Log(logmanager.ERROR, "Error in writing Position file", err)
	}

	// save changes
	err = file.Sync()
	if err != nil {
		logmanager.Log(logmanager.ERROR, "", err)
	}

}

func readPositionFile(filename string) (string, string) {
	var lineRead = 0
	var offsetValue = ""
	var dtValue = ""
	fileHandle, err := os.Open(filename)
	if err != nil {
		logmanager.Log(logmanager.ERROR, "Error in opening Position file", err)
	}

	defer fileHandle.Close()
	fileScanner := bufio.NewScanner(fileHandle)

	for fileScanner.Scan() {
		if lineRead == 0 {
			offsetValue = fileScanner.Text()
		}
		if lineRead == 1 {
			dtValue = fileScanner.Text()

		}
		lineRead++
	}
	return offsetValue, dtValue
}
