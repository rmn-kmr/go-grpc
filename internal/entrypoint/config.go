package entrypoint

import (
	"github.com/rmnkmr/go-common/config"
	lspConfig "github.com/rmnkmr/lsp/config"
)

type Configuration struct {
	config.Configuration `mapstructure:",squash"`
	Service              lspConfig.ServiceConfig `mapstructure:",squash" validate:"required"`
}
