package manager

import (
	"context"
	"fmt"
	"github.com/GCFactory/Patroni-BGP/pkg/bgp"
	api "github.com/osrg/gobgp/v3/api"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// Start will begin the Manager, which will start services and watch the configmap
func (sm *Manager) startBGP() error {
	var err error

	log.Info("Starting the BGP server to advertise VIP routes to BGP peers")
	sm.bgpServer, err = bgp.NewBGPServer(&sm.config.BGPConfig, func(p *api.WatchEventResponse_PeerEvent) {
		ipaddr := p.GetPeer().GetState().GetNeighborAddress()
		port := uint64(179)
		peerDescription := fmt.Sprintf("%s:%d", ipaddr, port)

		for stateName, stateValue := range api.PeerState_SessionState_value {
			metricValue := 0.0
			if stateValue == int32(p.GetPeer().GetState().GetSessionState().Number()) {
				metricValue = 1
			}

			sm.bgpSessionInfoGauge.With(prometheus.Labels{
				"state": stateName,
				"peer":  peerDescription,
			}).Set(metricValue)
		}
	})
	if err != nil {
		return err
	}

	// use a Go context so we can tell the leaderelection code when we
	// want to step down
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Defer a function to check if the bgpServer has been created and if so attempt to close it
	defer func() {
		if sm.bgpServer != nil {
			err := sm.bgpServer.Close()
			if err != nil {
				return
			}
		}
	}()

	// Shutdown function that will wait on this signal, unless we call it ourselves
	go func() {
		<-sm.signalChan
		sm.shutdownChan <- struct{}{}
		log.Info("Received termination, signaling shutdown")
		// Cancel the context, which will in turn cancel the leadership
		cancel()
	}()

	err = sm.patroniWatcher()
	if err != nil {
		return err
	}

	log.Infof("Shutting down Patroni-bgp")

	return nil
}
