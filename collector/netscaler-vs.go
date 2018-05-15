package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rokett/citrix-netscaler-exporter/netscaler"
)

var (
	virtualServersWaitingRequests = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_waiting_requests",
			Help: "Number of requests waiting on a specific virtual server",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	virtualServersHealth = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_health",
			Help: "Percentage of UP services bound to a specific virtual server",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	virtualServersInactiveServices = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_inactive_services",
			Help: "Number of inactive services bound to a specific virtual server",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	virtualServersActiveServices = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_active_services",
			Help: "Number of active services bound to a specific virtual server",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	virtualServersTotalHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "virtual_servers_total_hits",
			Help: "Total virtual server hits",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	virtualServersHitsRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_hits_rate",
			Help: "Number of hits/second to a specific virtual server",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	virtualServersTotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "virtual_servers_total_requests",
			Help: "Total virtual server requests",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	virtualServersRequestsRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_requests_rate",
			Help: "Number of requests/second to a specific virtual server",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	virtualServersTotalResponses = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "virtual_servers_total_responses",
			Help: "Total virtual server responses",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	virtualServersReponsesRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_responses_rate",
			Help: "Number of responses/second from a specific virtual server",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	virtualServersTotalRequestBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "virtual_servers_total_request_bytes",
			Help: "Total virtual server request bytes",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	virtualServersRequestBytesRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_request_bytes_rate",
			Help: "Number of request bytes/second to a specific virtual server",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	virtualServersTotalResponseBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "virtual_servers_total_response_bytes",
			Help: "Total virtual server response bytes",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	virtualServersReponseBytesRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_reponse_bytes_rate",
			Help: "Number of response bytes/second from a specific virtual server",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	virtualServersCurrentClientConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_current_client_connections",
			Help: "Number of current client connections on a specific virtual server",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)

	virtualServersCurrentServerConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "virtual_servers_current_server_connections",
			Help: "Number of current connections to the actual servers behind the specific virtual server.",
		},
		[]string{
			"ns_instance",
			"virtual_server",
		},
	)
)

func (e *Exporter) collectVirtualServerWaitingRequests(ns netscaler.NSAPIResponse) {
	e.virtualServersWaitingRequests.Reset()

	for _, vs := range ns.VirtualServerStats {
		waitingRequests, _ := strconv.ParseFloat(vs.WaitingRequests, 64)
		e.virtualServersWaitingRequests.WithLabelValues(e.nsInstance, vs.Name).Set(waitingRequests)
	}
}

func (e *Exporter) collectVirtualServerHealth(ns netscaler.NSAPIResponse) {
	e.virtualServersHealth.Reset()

	for _, vs := range ns.VirtualServerStats {
		health, _ := strconv.ParseFloat(vs.Health, 64)
		e.virtualServersHealth.WithLabelValues(e.nsInstance, vs.Name).Set(health)
	}
}

func (e *Exporter) collectVirtualServerInactiveServices(ns netscaler.NSAPIResponse) {
	e.virtualServersInactiveServices.Reset()

	for _, vs := range ns.VirtualServerStats {
		inactiveServices, _ := strconv.ParseFloat(vs.InactiveServices, 64)
		e.virtualServersInactiveServices.WithLabelValues(e.nsInstance, vs.Name).Set(inactiveServices)
	}
}

func (e *Exporter) collectVirtualServerActiveServices(ns netscaler.NSAPIResponse) {
	e.virtualServersActiveServices.Reset()

	for _, vs := range ns.VirtualServerStats {
		activeServices, _ := strconv.ParseFloat(vs.ActiveServices, 64)
		e.virtualServersActiveServices.WithLabelValues(e.nsInstance, vs.Name).Set(activeServices)
	}
}

func (e *Exporter) collectVirtualServerTotalHits(ns netscaler.NSAPIResponse) {
	e.virtualServersTotalHits.Reset()

	for _, vs := range ns.VirtualServerStats {
		totalHits, _ := strconv.ParseFloat(vs.TotalHits, 64)
		e.virtualServersTotalHits.WithLabelValues(e.nsInstance, vs.Name).Set(totalHits)
	}
}

func (e *Exporter) collectVirtualServerHitsRate(ns netscaler.NSAPIResponse) {
	e.virtualServersHitsRate.Reset()

	for _, vs := range ns.VirtualServerStats {
		e.virtualServersHitsRate.WithLabelValues(e.nsInstance, vs.Name).Set(vs.HitsRate)
	}
}

func (e *Exporter) collectVirtualServerTotalRequests(ns netscaler.NSAPIResponse) {
	e.virtualServersTotalRequests.Reset()

	for _, vs := range ns.VirtualServerStats {
		totalRequests, _ := strconv.ParseFloat(vs.TotalRequests, 64)
		e.virtualServersTotalRequests.WithLabelValues(e.nsInstance, vs.Name).Set(totalRequests)
	}
}

func (e *Exporter) collectVirtualServerRequestsRate(ns netscaler.NSAPIResponse) {
	e.virtualServersRequestsRate.Reset()

	for _, vs := range ns.VirtualServerStats {
		e.virtualServersRequestsRate.WithLabelValues(e.nsInstance, vs.Name).Set(vs.RequestsRate)
	}
}

func (e *Exporter) collectVirtualServerTotalResponses(ns netscaler.NSAPIResponse) {
	e.virtualServersTotalResponses.Reset()

	for _, vs := range ns.VirtualServerStats {
		totalResponses, _ := strconv.ParseFloat(vs.TotalResponses, 64)
		e.virtualServersTotalResponses.WithLabelValues(e.nsInstance, vs.Name).Set(totalResponses)
	}
}

func (e *Exporter) collectVirtualServerResponsesRate(ns netscaler.NSAPIResponse) {
	e.virtualServersReponsesRate.Reset()

	for _, vs := range ns.VirtualServerStats {
		e.virtualServersReponsesRate.WithLabelValues(e.nsInstance, vs.Name).Set(vs.ResponsesRate)
	}
}

func (e *Exporter) collectVirtualServerTotalRequestBytes(ns netscaler.NSAPIResponse) {
	e.virtualServersTotalRequestBytes.Reset()

	for _, vs := range ns.VirtualServerStats {
		totalRequestBytes, _ := strconv.ParseFloat(vs.TotalRequestBytes, 64)
		e.virtualServersTotalRequestBytes.WithLabelValues(e.nsInstance, vs.Name).Set(totalRequestBytes)
	}
}

func (e *Exporter) collectVirtualServerRequestBytesRate(ns netscaler.NSAPIResponse) {
	e.virtualServersRequestBytesRate.Reset()

	for _, vs := range ns.VirtualServerStats {
		e.virtualServersRequestBytesRate.WithLabelValues(e.nsInstance, vs.Name).Set(vs.RequestBytesRate)
	}
}

func (e *Exporter) collectVirtualServerTotalResponseBytes(ns netscaler.NSAPIResponse) {
	e.virtualServersTotalResponseBytes.Reset()

	for _, vs := range ns.VirtualServerStats {
		totalResponseBytes, _ := strconv.ParseFloat(vs.TotalResponseBytes, 64)
		e.virtualServersTotalResponseBytes.WithLabelValues(e.nsInstance, vs.Name).Set(totalResponseBytes)
	}
}

func (e *Exporter) collectVirtualServerResponseBytesRate(ns netscaler.NSAPIResponse) {
	e.virtualServersReponseBytesRate.Reset()

	for _, vs := range ns.VirtualServerStats {
		e.virtualServersReponseBytesRate.WithLabelValues(e.nsInstance, vs.Name).Set(vs.ResponseBytesRate)
	}
}

func (e *Exporter) collectVirtualServerCurrentClientConnections(ns netscaler.NSAPIResponse) {
	e.virtualServersCurrentClientConnections.Reset()

	for _, vs := range ns.VirtualServerStats {
		currentClientConnections, _ := strconv.ParseFloat(vs.CurrentClientConnections, 64)
		e.virtualServersCurrentClientConnections.WithLabelValues(e.nsInstance, vs.Name).Set(currentClientConnections)
	}
}

func (e *Exporter) collectVirtualServerCurrentServerConnections(ns netscaler.NSAPIResponse) {
	e.virtualServersCurrentServerConnections.Reset()

	for _, vs := range ns.VirtualServerStats {
		currentServerConnections, _ := strconv.ParseFloat(vs.CurrentServerConnections, 64)
		e.virtualServersCurrentServerConnections.WithLabelValues(e.nsInstance, vs.Name).Set(currentServerConnections)
	}
}