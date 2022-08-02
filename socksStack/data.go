package socksStack

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/goprotoextra"
)

type data struct {
	bytesSend         int
	requestId         int64
	ctx               context.Context
	goFunctionCounter GoFunctionCounter.IService
	handler           goCommsDefinitions.IRxNextHandler
}

func (self *data) SendError(err error) {
	if self.handler != nil {
		self.handler.OnError(err)
	}
}

func (self *data) SendPong(rws goprotoextra.ReadWriterSize) {
	if self.handler != nil {
		self.bytesSend += rws.Size()
		self.handler.OnSendData(rws)
	}
}

func (self *data) GetBytesSend() int {
	return self.bytesSend
}

func (self *data) SendRws(rws goprotoextra.ReadWriterSize) {
	if self.handler != nil {
		self.handler.OnSendData(rws)
	}
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

func (self *data) Start(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return nil
}

func (self *data) Stop() error {
	return nil
}
