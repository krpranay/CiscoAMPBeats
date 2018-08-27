// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

//import "time"
import (
	"errors"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/CiscoAMPBeats/logmanager"
	//"fmt"
)

type Config struct {
	Period      time.Duration `config:"period"`
	APIClientID string        `config:"api_client_ID"`
	APIKey      string        `config:"api_Key"`
	INIPath     string        `config:"inipath"`
	POSFILEPATH string        `config:"posfilepath"`
}

var DefaultConfig = Config{
	Period:      4 * time.Second,
	APIClientID: "f6cab9156394e0bc768b",
	APIKey:      "1fcd292c-2059-4235-8f76-62d1a6cf9db3",
	INIPath:     "CiscoAMPEndPoint.ini",
	POSFILEPATH: "CiscoAMPEndPoint.pos",
}

// Initializing different event types
// var EventTypes = map[string]int{
// 	"Adobe Reader Compromise":                1107296261,
// 	"Adobe Reader Launched a Shell":          1107296266,
// 	"All Fault Cleared":                      553648197,
// 	"APK Custom Threat Detected":             1090524041,
// 	"APK Threat Detected":                    1090524040,
// 	"Apple QuickTime Compromise":             1107296270,
// 	"Apple QuickTime Launched a Shell":       1107296271,
// 	"Application Authorized":                 570425398,
// 	"Application Deauthorized":               570425399,
// 	"Application Deregistered":               570425397,
// 	"Application Registered":                 570425396,
// 	"Attempting Quarantine Delete":           553648151,
// 	"Cisco AMP4EP - Base Rule ":              0, // value not defined in the configuration guide file, setting as zero
// 	"Cloud Recall Quarantine Attempt":        553648155,
// 	"Cloud Recall Quarantine Attempt Failed": 2164260893,
// 	"Cloud Recall Quarantine of	False Negative": 553648147,
// 	"Cloud Recall Quarantine Successful":           553648155,
// 	"Cloud Recall Restore from Quarantine":         553648154,
// 	"Cloud Recall Restore from Quarantine Failed":  2164260892,
// 	"Cloud Recall Restore of False Positive":       553648146,
// 	"Cognitive Incident":                           1107296285,
// 	"Connection to Suspicious Domain":              1107296277,
// 	"Critical Fault Raised":                        2164260931,
// 	"DFC Threat Detected ":                         1090519084,
// 	"Email Confirmation":                           1003,
// 	"Endpoint IOC Configuration Update Failure":    2164260911,
// 	"Endpoint IOC Configuration Update Success":    553648176,
// 	"Endpoint IOC Definition Update Failure":       2164260914,
// 	"Endpoint IOC Definition Update Success":       553648179,
// 	"Endpoint IOC Scan Completed With Detections":  1091567670,
// 	"Endpoint IOC Scan Completed, No Detections":   554696757,
// 	"Endpoint IOC Scan Detection Summary":          1090519089,
// 	"Endpoint IOC Scan Failed":                     2165309495,
// 	"Endpoint IOC Scan Started":                    554696756,
// 	"Executed Malware":                             1107296272,
// 	"Execution Blocked":                            553648168,
// 	"Exploit Prevention":                           1090519103,
// 	"Failed to Delete From Quarantine":             2164260889,
// 	"Fault Cleared":                                553648196,
// 	"File Fetch Completed":                         553648173,
// 	"File Fetch Failed":                            2164260910,
// 	"Forgotten Password Reset":                     1004,
// 	"Generic IOC":                                  1107296274,
// 	"Install Failure ":                             2164260895,
// 	"Install Started":                              553648158,
// 	"Java Compromise":                              1107296260,
// 	"Java Launched a Shell":                        1107296265,
// 	"Major Fault Raised":                           1090519107,
// 	"Microsoft Calculator Compromise":              1107296275,
// 	"Microsoft CHM Compromise":                     1107296281,
// 	"Microsoft Excel Compromise":                   1107296263,
// 	"Microsoft Excel Launched a Shell":             1107296268,
// 	"Microsoft Notepad":                            1107296276,
// 	"Microsoft PowerPoint Compromise":              1107296264,
// 	"Microsoft PowerPoint Launched a Shell":        1107296269,
// 	"Microsoft Word Compromise":                    1107296262,
// 	"Microsoft Word Launched a Shell":              1107296267,
// 	"Minor Fault Raised":                           553648195,
// 	"Multiple Infected Files":                      1107296257,
// 	"Password Has Been Reset":                      1005,
// 	"Policy Update":                                553648130,
// 	"Policy Update Failure":                        2164260866,
// 	"Potential Dropper Infection":                  1107296258,
// 	"Potential Ransomware ":                        1107296284,
// 	"Potential Webshell":                           1107296283,
// 	"Product Update Completed":                     553648136,
// 	"Product Update Failed":                        553648137,
// 	"Product Update Started":                       553648135,
// 	"Quarantine Failure":                           2164260880,
// 	"Quarantine Item Deleted":                      553648152,
// 	"Quarantine Item Restored":                     553648149,
// 	"Quarantine Request Failed To Be Delivered":    2181038130,
// 	"Quarantine Restore Failed":                    2164260884,
// 	"Quarantine Restore Requested":                 570425394,
// 	"Quarantine Restore Started":                   553648150,
// 	"Quarantined Item Deleted":                     553648152,
// 	"Reboot Completed":                             553648171,
// 	"Reboot Pending":                               553648170,
// 	"Rootkit Detection":                            1090519081,
// 	"Scan Completed With Detections":               1091567628,
// 	"Scan Completed, No Detections":                554696715,
// 	"Scan Failed":                                  2165309453,
// 	"Scan Started":                                 554696714,
// 	"Suspected Botnet Connection":                  1107296273,
// 	"Suspicious Cscript Launch":                    1107296282,
// 	"Suspicious Download":                          1107296280,
// 	"Threat Detected":                              1090519054,
// 	"Threat Detected in Exclusion":                 553648145,
// 	"Threat Detected In Low Prevalence Executable": 1107296278,
// 	"Threat Quarantined":                           553648143,
// 	"Uninstall":                                    553648166,
// 	"Uninstall Failure":                            2164260903,
// 	"Update: Reboot Advised":                       1090519097,
// 	"Update: Reboot Required":                      1090519096,
// 	"Update: Unexpected Reboot Required":           2164260922,
// 	"Vulnerable Application Detected":              1107296279,
// }

var re *regexp.Regexp
var pat = "[#].*\\n|\\s+\\n|\\S+[=]|.*\n"

func init() {
	re, _ = regexp.Compile(pat)
}

// function to read INI file, getting all fields in a map
func ReadINIFile(inipath string) (map[string]string, error) {
	var mymap = make(map[string]string)
	err := load(inipath, mymap)
	if err != nil {
		logmanager.Log(logmanager.ERROR, "Unable to read INI file - ", err)
	}
	return mymap, err
}

// Load adds or updates entries in an existing map with string keys
// and string values using a configuration file.
//
// The filename paramter indicates the configuration file to load ...
// the dest parameter is the map that will be updated.
//
// The configuration file entries should be constructed in key=value
// syntax.  A # symbol at the beginning of a line indicates a comment.
// Blank lines are ignored.
func load(filename string, dest map[string]string) error {
	fi, err := os.Stat(filename)
	if err != nil {
		return err
	}
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	buff := make([]byte, fi.Size())
	f.Read(buff)
	f.Close()
	str := string(buff)
	if !strings.HasSuffix(str, "\n") {
		str += "\n"
	}
	s2 := re.FindAllString(str, -1)

	for i := 0; i < len(s2); {
		if strings.HasPrefix(s2[i], "#") || strings.HasPrefix(s2[i], "[") {
			i++
		} else if strings.HasSuffix(s2[i], "=") {
			key := strings.ToLower(s2[i])[0 : len(s2[i])-1]
			i++
			if strings.HasSuffix(s2[i], "\n") {
				val := s2[i][0 : len(s2[i])-1]
				if strings.HasSuffix(val, "\r") {
					val = val[0 : len(val)-1]
				}
				i++
				dest[key] = val
			}
		} else if strings.Index(" \t\r\n", s2[i][0:1]) > -1 {
			i++
		} else {
			return errors.New("Unable to process line in cfg file containing " + s2[i])
		}
	}
	return nil
}
