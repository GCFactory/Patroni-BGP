package manager

import (
	"context"
	"github.com/GCFactory/Patroni-BGP/pkg/patroni"
	log "github.com/sirupsen/logrus"
)

func (sm *Manager) patroniWatcher(ctx context.Context) error {

	//id, err := os.Hostname()
	//if err != nil {
	//	return err
	//}

	// Use a restartable watcher, as this should help in the event of etcd or timeout issues
	rw := patroni.NewPatroniWatcher(sm.config.PatroniUrl)
	rw.Start()

	exitFunction := make(chan struct{})
	go func() {
		select {
		case <-sm.shutdownChan:
			log.Debug("(svcs) shutdown called")
			// Stop the retry watcher
			rw.Stop()
			return
		case <-exitFunction:
			log.Debug("(svcs) function ending")
			// Stop the retry watcher
			rw.Stop()
			err := sm.bgpServer.DelHost(sm.config.MasterAddress)
			if err != nil {
				log.Errorf("unable to remove host %s", sm.config.MasterAddress)
			}
			err = sm.bgpServer.DelHost(sm.config.ReplicaAddress)
			if err != nil {
				log.Errorf("unable to remove host %s", sm.config.ReplicaAddress)
			}
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

	close(exitFunction)
	log.Warnln("Stopping watching services for type: LoadBalancer in all namespaces")
	return nil
}
