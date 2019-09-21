package hardware

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"github.com/shirou/gopsutil/cpu"
	"net"
	"os/exec"
	"sort"
	"strings"
)

// 单例
var Hw *Hardware

func init() {
	Hw = &Hardware{
		Cpu: Cpu{},
		Mac: make([]string, 0),
	}
}

type Cpu struct {
	Model string `json:"model"`
	Num   int    `json:"num"`
}

type Hardware struct {
	Cpu  `json:"cpu"`
	Mac  []string `json:"mac"`
	Dmi  string   `json:"dmi"`
	code []byte
	md5  string
}

func (hw *Hardware) GetCode() (code []byte, err error) {
	hw.cpuInfo()
	hw.macInfo()
	hw.dmiInfo()
	code, err = json.Marshal(hw)
	hw.code = code
	return code, err
}

func (hw *Hardware) GetMd5() (string) {
	if len(hw.code) == 0 {
		if _, err := hw.GetCode(); err != nil {
			return ""
		}
	}
	if hw.md5 == "" {
		h := md5.New()
		h = sha256.New()
		h.Write(hw.code)
		hw.md5 = base64.StdEncoding.EncodeToString(h.Sum(nil))
	}
	return hw.md5
}

func (hw *Hardware) cpuInfo() {
	cpuStat, _ := cpu.Info()
	tmp := make(map[string]int)
	for _, c := range cpuStat {
		tmp[c.ModelName] += 1
	}
	for name, num := range tmp {
		hw.Cpu.Model = name
		hw.Cpu.Num = num
	}
}

func (hw *Hardware) macInfo() {
	interfaces, err := net.Interfaces()
	if err != nil {
		return
	}
	for _, inter := range interfaces {
		if inter.HardwareAddr.String() != "" {
			hw.Mac = append(hw.Mac, inter.HardwareAddr.String())
		}
	}
	sort.Strings(hw.Mac)
}

func (hw *Hardware) dmiInfo() {
	cmd := exec.Command("sh", "-c", "dmesg|grep DMI")
	if info, err := cmd.CombinedOutput(); err == nil {
		hw.Dmi = strings.TrimSpace(strings.Split(string(info), ":")[1])
	}
}
