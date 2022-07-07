package test

import (
	"fmt"
	"testing"

	"github.com/beevik/ntp"
	ntsc "ntsc.ac.cn/ta/ntp-server/internal/ntp"
)

func TestNNTPServer(t *testing.T) {
	s, err := ntsc.NewServer(&ntsc.Config{
		Listener: "0.0.0.0:1123",
	})
	if err != nil {
		t.Fatal(err)
	}
	err = <-s.Start()
	if err != nil {
		t.Fatal(err)
	}
}

func TestNTP(t *testing.T) {
	response, err := ntp.QueryWithOptions("127.0.0.1", ntp.QueryOptions{
		Port: 123,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("offset:", response.ClockOffset)
}
