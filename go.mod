module github.com/krewire/mdbind

go 1.27.1

require (
	github.com/krewire/krewire v0.1.0
	gopkg.in/yaml.v3 v3.0.1
)

require golang.org/x/net v0.59.0 // indirect

replace github.com/krewire/krewire => ../krewire
