package providers

import (
	"github.com/bhbosman/gocommon/messageRouter"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"io"
	"net"
)

type reactor struct {
	common.BaseConnectionReactor
	messageRouter *messageRouter.MessageRouter
	conn          *net.TCPConn
}

func (self *reactor) handleNetTCPConn(message *net.TCPConn) {
	self.conn = message
	go common.ReadFromIoReader(
		"dddd",
		self.conn,
		self.CancelCtx,
		self.ConnectionCancelFunc,
		RxHandlers.NewDefaultRxNextHandler(
			func(i interface{}) {
				if rws, ok := i.(*gomessageblock.ReaderWriter); ok {
					self.ToConnection(rws)
				}
			},
			func(err error) {
				self.ConnectionCancelFunc("sadassa", false, err)
			},
			func() {
				self.CancelFunc()
			},
		),
	)

}

func (self *reactor) handleEmptyQueue(message *messages.EmptyQueue) {
}

func (self *reactor) handleRws(message *gomessageblock.ReaderWriter) {
	_, _ = io.Copy(self.conn, message)
}

func newReactor(
	logger *zap.Logger,
	cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	connectionCancelFunc model.ConnectionCancelFunc,
	userContext interface{}) (*reactor, error) {
	result := &reactor{
		BaseConnectionReactor: common.NewBaseConnectionReactor(
			logger,
			cancelCtx,
			cancelFunc,
			connectionCancelFunc,
			userContext),
		messageRouter: messageRouter.NewMessageRouter(),
	}
	_ = result.messageRouter.Add(result.handleNetTCPConn)
	_ = result.messageRouter.Add(result.handleEmptyQueue)
	_ = result.messageRouter.Add(result.handleRws)

	return result, nil
}

func (self *reactor) Close() error {
	return self.BaseConnectionReactor.Close()
}

func (self *reactor) Init(
	onSend goprotoextra.ToConnectionFunc,
	toConnectionReactor goprotoextra.ToReactorFunc,
	onSendReplacement rxgo.NextFunc,
	toConnectionReactorReplacement rxgo.NextFunc) (rxgo.NextFunc, rxgo.ErrFunc, rxgo.CompletedFunc, error) {

	_, _, _, err := self.BaseConnectionReactor.Init(
		onSend,
		toConnectionReactor,
		onSendReplacement,
		toConnectionReactorReplacement)
	if err != nil {
		return nil, nil, nil, err
	}
	return func(i interface{}) {
			self.messageRouter.Route(i)
		},
		func(err error) {
			self.messageRouter.Route(err)
		}, func() {
			if self.conn != nil {

				self.conn.Close()
			}

		}, nil

}

func (self *reactor) Open() error {
	return nil
}
