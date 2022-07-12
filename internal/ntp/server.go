package ntp

import (
	"net"

	"github.com/gosnmp/gosnmp"
)

// Server NTP server
type Server struct {
	conf  *Config
	lnUDP *net.UDPConn
	lnTCP *net.TCPListener
}

// NewServer create NTP server
func NewServer(conf *Config) (*Server, error) {
	gosnmp.Default.Target = "127.0.0.1"
	gosnmp.Default.Port = 1169
	gosnmp.Default.Community = "1234qwer"
	server := Server{
		conf: conf,
	}
	return &server, nil
}

// Start start NTP server
func (s *Server) Start() chan error {
	errChan := make(chan error, 1)
	go s._startSNMP(errChan)
	go s._startUDP(errChan)
	go s._startTCP(errChan)
	return errChan
}

// Close close NTP server
func (s *Server) Close() error {
	s.lnUDP.Close()
	s.lnTCP.Close()
	return nil
}
