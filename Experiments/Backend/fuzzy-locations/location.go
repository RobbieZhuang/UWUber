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

func main() {
	fmt.Println("FUZZILY SEARCH - Get lat/long & city name for fuzzily searched location")

	locName := "Waterloo BK"

	locID := getValidLocation(locName)
	long, lat, cityName := getPlaceDetails(locID)

	fmt.Printf("Long: %f\nLat: %f\nCity Name: %s\n", long, lat, cityName)
}
