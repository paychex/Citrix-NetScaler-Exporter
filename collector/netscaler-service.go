package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rokett/citrix-netscaler-exporter/netscaler"
)

var (
	servicesThroughput = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "service_throughput",
			Help: "Number of bytes received or sent by this service (Mbps)",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesThroughputRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_throughput_rate",
			Help: "Rate (/s) counter for throughput",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesAvgTTFB = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_average_time_to_first_byte",
			Help: "Average TTFB between the NetScaler appliance and the server.TTFB is the time interval between sending the request packet to a service and receiving the first response from the service",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_state",
			Help: "Current state of the service",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesTotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "service_total_requests",
			Help: "Total number of requests received on this service",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesRequestsRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_request_rate",
			Help: "Rate (/s) counter for totalrequests",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesTotalResponses = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "service_total_responses",
			Help: "Total number of responses received on this service",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesResponsesRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_responses_rate",
			Help: "Rate (/s) counter for totalresponses",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesTotalRequestBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "service_total_request_bytes",
			Help: "Total number of request bytes received on this service",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesRequestBytesRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_request_bytes_rate",
			Help: "Rate (/s) counter for totalrequestbytes",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesTotalResponseBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "service_total_response_bytes",
			Help: "Total number of response bytes received on this service",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesResponseBytesRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_response_bytes_rate",
			Help: "Rate (/s) counter for totalresponsebytes",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesCurrentClientConns = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_current_client_connections",
			Help: "Number of current client connections",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesSurgeCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_surge_count",
			Help: "Number of requests in the surge queue",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesCurrentServerConns = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_current_server_connections",
			Help: "Number of current connections to the actual servers",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesServerEstablishedConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_server_established_connections",
			Help: "Number of server connections in ESTABLISHED state",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesCurrentReusePool = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_current_reuse_pool",
			Help: "Number of requests in the idle queue/reuse pool.",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesMaxClients = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_max_clients",
			Help: "Maximum open connections allowed on this service",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesCurrentLoad = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_current_load",
			Help: "Load on the service that is calculated from the bound load based monitor",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesVirtualServerServiceHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "service_virtual_server_service_hits",
			Help: "Number of times that the service has been provided",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesVirtualServerServiceHitsRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_virtual_server_service_hits_rate",
			Help: "Rate (/s) counter for vsvrservicehits",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	servicesActiveTransactions = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_active_transactions",
			Help: "Number of active transactions handled by this service. (Including those in the surge queue.) Active Transaction means number of transactions currently served by the server including those waiting in the SurgeQ",
		},
		[]string{
			"ns_instance",
			"service",
		},
	)

	serviceGroupsState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "servicegroup_state",
			Help: "Current state of the server",
		},
		[]string{
			"ns_instance",
			"servicegroup",
			"member",
		},
	)

	serviceGroupsAvgTTFB = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "servicegroup_average_time_to_first_byte",
			Help: "Average TTFB between the NetScaler appliance and the server.",
		},
		[]string{
			"ns_instance",
			"servicegroup",
			"member",
		},
	)

	serviceGroupsTotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "servicegroup_total_requests",
			Help: "Total number of requests received on this service",
		},
		[]string{
			"ns_instance",
			"servicegroup",
			"member",
		},
	)

	serviceGroupsRequestsRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "servicegroup_requests_rate",
			Help: "Rate (/s) counter for totalrequests",
		},
		[]string{
			"ns_instance",
			"servicegroup",
			"member",
		},
	)

	serviceGroupsTotalResponses = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "servicegroup_total_responses",
			Help: "Number of responses received on this service.",
		},
		[]string{
			"ns_instance",
			"servicegroup",
			"member",
		},
	)

	serviceGroupsResponsesRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "servicegroup_responses_rate",
			Help: "Rate (/s) counter for totalresponses",
		},
		[]string{
			"ns_instance",
			"servicegroup",
			"member",
		},
	)

	serviceGroupsTotalRequestBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "servicegroup_total_request_bytes",
			Help: "Total number of request bytes received on this service",
		},
		[]string{
			"ns_instance",
			"servicegroup",
			"member",
		},
	)

	serviceGroupsRequestBytesRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "servicegroup_request_bytes_rate",
			Help: "Rate (/s) counter for totalrequestbytes",
		},
		[]string{
			"ns_instance",
			"servicegroup",
			"member",
		},
	)

	serviceGroupsTotalResponseBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "servicegroup_total_response_bytes",
			Help: "Number of response bytes received by this service",
		},
		[]string{
			"ns_instance",
			"servicegroup",
			"member",
		},
	)

	serviceGroupsResponseBytesRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "servicegroup_response_bytes_rate",
			Help: "Rate (/s) counter for totalresponsebytes",
		},
		[]string{
			"ns_instance",
			"servicegroup",
			"member",
		},
	)

	serviceGroupsCurrentClientConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "servicegroup_current_client_connections",
			Help: "Number of current client connections.",
		},
		[]string{
			"ns_instance",
			"servicegroup",
			"member",
		},
	)

	serviceGroupsSurgeCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "servicegroup_surge_count",
			Help: "Number of requests in the surge queue.",
		},
		[]string{
			"ns_instance",
			"servicegroup",
			"member",
		},
	)

	serviceGroupsCurrentServerConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "servicegroup_current_server_connections",
			Help: "Number of current connections to the actual servers behind the virtual server.",
		},
		[]string{
			"ns_instance",
			"servicegroup",
			"member",
		},
	)

	serviceGroupsServerEstablishedConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "servicegroup_server_established_connections",
			Help: "Number of server connections in ESTABLISHED state.",
		},
		[]string{
			"ns_instance",
			"servicegroup",
			"member",
		},
	)

	serviceGroupsCurrentReusePool = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "servicegroup_current_reuse_pool",
			Help: "Number of requests in the idle queue/reuse pool.",
		},
		[]string{
			"ns_instance",
			"servicegroup",
			"member",
		},
	)

	serviceGroupsMaxClients = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "servicegroup_max_clients",
			Help: "Maximum open connections allowed on this service.",
		},
		[]string{
			"ns_instance",
			"servicegroup",
			"member",
		},
	)
)

func (e *Exporter) collectServicesThroughput(ns netscaler.NSAPIResponse) {
	e.servicesThroughput.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.Throughput, 64)
		e.servicesThroughput.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesThroughputRate(ns netscaler.NSAPIResponse) {
	e.servicesThroughputRate.Reset()

	for _, service := range ns.ServiceStats {
		e.servicesThroughputRate.WithLabelValues(e.nsInstance, service.Name).Set(service.ThroughputRate)
	}
}

func (e *Exporter) collectServicesAvgTTFB(ns netscaler.NSAPIResponse) {
	e.servicesAvgTTFB.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.AvgTimeToFirstByte, 64)
		e.servicesAvgTTFB.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesState(ns netscaler.NSAPIResponse) {
	e.servicesState.Reset()

	for _, service := range ns.ServiceStats {
		state := 0.0

		if service.State == "UP" {
			state = 1.0
		}

		e.servicesState.WithLabelValues(e.nsInstance, service.Name).Set(state)
	}
}

func (e *Exporter) collectServicesTotalRequests(ns netscaler.NSAPIResponse) {
	e.servicesTotalRequests.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.TotalRequests, 64)
		e.servicesTotalRequests.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesRequestsRate(ns netscaler.NSAPIResponse) {
	e.servicesRequestsRate.Reset()

	for _, service := range ns.ServiceStats {
		e.servicesRequestsRate.WithLabelValues(e.nsInstance, service.Name).Set(service.RequestsRate)
	}
}

func (e *Exporter) collectServicesTotalResponses(ns netscaler.NSAPIResponse) {
	e.servicesTotalResponses.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.TotalResponses, 64)
		e.servicesTotalResponses.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesResponsesRate(ns netscaler.NSAPIResponse) {
	e.servicesResponsesRate.Reset()

	for _, service := range ns.ServiceStats {
		e.servicesResponsesRate.WithLabelValues(e.nsInstance, service.Name).Set(service.ResponsesRate)
	}
}

func (e *Exporter) collectServicesTotalRequestBytes(ns netscaler.NSAPIResponse) {
	e.servicesTotalRequestBytes.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.TotalRequestBytes, 64)
		e.servicesTotalRequestBytes.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesRequestBytesRate(ns netscaler.NSAPIResponse) {
	e.servicesRequestBytesRate.Reset()

	for _, service := range ns.ServiceStats {
		e.servicesRequestBytesRate.WithLabelValues(e.nsInstance, service.Name).Set(service.RequestBytesRate)
	}
}

func (e *Exporter) collectServicesTotalResponseBytes(ns netscaler.NSAPIResponse) {
	e.servicesTotalResponseBytes.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.TotalResponseBytes, 64)
		e.servicesTotalResponseBytes.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesResponseBytesRate(ns netscaler.NSAPIResponse) {
	e.servicesResponseBytesRate.Reset()

	for _, service := range ns.ServiceStats {
		e.servicesResponseBytesRate.WithLabelValues(e.nsInstance, service.Name).Set(service.ResponseBytesRate)
	}
}

func (e *Exporter) collectServicesCurrentClientConns(ns netscaler.NSAPIResponse) {
	e.servicesCurrentClientConns.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentClientConnections, 64)
		e.servicesCurrentClientConns.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesSurgeCount(ns netscaler.NSAPIResponse) {
	e.servicesSurgeCount.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.SurgeCount, 64)
		e.servicesSurgeCount.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesCurrentServerConns(ns netscaler.NSAPIResponse) {
	e.servicesCurrentServerConns.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentServerConnections, 64)
		e.servicesCurrentServerConns.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesServerEstablishedConnections(ns netscaler.NSAPIResponse) {
	e.servicesServerEstablishedConnections.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.ServerEstablishedConnections, 64)
		e.servicesServerEstablishedConnections.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesCurrentReusePool(ns netscaler.NSAPIResponse) {
	e.servicesCurrentReusePool.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentReusePool, 64)
		e.servicesCurrentReusePool.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesMaxClients(ns netscaler.NSAPIResponse) {
	e.servicesMaxClients.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.MaxClients, 64)
		e.servicesMaxClients.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesCurrentLoad(ns netscaler.NSAPIResponse) {
	e.servicesCurrentLoad.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.CurrentLoad, 64)
		e.servicesCurrentLoad.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesVirtualServerServiceHits(ns netscaler.NSAPIResponse) {
	e.servicesVirtualServerServiceHits.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.ServiceHits, 64)
		e.servicesVirtualServerServiceHits.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServicesVirtualServerServiceHitsRate(ns netscaler.NSAPIResponse) {
	e.servicesVirtualServerServiceHitsRate.Reset()

	for _, service := range ns.ServiceStats {
		e.servicesVirtualServerServiceHitsRate.WithLabelValues(e.nsInstance, service.Name).Set(service.ServiceHitsRate)
	}
}

func (e *Exporter) collectServicesActiveTransactions(ns netscaler.NSAPIResponse) {
	e.servicesActiveTransactions.Reset()

	for _, service := range ns.ServiceStats {
		val, _ := strconv.ParseFloat(service.ActiveTransactions, 64)
		e.servicesActiveTransactions.WithLabelValues(e.nsInstance, service.Name).Set(val)
	}
}

func (e *Exporter) collectServiceGroupsState(ns netscaler.NSAPIResponse, sgName string, servername string) {
	e.serviceGroupsState.Reset()

	for _, sg := range ns.ServiceGroupMemberStats {
		state := 0.0

		if sg.State == "UP" {
			state = 1.0
		}

		e.serviceGroupsState.WithLabelValues(e.nsInstance, sgName, servername).Set(state)
	}
}

func (e *Exporter) collectServiceGroupsAvgTTFB(ns netscaler.NSAPIResponse, sgName string, servername string) {
	e.serviceGroupsAvgTTFB.Reset()

	for _, sg := range ns.ServiceGroupMemberStats {
		val, _ := strconv.ParseFloat(sg.AvgTimeToFirstByte, 64)
		e.serviceGroupsAvgTTFB.WithLabelValues(e.nsInstance, sgName, servername).Set(val)
	}
}

func (e *Exporter) collectServiceGroupsTotalRequests(ns netscaler.NSAPIResponse, sgName string, servername string) {
	e.serviceGroupsTotalRequests.Reset()

	for _, sg := range ns.ServiceGroupMemberStats {
		val, _ := strconv.ParseFloat(sg.TotalRequests, 64)
		e.serviceGroupsTotalRequests.WithLabelValues(e.nsInstance, sgName, servername).Set(val)
	}
}

func (e *Exporter) collectServiceGroupsRequestsRate(ns netscaler.NSAPIResponse, sgName string, servername string) {
	e.serviceGroupsRequestsRate.Reset()

	for _, sg := range ns.ServiceGroupMemberStats {
		e.serviceGroupsRequestsRate.WithLabelValues(e.nsInstance, sgName, servername).Set(sg.RequestsRate)
	}
}

func (e *Exporter) collectServiceGroupsTotalResponses(ns netscaler.NSAPIResponse, sgName string, servername string) {
	e.serviceGroupsTotalResponses.Reset()

	for _, sg := range ns.ServiceGroupMemberStats {
		val, _ := strconv.ParseFloat(sg.TotalResponses, 64)
		e.serviceGroupsTotalResponses.WithLabelValues(e.nsInstance, sgName, servername).Set(val)
	}
}

func (e *Exporter) collectServiceGroupsResponsesRate(ns netscaler.NSAPIResponse, sgName string, servername string) {
	e.serviceGroupsResponsesRate.Reset()

	for _, sg := range ns.ServiceGroupMemberStats {
		e.serviceGroupsResponsesRate.WithLabelValues(e.nsInstance, sgName, servername).Set(sg.ResponsesRate)
	}
}

func (e *Exporter) collectServiceGroupsTotalRequestBytes(ns netscaler.NSAPIResponse, sgName string, servername string) {
	e.serviceGroupsTotalRequestBytes.Reset()

	for _, sg := range ns.ServiceGroupMemberStats {
		val, _ := strconv.ParseFloat(sg.TotalRequestBytes, 64)
		e.serviceGroupsTotalRequestBytes.WithLabelValues(e.nsInstance, sgName, servername).Set(val)
	}
}

func (e *Exporter) collectServiceGroupsRequestBytesRate(ns netscaler.NSAPIResponse, sgName string, servername string) {
	e.serviceGroupsRequestBytesRate.Reset()

	for _, sg := range ns.ServiceGroupMemberStats {
		e.serviceGroupsRequestBytesRate.WithLabelValues(e.nsInstance, sgName, servername).Set(sg.RequestBytesRate)
	}
}

func (e *Exporter) collectServiceGroupsTotalResponseBytes(ns netscaler.NSAPIResponse, sgName string, servername string) {
	e.serviceGroupsTotalResponseBytes.Reset()

	for _, sg := range ns.ServiceGroupMemberStats {
		val, _ := strconv.ParseFloat(sg.TotalResponseBytes, 64)
		e.serviceGroupsTotalResponseBytes.WithLabelValues(e.nsInstance, sgName, servername).Set(val)
	}
}

func (e *Exporter) collectServiceGroupsResponseBytesRate(ns netscaler.NSAPIResponse, sgName string, servername string) {
	e.serviceGroupsResponseBytesRate.Reset()

	for _, sg := range ns.ServiceGroupMemberStats {
		e.serviceGroupsResponseBytesRate.WithLabelValues(e.nsInstance, sgName, servername).Set(sg.ResponseBytesRate)
	}
}

func (e *Exporter) collectServiceGroupsCurrentClientConnections(ns netscaler.NSAPIResponse, sgName string, servername string) {
	e.serviceGroupsCurrentClientConnections.Reset()

	for _, sg := range ns.ServiceGroupMemberStats {
		val, _ := strconv.ParseFloat(sg.CurrentClientConnections, 64)
		e.serviceGroupsCurrentClientConnections.WithLabelValues(e.nsInstance, sgName, servername).Set(val)
	}
}

func (e *Exporter) collectServiceGroupsSurgeCount(ns netscaler.NSAPIResponse, sgName string, servername string) {
	e.serviceGroupsSurgeCount.Reset()

	for _, sg := range ns.ServiceGroupMemberStats {
		val, _ := strconv.ParseFloat(sg.SurgeCount, 64)
		e.serviceGroupsSurgeCount.WithLabelValues(e.nsInstance, sgName, servername).Set(val)
	}
}

func (e *Exporter) collectServiceGroupsCurrentServerConnections(ns netscaler.NSAPIResponse, sgName string, servername string) {
	e.serviceGroupsCurrentServerConnections.Reset()

	for _, sg := range ns.ServiceGroupMemberStats {
		val, _ := strconv.ParseFloat(sg.CurrentServerConnections, 64)
		e.serviceGroupsCurrentServerConnections.WithLabelValues(e.nsInstance, sgName, servername).Set(val)
	}
}

func (e *Exporter) collectServiceGroupsServerEstablishedConnections(ns netscaler.NSAPIResponse, sgName string, servername string) {
	e.serviceGroupsServerEstablishedConnections.Reset()

	for _, sg := range ns.ServiceGroupMemberStats {
		val, _ := strconv.ParseFloat(sg.ServerEstablishedConnections, 64)
		e.serviceGroupsServerEstablishedConnections.WithLabelValues(e.nsInstance, sgName, servername).Set(val)
	}
}

func (e *Exporter) collectServiceGroupsCurrentReusePool(ns netscaler.NSAPIResponse, sgName string, servername string) {
	e.serviceGroupsCurrentReusePool.Reset()

	for _, sg := range ns.ServiceGroupMemberStats {
		val, _ := strconv.ParseFloat(sg.CurrentReusePool, 64)
		e.serviceGroupsCurrentReusePool.WithLabelValues(e.nsInstance, sgName, servername).Set(val)
	}
}

func (e *Exporter) collectServiceGroupsMaxClients(ns netscaler.NSAPIResponse, sgName string, servername string) {
	e.serviceGroupsMaxClients.Reset()

	for _, sg := range ns.ServiceGroupMemberStats {
		val, _ := strconv.ParseFloat(sg.MaxClients, 64)
		e.serviceGroupsMaxClients.WithLabelValues(e.nsInstance, sgName, servername).Set(val)
	}
}
