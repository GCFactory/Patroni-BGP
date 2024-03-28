package manager

import (
	"github.com/GCFactory/Patroni-BGP/pkg/bgp"
	patroni_bgp "github.com/GCFactory/Patroni-BGP/pkg/patroni-bgp"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Manager struct {
	configMap string
	config    *patroni_bgp.Config

	//BGP Manager, this is a singleton that manages all BGP advertisements
	bgpServer *bgp.Server

	// This channel is used to catch an OS signal and trigger a shutdown
	signalChan chan os.Signal

	// This channel is used to signal a shutdown
	shutdownChan chan struct{}

	// This is a prometheus gauge indicating the state of the sessions.
	// 1 means "ESTABLISHED", 0 means "NOT ESTABLISHED"
	bgpSessionInfoGauge *prometheus.GaugeVec

	// This mutex is to protect calls from various goroutines
	mutex sync.Mutex
}

// New will create a new managing object
func New(configMap string, config *patroni_bgp.Config) (*Manager, error) {

	// Flip this to something else
	// if config.DetectControlPlane {
	// 	log.Info("[k8s client] flipping to internal service account")
	// 	_, err = clientset.CoreV1().ServiceAccounts("kube-system").Apply(context.TODO(), kubevip.GenerateSA(), v1.ApplyOptions{FieldManager: "application/apply-patch"})
	// 	if err != nil {
	// 		return nil, fmt.Errorf("could not create k8s clientset from incluster config: %v", err)
	// 	}
	// 	_, err = clientset.RbacV1().ClusterRoles().Apply(context.TODO(), kubevip.GenerateCR(), v1.ApplyOptions{FieldManager: "application/apply-patch"})
	// 	if err != nil {
	// 		return nil, fmt.Errorf("could not create k8s clientset from incluster config: %v", err)
	// 	}
	// 	_, err = clientset.RbacV1().ClusterRoleBindings().Apply(context.TODO(), kubevip.GenerateCRB(), v1.ApplyOptions{FieldManager: "application/apply-patch"})
	// 	if err != nil {
	// 		return nil, fmt.Errorf("could not create k8s clientset from incluster config: %v", err)
	// 	}
	// }

	return &Manager{
		configMap: configMap,
		config:    config,
		bgpSessionInfoGauge: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "kube_vip",
			Subsystem: "manager",
			Name:      "bgp_session_info",
			Help:      "Display state of session by setting metric for label value with current state to 1",
		}, []string{"state", "peer"}),
	}, nil
}

// Start will begin the Manager, which will start services and watch the configmap
func (sm *Manager) Start() error {

	// listen for interrupts or the Linux SIGTERM signal and cancel
	// our context, which the leader election code will observe and
	// step down
	sm.signalChan = make(chan os.Signal, 1)
	// Add Notification for Userland interrupt
	signal.Notify(sm.signalChan, syscall.SIGINT)

	// Add Notification for SIGTERM (sent from Kubernetes)
	signal.Notify(sm.signalChan, syscall.SIGTERM)

	// All watchers and other goroutines should have an additional goroutine that blocks on this, to shut things down
	sm.shutdownChan = make(chan struct{})

	log.Infoln("Starting Patroni-bgp Manager with the BGP engine")
	return sm.startBGP()
}
