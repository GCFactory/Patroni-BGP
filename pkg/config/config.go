package config

import "github.com/GCFactory/Patroni-BGP/pkg/bgp"

// Config defines all of the settings for the Patroni-bgp Pod
type Config struct {
	// Logging, settings
	Logging int `yaml:"logging"`

	// AddPeersAsBackends, this will automatically add RAFT peers as backends to a loadbalancer
	AddPeersAsBackends bool `yaml:"addPeersAsBackends"`

	// VIP is the Virtual IP address exposed for the cluster (TODO: deprecate)
	VIP string `yaml:"vip"`

	// VipSubnet is the Subnet that is applied to the VIP
	VIPSubnet string `yaml:"vipSubnet"`

	// VIPCIDR is cidr range for the VIP (primarily needed for BGP)
	VIPCIDR string `yaml:"vipCidr"`

	// MasterAddress is the IP or DNS Name to use as a VirtualIP
	MasterAddress string `yaml:"address"`

	ReplicaAddress string `yaml:"replicaAddress"`

	// Listen port for the VirtualIP
	Port int `yaml:"port"`

	// SingleNode will start the cluster as a single Node (Raft disabled)
	SingleNode bool `yaml:"singleNode"`

	// Interface is the network interface to bind to (default: First Adapter)
	Interface string `yaml:"interface,omitempty"`

	// ServicesInterface is the network interface to bind to for services (optional)
	ServicesInterface string `yaml:"servicesInterface,omitempty"`

	// Routing Table ID for when using routing table mode
	RoutingTableID int `yaml:"routingTableID"`

	// Routing Table Type, what sort of route should be added to the routing table
	RoutingTableType int `yaml:"routingTableType"`

	// Routing Protocol, value that will be used as protocol when creating rutes
	RoutingProtocol int `yaml:"routingProtocol"`

	// Clean routing table of redundant routes on start
	CleanRoutingTable bool `yaml:"cleanRoutingTable"`

	// BGP Configuration
	BGPConfig     bgp.Config
	BGPPeerConfig bgp.Peer
	BGPPeers      []string

	// ProviderConfig, is the path to a provider configuration file
	ProviderConfig string

	// The hostport used to expose Prometheus metrics over an HTTP server
	PrometheusHTTPServer string `yaml:"prometheusHTTPServer,omitempty"`

	// Egress configuration

	// EgressPodCidr, this contains the pod cidr range to ignore Egress
	EgressPodCidr string

	// EgressServiceCidr, this contains the service cidr range to ignore
	EgressServiceCidr string

	// EgressWithNftables, this will use the iptables-nftables OVER iptables
	EgressWithNftables bool

	// DNSMode, this will set the mode DSN lookup will be performed (first, ipv4, ipv6, dual)
	DNSMode string `yaml:"dnsDualStackMode"`

	PatroniUrl string `yaml:"patroniAddress"`
}
