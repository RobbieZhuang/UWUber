package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func allTrips(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: All Trips Endpoint")
	allTrips, err := getAllTripsFromDB()
	if err != nil {
		log.Fatal("Could not get all trips", err)
	}
	json.NewEncoder(w).Encode(allTrips)
}

func getTrip(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["Id"]
	trip, err := getTripForIDFromDB(getIntFor(key))
	if err != nil {
		log.Fatal("TripID does not exist", err)
	}

	json.NewEncoder(w).Encode(trip)
}

func getTripForUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	value := vars["UserName"]

	trips, err := getAllTripsForUserFromDB(value)

	if err != nil {
		log.Fatal("TripID does not exist", err)
	}

	json.NewEncoder(w).Encode(trips)
}

/*
To search for a specific trip

POST to http://localhost:8081/trips/search

{
	"pickupLongitude":1,
	"pickupLatitude":2,
	"dropoffLongitude":3,
	"dropoffLatitude":4,
	"pickupTime":500
}

*/
func getTripWithFilter(w http.ResponseWriter, r *http.Request) {
	var tripRequest TripRequest
	_ = json.NewDecoder(r.Body).Decode(&tripRequest)

	trips, err := getAllTripsWithLocationAndTimeFromDB(tripRequest.PickupLongitude, tripRequest.PickupLatitude, tripRequest.DropoffLongitude, tripRequest.DropoffLatitude, tripRequest.PickupTime)
	if err != nil {
		log.Fatal("Could not find any trips with those specific inputs", err)
	}
	json.NewEncoder(w).Encode(trips)
}

/*
To create a trip, download postman

POST to http://localhost:8081/trips/{Id}

Put below in body
{
	"Id":101,
	"pickupLongitude":1,
	"pickupLatitude":2,
	"dropoffLongitude":3,
	"dropoffLatitude":4,
	"pickupTime":500,
	"driverName":
	"Bobby"
}

*/
func createTrip(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var trip Trip
	_ = json.NewDecoder(r.Body).Decode(&trip)

	trip.Id = getIntFor(params["Id"])

	addTripToDB(trip)

	allTrips, err := getAllTripsFromDB()
	if err != nil {
		log.Fatal("TripID does not exist", err)
	}
	json.NewEncoder(w).Encode(allTrips)
}

func updateTrip(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["Id"]

	var trip Trip
	_ = json.NewDecoder(r.Body).Decode(&trip)

	updatedId, err := updateTripInDB(getIntFor(key), trip)
	if err != nil {
		log.Fatal("Could not update trip with tripId", key, err)
	}
	updatedTrip, _ := getTripForIDFromDB(updatedId)
	json.NewEncoder(w).Encode(updatedTrip)
}

func deleteTrip(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	trips, err := deleteTripFromDB(getIntFor(params["Id"]))
	if err != nil {
		log.Fatal("TripID does not exist", err)
	}
	json.NewEncoder(w).Encode(trips)
}
