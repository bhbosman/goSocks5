package socksStack

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/goerrors"
	"github.com/reactivex/rxgo/v2"
)

func createStackState(
	ctx context.Context,
	goFunctionCounter GoFunctionCounter.IService,
) common.IStackState {
	return common.NewStackState(
		goCommsDefinitions.Socks5,
		false,
		func() (common.IStackCreateData, error) {
			return NewStackData(
				ctx,
				goFunctionCounter,
			), nil
		},
		func(connectionType model.ConnectionType, stackData common.IStackCreateData) error {
			if sd, ok := stackData.(*data); ok {
				return sd.Destroy()
			}
			return goerrors.InvalidType
		},
		func(
			conn common.IInputStreamForStack,
			stackData common.IStackCreateData,
			ToReactorFunc rxgo.NextFunc,
		) (common.IInputStreamForStack, error) {
			if sd, ok := stackData.(*data); ok {
				return conn, sd.Start(ctx)
			}
			return nil, goerrors.InvalidType
		},
		func(stackData interface{}) error {
			if sd, ok := stackData.(*data); ok {
				return sd.Stop()
			}
			return goerrors.InvalidType
		},
	)
}
