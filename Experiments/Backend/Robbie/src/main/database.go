package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// clean up repetitive db
// - Action (Query, Delete, Modify)
// - Query String to use
// - Multiple/Single
// - Input Value
// - Return Object(s) single or []

type DBAction string

const (
	Query  DBAction = ""
	Delete DBAction = ""
	Modify DBAction = ""
)

func StartDB() *sql.DB {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:8889)/test_db")
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func getTripForIDFromDB( /*db *sql.DB,*/ tripId int) (t Trip, e error) {

	// Start db acess - use a closure or some better way to deal with opening and closing the db
	db := StartDB()

	var (
		id               int
		pickupLongitude  int
		pickupLatitude   int
		dropoffLongitude int
		dropoffLatitude  int
		pickupTime       int
		driverName       string
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
		Id:               id,
		PickupLongitude:  pickupLongitude,
		PickupLatitude:   pickupLatitude,
		DropoffLongitude: dropoffLongitude,
		DropoffLatitude:  dropoffLatitude,
		PickupTime:       pickupTime,
		DriverName:       driverName,
	}
	fmt.Println(trip)

	return trip, err
}

func getAllTripsFromDB( /*db *sql.DB,*/ ) (t Trips, e error) {

	// Start db acess - use a closure or some better way to deal with opening and closing the db
	db := StartDB()

	var (
		id               int
		pickupLongitude  int
		pickupLatitude   int
		dropoffLongitude int
		dropoffLatitude  int
		pickupTime       int
		driverName       string
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
			Id:               id,
			PickupLongitude:  pickupLongitude,
			PickupLatitude:   pickupLatitude,
			DropoffLongitude: dropoffLongitude,
			DropoffLatitude:  dropoffLatitude,
			PickupTime:       pickupTime,
			DriverName:       driverName,
		}
		trips = append(trips, trip)
	}
	db.Close()

	return trips, err
}

func getAllTripsWithLocationAndTimeFromDB(aLong int, aLat int, bLong int, bLat int, pTime int) (t Trips, e error) {

	// Start db acess - use a closure or some better way to deal with opening and closing the db
	db := StartDB()

	var (
		id               int
		pickupLongitude  int
		pickupLatitude   int
		dropoffLongitude int
		dropoffLatitude  int
		pickupTime       int
		driverName       string
	)

	// Create a reusable query
	rows, err := db.Query("select id, pickupLongitude, pickupLatitude, dropoffLongitude, dropoffLatitude, pickupTime, driverName from trips where pickupLongitude = ? AND pickupLatitude = ? AND dropoffLongitude = ? AND dropoffLatitude = ? AND pickupTime = ?", aLong, aLat, bLong, bLat, pTime)
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
			Id:               id,
			PickupLongitude:  pickupLongitude,
			PickupLatitude:   pickupLatitude,
			DropoffLongitude: dropoffLongitude,
			DropoffLatitude:  dropoffLatitude,
			PickupTime:       pickupTime,
			DriverName:       driverName,
		}
		trips = append(trips, trip)
	}
	db.Close()

	return trips, err
}

func getAllTripsForUserFromDB(name string) (t Trips, e error) {

	// Start db acess - use a closure or some better way to deal with opening and closing the db
	db := StartDB()

	var (
		id               int
		pickupLongitude  int
		pickupLatitude   int
		dropoffLongitude int
		dropoffLatitude  int
		pickupTime       int
		driverName       string
	)

	// Create a reusable query
	rows, err := db.Query("select id, pickupLongitude, pickupLatitude, dropoffLongitude, dropoffLatitude, pickupTime, driverName from trips where driverName = ?", name)
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
			Id:               id,
			PickupLongitude:  pickupLongitude,
			PickupLatitude:   pickupLatitude,
			DropoffLongitude: dropoffLongitude,
			DropoffLatitude:  dropoffLatitude,
			PickupTime:       pickupTime,
			DriverName:       driverName,
		}
		trips = append(trips, trip)
	}
	db.Close()

	return trips, err
}

func addTripToDB( /*db *sql.DB,*/ trip Trip) (i int64, e error) {

	// Start db acess - use a closure or some better way to deal with opening and closing the db
	db := StartDB()

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

	return lastId, err
}

func updateTripInDB(id int, trip Trip) (i int, e error) {
	db := StartDB()

	stmt, err := db.Prepare("UPDATE trips SET pickupLongitude = ?, pickupLatitude = ?, dropoffLongitude = ?, dropoffLatitude = ?, pickupTime = ?, driverName = ? WHERE id= ?")
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(trip.PickupLongitude, trip.PickupLatitude, trip.DropoffLongitude, trip.DropoffLatitude, trip.PickupTime, trip.DriverName, id)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	return id, err
}

func deleteTripFromDB(id int) (t Trips, e error) {
	db := StartDB()

	delForm, err := db.Prepare("DELETE FROM trips WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(id)

	return getAllTripsFromDB()
}

func closeDB(db *sql.DB) {
	defer db.Close()
}
