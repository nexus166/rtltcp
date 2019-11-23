package main

import (
	"flag"
	"fmt"
	"hash"
	"io"
	"math/big"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/nexus166/rtltcp"
	"golang.org/x/crypto/sha3"
)

var (
	listenOn         = flag.String("path", "/dev/rrand", "the UNIX socket to listen on")
	dongleTCPAddr    = flag.String("dongle-tcp", "127.0.0.1:1234", "rtl_tcp -a")
	dongleStaticFreq = flag.Int("freq", 0, "static frequency to tune on (will randomly hop every 1s if unset)")
	daemon           = flag.Bool("daemon", false, "start process in background")
)

func init() {
	flag.Parse()
	args := os.Args[1:]
	// restart in background if required
	for i := 0; i < len(args); i++ {
		if strings.Contains("-"+"daemon", args[i]) {
			*daemon = true
			args[i] = "-daemon=false"
			break
		}
	}
	if *daemon {
		i, err := fork(os.Args[0], args...)
		if err != nil {
			panic(err)
		}
		pid := strconv.Itoa(i)
		fmt.Println(pid)
		os.Exit(0)
	}
}

func main() {
	x := &chaosSDR{
		hash: sha3.New512,
	}
	sdrRemote, err := net.ResolveTCPAddr("tcp", *dongleTCPAddr)
	if err != nil {
		panic(err)
	}
	if err = x.sdr.Connect(sdrRemote); err != nil {
		panic(err)
	}
	defer x.sdr.Close()
	if *dongleStaticFreq == 0 {
		go func() {
			rand.Seed(time.Now().Unix())
			for {
				if err = x.sdr.SetCenterFreq(uint32(rand.Intn(1500000000))); err != nil {
					panic(err)
				}
				time.Sleep(time.Second)
			}
		}()
	} else {
		if err = x.sdr.SetCenterFreq(uint32(*dongleStaticFreq)); err != nil {
			panic(err)
		}
	}
	errs := make(chan error)
	go fifoListener(*listenOn, x, errs)
	if !*daemon {
		fmt.Println(x.sdr.Info)
		for err = range errs {
			fmt.Println("error: ", err)
		}
	} else {
		for err = range errs {
			_ = err
		}
	}
}

type chaosSDR struct {
	sdr  rtltcp.SDR
	hash func() hash.Hash
}

func (r *chaosSDR) Read(p []byte) (int, error) {
	n, err := r.sdr.Read(p)
	if err != nil {
		return n, err
	}
	h := r.hash()
	_, _ = h.Write(p[:n])
	k := h.Sum(nil)
	h.Reset()
	l := len(k)
	for i := range p {
		p[i] = p[i] ^ k[i%l]
	}
	_, _ = h.Write(p)
	k = h.Sum(nil)
	for i := range p {
		p[i] = p[i] ^ k[i%l]
	}
	return n, err
}

func (r *chaosSDR) Int63() int64 {
	p := make([]byte, 64)
	n, err := r.Read(p)
	if err != nil {
		panic(err)
	}
	x := big.NewInt(0)
	x = x.SetBytes(p[:n])
	return x.Int64()
}

func (r *chaosSDR) Uint64() uint64 {
	p := make([]byte, 64)
	n, err := r.Read(p)
	if err != nil {
		panic(err)
	}
	x := big.NewInt(0)
	x = x.SetBytes(p[:n])
	return x.Uint64()
}

func (r *chaosSDR) Seed(ch int64) {
	if err := r.sdr.SetCenterFreq(uint32(ch)); err != nil {
		panic(err)
	}
}

func fifoListener(path string, sdr *chaosSDR, errC chan error) {
	if err := os.RemoveAll(path); err != nil {
		panic(err)
	}
	err := syscall.Mkfifo(path, 0644)
	if err != nil {
		panic(err)
	}
	fifoL, err := os.OpenFile(path, os.O_WRONLY, os.ModeNamedPipe)
	if err != nil {
		panic(err)
	}
	for {
		_, err := io.Copy(fifoL, rand.New(sdr))
		if err != nil {
			errC <- err
		}
	}
}

func fork(bin string, args ...string) (int, error) {
	cmd := exec.Command(bin, args...)
	//cmd.Env = os.Environ()
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.ExtraFiles = nil
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}
	if err := cmd.Start(); err != nil {
		return 0, err
	}
	return cmd.Process.Pid, nil
}
