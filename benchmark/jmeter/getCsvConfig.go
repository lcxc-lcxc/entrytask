package main

/**
this is for create a jmeter csv config file
*/

import (
	"os"
	"strconv"
)

var username string = "a_username_"

// create 200 unique user
func main() {
	csvConfig, err := os.OpenFile("benchmark/jmeter/csvConfig.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return
	}

	var i int64
	for i = 0; i < 200; i++ {
		csvConfig.Write([]byte(username + "_" + strconv.FormatInt(i, 10)))
		csvConfig.Write([]byte(","))
		csvConfig.Write([]byte(username + "_" + strconv.FormatInt(i, 10)))
		csvConfig.Write([]byte("\n"))
	}

}
