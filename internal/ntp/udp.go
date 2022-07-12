package ntp

import (
	"net"
	"time"

	"github.com/sirupsen/logrus"
)

func (s *Server) _startUDP(errChan chan error) {
	udpAddr, err := net.ResolveUDPAddr("udp", s.conf.Listener)
	if err != nil {
		errChan <- err
		return
	}
	ln, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		errChan <- err
		return
	}
	s.lnUDP = ln
	logrus.WithField("prefix", "ntp.udp").
		Infof("ntp listener: %s", s.conf.Listener)
	buf := make([]byte, ntpBufSize)
	for {
		n, addr, err := s.lnUDP.ReadFromUDP(buf)
		if err != nil {
			logrus.WithField("prefix", "ntp.handler").WithError(err).
				Errorf("read UDP packet failed")
			continue
		}
		data := buf[:n]
		logrus.WithField("prefix", "ntp.handler").
			Infof("read UDP packet from: [%s]", addr)
		now, _, err := getTimeWithTANode()
		if err != nil {
			logrus.WithField("prefix", "ntp.handler").Errorf("get TA time failed: %s", err.Error())
			continue
		}
		go s.udpNTPHandler(now, addr, data)
	}
}
func (s *Server) udpNTPHandler(reciveTime time.Time, addr *net.UDPAddr, data []byte) {
	if !validFormat(data) {
		return
	}
	td, _, err := genResp(reciveTime, data)
	if err != nil {
		return
	}
	if _, err = s.lnUDP.WriteTo(td, addr); err != nil {
		return
	}
}
