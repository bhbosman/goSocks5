package socksStack

import (
	"context"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/Services/interfaces"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/intf"
	"github.com/cskr/pubsub"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func ProvideConnectionReactor() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Target: func(
				params struct {
					fx.In
					PubSub                 *pubsub.PubSub `name:"Application"`
					CancelCtx              context.Context
					CancelFunc             context.CancelFunc
					ConnectionCancelFunc   model.ConnectionCancelFunc
					Logger                 *zap.Logger
					GoFunctionCounter      GoFunctionCounter.IService
					UniqueReferenceService interfaces.IUniqueReferenceService
				},
			) (intf.IConnectionReactor, error) {
				return newReactor(
					params.Logger,
					params.CancelCtx,
					params.CancelFunc,
					params.ConnectionCancelFunc,
					params.GoFunctionCounter,
					params.UniqueReferenceService,
					params.PubSub,
				)
			},
		},
	)
}
