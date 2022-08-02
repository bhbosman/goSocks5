package providers

import (
	"fmt"
	"github.com/bhbosman/goCommsNetListener"
	"github.com/bhbosman/goSocks5/socksStack"
	"github.com/bhbosman/gocommon/fx/PubSub"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocomms/common"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
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
					Name,
					false,
					nil,
					sshConnectionUrl,
					common.MaxConnectionsSetting(1024),
					common.NewConnectionInstanceOptions(
						PubSub.ProvidePubSubInstance("Application", params.PubSub),
						socksStack.ProvideConnectionReactor(),
						socksStack.ProvideTransportFactoryForSocks5(),
					),
				)
				return f(params.NetAppFuncInParams), nil
			},
		},
	)
}
