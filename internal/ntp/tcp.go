package ntp

import (
	"fmt"
	"io"
	"net"

	"github.com/sirupsen/logrus"
)

func (s *Server) _startTCP(errChan chan error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", s.conf.Listener)
	if err != nil {
		errChan <- err
		return
	}
	ln, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		errChan <- err
		return
	}
	s.lnTCP = ln
	logrus.WithField("prefix", "ntp.tcp").
		Infof("ntp listener: %s", s.conf.Listener)
	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			logrus.WithField("prefix", "ntp.tcp").
				Errorf("failed to accept: %v", err)
			continue
		}
		client := &tcpClient{
			conn: conn,
		}
		go func() {
			if err := client.start(); err != nil {
				conn.Close()
			}
		}()
	}
}

type tcpClient struct {
	conn *net.TCPConn
}

func (tc *tcpClient) start() error {
	buf := make([]byte, ntpBufSize)
	for {
		n, err := tc.conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				tc.conn.Close()
				return nil
			}
			logrus.WithField("prefix", "ntp.handler").WithError(err).
				Errorf("read tcp packet failed")
			return err
		}
		data := buf[:n]
		logrus.WithField("prefix", "ntp.handler").
			Infof("read tcp packet from: [%s]", tc.conn.RemoteAddr().String())
		reciveTime, _, err := getTimeWithTANode()
		if err != nil {
			logrus.WithField("prefix", "ntp.handler").
				Errorf("get TA time failed: %s", err.Error())
			return err
		}
		if !validFormat(data) {
			return fmt.Errorf("failed to validate format")
		}
		td, _, err := genResp(reciveTime, data)
		if err != nil {
			return err
		}
		if _, err = tc.conn.Write(td); err != nil {
			logrus.WithField("prefix", "ntp.handler").
				Errorf("write time response failed: %s", err.Error())
			return err
		}
	}
}
