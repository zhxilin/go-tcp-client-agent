package middleware

import (
	"bytes"
	"encoding/binary"
	"errors"
	msg "go-tcp-client-agent/core/proto"
	"time"

	"google.golang.org/protobuf/proto"
)

const (
	cMsgHeadLen = 18
)

type gtaMsgHead struct {
	msgID     uint16
	seq       uint32
	ts        uint64
	totalSize uint32
}

type GtaMsgParser struct {
	msgSeq int32
}

func NewMsgParser() *GtaMsgParser {
	return &GtaMsgParser{}
}

func getNowTs() int64 {
	return time.Now().UnixNano() / 1e6
}

func (p *GtaMsgParser) packMsgBase(m proto.Message) ([]byte, int) {
	pb, err := proto.Marshal(m)
	if err != nil {
		return nil, 0
	}

	data := &msg.MsgBase{
		Payload: pb,
	}
	buffer, _ := proto.Marshal(data)
	len := binary.Size(buffer)

	return buffer, len
}

func (p *GtaMsgParser) Pack(msgID uint16, m proto.Message) []byte {

	body, bodyLen := p.packMsgBase(m)

	p.msgSeq++

	head := &gtaMsgHead{
		msgID:     msgID,
		seq:       uint32(p.msgSeq),
		ts:        uint64(getNowTs()),
		totalSize: uint32(bodyLen + cMsgHeadLen),
	}

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, head)
	binary.Write(buf, binary.BigEndian, body)

	return buf.Bytes()
}

func (cli *GtaMsgParser) UnPack(buffer []byte, len int) (uint16, []byte, int, error) {
	if len < cMsgHeadLen {
		return 0, nil, 0, errors.New("invalid msg length")
	}

	head := &gtaMsgHead{}
	buf := bytes.NewBuffer(buffer)

	binary.Read(buf, binary.BigEndian, &head.msgID)
	binary.Read(buf, binary.BigEndian, &head.seq)
	binary.Read(buf, binary.BigEndian, &head.ts)
	binary.Read(buf, binary.BigEndian, &head.totalSize)

	body := make([]byte, head.totalSize-cMsgHeadLen)

	binary.Read(buf, binary.BigEndian, &body)

	data := &msg.MsgBase{}
	proto.Unmarshal(body, data)

	payload := data.Payload

	return head.msgID, payload, int(head.totalSize), nil
}
