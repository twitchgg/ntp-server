package ntp

import "time"

type ntpTimeShort uint32
type ntpTime uint64
type mode uint8

// The LeapIndicator is used to warn if a leap second should be inserted
// or deleted in the last minute of the current month.
type LeapIndicator uint8

var (
	ntpEpoch = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
)

const (
	ntpBufSize = 512
	// LiNoWarning LI_NO_WARNING
	LiNoWarning = 0
	// LeapAddSecond indicates the last minute of the day has 61 seconds.
	LeapAddSecond = 1
	// LeapDelSecond indicates the last minute of the day has 59 seconds.
	LeapDelSecond = 2
	// LeapNotInSync indicates an unsynchronized leap second.
	LeapNotInSync = 3
	// LiAlarmCondition LI_ALARM_CONDITION
	LiAlarmCondition = 3
	// VnFirst VN_FIRST
	VnFirst = 1
	// VnLast VN_LAST
	VnLast = 4
	// ModeClient MODE_CLIENT
	ModeClient = 3
	// From1900To1970 FROM_1900_TO_1970
	From1900To1970 = 2208988800

	defaultNtpVersion = 4
	nanoPerSec        = 1000000000
	maxStratum        = 16
	defaultTimeout    = 5 * time.Second
	maxPollInterval   = (1 << 17) * time.Second
	maxDispersion     = 16 * time.Second
)

type msg struct {
	LiVnMode       uint8 // Leap Indicator (2) + Version (3) + Mode (3)
	Stratum        uint8
	Poll           int8
	Precision      int8
	RootDelay      ntpTimeShort
	RootDispersion ntpTimeShort
	ReferenceID    uint32
	ReferenceTime  ntpTime
	OriginTime     ntpTime
	ReceiveTime    ntpTime
	TransmitTime   ntpTime
}

type packet struct {
	Settings       uint8  // leap yr indicator, ver number, and mode
	Stratum        uint8  // stratum of local clock
	Poll           int8   // poll exponent
	Precision      int8   // precision exponent
	RootDelay      uint32 // root delay
	RootDispersion uint32 // root dispersion
	ReferenceID    uint32 // reference id
	RefTimeSec     uint32 // reference timestamp sec
	RefTimeFrac    uint32 // reference timestamp fractional
	OrigTimeSec    uint32 // origin time secs
	OrigTimeFrac   uint32 // origin time fractional
	RxTimeSec      uint32 // receive time secs
	RxTimeFrac     uint32 // receive time frac
	TxTimeSec      uint32 // transmit time secs
	TxTimeFrac     uint32 // transmit time frac
}

func validFormat(req []byte) bool {
	var l = req[0] >> 6
	var v = (req[0] << 2) >> 5
	var m = (req[0] << 5) >> 5
	if (l == LiNoWarning) || (l == LiAlarmCondition) {
		if (v >= VnFirst) && (v <= VnLast) {
			if m == ModeClient {
				return true
			}
		}
	}
	return false
}

// setVersion sets the NTP protocol version on the message.
func (m *msg) setVersion(v int) {
	m.LiVnMode = (m.LiVnMode & 0xc7) | uint8(v)<<3
}

// setMode sets the NTP protocol mode on the message.
func (m *msg) setMode(md mode) {
	m.LiVnMode = (m.LiVnMode & 0xf8) | uint8(md)
}

func toNtpTime(t time.Time) ntpTime {
	nsec := uint64(t.Sub(ntpEpoch))
	sec := nsec / nanoPerSec
	// Round up the fractional component so that repeated conversions
	// between time.Time and ntpTime do not yield continually decreasing
	// results.
	frac := (((nsec - sec*nanoPerSec) << 32) + nanoPerSec - 1) / nanoPerSec
	return ntpTime(sec<<32 | frac)
}

func (t *ntpTime) Time() time.Time {
	return ntpEpoch.Add(t.Duration())
}

// Duration interprets the fixed-point ntpTime as a number of elapsed seconds
// and returns the corresponding time.Duration value.
func (t *ntpTime) Duration() time.Duration {
	sec := (*t >> 32) * nanoPerSec
	frac := (*t & 0xffffffff) * nanoPerSec >> 32
	return time.Duration(sec + frac)
}
