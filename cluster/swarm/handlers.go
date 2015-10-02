package swarm

import (
	"bytes"
	"encoding/gob"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/swarm/cluster"
	"github.com/hashicorp/serf/serf"
)

func (c *Cluster) handlers() map[string]func(e serf.UserEvent) error {
	clusterEventHandlers := map[string]func(e serf.UserEvent) error{
		"node-join":    c.clusterNodeJoinHandler,
		"node-leave":   c.clusterNodeLeaveHandler,
		"engine-join":  c.clusterEngineJoinHandler,
		"engine-leave": c.clusterEngineLeaveHandler,
		"engine-sync":  c.clusterSyncEnginesHandler,
	}

	return clusterEventHandlers
}

// handlers

func (c *Cluster) applyCommand(cmd *cluster.Command, timeout time.Duration) error {
	if c.IsReady() {
		buf := bytes.NewBuffer(nil)
		enc := gob.NewEncoder(buf)

		if err := enc.Encode(cmd); err != nil {
			return err
		}

		if err := c.discover.Apply(buf.Bytes(), time.Second*1); err != nil {
			return err
		}

	}

	return nil
}

func (c *Cluster) clusterNodeJoinHandler(e serf.UserEvent) error {
	log.Debug("cluster event: engine-join")
	if err := c.discover.SendEvent("engine-join", []byte(c.engineAddr)); err != nil {
		return err
	}

	return nil
}

func (c *Cluster) clusterNodeLeaveHandler(e serf.UserEvent) error {
	log.Debug("cluster event: engine-leave")
	if err := c.discover.SendEvent("engine-leave", []byte(c.engineAddr)); err != nil {
		return err
	}

	return nil
}

func (c *Cluster) clusterEngineJoinHandler(e serf.UserEvent) error {
	addr := string(e.Payload)

	cmd := &cluster.Command{
		CmdType: cluster.CommandNodeJoin,
		Key:     addr,
		Value:   addr,
	}

	if err := c.applyCommand(cmd, time.Second*1); err != nil {
		return err
	}

	// trigger sync
	if err := c.discover.SendEvent("engine-sync", nil); err != nil {
		return err
	}

	return nil
}

func (c *Cluster) clusterEngineLeaveHandler(e serf.UserEvent) error {
	addr := string(e.Payload)

	cmd := &cluster.Command{
		CmdType: cluster.CommandNodeLeave,
		Key:     addr,
		Value:   addr,
	}

	if err := c.applyCommand(cmd, time.Second*1); err != nil {
		return err
	}

	// trigger sync
	if err := c.discover.SendEvent("engine-sync", nil); err != nil {
		return err
	}

	return nil
}

// syncEnginesHandler syncs the cluster engine states
func (c *Cluster) clusterSyncEnginesHandler(e serf.UserEvent) error {
	if c.IsReady() {
		// allow for replication
		time.Sleep(time.Second * 2)

		// trigger sync
		c.syncEngines()
	}

	return nil
}
