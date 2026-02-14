package api

//go:generate go run github.com/ogen-go/ogen/cmd/ogen --config ogen-server.yml --clean --target . --package api api.yml
//go:generate go run github.com/ogen-go/ogen/cmd/ogen --config ogen-client.yml --clean --target ../../client --package client api.yml
