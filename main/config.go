package main

import (
	"encoding/json"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	bhbibimp "github.com/cppforlife/bosh-hub/bosh-init-bin/importer"
	bhimporter "github.com/cppforlife/bosh-hub/release/importer"
	bhnoteimporter "github.com/cppforlife/bosh-hub/release/noteimporter"
	bhwatcher "github.com/cppforlife/bosh-hub/release/watcher"
	bhstemsimp "github.com/cppforlife/bosh-hub/stemcell/importer"
)

type Config struct {
	Repos ReposOptions

	APIKey string

	Analytics AnalyticsConfig

	// Does not start web server; just does background work
	ActAsWorker bool

	Watcher      bhwatcher.FactoryOptions
	Importer     bhimporter.FactoryOptions
	NoteImporter bhnoteimporter.FactoryOptions

	StemcellImporter    bhstemsimp.FactoryOptions
	BoshInitBinImporter bhbibimp.FactoryOptions
}

type AnalyticsConfig struct {
	GoogleAnalyticsID string
}

func NewConfigFromPath(path string, fs boshsys.FileSystem) (Config, error) {
	var config Config

	bytes, err := fs.ReadFile(path)
	if err != nil {
		return config, bosherr.WrapError(err, "Reading config %s", path)
	}

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return config, bosherr.WrapError(err, "Unmarshalling config")
	}

	return config, nil
}
