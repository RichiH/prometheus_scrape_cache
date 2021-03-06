# Status

NOT YET WORKING


# Intro

This is basically a caching reverse proxy which adds timestamps to scraped metrics.
Think of it as a pullgateway instead of a [pushgateway](https://github.com/prometheus/pushgateway).
I might rename it to pullgateway; we will find out together.


# Why?

Prometheus 1.x and 2.0.x have globally synchronized scrape offsets as long as they are fed from the same service discovery.
This will likely change with 2.3.x or later.
Currently, a standard caching proxy will do the right thing even if you're not aware of the Prometheus-internal scrape implementation.
Once that's changed, a normal caching proxy will not suffice any more and the load on your scrape targets will multiply by up to the amount of Prometheus servers you run.
This tool is written to avoid the fiery death of SNMP devices attacked via the [snmp_exporter](https://github.com/prometheus/snmp_exporter), but I can't rule out that it will have other valid uses.

As a general rule of thumb, you should only use this program if your load tests confirmed that you need it.

That being said, you can also use the scrape cache to sync the results stored within a HA pair if you so choose.


# Usage

The scrape cache is completely transparent and will forward both HTTP code and body; only `/prometheus_scrape_cache/*`is reserved for its own use.

You will need one instance of the scrape cache per exporter/instrumentation.
This keeps the implementation sane, scales well, and follows the Unix^w microservice mantra of "one thing well".


# Current caveats

* You might introduce delay in your scrapes
* New staleness in Prometheus 2.x is not yet possible with this tool, but that's on the TODO list
* It's currently impossible to get at the metrics of a scrape cache behind another scrape cache
* If you scrape anything which already has timestamps, things go boom for now

# Further reading

* Slices are flushed out after 0.5 * slice length, per default 1h: https://github.com/prometheus/tsdb/blob/494acd307058387ced7646f9996b0f7372eaa558/db.go#L377-L392 
