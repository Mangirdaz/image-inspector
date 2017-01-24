package storage

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/mangirdaz/image-inspector/pkg/config"

	boltdb "github.com/docker/libkv/store/boltdb"
	"github.com/docker/libkv/store/consul"

	log "github.com/Sirupsen/logrus"
)

// NotFound - message for objects that could not be found
const NotFound = "Not found"

// Namespace - namespace for storage
const Namespace = "notes/"

// backend related errors
const (
	ErrFailedToLock string = "error, failed to acquire lock"
	ErrLockNotFound string = "error, lock not found"
)

func addNamespace(key string) string {
	var buffer bytes.Buffer
	buffer.WriteString(Namespace)
	if key != "/" {
		buffer.WriteString(key)
	}
	return buffer.String()
}

func removeNamespace(key string) string {
	return strings.TrimPrefix(key, Namespace)
}

// KVPair - kv pair to store results from kv backends
type KVPair struct {
	Key       string
	Value     []byte
	LastIndex uint64
}

// KVBackend - generic kv backend interface
type KVBackend interface {
	Put(key string, value []byte) error
	Get(key string) (*KVPair, error)
}

// LibKVBackend - libkv container
type LibKVBackend struct {
	Options *store.Config
	Backend string
	Addrs   []string
	Store   store.Store
	Locks   map[string]store.Locker

	MutexLocks map[string]*sync.Mutex

	TTLLocks map[string]chan struct{}
}

func init() {
	consul.Register()
	boltdb.Register()
	// etcd.Register()
	log.Debug("stores registered")
}

// available backend stores
const (
	BoltDBBackend    = "boltdb"
	ZookeeperBackend = "zk"
	EtcdBackend      = "etcd"
	ConsulBackend    = "consul"
)

// NewLibKVBackend - returns new libkv based backend (https://github.com/docker/libkv)
func NewLibKVBackend(backend, bucket string, addrs []string) (*LibKVBackend, error) {
	// default backend - consul
	var backendStore store.Backend

	// TODO: options could include TLS details (certs) or bucket if backend is BoltDB
	options := &store.Config{
		ConnectionTimeout: 10 * time.Second,
	}

	if backend == EtcdBackend {
		backendStore = store.ETCD
		return nil, fmt.Errorf("this backend is currently not supported")
	} else if backend == ZookeeperBackend {
		backendStore = store.ZK
		return nil, fmt.Errorf("this backend is currently not supported")
	} else if backend == BoltDBBackend {
		backendStore = store.BOLTDB
		// setting bucket name to address
		options.Bucket = bucket
	} else {
		backendStore = store.CONSUL
	}

	log.WithFields(log.Fields{
		"backend": backendStore,
		"addrs":   addrs,
	}).Info("initiating store")

	libstore, err := libkv.NewStore(backendStore, addrs, options)
	if err != nil {
		log.WithFields(log.Fields{
			"backend": backendStore,
			"addrs":   addrs,
		}).Error("failed to create store")
		return nil, err
	}

	locks := make(map[string]store.Locker)

	ttlLocks := make(map[string]chan struct{})

	// preparing mutex locks, only used by BoltDB
	ml := make(map[string]*sync.Mutex)

	return &LibKVBackend{
		Backend:    backend,
		Addrs:      addrs,
		Options:    options,
		Store:      libstore,
		Locks:      locks,
		TTLLocks:   ttlLocks,
		MutexLocks: ml,
	}, nil
}

// Put - puts object into kv store
func (l *LibKVBackend) Put(key string, value []byte) error {
	log.Debugf("Put to DB [%s]", fmt.Sprintf("%s%s", Namespace, key))
	return l.Store.Put(fmt.Sprintf("%s%s", Namespace, key), value, nil)
}

// Get - gets object from kv store
func (l *LibKVBackend) Get(key string) (*KVPair, error) {
	kvPair, err := l.Store.Get(addNamespace(key))
	if err != nil {
		// replacing error with locally defined
		if err == store.ErrKeyNotFound {
			return nil, fmt.Errorf(NotFound)
		}

		return nil, err
	}
	return &KVPair{
		Key:       removeNamespace(kvPair.Key),
		Value:     kvPair.Value,
		LastIndex: kvPair.LastIndex,
	}, nil
}

func InitKVStorage() *LibKVBackend {

	switch storagebackend := config.Get("EnvDefaultKVBackend"); storagebackend {
	case config.StorageBoltDB:
		return generateDefaultStorageBackend()
	case config.StorageConsul:
		consul.Register()
		db := config.Get("EnvDefaultConsulAddr")
		kv, err := NewLibKVBackend(config.StorageConsul, "default", []string{db})
		if err != nil {
			log.WithFields(log.Fields{
				"path":  db,
				"error": err,
			}).Fatal("failed to create new consult connection")
		}
		return kv
	default:
		return generateDefaultStorageBackend()
	}
}

func generateDefaultStorageBackend() *LibKVBackend {
	db := config.Get("EnvDatabasePath")
	kv, err := NewLibKVBackend(config.DefaultStorageBackend, "default", []string{db})
	if err != nil {
		log.WithFields(log.Fields{
			"path":  db,
			"error": err,
		}).Fatal("failed to create new database")
	}
	return kv
}
