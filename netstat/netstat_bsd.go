package netstat

import (
	"fmt"
	"log"
	"unsafe"
)

// Socket states
const (
	Established SkState = 0x01
	SynSent             = 0x02
	SynRecv             = 0x03
	FinWait1            = 0x04
	FinWait2            = 0x05
	TimeWait            = 0x06
	Close               = 0x07
	CloseWait           = 0x08
	LastAck             = 0x09
	Listen              = 0x0a
	Closing             = 0x0b
)

var skStates = [...]string{
	"UNKNOWN",
	"ESTABLISHED",
	"SYN_SENT",
	"SYN_RECV",
	"FIN_WAIT1",
	"FIN_WAIT2",
	"TIME_WAIT",
	"", // CLOSE
	"CLOSE_WAIT",
	"LAST_ACK",
	"LISTEN",
	"CLOSING",
}

func osTCPSocks(accept AcceptFn) ([]SockTabEntry, error) {
	var s = "net.inet.tcp.pcblist"
	var retry = 5
	var xig, exig *Xinpgen
	var buf []byte
	for {
		b, err := SysctlByName(s)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		xig = (*Xinpgen)(unsafe.Pointer(&b[0]))
		sxig := unsafe.Sizeof(*xig)
		eoff := uintptr(len(b)) - sxig
		exig = (*Xinpgen)(unsafe.Pointer(&b[eoff]))
		if uint64(xig.Len) != uint64(sxig) || uint64(exig.Len) != uint64(sxig) {
			log.Fatalf("xinpgen size mismatch, xig: %d, exig: %v", xig, exig)
		}
		fmt.Printf("xig: %v, buflen: %d, eoff: %d\n", xig, len(b), eoff)
		if !(xig.Gen != exig.Gen && retry > 0) {
			buf = b
			break
		}
		retry--
	}
	var index uint64
	index = uint64(index) + uint64(xig.Len)
	for {
		xig = (*Xinpgen)(unsafe.Pointer(&buf[index]))
		if uintptr(unsafe.Pointer(xig)) >= uintptr(unsafe.Pointer(exig)) {
			break
		}
		xtp := (*Xtcpcb)(unsafe.Pointer(xig))
		fmt.Printf("Proto: %d\n", xtp.Socket.Xso_protocol)
		index += uint64(xig.Len)
	}
	return nil, nil
}

func osTCP6Socks(accept AcceptFn) ([]SockTabEntry, error) {
	return nil, nil
}

func osUDPSocks(accept AcceptFn) ([]SockTabEntry, error) {
	return nil, nil
}

func osUDP6Socks(accept AcceptFn) ([]SockTabEntry, error) {
	return nil, nil
}
