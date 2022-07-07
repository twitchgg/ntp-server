package ntp

import (
	"bytes"
	"encoding/binary"
	"net"
	"time"

	"github.com/sirupsen/logrus"
)

// Server NTP server
type Server struct {
	conf *Config
	ln   *net.UDPConn
}

// NewServer create NTP server
func NewServer(conf *Config) (*Server, error) {
	server := Server{
		conf: conf,
	}
	return &server, nil
}

// Start start NTP server
func (s *Server) Start() chan error {
	errChan := make(chan error, 1)
	udpAddr, err := net.ResolveUDPAddr("udp", s.conf.Listener)
	if err != nil {
		errChan <- err
		return errChan
	}
	ln, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		errChan <- err
		return errChan
	}
	s.ln = ln
	go func() {
		logrus.WithField("prefix", "ntp.udp").Infof("ntp listener: %s", s.conf.Listener)
		buf := make([]byte, ntpBufSize)
		for {
			n, addr, err := s.ln.ReadFromUDP(buf)
			if err != nil {
				logrus.WithField("prefix", "ntp.handler").WithError(err).
					Errorf("read UDP packet failed")
				continue
			}
			data := buf[:n]
			logrus.WithField("prefix", "ntp.handler").
				Infof("read UDP packet from: [%s]", addr)
			now, _, err := s.getTimeWithTANode()
			if err != nil {
				logrus.WithField("prefix", "ntp.handler").Errorf("get TA time failed: %s", err.Error())
				continue
			}
			go s.ntpHandler(now, addr, data)
		}
	}()
	return errChan
}

// Close close NTP server
func (s *Server) Close() error {
	return s.ln.Close()
}

func (s *Server) genResp(reciveTime time.Time,
	addr *net.UDPAddr, req []byte) ([]byte, *msg, error) {
	var accuracy int8
	now, accuracy, err := s.getTimeWithTANode()
	if err != nil {
		logrus.WithField("prefix", "ntp.handler").Errorf("get TA time failed: %s", err.Error())
		return nil, nil, err
	}
	var buf bytes.Buffer
	buf.Write(req)
	recvMsg := new(msg)
	testMsg := new(packet)
	if err = binary.Read(&buf, binary.BigEndian, recvMsg); err != nil {
		return nil, nil, err
	}
	buf.Reset()
	buf.Write(req)
	if err = binary.Read(&buf, binary.BigEndian, testMsg); err != nil {
		return nil, nil, err
	}

	respMsg := new(msg)
	respMsg.ReceiveTime = toNtpTime(reciveTime)
	respMsg.ReferenceID = uint32(0001)
	respMsg.Stratum = 2
	respMsg.Precision = accuracy
	respMsg.setMode(4)
	respMsg.setVersion(3)
	respMsg.Poll = 0

	respMsg.OriginTime = recvMsg.TransmitTime
	respMsg.ReferenceTime = toNtpTime(now)
	now, _, err = s.getTimeWithTANode()
	if err != nil {
		logrus.WithField("prefix", "ntp.handler").Errorf("get TA time failed: %s", err.Error())
		return nil, nil, err
	}
	respMsg.TransmitTime = toNtpTime(now)
	buf.Reset()
	if err := binary.Write(&buf, binary.BigEndian, respMsg); err != nil {
		return nil, nil, err
	}
	resp := buf.Bytes()
	return resp, respMsg, nil
}

func (s *Server) getTimeWithTANode() (time.Time, int8, error) {
	t := time.Now()
	return t, -29, nil
}
