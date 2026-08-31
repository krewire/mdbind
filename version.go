package mdbind

import "github.com/krewire/libs/core"

// Version is the mdbind module version.
var Version = core.MustParseVersion("0.2.0")

// EcosystemRequires declares the minimum versions of the modules mdbind depends on.
// mdbind depends only on libs (libs/markdown), not on framework/web.
var EcosystemRequires = map[core.ModuleName]core.Version{
	core.ModuleLibs: core.MustParseVersion("0.3.0"),
}
