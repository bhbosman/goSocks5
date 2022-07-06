package providers

import (
	"fmt"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goCommsNetListener"
	"github.com/bhbosman/goCommsStacks/bottom"
	"github.com/bhbosman/goCommsStacks/topStack"
	pingPong "github.com/bhbosman/goSocks5/SocksStack"
	"github.com/bhbosman/gocommon/fx/PubSub"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gocomms/intf"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"net/url"
)

func ProvideSocksConnection() fx.Option {
	sshConnectionUrl, err := url.Parse("tcp4://127.0.0.1:1080")
	if err != nil {
		return fx.Error(err)
	}
	Name := fmt.Sprintf("SocksConnection:%v", sshConnectionUrl.Port())
	return fx.Provide(
		fx.Annotated{
			Group: "Apps",
			Target: func(params struct {
				fx.In
				PubSub             *pubsub.PubSub `name:"Application"`
				NetAppFuncInParams common.NetAppFuncInParams
			}) (messages.CreateAppCallback, error) {
				f := goCommsNetListener.NewNetListenApp(
					Name,
					0,
					0,
					Name,
					false,
					nil,
					sshConnectionUrl,
					goCommsDefinitions.TransportFactorySocks5,
					common.MaxConnectionsSetting(1024),

					common.NewConnectionInstanceOptions(
						PubSub.ProvidePubSubInstance("Application", params.PubSub),
						ProvideConnectionReactor(),
						ProvideTransportFactoryForSocks5(
							bottom.Provide(),
							pingPong.ProvideSocks5Stack(),
							topStack.ProvideTopStack(),
						),
					),
				)
				return f(params.NetAppFuncInParams), nil
			},
		},
	)
}

func ProvideConnectionReactor() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Target: func(
				params struct {
					fx.In
					PubSub               *pubsub.PubSub `name:"Application"`
					CancelCtx            context.Context
					CancelFunc           context.CancelFunc
					ConnectionCancelFunc model.ConnectionCancelFunc
					Logger               *zap.Logger
					ClientContext        interface{} `name:"UserContext"`
				},
			) (intf.IConnectionReactor, error) {
				return newReactor(
					params.Logger,
					params.CancelCtx,
					params.CancelFunc,
					params.ConnectionCancelFunc,
					params.ClientContext,
				)

			},
		},
	)
}

func ProvideTransportFactoryForSocks5(
	topStackProvider fx.Option,
	socks5Provider fx.Option,
	bottomStackProvider fx.Option,
) fx.Option {
	var fxOption = []fx.Option{
		goCommsDefinitions.ProvideStackName(goCommsDefinitions.TransportFactorySocks5),
		fx.Provide(
			fx.Annotated{
				Group: "TransportFactory",
				Target: goCommsDefinitions.NewTransportFactory(
					goCommsDefinitions.TransportFactorySocks5,
					//goCommsDefinitions.TopStackName,
					goCommsDefinitions.Socks5,
					//goCommsDefinitions.BottomStackStackName,
				),
			},
		),
	}
	if topStackProvider != nil {
		fxOption = append(fxOption, topStackProvider)
	}
	if socks5Provider != nil {
		fxOption = append(fxOption, socks5Provider)
	}
	if bottomStackProvider != nil {
		fxOption = append(fxOption, bottomStackProvider)
	}
	return fx.Options(fxOption...)
}
