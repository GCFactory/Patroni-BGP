package patroniBgp

import (
	"os"
	"strconv"

	"github.com/GCFactory/Patroni-BGP/pkg/bgp"
	"github.com/GCFactory/Patroni-BGP/pkg/detector"
)

// ParseEnvironment - will popultate the configuration from environment variables
func ParseEnvironment(c *Config) error {
	if c == nil {
		return nil
	}
	// Ensure that logging is set through the environment variables
	env := os.Getenv(vipLogLevel)
	// Set default value
	if env == "" {
		env = "4"
	}

	if env != "" {
		logLevel, err := strconv.ParseUint(env, 10, 32)
		if err != nil {
			panic("Unable to parse environment variable [vip_loglevel], should be int")
		}
		c.Logging = int(logLevel)
	}

	// Find interface
	env = os.Getenv(vipInterface)
	if env != "" {
		c.Interface = env
	}

	// Find interface
	env = os.Getenv(patroniURL)
	if env != "" {
		c.PatroniURL = env
	}

	// Find interface
	env = os.Getenv(primaryAddress)
	if env != "" {
		c.PrimaryAddress = env
	}

	// Find interface
	env = os.Getenv(syncReplicaAddress)
	if env != "" {
		c.SyncReplicaAddress = env
	}

	// Find interface
	env = os.Getenv(asyncReplicaAddress)
	if env != "" {
		c.AsyncReplicaAddress = env
	}

	// Find (services) interface
	env = os.Getenv(vipServicesInterface)
	if env != "" {
		c.ServicesInterface = env
	}

	// Find vip address
	env = os.Getenv(vipAddress)
	if env != "" {
		// TODO - parse address net.Host()
		c.VIP = env
		// } else {
		// 	c.VIP = os.Getenv(address)
	}

	// Find vip address cidr range
	env = os.Getenv(vipCidr)
	if env != "" {
		c.VIPCIDR = env
	}

	// Find vip address subnet
	env = os.Getenv(vipSubnet)
	if env != "" {
		c.VIPSubnet = env
	}

	// Routing Table ID
	env = os.Getenv(vipRoutingTableID)
	if env != "" {
		i, err := strconv.ParseInt(env, 10, 32)
		if err != nil {
			return err
		}
		c.RoutingTableID = int(i)
	}

	// Routing Table Type
	env = os.Getenv(vipRoutingTableType)
	if env != "" {
		i, err := strconv.ParseInt(env, 10, 32)
		if err != nil {
			return err
		}
		c.RoutingTableType = int(i)
	}

	// Routing protocol
	env = os.Getenv(vipRoutingProtocol)
	if env != "" {
		i, err := strconv.ParseInt(env, 10, 32)
		if err != nil {
			return err
		}
		c.RoutingProtocol = int(i)
	}

	// Clean routing table
	env = os.Getenv(vipCleanRoutingTable)
	if env != "" {
		b, err := strconv.ParseBool(env)
		if err != nil {
			return err
		}
		c.CleanRoutingTable = b
	}

	// DNS mode
	env = os.Getenv(dnsMode)
	if env != "" {
		c.DNSMode = env
	}

	// BGP Router interface determines an interface that we can use to find an address for
	env = os.Getenv(bgpRouterInterface)
	if env != "" {
		_, address, err := detector.FindIPAddress(env)
		if err != nil {
			return err
		}
		c.BGPConfig.RouterID = address
	}

	// RouterID
	env = os.Getenv(bgpRouterID)
	if env != "" {
		c.BGPConfig.RouterID = env
	}

	// AS
	env = os.Getenv(bgpRouterAS)
	if env != "" {
		u64, err := strconv.ParseUint(env, 10, 32)
		if err != nil {
			return err
		}
		c.BGPConfig.AS = uint32(u64)
	}

	// Peer AS
	env = os.Getenv(bgpPeerAS)
	if env != "" {
		u64, err := strconv.ParseUint(env, 10, 32)
		if err != nil {
			return err
		}
		c.BGPPeerConfig.AS = uint32(u64)
	}

	// Peer AS
	env = os.Getenv(bgpPeers)
	if env != "" {
		peers, err := bgp.ParseBGPPeerConfig(env)
		if err != nil {
			return err
		}
		c.BGPConfig.Peers = peers
	}

	// BGP Peer mutlihop
	env = os.Getenv(bgpMultiHop)
	if env != "" {
		b, err := strconv.ParseBool(env)
		if err != nil {
			return err
		}
		c.BGPPeerConfig.MultiHop = b
	}

	// BGP Peer password
	env = os.Getenv(bgpPeerPassword)
	if env != "" {
		c.BGPPeerConfig.Password = env
	}

	// BGP Source Interface
	env = os.Getenv(bgpSourceIF)
	if env != "" {
		c.BGPConfig.SourceIF = env
	}

	// BGP Source PrimaryAddress
	env = os.Getenv(bgpSourceIP)
	if env != "" {
		c.BGPConfig.SourceIP = env
	}

	// BGP Peer options, add them if relevant
	env = os.Getenv(bgpPeerAddress)
	if env != "" {
		c.BGPPeerConfig.Address = env
		// If we've added in a peer configuration, then we should add it to the BGP configuration
		c.BGPConfig.Peers = append(c.BGPConfig.Peers, c.BGPPeerConfig)
	}

	// BGP Timers options
	env = os.Getenv(bgpHoldTime)
	if env != "" {
		u64, err := strconv.ParseUint(env, 10, 32)
		if err != nil {
			return err
		}
		c.BGPConfig.HoldTime = u64
	}
	env = os.Getenv(bgpKeepaliveInterval)
	if env != "" {
		u64, err := strconv.ParseUint(env, 10, 32)
		if err != nil {
			return err
		}
		c.BGPConfig.KeepaliveInterval = u64
	}

	// Find Prometheus configuration
	env = os.Getenv(prometheusServer)
	if env != "" {
		c.PrometheusHTTPServer = env
	}

	return nil
}
