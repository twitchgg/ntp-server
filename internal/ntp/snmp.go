package ntp

import (
	"fmt"
	"time"

	"github.com/gosnmp/gosnmp"
	"github.com/sirupsen/logrus"
)

func (s *Server) _startSNMP(errChan chan error) {
	err := gosnmp.Default.Connect()
	if err != nil {
		errChan <- fmt.Errorf("failed to connect snmp trap server: %v", err)
		return
	}
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
}
