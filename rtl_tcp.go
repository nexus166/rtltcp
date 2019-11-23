package rtltcp

// https://github.com/bemasher/rtltcp

import (
	"encoding/binary"
	"fmt"
	"net"
)

// Command constants defined in rtl_tcp.c
const (
	centerFreq = iota + 1
	sampleRate
	tunerGainMode
	tunerGain
	freqCorrection
	tunerIfGain
	testMode
	agcMode
	directSampling
	offsetTuning
	rtlXtalFreq
	tunerXtalFreq
	gainByIndex
)

// Give an address of the form "127.0.0.1:1234" connects to the spectrum
// server at the given address or returns an error. The user is responsible
// for closing this connection. If addr is nil, use "127.0.0.1:1234" or
// command line flag value.
func (sdr *SDR) Connect(addr *net.TCPAddr) (err error) {
	if addr == nil {
		if sdr.Flags.ServerAddr == "" {
			sdr.Flags.ServerAddr = "127.0.0.1:1234"
		}

		// Parse and resolve rtl_tcp server address.
		addr, err = net.ResolveTCPAddr("tcp", sdr.Flags.ServerAddr)
		if err != nil {
			return
		}
	}

	sdr.TCPConn, err = net.DialTCP("tcp", nil, addr)
	if err != nil {
		err = fmt.Errorf("Error connecting to spectrum server: %s", err)
		return
	}

	// If we exit this function due to an error, close the connection.
	defer func() {
		if err != nil {
			sdr.Close()
		}
	}()

	err = binary.Read(sdr.TCPConn, binary.BigEndian, &sdr.Info)
	if err != nil {
		err = fmt.Errorf("Error getting dongle information: %s", err)
		return
	}

	if !sdr.Info.Valid() {
		err = fmt.Errorf("Invalid magic number: expected %q received %q", dongleMagic, sdr.Info.Magic)
	}

	return
}

// Set the center frequency in Hz.
func (sdr SDR) SetCenterFreq(freq uint32) (err error) {
	return sdr.execute(command{centerFreq, freq})
}

// Set the sample rate in Hz.
func (sdr SDR) SetSampleRate(rate uint32) (err error) {
	return sdr.execute(command{sampleRate, rate})
}

// Set gain in tenths of dB. (197 => 19.7dB)
func (sdr SDR) SetGain(gain uint32) (err error) {
	return sdr.execute(command{tunerGain, gain})
}

// Set the Tuner AGC, true to enable.
func (sdr SDR) SetGainMode(state bool) (err error) {
	if state {
		return sdr.execute(command{tunerGainMode, 0})
	}
	return sdr.execute(command{tunerGainMode, 1})
}

// Set gain by index, must be <= DongleInfo.GainCount
func (sdr SDR) SetGainByIndex(idx uint32) (err error) {
	if idx > sdr.Info.GainCount {
		return fmt.Errorf("invalid gain index: %d", idx)
	}
	return sdr.execute(command{gainByIndex, idx})
}

// Set frequency correction in ppm.
func (sdr SDR) SetFreqCorrection(ppm uint32) (err error) {
	return sdr.execute(command{freqCorrection, ppm})
}

// Set tuner intermediate frequency stage and gain.
func (sdr SDR) SetTunerIfGain(stage, gain uint16) (err error) {
	return sdr.execute(command{tunerIfGain, (uint32(stage) << 16) | uint32(gain)})
}

// Set test mode, true for enabled.
func (sdr SDR) SetTestMode(state bool) (err error) {
	if state {
		return sdr.execute(command{testMode, 1})
	}
	return sdr.execute(command{testMode, 0})
}

// Set RTL AGC mode, true for enabled.
func (sdr SDR) SetAGCMode(state bool) (err error) {
	if state {
		return sdr.execute(command{agcMode, 1})
	}
	return sdr.execute(command{agcMode, 0})
}

// Set direct sampling mode.
func (sdr SDR) SetDirectSampling(state bool) (err error) {
	if state {
		return sdr.execute(command{directSampling, 1})
	}
	return sdr.execute(command{directSampling, 0})
}

// Set offset tuning, true for enabled.
func (sdr SDR) SetOffsetTuning(state bool) (err error) {
	if state {
		return sdr.execute(command{offsetTuning, 1})
	}
	return sdr.execute(command{offsetTuning, 0})
}

// Set RTL xtal frequency.
func (sdr SDR) SetRTLXtalFreq(freq uint32) (err error) {
	return sdr.execute(command{rtlXtalFreq, freq})
}

// Set tuner xtal frequency.
func (sdr SDR) SetTunerXtalFreq(freq uint32) (err error) {
	return sdr.execute(command{tunerXtalFreq, freq})
}
