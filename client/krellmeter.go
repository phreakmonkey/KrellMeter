package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

  "gopkg.in/ini.v1"
	"github.com/shirou/gopsutil/cpu"
	"github.com/tarm/serial"
)

var (
	SERIALPORT string  = ""
	SERIALBAUD int     = 115200
	AMAX       float64 = 225
	BMAX       float64 = 225
	afactor    float64 = AMAX / 100.0
	bfactor    float64 = BMAX / 100.0
)

func send(s *serial.Port, a_percent int, b_percent int) {
	buf := make([]byte, 128)
	astr := "A" + strconv.FormatFloat(float64(a_percent)*afactor, 'f', 0, 64) + "\n"
	bstr := "B" + strconv.FormatFloat(float64(b_percent)*bfactor, 'f', 0, 64) + "\n"
	s.Write([]byte(astr + bstr))
	_, err := s.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
}

func npipe() *os.File {
	fmt.Println("Launching nvidia-smi")
	cmd := exec.Command("nvidia-smi", "-l", "2")
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()
	return stdout.(*os.File)
}

func parse_nvidia(buf []byte, l int) int {
	var pct int = -1
	if l > 100 {
		if buf[0] == 0x7C && buf[1] == 0x20 {
			s := strings.Split(string(buf), "\n")
			switch s[1][60] {
			case 0x20:
				pct, _ = strconv.Atoi(s[1][61:63])
			case 0x31:
				pct = 100
			}
		}
	}
	return pct
}

func main() {
	var gpu     int = 0
  var inifile string = "krellmeter.ini"

  if len(os.Args) > 1 {
    inifile = os.Args[1]
  }
  cfg, err := ini.Load(inifile)
  if err != nil {
    fmt.Printf("Fail to read ini file: %v", err)
    os.Exit(1)
  }

  SERIALPORT = cfg.Section("serial").Key("port").String()
  SERIALBAUD = cfg.Section("serial").Key("baud").MustInt()
  AMAX = float64(cfg.Section("meter1").Key("max").MustInt(225))
  BMAX = float64(cfg.Section("meter2").Key("max").MustInt(225))
  afactor = AMAX / 100.0
  bfactor = BMAX / 100.0

	c := &serial.Config{Name: SERIALPORT, Baud: SERIALBAUD, ReadTimeout: time.Millisecond * 5}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	// Zero meters
	send(s, 0, 0)
	time.Sleep(250 * time.Millisecond)

	// Run up & run down:
	for a := 0; a <= 100; a++ {
		send(s, a, a)
		time.Sleep(25 * time.Millisecond)
		if a == 25 || a == 50 || a == 75 || a == 100 {
			time.Sleep(250 * time.Millisecond)
		}
	}
	for a := 100; a >= 0; a-- {
		send(s, a, a)
		time.Sleep(25 * time.Millisecond)
	}

	nvidia_output := npipe()
	buf := make([]byte, 8192)
	for {
		i, _ := nvidia_output.Read(buf)
		j := parse_nvidia(buf, i)
		if j > -1 {
			gpu = j
		}
		p, _ := cpu.Percent(0, false)
		send(s, int(p[0]), gpu)
		time.Sleep(250 * time.Millisecond)
	}
}
