package socksStack

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/common"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func ProvideSocks5Stack() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Group: "StackDefinition",
				Target: func(
					params struct {
						fx.In
						ConnectionCancelFunc model.ConnectionCancelFunc
						Opts                 []rxgo.Option
						Logger               *zap.Logger
						Ctx                  context.Context
						GoFunctionCounter    GoFunctionCounter.IService
					},
				) (common.IStackDefinition, error) {
					return common.NewStackDefinition(
						goCommsDefinitions.Socks5,
						inbound(
							params.ConnectionCancelFunc,
							params.Logger,
							params.Ctx,
							params.GoFunctionCounter,
							params.Opts...,
						),
						outbound(
							params.ConnectionCancelFunc,
							params.Logger,
							params.Ctx,
							params.GoFunctionCounter,
							params.Opts...,
						),
						createStackState(
							params.Ctx,
							params.GoFunctionCounter,
						),
					)
				},
			},
		),
	)
}
