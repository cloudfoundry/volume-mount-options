module code.cloudfoundry.org/volume-mount-options

go 1.23.0

toolchain go1.23.6

require (
	github.com/google/gofuzz v1.2.0
	github.com/maxbrunsfeld/counterfeiter/v6 v6.8.1
	github.com/onsi/ginkgo/v2 v2.23.0
	github.com/onsi/gomega v1.36.2
)

require (
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/pprof v0.0.0-20250302191652-9094ed2288e7 // indirect
	github.com/kr/pretty v0.2.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	golang.org/x/mod v0.24.0 // indirect
	golang.org/x/net v0.37.0 // indirect
	golang.org/x/sync v0.12.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	golang.org/x/tools v0.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

retract (
	v1.1.1 // Contains retractions
	v1.1.0 // Published before v0.3.0, presumably by accident
	v1.0.0 // Published before v0.1.0, presumably by accident
)
