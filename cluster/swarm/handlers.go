package swarm

import (
	"bytes"
	"encoding/gob"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/swarm/cluster"
	"github.com/hashicorp/serf/serf"
)

type SyncAction int

const (
	SyncActionAdd SyncAction = iota
	SyncActionRemove
)

type SyncOp struct {
	Action SyncAction
	Node   string
}

func (c *Cluster) handlers() map[string]func(e serf.UserEvent) error {
	clusterEventHandlers := map[string]func(e serf.UserEvent) error{
		"node-join":    c.clusterNodeJoinHandler,
		"node-leave":   c.clusterNodeLeaveHandler,
		"engine-join":  c.clusterEngineJoinHandler,
		"engine-leave": c.clusterEngineLeaveHandler,
	}

	return clusterEventHandlers
}

func encode(data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)

	if err := enc.Encode(data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// handlers

func (c *Cluster) applyCommand(cmd *cluster.Command, timeout time.Duration) error {
	if c.IsReady() {
		data, err := encode(cmd)
		if err != nil {
			return err
		}

		if err := c.discover.Apply(data, time.Second*1); err != nil {
			return err
		}
	}

	return nil
}

func (c *Cluster) clusterNodeJoinHandler(e serf.UserEvent) error {
	log.Debug("cluster event: engine-join")
	if c.IsReady() && c.discover.IsLeader() {
		if err := c.discover.SendEvent("engine-join", []byte(c.engineAddr)); err != nil {
			return err
		}
	}

	return nil
}

func (c *Cluster) clusterNodeLeaveHandler(e serf.UserEvent) error {
	log.Debug("cluster event: engine-leave")
	if c.IsReady() && c.discover.IsLeader() {
		if err := c.discover.SendEvent("engine-leave", []byte(c.engineAddr)); err != nil {
			return err
		}
	}

	return nil
}

func (c *Cluster) clusterEngineJoinHandler(e serf.UserEvent) error {
	if c.IsReady() {
		addr := string(e.Payload)

		cmd := &cluster.Command{
			CmdType: cluster.CommandNodeJoin,
			Key:     addr,
			Value:   addr,
		}

		if err := c.applyCommand(cmd, time.Second*1); err != nil {
			return err
		}
	}

	return nil
}

func (c *Cluster) clusterEngineLeaveHandler(e serf.UserEvent) error {
	if c.IsReady() {
		addr := string(e.Payload)

		cmd := &cluster.Command{
			CmdType: cluster.CommandNodeLeave,
			Key:     addr,
			Value:   addr,
		}

		if err := c.applyCommand(cmd, time.Second*1); err != nil {
			return err
		}

	}

	return nil
}
