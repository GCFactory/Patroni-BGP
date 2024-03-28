package manager

import (
	"context"
	"fmt"
	"github.com/GCFactory/Patroni-BGP/pkg/vip"
	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"time"
)

// Start will begin the Manager, which will start services and watch the configmap
func (sm *Manager) startTableMode() error {
	var err error

	// use a Go context so we can tell the leaderelection code when we
	// want to step down
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	log.Infof("all routing table entries will exist in table [%d] with protocol [%d]", sm.config.RoutingTableID, sm.config.RoutingProtocol)

	if sm.config.CleanRoutingTable {
		go func() {
			// we assume that after 10s all services should be configured so we can delete redundant routes
			time.Sleep(time.Second * 10)
			if err := sm.cleanRoutes(); err != nil {
				log.Errorf("error checking for old routes: %v", err)
			}
		}()
	}

	// Shutdown function that will wait on this signal, unless we call it ourselves
	go func() {
		<-sm.signalChan
		log.Info("Received termination, signaling shutdown")

		// Cancel the context, which will in turn cancel the leadership
		cancel()
	}()

	log.Infof("beginning watching services without leader election")
	err = sm.patroniWatcher(ctx)
	if err != nil {
		log.Errorf("Cannot watch services, %v", err)
	}
	return nil
}

func (sm *Manager) cleanRoutes() error {
	routes, err := vip.ListRoutes(sm.config.RoutingTableID, sm.config.RoutingProtocol)
	if err != nil {
		return fmt.Errorf("error getting routes: %w", err)
	}

	for i := range routes {
		found := false
		//for _, instance := range sm.serviceInstances {
		//	for _, cluster := range instance.clusters {
		//		for n := range cluster.Network {
		//			r := cluster.Network[n].PrepareRoute()
		//			if r.Dst.String() == routes[i].Dst.String() {
		//				found = true
		//			}
		//		}
		//	}
		//}
		if !found {
			err = netlink.RouteDel(&(routes[i]))
			if err != nil {
				log.Errorf("[route] error deleting route: %v", routes[i])
			}
			log.Debugf("[route] deleted route: %v", routes[i])
		}
	}
	return nil
}
