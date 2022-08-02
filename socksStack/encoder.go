package socksStack

import (
	"encoding/binary"
	"github.com/bhbosman/gomessageblock"
)

type encoder struct {
	rws *gomessageblock.ReaderWriter
}

func (self *encoder) writeByte(ver byte) error {
	b := [1]byte{ver}
	_, _ = self.rws.Write(b[:])
	return nil
}

func (self *encoder) writeString(value string) {
	b := []byte(value)
	_ = self.writeByte(byte(len(b)))
	_, _ = self.rws.Write(b)
}

func (self *encoder) writeUInt16(port uint16) {
	b := [2]byte{}
	binary.BigEndian.PutUint16(b[:], port)
	_, _ = self.rws.Write(b[:])
}

func (self *encoder) writeBytes(v []byte) {
	_, _ = self.rws.Write(v)
}
