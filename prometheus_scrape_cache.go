package main

import (
	"fmt"
	"io/ioutil"
//	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

// TODO
// cli opt base url
// cli opt timeout
// metrics endpoint
// logs
// systemd unit file

func main() {
	resp, err_get := http.Get("http://demo.robustperception.io:9090/metrics")
	if err_get != nil {
		// TODO handle error
	}
	epoch := time.Now().Unix()
	fmt.Println(epoch)


	// close connection
	defer resp.Body.Close()

	if resp.StatusCode == 200 { // OK
		bodyBytes, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			// TODO handle error
		}
		bodyString := string(bodyBytes)
//		fmt.Println(bodyString)
		// build regexp
		var re = regexp.MustCompile("(?m)(^[^#].*$)")
		reply_string := re.ReplaceAllString(bodyString, `$1 ` + strconv.Itoa(int(epoch)))
		fmt.Println(reply_string)
	}



}
