package internal

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/rxOverride"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/goerrors"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/zap"
)

func Inbound(
	ConnectionCancelFunc model.ConnectionCancelFunc,
	logger *zap.Logger,
	ctx context.Context,
	goFunctionCounter GoFunctionCounter.IService,
	opts ...rxgo.Option,
) func() (common.IStackBoundDefinition, error) {
	return func() (common.IStackBoundDefinition, error) {
		return common.NewBoundDefinition(
				func(
					stackData common.IStackCreateData,
					pipeData common.IPipeCreateData,
					obs rxgo.Observable,
				) (rxgo.Observable, error) {
					if stackDataActual, ok := stackData.(*data); ok {
						inBoundChannel := make(chan rxgo.Item)
						inboundStackHandler := NewInboundStackHandler(stackDataActual)

						createSendData, createSendError, createComplete, err := RxHandlers.All(
							goCommsDefinitions.Socks5,
							model.StreamDirectionUnknown,
							inBoundChannel,
							logger,
							ctx,
						)
						if err != nil {
							return nil, err
						}

						rxNextHandler, err := RxHandlers.NewRxNextHandler(
							goCommsDefinitions.Socks5,
							ConnectionCancelFunc,
							inboundStackHandler,
							createSendData,
							createSendError,
							createComplete,
							logger)
						if err != nil {
							return nil, err
						}

						_ = rxOverride.ForEach(
							goCommsDefinitions.Socks5,
							model.StreamDirectionUnknown,
							obs,
							ctx,
							goFunctionCounter,
							rxNextHandler.OnSendData,
							rxNextHandler.OnError,
							rxNextHandler.OnComplete,
							false,
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
