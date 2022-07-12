package ntp

import (
	"bytes"
	"encoding/binary"
	"time"

	"github.com/sirupsen/logrus"
)

func genResp(reciveTime time.Time,
	req []byte) ([]byte, *msg, error) {
	var accuracy int8
	now, accuracy, err := getTimeWithTANode()
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
	now, _, err = getTimeWithTANode()
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

func getTimeWithTANode() (time.Time, int8, error) {
	t := time.Now()
	return t, -29, nil
}
