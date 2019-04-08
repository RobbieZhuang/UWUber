// Structs copied from Robbie's Experiments/Backend/Trips/src/main/trip.go
// I added some new variables for matching trips


package main

import "time"


type Trip struct {
	Id int `json:"Id"`
	PickupLongitude int `json:"pickupLongitude"`
	PickupLatitude  int `json:"pickupLatitude"`
	DropoffLongitude  int    `json:"dropoffLongitude"`
	DropoffLatitude   int    `json:"dropoffLatitude"`
	DriverName        string `json:"driverName"`
    PickupLocation    string  // Added
    DropoffLocation   string  // Added
	SpotsAvailable    int  // Added
	PickupTimeData TimeData  // Added
}

type Trips []Trip

type TripRequest struct {
	Id int `json:"Id"` // Added this
	PickupLongitude int `json:"pickupLongitude"`
	PickupLatitude  int `json:"pickupLatitude"`
	DropoffLongitude int   `json:"dropoffLongitude"`
	DropoffLatitude  int   `json:"dropoffLatitude"`
    PickupLocation    string  // Added
    DropoffLocation   string  // Added
	PickupTimeData TimeData  // Added
}

type TripRequests []TripRequests

// Return to db
type AddressData struct {
    Lat float64 `json:"lat"`
    Lng float64 `json:"lng"`
    City string `json:"city"`
}

// Return to db
type TimeData struct {
    TimePrecise string `json:"timePrecise"`
    TimeDescription string `json:"timeDescription"`
    Date string `json:"date"`
	Time time.Time  // Added
}
