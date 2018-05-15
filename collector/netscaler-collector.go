package collector

import (
	"errors"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	logger "github.com/prometheus/common/log"
	"github.com/rokett/citrix-netscaler-exporter/netscaler"
)

var (
	nsInstance string

	modelID = prometheus.NewDesc(
		"model_id",
		"NetScaler model - reflects the bandwidth available; for example VPX 10 would report as 10.",
		[]string{
			"ns_instance",
		},
		nil,
	)

	mgmtCPUUsage = prometheus.NewDesc(
		"mgmt_cpu_usage",
		"Current CPU utilisation for management",
		[]string{
			"ns_instance",
		},
		nil,
	)

	pktCPUUsage = prometheus.NewDesc(
		"pkt_cpu_usage",
		"Current CPU utilisation for packet engines, excluding management",
		[]string{
			"ns_instance",
		},
		nil,
	)

	memUsage = prometheus.NewDesc(
		"mem_usage",
		"Current memory utilisation",
		[]string{
			"ns_instance",
		},
		nil,
	)

	flashPartitionUsage = prometheus.NewDesc(
		"flash_partition_usage",
		"Used space in /flash partition of the disk, as a percentage.",
		[]string{
			"ns_instance",
		},
		nil,
	)

	varPartitionUsage = prometheus.NewDesc(
		"var_partition_usage",
		"Used space in /var partition of the disk, as a percentage. ",
		[]string{
			"ns_instance",
		},
		nil,
	)

	rxMbPerSec = prometheus.NewDesc(
		"received_mb_per_second",
		"Number of Megabits received by the NetScaler appliance per second",
		[]string{
			"ns_instance",
		},
		nil,
	)

	txMbPerSec = prometheus.NewDesc(
		"transmit_mb_per_second",
		"Number of Megabits transmitted by the NetScaler appliance per second",
		[]string{
			"ns_instance",
		},
		nil,
	)

	httpRequestsRate = prometheus.NewDesc(
		"http_requests_rate",
		"HTTP requests received per second",
		[]string{
			"ns_instance",
		},
		nil,
	)

	httpResponsesRate = prometheus.NewDesc(
		"http_responses_rate",
		"HTTP requests sent per second",
		[]string{
			"ns_instance",
		},
		nil,
	)

	tcpCurrentClientConnections = prometheus.NewDesc(
		"tcp_current_client_connections",
		"Client connections, including connections in the Opening, Established, and Closing state.",
		[]string{
			"ns_instance",
		},
		nil,
	)

	tcpCurrentClientConnectionsEstablished = prometheus.NewDesc(
		"tcp_current_client_connections_established",
		"Current client connections in the Established state, which indicates that data transfer can occur between the NetScaler and the client.",
		[]string{
			"ns_instance",
		},
		nil,
	)

	tcpCurrentServerConnections = prometheus.NewDesc(
		"tcp_current_server_connections",
		"Server connections, including connections in the Opening, Established, and Closing state.",
		[]string{
			"ns_instance",
		},
		nil,
	)

	tcpCurrentServerConnectionsEstablished = prometheus.NewDesc(
		"tcp_current_server_connections_established",
		"Current server connections in the Established state, which indicates that data transfer can occur between the NetScaler and the server.",
		[]string{
			"ns_instance",
		},
		nil,
	)

	//
	// Metrics about the SNMP exporter itself.
	nsScrapeDuration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "nscollector_scrape_duration_seconds",
			Help: "Duration of scrape by the Netscaler exporter",
		},
		[]string{"ns_instance"},
	)
	nsCollectionRequestErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "nscollector_request_errors_total",
			Help: "Errors in requests to the Netscaler exporter",
		},
	)

	nsCollectorScrapeStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "nscollector_scrape_status",
			Help: "Status of the collector scrape run",
		},
		[]string{"ns_instance"},
	)
)

// Exporter represents the metrics exported to Prometheus
type Exporter struct {
	modelID                                   *prometheus.Desc
	mgmtCPUUsage                              *prometheus.Desc
	memUsage                                  *prometheus.Desc
	pktCPUUsage                               *prometheus.Desc
	flashPartitionUsage                       *prometheus.Desc
	varPartitionUsage                         *prometheus.Desc
	rxMbPerSec                                *prometheus.Desc
	txMbPerSec                                *prometheus.Desc
	httpRequestsRate                          *prometheus.Desc
	httpResponsesRate                         *prometheus.Desc
	tcpCurrentClientConnections               *prometheus.Desc
	tcpCurrentClientConnectionsEstablished    *prometheus.Desc
	tcpCurrentServerConnections               *prometheus.Desc
	tcpCurrentServerConnectionsEstablished    *prometheus.Desc
	interfacesRxBytesPerSecond                *prometheus.GaugeVec
	interfacesTxBytesPerSecond                *prometheus.GaugeVec
	interfacesRxPacketsPerSecond              *prometheus.GaugeVec
	interfacesTxPacketsPerSecond              *prometheus.GaugeVec
	interfacesJumboPacketsRxPerSecond         *prometheus.GaugeVec
	interfacesJumboPacketsTxPerSecond         *prometheus.GaugeVec
	interfacesErrorPacketsRxPerSecond         *prometheus.GaugeVec
	virtualServersWaitingRequests             *prometheus.GaugeVec
	virtualServersHealth                      *prometheus.GaugeVec
	virtualServersInactiveServices            *prometheus.GaugeVec
	virtualServersActiveServices              *prometheus.GaugeVec
	virtualServersTotalHits                   *prometheus.CounterVec
	virtualServersHitsRate                    *prometheus.GaugeVec
	virtualServersTotalRequests               *prometheus.CounterVec
	virtualServersRequestsRate                *prometheus.GaugeVec
	virtualServersTotalResponses              *prometheus.CounterVec
	virtualServersReponsesRate                *prometheus.GaugeVec
	virtualServersTotalRequestBytes           *prometheus.CounterVec
	virtualServersRequestBytesRate            *prometheus.GaugeVec
	virtualServersTotalResponseBytes          *prometheus.CounterVec
	virtualServersReponseBytesRate            *prometheus.GaugeVec
	virtualServersCurrentClientConnections    *prometheus.GaugeVec
	virtualServersCurrentServerConnections    *prometheus.GaugeVec
	servicesThroughput                        *prometheus.CounterVec
	servicesThroughputRate                    *prometheus.GaugeVec
	servicesAvgTTFB                           *prometheus.GaugeVec
	servicesState                             *prometheus.GaugeVec
	servicesTotalRequests                     *prometheus.CounterVec
	servicesRequestsRate                      *prometheus.GaugeVec
	servicesTotalResponses                    *prometheus.CounterVec
	servicesResponsesRate                     *prometheus.GaugeVec
	servicesTotalRequestBytes                 *prometheus.CounterVec
	servicesRequestBytesRate                  *prometheus.GaugeVec
	servicesTotalResponseBytes                *prometheus.CounterVec
	servicesResponseBytesRate                 *prometheus.GaugeVec
	servicesCurrentClientConns                *prometheus.GaugeVec
	servicesSurgeCount                        *prometheus.GaugeVec
	servicesCurrentServerConns                *prometheus.GaugeVec
	servicesServerEstablishedConnections      *prometheus.GaugeVec
	servicesCurrentReusePool                  *prometheus.GaugeVec
	servicesMaxClients                        *prometheus.GaugeVec
	servicesCurrentLoad                       *prometheus.GaugeVec
	servicesVirtualServerServiceHits          *prometheus.CounterVec
	servicesVirtualServerServiceHitsRate      *prometheus.GaugeVec
	servicesActiveTransactions                *prometheus.GaugeVec
	serviceGroupsState                        *prometheus.GaugeVec
	serviceGroupsAvgTTFB                      *prometheus.GaugeVec
	serviceGroupsTotalRequests                *prometheus.CounterVec
	serviceGroupsRequestsRate                 *prometheus.GaugeVec
	serviceGroupsTotalResponses               *prometheus.CounterVec
	serviceGroupsResponsesRate                *prometheus.GaugeVec
	serviceGroupsTotalRequestBytes            *prometheus.CounterVec
	serviceGroupsRequestBytesRate             *prometheus.GaugeVec
	serviceGroupsTotalResponseBytes           *prometheus.CounterVec
	serviceGroupsResponseBytesRate            *prometheus.GaugeVec
	serviceGroupsCurrentClientConnections     *prometheus.GaugeVec
	serviceGroupsSurgeCount                   *prometheus.GaugeVec
	serviceGroupsCurrentServerConnections     *prometheus.GaugeVec
	serviceGroupsServerEstablishedConnections *prometheus.GaugeVec
	serviceGroupsCurrentReusePool             *prometheus.GaugeVec
	serviceGroupsMaxClients                   *prometheus.GaugeVec
	nsScrapeDuration                          *prometheus.SummaryVec
	nsCollectorScrapeStatus                   *prometheus.GaugeVec
	nsCollectionRequestErrors                 prometheus.Counter
	userName                                  *string
	password                                  *string
	url                                       *string
	nsInstance                                string
}

// NewExporter initialises the exporter
func NewExporter(url *string, user *string, pass *string) (*Exporter, error) {
	if *url == "" {
		return nil, errors.New("No Url Specified")
	}
	if *user == "" {
		return nil, errors.New("No Username Specified")
	}
	if *pass == "" {
		return nil, errors.New("No Password Specified")
	}
	return &Exporter{
		modelID:                                   modelID,
		mgmtCPUUsage:                              mgmtCPUUsage,
		memUsage:                                  memUsage,
		pktCPUUsage:                               pktCPUUsage,
		flashPartitionUsage:                       flashPartitionUsage,
		varPartitionUsage:                         varPartitionUsage,
		rxMbPerSec:                                rxMbPerSec,
		txMbPerSec:                                txMbPerSec,
		httpRequestsRate:                          httpRequestsRate,
		httpResponsesRate:                         httpResponsesRate,
		tcpCurrentClientConnections:               tcpCurrentClientConnections,
		tcpCurrentClientConnectionsEstablished:    tcpCurrentClientConnectionsEstablished,
		tcpCurrentServerConnections:               tcpCurrentServerConnections,
		tcpCurrentServerConnectionsEstablished:    tcpCurrentServerConnectionsEstablished,
		interfacesRxBytesPerSecond:                interfacesRxBytesPerSecond,
		interfacesTxBytesPerSecond:                interfacesTxBytesPerSecond,
		interfacesRxPacketsPerSecond:              interfacesRxPacketsPerSecond,
		interfacesTxPacketsPerSecond:              interfacesTxPacketsPerSecond,
		interfacesJumboPacketsRxPerSecond:         interfacesJumboPacketsRxPerSecond,
		interfacesJumboPacketsTxPerSecond:         interfacesJumboPacketsTxPerSecond,
		interfacesErrorPacketsRxPerSecond:         interfacesErrorPacketsRxPerSecond,
		virtualServersWaitingRequests:             virtualServersWaitingRequests,
		virtualServersHealth:                      virtualServersHealth,
		virtualServersInactiveServices:            virtualServersInactiveServices,
		virtualServersActiveServices:              virtualServersActiveServices,
		virtualServersTotalHits:                   virtualServersTotalHits,
		virtualServersHitsRate:                    virtualServersHitsRate,
		virtualServersTotalRequests:               virtualServersTotalRequests,
		virtualServersRequestsRate:                virtualServersRequestsRate,
		virtualServersTotalResponses:              virtualServersTotalResponses,
		virtualServersReponsesRate:                virtualServersReponsesRate,
		virtualServersTotalRequestBytes:           virtualServersTotalRequestBytes,
		virtualServersRequestBytesRate:            virtualServersRequestBytesRate,
		virtualServersTotalResponseBytes:          virtualServersTotalResponseBytes,
		virtualServersReponseBytesRate:            virtualServersReponseBytesRate,
		virtualServersCurrentClientConnections:    virtualServersCurrentClientConnections,
		virtualServersCurrentServerConnections:    virtualServersCurrentServerConnections,
		servicesThroughput:                        servicesThroughput,
		servicesThroughputRate:                    servicesThroughputRate,
		servicesAvgTTFB:                           servicesAvgTTFB,
		servicesState:                             servicesState,
		servicesTotalRequests:                     servicesTotalRequests,
		servicesRequestsRate:                      servicesRequestsRate,
		servicesTotalResponses:                    servicesTotalResponses,
		servicesResponsesRate:                     servicesResponsesRate,
		servicesTotalRequestBytes:                 servicesTotalRequestBytes,
		servicesRequestBytesRate:                  servicesRequestBytesRate,
		servicesTotalResponseBytes:                servicesTotalResponseBytes,
		servicesResponseBytesRate:                 servicesResponseBytesRate,
		servicesCurrentClientConns:                servicesCurrentClientConns,
		servicesSurgeCount:                        servicesSurgeCount,
		servicesCurrentServerConns:                servicesCurrentServerConns,
		servicesServerEstablishedConnections:      servicesServerEstablishedConnections,
		servicesCurrentReusePool:                  servicesCurrentReusePool,
		servicesMaxClients:                        servicesMaxClients,
		servicesCurrentLoad:                       servicesCurrentLoad,
		servicesVirtualServerServiceHits:          servicesVirtualServerServiceHits,
		servicesVirtualServerServiceHitsRate:      servicesVirtualServerServiceHitsRate,
		servicesActiveTransactions:                servicesActiveTransactions,
		serviceGroupsState:                        serviceGroupsState,
		serviceGroupsAvgTTFB:                      serviceGroupsAvgTTFB,
		serviceGroupsTotalRequests:                serviceGroupsTotalRequests,
		serviceGroupsRequestsRate:                 serviceGroupsRequestsRate,
		serviceGroupsTotalResponses:               serviceGroupsTotalResponses,
		serviceGroupsResponsesRate:                serviceGroupsResponsesRate,
		serviceGroupsTotalRequestBytes:            serviceGroupsTotalRequestBytes,
		serviceGroupsRequestBytesRate:             serviceGroupsRequestBytesRate,
		serviceGroupsTotalResponseBytes:           serviceGroupsTotalResponseBytes,
		serviceGroupsResponseBytesRate:            serviceGroupsResponseBytesRate,
		serviceGroupsCurrentClientConnections:     serviceGroupsCurrentClientConnections,
		serviceGroupsSurgeCount:                   serviceGroupsSurgeCount,
		serviceGroupsCurrentServerConnections:     serviceGroupsCurrentServerConnections,
		serviceGroupsServerEstablishedConnections: serviceGroupsServerEstablishedConnections,
		serviceGroupsCurrentReusePool:             serviceGroupsCurrentReusePool,
		serviceGroupsMaxClients:                   serviceGroupsMaxClients,
		nsScrapeDuration:                          nsScrapeDuration,
		nsCollectionRequestErrors:                 nsCollectionRequestErrors,
		nsCollectorScrapeStatus:                   nsCollectorScrapeStatus,
		userName:                                  user,
		password:                                  pass,
		url:                                       url,
	}, nil
}

// Describe implements Collector
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- modelID
	ch <- mgmtCPUUsage
	ch <- memUsage
	ch <- pktCPUUsage
	ch <- flashPartitionUsage
	ch <- varPartitionUsage
	ch <- rxMbPerSec
	ch <- txMbPerSec
	ch <- httpRequestsRate
	ch <- httpResponsesRate
	ch <- tcpCurrentClientConnections
	ch <- tcpCurrentClientConnectionsEstablished
	ch <- tcpCurrentServerConnections
	ch <- tcpCurrentServerConnectionsEstablished

	e.interfacesRxBytesPerSecond.Describe(ch)
	e.interfacesTxBytesPerSecond.Describe(ch)
	e.interfacesRxPacketsPerSecond.Describe(ch)
	e.interfacesTxPacketsPerSecond.Describe(ch)
	e.interfacesJumboPacketsRxPerSecond.Describe(ch)
	e.interfacesJumboPacketsTxPerSecond.Describe(ch)
	e.interfacesErrorPacketsRxPerSecond.Describe(ch)

	e.virtualServersWaitingRequests.Describe(ch)
	e.virtualServersHealth.Describe(ch)
	e.virtualServersInactiveServices.Describe(ch)
	e.virtualServersActiveServices.Describe(ch)
	e.virtualServersTotalHits.Describe(ch)
	e.virtualServersHitsRate.Describe(ch)
	e.virtualServersTotalRequests.Describe(ch)
	e.virtualServersRequestsRate.Describe(ch)
	e.virtualServersTotalResponses.Describe(ch)
	e.virtualServersReponsesRate.Describe(ch)
	e.virtualServersTotalRequestBytes.Describe(ch)
	e.virtualServersRequestBytesRate.Describe(ch)
	e.virtualServersTotalResponseBytes.Describe(ch)
	e.virtualServersReponseBytesRate.Describe(ch)
	e.virtualServersCurrentClientConnections.Describe(ch)
	e.virtualServersCurrentServerConnections.Describe(ch)

	e.servicesThroughput.Describe(ch)
	e.servicesThroughputRate.Describe(ch)
	e.servicesAvgTTFB.Describe(ch)
	e.servicesState.Describe(ch)
	e.servicesTotalRequests.Describe(ch)
	e.servicesRequestsRate.Describe(ch)
	e.servicesTotalResponses.Describe(ch)
	e.servicesResponsesRate.Describe(ch)
	e.servicesTotalRequestBytes.Describe(ch)
	e.servicesRequestBytesRate.Describe(ch)
	e.servicesTotalResponseBytes.Describe(ch)
	e.servicesResponseBytesRate.Describe(ch)
	e.servicesCurrentClientConns.Describe(ch)
	e.servicesSurgeCount.Describe(ch)
	e.servicesCurrentServerConns.Describe(ch)
	e.servicesServerEstablishedConnections.Describe(ch)
	e.servicesCurrentReusePool.Describe(ch)
	e.servicesMaxClients.Describe(ch)
	e.servicesCurrentLoad.Describe(ch)
	e.servicesVirtualServerServiceHits.Describe(ch)
	e.servicesVirtualServerServiceHitsRate.Describe(ch)
	e.servicesActiveTransactions.Describe(ch)

	e.serviceGroupsState.Describe(ch)
	e.serviceGroupsAvgTTFB.Describe(ch)
	e.serviceGroupsTotalRequests.Describe(ch)
	e.serviceGroupsRequestsRate.Describe(ch)
	e.serviceGroupsTotalResponses.Describe(ch)
	e.serviceGroupsResponsesRate.Describe(ch)
	e.serviceGroupsTotalRequestBytes.Describe(ch)
	e.serviceGroupsRequestBytesRate.Describe(ch)
	e.serviceGroupsTotalResponseBytes.Describe(ch)
	e.serviceGroupsResponseBytesRate.Describe(ch)
	e.serviceGroupsCurrentClientConnections.Describe(ch)
	e.serviceGroupsSurgeCount.Describe(ch)
	e.serviceGroupsCurrentServerConnections.Describe(ch)
	e.serviceGroupsServerEstablishedConnections.Describe(ch)
	e.serviceGroupsCurrentReusePool.Describe(ch)
	e.serviceGroupsMaxClients.Describe(ch)

	e.nsScrapeDuration.Describe(ch)
	e.nsCollectionRequestErrors.Describe(ch)
}

func (e *Exporter) recordScrapeDuration(start time.Time) {

	duration := float64(time.Since(start).Seconds())
	e.nsScrapeDuration.WithLabelValues(e.nsInstance).Observe(duration)
	logger.Debugf("Scrape of target '%s' took %f seconds", e.nsInstance, duration)
}

// Collect is initiated by the Prometheus handler and gathers the metrics
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	// get start time for scrape stats
	scrapeStart := time.Now()

	nsClient, err := netscaler.NewNitroClient(*e.url, *e.userName, *e.password)
	if err != nil {
		e.nsCollectionRequestErrors.Inc()
		e.nsCollectorScrapeStatus.WithLabelValues(e.nsInstance).Set(float64(0))
		e.nsCollectorScrapeStatus.Collect(ch)
		e.nsCollectionRequestErrors.Collect(ch)
		e.recordScrapeDuration(scrapeStart)
		e.nsScrapeDuration.Collect(ch)
		logger.Error(err)
		return
	}

	err = netscaler.Connect(nsClient)
	if err != nil {
		e.nsCollectionRequestErrors.Inc()
		e.nsCollectorScrapeStatus.WithLabelValues(e.nsInstance).Set(float64(0))
		e.nsCollectorScrapeStatus.Collect(ch)
		e.nsCollectionRequestErrors.Collect(ch)
		e.recordScrapeDuration(scrapeStart)
		e.nsScrapeDuration.Collect(ch)
		logger.Error(err)
		return
	}

	nslicense, err := netscaler.GetNSLicense(nsClient, "")
	if err != nil {
		e.nsCollectionRequestErrors.Inc()
		logger.Error(err)
	}

	ns, err := netscaler.GetNSStats(nsClient, "")
	if err != nil {
		e.nsCollectionRequestErrors.Inc()
		logger.Error(err)
	}

	interfaces, err := netscaler.GetInterfaceStats(nsClient, "")
	if err != nil {
		e.nsCollectionRequestErrors.Inc()
		logger.Error(err)
	}

	virtualServers, err := netscaler.GetVirtualServerStats(nsClient, "")
	if err != nil {
		e.nsCollectionRequestErrors.Inc()
		logger.Error(err)
	}

	services, err := netscaler.GetServiceStats(nsClient, "")
	if err != nil {
		e.nsCollectionRequestErrors.Inc()
		logger.Error(err)
	}

	fltModelID, _ := strconv.ParseFloat(nslicense.NSLicense.ModelID, 64)

	fltTCPCurrentClientConnections, _ := strconv.ParseFloat(ns.NSStats.TCPCurrentClientConnections, 64)
	fltTCPCurrentClientConnectionsEstablished, _ := strconv.ParseFloat(ns.NSStats.TCPCurrentClientConnectionsEstablished, 64)
	fltTCPCurrentServerConnections, _ := strconv.ParseFloat(ns.NSStats.TCPCurrentServerConnections, 64)
	fltTCPCurrentServerConnectionsEstablished, _ := strconv.ParseFloat(ns.NSStats.TCPCurrentServerConnectionsEstablished, 64)

	ch <- prometheus.MustNewConstMetric(
		modelID, prometheus.GaugeValue, fltModelID, e.nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		mgmtCPUUsage, prometheus.GaugeValue, ns.NSStats.MgmtCPUUsagePcnt, e.nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		memUsage, prometheus.GaugeValue, ns.NSStats.MemUsagePcnt, e.nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		pktCPUUsage, prometheus.GaugeValue, ns.NSStats.PktCPUUsagePcnt, e.nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		flashPartitionUsage, prometheus.GaugeValue, ns.NSStats.FlashPartitionUsage, e.nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		varPartitionUsage, prometheus.GaugeValue, ns.NSStats.VarPartitionUsage, e.nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		rxMbPerSec, prometheus.GaugeValue, ns.NSStats.ReceivedMbPerSecond, e.nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		txMbPerSec, prometheus.GaugeValue, ns.NSStats.TransmitMbPerSecond, e.nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		httpRequestsRate, prometheus.GaugeValue, ns.NSStats.HTTPRequestsRate, e.nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		httpResponsesRate, prometheus.GaugeValue, ns.NSStats.HTTPResponsesRate, e.nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		tcpCurrentClientConnections, prometheus.GaugeValue, fltTCPCurrentClientConnections, e.nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		tcpCurrentClientConnectionsEstablished, prometheus.GaugeValue, fltTCPCurrentClientConnectionsEstablished, e.nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		tcpCurrentServerConnections, prometheus.GaugeValue, fltTCPCurrentServerConnections, e.nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		tcpCurrentServerConnectionsEstablished, prometheus.GaugeValue, fltTCPCurrentServerConnectionsEstablished, e.nsInstance,
	)

	e.collectInterfacesRxBytesPerSecond(interfaces)
	e.interfacesRxBytesPerSecond.Collect(ch)

	e.collectInterfacesTxBytesPerSecond(interfaces)
	e.interfacesTxBytesPerSecond.Collect(ch)

	e.collectInterfacesRxPacketsPerSecond(interfaces)
	e.interfacesRxPacketsPerSecond.Collect(ch)

	e.collectInterfacesTxPacketsPerSecond(interfaces)
	e.interfacesTxPacketsPerSecond.Collect(ch)

	e.collectInterfacesJumboPacketsRxPerSecond(interfaces)
	e.interfacesJumboPacketsRxPerSecond.Collect(ch)

	e.collectInterfacesJumboPacketsTxPerSecond(interfaces)
	e.interfacesJumboPacketsTxPerSecond.Collect(ch)

	e.collectInterfacesErrorPacketsRxPerSecond(interfaces)
	e.interfacesErrorPacketsRxPerSecond.Collect(ch)

	e.collectVirtualServerWaitingRequests(virtualServers)
	e.virtualServersWaitingRequests.Collect(ch)

	e.collectVirtualServerHealth(virtualServers)
	e.virtualServersHealth.Collect(ch)

	e.collectVirtualServerInactiveServices(virtualServers)
	e.virtualServersInactiveServices.Collect(ch)

	e.collectVirtualServerActiveServices(virtualServers)
	e.virtualServersActiveServices.Collect(ch)

	e.collectVirtualServerTotalHits(virtualServers)
	e.virtualServersTotalHits.Collect(ch)

	e.collectVirtualServerHitsRate(virtualServers)
	e.virtualServersHitsRate.Collect(ch)

	e.collectVirtualServerTotalRequests(virtualServers)
	e.virtualServersTotalRequests.Collect(ch)

	e.collectVirtualServerRequestsRate(virtualServers)
	e.virtualServersRequestsRate.Collect(ch)

	e.collectVirtualServerTotalResponses(virtualServers)
	e.virtualServersTotalResponses.Collect(ch)

	e.collectVirtualServerResponsesRate(virtualServers)
	e.virtualServersReponsesRate.Collect(ch)

	e.collectVirtualServerTotalRequestBytes(virtualServers)
	e.virtualServersTotalRequestBytes.Collect(ch)

	e.collectVirtualServerRequestBytesRate(virtualServers)
	e.virtualServersRequestBytesRate.Collect(ch)

	e.collectVirtualServerTotalResponseBytes(virtualServers)
	e.virtualServersTotalResponseBytes.Collect(ch)

	e.collectVirtualServerResponseBytesRate(virtualServers)
	e.virtualServersReponseBytesRate.Collect(ch)

	e.collectVirtualServerCurrentClientConnections(virtualServers)
	e.virtualServersCurrentClientConnections.Collect(ch)

	e.collectVirtualServerCurrentServerConnections(virtualServers)
	e.virtualServersCurrentServerConnections.Collect(ch)

	e.collectServicesThroughput(services)
	e.servicesThroughput.Collect(ch)

	e.collectServicesThroughputRate(services)
	e.servicesThroughputRate.Collect(ch)

	e.collectServicesAvgTTFB(services)
	e.servicesAvgTTFB.Collect(ch)

	e.collectServicesState(services)
	e.servicesState.Collect(ch)

	e.collectServicesTotalRequests(services)
	e.servicesTotalRequests.Collect(ch)

	e.collectServicesRequestsRate(services)
	e.servicesRequestsRate.Collect(ch)

	e.collectServicesTotalResponses(services)
	e.servicesTotalResponses.Collect(ch)

	e.collectServicesResponsesRate(services)
	e.servicesResponsesRate.Collect(ch)

	e.collectServicesTotalRequestBytes(services)
	e.servicesTotalRequestBytes.Collect(ch)

	e.collectServicesRequestBytesRate(services)
	e.servicesRequestBytesRate.Collect(ch)

	e.collectServicesTotalResponseBytes(services)
	e.servicesTotalResponseBytes.Collect(ch)

	e.collectServicesResponseBytesRate(services)
	e.servicesResponseBytesRate.Collect(ch)

	e.collectServicesCurrentClientConns(services)
	e.servicesCurrentClientConns.Collect(ch)

	e.collectServicesSurgeCount(services)
	e.servicesSurgeCount.Collect(ch)

	e.collectServicesCurrentServerConns(services)
	e.servicesCurrentServerConns.Collect(ch)

	e.collectServicesServerEstablishedConnections(services)
	e.servicesServerEstablishedConnections.Collect(ch)

	e.collectServicesCurrentReusePool(services)
	e.servicesCurrentReusePool.Collect(ch)

	e.collectServicesMaxClients(services)
	e.servicesMaxClients.Collect(ch)

	e.collectServicesCurrentLoad(services)
	e.servicesCurrentLoad.Collect(ch)

	e.collectServicesVirtualServerServiceHits(services)
	e.servicesVirtualServerServiceHits.Collect(ch)

	e.collectServicesVirtualServerServiceHitsRate(services)
	e.servicesVirtualServerServiceHitsRate.Collect(ch)

	e.collectServicesActiveTransactions(services)
	e.servicesActiveTransactions.Collect(ch)

	servicegroups, err := netscaler.GetServiceGroups(nsClient, "attrs=servicegroupname")
	if err != nil {
		logger.Error(err)
	}

	for _, sg := range servicegroups.ServiceGroups {
		bindings, err2 := netscaler.GetServiceGroupMemberBindings(nsClient, sg.Name)
		if err2 != nil {
			logger.Error(err2)
		}

		for _, member := range bindings.ServiceGroupMemberBindings {
			// NetScaler API has a bug which means it throws errors if you try to retrieve stats for a wildcard port (* in GUI, 65535 in API and CLI).
			// Until Citrix resolve the issue we skip attempting to retrieve stats for those service groups.
			if member.Port != 65535 {
				port := strconv.FormatInt(member.Port, 10)

				qs := "args=servicegroupname:" + sg.Name + ",servername:" + member.ServerName + ",port:" + port
				stats, err2 := netscaler.GetServiceGroupMemberStats(nsClient, qs)
				if err2 != nil {
					logger.Error(err2)
				}

				e.collectServiceGroupsState(stats, sg.Name, member.ServerName)
				e.serviceGroupsState.Collect(ch)

				e.collectServiceGroupsAvgTTFB(stats, sg.Name, member.ServerName)
				e.serviceGroupsAvgTTFB.Collect(ch)

				e.collectServiceGroupsTotalRequests(stats, sg.Name, member.ServerName)
				e.serviceGroupsTotalRequests.Collect(ch)

				e.collectServiceGroupsRequestsRate(stats, sg.Name, member.ServerName)
				e.serviceGroupsRequestsRate.Collect(ch)

				e.collectServiceGroupsTotalResponses(stats, sg.Name, member.ServerName)
				e.serviceGroupsTotalResponses.Collect(ch)

				e.collectServiceGroupsResponsesRate(stats, sg.Name, member.ServerName)
				e.serviceGroupsResponsesRate.Collect(ch)

				e.collectServiceGroupsTotalRequestBytes(stats, sg.Name, member.ServerName)
				e.serviceGroupsTotalRequestBytes.Collect(ch)

				e.collectServiceGroupsRequestBytesRate(stats, sg.Name, member.ServerName)
				e.serviceGroupsRequestBytesRate.Collect(ch)

				e.collectServiceGroupsTotalResponseBytes(stats, sg.Name, member.ServerName)
				e.serviceGroupsTotalResponseBytes.Collect(ch)

				e.collectServiceGroupsResponseBytesRate(stats, sg.Name, member.ServerName)
				e.serviceGroupsResponseBytesRate.Collect(ch)

				e.collectServiceGroupsCurrentClientConnections(stats, sg.Name, member.ServerName)
				e.serviceGroupsCurrentClientConnections.Collect(ch)

				e.collectServiceGroupsSurgeCount(stats, sg.Name, member.ServerName)
				e.serviceGroupsSurgeCount.Collect(ch)

				e.collectServiceGroupsCurrentServerConnections(stats, sg.Name, member.ServerName)
				e.serviceGroupsCurrentServerConnections.Collect(ch)

				e.collectServiceGroupsServerEstablishedConnections(stats, sg.Name, member.ServerName)
				e.serviceGroupsServerEstablishedConnections.Collect(ch)

				e.collectServiceGroupsCurrentReusePool(stats, sg.Name, member.ServerName)
				e.serviceGroupsCurrentReusePool.Collect(ch)

				e.collectServiceGroupsMaxClients(stats, sg.Name, member.ServerName)
				e.serviceGroupsMaxClients.Collect(ch)
			}
		}
	}

	err = netscaler.Disconnect(nsClient)
	if err != nil {
		e.nsCollectionRequestErrors.Inc()
		logger.Error(err)
	}
	e.recordScrapeDuration(scrapeStart)
	e.nsScrapeDuration.Collect(ch)
	e.nsCollectionRequestErrors.Collect(ch)
	e.nsCollectorScrapeStatus.WithLabelValues(e.nsInstance).Set(float64(1))
	e.nsCollectorScrapeStatus.Collect(ch)
}
