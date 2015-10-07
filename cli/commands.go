package cli

import "github.com/codegangsta/cli"

var (
	commands = []cli.Command{
		{
			Name:      "create",
			ShortName: "c",
			Usage:     "Create a cluster",
			Action:    create,
		},
		{
			Name:      "list",
			ShortName: "l",
			Usage:     "List nodes in a cluster",
			Flags:     []cli.Flag{flTimeout},
			Action:    list,
		},
		{
			Name:      "manage",
			ShortName: "m",
			Usage:     "Manage a docker cluster",
			Flags: []cli.Flag{
				flStrategy, flFilter,
				flHosts,
				flLeaderElection, flLeaderTTL, flManageAdvertise,
				flTLS, flTLSCaCert, flTLSCert, flTLSKey, flTLSVerify,
				flHeartBeat,
				flEnableCors,
				flCluster, flClusterOpt},
			Action: manage,
		},
		{
			Name:      "join",
			ShortName: "j",
			Usage:     "join a docker cluster",
			Flags:     []cli.Flag{flJoinAdvertise, flHeartBeat, flTTL},
			Action:    join,
		},
		{
			Name:  "cluster",
			Usage: "Start a Swarm Cluster",
			Flags: []cli.Flag{
				flClusterNodeName, flClusterEngineAddr, flClusterBindAddr, flClusterAdvertiseAddr, flClusterRaftBindAddr, flClusterRaftAdvertiseAddr, flClusterStorePath, flClusterJoin, flClusterDebug,
				flStrategy, flFilter,
				flHosts,
				flTLS, flTLSCaCert, flTLSCert, flTLSKey, flTLSVerify,
				flEnableCors,
				flCluster, flClusterOpt},
			Action: clusterRun,
		},
	}
)
