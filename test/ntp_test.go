package test

import (
	"fmt"
	"testing"
	"time"

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

func TestTimeFormat(t *testing.T) {
	str := "2022-07-14T00:03:57.296714700+08:00"
	s, err := time.Parse(time.RFC3339Nano, str)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(s)
}
