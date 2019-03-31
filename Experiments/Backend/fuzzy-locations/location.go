package main

import (
	"encoding/json"
	"fmt"
	"net/url"

	"gopkg.in/resty.v1"
)

const token = "AIzaSyDe9KBNpY2cZ8ghI-hTNcRoXHOVDYqQdvA"

type LocationResponse struct {
	Predictions     []*Predictions `json:"predictions"`
	PickupLongitude int            `json:"pickupLongitude"`
}

type Predictions struct {
	Description string `json:"description"`
	ID          string `json:"id"`
	PlaceID     string `json:"place_id"`
}

type ResultResult struct {
	Result *PlaceResponseResult `json:"result"`
}

type PlaceResponseResult struct {
	AddressComponents []*AddressComponents `json:"address_components"`
	Geometry          *Geometry            `json:"geometry"`
}

type AddressComponents struct {
	LongName  string   `json:"long_name"`
	ShortName string   `json:"short_name"`
	Types     []string `json:"types"`
}

type Geometry struct {
	Location *LongLat `json:"location"`
}

type LongLat struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Time struct {
	ExactTime  string `json:"exactTime"`
	TimeString string `json:"timeString"`
	Date       string `json:"date"`
}

type Driver struct {
	Name       string `json:"name"`
	ProfileURL string `json:"profileURL"`
	FBContact  string `json:"fbContact"`
}

type Trip struct {
	PostInformation *Post        `json:"postInformation"`
	PickupLocation  *AddressData `json:"pickupLocation"`
	DropoffLocation *AddressData `json:"dropoffLocation"`
	PickupTime      *Time        `json:"pickupTime"`
	Driver          *Driver      `json:"driver"`
	FBPosting       string       `json:"fbPosting"`
	SpotsAvailable  string       `json:"spotsAvailable"`
	Price           string       `json:"price"`
}

// Sauga/Sq1 works on google maps, but not our places endpoint? Try Waterloo BK
// https://maps.googleapis.com/maps/api/place/findplacefromtext/json?input=Sauga+Sq1&key=AIzaSyDe9KBNpY2cZ8ghI-hTNcRoXHOVDYqQdvA

func getValidLocation(str string) (id string) {
	resp, _ := resty.R().
		SetQueryParams(map[string]string{
			"key":   token,
			"input": url.QueryEscape(str + " Canada"),
		}).
		Get("https://maps.googleapis.com/maps/api/place/queryautocomplete/json")

	var locResp LocationResponse

	_ = json.Unmarshal(resp.Body(), &locResp)

	if len(locResp.Predictions) > 0 {
		return locResp.Predictions[0].PlaceID
	} else {
		fmt.Printf("\nno predictions found")
		return ""
	}
}

func getPlaceDetails(placeID string) (long float64, lat float64, city string) {
	resp, _ := resty.R().SetQueryParams(map[string]string{
		"placeid": placeID,
		"key":     token,
	}).Get("https://maps.googleapis.com/maps/api/place/details/json")

	var coords ResultResult

	_ = json.Unmarshal(resp.Body(), &coords)

	var countryName = ""

	for _, addressComp := range coords.Result.AddressComponents {
		if contains(addressComp.Types, "locality") || contains(addressComp.Types, "political") {
			countryName = addressComp.ShortName
			break
		}
	}
	return coords.Result.Geometry.Location.Lng, coords.Result.Geometry.Location.Lat, countryName
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

type AddressData struct {
	Lat  string `json:"lat"`
	Lng  string `json:"lng"`
	City string `json:"city"`
}

// Call this method to get AddressData object back!
func getAddress(s string) AddressData {
	couldNotParse := "CAN NOT PARSE, HUMAN VERIFICATION REQUIRED"
	locID := getValidLocation(s)
	if locID == "" {
		return AddressData{Lat: couldNotParse, Lng: couldNotParse, City: couldNotParse}
	}
	long, lat, cityName := getPlaceDetails(locID)
	return AddressData{Lat: fmt.Sprintf("%f", lat), Lng: fmt.Sprintf("%f", long), City: cityName}
}

type TimeData struct {
	TimePrecise     string `json:"timePrecise"`
	TimeDescription string `json:"timeDescription"`
	Date            string `json:"date"`
}

func main() {
	sampleJson1 := []byte(`{"id":"id1", "username":"Brendan Zhang", "message":"Looking for ride to Union/Finch from Waterloo bk on Sunday (10th) after 4pm.", "updatedTime" : "2018-02-31T05:33:31+0000"}`)
	sampleJson2 := []byte(`{"id":"id2", "username":"Daniell Yang", "message":"Looking for a ride from Brampton to Waterloo on 10th March (Sunday).", "updatedTime":"2018-02-31T05:33:31+0000"}`)
	sampleJson3 := []byte(`{"id":"id3", "username":"Bimesh DeSilva", "message":"Driving London -> Waterloo @ 1 pm on Sunday March 10th, $20", "updatedTime" : "2018-02-31T05:33:31+0000"}`)
	sampleJson4 := []byte(`{"id":"id4", "username":"Max Gao", "message":"driving richmond hill freshco plaza to waterloo bk plaza at 1pm sunday march 10, no middle seat, taking 407, $20 a seat", "updatedTime" : "2018-02-31T05:33:31+0000"}`)
	shitpost := []byte(`{"id":"id5", "username":"shitposter", "message":"Shitpost", "updatedTime" : "2018-02-31T05:33:31+0000"}`)

	fmt.Println("FUZZILY SEARCH - Get lat/long & city name for fuzzily searched location")

	locName := "Waterloo BK"

	locID := getValidLocation(locName)
	long, lat, cityName := getPlaceDetails(locID)

	fmt.Printf("Long: %f\nLat: %f\nCity Name: %s\n", long, lat, cityName)

	ad := getAddress(locName)
	fmt.Printf("\nLong: %s\nLat: %s\nCity name: %s\n", ad.Lng, ad.Lat, ad.City)
}
