package rtltcp

// https://github.com/bemasher/rtltcp

import "flag"

type Flags struct {
	ServerAddr     string
	CenterFreq     ScientificNotation
	SampleRate     ScientificNotation
	TunerGainMode  bool
	TunerGain      float64
	FreqCorrection int
	TestMode       bool
	AgcMode        bool
	DirectSampling bool
	OffsetTuning   bool
	RtlXtalFreq    uint
	TunerXtalFreq  uint
	GainByIndex    uint
}

// Registers command line flags for rtltcp commands.
func (sdr *SDR) RegisterFlags() {
	flag.StringVar(&sdr.Flags.ServerAddr, "server", "127.0.0.1:1234", "address or hostname of rtl_tcp instance")
	flag.Var(&sdr.Flags.CenterFreq, "centerfreq", "center frequency to receive on")
	flag.Lookup("centerfreq").DefValue = "100M"
	flag.Var(&sdr.Flags.SampleRate, "samplerate", "sample rate")
	flag.Lookup("samplerate").DefValue = "2.4M"
	flag.BoolVar(&sdr.Flags.TunerGainMode, "tunergainmode", false, "enable/disable tuner gain")
	flag.Float64Var(&sdr.Flags.TunerGain, "tunergain", 0.0, "set tuner gain in dB")
	flag.IntVar(&sdr.Flags.FreqCorrection, "freqcorrection", 0, "frequency correction in ppm")
	flag.BoolVar(&sdr.Flags.TestMode, "testmode", false, "enable/disable test mode")
	flag.BoolVar(&sdr.Flags.AgcMode, "agcmode", false, "enable/disable rtl agc")
	flag.BoolVar(&sdr.Flags.DirectSampling, "directsampling", false, "enable/disable direct sampling")
	flag.BoolVar(&sdr.Flags.OffsetTuning, "offsettuning", false, "enable/disable offset tuning")
	flag.UintVar(&sdr.Flags.RtlXtalFreq, "rtlxtalfreq", 0, "set rtl xtal frequency")
	flag.UintVar(&sdr.Flags.TunerXtalFreq, "tunerxtalfreq", 0, "set tuner xtal frequency")
	flag.UintVar(&sdr.Flags.GainByIndex, "gainbyindex", 0, "set gain by index")
}

// Parses flags and executes commands associated with each flag. Should only
// be called once connected to rtl_tcp.
func (sdr SDR) HandleFlags() (err error) {
	// Catch any errors panicked while visiting flags.
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	flag.CommandLine.Visit(func(f *flag.Flag) {
		var err error
		switch f.Name {
		case "centerfreq":
			err = sdr.SetCenterFreq(uint32(sdr.Flags.CenterFreq))
		case "samplerate":
			err = sdr.SetSampleRate(uint32(sdr.Flags.SampleRate))
		case "tunergainmode":
			err = sdr.SetGainMode(sdr.Flags.TunerGainMode)
		case "tunergain":
			err = sdr.SetGain(uint32(sdr.Flags.TunerGain * 10.0))
		case "freqcorrection":
			err = sdr.SetFreqCorrection(uint32(sdr.Flags.FreqCorrection))
		case "testmode":
			err = sdr.SetTestMode(sdr.Flags.TestMode)
		case "agcmode":
			err = sdr.SetAGCMode(sdr.Flags.AgcMode)
		case "directsampling":
			err = sdr.SetDirectSampling(sdr.Flags.DirectSampling)
		case "offsettuning":
			err = sdr.SetOffsetTuning(sdr.Flags.OffsetTuning)
		case "rtlxtalfreq":
			err = sdr.SetRTLXtalFreq(uint32(sdr.Flags.RtlXtalFreq))
		case "tunerxtalfreq":
			err = sdr.SetTunerXtalFreq(uint32(sdr.Flags.TunerXtalFreq))
		case "gainbyindex":
			err = sdr.SetGainByIndex(uint32(sdr.Flags.GainByIndex))
		}

		// If we encounter an error, panic to catch in parent scope.
		if err != nil {
			panic(err)
		}
	})

	return
}
