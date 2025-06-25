package systemd

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

var mu sync.Mutex
var cond = sync.NewCond(&mu)

// StdinStdoutListener implements the net.Listener interface
type StdinStdoutListener struct {
	conn *stdinStdoutConn
}

func (l *StdinStdoutListener) Accept() (net.Conn, error) {
	// Accept always returns a new connection from stdin/stdout
	if l.conn == nil {
		cond.L.Lock()
		l.conn = &stdinStdoutConn{}

		return l.conn, nil
	}
	cond.L.Lock()
	return nil, io.EOF
}

func (l *StdinStdoutListener) Close() error {
	// Closing stdin/stdout is not necessary for this use case
	cond.Signal()
	cond.L.Unlock()
	return nil
}

func (l *StdinStdoutListener) Addr() net.Addr {
	return dummyAddr{}
}

// stdinStdoutConn implements the net.Conn interface
type stdinStdoutConn struct{}

func (c *stdinStdoutConn) Read(b []byte) (n int, err error) {
	// Read from stdin
	nr, err := os.Stdin.Read(b)

	if err != nil {
		cond.Signal()
		cond.L.Unlock()
	}

	return nr, err
}

func (c *stdinStdoutConn) Write(b []byte) (n int, err error) {
	// Write to stdout
	nw, err := os.Stdout.Write(b)

	if err != nil {
		cond.Signal()
		cond.L.Unlock()
	}

	return nw, err
}

func (c *stdinStdoutConn) Close() error {
	// Closing stdin/stdout is not necessary for this use case
	cond.Signal()
	cond.L.Unlock()
	return nil
}

func (c *stdinStdoutConn) LocalAddr() net.Addr {
	localAddr, err := resolveLocalAddr()
	if err != nil {
		// Handle the case where LOCAL_ADDR and LOCAL_PORT are not set or invalid.
		log.Printf("Warning: could not resolve local address: %v\n", err) //
		// Provide a default address or handle the error gracefully.
		localAddr = &net.UnixAddr{Net: "unix", Name: "local-unknown"}
	}
	return localAddr
}

func (c *stdinStdoutConn) RemoteAddr() net.Addr {
	// Provide a dummy address
	remoteAddr, err := resolveRemoteAddr()
	if err != nil {
		// Handle the case where REMOTE_ADDR and REMOTE_PORT are not set or invalid.
		// You might want to log a warning or return a different error.
		log.Printf("Warning: could not resolve remote address: %v\n", err) //
		// Provide a default address or handle the error gracefully.
		remoteAddr = &net.UnixAddr{Net: "unix", Name: "remote-unknown"}
	}
	return remoteAddr
}

func (c *stdinStdoutConn) SetDeadline(t time.Time) error {
	return nil // Not applicable for stdin/stdout
}

func (c *stdinStdoutConn) SetReadDeadline(t time.Time) error {
	return nil // Not applicable for stdin/stdout
}

func (c *stdinStdoutConn) SetWriteDeadline(t time.Time) error {
	return nil // Not applicable for stdin/stdout
}

// dummyAddr implements the net.Addr interface
type dummyAddr struct{}

func (a dummyAddr) Network() string {
	return "stdin/stdout"
}

func (a dummyAddr) String() string {
	return "stdin/stdout"
}

// resolveRemoteAddr reads REMOTE_ADDR and REMOTE_PORT environment variables
// and constructs a net.Addr.
func resolveRemoteAddr() (net.Addr, error) {
	remoteAddrStr := os.Getenv("REMOTE_ADDR")
	remotePortStr := os.Getenv("REMOTE_PORT")

	if remoteAddrStr == "" || remotePortStr == "" {
		return nil, fmt.Errorf("REMOTE_ADDR or REMOTE_PORT not set") //
	}

	address := fmt.Sprintf("%s:%s", remoteAddrStr, remotePortStr)
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve TCP address for remote: %w", err)
	}
	return addr, nil
}

// resolveLocalAddr reads LOCAL_ADDR and LOCAL_PORT environment variables
// and constructs a net.Addr.
func resolveLocalAddr() (net.Addr, error) {
	localAddrStr := os.Getenv("LOCAL_ADDR")
	localPortStr := os.Getenv("LOCAL_PORT")

	if localAddrStr == "" || localPortStr == "" {
		return nil, fmt.Errorf("LOCAL_ADDR or LOCAL_PORT not set") //
	}

	address := fmt.Sprintf("%s:%s", localAddrStr, localPortStr)
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve TCP address for local: %w", err)
	}
	return addr, nil
}
