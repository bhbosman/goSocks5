package socksStack

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/gocommon"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/rxOverride"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/goerrors"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/zap"
)

func inbound(
	ConnectionCancelFunc model.ConnectionCancelFunc,
	logger *zap.Logger,
	ctx context.Context,
	goFunctionCounter GoFunctionCounter.IService,
	opts ...rxgo.Option,
) func() (common.IStackBoundFactory, error) {
	return func() (common.IStackBoundFactory, error) {
		return common.NewStackBoundDefinition(
				func(
					stackData common.IStackCreateData,
					pipeData common.IPipeCreateData,
					obs gocommon.IObservable,
				) (gocommon.IObservable, error) {
					if stackDataActual, ok := stackData.(*data); ok {
						inBoundChannel := make(chan rxgo.Item)
						inboundStackHandlerInstance := newInboundStackHandler(stackDataActual)

						handler, err := RxHandlers.All2(
							goCommsDefinitions.Socks5,
							model.StreamDirectionUnknown,
							inBoundChannel,
							logger,
							ctx,
							true,
						)
						if err != nil {
							return nil, err
						}

						rxNextHandler, err := RxHandlers.NewRxNextHandler2(
							goCommsDefinitions.Socks5,
							ConnectionCancelFunc,
							inboundStackHandlerInstance,
							handler,
							logger)
						if err != nil {
							return nil, err
						}

						_ = rxOverride.ForEach2(
							goCommsDefinitions.Socks5,
							model.StreamDirectionUnknown,
							obs,
							ctx,
							goFunctionCounter,
							rxNextHandler,
							opts...)

						resultObs := rxgo.FromChannel(inBoundChannel, opts...)
						return resultObs, nil
					}
					return nil, goerrors.InvalidType
				},
				nil),
			nil
	}
}
