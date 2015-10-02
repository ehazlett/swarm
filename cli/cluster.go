package cli

import (
	"crypto/tls"
	"net"
	"os"
	"os/signal"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/swarm/api"
	//"github.com/docker/swarm/cluster/mesos"
	"github.com/docker/swarm/cluster/swarm"
	"github.com/docker/swarm/scheduler"
	"github.com/docker/swarm/scheduler/filter"
	"github.com/docker/swarm/scheduler/strategy"
	"github.com/ehazlett/libdiscover"
)

func waitForInterrupt(f func()) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	for _ = range sigChan {
		// cleanup
		f()
	}
}

func clusterRun(c *cli.Context) {
	var (
		tlsConfig *tls.Config
		err       error
	)

	// If either --tls or --tlsverify are specified, load the certificates.
	if c.Bool("tls") || c.Bool("tlsverify") {
		if !c.IsSet("tlscert") || !c.IsSet("tlskey") {
			log.Fatal("--tlscert and --tlskey must be provided when using --tls")
		}
		if c.Bool("tlsverify") && !c.IsSet("tlscacert") {
			log.Fatal("--tlscacert must be provided when using --tlsverify")
		}
		tlsConfig, err = loadTLSConfig(
			c.String("tlscacert"),
			c.String("tlscert"),
			c.String("tlskey"),
			c.Bool("tlsverify"))
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// Otherwise, if neither --tls nor --tlsverify are specified, abort if
		// the other flags are passed as they will be ignored.
		if c.IsSet("tlscert") || c.IsSet("tlskey") || c.IsSet("tlscacert") {
			log.Fatal("--tlscert, --tlskey and --tlscacert require the use of either --tls or --tlsverify")
		}
	}

	s, err := strategy.New(c.String("strategy"))
	if err != nil {
		log.Fatal(err)
	}

	// see https://github.com/codegangsta/cli/issues/160
	names := c.StringSlice("filter")
	if c.IsSet("filter") || c.IsSet("f") {
		names = names[DefaultFilterNumber:]
	}
	fs, err := filter.New(names)
	if err != nil {
		log.Fatal(err)
	}

	raftAdvAddr, err := net.ResolveTCPAddr("tcp", c.String("cluster-raft-advertise-addr"))
	if err != nil {
		log.Fatal(err)
	}

	// start discover
	discoverConfig := &libdiscover.DiscoverConfig{
		Name:              c.String("cluster-node-name"),
		BindAddr:          c.String("cluster-bind-addr"),
		BindPort:          c.Int("cluster-bind-port"),
		AdvertiseAddr:     c.String("cluster-advertise-addr"),
		AdvertisePort:     c.Int("cluster-advertise-port"),
		RaftBindAddr:      c.String("cluster-raft-bind-addr"),
		RaftAdvertiseAddr: raftAdvAddr,
		JoinAddr:          c.String("cluster-join"),
		StorePath:         c.String("cluster-store-path"),
		Debug:             c.Bool("cluster-debug"),
	}

	sched := scheduler.New(s, fs)
	cl, err := swarm.NewCluster(sched, tlsConfig, discoverConfig, c.StringSlice("cluster-opt"), c.String("cluster-engine-addr"))
	if err != nil {
		log.Fatal(err)
	}

	// see https://github.com/codegangsta/cli/issues/160
	hosts := c.StringSlice("host")
	if c.IsSet("host") || c.IsSet("H") {
		hosts = hosts[1:]
	}

	server := api.NewServer(hosts, tlsConfig)
	server.SetHandler(api.NewPrimary(cl, tlsConfig, &statusHandler{cl, nil, nil}, c.Bool("cors")))

	go server.ListenAndServe()

	// catch an interrupt so we can stop the cluster gracefully
	waitForInterrupt(func() {
		log.Info("shutting down")

		// exit
		if err := cl.Stop(); err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	})
}
