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
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
)

// TODO
// cli opt base url
// cli opt timeout
// cli opt for own endpoint to allow layered scraping caches
// logs
// systemd unit file
// proper match & append instead ReplaceAllString .* then $1
// always pass through HTTP return code
// handle already-existing timestamps (positive match on anything which is a valid line, but without timestamp)
// https://golang.org/pkg/net/http/#ServeMux
// event & timeout driven function routing
// put in X-header for exposition format

func main() {
	var (
		showVersion = flag.Bool("version", false, "Print version information.")
		listenAddress = flag.String("web.listen-address", ":8080", "Address to listen on for web interface and telemetry.")
		metricsPath   = flag.String("web.telemetry-path", "/prometheus_scrape_cache/metrics", "Path under which to expose metrics.")
		baseUrl = flag.String("base.url", "http://localhost:8080/metrics", "Base URL to scrape from")
	)

	flag.Parse()

	if *showVersion {
		fmt.Fprintln(os.Stdout, version.Print("prometheus_scrape_cache"))
		os.Exit(0)
	}

	resp, err_get := http.Get(*baseUrl)
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
		reply_string := re.ReplaceAllString(bodyString, `$1 `+strconv.Itoa(int(epoch)))
		fmt.Println(reply_string)
	}

	// Set up Prometheus metrics endpoint
	http.Handle(*metricsPath, promhttp.Handler())
	log.Fatal(http.ListenAndServe(*listenAddress, nil))

}
