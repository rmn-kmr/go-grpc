package main

import (
	networkHack "github.com/rmnkmr/go-common/hack"
	"github.com/rmnkmr/lsp/internal/entrypoint"
)

func main() {
	networkHack.WaitForNetwork()
	entrypoint.Execute()
}
