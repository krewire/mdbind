package mdbind

import "github.com/krewire/libs/core"

// Version is the mdbind module version.
var Version = core.MustParseVersion("0.1.0")

// EcosystemRequires declares compatibility.
var EcosystemRequires = map[core.ModuleName]core.Version{
	core.ModuleFramework: core.MustParseVersion("0.1.0"),
	core.ModuleLibs:      core.MustParseVersion("0.1.0"),
}
