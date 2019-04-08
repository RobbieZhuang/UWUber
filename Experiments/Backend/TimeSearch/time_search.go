// Search for matching trip and trip request pairs given a list of trips and and a list of requests

package main

import (
    "time"
)


// Any request and trip within <= 15 min. of each other will be considered matching in time
// Change this for different time thresholds
const TIME_THRESHOLD = 15 * time.Minute

// Test if the times of trip and requests are within a certain threshold
func timeMatches(tripTimeData TimeData, requestTimeData TimeData) (bool) {
    duration := tripTimeData.Time.Sub(requestTimeData.Time)
    if duration < 0 {
        duration = -duration
    }
    return duration <= TIME_THRESHOLD
}

// Test if the trip and the request satisfy matching requirements
// Change this for different requirements for matches
func tripMatches(trip Trip, request TripRequest) (bool) {
    return trip.PickupLocation == request.PickupLocation &&
        trip.DropoffLocation == request.DropoffLocation &&
        trip.SpotsAvailable > 0 &&
        timeMatches(trip.PickupTimeData, request.PickupTimeData)
}

// Given a list of trips and and a list of requests, output a map of matching request ids to a slice of trips
func TimeSearch(trips []Trip, tripRequests []TripRequest) (map[int][]int) {
    requestsToTrips := make(map[int][]int)

    // For each trip, find the list of requests that satisfy its requirements
    for _, request := range tripRequests {
        requestsToTrips[request.Id] = []int{}

        for _, trip := range trips {
            // Ignore if the trip does not match the request
            if !tripMatches(trip, request) {
                continue
            }

            // Trip and request match, add to the maps
            if slice, ok := requestsToTrips[request.Id]; ok {
                // Add the trip to the map
                slice = append(slice, trip.Id)
                requestsToTrips[request.Id] = slice
            } else {
                // Request not yet added to map, add request as key to map
                slice := []int{trip.Id}
                requestsToTrips[request.Id] = slice
            }
        }
    }

    return requestsToTrips
}
