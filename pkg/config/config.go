package config

import (
	"os"

	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/google/uuid"
)

//default values for application
const (
	//api const
	DefaultAPIPort = "8000"
	DefaultAPIIP   = "0.0.0.0"
	//basic auth enabled or not
	DefaultBasicAuthentication = false

	//
	DefaultDockerURL = "unix:///var/run/docker.sock"
	DefaultImagePath = "/tmp/image-content"

	//libkv
	DefaultStorageBackend = StorageBoltDB
	StorageBoltDB         = "boltdb"
	StorageConsul         = "consul"
	StorageETCD           = "etcd"
	//storage defaults
	DefaultKVStorageIp   = "0.0.0.0"
	DefaultKVStoragePort = "8500"

	//ENV
	EnvAPIPort      = "API_PORT"
	EnvAPIIP        = "API_IP"
	EnvBasicAuth    = "API_BASIC_AUTH"
	EnvDockerURL    = "DOCKER_URL"
	EnvImagePath    = "IMAGE_PATH"
	EnvDatabasePath = "BOLTDB_LOCATION"
	//consul defaults
	EnvKVStorageIp      = "KEYVAL_STORAGE_IP"
	EnvKVStoragePort    = "KEYVAL_STORAGE_PORT"
	EnvDefaultKVBackend = "STORAGE_BACKEND"
)

//Options structures for application default and configuration
type Options struct {
	Default     string
	Environment string
}

// GenerateID for Note
func GenerateID() (id string) {
	return uuid.New().String()
}

// Get - gets specified variable from either environment or default one
func Get(variable string) string {

	var config = map[string]Options{
		"EnvAPIPort": {
			Default:     DefaultAPIPort,
			Environment: EnvAPIPort,
		},
		"EnvAPIIP": {
			Default:     DefaultAPIIP,
			Environment: EnvAPIIP,
		},
		"EnvDockerURL": {
			Default:     DefaultDockerURL,
			Environment: EnvDockerURL,
		},
		"EnvImagePath": {
			Default:     DefaultImagePath,
			Environment: EnvImagePath,
		},
		"EnvBasicAuth": {
			Default:     strconv.FormatBool(DefaultBasicAuthentication),
			Environment: EnvBasicAuth,
		},

		"EnvKVStorageIp": {
			Default:     DefaultKVStorageIp,
			Environment: EnvKVStorageIp,
		},
		"EnvKVStoragePort": {
			Default:     DefaultKVStoragePort,
			Environment: EnvKVStoragePort,
		},
		"EnvDefaultKVBackend": {
			Default:     DefaultStorageBackend,
			Environment: EnvDefaultKVBackend,
		},
		"EnvDatabasePath": {
			Default:     "data/database.db",
			Environment: EnvDatabasePath,
		},
	}

	for k, v := range config {
		if k == variable {
			if os.Getenv(v.Environment) != "" {
				log.WithFields(log.Fields{
					"key":   k,
					"value": v.Environment,
				}).Debug("config: setting configuration")
				return os.Getenv(v.Environment)
			}
			log.WithFields(log.Fields{
				"key":   k,
				"value": v.Default,
			}).Debug("config: setting configuration")
			return v.Default

		}
	}
	return ""
}
