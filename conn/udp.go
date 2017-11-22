package conn

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/google/btree"
	"github.com/skycoin/net/msg"
	"hash/crc32"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

const (
	MAX_UDP_PACKAGE_SIZE = 1200
)

type UDPConn struct {
	ConnCommonFields
	*UDPPendingMap
	*streamQueue
	UdpConn *net.UDPConn
	addr    *net.UDPAddr

	// write loop with ping
	SendPing bool
	rto      time.Duration
	rtt      time.Duration

	rtoResendCount  uint32
	lossResendCount uint32
	ackCount        uint32
	overAckCount    uint32

	// congestion algorithm
	*ca
}

type mode int

const (
	startup mode = iota
	drain
	probeBW
)

// used for server spawn udp conn
func NewUDPConn(c *net.UDPConn, addr *net.UDPAddr) *UDPConn {
	conn := &UDPConn{
		UdpConn:          c,
		addr:             addr,
		ConnCommonFields: NewConnCommonFileds(),
		UDPPendingMap:    NewUDPPendingMap(),
		streamQueue:      newStreamQueue(),
		rto:              300 * time.Millisecond,

		ca: newCA(),
	}
	return conn
}

func (c *UDPConn) ReadLoop() error {
	return nil
}

func (c *UDPConn) WriteLoop() (err error) {
	if c.SendPing {
		err = c.writeLoopWithPing()
	} else {
		err = c.writeLoop()
	}
	c.GetContextLogger().Debugf("%s", c.String())
	return
}

func (c *UDPConn) writeLoop() (err error) {
	defer func() {
		if err != nil {
			c.SetStatusToError(err)
		}
	}()
	for {
		select {
		case m, ok := <-c.Out:
			if !ok {
				c.GetContextLogger().Debug("udp conn closed")
				return nil
			}
			err := c.Write(m)
			if err != nil {
				c.GetContextLogger().Debugf("write msg is failed %v", err)
				return err
			}
		}
	}
}

func (c *UDPConn) writeLoopWithPing() (err error) {
	ticker := time.NewTicker(time.Second * UDP_PING_TICK_PERIOD)
	defer func() {
		ticker.Stop()
		if err != nil {
			c.SetStatusToError(err)
		}
	}()

	for {
		select {
		case <-ticker.C:
			nowUnix := time.Now().Unix()
			if nowUnix-c.GetLastTime() >= UDP_GC_PERIOD {
				c.Close()
				return errors.New("timeout")
			}
			err := c.Ping()
			if err != nil {
				return err
			}
		case m, ok := <-c.Out:
			if !ok {
				c.GetContextLogger().Debug("udp conn closed")
				return nil
			}
			err := c.Write(m)
			if err != nil {
				c.GetContextLogger().Debugf("write msg is failed %v", err)
				return err
			}
		}
	}
}

func (c *UDPConn) Write(bytes []byte) (err error) {
	s := c.GetNextSeq()
	m := msg.NewUDP(msg.TYPE_NORMAL, s, bytes, c.getDelivered(), c.getDeliveryTime())
	c.AddMsg(s, m)
	ok := c.addBytesInFlight(m)
	if !ok {
		return nil
	}
	c.GetContextLogger().Debugf("new msg seq %d", m.Seq)
	err = c.WriteBytes(m.PkgBytes())
	c.transmitted(m)
	return
}

func (c *UDPConn) transmitted(m *msg.UDPMessage) {
	m.Transmitted()
	m.SetRTO(c.getRTO(), func() (err error) {
		c.AddLossResendCount()
		err = c.WriteBytes(m.PkgBytes())
		if err != nil {
			c.SetStatusToError(err)
			c.Close()
		}
		return
	})
}

func (c *UDPConn) resendMsg(m *msg.UDPMessage) error {
	if !c.UDPPendingMap.exists(m.Seq) {
		return nil
	}
	c.GetContextLogger().Debugf("resendMsg %s", m)
	return c.WriteBytes(m.PkgBytes())
}

func (c *UDPConn) writePendingMsgs() error {
	for {
		m := c.ca.popMessage()
		if m == nil {
			return nil
		}
		ok := c.addBytesInFlight(m)
		if !ok {
			return nil
		}
		c.GetContextLogger().Debugf("new msg seq %d", m.Seq)
		err := c.WriteBytes(m.PkgBytes())
		if err != nil {
			return err
		}
		c.transmitted(m)
	}
}

func (c *UDPConn) WriteBytes(bytes []byte) error {
	l := len(bytes)
	c.AddSentBytes(l)
	c.WriteMutex.Lock()
	defer c.WriteMutex.Unlock()
	n, err := c.UdpConn.WriteToUDP(bytes, c.addr)
	if err == nil && n != l {
		return errors.New("nothing was written")
	}
	return err
}

func (c *UDPConn) Ack(seq uint32) error {
	nSeq := c.getNextAckSeq()
	c.GetContextLogger().Debugf("ack %d, next %d", seq, nSeq)
	p := make([]byte, msg.ACK_HEADER_SIZE+msg.PKG_HEADER_SIZE)
	m := p[msg.PKG_HEADER_SIZE:]
	m[msg.ACK_TYPE_BEGIN] = msg.TYPE_ACK
	binary.BigEndian.PutUint32(m[msg.ACK_SEQ_BEGIN:], seq)
	binary.BigEndian.PutUint32(m[msg.ACK_NEXT_SEQ_BEGIN:], nSeq)
	checksum := crc32.ChecksumIEEE(m)
	binary.BigEndian.PutUint32(p[msg.PKG_CRC32_BEGIN:], checksum)
	return c.WriteBytes(p)
}

func (c *UDPConn) RecvAck(m []byte) (err error) {
	if len(m) < msg.ACK_HEADER_SIZE {
		return fmt.Errorf("invalid ack msg %x", m)
	}
	seq := binary.BigEndian.Uint32(m[msg.ACK_SEQ_BEGIN:msg.ACK_SEQ_END])
	ns := binary.BigEndian.Uint32(m[msg.ACK_NEXT_SEQ_BEGIN:msg.ACK_NEXT_SEQ_END])
	n, ok := c.getMinUnAckSeq()

	c.GetContextLogger().Debugf("recv ack %d, next %d, min unack %t %d", seq, ns, ok, n)
	if ok {
		for ; ns > n+1; n++ {
			c.GetContextLogger().Debugf("ignore ack %d", n)
			err = c.delMsg(n, true)
			if err != nil {
				return
			}
		}
	}
	err = c.delMsg(seq, false)
	return
}

func (c *UDPConn) Ping() error {
	c.GetContextLogger().Debug("ping")
	p := make([]byte, msg.PING_MSG_HEADER_SIZE+msg.PKG_HEADER_SIZE)
	m := p[msg.PKG_HEADER_SIZE:]
	m[msg.PING_MSG_TYPE_BEGIN] = msg.TYPE_PING
	binary.BigEndian.PutUint64(m[msg.PING_MSG_TIME_BEGIN:], msg.UnixMillisecond())
	checksum := crc32.ChecksumIEEE(m)
	binary.BigEndian.PutUint32(p[msg.PKG_CRC32_BEGIN:], checksum)
	return c.WriteBytes(p)
}

func (c *UDPConn) GetChanOut() chan<- []byte {
	return c.Out
}

func (c *UDPConn) GetChanIn() <-chan []byte {
	return c.In
}

func (c *UDPConn) GetNextSeq() uint32 {
	return atomic.AddUint32(&c.seq, 1)
}

func (c *UDPConn) Close() {
	c.ConnCommonFields.Close()
}

func (c *UDPConn) String() string {
	return fmt.Sprintf(
		`udp connection(%s):
			rtoResend:%d,
			lossResend:%d,
			ack:%d,
			overAck:%d,`,
		c.GetRemoteAddr().String(),
		atomic.LoadUint32(&c.rtoResendCount),
		atomic.LoadUint32(&c.lossResendCount),
		atomic.LoadUint32(&c.ackCount),
		atomic.LoadUint32(&c.overAckCount),
	)
}

func (c *UDPConn) GetRemoteAddr() net.Addr {
	return c.addr
}

func (c *UDPConn) getRTO() (rto time.Duration) {
	c.FieldsMutex.RLock()
	rto = c.rto
	c.FieldsMutex.RUnlock()
	return
}

func (c *UDPConn) setRTO(rto time.Duration) {
	c.GetContextLogger().Debugf("setRTO %d", rto)
	c.FieldsMutex.Lock()
	c.rto = rto
	c.FieldsMutex.Unlock()
}

func (c *UDPConn) AddMsg(k uint32, v *msg.UDPMessage) {
	c.UDPPendingMap.AddMsg(k, v)
}

func (c *UDPConn) delMsg(seq uint32, ignore bool) error {
	c.AddAckCount()
	ok, um, msgs := c.DelMsgAndGetLossMsgs(seq)
	if ok {
		if !ignore {
			//c.updateRTT(um.GetRTT())
			//c.updateDeliveryRate(um)
		}
		c.subBytesInFlight(um.TotalSize())
		if len(msgs) > 1 {
			c.GetContextLogger().Debugf("resend loss msgs %v", msgs)
			for _, msg := range msgs {
				err := c.resendMsg(msg)
				if err != nil {
					c.SetStatusToError(err)
					c.Close()
					return err
				}
				c.AddLossResendCount()
			}
		}
		c.UpdateLastAck(seq)
	} else {
		c.GetContextLogger().Debugf("over ack %s", c)
		c.AddOverAckCount()
	}
	return c.writePendingMsgs()
}

func (c *UDPConn) AddLossResendCount() {
	atomic.AddUint32(&c.lossResendCount, 1)
}

func (c *UDPConn) AddResendCount() {
	atomic.AddUint32(&c.rtoResendCount, 1)
}

func (c *UDPConn) AddAckCount() {
	atomic.AddUint32(&c.ackCount, 1)
}

func (c *UDPConn) AddOverAckCount() {
	atomic.AddUint32(&c.overAckCount, 1)
}

func (c *UDPConn) IsTCP() bool {
	return false
}

func (c *UDPConn) IsUDP() bool {
	return true
}

func (c *UDPConn) getRTT() time.Duration {
	return time.Duration(atomic.LoadInt64((*int64)(&c.rtt)))
}

type rttSampler struct {
	tree  *btree.BTree
	ring  []rtt
	mask  int
	index int
}

type rtt time.Duration

func (a rtt) Less(b btree.Item) bool {
	return a < b.(rtt)
}

// size should be power of 2
func newRttSampler(size int) *rttSampler {
	if size < 2 || (size&(size-1)) > 0 {
		var n uint
		for size > 0 {
			size >>= 1
			n++
		}
		size = 1 << n
	}
	return &rttSampler{
		ring: make([]rtt, size),
		mask: size - 1,
		tree: btree.New(2),
	}
}

func (t *rttSampler) push(r rtt) rtt {
	if r <= 0 {
		panic("push rtt <= 0")
	}
	or := t.ring[t.index]
	if or > 0 {
		t.tree.Delete(or)
	}
	t.ring[t.index] = r
	t.tree.ReplaceOrInsert(r)
	t.index = (t.index + 1) & t.mask
	return t.tree.Min().(rtt)
}

func (t *rttSampler) getMin() rtt {
	return t.tree.Min().(rtt)
}

func (c *UDPConn) updateRTT(t time.Duration) {
	if t <= 0 {
		panic("updateRTT t <= 0")
	}
	r := c.rttSamples.push(rtt(t))
	for {
		ot := c.getRTT()
		c.GetContextLogger().Debugf("updateRTT ot %d new %d", ot, r)
		if ot == 0 || time.Duration(r) < ot {
			ok := atomic.CompareAndSwapInt64((*int64)(&c.rtt), int64(ot), int64(r))
			if !ok {
				continue
			}
			c.setRTO(t * 4)
		}
		break
	}
	c.GetContextLogger().Debug("updateRTT end")
}

type rateSampler struct {
	tree  *btree.BTree
	ring  []rate
	mask  int
	index int
}

type rate uint64

func (a rate) Less(b btree.Item) bool {
	return a < b.(rate)
}

// size should be power of 2
func newRateSampler(size int) *rateSampler {
	if size < 2 || (size&(size-1)) > 0 {
		var n uint
		for size > 0 {
			size >>= 1
			n++
		}
		size = 1 << n
	}
	return &rateSampler{
		ring: make([]rate, size),
		mask: size - 1,
		tree: btree.New(2),
	}
}

func (t *rateSampler) push(r rate) rate {
	if r <= 0 {
		panic("push rate <= 0")
	}
	or := t.ring[t.index]
	if or > 0 {
		t.tree.Delete(or)
	}
	t.ring[t.index] = r
	t.tree.ReplaceOrInsert(r)
	t.index = (t.index + 1) & t.mask
	return t.tree.Max().(rate)
}

func (c *UDPConn) updateDeliveryRate(m *msg.UDPMessage) {
	c.ca.Lock()
	defer c.ca.Unlock()
	c.delivered += uint64(m.TotalSize())
	if c.deliveryTime.IsZero() {
		c.deliveryTime = time.Now()
		return
	}
	d := c.delivered
	if d <= 0 {
		return
	}
	dt := time.Now().Sub(c.deliveryTime) / time.Millisecond
	c.GetContextLogger().Debugf("drate d %d dt %d", d, dt)
	drate := d / uint64(dt)
	if drate <= 0 {
		return
	}
	max := uint64(c.rateSamples.push(rate(drate)))
	if max != drate && c.ca.mode == startup {
		c.fullCnt++
		if c.fullCnt > 3 {
			c.ca.mode = drain
			c.ca.pacingGain = drainGain
		}
	}
	rtt := uint64(c.rttSamples.getMin()) / uint64(time.Millisecond)
	cwnd := int(float64(max*rtt) * c.pacingGain)
	c.GetContextLogger().Debugf("mode %d, max bytes %d rtt %d: cwnd %d", c.mode, max, rtt, cwnd)
	//if c.ca.mode == startup && cwnd > c.cwnd {
	//	c.cwnd = cwnd
	//} else if c.ca.mode == drain {
	//	pcwnd := int(max * rtt)
	//	if c.ca.bytesInFlight < int(pcwnd) {
	//		c.ca.mode = probeBW
	//		cwnd = pcwnd
	//	}
	//	c.cwnd = cwnd
	//}
}

type ca struct {
	delivered    uint64
	deliveryTime time.Time
	rttSamples   *rttSampler
	rateSamples  *rateSampler
	cwnd         int
	mode
	pacingGain    float64
	fullCnt       uint
	bytesInFlight int

	buff *btree.BTree
	sync.RWMutex
}

func newCA() *ca {
	c := &ca{
		rttSamples:  newRttSampler(16),
		rateSamples: newRateSampler(16),
		pacingGain:  highGain,
		cwnd:        20480,
	}
	c.buff = btree.New(2)
	return c
}

func (ca *ca) getDelivered() (d uint64) {
	ca.RLock()
	d = ca.delivered
	ca.RUnlock()
	return
}

func (ca *ca) getDeliveryTime() (d time.Time) {
	ca.RLock()
	d = ca.deliveryTime
	ca.RUnlock()
	return
}

func (ca *ca) addBytesInFlight(m *msg.UDPMessage) (ok bool) {
	bytes := m.PkgBytes()
	ca.Lock()
	defer ca.Unlock()

	min := ca.buff.Min()
	if (min != nil && min.Less(m)) || ca.cwnd < ca.bytesInFlight+len(bytes) {
		ca.buff.ReplaceOrInsert(m)
		return
	}
	ca.bytesInFlight += len(bytes)
	ok = true
	return
}

func (ca *ca) subBytesInFlight(s int) {
	ca.Lock()
	if ca.bytesInFlight < s {
		panic("subBytesInFlight ca.bytesInFlight < s")
	}
	ca.bytesInFlight -= s
	ca.Unlock()
}

func (ca *ca) popMessage() (m *msg.UDPMessage) {
	ca.Lock()
	defer ca.Unlock()

	if ca.buff.Len() < 1 {
		return
	}
	element := ca.buff.Min()
	if element == nil {
		return
	}
	ca.buff.DeleteMin()
	m = element.(*msg.UDPMessage)
	return
}

func (ca *ca) getBytesInFlight() (r int) {
	ca.RLock()
	r = ca.bytesInFlight
	ca.RUnlock()
	return
}
