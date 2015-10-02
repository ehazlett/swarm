package cli

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/docker/swarm/scheduler/filter"
	"github.com/docker/swarm/scheduler/strategy"
)

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}

	addr := "127.0.0.1"
	// automatically pick the second interface
	if len(addrs) > 1 {
		a := addrs[1].String()
		ipnet := strings.Split(a, "/")
		addr = ipnet[0]
	}

	return addr
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}

	return hostname
}

func homepath(p string) string {
	home := os.Getenv("HOME")
	if runtime.GOOS == "windows" {
		home = os.Getenv("USERPROFILE")
	}
	return filepath.Join(home, p)
}

func getDiscovery(c *cli.Context) string {
	if len(c.Args()) == 1 {
		return c.Args()[0]
	}
	return os.Getenv("SWARM_DISCOVERY")
}

var (
	flJoinAdvertise = cli.StringFlag{
		Name:   "advertise, addr",
		Usage:  "Address of the Docker Engine joining the cluster. Swarm manager(s) MUST be able to reach the Docker Engine at this address.",
		EnvVar: "SWARM_ADVERTISE",
	}
	flManageAdvertise = cli.StringFlag{
		Name:   "advertise, addr",
		Usage:  "Address of the swarm manager joining the cluster. Other swarm manager(s) MUST be able to reach the swarm manager at this address.",
		EnvVar: "SWARM_ADVERTISE",
	}
	// hack for go vet
	flHostsValue = cli.StringSlice([]string{"tcp://128.0.0.1:2375"})

	flHosts = cli.StringSliceFlag{
		Name:   "host, H",
		Value:  &flHostsValue,
		Usage:  "ip/socket to listen on",
		EnvVar: "SWARM_HOST",
	}
	flClusterNodeName = cli.StringFlag{
		Name:  "cluster-node-name",
		Usage: "name of node in cluster",
		Value: getHostname(),
	}
	flClusterEngineAddr = cli.StringFlag{
		Name:  "cluster-engine-addr",
		Usage: "address to the docker engine",
		Value: "127.0.0.1:2375",
	}
	flClusterBindAddr = cli.StringFlag{
		Name:  "cluster-bind-addr",
		Usage: "bind address",
		Value: "0.0.0.0",
	}
	flClusterBindPort = cli.IntFlag{
		Name:  "cluster-bind-port",
		Usage: "bind port",
		Value: 7946,
	}
	flClusterAdvertiseAddr = cli.StringFlag{
		Name:  "cluster-advertise-addr",
		Usage: "advertise address",
		Value: getLocalIP(),
	}
	flClusterAdvertisePort = cli.IntFlag{
		Name:  "cluster-advertise-port",
		Usage: "advertise port",
		Value: 7946,
	}
	flClusterRaftBindAddr = cli.StringFlag{
		Name:  "cluster-raft-bind-addr",
		Usage: "raft bind address",
		Value: "0.0.0.0:8746",
	}
	flClusterRaftAdvertiseAddr = cli.StringFlag{
		Name:  "cluster-raft-advertise-addr",
		Usage: "raft advertise address",
		Value: fmt.Sprintf("%s:8746", getLocalIP()),
	}
	flClusterStorePath = cli.StringFlag{
		Name:  "cluster-store-path",
		Usage: "cluster storage path",
		Value: filepath.Join(os.TempDir(), "grid"),
	}
	flClusterJoin = cli.StringFlag{
		Name:  "cluster-join",
		Usage: "join an existing cluster",
		Value: "",
	}
	flClusterDebug = cli.BoolFlag{
		Name:  "cluster-debug",
		Usage: "enable debug logging for cluster",
	}
	flHeartBeat = cli.StringFlag{
		Name:  "heartbeat",
		Value: "20s",
		Usage: "period between each heartbeat",
	}
	flTTL = cli.StringFlag{
		Name:  "ttl",
		Value: "60s",
		Usage: "sets the expiration of an ephemeral node",
	}
	flTimeout = cli.StringFlag{
		Name:  "timeout",
		Value: "10s",
		Usage: "timeout period",
	}
	flEnableCors = cli.BoolFlag{
		Name:  "api-enable-cors, cors",
		Usage: "enable CORS headers in the remote API",
	}
	flTLS = cli.BoolFlag{
		Name:  "tls",
		Usage: "use TLS; implied by --tlsverify=true",
	}
	flTLSCaCert = cli.StringFlag{
		Name:  "tlscacert",
		Usage: "trust only remotes providing a certificate signed by the CA given here",
	}
	flTLSCert = cli.StringFlag{
		Name:  "tlscert",
		Usage: "path to TLS certificate file",
	}
	flTLSKey = cli.StringFlag{
		Name:  "tlskey",
		Usage: "path to TLS key file",
	}
	flTLSVerify = cli.BoolFlag{
		Name:  "tlsverify",
		Usage: "use TLS and verify the remote",
	}
	flStrategy = cli.StringFlag{
		Name:  "strategy",
		Usage: "placement strategy to use [" + strings.Join(strategy.List(), ", ") + "]",
		Value: strategy.List()[0],
	}

	// hack for go vet
	flFilterValue = cli.StringSlice(filter.List())
	// DefaultFilterNumber is exported
	DefaultFilterNumber = len(flFilterValue)

	flFilter = cli.StringSliceFlag{
		Name:  "filter, f",
		Usage: "filter to use [" + strings.Join(filter.List(), ", ") + "]",
		Value: &flFilterValue,
	}

	flCluster = cli.StringFlag{
		Name:  "cluster-driver, c",
		Usage: "cluster driver to use [swarm, mesos-experimental]",
		Value: "swarm",
	}
	flClusterOpt = cli.StringSliceFlag{
		Name:  "cluster-opt",
		Usage: "cluster driver options",
		Value: &cli.StringSlice{},
	}

	flLeaderElection = cli.BoolFlag{
		Name:  "replication",
		Usage: "Enable Swarm manager replication",
	}
	flLeaderTTL = cli.StringFlag{
		Name:  "replication-ttl",
		Value: "30s",
		Usage: "Leader lock release time on failure",
	}
)
