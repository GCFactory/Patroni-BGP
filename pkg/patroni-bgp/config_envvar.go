package patroni_bgp

// Environment variables
const (
	// vipLogLevel - defines the level of logging to produce (5 being the most verbose)
	vipLogLevel = "vip_loglevel"

	// patroniUrl - defines URL of patroni REST API
	patroniUrl = "patroni_url"

	// primaryAddress - defines what address will be announced for Primary instance
	primaryAddress = "patroni_primary_address"
	// syncReplicaAddress - defines what address will be announced for replica instance in sync mode
	syncReplicaAddress = "patroni_syncreplica_address"
	// asyncReplicaAddress - defines what address will be announced for replica instance in async mode
	asyncReplicaAddress = "patroni_asyncreplica_address"

	// vipInterface - defines the interface that the vip should bind too
	vipInterface = "vip_interface"

	// vipServicesInterface - defines the interface that the service vips should bind too
	vipServicesInterface = "vip_servicesinterface"

	vipCidr = "vip_cidr"

	// vipCidr - defines the cidr that the vip will use (for BGP)
	// vipSubnet - defines the subnet that the vip will use
	/////////////////////////////////////
	// TO DO:
	// Determine how to tidy this mess up
	/////////////////////////////////////
	// vipAddress - defines the address that the vip will expose
	// DEPRECATED: will be removed in a next release

	vipSubnet  = "vip_subnet"
	vipAddress = "vip_address"

	// bgpRouterID defines the routerID for the BGP server
	bgpRouterID = "bgp_routerid"
	// bgpRouterInterface defines the interface that we can find the address for
	bgpRouterInterface = "bgp_routerinterface"
	// bgpRouterAS defines the AS for the BGP server
	bgpRouterAS = "bgp_as"
	// bgpPeerAddress defines the address for a BGP peer
	bgpPeerAddress = "bgp_peeraddress"
	// bgpPeers defines the address for a BGP peer
	bgpPeers = "bgp_peers"
	// bgpPeerAS defines the AS for a BGP peer
	bgpPeerAS = "bgp_peeras"
	// bgpPeerAS defines the AS for a BGP peer
	bgpPeerPassword = "bgp_peerpass" // nolint
	// bgpMultiHop enables mulithop routing
	bgpMultiHop = "bgp_multihop"
	// bgpSourceIF defines the source interface for BGP peering
	bgpSourceIF = "bgp_sourceif"
	// bgpSourceIP defines the source address for BGP peering
	bgpSourceIP = "bgp_sourceip"
	// bgpHoldTime defines bgp timers hold time
	bgpHoldTime = "bgp_hold_time"
	// bgpKeepaliveInterval defines bgp timers keepalive interval
	bgpKeepaliveInterval = "bgp_keepalive_interval"

	// vipRoutingTable - defines if table mode will be used for vips
	vipRoutingTable = "vip_routingtable" //nolint

	// vipRoutingTableID - defines which table mode will be used for vips
	vipRoutingTableID = "vip_routingtableid" //nolint

	// vipRoutingTableType - defines which table type will be used for vip routes
	// 						 valid values for this variable can be found in:
	//						 https://pkg.go.dev/golang.org/x/sys/unix#RTN_UNSPEC
	//						 Note that route type have the prefix `RTN_`, and you
	//						 specify the integer value, not the name. For example:
	//						 you should say `vip_routingtabletype=2` for RTN_LOCAL
	vipRoutingTableType = "vip_routingtabletype" //nolint

	// vipRoutingProtocol - defines what value will be used as protocol when creating routes
	vipRoutingProtocol = "vip_routingprotocol" //nolint

	// vipCleanRoutingTable - defines if routing table will be cleaned of redundant routes on patroni-bgp's start
	vipCleanRoutingTable = "vip_cleanroutingtable" //nolint

	// prometheusServer defines the address prometheus listens on
	prometheusServer = "prometheus_server"

	// dnsMode defines mode that DNS lookup will be performed with (first, ipv4, ipv6, dual)
	dnsMode = "dns_mode"
)
