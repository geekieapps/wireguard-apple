package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net"
	"net/netip"
	"sync"
	"time"

	sing_tun "github.com/sagernet/sing-tun"
	"github.com/sagernet/sing/common/buf"
	"github.com/sagernet/sing/common/logger"
	M "github.com/sagernet/sing/common/metadata"
	N "github.com/sagernet/sing/common/network"
	"github.com/sagernet/sing/common/uot"
	"golang.org/x/net/proxy"
)

// MARK: - State

var (
	singMu    sync.Mutex
	singStack sing_tun.Stack
	singDev   sing_tun.Tun
)

// MARK: - Handler

type singTunHandler struct {
	proxyAddr string // 127.0.0.1:10808
}

// PrepareConnection is called by gVisor before NewConnectionEx/NewPacketConnectionEx.
// Returning nil allows all connections through.
func (h *singTunHandler) PrepareConnection(
	network string,
	source M.Socksaddr,
	destination M.Socksaddr,
	routeContext sing_tun.DirectRouteContext,
	timeout time.Duration,
) (sing_tun.DirectRouteDestination, error) {
	return nil, nil
}

// NewConnectionEx forwards TCP connections from gVisor → SOCKS5 → XRay.
func (h *singTunHandler) NewConnectionEx(
	ctx context.Context,
	conn net.Conn,
	source M.Socksaddr,
	destination M.Socksaddr,
	onClose N.CloseHandlerFunc,
) {
	defer func() {
		conn.Close()
		if onClose != nil {
			onClose(nil)
		}
	}()

	dialer, err := proxy.SOCKS5("tcp", h.proxyAddr, nil, proxy.Direct)
	if err != nil {
		return
	}
	proxyConn, err := dialer.Dial("tcp", destination.String())
	if err != nil {
		return
	}
	defer proxyConn.Close()
	relay(conn, proxyConn)
}

// NewPacketConnectionEx forwards UDP via UoT (UDP-over-TCP) → SOCKS5 → XRay.
// UoT avoids flaky SOCKS5 UDP ASSOCIATE by tunneling UDP inside a TCP connection.
func (h *singTunHandler) NewPacketConnectionEx(
	ctx context.Context,
	conn N.PacketConn,
	source M.Socksaddr,
	destination M.Socksaddr,
	onClose N.CloseHandlerFunc,
) {
	defer func() {
		conn.Close()
		if onClose != nil {
			onClose(nil)
		}
	}()

	// Connect to SOCKS5 proxy using the UoT magic address so XRay's inbound
	// recognises this as a UDP-over-TCP session.
	dialer, err := proxy.SOCKS5("tcp", h.proxyAddr, nil, proxy.Direct)
	if err != nil {
		return
	}
	tcpConn, err := dialer.Dial("tcp", uot.RequestDestination(uot.Version).String())
	if err != nil {
		return
	}
	defer tcpConn.Close()

	uotConn := uot.NewConn(tcpConn, uot.Request{
		IsConnect:   false,
		Destination: destination,
	})
	defer uotConn.Close()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			buffer := buf.New()
			addr, err := conn.ReadPacket(buffer)
			if err != nil {
				buffer.Release()
				return
			}
			if err := uotConn.WritePacket(buffer, addr); err != nil {
				buffer.Release()
				return
			}
			buffer.Release()
		}
	}()

	for {
		buffer := buf.New()
		addr, err := uotConn.ReadPacket(buffer)
		if err != nil {
			buffer.Release()
			return
		}
		if err := conn.WritePacket(buffer, addr); err != nil {
			buffer.Release()
			return
		}
		buffer.Release()
	}
}

// MARK: - C exports

// SingTunStart starts a gVisor TUN stack forwarding traffic to proxyAddr via SOCKS5.
// tunFd is the utun file descriptor obtained from the iOS Network Extension.
//
//export SingTunStart
func SingTunStart(tunFd C.int, mtu C.int, proxyAddr *C.char) *C.char {
	singMu.Lock()
	defer singMu.Unlock()

	if singStack != nil {
		return singOK()
	}

	handler := &singTunHandler{proxyAddr: C.GoString(proxyAddr)}

	tunOptions := sing_tun.Options{
		FileDescriptor: int(tunFd),
		MTU:            uint32(mtu),
		Inet4Address:   []netip.Prefix{netip.MustParsePrefix("198.18.0.1/16")},
	}

	dev, err := sing_tun.New(tunOptions)
	if err != nil {
		return singErr("create tun: " + err.Error())
	}

	stack, err := sing_tun.NewStack("gvisor", sing_tun.StackOptions{
		Context:    context.Background(),
		Tun:        dev,
		TunOptions: tunOptions,
		Handler:    handler,
		Logger:     logger.NOP(),
	})
	if err != nil {
		dev.Close()
		return singErr("create stack: " + err.Error())
	}

	if err := stack.Start(); err != nil {
		stack.Close()
		dev.Close()
		return singErr("start stack: " + err.Error())
	}

	singStack = stack
	singDev = dev
	return singOK()
}

// SingTunStop tears down the gVisor TUN stack.
//
//export SingTunStop
func SingTunStop() *C.char {
	singMu.Lock()
	defer singMu.Unlock()

	if singStack != nil {
		singStack.Close()
		singStack = nil
	}
	if singDev != nil {
		singDev.Close()
		singDev = nil
	}
	return singOK()
}

// MARK: - Helpers

func relay(a, b io.ReadWriteCloser) {
	done := make(chan struct{})
	go func() {
		io.Copy(a, b)
		close(done)
	}()
	io.Copy(b, a)
	<-done
}

type singResult struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

func singOK() *C.char {
	data, _ := json.Marshal(singResult{Success: true})
	return C.CString(base64.StdEncoding.EncodeToString(data))
}

func singErr(msg string) *C.char {
	data, _ := json.Marshal(singResult{Success: false, Error: msg})
	return C.CString(base64.StdEncoding.EncodeToString(data))
}
