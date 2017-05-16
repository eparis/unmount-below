package mounts

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var (
	slashRegexp = regexp.MustCompile("/")
)

// 17 60 0:17 / /sys rw,nosuid,nodev,noexec,relatime shared:6 - sysfs sysfs rw,seclabel
type Mount struct {
	field1      string
	field2      string
	field3      string
	source      string
	target      string
	opts        string
	propagation string
	dash        string
	fstype      string
	field10     string
	moreOpts    string
}

type Mounts []Mount

func (mnts Mounts) Len() int      { return len(mnts) }
func (mnts Mounts) Swap(i, j int) { mnts[i], mnts[j] = mnts[j], mnts[i] }
func (mnts Mounts) Less(i, j int) bool {
	iPath := mnts[i].target
	jPath := mnts[j].target

	iSlashes := len(slashRegexp.FindAllStringIndex(iPath, -1))
	jSlashes := len(slashRegexp.FindAllStringIndex(jPath, -1))
	if iSlashes < jSlashes {
		return true
	} else if iSlashes > jSlashes {
		return false
	}

	return iPath < jPath
}

func (mnts Mounts) mountsUnder(path string) (Mounts, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	out := Mounts{}
	for _, mnt := range mnts {
		if strings.Contains(mnt.target, path) && mnt.target != path {
			out = append(out, mnt)
		}
	}
	return out, err
}

func LongestMountUnder(path string) (string, error) {
	mnts, err := parseMounts()
	if err != nil {
		return "", err
	}

	mntsUnder, err := mnts.mountsUnder(path)
	if err != nil {
		return "", err
	}
	if len(mntsUnder) == 0 {
		return "", os.ErrNotExist
	}
	sort.Sort(mntsUnder)
	return mntsUnder[len(mntsUnder)-1].target, nil
}

func parseLine(in string) (Mount, error) {
	mount := Mount{}
	fields := strings.Fields(in)
	if len(fields) < 11 {
		return mount, fmt.Errorf("Unable to parse mountinfo line: %q", in)
	}
	if len(fields) > 11 {
		log.Printf("Potentially mis-parsed mount line. Len is %d not 10: %s", len(fields), in)
	}
	mount.field1 = fields[0]
	mount.field2 = fields[1]
	mount.field3 = fields[2]
	mount.source = fields[3]
	mount.target = fields[4]
	mount.opts = fields[5]
	mount.propagation = fields[6]
	mount.dash = fields[7]
	mount.fstype = fields[9]
	mount.field10 = fields[8]
	mount.moreOpts = fields[10]

	if mount.dash != "-" {
		log.Printf("Potentially mis-parsed mount line. Dash is %q not \"-\": %q", mount.dash, in)
	}
	return mount, nil
}

func parseMounts() (Mounts, error) {
	mnts := Mounts{}

	file, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return mnts, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		mount, err := parseLine(scanner.Text())
		if err != nil {
			return mnts, err
		}
		mnts = append(mnts, mount)
	}
	if err := scanner.Err(); err != nil {
		return mnts, err
	}
	return mnts, nil
}
