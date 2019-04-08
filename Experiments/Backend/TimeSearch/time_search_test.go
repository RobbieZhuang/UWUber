// To test time_search.go, go to this directory and run `go test`

package main

import (
    "testing"
    "reflect"
    "time"
)

func TestSimple(t *testing.T) {
    trips := []Trip{
        Trip{
            Id: 123,
            PickupLocation: "Waterloo",
            DropoffLocation: "Toronto",
            PickupTimeData: TimeData{
                Time: time.Date(2019, 4, 7, 3, 43, 0, 0, time.UTC),
            },
            SpotsAvailable: 2,
        },
        Trip{
            Id: 233,
            PickupLocation: "Waterloo",
            DropoffLocation: "Toronto",
            PickupTimeData: TimeData{
                Time: time.Date(2019, 4, 7, 3, 20, 0, 0, time.UTC),
            },
            SpotsAvailable: 1,
        },
    }

    requests := []TripRequest{
        TripRequest{
            Id: 337,
            PickupLocation: "Waterloo",
            DropoffLocation: "Toronto",
            PickupTimeData: TimeData{
                Time: time.Date(2019, 4, 7, 3, 34, 0, 0, time.UTC),
            },
        },
    }

    requestsToTrips := TimeSearch(trips, requests)
    if len(requestsToTrips) != 1 || requestsToTrips[337][0] != 123 {
        t.Errorf("Failed TestSimple")
    }
}

// The method TimeSearch is supposed to return a slice of matched trip ids for each request id
func TestMultiMatches(t *testing.T) {
    time10 := TimeData{Time: time.Date(2019, 4, 7, 3, 10, 0, 0, time.UTC)}
    time34 := TimeData{Time: time.Date(2019, 4, 7, 3, 34, 0, 0, time.UTC)}
    time45 := TimeData{Time: time.Date(2019, 4, 7, 3, 45, 0, 0, time.UTC)}
    time20 := TimeData{Time: time.Date(2019, 4, 7, 3, 20, 0, 0, time.UTC)}
    time100 := TimeData{Time: time.Date(2019, 4, 7, 4, 40, 0, 0, time.UTC)}
    time101 := TimeData{Time: time.Date(2019, 4, 7, 4, 41, 0, 0, time.UTC)}

    trips := []Trip{
        Trip{Id: 123, PickupLocation: "Waterloo", DropoffLocation: "Toronto", PickupTimeData: time45, SpotsAvailable: 2},
        Trip{Id: 233, PickupLocation: "Waterloo", DropoffLocation: "Toronto", PickupTimeData: time20, SpotsAvailable: 1},
        Trip{Id: 235, PickupLocation: "McMaster", DropoffLocation: "Toronto", PickupTimeData: time100, SpotsAvailable: 3},
        Trip{Id: 236, PickupLocation: "McMaster", DropoffLocation: "Toronto", PickupTimeData: time100, SpotsAvailable: 0},
        Trip{Id: 238, PickupLocation: "Toronto", DropoffLocation: "Waterloo", PickupTimeData: time20, SpotsAvailable: 1},
    }

    requests := []TripRequest{
        TripRequest{Id: 337, PickupLocation: "McMaster", DropoffLocation: "Toronto", PickupTimeData: time101}, // Matched
        TripRequest{Id: 338, PickupLocation: "Waterloo", DropoffLocation: "Toronto", PickupTimeData: time34},  // Matched x2
        TripRequest{Id: 339, PickupLocation: "Waterloo", DropoffLocation: "Toronto", PickupTimeData: time34},  // Matched x2
        TripRequest{Id: 343, PickupLocation: "Toronto", DropoffLocation: "Waterloo", PickupTimeData: time10},  // Matched
        TripRequest{Id: 359, PickupLocation: "McMaster", DropoffLocation: "Waterloo", PickupTimeData: time34},
    }

    requestsToTrips := TimeSearch(trips, requests)
    if len(requestsToTrips) != 5 {
        t.Errorf("Failed TestMultiMatches - wrong number of keys in map")
    }
    if !reflect.DeepEqual(requestsToTrips[337], []int{235}) {
        t.Errorf("Failed TestMultiMatches - TripRequest 337")
    }
    if !reflect.DeepEqual(requestsToTrips[338], []int{123, 233}) {
        t.Errorf("Failed TestMultiMatches - TripRequest 338")
    }
    if !reflect.DeepEqual(requestsToTrips[339], []int{123, 233}) {
        t.Errorf("Failed TestMultiMatches - TripRequest 339")
    }
    if !reflect.DeepEqual(requestsToTrips[343], []int{238}) {
        t.Errorf("Failed TestMultiMatches - TripRequest 343")
    }
    if !reflect.DeepEqual(requestsToTrips[359], []int{}) {
        t.Errorf("Failed TestMultiMatches - TripRequest 359")
    }
}
