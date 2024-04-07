package manager

import (
	patroniBgp "github.com/GCFactory/Patroni-BGP/pkg/patroni-bgp"
	"github.com/GCFactory/Patroni-BGP/pkg/vip"
	log "github.com/sirupsen/logrus"
)

func (sm *Manager) patroniWatcher() error {
	// Use a restartable watcher, as this should help in the event of etcd or timeout issues
	rw := patroniBgp.NewPatroniWatcher(sm.config.PatroniURL)
	if sm.config.PrimaryAddress != "" && vip.IsIP(sm.config.PrimaryAddress) {
		rw.EnablePrimaryAddress()
	}
	if sm.config.SyncReplicaAddress != "" && vip.IsIP(sm.config.SyncReplicaAddress) {
		rw.EnableSyncReplicaAddress()
	}
	if sm.config.AsyncReplicaAddress != "" && vip.IsIP(sm.config.AsyncReplicaAddress) {
		rw.EnableAsyncReplicaAddress()
	}
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
		case patroniBgp.PatroniStateMaster:
			err := sm.bgpServer.DelHost(sm.config.SyncReplicaAddress)
			if err != nil {
				log.Errorf("unable to remove host %s", sm.config.SyncReplicaAddress)
			}
			err = sm.bgpServer.AddHost(sm.config.PrimaryAddress)
			if err != nil {
				log.Errorf("unable to add host %s", sm.config.PrimaryAddress)
			}
		case patroniBgp.PatroniStateSyncReplica:
			err := sm.bgpServer.DelHost(sm.config.PrimaryAddress)
			if err != nil {
				log.Errorf("unable to remove host %s", sm.config.PrimaryAddress)
			}
			err = sm.bgpServer.AddHost(sm.config.SyncReplicaAddress)
			if err != nil {
				log.Errorf("unable to add host %s", sm.config.SyncReplicaAddress)
			}
		case patroniBgp.PatroniStateAsyncReplica:
			err := sm.bgpServer.DelHost(sm.config.PrimaryAddress)
			if err != nil {
				log.Errorf("unable to remove host %s", sm.config.PrimaryAddress)
			}
			err = sm.bgpServer.AddHost(sm.config.AsyncReplicaAddress)
			if err != nil {
				log.Errorf("unable to add host %s", sm.config.AsyncReplicaAddress)
			}
		case patroniBgp.PatroniStateUndefined:
			log.Warnln("undefined state")
		case patroniBgp.PatroniStateError:
		default:
			log.Warnln("error state")
		}
	}
	log.Warnln("Stopping bgp announce")
	return nil
}
