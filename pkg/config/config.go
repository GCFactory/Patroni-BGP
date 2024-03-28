package config

import "github.com/GCFactory/Patroni-BGP/pkg/bgp"

// Config defines all of the settings for the Kube-Vip Pod
type Config struct {
	// Logging, settings
	Logging int `yaml:"logging"`

	// EnableARP, will use ARP to advertise the VIP address
	EnableARP bool `yaml:"enableARP"`

	// EnableBGP, will use BGP to advertise the VIP address
	EnableBGP bool `yaml:"enableBGP"`

	// EnableRoutingTable, will use the routing table to advertise the VIP address
	EnableRoutingTable bool `yaml:"enableRoutingTable"`

	// EnableServiceSecurity, will enable the use of iptables to secure services
	EnableServiceSecurity bool `yaml:"EnableServiceSecurity"`

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

	// StartAsLeader, this will start this node as the leader before other nodes connect
	StartAsLeader bool `yaml:"startAsLeader"`

	// Interface is the network interface to bind to (default: First Adapter)
	Interface string `yaml:"interface,omitempty"`

	// ServicesInterface is the network interface to bind to for services (optional)
	ServicesInterface string `yaml:"servicesInterface,omitempty"`

	// EnableLoadBalancer, provides the flexibility to make the load-balancer optional
	EnableLoadBalancer bool `yaml:"enableLoadBalancer"`

	// Listen port for the IPVS Service
	LoadBalancerPort int `yaml:"lbPort"`

	// Forwarding method for the IPVS Service
	LoadBalancerForwardingMethod string `yaml:"lbForwardingMethod"`

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

	// ServicesLeaseName, this will set the lease name for services leader in arp mode
	ServicesLeaseName string `yaml:"servicesLeaseName"`

	// DNSMode, this will set the mode DSN lookup will be performed (first, ipv4, ipv6, dual)
	DNSMode string `yaml:"dnsDualStackMode"`

	PatroniUrl string `yaml:"patroniAddress"`
}

// KubernetesLeaderElection defines all of the settings for Kubernetes KubernetesLeaderElection
type KubernetesLeaderElection struct {
	// EnableLeaderElection will use the Kubernetes leader election algorithm
	EnableLeaderElection bool `yaml:"enableLeaderElection"`

	// LeaseName - name of the lease for leader election
	LeaseName string `yaml:"leaseName"`

	// Lease Duration - length of time a lease can be held for
	LeaseDuration int

	// RenewDeadline - length of time a host can attempt to renew its lease
	RenewDeadline int

	// RetryPerion - Number of times the host will retry to hold a lease
	RetryPeriod int

	// LeaseAnnotations - annotations which will be given to the lease object
	LeaseAnnotations map[string]string
}

// Etcd defines all the settings for the etcd client.
type Etcd struct {
	CAFile         string
	ClientCertFile string
	ClientKeyFile  string
	Endpoints      []string
}

// LoadBalancer contains the configuration of a load balancing instance
type LoadBalancer struct {
	// Name of a LoadBalancer
	Name string `yaml:"name"`

	// Type of LoadBalancer, either TCP of HTTP(s)
	Type string `yaml:"type"`

	// Listening frontend port of this LoadBalancer instance
	Port int `yaml:"port"`

	// BindToVip will bind the load balancer port to the VIP itself
	BindToVip bool `yaml:"bindToVip"`

	// Forwarding method of LoadBalancer, either Local, Tunnel, DirectRoute or Bypass
	ForwardingMethod string `yaml:"forwardingMethod"`
}
