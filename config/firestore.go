package config

import (
	"encoding/json"
	"io/ioutil"

	"google.golang.org/api/option"
)

// A FirestoreConfig holds the configuration necessary for the server to talk
// to firestore
type FirestoreConfig struct {
	ProjectID         string
	ServiceAccountKey option.ClientOption
}

// LoadFirestoreConfig loads a firestore config from a JSON file in the filesystem
func LoadFirestoreConfig(paths Paths) (*FirestoreConfig, error) {
	bytes, err := ioutil.ReadFile(paths.FirestorePath)
	if err != nil {
		return nil, err
	}
	cfg := &FirestoreConfig{}
	err = json.Unmarshal(bytes, cfg)
	if err != nil {
		return nil, err
	}
	cfg.ServiceAccountKey = option.WithCredentialsFile(paths.ServiceAccountKey)
	return cfg, nil
}
