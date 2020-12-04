
install:
	cd gen/templates && wtool embed . -f
	CGO_ENABLED=0 go install ./cmd/wproto

