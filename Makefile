
install:
	cd gen/templates && wbin . -f
	CGO_ENABLED=0 go install ./cmd/wproto

