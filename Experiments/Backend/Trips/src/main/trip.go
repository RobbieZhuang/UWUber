package main

// Functions & structs for allTrips
// type Address struct {
// 	Longitude int `json:"longitude"`
// 	Latitude int `json:"latitude"`
// }

type Trip struct {
	Id int `json:"Id"`
	// Pickup *Address `json:"pickupLocation"`
	PickupLongitude int `json:"pickupLongitude"`
	PickupLatitude  int `json:"pickupLatitude"`
	// Dropoff *Address `json:"dropoffLocation"`
	DropoffLongitude int    `json:"dropoffLongitude"`
	DropoffLatitude  int    `json:"dropoffLatitude"`
	PickupTime       int    `json:"pickupTime"`
	DriverName       string `json:"driverName"`
}

type Trips []Trip

type TripRequest struct {
	// Pickup *Address `json:"pickupLocation"`
	PickupLongitude int `json:"pickupLongitude"`
	PickupLatitude  int `json:"pickupLatitude"`
	// Dropoff *Address `json:"dropoffLocation"`
	DropoffLongitude int `json:"dropoffLongitude"`
	DropoffLatitude  int `json:"dropoffLatitude"`
	PickupTime       int `json:"pickupTime"`
}

// var a1 = Address{
// 	Longitude:1,
// 	Latitude:2,
// }

// var a2 = Address{
// 	Longitude:3,
// 	Latitude:4,
// }
