package manager

import (
	"github.com/GCFactory/Patroni-BGP/pkg/patroni"
	log "github.com/sirupsen/logrus"
)

func (sm *Manager) patroniWatcher() error {
	// Use a restartable watcher, as this should help in the event of etcd or timeout issues
	rw := patroni.NewPatroniWatcher(sm.config.PatroniUrl)
	rw.Start()

	go func() {
		select {
		case <-sm.shutdownChan:
			log.Debug("(svcs) shutdown called")
			// Stop the retry watcher
			rw.Stop()
			return
		}
	}()
	ch := rw.ResultChan()

	for state := range *ch {
		switch state {
		case patroni.PatroniStateMaster:
			err := sm.bgpServer.DelHost(sm.config.ReplicaAddress)
			if err != nil {
				log.Errorf("unable to remove host %s", sm.config.ReplicaAddress)
			}
			err = sm.bgpServer.AddHost(sm.config.MasterAddress)
			if err != nil {
				log.Errorf("unable to add host %s", sm.config.MasterAddress)
			}
		case patroni.PatroniStateReplica:
			err := sm.bgpServer.DelHost(sm.config.ReplicaAddress)
			if err != nil {
				log.Errorf("unable to remove host %s", sm.config.ReplicaAddress)
			}
			err = sm.bgpServer.AddHost(sm.config.MasterAddress)
			if err != nil {
				log.Errorf("unable to add host %s", sm.config.MasterAddress)
			}
		case patroni.PatroniStateUndefined:
			log.Warnln("undefined state")
		case patroni.PatroniStateError:
		default:
			log.Warnln("error state")
		}
	}
	log.Warnln("Stopping bgp announce")
	return nil
}
