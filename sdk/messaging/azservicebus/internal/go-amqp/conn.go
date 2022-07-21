// Copyright (C) 2017 Kale Blankenship
// Portions Copyright (c) Microsoft Corporation
package amqp

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"net"
	"net/url"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/go-amqp/internal/bitmap"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/go-amqp/internal/buffer"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/go-amqp/internal/encoding"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/go-amqp/internal/frames"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/go-amqp/internal/log"
)

// Default connection options
const (
	defaultIdleTimeout  = 1 * time.Minute
	defaultMaxFrameSize = 65536
	defaultMaxSessions  = 65536
)

// ConnOptions contains the optional settings for configuring an AMQP connection.
type ConnOptions struct {
	// ContainerID sets the container-id to use when opening the connection.
	//
	// A container ID will be randomly generated if this option is not used.
	ContainerID string

	// HostName sets the hostname sent in the AMQP
	// Open frame and TLS ServerName (if not otherwise set).
	HostName string

	// IdleTimeout specifies the maximum period in milliseconds between
	// receiving frames from the peer.
	//
	// Specify a value less than zero to disable idle timeout.
	//
	// Default: 1 minute.
	IdleTimeout time.Duration

	// MaxFrameSize sets the maximum frame size that
	// the connection will accept.
	//
	// Must be 512 or greater.
	//
	// Default: 512.
	MaxFrameSize uint32

	// MaxSessions sets the maximum number of channels.
	// The value must be greater than zero.
	//
	// Default: 65535.
	MaxSessions uint16

	// Properties sets an entry in the connection properties map sent to the server.
	Properties map[string]interface{}

	// SASLType contains the specified SASL authentication mechanism.
	SASLType SASLType

	// Timeout configures how long to wait for the
	// server during connection establishment.
	//
	// Once the connection has been established, IdleTimeout
	// applies. If duration is zero, no timeout will be applied.
	//
	// Default: 0.
	Timeout time.Duration

	// TLSConfig sets the tls.Config to be used during
	// TLS negotiation.
	//
	// This option is for advanced usage, in most scenarios
	// providing a URL scheme of "amqps://" is sufficient.
	TLSConfig *tls.Config

	// test hook
	dialer dialer
}

// used to abstract the underlying dialer for testing purposes
type dialer interface {
	NetDialerDial(c *conn, host, port string) error
	TLSDialWithDialer(c *conn, host, port string) error
}

// conn is an AMQP connection.
// only exported fields and methods are part of public surface area,
// all others are considered to be internal implementation details.
type conn struct {
	net            net.Conn      // underlying connection
	connectTimeout time.Duration // time to wait for reads/writes during conn establishment
	dialer         dialer        // used for testing purposes, it allows faking dialing TCP/TLS endpoints

	// TLS
	tlsNegotiation bool        // negotiate TLS
	tlsComplete    bool        // TLS negotiation complete
	tlsConfig      *tls.Config // TLS config, default used if nil (ServerName set to Client.hostname)

	// SASL
	saslHandlers map[encoding.Symbol]stateFunc // map of supported handlers keyed by SASL mechanism, SASL not negotiated if nil
	saslComplete bool                          // SASL negotiation complete; internal *except* for SASL auth methods

	// local settings
	maxFrameSize uint32                          // max frame size to accept
	channelMax   uint16                          // maximum number of channels to allow
	hostname     string                          // hostname of remote server (set explicitly or parsed from URL)
	idleTimeout  time.Duration                   // maximum period between receiving frames
	properties   map[encoding.Symbol]interface{} // additional properties sent upon connection open
	containerID  string                          // set explicitly or randomly generated

	// peer settings
	peerIdleTimeout  time.Duration // maximum period between sending frames
	PeerMaxFrameSize uint32        // maximum frame size peer will accept

	// conn state
	Done    chan struct{} // indicates the connection has terminated
	doneErr error         // contains the error state returned from Close() and Err()

	// connReader and connWriter management
	rxtxExit  chan struct{} // signals connReader and connWriter to exit
	closeOnce sync.Once     // ensures that close() is only called once

	// session tracking
	channels            *bitmap.Bitmap
	sessionsByChannel   map[uint16]*Session
	sessionsByChannelMu sync.RWMutex

	// connReader
	rxBuf         buffer.Buffer // incomes bytes buffer
	rxDone        chan struct{} // closed when connReader exits
	rxErr         chan error    // contains last error reading from c.net
	connReaderRun chan func()   // functions to be run by conn reader (set deadline on conn to run)

	// connWriter
	txFrame chan frames.Frame // AMQP frames to be sent by connWriter
	txBuf   buffer.Buffer     // buffer for marshaling frames before transmitting
	txDone  chan struct{}     // closed when connWriter exits
	txErr   chan error        // contains last error writing to c.net
}

// implements the dialer interface
type defaultDialer struct{}

func (defaultDialer) NetDialerDial(c *conn, host, port string) (err error) {
	dialer := &net.Dialer{Timeout: c.connectTimeout}
	c.net, err = dialer.Dial("tcp", net.JoinHostPort(host, port))
	return
}

func (defaultDialer) TLSDialWithDialer(c *conn, host, port string) (err error) {
	dialer := &net.Dialer{Timeout: c.connectTimeout}
	c.net, err = tls.DialWithDialer(dialer, "tcp", net.JoinHostPort(host, port), c.tlsConfig)
	return
}

func dialConn(addr string, opts *ConnOptions) (*conn, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}
	host, port := u.Hostname(), u.Port()
	if port == "" {
		port = "5672"
		if u.Scheme == "amqps" || u.Scheme == "amqp+ssl" {
			port = "5671"
		}
	}

	var cp ConnOptions
	if opts != nil {
		cp = *opts
	}

	// prepend SASL credentials when the user/pass segment is not empty
	if u.User != nil {
		pass, _ := u.User.Password()
		cp.SASLType = SASLTypePlain(u.User.Username(), pass)
	}

	if cp.HostName == "" {
		cp.HostName = host
	}

	c, err := newConn(nil, &cp)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case "amqp", "":
		err = c.dialer.NetDialerDial(c, host, port)
	case "amqps", "amqp+ssl":
		c.initTLSConfig()
		c.tlsNegotiation = false
		err = c.dialer.TLSDialWithDialer(c, host, port)
	default:
		err = fmt.Errorf("unsupported scheme %q", u.Scheme)
	}

	if err != nil {
		return nil, err
	}
	return c, nil
}

func newConn(netConn net.Conn, opts *ConnOptions) (*conn, error) {
	c := &conn{
		dialer:            defaultDialer{},
		net:               netConn,
		maxFrameSize:      defaultMaxFrameSize,
		PeerMaxFrameSize:  defaultMaxFrameSize,
		channelMax:        defaultMaxSessions - 1, // -1 because channel-max starts at zero
		idleTimeout:       defaultIdleTimeout,
		containerID:       randString(40),
		Done:              make(chan struct{}),
		rxtxExit:          make(chan struct{}),
		rxErr:             make(chan error, 1), // buffered to ensure connReader won't block
		txErr:             make(chan error, 1), // buffered to ensure connWriter won't block
		rxDone:            make(chan struct{}),
		connReaderRun:     make(chan func(), 1), // buffered to allow queueing function before interrupt
		txFrame:           make(chan frames.Frame),
		txDone:            make(chan struct{}),
		sessionsByChannel: map[uint16]*Session{},
	}

	// apply options
	if opts == nil {
		opts = &ConnOptions{}
	}

	if opts.ContainerID != "" {
		c.containerID = opts.ContainerID
	}
	if opts.HostName != "" {
		c.hostname = opts.HostName
	}
	if opts.IdleTimeout > 0 {
		c.idleTimeout = opts.IdleTimeout
	} else if opts.IdleTimeout < 0 {
		c.idleTimeout = 0
	}
	if opts.MaxFrameSize > 0 && opts.MaxFrameSize < 512 {
		return nil, fmt.Errorf("invalid MaxFrameSize value %d", opts.MaxFrameSize)
	} else if opts.MaxFrameSize > 512 {
		c.maxFrameSize = opts.MaxFrameSize
	}
	if opts.MaxSessions > 0 {
		c.channelMax = opts.MaxSessions
	}
	if opts.SASLType != nil {
		if err := opts.SASLType(c); err != nil {
			return nil, err
		}
	}
	if opts.Timeout > 0 {
		c.connectTimeout = opts.Timeout
	}
	if opts.Properties != nil {
		c.properties = make(map[encoding.Symbol]interface{})
		for key, val := range opts.Properties {
			c.properties[encoding.Symbol(key)] = val
		}
	}
	if opts.TLSConfig != nil {
		c.tlsConfig = opts.TLSConfig.Clone()
	}
	if opts.dialer != nil {
		c.dialer = opts.dialer
	}
	return c, nil
}

func (c *conn) initTLSConfig() {
	// create a new config if not already set
	if c.tlsConfig == nil {
		c.tlsConfig = new(tls.Config)
	}

	// TLS config must have ServerName or InsecureSkipVerify set
	if c.tlsConfig.ServerName == "" && !c.tlsConfig.InsecureSkipVerify {
		c.tlsConfig.ServerName = c.hostname
	}
}

// Start establishes the connection and begins multiplexing network IO.
// It is an error to call Start() on a connection that's been closed.
func (c *conn) Start() error {
	// run connection establishment state machine
	for state := c.negotiateProto; state != nil; {
		var err error
		state, err = state()
		// check if err occurred
		if err != nil {
			close(c.txDone) // close here since connWriter hasn't been started yet
			close(c.rxDone)
			_ = c.Close()
			return err
		}
	}

	// we can't create the channel bitmap until the connection has been established.
	// this is because our peer can tell us the max channels they support.
	c.channels = bitmap.New(uint32(c.channelMax))

	go c.connWriter()
	go c.connReader()

	return nil
}

// Close closes the connection.
func (c *conn) Close() error {
	c.closeOnce.Do(func() { c.close() })

	var connErr *ConnectionError
	if errors.As(c.doneErr, &connErr) && connErr.inner == nil {
		// an empty ConnectionError means the connection was closed by the caller
		return nil
	}

	// there was an error during shut-down or connReader/connWriter
	// experienced a terminal error
	return c.doneErr
}

// close is called once, either from Close() or when connReader/connWriter exits
func (c *conn) close() {
	defer close(c.Done)

	close(c.rxtxExit)

	// wait for writing to stop, allows it to send the final close frame
	<-c.txDone

	// drain pending TX error
	var txErr error
	select {
	case txErr = <-c.txErr:
		// there was an error in connWriter
	default:
		// no pending write error
	}

	closeErr := c.net.Close()

	// check rxDone after closing net, otherwise may block
	// for up to c.idleTimeout
	<-c.rxDone

	// drain pending RX error
	var rxErr error
	select {
	case rxErr = <-c.rxErr:
		// there was an error in connReader
		if errors.Is(rxErr, net.ErrClosed) {
			// this is the expected error when the connection is closed, swallow it
			rxErr = nil
		}
	default:
		// no pending read error
	}

	if txErr == nil && rxErr == nil && closeErr == nil {
		// if there are no errors, it means user initiated close() and we shut down cleanly
		c.doneErr = &ConnectionError{}
	} else if amqpErr, ok := rxErr.(*Error); ok {
		// we experienced a peer-initiated close that contained an Error.  return it
		c.doneErr = &ConnectionError{inner: amqpErr}
	} else if txErr != nil {
		c.doneErr = &ConnectionError{inner: txErr}
	} else if rxErr != nil {
		c.doneErr = &ConnectionError{inner: rxErr}
	} else {
		c.doneErr = &ConnectionError{inner: closeErr}
	}
}

// Err returns the connection's error state after it's been closed.
// Calling this on an open connection will block until the connection is closed.
func (c *conn) Err() error {
	<-c.Done
	return c.doneErr
}

func (c *conn) NewSession() (*Session, error) {
	c.sessionsByChannelMu.Lock()
	defer c.sessionsByChannelMu.Unlock()

	// create the next session to allocate
	// note that channel always start at 0
	channel, ok := c.channels.Next()
	if !ok {
		return nil, fmt.Errorf("reached connection channel max (%d)", c.channelMax)
	}
	session := newSession(c, uint16(channel))
	ch := session.channel
	c.sessionsByChannel[ch] = session
	return session, nil
}

func (c *conn) DeleteSession(s *Session) {
	c.sessionsByChannelMu.Lock()
	defer c.sessionsByChannelMu.Unlock()

	delete(c.sessionsByChannel, s.channel)
	c.channels.Remove(uint32(s.channel))
}

// connReader reads from the net.Conn, decodes frames, and either handles
// them here as appropriate or sends them to the session.rx channel.
func (c *conn) connReader() {
	defer func() {
		close(c.rxDone)
		c.closeOnce.Do(func() { c.close() })
	}()

	var sessionsByRemoteChannel = make(map[uint16]*Session)
	var err error
	for {
		if err != nil {
			log.Debug(1, "connReader terminal error: %v", err)
			c.rxErr <- err
			return
		}

		if c.idleTimeout > 0 {
			_ = c.net.SetReadDeadline(time.Now().Add(c.idleTimeout))
		}

		var fr frames.Frame
		fr, err = c.readFrame()
		if err != nil {
			continue
		}

		var (
			session *Session
			ok      bool
		)

		switch body := fr.Body.(type) {
		// Server initiated close.
		case *frames.PerformClose:
			// connWriter will send the close performative ack on its way out.
			// it's a SHOULD though, not a MUST.
			log.Debug(3, "RX (connReader): %s", body)
			if body.Error == nil {
				return
			}
			err = body.Error
			continue

		// RemoteChannel should be used when frame is Begin
		case *frames.PerformBegin:
			if body.RemoteChannel == nil {
				// since we only support remotely-initiated sessions, this is an error
				// TODO: it would be ideal to not have this kill the connection
				err = fmt.Errorf("%T: nil RemoteChannel", fr.Body)
				continue
			}
			c.sessionsByChannelMu.RLock()
			session, ok = c.sessionsByChannel[*body.RemoteChannel]
			c.sessionsByChannelMu.RUnlock()
			if !ok {
				err = fmt.Errorf("unexpected remote channel number %d", *body.RemoteChannel)
				continue
			}

			session.remoteChannel = fr.Channel
			sessionsByRemoteChannel[fr.Channel] = session

		case *frames.PerformEnd:
			session, ok = sessionsByRemoteChannel[fr.Channel]
			if !ok {
				err = fmt.Errorf("%T: didn't find channel %d in sessionsByRemoteChannel (PerformEnd)", fr.Body, fr.Channel)
				continue
			}
			// we MUST remove the remote channel from our map as soon as we receive
			// the ack (i.e. before passing it on to the session mux) on the session
			// ending since the numbers are recycled.
			delete(sessionsByRemoteChannel, fr.Channel)

		default:
			// pass on performative to the correct session
			session, ok = sessionsByRemoteChannel[fr.Channel]
			if !ok {
				err = fmt.Errorf("%T: didn't find channel %d in sessionsByRemoteChannel", fr.Body, fr.Channel)
				continue
			}
		}

		select {
		case session.rx <- fr:
		case <-c.rxtxExit:
			return
		}
	}
}

// readFrame reads a complete frame from c.net.
// it assumes that any read deadline has already been applied.
// used externally by SASL only.
func (c *conn) readFrame() (frames.Frame, error) {
	switch {
	// Cheaply reuse free buffer space when fully read.
	case c.rxBuf.Len() == 0:
		c.rxBuf.Reset()

	// Prevent excessive/unbounded growth by shifting data to beginning of buffer.
	case int64(c.rxBuf.Size()) > int64(c.maxFrameSize):
		c.rxBuf.Reclaim()
	}

	var (
		currentHeader   frames.Header // keep track of the current header, for frames split across multiple TCP packets
		frameInProgress bool          // true if in the middle of receiving data for currentHeader
	)

	for {
		// need to read more if buf doesn't contain the complete frame
		// or there's not enough in buf to parse the header
		if frameInProgress || c.rxBuf.Len() < frames.HeaderSize {
			err := c.rxBuf.ReadFromOnce(c.net)
			if err != nil {
				log.Debug(1, "readFrame error: %v", err)
				select {
				// if there is a pending connReaderRun function, execute it
				case f := <-c.connReaderRun:
					f()
					continue

				// return error to caller
				default:
					return frames.Frame{}, err
				}
			}
		}

		// read more if buf doesn't contain enough to parse the header
		if c.rxBuf.Len() < frames.HeaderSize {
			continue
		}

		// parse the header if a frame isn't in progress
		if !frameInProgress {
			var err error
			currentHeader, err = frames.ParseHeader(&c.rxBuf)
			if err != nil {
				return frames.Frame{}, err
			}
			frameInProgress = true
		}

		// check size is reasonable
		if currentHeader.Size > math.MaxInt32 { // make max size configurable
			return frames.Frame{}, errors.New("payload too large")
		}

		bodySize := int64(currentHeader.Size - frames.HeaderSize)

		// the full frame hasn't been received, keep reading
		if int64(c.rxBuf.Len()) < bodySize {
			continue
		}
		frameInProgress = false

		// check if body is empty (keepalive)
		if bodySize == 0 {
			continue
		}

		// parse the frame
		b, ok := c.rxBuf.Next(bodySize)
		if !ok {
			return frames.Frame{}, fmt.Errorf("buffer EOF; requested bytes: %d, actual size: %d", bodySize, c.rxBuf.Len())
		}

		parsedBody, err := frames.ParseBody(buffer.New(b))
		if err != nil {
			return frames.Frame{}, err
		}

		return frames.Frame{Channel: currentHeader.Channel, Body: parsedBody}, nil
	}
}

func (c *conn) connWriter() {
	defer func() {
		close(c.txDone)
		c.closeOnce.Do(func() { c.close() })
	}()

	// disable write timeout
	if c.connectTimeout != 0 {
		c.connectTimeout = 0
		_ = c.net.SetWriteDeadline(time.Time{})
	}

	var (
		// keepalives are sent at a rate of 1/2 idle timeout
		keepaliveInterval = c.peerIdleTimeout / 2
		// 0 disables keepalives
		keepalivesEnabled = keepaliveInterval > 0
		// set if enable, nil if not; nil channels block forever
		keepalive <-chan time.Time
	)

	if keepalivesEnabled {
		ticker := time.NewTicker(keepaliveInterval)
		defer ticker.Stop()
		keepalive = ticker.C
	}

	var err error
	for {
		if err != nil {
			log.Debug(1, "connWriter terminal error: %v", err)
			c.txErr <- err
			return
		}

		select {
		// frame write request
		case fr := <-c.txFrame:
			err = c.writeFrame(fr)
			if err == nil && fr.Done != nil {
				close(fr.Done)
			}

		// keepalive timer
		case <-keepalive:
			log.Debug(3, "sending keep-alive frame")
			_, err = c.net.Write(keepaliveFrame)
			// It would be slightly more efficient in terms of network
			// resources to reset the timer each time a frame is sent.
			// However, keepalives are small (8 bytes) and the interval
			// is usually on the order of minutes. It does not seem
			// worth it to add extra operations in the write path to
			// avoid. (To properly reset a timer it needs to be stopped,
			// possibly drained, then reset.)

		// connection complete
		case <-c.rxtxExit:
			// send close performative.  note that the spec says we
			// SHOULD wait for the ack but we don't HAVE to, in order
			// to be resilient to bad actors etc.  so we just send
			// the close performative and exit.
			cls := &frames.PerformClose{}
			log.Debug(1, "TX (connWriter): %s", cls)
			c.txErr <- c.writeFrame(frames.Frame{
				Type: frameTypeAMQP,
				Body: cls,
			})
			return
		}
	}
}

// writeFrame writes a frame to the network.
// used externally by SASL only.
func (c *conn) writeFrame(fr frames.Frame) error {
	if c.connectTimeout != 0 {
		_ = c.net.SetWriteDeadline(time.Now().Add(c.connectTimeout))
	}

	// writeFrame into txBuf
	c.txBuf.Reset()
	err := writeFrame(&c.txBuf, fr)
	if err != nil {
		return err
	}

	// validate the frame isn't exceeding peer's max frame size
	requiredFrameSize := c.txBuf.Len()
	if uint64(requiredFrameSize) > uint64(c.PeerMaxFrameSize) {
		return fmt.Errorf("%T frame size %d larger than peer's max frame size %d", fr, requiredFrameSize, c.PeerMaxFrameSize)
	}

	// write to network
	_, err = c.net.Write(c.txBuf.Bytes())
	return err
}

// writeProtoHeader writes an AMQP protocol header to the
// network
func (c *conn) writeProtoHeader(pID protoID) error {
	if c.connectTimeout != 0 {
		_ = c.net.SetWriteDeadline(time.Now().Add(c.connectTimeout))
	}
	_, err := c.net.Write([]byte{'A', 'M', 'Q', 'P', byte(pID), 1, 0, 0})
	return err
}

// keepaliveFrame is an AMQP frame with no body, used for keepalives
var keepaliveFrame = []byte{0x00, 0x00, 0x00, 0x08, 0x02, 0x00, 0x00, 0x00}

// SendFrame is used by sessions and links to send frames across the network.
func (c *conn) SendFrame(fr frames.Frame) error {
	select {
	case c.txFrame <- fr:
		return nil
	case <-c.Done:
		return c.Err()
	}
}

// stateFunc is a state in a state machine.
//
// The state is advanced by returning the next state.
// The state machine concludes when nil is returned.
type stateFunc func() (stateFunc, error)

// negotiateProto determines which proto to negotiate next.
// used externally by SASL only.
func (c *conn) negotiateProto() (stateFunc, error) {
	// in the order each must be negotiated
	switch {
	case c.tlsNegotiation && !c.tlsComplete:
		return c.exchangeProtoHeader(protoTLS)
	case c.saslHandlers != nil && !c.saslComplete:
		return c.exchangeProtoHeader(protoSASL)
	default:
		return c.exchangeProtoHeader(protoAMQP)
	}
}

type protoID uint8

// protocol IDs received in protoHeaders
const (
	protoAMQP protoID = 0x0
	protoTLS  protoID = 0x2
	protoSASL protoID = 0x3
)

// exchangeProtoHeader performs the round trip exchange of protocol
// headers, validation, and returns the protoID specific next state.
func (c *conn) exchangeProtoHeader(pID protoID) (stateFunc, error) {
	// write the proto header
	if err := c.writeProtoHeader(pID); err != nil {
		return nil, err
	}

	// read response header
	p, err := c.readProtoHeader()
	if err != nil {
		return nil, err
	}

	if pID != p.ProtoID {
		return nil, fmt.Errorf("unexpected protocol header %#00x, expected %#00x", p.ProtoID, pID)
	}

	// go to the proto specific state
	switch pID {
	case protoAMQP:
		return c.openAMQP, nil
	case protoTLS:
		return c.startTLS, nil
	case protoSASL:
		return c.negotiateSASL, nil
	default:
		return nil, fmt.Errorf("unknown protocol ID %#02x", p.ProtoID)
	}
}

// readProtoHeader reads a protocol header packet from c.rxProto.
func (c *conn) readProtoHeader() (protoHeader, error) {
	// only read from the network once our buffer has been exhausted.
	// TODO: this preserves existing behavior as some tests rely on this
	// implementation detail (it lets you replay a stream of bytes). we
	// might want to consider removing this and fixing the tests as the
	// protocol doesn't actually work this way.
	if c.rxBuf.Len() == 0 {
		for {
			if c.connectTimeout != 0 {
				_ = c.net.SetReadDeadline(time.Now().Add(c.connectTimeout))
			}

			err := c.rxBuf.ReadFromOnce(c.net)
			if err != nil {
				return protoHeader{}, err
			}

			// read more if buf doesn't contain enough to parse the header
			if c.rxBuf.Len() >= protoHeaderSize {
				break
			}
		}

		// reset outside the loop
		if c.connectTimeout != 0 {
			_ = c.net.SetReadDeadline(time.Time{})
		}
	}

	p, err := parseProtoHeader(&c.rxBuf)
	if err != nil {
		return protoHeader{}, err
	}
	return p, nil
}

// startTLS wraps the conn with TLS and returns to Client.negotiateProto
func (c *conn) startTLS() (stateFunc, error) {
	c.initTLSConfig()

	// buffered so connReaderRun won't block
	done := make(chan error, 1)

	// this function will be executed by connReader
	c.connReaderRun <- func() {
		defer close(done)
		_ = c.net.SetReadDeadline(time.Time{}) // clear timeout

		// wrap existing net.Conn and perform TLS handshake
		tlsConn := tls.Client(c.net, c.tlsConfig)
		if c.connectTimeout != 0 {
			_ = tlsConn.SetWriteDeadline(time.Now().Add(c.connectTimeout))
		}
		done <- tlsConn.Handshake()
		// TODO: return?

		// swap net.Conn
		c.net = tlsConn
		c.tlsComplete = true
	}

	// set deadline to interrupt connReader
	_ = c.net.SetReadDeadline(time.Time{}.Add(1))

	if err := <-done; err != nil {
		return nil, err
	}

	// go to next protocol
	return c.negotiateProto, nil
}

// openAMQP round trips the AMQP open performative
func (c *conn) openAMQP() (stateFunc, error) {
	// send open frame
	open := &frames.PerformOpen{
		ContainerID:  c.containerID,
		Hostname:     c.hostname,
		MaxFrameSize: c.maxFrameSize,
		ChannelMax:   c.channelMax,
		IdleTimeout:  c.idleTimeout / 2, // per spec, advertise half our idle timeout
		Properties:   c.properties,
	}
	log.Debug(1, "TX (openAMQP): %s", open)
	err := c.writeFrame(frames.Frame{
		Type:    frameTypeAMQP,
		Body:    open,
		Channel: 0,
	})
	if err != nil {
		return nil, err
	}

	// get the response
	fr, err := c.readSingleFrame()
	if err != nil {
		return nil, err
	}
	o, ok := fr.Body.(*frames.PerformOpen)
	if !ok {
		return nil, fmt.Errorf("openAMQP: unexpected frame type %T", fr.Body)
	}
	log.Debug(1, "RX (openAMQP): %s", o)

	// update peer settings
	if o.MaxFrameSize > 0 {
		c.PeerMaxFrameSize = o.MaxFrameSize
	}
	if o.IdleTimeout > 0 {
		// TODO: reject very small idle timeouts
		c.peerIdleTimeout = o.IdleTimeout
	}
	if o.ChannelMax < c.channelMax {
		c.channelMax = o.ChannelMax
	}

	// connection established, exit state machine
	return nil, nil
}

// negotiateSASL returns the SASL handler for the first matched
// mechanism specified by the server
func (c *conn) negotiateSASL() (stateFunc, error) {
	// read mechanisms frame
	fr, err := c.readSingleFrame()
	if err != nil {
		return nil, err
	}
	sm, ok := fr.Body.(*frames.SASLMechanisms)
	if !ok {
		return nil, fmt.Errorf("negotiateSASL: unexpected frame type %T", fr.Body)
	}
	log.Debug(1, "RX (negotiateSASL): %s", sm)

	// return first match in c.saslHandlers based on order received
	for _, mech := range sm.Mechanisms {
		if state, ok := c.saslHandlers[mech]; ok {
			return state, nil
		}
	}

	// no match
	return nil, fmt.Errorf("no supported auth mechanism (%v)", sm.Mechanisms) // TODO: send "auth not supported" frame?
}

// saslOutcome processes the SASL outcome frame and return Client.negotiateProto
// on success.
//
// SASL handlers return this stateFunc when the mechanism specific negotiation
// has completed.
// used externally by SASL only.
func (c *conn) saslOutcome() (stateFunc, error) {
	// read outcome frame
	fr, err := c.readSingleFrame()
	if err != nil {
		return nil, err
	}
	so, ok := fr.Body.(*frames.SASLOutcome)
	if !ok {
		return nil, fmt.Errorf("saslOutcome: unexpected frame type %T", fr.Body)
	}
	log.Debug(1, "RX (saslOutcome): %s", so)

	// check if auth succeeded
	if so.Code != encoding.CodeSASLOK {
		return nil, fmt.Errorf("SASL PLAIN auth failed with code %#00x: %s", so.Code, so.AdditionalData) // implement Stringer for so.Code
	}

	// return to c.negotiateProto
	c.saslComplete = true
	return c.negotiateProto, nil
}

// readSingleFrame is used during connection establishment to read a single frame.
//
// After setup, conn.connReader handles incoming frames.
func (c *conn) readSingleFrame() (frames.Frame, error) {
	if c.connectTimeout != 0 {
		_ = c.net.SetDeadline(time.Now().Add(c.connectTimeout))
		defer func() { _ = c.net.SetDeadline(time.Time{}) }()
	}

	fr, err := c.readFrame()
	if err != nil {
		return frames.Frame{}, err
	}

	return fr, nil
}

const protoHeaderSize = 8

type protoHeader struct {
	ProtoID  protoID
	Major    uint8
	Minor    uint8
	Revision uint8
}

// parseProtoHeader reads the proto header from r and returns the results
//
// An error is returned if the protocol is not "AMQP" or if the version is not 1.0.0.
func parseProtoHeader(r *buffer.Buffer) (protoHeader, error) {
	buf, ok := r.Next(protoHeaderSize)
	if !ok {
		return protoHeader{}, errors.New("invalid protoHeader")
	}
	_ = buf[7]

	if !bytes.Equal(buf[:4], []byte{'A', 'M', 'Q', 'P'}) {
		return protoHeader{}, fmt.Errorf("unexpected protocol %q", buf[:4])
	}

	p := protoHeader{
		ProtoID:  protoID(buf[4]),
		Major:    buf[5],
		Minor:    buf[6],
		Revision: buf[7],
	}

	if p.Major != 1 || p.Minor != 0 || p.Revision != 0 {
		return p, fmt.Errorf("unexpected protocol version %d.%d.%d", p.Major, p.Minor, p.Revision)
	}
	return p, nil
}

// writesFrame encodes fr into buf.
// split out from conn.WriteFrame for testing purposes.
func writeFrame(buf *buffer.Buffer, fr frames.Frame) error {
	// write header
	buf.Append([]byte{
		0, 0, 0, 0, // size, overwrite later
		2,       // doff, see frameHeader.DataOffset comment
		fr.Type, // frame type
	})
	buf.AppendUint16(fr.Channel) // channel

	// write AMQP frame body
	err := encoding.Marshal(buf, fr.Body)
	if err != nil {
		return err
	}

	// validate size
	if uint(buf.Len()) > math.MaxUint32 {
		return errors.New("frame too large")
	}

	// retrieve raw bytes
	bufBytes := buf.Bytes()

	// write correct size
	binary.BigEndian.PutUint32(bufBytes, uint32(len(bufBytes)))
	return nil
}
