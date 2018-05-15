package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rokett/citrix-netscaler-exporter/netscaler"
)

var (
	interfacesRxBytesPerSecond = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_received_bytes_per_second",
			Help: "Number of bytes received per second by specific interfaces",
		},
		[]string{
			"ns_instance",
			"interface",
			"alias",
		},
	)

	interfacesTxBytesPerSecond = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_transmitted_bytes_per_second",
			Help: "Number of bytes transmitted per second by specific interfaces",
		},
		[]string{
			"ns_instance",
			"interface",
			"alias",
		},
	)

	interfacesRxPacketsPerSecond = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_received_packets_per_second",
			Help: "Number of packets received per second by specific interfaces",
		},
		[]string{
			"ns_instance",
			"interface",
			"alias",
		},
	)

	interfacesTxPacketsPerSecond = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_transmitted_packets_per_second",
			Help: "Number of packets transmitted per second by specific interfaces",
		},
		[]string{
			"ns_instance",
			"interface",
			"alias",
		},
	)

	interfacesJumboPacketsRxPerSecond = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_jumbo_packets_received_per_second",
			Help: "Number of bytes received per second by specific interfaces",
		},
		[]string{
			"ns_instance",
			"interface",
			"alias",
		},
	)

	interfacesJumboPacketsTxPerSecond = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_jumbo_packets_transmitted_per_second",
			Help: "Number of jumbo packets transmitted per second by specific interfaces",
		},
		[]string{
			"ns_instance",
			"interface",
			"alias",
		},
	)

	interfacesErrorPacketsRxPerSecond = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "interfaces_error_packets_received_per_second",
			Help: "Number of error packets received per second by specific interfaces",
		},
		[]string{
			"ns_instance",
			"interface",
			"alias",
		},
	)
)

func (e *Exporter) collectInterfacesRxBytesPerSecond(ns netscaler.NSAPIResponse) {
	e.interfacesRxBytesPerSecond.Reset()

	for _, iface := range ns.InterfaceStats {
		e.interfacesRxBytesPerSecond.WithLabelValues(e.nsInstance, iface.ID, iface.Alias).Set(iface.ReceivedBytesPerSecond)
	}
}

func (e *Exporter) collectInterfacesTxBytesPerSecond(ns netscaler.NSAPIResponse) {
	e.interfacesTxBytesPerSecond.Reset()

	for _, iface := range ns.InterfaceStats {
		e.interfacesTxBytesPerSecond.WithLabelValues(e.nsInstance, iface.ID, iface.Alias).Set(iface.TransmitBytesPerSecond)
	}
}

func (e *Exporter) collectInterfacesRxPacketsPerSecond(ns netscaler.NSAPIResponse) {
	e.interfacesRxPacketsPerSecond.Reset()

	for _, iface := range ns.InterfaceStats {
		e.interfacesRxPacketsPerSecond.WithLabelValues(e.nsInstance, iface.ID, iface.Alias).Set(iface.ReceivedPacketsPerSecond)
	}
}

func (e *Exporter) collectInterfacesTxPacketsPerSecond(ns netscaler.NSAPIResponse) {
	e.interfacesTxPacketsPerSecond.Reset()

	for _, iface := range ns.InterfaceStats {
		e.interfacesTxPacketsPerSecond.WithLabelValues(e.nsInstance, iface.ID, iface.Alias).Set(iface.TransmitPacketsPerSecond)
	}
}

func (e *Exporter) collectInterfacesJumboPacketsRxPerSecond(ns netscaler.NSAPIResponse) {
	e.interfacesJumboPacketsRxPerSecond.Reset()

	for _, iface := range ns.InterfaceStats {
		e.interfacesJumboPacketsRxPerSecond.WithLabelValues(e.nsInstance, iface.ID, iface.Alias).Set(iface.JumboPacketsReceivedPerSecond)
	}
}

func (e *Exporter) collectInterfacesJumboPacketsTxPerSecond(ns netscaler.NSAPIResponse) {
	e.interfacesJumboPacketsTxPerSecond.Reset()

	for _, iface := range ns.InterfaceStats {
		e.interfacesJumboPacketsTxPerSecond.WithLabelValues(e.nsInstance, iface.ID, iface.Alias).Set(iface.JumboPacketsTransmittedPerSecond)
	}
}

func (e *Exporter) collectInterfacesErrorPacketsRxPerSecond(ns netscaler.NSAPIResponse) {
	e.interfacesErrorPacketsRxPerSecond.Reset()

	for _, iface := range ns.InterfaceStats {
		e.interfacesErrorPacketsRxPerSecond.WithLabelValues(e.nsInstance, iface.ID, iface.Alias).Set(iface.ErrorPacketsReceivedPerSecond)
	}
}
