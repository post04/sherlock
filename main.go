package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/Schmenn/sherlock/requests"
	"github.com/fatih/color"
)

// SitesStruct is what's seen in the sites.json file
type SitesStruct struct {
	URL         string      `json:"url"`
	PrintURL    string      `json:"printURL"`
	Method      string      `json:"method"`
	Check       string      `json:"check"`
	Body        string      `json:"body"`
	Match       interface{} `json:"match"`
	StructMatch interface{} `json:"structMatch"`
	Headers     [][]string  `json:"headers"`
}

func main() {

	if len(os.Args) < 2 {
		color.HiRed("Please specify a username")
		os.Exit(1)
	}
	name := os.Args[1]
	f, err := os.ReadFile("sites.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	var sites []SitesStruct
	err = json.Unmarshal(f, &sites)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
	var wg sync.WaitGroup
	for _, site := range sites {
		wg.Add(1)
		go func(site SitesStruct) {
			defer wg.Done()
			req := requests.NewRequest(site.Method, site.URL, name, site.Body, site.Headers)
			if site.PrintURL != "" {
				site.URL = site.PrintURL
			}
			if req == nil {
				color.HiRed("[ ] " + strings.Split(strings.Split(site.URL, "//")[1], "/")[0])
				return
			}
			resp := requests.MakeRequest(req, &http.Client{})
			if resp == nil {
				color.HiRed("[ ] " + strings.Split(strings.Split(site.URL, "//")[1], "/")[0])
				return
			}
			if requests.CheckValidity(resp, site.Check, site.Match, site.StructMatch) {
				color.HiGreen(fmt.Sprintf("[X] %s", strings.Replace(site.URL, "{}", name, 1)))
			} else {
				color.HiRed("[ ] " + strings.Split(strings.Split(site.URL, "//")[1], "/")[0])
			}
		}(site)
	}
	wg.Wait()
}
