package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"

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

var monthMap = map[string]string{
	"JAN":  "01",
	"FEB":  "02",
	"MAR":  "03",
	"APR":  "04",
	"MAY":  "05",
	"JUN":  "06",
	"JUL":  "07",
	"AUG":  "08",
	"SEPT": "09",
	"OCT":  "10",
	"NOV":  "11",
	"DEC":  "12"}

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
	TripID          string       `json:"TripID"`
	PostInformation *Post        `json:"postInformation"`
	PickupLocation  *AddressData `json:"pickupLocation"`
	DropoffLocation *AddressData `json:"dropoffLocation"`
	PickupTime      *TimeData    `json:"pickupTime"`
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
	TimePrecise string `json:"timePrecise"`
	TimeRange   string `json:"timeRange"`
	Date        string `json:"date"`
}

func getTimeObject(message string, postTime string) TimeData {
	// this code is garbage sorry lmao
	couldNotParse := "CAN NOT PARSE, HUMAN VERIFICATION REQUIRED"

	currDate := postTime[0:10]
	upperMessage := strings.ToUpper(message)
	dateFromPost := ""
	timePrecise := ""
	timeDescription := ""
	for key, value := range monthMap {
		if strings.Contains(upperMessage, key) {
			pattern, _ := regexp.Compile("[0-9]+")
			day := pattern.FindString(upperMessage[strings.Index(upperMessage, key):])
			if len(day) != 2 {
				day = "0" + day
			}
			dateFromPost = currDate[0:4] + "-" + value + "-" + day
		}
	}

	pattern, _ := regexp.Compile("[0-9]+.*?(AM|PM)")
	allFoundTime := pattern.FindAllString(upperMessage, 2)
	if len(allFoundTime) == 2 {
		timeDescription = allFoundTime[0] + " to " + allFoundTime[1]
	} else if len(allFoundTime) == 1 {
		timePrecise = allFoundTime[0]
	}

	if len(dateFromPost) == 0 {
		dateFromPost = couldNotParse
	}
	if len(timePrecise) == 0 {
		timePrecise = couldNotParse
	}
	if len(timeDescription) == 0 {
		timeDescription = couldNotParse
	}

	return TimeData{
		TimePrecise: timePrecise,
		TimeRange:   timeDescription,
		Date:        dateFromPost,
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
