package main

import (
	"fmt"
	"log"
	"strconv"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/google/go-cmp/cmp"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// Functions & structs for allTrips
// type Address struct {
// 	Longitude int `json:"longitude"`
// 	Latitude int `json:"latitude"`
// }

type Trip struct {
	Id int `json:"Id"`
	// Pickup *Address `json:"pickupLocation"`
	PickupLongitude int `json:"pickupLongitude"`
	PickupLatitude int `json:"pickupLatitude"`
	// Dropoff *Address `json:"dropoffLocation"`
	DropoffLongitude int `json:"dropoffLongitude"`
	DropoffLatitude int `json:"dropoffLatitude"`
	PickupTime int `json:"pickupTime"`
	DriverName string `json:"driverName"`
}

type Trips []Trip

type TripRequest struct {
	// Pickup *Address `json:"pickupLocation"`
	PickupLongitude int `json:"pickupLongitude"`
	PickupLatitude int `json:"pickupLatitude"`
	// Dropoff *Address `json:"dropoffLocation"`
	DropoffLongitude int `json:"dropoffLongitude"`
	DropoffLatitude int `json:"dropoffLatitude"`
	PickupTime int `json:"pickupTime"`
}

// var a1 = Address{
// 	Longitude:1, 
// 	Latitude:2,
// }

// var a2 = Address{
// 	Longitude:3, 
// 	Latitude:4,
// }

var trips = Trips{
	Trip{
		Id: 100,
		PickupLongitude: 1,
		PickupLatitude: 2,
		DropoffLongitude: 3,
		DropoffLatitude: 4,  
		// Pickup: &a1, 
		// Dropoff: &a2,
		PickupTime: 500,
		DriverName: "Robbie",
	},
	Trip{
		Id: 101,
		PickupLongitude: 1,
		PickupLatitude: 2,
		DropoffLongitude: 3,
		DropoffLatitude: 4,  
		// Pickup: &a1, 
		// Dropoff: &a2,
		PickupTime: 500,
		DriverName: "Bobby",
	},
}

func allTrips(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: All Trips Endpoint")
	allTrips := getAllTripsFromDB()
	json.NewEncoder(w).Encode(allTrips)
}

// Homepage
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Homepage Endpoint Hit")
}

func getTrip(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["Id"]
	trip, err := getTripForIDFromDB(getIntFor(key))
	if err != nil {
		log.Fatal("TripID does not exist", err)
	}
	
	json.NewEncoder(w).Encode(trip)

	// for _, item := range trips {
	// 	if item.Id == getIntFor(key) {
	// 		json.NewEncoder(w).Encode(item)
	// 	}
	// }
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

	for _, item := range trips {
		if cmp.Equal(item.PickupLatitude, tripRequest.PickupLatitude) && 
			cmp.Equal(item.PickupLongitude, tripRequest.PickupLongitude) && 
			cmp.Equal(item.DropoffLatitude, tripRequest.DropoffLatitude) && 
			cmp.Equal(item.DropoffLongitude, tripRequest.DropoffLongitude) && 
			item.PickupTime == tripRequest.PickupTime {
			json.NewEncoder(w).Encode(item)
		}
	}
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

	allTrips := getAllTripsFromDB()
	json.NewEncoder(w).Encode(allTrips)
}


func deleteTrip(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var newTrips Trips
	for _, item := range trips {
		if item.Id != getIntFor(params["Id"]) {
			newTrips = append(newTrips, item)
		}
	}
	trips = newTrips
	json.NewEncoder(w).Encode(trips)
}

func getIntFor(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}
	return i
}


// DB
func StartDB() *sql.DB {
	db, err := sql.Open("mysql",
			"root:root@tcp(localhost:8889)/test_db")
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func getTripForIDFromDB(/*db *sql.DB,*/ tripId int) (t Trip, e error) {

	// Start db acess - use a closure or some better way to deal with opening and closing the db
	db, err := sql.Open("mysql", "root:root@tcp(localhost:8889)/test_db")
	if err != nil {
		log.Fatal(err)
	}

	var (
		id int
		pickupLongitude int
		pickupLatitude int 
		dropoffLongitude int 
		dropoffLatitude int 
		pickupTime int
		driverName string
	)

	// Create a reusable query
	stmt, err := db.Prepare("select id, pickupLongitude, pickupLatitude, dropoffLongitude, dropoffLatitude, pickupTime, driverName from trips where id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(tripId).Scan(&id, &pickupLongitude, &pickupLatitude, &dropoffLongitude, &dropoffLatitude, &pickupTime, &driverName)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(id, pickupLongitude, pickupLatitude, dropoffLongitude, dropoffLatitude, pickupTime, driverName)

	db.Close()

	var trip = Trip{
		Id: id,
		PickupLongitude: pickupLongitude,
		PickupLatitude: pickupLatitude,
		DropoffLongitude: dropoffLongitude,
		DropoffLatitude: dropoffLatitude,
		PickupTime: pickupTime,
		DriverName: driverName,
	}
	fmt.Println(trip)

	return trip, err
}

func getAllTripsFromDB(/*db *sql.DB,*/) Trips {

	// Start db acess - use a closure or some better way to deal with opening and closing the db
	db, err := sql.Open("mysql", "root:root@tcp(localhost:8889)/test_db")
	if err != nil {
		log.Fatal(err)
	}

	var (
		id int
		pickupLongitude int
		pickupLatitude int 
		dropoffLongitude int 
		dropoffLatitude int 
		pickupTime int
		driverName string
	)

	// Create a reusable query
	rows, err := db.Query("select id, pickupLongitude, pickupLatitude, dropoffLongitude, dropoffLatitude, pickupTime, driverName from trips")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var trips Trips
	for rows.Next() {
		err = rows.Scan(&id, &pickupLongitude, &pickupLatitude, &dropoffLongitude, &dropoffLatitude, &pickupTime, &driverName)
		if err != nil {
			log.Fatal(err)
		}
		var trip = Trip{
			Id: id,
			PickupLongitude: pickupLongitude,
			PickupLatitude: pickupLatitude,
			DropoffLongitude: dropoffLongitude,
			DropoffLatitude: dropoffLatitude,
			PickupTime: pickupTime,
			DriverName: driverName,
		}
		trips = append(trips, trip)
	}
	db.Close()

	return trips
}

func addTripToDB(/*db *sql.DB,*/ trip Trip) int64 {

	// Start db acess - use a closure or some better way to deal with opening and closing the db
	db, err := sql.Open("mysql", "root:root@tcp(localhost:8889)/test_db")
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare("INSERT INTO trips (pickupLongitude, pickupLatitude, dropoffLongitude, dropoffLatitude, pickupTime, driverName) VALUES(?,?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(&trip.PickupLongitude, &trip.PickupLatitude, &trip.DropoffLongitude, &trip.DropoffLatitude, &trip.PickupTime, &trip.DriverName)
	if err != nil {
		log.Fatal(err)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)

	db.Close()

	return lastId
}

func closeDB(db *sql.DB) {
	defer db.Close()
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
	// db := StartDB()
	// closeDB(db)

	fmt.Println("Rest API v1.0 - With Mux Routers | Big Car Tings")
	handleRequests()
}