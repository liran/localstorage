lint:
	staticcheck ./...

test:
	go test -timeout 45s ./... -race -tags all -v
