package bootstrap

import (
	"everlasting/src/infrastructure/pkg"

	"github.com/sarulabs/di"
)

// InitContainer - Initialize all container
func InitializeContainer(config *pkg.Config) di.Container {
	builder, _ := di.NewBuilder()
	loadAdapter(builder, config)
	loadPkg(builder, config)
	loadPersistence(builder, config)

	return builder.Build()
}
