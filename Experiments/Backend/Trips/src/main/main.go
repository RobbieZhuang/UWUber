package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// Homepage
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Homepage Endpoint Hit")
}

func handleRequests() {

	nrouter := mux.NewRouter().StrictSlash(true)

	// Always put in order of priority to be hit
	nrouter.HandleFunc("/trips/user/{UserName}", getTripForUser).Methods("GET")
	nrouter.HandleFunc("/trips/search", getTripWithFilter).Methods("POST")
	nrouter.HandleFunc("/trips/{Id}", getTrip).Methods("GET")
	nrouter.HandleFunc("/trips/{Id}", createTrip).Methods("POST")
	nrouter.HandleFunc("/trips/{Id}", deleteTrip).Methods("DELETE")
	nrouter.HandleFunc("/trips/update/{Id}", updateTrip).Methods("POST")
	nrouter.HandleFunc("/trips", allTrips).Methods("GET")
	nrouter.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe(":8081", nrouter))
}

func main() {
	fmt.Println("Rest API v1.0 - With Mux Routers | Big Car Tings")
	handleRequests()
}
