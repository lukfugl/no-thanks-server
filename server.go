package main

import (
	"context"
	"log"
	"net/http"

	"lukfugl/no-thanks-server/api"
	"lukfugl/no-thanks-server/config"
	"lukfugl/no-thanks-server/firestore"
)

var configPaths = config.Paths{
	FirestorePath:     "config/firestore.json",
	ServiceAccountKey: "certs/googleServiceAccountKey.json",
	SSLCertPath:       "certs/https-server.crt",
	SSLKeyPath:        "certs/https-server.key",
}

func main() {
	ctx := context.Background()

	fsCfg, err := config.LoadFirestoreConfig(configPaths)
	if err != nil {
		log.Fatal(err)
		return
	}

	fs, err := firestore.MakeClient(ctx, fsCfg)
	if err != nil {
		log.Fatal(err)
		return
	}

	http.Handle("/", &api.HTTPHandler{
		Firestore: fs,
	})
	log.Println("** Service starting on port 8080 **")
	err = http.ListenAndServeTLS(
		":8080",
		configPaths.SSLCertPath,
		configPaths.SSLKeyPath,
		nil,
	)
	if err != nil {
		log.Fatal(err)
		return
	}
}
