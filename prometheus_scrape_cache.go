package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
)

// TODO
// cli opt timeout
// cli opt for own endpoint to allow layered scraping caches
// systemd unit file
// proper match & append instead ReplaceAllString .* then $1
// always pass through HTTP return code
// handle already-existing timestamps (positive match on anything which is a valid line, but without timestamp)
// https://golang.org/pkg/net/http/#ServeMux
// event & timeout driven function routing
// put in X-header for exposition format

func main() {
	var (
		listenAddress = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":8080").String()
		metricsPath   = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/prometheus_scrape_cache/metrics").String()
		baseUrl       = kingpin.Flag("base.url", "Base URL to scrape from").Default("http://demo.robustperception.io:9090/metrics").String()
		//TODO remove demo URL and force user to set a value
	)

	log.AddFlags(kingpin.CommandLine)
	kingpin.Version(version.Print("prometheus_scrape_cache"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	log.Infoln("Starting prometheus_scrape_cache", version.Info())
	log.Infoln("Build context", version.BuildContext())

	resp, err := http.Get(*baseUrl)
	if err != nil {
		log.Fatalf("Couldn't scrape metrics: %s", err)
	}
	epoch := time.Now().Unix()

	// close connection
	defer resp.Body.Close()

	if resp.StatusCode == 200 { // OK
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			// TODO For production, we will need to save the text & error code and return that for any scrapes
			log.Fatalf("Couldn't read body: %s", err)
		}
		bodyString := string(bodyBytes)
		// build regexp
		var re = regexp.MustCompile("(?m)(^[^#].*$)")
		reply_string := re.ReplaceAllString(bodyString, `$1 `+strconv.Itoa(int(epoch)))
		fmt.Println(reply_string)
	}

	// Set up Prometheus metrics endpoint
	http.Handle(*metricsPath, promhttp.Handler())
	log.Fatal(http.ListenAndServe(*listenAddress, nil))

}
