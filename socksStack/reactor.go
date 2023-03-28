package socksStack

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/messageRouter"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocommon/services/interfaces"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/gomessageblock"
	"github.com/cskr/pubsub"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/zap"
	"io"
	"net"
)

type reactor struct {
	common.BaseConnectionReactor
	messageRouter     messageRouter.IMessageRouter
	conn              *net.TCPConn
	goFunctionCounter GoFunctionCounter.IService
}

func (self *reactor) handleNetTCPConn(message *net.TCPConn) {
	self.conn = message
	go common.ReadFromIoReader(
		self.conn,
		self.CancelCtx,
		self.CancelFunc,
		func() goCommsDefinitions.IRxNextHandler {
			r, _ := goCommsDefinitions.NewDefaultRxNextHandler(
				func(i interface{}) {
					if rws, ok := i.(*gomessageblock.ReaderWriter); ok {
						self.OnSendToConnection(rws)
					}
				},
				func(i interface{}) bool {
					return false
				},
				func(err error) {
					self.ConnectionCancelFunc("sadassa", false, err)
				},
				func() {
					self.CancelFunc()
				},
				func() bool {
					return true
				},
			)
			return r
		}(),
	)
}

func (self *reactor) handleEmptyQueue(_ *messages.EmptyQueue) {
}

func (self *reactor) handleRws(message *gomessageblock.ReaderWriter) {
	_, _ = io.Copy(self.conn, message)
}

func newReactor(
	logger *zap.Logger,
	cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	connectionCancelFunc model.ConnectionCancelFunc,
	goFunctionCounter GoFunctionCounter.IService,
	UniqueReferenceService interfaces.IUniqueReferenceService,
	PubSub *pubsub.PubSub,
) (intf.IConnectionReactor, error) {
	result := &reactor{
		BaseConnectionReactor: common.NewBaseConnectionReactor(
			logger,
			cancelCtx,
			cancelFunc,
			connectionCancelFunc,
			UniqueReferenceService.Next("ConnectionReactor"),
			PubSub,
			goFunctionCounter,
		),
		messageRouter:     messageRouter.NewMessageRouter(),
		goFunctionCounter: goFunctionCounter,
	}
	_ = result.messageRouter.Add(result.handleNetTCPConn)
	_ = result.messageRouter.Add(result.handleEmptyQueue)
	_ = result.messageRouter.Add(result.handleRws)

	return result, nil
}

func (self *reactor) Close() error {
	return self.BaseConnectionReactor.Close()
}

func (self *reactor) Init(params intf.IInitParams) (rxgo.NextFunc, rxgo.ErrFunc, rxgo.CompletedFunc, error) {
	_, _, _, err := self.BaseConnectionReactor.Init(params)
	if err != nil {
		return nil, nil, nil, err
	}
	return func(i interface{}) {
			self.messageRouter.Route(i)
		},
		func(err error) {
			self.messageRouter.Route(err)
		}, func() {
			self.CancelFunc()
		},
		nil
}

func (self *reactor) Open() error {
	return nil
}
