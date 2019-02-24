package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/google/go-cmp/cmp"
)

// Functions & structs for allTrips
type Address struct {
	Longitude int `json:"longitude"`
	Latitude int `json:"latitude"`
}

type Trip struct {
	Id string `json:"Id"`
	Pickup *Address `json:"pickupLocation"`
	Dropoff *Address `json:"dropoffLocation"`
	PickupTime int `json:"pickupTime"`
	DriverName string `json:"driverName"`
}

type Trips []Trip

type TripRequest struct {
	Pickup *Address `json:"pickupLocation"`
	Dropoff *Address `json:"dropoffLocation"`
	PickupTime int `json:"pickupTime"`
}

var a1 = Address{
	Longitude:1, 
	Latitude:2,
}

var a2 = Address{
	Longitude:3, 
	Latitude:4,
}

var trips = Trips{
	Trip{
		Id: "100",
		Pickup: &a1, 
		Dropoff: &a2,
		PickupTime: 500,
		DriverName: "Robbie",
	},
	Trip{
		Id: "101",
		Pickup: &a1, 
		Dropoff: &a2,
		PickupTime: 500,
		DriverName: "Bobby",
	},
}

func allTrips(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: All Trips Endpoint")
	json.NewEncoder(w).Encode(trips)
}

// Homepage
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Homepage Endpoint Hit")
}

func getTrip(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["Id"]

	for _, item := range trips {
		if item.Id == key {
			json.NewEncoder(w).Encode(item)
		}
	}
}

func getTripForUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["UserName"]

	for _, item := range trips {
		if item.DriverName == key {
			json.NewEncoder(w).Encode(item)
		}
	}
}

/*
To search for a specific trip

POST to http://localhost:8081/trips/search
{
"pickupLocation":{"longitude":1,"latitude":2},
"dropoffLocation":{"longitude":3,"latitude":4},
"pickupTime":500
}
*/
func getTripWithFilter(w http.ResponseWriter, r *http.Request) {
	var tripRequest TripRequest
	_ = json.NewDecoder(r.Body).Decode(&tripRequest)

	for _, item := range trips {
		if cmp.Equal(item.Pickup, tripRequest.Pickup) && cmp.Equal(item.Dropoff, tripRequest.Dropoff) && item.PickupTime == tripRequest.PickupTime {
			json.NewEncoder(w).Encode(item)
		}
	}
}

/*
To create a trip, download postman

POST to http://localhost:8081/trips/{Id}

Put below in body
{
"Id": "100",
"pickupLocation": {
"longitude": 1,
"latitude": 2
},
"dropoffLocation": {
"longitude": 3,
"latitude": 4
},
"pickupTime": 500,
"driverName": "Big Boi"
},

*/
func createTrip(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var trip Trip
	_ = json.NewDecoder(r.Body).Decode(&trip)
	trip.Id = params["Id"]

	trips = append(trips, trip)
	json.NewEncoder(w).Encode(trips)
}


func deleteTrip(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var newTrips Trips
	for _, item := range trips {
		if item.Id != params["Id"] {
			newTrips = append(newTrips, item)
		}
	}
	trips = newTrips
	json.NewEncoder(w).Encode(trips)
}


func handleRequests() {

	nrouter := mux.NewRouter().StrictSlash(true)	

	// Always put in order of priority to be hit
	nrouter.HandleFunc("/trips/user/{UserName}", getTripForUser).Methods("GET")
	nrouter.HandleFunc("/trips/search", getTripWithFilter).Methods("POST")
	nrouter.HandleFunc("/trips/{Id}", getTrip).Methods("GET")
	nrouter.HandleFunc("/trips/{Id}", createTrip).Methods("POST")
	nrouter.HandleFunc("/trips/{Id}", deleteTrip).Methods("DELETE")
	nrouter.HandleFunc("/trips", allTrips).Methods("GET")
	nrouter.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe(":8081", nrouter))
}

func main() {
	fmt.Println("Rest API v1.0 - With Mux Routers | Big Car Tings")
	handleRequests()
}