package main

import (
	"encoding/json"
	"log"
	"regexp"
	"strings"
	"fmt"
	"crypto/rand"
)

func createUniqueTripId() (uniqueId string){
	b := make([]byte, 16)
	_, err:= rand.Read(b)
	if err != nil {
		log.Println("Error: ", err)
	}

	uniqueId = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return
}

func createShitPost(id string, username string, message string, updatedTime string) {
	//couldNotParse := "CAN NOT PARSE, HUMAN VERIFICATION REQUIRED"
	//We will do something so we can manually create these, populate db with raw or output to csv?
}

func createDriver(id string, username string) *Driver {
	var driver Driver
	driver.Name = username
	driver.ProfileURL = ""
	driver.FBContact = ""
	return &driver
}

func createTrip(id string, username string, message string, updatedTime string, post *Post) {
	var trip Trip
	trip.PostInformation = post
	pattern, _ := regexp.Compile("(DRIVING|FROM)(.*)( TO | -> )")
	upperMessage := strings.ToUpper(message)
	pickupSegment := pattern.FindString(upperMessage)
	dropoffSegment := strings.Split(upperMessage, pickupSegment)[1]
	pickupCity := ""
	dropoffCity := ""
	for key, value := range locationsMap {
		if strings.Contains(pickupSegment, key) {
			pickupCity = value
		}
		if strings.Contains(dropoffSegment, key) {
			dropoffCity = value
		}
		if len(pickupCity) > 0 && len(dropoffCity) > 0 {
			break
		}
	}
	trip.TripId = createUniqueTripId()
	pickupLocation := getAddressObject(pickupCity)
	trip.PickupLocation = &pickupLocation
	dropoffLocation := getAddressObject(dropoffCity)
	trip.DropoffLocation = &dropoffLocation
	trip.Driver = createDriver(id, username)
	timeData := getTimeObject(message, updatedTime)
	trip.PickupTime = &timeData
	trip.FBPosting = message
	trip.SpotsAvailable = ""
	trip.Price = ""
	
	persistToDB(trip)
	
	fmt.Printf("Pickup Location\n")
	fmt.Printf("Long: %s\nLat: %s\nCity name: %s\nLong name: %s\n", trip.PickupLocation.Lng, trip.PickupLocation.Lat, trip.PickupLocation.City, trip.PickupLocation.FormattedAddress)	
	fmt.Printf("Dropoff Location\n")
	fmt.Printf("Long: %s\nLat: %s\nCity name: %s\nLong name: %s\n", trip.DropoffLocation.Lng, trip.DropoffLocation.Lat, trip.DropoffLocation.City, trip.DropoffLocation.FormattedAddress)	
	fmt.Printf("Approximate Time Info\n")
	fmt.Printf("TimePrecise: %s\nTimeRange: %s\nDate: %s\n", trip.PickupTime.TimePrecise, trip.PickupTime.TimeRange, trip.PickupTime.Date)

	//populate db
}

func createTripRequest(id string, username string, message string, updatedTime string) {
	//Eventually will do something once we get more trips
}

func parseMessage(id string, username string, message string, updatedTime string, post *Post) {
	pattern, _ := regexp.Compile("(DRIV|LOOK)")
	match := pattern.FindString(strings.ToUpper(message))
	if match == "DRIV" {
		createTrip(id, username, message, updatedTime, post)
	} else if match == "LOOK" {
		createTripRequest(id, username, message, updatedTime)
	} else {
		createShitPost(id, message, username, updatedTime)
	}
}

func parseJson(post []byte) {
	var parsePost Post
	err := json.Unmarshal(post, &parsePost)
	if err != nil {
		log.Fatal("Why does Golang compile error if you dont use variables", err)
	}
	parseMessage(parsePost.Id, parsePost.Username, parsePost.Message, parsePost.UpdatedTime, &parsePost)
}
