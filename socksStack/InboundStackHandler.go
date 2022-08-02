package socksStack

import (
	"encoding/binary"
	"fmt"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"io"
	"net"
	"net/url"
)

type OnStateCallback func(
	onNextData func(data []byte) error,
	onNextMessage func(interface{}) error,
) (canContinue bool, err error)

type inboundStackHandler struct {
	Rw           *gomessageblock.ReaderWriter
	errorState   error
	stackData    *data
	currentState OnStateCallback
}

func (self *inboundStackHandler) GetAdditionalBytesIncoming() int {
	return 0
}

func (self *inboundStackHandler) GetAdditionalBytesSend() int {
	return 0
}

func (self *inboundStackHandler) ReadMessage(_ interface{}) (interface{}, bool, error) {
	return nil, false, nil
}

func (self *inboundStackHandler) Close() error {
	return nil
}

func (self *inboundStackHandler) OnError(err error) {
	self.errorState = err
}

func (self *inboundStackHandler) inboundState(
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

func (self *inboundStackHandler) NextReadWriterSize(
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

func (self *inboundStackHandler) OnComplete() {
	if self.errorState == nil {
		self.errorState = RxHandlers.RxHandlerComplete
	}
}

func (self *inboundStackHandler) ReadSocksHeader(next func(byte) OnStateCallback) OnStateCallback {
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

func (self *inboundStackHandler) ReadMethodCount(version byte) OnStateCallback {
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

func (self *inboundStackHandler) ReadMethods(version byte, count byte) OnStateCallback {
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

			self.currentState = self.SendServerMessage(version, b)
			return true, nil
		}
		return false, nil
	}
}

func (self *inboundStackHandler) SendServerMessage(version byte, methods []byte) OnStateCallback {
	return func(
		onNextData func(data []byte) error,
		onNextMessage func(interface{}) error,
	) (canContinue bool, err error) {
		for _, method := range methods {
			if method == 0 {
				b := [2]byte{}
				b[0] = version
				b[1] = 0
				self.stackData.handler.OnSendData(gomessageblock.NewReaderWriterBlock(b[:]))
				break
			}
		}
		self.currentState = self.ReadRequestCommand(version)
		return true, err
	}
}

func (self *inboundStackHandler) ReadRequestCommand(versionRequired byte) OnStateCallback {
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
			self.currentState = self.switchAddressType(b[0], b[1], b[3])
			return true, nil
		}
		return false, nil
	}
}

func (self *inboundStackHandler) switchAddressType(version byte, command byte, addressType byte) OnStateCallback {
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

func (self *inboundStackHandler) CreateDialer(
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
		dd := dialer.RemoteAddr()
		if ipAddress, ok := dd.(*net.TCPAddr); ok {
			reply := Reply{
				Ver:         version,
				Rep:         0,
				Rsv:         0,
				Type:        1,
				BindAddress: ipAddress.IP,
				BindPort:    uint16(ipAddress.Port),
			}
			err := reply.encode(enc)
			if err != nil {
				return false, err
			}
			self.stackData.handler.OnSendData(rws)
			self.currentState = self.PassData()
			return true, nil
		}
		return false, fmt.Errorf("invalid type")
	}
}

func (self *inboundStackHandler) PassData() OnStateCallback {
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

func (self *inboundStackHandler) readIpAddressAndPort(version byte) OnStateCallback {
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

func (self *inboundStackHandler) readDomainNameAndPorts(version byte) OnStateCallback {
	var state func() (bool, error)

	readDestAddressAndPort := func(l byte) func() (bool, error) {
		return func() (bool, error) {
			if self.Rw.Size() >= int(l)+2 {
				b01 := make([]byte, l, l)
				_, _ = self.Rw.Read(b01)
				s := string(b01)
				b02 := [2]byte{}
				_, _ = self.Rw.Read(b02[:])
				portNumber := binary.BigEndian.Uint16(b02[:])
				urlToConnect, err := url.Parse(fmt.Sprintf("tcp4://%v:%d", s, portNumber))
				if err != nil {
					return false, err
				}
				self.currentState = self.CreateDialer(version, urlToConnect)
				return false, nil
			}
			return false, nil
		}

	}
	readLength := func() (bool, error) {
		if self.Rw.Size() >= 1 {
			b := [1]byte{}
			_, _ = self.Rw.Read(b[:])
			l := b[0]
			state = readDestAddressAndPort(l)
			return true, nil

		}
		return false, nil
	}

	state = readLength
	return func(
		onNextData func(data []byte) error,
		onNextMessage func(interface{}) error,
	) (bool, error) {
		var err error
		canContinue := true
		for canContinue {
			canContinue, err = state()
			if err != nil {
				return false, err
			}
		}
		return true, nil
	}
}

type Reply struct {
	Ver         byte
	Rep         byte
	Rsv         byte
	Type        byte
	BindAddress interface{}
	BindPort    uint16
}

func (self Reply) encode(encoder encoder) error {
	_ = encoder.writeByte(self.Ver)
	_ = encoder.writeByte(self.Rep)
	_ = encoder.writeByte(self.Rsv)
	_ = encoder.writeByte(self.Type)
	switch v := self.BindAddress.(type) {
	case string:
		encoder.writeString(v)
		break
	case []byte:
		encoder.writeBytes(v)
		break
	case net.IP:
		encoder.writeBytes(v)
		break
	default:
		return fmt.Errorf("sdsds")
	}
	encoder.writeUInt16(self.BindPort)
	return nil
}

func newInboundStackHandler(stackData *data) RxHandlers.IRxNextStackHandler {
	result := &inboundStackHandler{
		Rw:         gomessageblock.NewReaderWriter(),
		errorState: nil,
		stackData:  stackData,
	}
	result.currentState = result.ReadSocksHeader(result.ReadMethodCount)
	return result
}
