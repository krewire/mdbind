package mdbind

import (
	"github.com/krewire/krewire/packages/model"
	"github.com/krewire/krewire/packages/version"
)

// Version is the mdbind module version.
var Version = version.MustParseVersion("0.1.0")

// Tier is the ecosystem tier assigned to the mdbind repository.
var Tier = model.TierFree

// EcosystemRequires declares the minimum version of each Krewire module this one
// was built against.
var EcosystemRequires = map[string]version.Version{
	"version": version.MustParseVersion("0.1.0"),
	"model":   version.MustParseVersion("0.1.0"),
}
