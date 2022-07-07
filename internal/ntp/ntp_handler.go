package ntp

import (
	"net"
	"time"
)

func (s *Server) ntpHandler(reciveTime time.Time, addr *net.UDPAddr, data []byte) {
	if !validFormat(data) {
		return
	}
	td, _, err := s.genResp(reciveTime, addr, data)
	if err != nil {
		return
	}
	if _, err = s.ln.WriteTo(td, addr); err != nil {
		return
	}
}
