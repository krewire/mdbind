package mdbind

import "github.com/krewire/krewire/packages/kern"

// Version is the mdbind module version.
var Version = kern.MustParseVersion("0.1.0")

// EcosystemRequires declares the minimum version of each Krewire module this one
// was built against. mdbind uses the Markdown renderer from libs, so it depends
// on the kernel directly and on libs.
//
// The kernel names no modules, so this module states its dependencies by name.
var EcosystemRequires = map[string]kern.Version{
	"kern": kern.MustParseVersion("0.1.0"),
	"libs": kern.MustParseVersion("0.1.0"),
}
