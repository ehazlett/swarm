package cli

import (
	"encoding/json"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/swarm/config"
	"github.com/samalba/dockerclient"
)

func spec(c *cli.Context) {
	cfg := &config.Service{
		Name:      "example",
		Instances: 1,
		Pods: []config.Pod{
			{
				Name: "shell",
				Containers: []config.Container{
					{
						Config: &dockerclient.ContainerConfig{
							Image:        "busybox",
							Cmd:          []string{"sh"},
							Tty:          true,
							OpenStdin:    true,
							AttachStdin:  true,
							AttachStdout: true,
							AttachStderr: true,
							Hostname:     "shell",
						},
					},
				},
			},
		},
	}

	if err := json.NewEncoder(os.Stdout).Encode(&cfg); err != nil {
		log.Fatal(err)
	}
}
