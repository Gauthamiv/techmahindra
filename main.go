package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"encoding/json"

	"gopkg.in/yaml.v2"

	"errors"
	"net/http"
	"os"
	"strconv"
	"sync"

	"gopkg.in/robfig/cron.v2"
)

type ConfigDetails struct {
	Frequency int    `json:"frequency,omitempty"`
	Directory string `json:"data-directory,omitempty"`
}
type ResponseDetails struct {
	Name      string `json:"name,omitempty"`
	Url       string `json:"url,omitempty"`
	Site      string `json:"site,omitempty"`
	In24Hour  string `json:"in_24_hours,omitempty"`
	Status    string `json:"status,omitempty"`
	StartTime string `json:"start_time,omitempty"`
	EndTime   string `json:"end_time,omitempty"`
	Duration  string `json:"duration,omitempty"`
}
type ResponseDetailsArray struct {
	Resp []ResponseDetails
}

var config ConfigDetails
var errLog *log.Logger
var reqlog *log.Logger

func init() {

	var err error

	e, err := os.OpenFile("error.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Printf("error opening file: %v", err)
		os.Exit(1)
	}
	errLog = log.New(e, "", log.Ldate|log.Ltime)

	reql, err := os.OpenFile("custom.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Printf("error opening file: %v", err)
		os.Exit(1)
	}
	reqlog = log.New(reql, "", log.Ldate|log.Ltime)
	config, err = loadConfig("config.yml")
	if err != nil {
		os.Exit(1)
	}

}
func main() {
	reqlog.Println("In main function")
	var wg sync.WaitGroup
	wg.Add(1)
	go calling_cron(&wg)
	wg.Wait()

}
func calling_cron(wg *sync.WaitGroup) {
	reqlog.Println("Calling cron job")
	c := cron.New()
	c.AddFunc("@every "+strconv.Itoa(config.Frequency)+"h", get_response)
	c.Start()
}
func loadConfig(path string) (ConfigDetails, error) {
	reqlog.Println("Loading config file")
	content, err := ioutil.ReadFile(path)
	if err != nil {
		errLog.Println(err)
		return ConfigDetails{}, err
	}
	configYML := string(content)

	config := ConfigDetails{}
	err = yaml.Unmarshal([]byte(configYML), &config)
	if err != nil {
		errLog.Println(err)
		return ConfigDetails{}, err
	}

	return config, nil
}
func get_response() {
	reqlog.Println("Getting response")
	req, err := http.NewRequest("GET", "https://kontests.net/api/v1/all", nil)
	cli := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}

	resp, err := cli.Do(req)
	if err != nil || resp.StatusCode != 200 {

		fmt.Printf("Error is nil but status code is %d", resp.StatusCode)
		errLog.Printf("Error is nil but status code is %d", resp.StatusCode)
		err = errors.New("Error is nil but status code is not 200")
		return
	} else if err != nil {
		errLog.Println("Error while fetching response using API is: ", err)
		return
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errLog.Println("Error while reading response is: ", err)
		return
	}
	var content *[]ResponseDetails
	err = json.Unmarshal(bodyBytes, &content)
	if err != nil {
		errLog.Println("unmarshalling response failed", err)
		return
	}
	err = writting_into_file(content)
	if err != nil {
		return
	}
	return
}
func writting_into_file(content *[]ResponseDetails) error {
	reqlog.Println("Writting into json file")
	t1 := time.Now()

	t2 := t1.Format("2006-01-02T15-04-05")
	filew, err := os.OpenFile("response"+t2+".json", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Printf("error opening file: %v", err)
		errLog.Println("Error while opening file to write response")
		return err
	}
	_, err = filew.WriteString("")
	if err != nil {
		errLog.Println("Error while writting response into file", err)
		return err
	}
	c1 := *content

	file, err := json.MarshalIndent(c1, "", " ")
	if err != nil {
		fmt.Println("Error is: ", err)
	}

	err = ioutil.WriteFile(config.Directory+"response"+t2+".json", file, 0644)
	if err != nil {
		fmt.Println("Error is: ", err)
	}
	return nil
}
