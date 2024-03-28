package main

import (
	"context"
	"errors"
	"github.com/GCFactory/Patroni-BGP/pkg/manager"
	patroni_bgp "github.com/GCFactory/Patroni-BGP/pkg/patroni-bgp"
	"github.com/GCFactory/Patroni-BGP/pkg/vip"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

var (
	initConfig patroni_bgp.Config
)

// ConfigMap name within a Kubernetes cluster
var configMap string

func main() {
	ctx := context.TODO()
	// parse environment variables, these will overwrite anything loaded or flags
	err := patroni_bgp.ParseEnvironment(&initConfig)
	if err != nil {
		log.Fatalln(err)
	}

	// Set the logging level for all subsequent functions
	log.SetLevel(log.Level(initConfig.Logging))

	// start prometheus server
	if initConfig.PrometheusHTTPServer != "" {
		go servePrometheusHTTPServer(ctx, PrometheusHTTPServerConfig{
			Addr: initConfig.PrometheusHTTPServer,
		})
	}

	// Check if the interface needs auto-detecting
	if initConfig.Interface == "" {
		log.Infof("No interface is specified for VIP in config, auto-detecting default Interface")
		defaultIF, err := vip.GetDefaultGatewayInterface()
		if err != nil {
			log.Fatalf("unable to detect default interface -> [%v]", err)
		}
		initConfig.Interface = defaultIF.Name
		log.Infof("patroni-bgp will bind to interface [%s]", initConfig.Interface)

		go func() {
			if err := vip.MonitorDefaultInterface(context.TODO(), defaultIF); err != nil {
				log.Fatalf("crash: %s", err.Error())
			}
		}()
	}

	// Perform a check on th state of the interface
	if err := initConfig.CheckInterface(); err != nil {
		log.Fatalln(err)
	}

	// User Environment variables as an option to make manifest clearer
	envConfigMap := os.Getenv("vip_configmap")
	if envConfigMap != "" {
		configMap = envConfigMap
	}

	// Define the new service manager
	mgr, err := manager.New(configMap, &initConfig)
	if err != nil {
		log.Fatalf("configuring new Manager error -> %v", err)
	}

	prometheus.MustRegister(mgr.PrometheusCollector()...)

	// Start the service manager, this will watch the config Map and construct patroni-bgp services for it
	err = mgr.Start()
	if err != nil {
		log.Fatalf("starting new Manager error -> %v", err)
	}
}

// PrometheusHTTPServerConfig defines the Prometheus server configuration.
type PrometheusHTTPServerConfig struct {
	// Addr sets the http server address used to expose the metric endpoint
	Addr string
}

func servePrometheusHTTPServer(ctx context.Context, config PrometheusHTTPServerConfig) {
	var err error
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`<html>
			<head><title>patroni-bgp</title></head>
			<body>
			<h1>patroni-bgp Metrics</h1>
			<p><a href="` + "/metrics" + `">Metrics</a></p>
			</body>
			</html>`))
	})

	srv := &http.Server{
		Addr:              config.Addr,
		Handler:           mux,
		ReadHeaderTimeout: 2 * time.Second,
	}

	go func() {
		if err = srv.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			log.Fatalf("listen:%+s\n", err)
		}
	}()

	log.Printf("prometheus HTTP server started")

	<-ctx.Done()

	log.Printf("prometheus HTTP server stopped")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err = srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("server Shutdown Failed:%+s", err)
	}

	if err == http.ErrServerClosed {
		err = nil
	}
}
