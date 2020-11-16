module github.com/shupkg/wproto

go 1.15

replace github.com/shupkg/wbin => ../wbin

require (
	github.com/shupkg/wbin v0.0.0-00010101000000-000000000000
	github.com/spf13/pflag v1.0.5
	google.golang.org/protobuf v1.25.0
)
