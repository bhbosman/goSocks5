package socksStack

import (
	"github.com/bhbosman/goCommsDefinitions"
	"go.uber.org/fx"
)

func ProvideTransportFactoryForSocks5() fx.Option {
	stackName := goCommsDefinitions.TransportFactorySocks5
	return fx.Options(
		ProvideSocks5Stack(),
		goCommsDefinitions.ProvideStackName(stackName),
		fx.Provide(
			fx.Annotated{
				Group:  "TransportFactory",
				Target: goCommsDefinitions.NewTransportFactory(stackName, goCommsDefinitions.Socks5),
			},
		),
	)
}
