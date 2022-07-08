package ntp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"

	"github.com/gosnmp/gosnmp"
	"github.com/sirupsen/logrus"
)

// Server NTP server
type Server struct {
	conf *Config
	ln   *net.UDPConn
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
	err := gosnmp.Default.Connect()
	if err != nil {
		errChan <- fmt.Errorf("failed to connect snmp trap server: %v", err)
		return errChan
	}

	go func() {
		pdu1 := gosnmp.SnmpPDU{
			Name:  "1.3.6.1.4.1.326.2.6.1.1",
			Type:  gosnmp.Integer,
			Value: 0,
		}
		pdu2 := gosnmp.SnmpPDU{
			Name:  "1.3.6.1.4.1.326.2.6.1.2",
			Type:  gosnmp.Integer,
			Value: 10,
		}
		pdu3 := gosnmp.SnmpPDU{
			Name:  "1.3.6.1.4.1.326.2.6.1.3",
			Type:  gosnmp.Integer,
			Value: 20,
		}
		pdu4 := gosnmp.SnmpPDU{
			Name:  "1.3.6.1.4.1.326.2.6.1.4",
			Type:  gosnmp.Integer,
			Value: 30,
		}
		trap := gosnmp.SnmpTrap{
			Variables: []gosnmp.SnmpPDU{pdu1, pdu2, pdu3, pdu4},
		}
		for {
			if _, err := gosnmp.Default.SendTrap(trap); err != nil {
				logrus.WithField("prefix", "ntp.udp").
					Errorf("failed to send snmp data: %v", err)
			}
			time.Sleep(time.Second * 3)
		}
	}()
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
