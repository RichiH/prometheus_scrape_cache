package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// TODO
// cli opt base url
// cli opt timeout
// cli opt for own endpoint to allow layered scraping caches
// metrics endpoint
// logs
// systemd unit file
// proper match & append instead ReplaceAllString .* then $1
// always pass through HTTP return code
// handle already-existing timestamps (positive match on anything which is a valid line, but without timestamp)
// https://golang.org/pkg/net/http/#ServeMux
// event & timeout driven function routing
// put in X-header for exposition format

func main() {
	flag.Parse()

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



	// Set up Prometheus metrics endpoint
	var addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
	http.Handle("/prometheus_scrape_cache/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
