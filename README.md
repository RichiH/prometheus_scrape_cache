# Intro

This is basically a caching reverse proxy which adds timestamps to scraped metrics.
Think of it as a pullgateway instead of a [pushgateway](https://github.com/prometheus/pushgateway).
I might rename it pullgateway; TBD.

# Why?

Prometheus 1.x and 2.0.x have globally synchronized scrape offsets as long as they are fed with the same service discovery.
This will likely change with 2.1.x and later.
Currently, a standard caching proxy will do the right thing even if you're not aware of the Prometheus-internal scrape implementation.
After that, the load on your scrape targets will multiply.
This tool is written to avoid the fiery death of SNMP devices attacked via the [snmp_exporter](https://github.com/prometheus/snmp_exporter), but I can't rule out that it will have other valid uses.

As a general rule of thumb, you should only use this program if your load tests confirmed that you need it.

That being said, you can also use the scrape cache to sync the results stored within a HA pair if you so choose.

# Usage

The scrape cache is completely transparent and will forward both HTTP code and body; only `/prometheus_scrape_cache/*`is reserved for its own use.
This means it's currently impossible to get at the metrics of a scrape cache behind another scrape cache.

# Caveats

* You might introduce delay in your scrapes
* New staleness in Prometheus 2.x is not yet possible with this tool, but that's on the TODO list
