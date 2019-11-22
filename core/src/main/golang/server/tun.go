package server

import (
	"encoding/binary"
	"net"

	"github.com/Dreamacro/clash/log"
	"github.com/kr328/clash/tun"
	"golang.org/x/sys/unix"
)

const (
	tunCommandEnd = 0x243
)

func handleTunStart(client *net.UnixConn) {
	buffer := make([]byte, unix.CmsgLen(4*1))

	_, noob, _, _, err := client.ReadMsgUnix(nil, buffer)
	if err != nil {
		log.Warnln("Read tun socket failure, %s", err.Error())
		return
	}

	msg, err := unix.ParseSocketControlMessage(buffer[:noob])
	if err != nil || len(msg) != 1 {
		log.Warnln("Parse tun socket failure, %s", err.Error())
		return
	}

	fds, err := unix.ParseUnixRights(&msg[0])
	if err != nil {
		log.Warnln("Parse tun socket failure, %s", err.Error())
		return
	}

	var mtu uint32
	var end uint32

	binary.Read(client, binary.BigEndian, &mtu)
	binary.Read(client, binary.BigEndian, &end)

	if end != tunCommandEnd {
		log.Warnln("Invalid tun command end")
		return
	}

	tun.StartTunProxy(fds[0], int(mtu))
}

func handleTunStop(client *net.UnixConn) {
	buf, _ := readCommandPacket(client)

	if tunCommandEnd != binary.BigEndian.Uint32(buf) {
		log.Warnln("Invalid tun command end")
	}

	tun.StopTunProxy()
}