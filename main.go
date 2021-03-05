package main

import (
	"log"
	"os"
	"strconv"
)

func main() {
	reseedStr := os.Getenv("PM_RESEED")

	reseed, err := strconv.ParseBool(reseedStr)

	if err != nil {
		log.Println("REseed string not found")
		reseed = false
	}

	log.Println(reseedStr)

	Run(reseed)
}
