package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var username string = "a_username_"

//With 10 million users in the database, the service should support 4000 concurrent logins per second,
//with at least 200 unique users in the concurrent requests.
func main() {

	var i int64
	for i = 0; i < 200; i++ {
		statusCode, err := sendRegisterPostRequest(username+"_"+strconv.FormatInt(i, 10), username+"_"+strconv.FormatInt(i, 10))
		if err != nil {
			fmt.Println(err)
		}
		if statusCode != 200 {
			fmt.Println("statusCode: ", statusCode)
		}
		if i%5 == 0 {
			time.Sleep(1000 * time.Millisecond)
		}
	}
}

func sendRegisterPostRequest(username, password string) (int, error) {
	res, err := http.PostForm("http://127.0.0.1:80/api/users", url.Values{
		"username": {username},
		"password": {password},
	})

	if err != nil {
		return 0, err
	}
	return res.StatusCode, nil

}
