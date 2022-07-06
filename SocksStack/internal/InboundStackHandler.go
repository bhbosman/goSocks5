package internal

import (
	"encoding/binary"
	"fmt"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"io"
	"net"
	"net/url"
	"strconv"
)

type OnStateCallback func(
	onNextData func(data []byte) error,
	onNextMessage func(interface{}) error,
) (canContinue bool, err error)

type InboundStackHandler struct {
	Rw           *gomessageblock.ReaderWriter
	errorState   error
	stackData    *data
	currentState OnStateCallback
}

func (self *InboundStackHandler) GetAdditionalBytesIncoming() int {
	return 0
}

func (self *InboundStackHandler) GetAdditionalBytesSend() int {
	return 0
}

func (self *InboundStackHandler) ReadMessage(_ interface{}) error {
	return nil
}

func (self *InboundStackHandler) Close() error {
	return nil
}

func (self *InboundStackHandler) OnError(err error) {
	self.errorState = err
}

func (self *InboundStackHandler) inboundState(
	onNextData func(data []byte) error,
	onNextMessage func(interface{}) error,
) error {
	var err error
	canContinue := true
	for canContinue {
		canContinue, err = self.currentState(onNextData, onNextMessage)
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *InboundStackHandler) NextReadWriterSize(
	rws goprotoextra.ReadWriterSize,
	onNextRws func(rws goprotoextra.ReadWriterSize) error,
	onNextMessage func(interface{}) error,
	_ func(size int) error,
) error {
	if self.errorState != nil {
		return self.errorState
	}
	err := self.Rw.SetNext(rws)
	if err != nil {
		return err
	}
	return self.inboundState(
		func(data []byte) error {
			return onNextRws(gomessageblock.NewReaderWriterBlock(data))
		},
		func(i interface{}) error {
			return onNextMessage(i)
		},
	)
}

func (self *InboundStackHandler) OnComplete() {
	if self.errorState == nil {
		self.errorState = RxHandlers.RxHandlerComplete
	}
}

func (self *InboundStackHandler) ReadSocksHeader(next func(byte) OnStateCallback) OnStateCallback {
	return func(
		onNextData func(data []byte) error,
		onNextMessage func(interface{}) error,
	) (bool, error) {
		if self.Rw.Size() >= 1 {
			b := [1]byte{}
			n, err := self.Rw.Read(b[:])
			if err != nil {
				return false, err
			}
			if n != 1 {
				return false, io.ErrShortBuffer
			}
			if b[0] == 0x05 {
				self.currentState = next(b[0])
				return true, nil
			}
			return false, fmt.Errorf("wrong socks version. expected 5")
		}
		return false, nil
	}
}

func (self *InboundStackHandler) ReadMethodCount(version byte) OnStateCallback {
	return func(
		onNextData func(data []byte) error,
		onNextMessage func(interface{}) error,
	) (canContinue bool, err error) {
		if self.Rw.Size() >= 1 {
			b := [1]byte{}
			n, err := self.Rw.Read(b[:])
			if err != nil {
				return false, err
			}
			if n != 1 {
				return false, io.ErrShortBuffer
			}

			self.currentState = self.ReadMethods(version, b[0])
			return true, nil
		}
		return false, nil
	}
}

func (self *InboundStackHandler) ReadMethods(version byte, count byte) OnStateCallback {
	return func(
		onNextData func(data []byte) error,
		onNextMessage func(interface{}) error,
	) (canContinue bool, err error) {
		if self.Rw.Size() >= int(count) {
			b := make([]byte, count)
			n, err := self.Rw.Read(b)
			if err != nil {
				return false, err
			}
			if n != int(count) {
				return false, io.ErrShortBuffer
			}

			self.currentState = self.SendServerMessage(version, count, b)
			return true, nil
		}
		return false, nil
	}
}

func (self *InboundStackHandler) SendServerMessage(version byte, count byte, methods []byte) OnStateCallback {
	return func(
		onNextData func(data []byte) error,
		onNextMessage func(interface{}) error,
	) (canContinue bool, err error) {
		for _, method := range methods {
			if method == 0 {
				b := [2]byte{}
				b[0] = version
				b[1] = 0
				self.stackData.onRxSendData(gomessageblock.NewReaderWriterBlock(b[:]))
				break
			}
		}
		self.currentState = self.ReadRequestCommand(version)
		return true, err
	}
}

func (self *InboundStackHandler) ReadRequestCommand(versionRequired byte) OnStateCallback {
	return func(
		onNextData func(data []byte) error,
		onNextMessage func(interface{}) error,
	) (canContinue bool, err error) {
		if self.Rw.Size() >= 4 {
			b := [4]byte{}
			n, err := self.Rw.Read(b[:])
			if err != nil {
				return false, err
			}
			if n != 4 {
				return false, io.ErrShortBuffer
			}
			if b[0] != versionRequired {
				return false, fmt.Errorf("wrong socks version. expected 5")
			}
			self.currentState = self.switchAddressType(b[0], b[1], b[2], b[3])
			return true, nil
		}
		return false, nil
	}
}

func (self *InboundStackHandler) switchAddressType(
	version byte,
	command byte,
	reserved byte,
	addressType byte,
) OnStateCallback {
	return func(
		onNextData func(data []byte) error,
		onNextMessage func(interface{}) error,
	) (canContinue bool, err error) {
		switch command {
		case 1:
			switch addressType {
			case 1:
				self.currentState = self.readIpAddressAndPort(version)
				return true, nil
			case 3:
				self.currentState = self.readDomainNameAndPorts(version)
				return true, nil
			default:
				return false, fmt.Errorf("not supported")
			}
		default:
			return false, fmt.Errorf("not supported")
		}
	}
}

func (self *InboundStackHandler) CreateDialer(
	version byte,
	connect *url.URL) OnStateCallback {
	return func(
		onNextData func(data []byte) error,
		onNextMessage func(interface{}) error,
	) (canContinue bool, err error) {
		var dialer net.Conn
		dialer, err = net.Dial(connect.Scheme, connect.Host)
		if err != nil {
			return false, err
		}
		_ = onNextMessage(dialer)
		rws := gomessageblock.NewReaderWriter()
		enc := encoder{rws: rws}
		u, _ := url.Parse(fmt.Sprintf("%v://%v", dialer.RemoteAddr().Network(), dialer.RemoteAddr().String()))
		p, _ := strconv.Atoi(u.Port())
		reply := Reply{
			Ver:         version,
			Rep:         0,
			Rsv:         0,
			Type:        1,
			BindAddress: u.Host,
			BindPort:    uint16(p),
		}
		reply.encode(enc)
		self.stackData.onRxSendData(rws)
		self.currentState = self.PassData()
		return true, nil
	}
}

func (self *InboundStackHandler) PassData() OnStateCallback {
	return func(
		onNextData func(data []byte) error,
		onNextMessage func(interface{}) error,
	) (bool, error) {
		length := self.Rw.Size()
		if length > 0 {
			dataBlock := make([]byte, length)
			_, err := self.Rw.Read(dataBlock)
			if err != nil {
				self.errorState = err
				return false, self.errorState
			}
			err = onNextData(dataBlock)
		}
		return self.Rw.Size() > 0, nil
	}
}

func (self *InboundStackHandler) readIpAddressAndPort(version byte) OnStateCallback {
	return func(
		onNextData func(data []byte) error,
		onNextMessage func(interface{}) error,
	) (canContinue bool, err error) {
		if self.Rw.Size() >= 6 {
			b := [4]byte{}
			_, _ = self.Rw.Read(b[:])
			s := net.IPv4(b[0], b[1], b[2], b[3]).String()
			_, _ = self.Rw.Read(b[0:2])
			portNumber := binary.BigEndian.Uint16(b[0:2])
			var urlToConnect *url.URL
			urlToConnect, err = url.Parse(fmt.Sprintf("tcp4://%v:%d", s, portNumber))
			self.currentState = self.CreateDialer(version, urlToConnect)
			return true, nil
		}
		return false, err

	}
}

func (self *InboundStackHandler) readDomainNameAndPorts(version byte) OnStateCallback {
	return func(
		onNextData func(data []byte) error,
		onNextMessage func(interface{}) error,
	) (canContinue bool, err error) {
		return false, err
	}
}

type Reply struct {
	Ver         byte
	Rep         byte
	Rsv         byte
	Type        byte
	BindAddress string
	BindPort    uint16
}

func (self Reply) encode(encoder encoder) {
	encoder.WriteByte(self.Ver)
	encoder.WriteByte(self.Rep)
	encoder.WriteByte(self.Rsv)
	encoder.WriteByte(self.Type)
	encoder.WriteString(self.BindAddress)
	encoder.WriteUInt16(self.BindPort)
}

type encoder struct {
	rws *gomessageblock.ReaderWriter
}

func (self *encoder) WriteByte(ver byte) {
	b := [1]byte{ver}
	_, _ = self.rws.Write(b[:])
}

func (self *encoder) WriteString(value string) {
	b := []byte(value)
	self.WriteByte(byte(len(b)))
	_, _ = self.rws.Write(b)
}

func (self *encoder) WriteUInt16(port uint16) {
	b := [2]byte{}
	binary.BigEndian.PutUint16(b[:], port)
	_, _ = self.rws.Write(b[:])
}

func NewInboundStackHandler(stackData *data) RxHandlers.IRxNextStackHandler {
	result := &InboundStackHandler{
		Rw:         gomessageblock.NewReaderWriter(),
		errorState: nil,
		stackData:  stackData,
	}
	result.currentState = result.ReadSocksHeader(result.ReadMethodCount)
	return result
}
