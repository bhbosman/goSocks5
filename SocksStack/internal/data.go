package internal

import (
	"context"
	pingpong "github.com/bhbosman/goMessages/pingpong/stream"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/stream"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/goprotoextra"
	"github.com/reactivex/rxgo/v2"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type data struct {
	bytesSend         int
	requestId         int64
	ctx               context.Context
	onRxSendData      func(data interface{})
	onRxSendError     func(err error)
	onRxComplete      func()
	goFunctionCounter GoFunctionCounter.IService
}

func (self *data) SendError(err error) {
	if self.onRxSendError != nil {
		self.onRxSendError(err)
	}
}

func (self *data) SendPong(rws goprotoextra.ReadWriterSize) {
	if self.onRxSendData != nil {
		self.bytesSend += rws.Size()
		self.onRxSendData(rws)
	}
}

func (self *data) GetBytesSend() int {
	return self.bytesSend
}

func (self *data) SendRws(rws goprotoextra.ReadWriterSize) {
	if self.onRxSendData != nil {
		self.onRxSendData(rws)
	}
}

func (self *data) PongReceived(_ *pingpong.Pong) {
}

func NewStackData(
	ctx context.Context,
	goFunctionCounter GoFunctionCounter.IService,
) *data {
	return &data{
		requestId:         0,
		ctx:               ctx,
		goFunctionCounter: goFunctionCounter,
	}
}

func (self *data) Destroy() error {
	return nil
}

func (self *data) Ping() {
	sendTime := time.Now()
	msg := &pingpong.Ping{
		RequestId: 0,
		RequestTimeStamp: &timestamppb.Timestamp{
			Seconds: sendTime.Unix(),
			Nanos:   int32(sendTime.Nanosecond()),
		},
	}
	rws, err := stream.Marshall(msg)
	if err != nil {
		return
	}
	self.bytesSend += rws.Size()
	if self.onRxSendData != nil {
		self.onRxSendData(rws)
	}
}

func (self *data) Start(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return nil
}

func (self *data) Stop() error {
	return nil
}

func (self *data) SetOnRxSendData(data rxgo.NextFunc) error {
	if data == nil {
		return goerrors.InvalidParam
	}
	self.onRxSendData = data
	return nil
}

func (self *data) setOnRxSendError(sendError rxgo.ErrFunc) error {
	if sendError == nil {
		return goerrors.InvalidParam
	}
	self.onRxSendError = sendError
	return nil
}

func (self *data) setOnRxComplete(complete rxgo.CompletedFunc) error {
	if complete == nil {
		return goerrors.InvalidParam
	}
	self.onRxComplete = complete
	return nil
}
