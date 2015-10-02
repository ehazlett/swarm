package cluster

import (
	"github.com/hashicorp/raft"
)

type SwarmFSMSnapshot struct{}

func (s SwarmFSMSnapshot) Persist(sink raft.SnapshotSink) error {
	return nil
}

func (s SwarmFSMSnapshot) Release() {}
