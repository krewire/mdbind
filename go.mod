module github.com/krewire/mdbind

go 1.22

require (
	github.com/krewire/framework v0.5.1
	github.com/krewire/libs v0.1.0
	github.com/yuin/goldmark v1.8.5
)

require (
	github.com/gomarkdown/markdown v0.0.0-20260818103853-6d1f24fc3a11 // indirect
	golang.org/x/net v0.30.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/krewire/libs => ../libs

replace github.com/krewire/framework => ../framework
