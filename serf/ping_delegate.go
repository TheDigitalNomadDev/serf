package serf

import (
	"bytes"
	"fmt"
	"time"

	"github.com/hashicorp/go-msgpack/codec"
	"github.com/hashicorp/memberlist"
	"github.com/hashicorp/serf/coordinate"
)

type pingDelegate struct {
	serf *Serf
}

func (self *pingDelegate) AckPayload() []byte {
	var buf bytes.Buffer
	fmt.Printf("coord: %v\n", self.serf.coord)
	enc := codec.NewEncoder(&buf, &codec.MsgpackHandle{})
	if err := enc.Encode(self.serf.coord); err != nil {
		panic(fmt.Sprintf("[ERR] serf: Failed to encode coordinate: %v\n", err))
	}
	fmt.Printf("encoded bytes: %v\n", buf.Bytes())
	return buf.Bytes()
}

func (self *pingDelegate) NotifyPingComplete(other *memberlist.Node, rtt time.Duration, payload []byte) {
	var coord coordinate.Client
	r := bytes.NewReader(payload)
	dec := codec.NewDecoder(r, &codec.MsgpackHandle{})
	if err := dec.Decode(&coord); err != nil {
		panic(fmt.Sprintf("[ERR] serf: Failed to decode coordinate: %v", err))
	}
	self.serf.coord.Update(&coord, rtt)
}
