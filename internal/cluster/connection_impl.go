// Copyright (c) 2008-2020, Hazelcast, Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License")
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cluster

import (
	"bytes"
	"fmt"
	"github.com/hazelcast/hazelcast-go-client/v4/internal/event"
	"net"
	"sync/atomic"
	"time"

	"crypto/tls"

	"github.com/hazelcast/hazelcast-go-client/v4/internal/config"
	"github.com/hazelcast/hazelcast-go-client/v4/internal/core"
	"github.com/hazelcast/hazelcast-go-client/v4/internal/core/logger"
	"github.com/hazelcast/hazelcast-go-client/v4/internal/proto"
	"github.com/hazelcast/hazelcast-go-client/v4/internal/util/timeutil"
	"github.com/hazelcast/hazelcast-go-client/v4/internal/util/versionutil"
)

const (
	kb              = 1024
	bufferSize      = 128 * kb
	protocolStarter = "CP2"
	IntMask         = 0xffff
)

type ResponseHandler func(msg *proto.ClientMessage)

type Connection interface {
}

type ConnectionImpl struct {
	responseCh                chan<- *proto.ClientMessage
	pending                   chan *proto.ClientMessage
	received                  chan *proto.ClientMessage
	socket                    net.Conn
	closed                    chan struct{}
	endpoint                  atomic.Value
	status                    int32
	isOwnerConnection         bool
	lastRead                  atomic.Value
	lastWrite                 atomic.Value
	closedTime                atomic.Value
	readBuffer                []byte
	clientMessageReader       *clientMessageReader
	connectionID              int64
	eventDispatcher           event.DispatchService
	connectedServerVersion    int32
	connectedServerVersionStr string
	startTime                 int64
	logger                    logger.Logger
}

func (c *ConnectionImpl) sendProtocolStarter() error {
	_, err := c.socket.Write([]byte(protocolStarter))
	return err
}

func (c *ConnectionImpl) createSocket(networkCfg NetworkConfig, address *core.Address) (net.Conn, error) {
	conTimeout := timeutil.GetPositiveDurationOrMax(networkCfg.ConnectionTimeout())
	socket, err := c.dialToAddressWithTimeout(address, conTimeout)
	if err != nil {
		return nil, err
	}
	//if networkCfg.SSLConfig().Enabled() {
	//	socket, err = c.openTLSConnection(networkCfg.SSLConfig(), socket)
	//}
	return socket, err
}

func (c *ConnectionImpl) dialToAddressWithTimeout(address *core.Address, conTimeout time.Duration) (net.Conn, error) {
	return net.DialTimeout("tcp", address.String(), conTimeout)
}

func (c *ConnectionImpl) init() {
	c.lastWrite.Store(time.Time{})
	c.closedTime.Store(time.Time{})
	c.startTime = timeutil.GetCurrentTimeInMilliSeconds()
	c.lastRead.Store(time.Now())
}

func (c *ConnectionImpl) openTLSConnection(sslCfg *config.SSLConfig, conn net.Conn) (net.Conn, error) {
	tlsCon := tls.Client(conn, sslCfg.Config)
	err := tlsCon.Handshake()
	return tlsCon, err
}

func (c *ConnectionImpl) isAlive() bool {
	return atomic.LoadInt32(&c.status) == 0
}

func (c *ConnectionImpl) writePool() {
	//Writer process
	for {
		select {
		case request, ok := <-c.pending:
			if !ok {
				continue
			}
			if err := c.write(request); err != nil {
				// XXX: create a new client message?
				request.Err = err
				c.responseCh <- request
			} else {
				c.lastWrite.Store(time.Now())
			}
		case <-c.closed:
			return
		}
	}
}

func (c *ConnectionImpl) send(inv ConnectionBoundInvocation) bool {
	select {
	case <-c.closed:
		return false
	case c.pending <- inv.Request():
		inv.StoreSentConnection(c)
		return true
	}
}

func (c *ConnectionImpl) write(clientMessage *proto.ClientMessage) error {
	buf := make([]byte, clientMessage.GetTotalLength())
	clientMessage.GetBytes(buf)
	_, err := c.socket.Write(buf)
	return err
}

func (c *ConnectionImpl) readPool() {
	for {
		c.socket.SetDeadline(time.Now().Add(2 * time.Second))
		buf := make([]byte, 4096)
		n, err := c.socket.Read(buf)
		if !c.isAlive() {
			return
		}
		if err != nil {
			if c.isTimeoutError(err) {
				continue
			}
			if err.Error() == "EOF" {
				continue
			}
			c.close(err)
			return
		}
		if n == 0 {
			continue
		}

		c.clientMessageReader.Append(bytes.NewBuffer(buf[:n]))
		c.receiveMessage()
		c.clientMessageReader.Reset()
	}
}

func (c *ConnectionImpl) isTimeoutError(err error) bool {
	netErr, ok := err.(net.Error)
	return ok && netErr.Timeout()
}

func (c *ConnectionImpl) StartTime() int64 {
	return c.startTime
}

func (c *ConnectionImpl) receiveMessage() {
	clientMessage := c.clientMessageReader.Read()
	if clientMessage != nil && clientMessage.StartFrame.HasUnFragmentedMessageFlags() {
		c.responseCh <- clientMessage
	}
}

func (c *ConnectionImpl) localAddress() net.Addr {
	return c.socket.LocalAddr()
}

func (c *ConnectionImpl) setConnectedServerVersion(connectedServerVersion string) {
	c.connectedServerVersionStr = connectedServerVersion
	c.connectedServerVersion = versionutil.CalculateVersion(connectedServerVersion)
}

func (c *ConnectionImpl) close(closeErr error) {
	if !atomic.CompareAndSwapInt32(&c.status, 0, 1) {
		return
	}
	//c.logger.Warn("ConnectionImpl :", c, " closed, err: ", closeErr)
	close(c.closed)
	c.socket.Close()
	c.closedTime.Store(time.Now())
	c.eventDispatcher.Publish(NewConnectionClosed(c, closeErr))
}

func (c *ConnectionImpl) String() string {
	return fmt.Sprintf("ClientConnection{isAlive=%t, connectionID=%d, endpoint=%s, lastReadTime=%s, lastWriteTime=%s, closedTime=%s, connected server version=%s",
		c.isAlive(), c.connectionID, c.endpoint.Load(), c.lastRead.Load(), c.lastWrite.Load(), c.closedTime.Load(), c.connectedServerVersionStr)
}
