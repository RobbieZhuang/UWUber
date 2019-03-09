package main

import (
	"log"
	"strconv"
)

func getIntFor(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}
	return i
}
