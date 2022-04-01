package configuration

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	SocketPath        string
	ClientConfPath    string
	ClientKeyringPath string
}

func ParseConfigFile(configPath string) (*Configuration, error) {
	// prepare a default configuration
	configuration := Configuration{
		SocketPath:        "/tmp/ceph-proxy.socket",
		ClientConfPath:    "/clients/conf",
		ClientKeyringPath: "/clients/keyring",
	}

	file, err := os.Open(configPath)
	if err != nil {
		return &configuration, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	if err != nil {
		return &configuration, err
	}

	return &configuration, nil
}
