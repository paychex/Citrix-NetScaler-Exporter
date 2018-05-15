package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rokett/citrix-netscaler-exporter/collector"

	logger "github.com/prometheus/common/log"

	"github.com/jamiealquiza/envy"
)

var (
	app        = "Citrix-NetScaler-Exporter"
	version    string
	build      string
	nsURL      = flag.String("url", "", "Base URL of the NetScaler management interface.  Normally something like https://my-netscaler.something.x")
	username   = flag.String("username", "", "Username with which to connect to the NetScaler API")
	password   = flag.String("password", "", "Password with which to connect to the NetScaler API")
	bindPort   = flag.Int("bind_port", 9280, "Port to bind the exporter endpoint to")
	versionFlg = flag.Bool("version", false, "Display application version")
	multiQuery = flag.Bool("multi", false, "Enable query endpoint")
)

func queryHandler(w http.ResponseWriter, r *http.Request) {
	queryURL := r.URL.Query().Get("target")
	if queryURL == "" {
		http.Error(w, "'target' parameter must be specified", 400)
		return
	}
	queryURL = "http://" + queryURL

	logger.Info("Scraping target " + queryURL)

	exporter, err := collector.NewExporter(&queryURL, username, password)
	if err != nil {
		http.Error(w, "Error creating exporter"+err.Error(), 400)
		logger.Error(err)
		return
	}

	registry := prometheus.NewRegistry()
	registry.MustRegister(exporter)

	// Delegate http serving to Promethues client library, which will call collector.Collect.
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

func main() {

	// This makes it work better with env variables
	// looks for NETSCALER_USERNAME & NETSCALER_PASSWORD & NETSCALER_URL & NETSCALER_MULTI

	flag.Parse()
	envy.Parse("NETSCALER")

	if *versionFlg {
		fmt.Println(app + " v" + version + " build " + build)
		os.Exit(0)
	}

	if *username == "" || *password == "" {
		fmt.Println("Missing username or password")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *nsURL == "" && !*multiQuery {
		fmt.Println("Missing URL or multiquery flag")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// This can go one of two ways
	// either just monitor one device or go into a query mode based on flag/env variable "multiquery"
	// to allow for multiple systems querying
	if *multiQuery {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`<html>
				<head>
				<title>Citrix Netscaler Exporter</title>
				<style>
				label{
				display:inline-block;
				width:75px;
				}
				form label {
				margin: 10px;
				}
				form input {
				margin: 10px;
				}
				</style>
				</head>
				<body>
				<h1>Netscaler Exporter</h1>
				<form action="/query">
				<label>Target:</label> <input type="text" name="target" placeholder="X.X.X.X" value="1.2.3.4"><br>
				<input type="submit" value="Submit">
				</form>
				</html>`))
		})

		http.HandleFunc("/query", queryHandler)     // Endpoint to do specific cluster scrapes.
		http.Handle("/metrics", promhttp.Handler()) // endpoint for exporter stats
	} else {

		u, err := url.Parse(*nsURL)
		if err != nil {
			logger.Fatal(err)
		}

		exporter, _ := collector.NewExporter(&u.Host, username, password)
		prometheus.MustRegister(exporter)

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`<html>
			<head><title>Citrix NetScaler Exporter</title></head>
			<body>
			<h1>Citrix NetScaler Exporter</h1>
			<p><a href="/metrics">Metrics</a></p>
			</body>
			</html>`))
		})

		http.Handle("/metrics", promhttp.Handler())
	}

	listeningPort := ":" + strconv.Itoa(*bindPort)
	logger.Info("Listening on port: " + listeningPort)

	err := http.ListenAndServe(listeningPort, nil)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}
