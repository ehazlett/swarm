package cluster

import (
	"bytes"
	"encoding/gob"
	"errors"
	"io"

	log "github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"github.com/hashicorp/raft"
)

type CommandType uint8

const (
	CommandNodeJoin CommandType = iota
	CommandNodeLeave
)

type Command struct {
	CmdType CommandType
	Key     string
	Value   string
}

var (
	BucketNodes        = []byte("swarm")
	ErrKeyDoesNotExist = errors.New("key does not exist")
)

type SwarmFSM struct {
	Store *bolt.DB
}

func (f SwarmFSM) Apply(l *raft.Log) interface{} {
	log.Debug("fsm: apply")
	buf := bytes.NewBuffer(l.Data)

	dec := gob.NewDecoder(buf)

	var cmd Command
	if err := dec.Decode(&cmd); err != nil {
		return err
	}

	switch cmd.CmdType {
	case CommandNodeJoin:
		f.Store.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte(BucketNodes))
			if err != nil {
				return err
			}

			if err := b.Put([]byte(cmd.Key), []byte(cmd.Value)); err != nil {
				return err
			}

			log.Debug("store: added key: name=%s", cmd.Key)

			return nil
		})
	case CommandNodeLeave:
		f.Store.Update(func(tx *bolt.Tx) error {
			return tx.Bucket([]byte(BucketNodes)).Delete([]byte(cmd.Key))
		})

		log.Debug("store: removed key: name=%s", cmd.Key)
	}

	return nil
}

func (f SwarmFSM) Snapshot() (raft.FSMSnapshot, error) {
	log.Debug("fsm: snapshot")
	return SwarmFSMSnapshot{}, nil
}

func (f SwarmFSM) Restore(r io.ReadCloser) error {
	log.Debug("fsm: restore")
	return nil
}
