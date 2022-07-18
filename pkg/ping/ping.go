package ping

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

var (
	echoMsg = []byte("PING")

	BadReplyType  = fmt.Errorf("bad replay type")
	BadReplayBody = fmt.Errorf("bad replay body")
)

const (
	maxBufLen = 1500

	ipv4proto = 1
	ipv6proto = 58
)

func ping(ip, listen, network string, proto int, t icmp.Type, rt icmp.Type, timeout time.Duration) (time.Duration, error) {
	con, err := icmp.ListenPacket(network, listen)
	if err != nil {
		return 0, err
	}
	defer con.Close()

	err = con.SetDeadline(time.Now().Add(timeout))
	if err != nil {
		return 0, err
	}

	wm := icmp.Message{
		Type: t,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getegid() & 0xffff,
			Seq:  1,
			Data: echoMsg,
		},
	}
	b, err := wm.Marshal(nil)
	if err != nil {
		return 0, err
	}

	start := time.Now()
	_, err = con.WriteTo(b, &net.UDPAddr{
		IP:   net.ParseIP(ip),
		Zone: "en0",
	})
	if err != nil {
		return 0, err
	}

	rb := make([]byte, maxBufLen)
	n, _, err := con.ReadFrom(rb)
	if err != nil {
		return 0, err
	}
	reply, err := icmp.ParseMessage(proto, rb[:n])
	if err != nil {
		return 0, err
	}

	if reply.Type != rt {
		return 0, BadReplyType
	}

	rb, err = reply.Body.Marshal(proto)
	if err != nil {
		return 0, err
	}

	if !bytes.Equal(rb[len(echoMsg):], echoMsg) {
		return 0, BadReplayBody
	}

	return time.Since(start), nil
}

func Ping(ip string, timeout time.Duration) (time.Duration, error) {
	if strings.Contains(ip, "::") {
		return ping(ip, "::1", "udp6", ipv6proto, ipv6.ICMPTypeEchoRequest, ipv6.ICMPTypeEchoReply, timeout)
	}

	return ping(ip, "0.0.0.0", "udp4", ipv4proto, ipv4.ICMPTypeEcho, ipv4.ICMPTypeEchoReply, timeout)
}
