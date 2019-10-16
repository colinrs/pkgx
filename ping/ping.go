package ping

import (
	"encoding/binary"
	"math/rand"
	"net"
	"time"
	"os"
	"fmt"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const (

	// IPv4 ...
	IPv4 = 1
 
)


// IcmpPkg ...
type IcmpPkg struct {
	conn     net.PacketConn
	ipv4conn *ipv4.PacketConn
	msg      icmp.Message
	netmsg   []byte
	id       int
	seq      int
	maxrtt   time.Duration
	dest     net.Addr
}

// ICMP ...
type ICMP struct {
	Addr    net.Addr
	RTT     time.Duration
	MaxRTT  time.Duration
	MinRTT  time.Duration
	AvgRTT  time.Duration
	Final   bool
	Timeout bool
	Down    bool
	Error   error
}

// Send ...
func (icmpPkg *IcmpPkg) Send(ttl int) (hicmp ICMP) {



	icmpPkg.conn, hicmp.Error = net.ListenPacket("ip4:icmp", "0.0.0.0")
	if nil != hicmp.Error {
		return
	}
	defer icmpPkg.conn.Close()
	icmpPkg.ipv4conn = ipv4.NewPacketConn(icmpPkg.conn)
	defer icmpPkg.ipv4conn.Close()
	hicmp.Error = icmpPkg.conn.SetReadDeadline(time.Now().Add(icmpPkg.maxrtt))
	if nil != hicmp.Error {
		return
	}
	if nil != icmpPkg.ipv4conn {
		hicmp.Error = icmpPkg.ipv4conn.SetTTL(ttl)
	}
	if nil != hicmp.Error {
		return
	}
	sendOn := time.Now()
	if nil != icmpPkg.ipv4conn {
		_, hicmp.Error = icmpPkg.conn.WriteTo(icmpPkg.netmsg, icmpPkg.dest)
	}
	if nil != hicmp.Error {
		return
	}
	buf := make([]byte, 1500)
	for {
		var readLen int
		readLen, hicmp.Addr, hicmp.Error = icmpPkg.conn.ReadFrom(buf)
		if nerr, ok := hicmp.Error.(net.Error); ok && nerr.Timeout() {
			hicmp.Timeout = true
			return
		}
		if nil != hicmp.Error {
			return
		}
		var result *icmp.Message
		if nil != icmpPkg.ipv4conn {
			result, hicmp.Error = icmp.ParseMessage(IPv4, buf[:readLen])
		}
		if nil != hicmp.Error {
			return
		}
		switch result.Type {
		case ipv4.ICMPTypeEchoReply:
			if rply, ok := result.Body.(*icmp.Echo); ok {
				if icmpPkg.id == rply.ID && icmpPkg.seq == rply.Seq {
					hicmp.Final = true
					hicmp.RTT = time.Since(sendOn)
					return
				}

			}
		case ipv4.ICMPTypeTimeExceeded:
			if rply, ok := result.Body.(*icmp.TimeExceeded); ok {
				if len(rply.Data) > 24 {
					if uint16(icmpPkg.id) == binary.BigEndian.Uint16(rply.Data[24:26]) {
						hicmp.RTT = time.Since(sendOn)
						return
					}
				}
			}
		case ipv4.ICMPTypeDestinationUnreachable:
			if rply, ok := result.Body.(*icmp.Echo); ok {
				if icmpPkg.id == rply.ID && icmpPkg.seq == rply.Seq {
					hicmp.Down = true
					hicmp.RTT = time.Since(sendOn)
					return
				}
			}
		}
	}
}

// Ping ...
func Ping(IPAddr string, maxrtt time.Duration) (rtt float64, err error) {

	var ip *net.IPAddr
	ip ,err = net.ResolveIPAddr("ip", IPAddr)
	if err != nil{
		err = fmt.Errorf("invaild ip")
		return
	}
	if os.Geteuid() != 0 {
		err = fmt.Errorf("skipping ping, root permissions missing")
		return
	}
	icmpPkg := new(IcmpPkg)
	icmpPkg.dest = ip
	icmpPkg.maxrtt = maxrtt
	icmpPkg.id = rand.Intn(65535)
	icmpPkg.seq = os.Getpid()&0xffff
	icmpPkg.msg = icmp.Message{
		Type: ipv4.ICMPTypeEcho, 
		Code: 0, 
		Body: &icmp.Echo{
			ID: icmpPkg.id, 
			Seq: icmpPkg.seq,
			Data: []byte("ping-ping-ping"),
		}}
	icmpPkg.netmsg, err = icmpPkg.msg.Marshal(nil)
	if err != nil {
		return
	}

	pingRsult := icmpPkg.Send(64)
	rtt = float64(pingRsult.RTT.Nanoseconds()) / 1e6
	err = pingRsult.Error
	return
}
