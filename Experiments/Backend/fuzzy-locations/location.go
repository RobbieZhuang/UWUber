package main

import (
	"encoding/json"
	"fmt"
	"net/url"

	"gopkg.in/resty.v1"
)

const token = "[INSERT API KEY HERE]"

type Predictions struct {
	Description string `json:"description"`
	ID          string `json:"id"`
}

type LocationResponse struct {
	Predictions     []*Predictions `json:"predictions"`
	PickupLongitude int            `json:"pickupLongitude"`
}

// Sauga/Sq1 works on google maps, but not our places endpoint? Try Waterloo BK
// https://maps.googleapis.com/maps/api/place/findplacefromtext/json?input=Sauga+Sq1&key=AIzaSyDe9KBNpY2cZ8ghI-hTNcRoXHOVDYqQdvA

func getValidLocation(str string) {
	resp, _ := resty.R().
		SetQueryParams(map[string]string{
			"key":   token,
			"input": url.QueryEscape(str + " Canada"),
		}).
		Get("https://maps.googleapis.com/maps/api/place/queryautocomplete/json")

	var locResp LocationResponse

	_ = json.Unmarshal(resp.Body(), &locResp)

	fmt.Printf(url.QueryEscape(str))
	fmt.Printf("\nResponse Body: %v", resp)

	if len(locResp.Predictions) > 0 {
		fmt.Printf("\n%s", locResp.Predictions[0].Description)
	} else {
		fmt.Printf("\nno predictions found")
	}
}

func main() {
	fmt.Println("Rest API v1.0 - With Mux Routers | Big Car Tings")
	getValidLocation("Square One")
}
