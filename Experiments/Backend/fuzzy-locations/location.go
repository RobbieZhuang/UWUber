package main

import (
	"encoding/json"
	"fmt"
	"net/url"

	"gopkg.in/resty.v1"
)

const token = "AIzaSyDe9KBNpY2cZ8ghI-hTNcRoXHOVDYqQdvA"

var locationsMap = map[string]string{
	"TDOT":          "Toronto",
	"TORONTO":       "Toronto",
	"RH":            "Richmond Hill",
	"RICHMOND HILL": "Richmond Hill",
	"MARKHAM":       "Richmond Hill",
	"LOO":           "Waterloo",
	"WATERLOO":      "Waterloo",
	"BK PLAZA":      "Waterloo",
	"BURGER KING":   "Waterloo",
	"UW":            "Waterloo",
	"STC":           "Scarborough",
	"SCARBOROUGH":   "Scarborough",
	"BRAMPTON":      "Brampton",
	"LONDON":        "London",
	"WESTERN":       "London",
	"PM":            "Pacific Mall",
	"PACIFIC MALL":  "Pacific Mall",
	"FINCH":         "Finch Station"}

type Time struct {
	TimePrecise  string `json:"timePrecise"`
	TimeDescription string `json:"timeDescription"`
	Date       string `json:"date"`
}

type Driver struct {
	Name       string `json:"name"`
	ProfileURL string `json:"profileURL"`
	FBContact  string `json:"fbContact"`
}

type Post struct {
	Id          string `json:"id"`
	Username    string `json:”username”`
	Message     string `json:"message"`
	UpdatedTime string `json:"updatedTime"`
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

// MARK: Location Parsing

// Google Places API Struct
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
	FormattedAddress  string               `json:"formatted_address"`
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

// Location Parsing
type Error struct {
	HasErr bool
	ErrMsg string
}

type LocationID struct {
	Id  string
	Err Error
}

type AddressDataResponse struct {
	Lat      float64
	Lng      float64
	City     string
	LongName string
	Err      Error
}

func getLocationID(str string) (locationID LocationID) {

	// Sauga/Sq1 works on google maps, but not our places endpoint? Try Waterloo BK
	// https://maps.googleapis.com/maps/api/place/findplacefromtext/json?input=Sauga+Sq1&key=AIzaSyDe9KBNpY2cZ8ghI-hTNcRoXHOVDYqQdvA

	resp, err := resty.R().
		SetQueryParams(map[string]string{
			"key":   token,
			"input": url.QueryEscape(str + " Canada"),
		}).
		Get("https://maps.googleapis.com/maps/api/place/queryautocomplete/json")

	if err != nil {
		return LocationID{
			Id: "",
			Err: Error{
				HasErr: true,
				ErrMsg: "Could not get ID for place.",
			},
		}
	}

	var locResp LocationResponse

	_ = json.Unmarshal(resp.Body(), &locResp)

	if len(locResp.Predictions) > 0 {
		return LocationID{
			Id: locResp.Predictions[0].PlaceID,
			Err: Error{
				HasErr: false,
				ErrMsg: "",
			},
		}
	}
	return LocationID{
		Id: "",
		Err: Error{
			HasErr: true,
			ErrMsg: "Could not parse location ID.",
		},
	}
}

func getPlaceDetails(locationID LocationID) (addressDataResponse AddressDataResponse) {
	if locationID.Err.HasErr {
		return AddressDataResponse{
			Lat:      0,
			Lng:      0,
			City:     "",
			LongName: "",
			Err: Error{
				HasErr: true,
				ErrMsg: "Invalid placeID, could not get address data.",
			},
		}
	}

	resp, _ := resty.R().SetQueryParams(map[string]string{
		"placeid": locationID.Id,
		"key":     token,
	}).Get("https://maps.googleapis.com/maps/api/place/details/json")

	var placeDetails ResultResult

	err := json.Unmarshal(resp.Body(), &placeDetails)

	if err != nil {
		return AddressDataResponse{
			Lat:      0,
			Lng:      0,
			City:     "",
			LongName: "",
			Err: Error{
				HasErr: true,
				ErrMsg: "Invalid placeID, could not get address data.",
			},
		}
	}
	var countryName = ""

	for _, addressComp := range placeDetails.Result.AddressComponents {
		if contains(addressComp.Types, "locality") || contains(addressComp.Types, "political") {
			countryName = addressComp.ShortName
			break
		}
	}

	return AddressDataResponse{
		Lat:      placeDetails.Result.Geometry.Location.Lat,
		Lng:      placeDetails.Result.Geometry.Location.Lng,
		City:     countryName,
		LongName: placeDetails.Result.FormattedAddress,
		Err: Error{
			HasErr: false,
			ErrMsg: "",
		},
	}
}

type AddressData struct {
	Lat              string `json:"lat"`
	Lng              string `json:"lng"`
	City             string `json:"city"`
	FormattedAddress string `json:"formattedAddress"`
}

// Call this method to get AddressData object back!
func getAddressObject(s string) AddressData {
	couldNotParse := "CAN NOT PARSE, HUMAN VERIFICATION REQUIRED"
	locID := getLocationID(s)
	addressData := getPlaceDetails(locID)
	if locID.Err.HasErr || addressData.Err.HasErr {
		return AddressData{
			Lat:              couldNotParse,
			Lng:              couldNotParse,
			City:             couldNotParse,
			FormattedAddress: couldNotParse,
		}
	}
	return AddressData{
		Lat:              fmt.Sprintf("%f", addressData.Lat),
		Lng:              fmt.Sprintf("%f", addressData.Lng),
		City:             addressData.City,
		FormattedAddress: addressData.LongName,
	}
}

// MARK: Time Parsing

type TimeData struct {
	TimePrecise     string `json:"timePrecise"`
	TimeDescription string `json:"timeDescription"`
	Date            string `json:"date"`
}

func getTimeObject(s string, postTime string) TimeData {
	couldNotParse := "CAN NOT PARSE, HUMAN VERIFICATION REQUIRED"
	return TimeData{
		TimePrecise:     couldNotParse,
		TimeDescription: couldNotParse,
		Date:            couldNotParse,
	}
}

// MARK: Helper Functions

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func main() {
	//sampleJson1 := []byte(`{"id":"id1", "username":"Brendan Zhang", "message":"Looking for ride to Union/Finch from Waterloo bk on Sunday (10th) after 4pm.", "updatedTime" : "2018-02-31T05:33:31+0000"}`)
	//sampleJson2 := []byte(`{"id":"id2", "username":"Daniell Yang", "message":"Looking for a ride from Brampton to Waterloo on 10th March (Sunday).", "updatedTime":"2018-02-31T05:33:31+0000"}`)
	 sampleJson3 := []byte(`{"id":"id3", "username":"Bimesh DeSilva", "message":"Driving London -> Waterloo @ 1 pm on Sunday March 10th, $20", "updatedTime" : "2018-02-31T05:33:31+0000"}`)
	//sampleJson4 := []byte(`{"id":"id4", "username":"Max Gao", "message":"driving richmond hill freshco plaza to waterloo bk plaza at 1pm sunday march 10, no middle seat, taking 407, $20 a seat", "updatedTime" : "2018-02-31T05:33:31+0000"}`)
	//shitpost := []byte(`{"id":"id5", "username":"shitposter", "message":"Shitpost", "updatedTime" : "2018-02-31T05:33:31+0000"}`)

	locName := "Waterloo BK"

	time := "This afternoon"
	timeContext := "04/13/2019"

	ad := getAddressObject(locName)
	fmt.Printf("Long: %s\nLat: %s\nCity name: %s\nLong name: %s\n", ad.Lng, ad.Lat, ad.City, ad.FormattedAddress)

	fmt.Println()

	tm := getTimeObject(time, timeContext)
	fmt.Printf("TimePrecise: %s\nTimeDescription: %s\nDate: %s\n", tm.TimePrecise, tm.TimeDescription, tm.Date)

	 parseJson(sampleJson3)
}
