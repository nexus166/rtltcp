package rtltcp

// https://github.com/bemasher/rtltcp

import (
	"encoding/binary"
	"fmt"
	"net"
)

var dongleMagic = [...]byte{'R', 'T', 'L', '0'}

// Contains dongle information and an embedded tcp connection to the spectrum server
type SDR struct {
	*net.TCPConn
	Flags Flags
	Info  DongleInfo
}

// Contains the Magic number, tuner information and the number of valid gain values.
type DongleInfo struct {
	Magic     [4]byte
	Tuner     Tuner
	GainCount uint32 // Useful for setting gain by index
}

func (d DongleInfo) String() string {
	return fmt.Sprintf("{Magic:%q Tuner:%s GainCount:%d}", d.Magic, d.Tuner, d.GainCount)
}

// Checks that the magic number received matches the expected byte string 'RTL0'.
func (d DongleInfo) Valid() bool {
	return d.Magic == dongleMagic
}

// Provides mapping of tuner value to tuner string.
type Tuner uint32

func (t Tuner) String() string {
	switch t {
	case 1:
		return "E4000"
	case 2:
		return "FC0012"
	case 3:
		return "FC0013"
	case 4:
		return "FC2580"
	case 5:
		return "R820T"
	case 6:
		return "R828D"
	}
	return "UNKNOWN"
}

func (sdr SDR) execute(cmd command) (err error) {
	return binary.Write(sdr.TCPConn, binary.BigEndian, cmd)
}

type command struct {
	command   uint8
	Parameter uint32
}
