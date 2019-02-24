package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"strconv"
)

// Functions & structs for allTrips
type Address struct {
	Longitude int `json:"longitude"`
	Latitude int `json:"latitude"`
}

type Trip struct {
	Id int `json:"Id"`
	Pickup Address `json:"pickupLocation"`
	Dropoff Address `json:"dropoffLocation"`
	PickupTime int `json:"pickupTime"`
	DriverName string `json:"driverName"`
}

type Trips []Trip

func allTrips(w http.ResponseWriter, r *http.Request) {
	a1 := Address{
		Longitude:1, 
		Latitude:2,
	}
	a2 := Address{
		Longitude:3, 
		Latitude:4,
	}
	
	trips := Trips{
		Trip{
			Id: 100,
			Pickup: a1, 
			Dropoff: a2,
			PickupTime: 500,
			DriverName: "Robbie",
		},
		Trip{
			Id: 101,
			Pickup: a1, 
			Dropoff: a2,
			PickupTime: 500,
			DriverName: "Bobby",
		},
	}

	fmt.Println("Endpoint Hit: All Trips Endpoint")
	json.NewEncoder(w).Encode(trips)
}

// Homepage
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Homepage Endpoint Hit")
}

func returnSingleTrip(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["Id"]

	a1 := Address{
		Longitude:1, 
		Latitude:2,
	}
	a2 := Address{
		Longitude:3, 
		Latitude:4,
	}
	
	trips := Trips{
		Trip{
			Id: 100,
			Pickup: a1, 
			Dropoff: a2,
			PickupTime: 500,
			DriverName: "Robbie",
		},
		Trip{
			Id: 101,
			Pickup: a1, 
			Dropoff: a2,
			PickupTime: 500,
			DriverName: "Bobby",
		},
	}

	// json.NewEncoder(w).Encode()
	for i := 0; i < len(trips); i++ {
		curId, _ := strconv.ParseInt(key, 0, 32)
		if trips[i].Id == int(curId) {
			json.NewEncoder(w).Encode(trips[i])
		}
	}
}

// match URL path hit with a defined function
func handleRequests() {

	nrouter := mux.NewRouter().StrictSlash(true)
	nrouter.HandleFunc("/trips/{Id}", returnSingleTrip)
	nrouter.HandleFunc("/trips", allTrips)
	nrouter.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe(":8081", nrouter))
}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	handleRequests()
}