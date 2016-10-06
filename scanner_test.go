package lmsensors

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestScannerScan(t *testing.T) {
	tests := []struct {
		name    string
		fs      filesystem
		devices []*Device
	}{
		{
			name: "coretemp device",
			fs: &memoryFilesystem{
				matches: []string{"/sys/devices/platform/coretemp.0/name"},
				files: map[string]memoryFile{
					"/sys/devices/platform/coretemp.0": {
						info: &memoryFileInfo{
							isDir: true,
						},
					},
					"/sys/devices/platform/coretemp.0/name": {
						contents: "coretemp",
					},
					"/sys/devices/platform/coretemp.0/temp1_crit": {
						contents: "100000",
					},
					"/sys/devices/platform/coretemp.0/temp1_crit_alarm": {
						contents: "0",
					},
					"/sys/devices/platform/coretemp.0/temp1_input": {
						contents: "40000",
					},
					"/sys/devices/platform/coretemp.0/temp1_label": {
						contents: "Core 0",
					},
					"/sys/devices/platform/coretemp.0/temp1_max": {
						contents: "80000",
					},
					"/sys/devices/platform/coretemp.0/temp2_crit": {
						contents: "100000",
					},
					"/sys/devices/platform/coretemp.0/temp2_crit_alarm": {
						contents: "0",
					},
					"/sys/devices/platform/coretemp.0/temp2_input": {
						contents: "42000",
					},
					"/sys/devices/platform/coretemp.0/temp2_label": {
						contents: "Core 1",
					},
					"/sys/devices/platform/coretemp.0/temp2_max": {
						contents: "80000",
					},
				},
			},
			devices: []*Device{{
				Name: "coretemp",
				Sensors: []Sensor{
					&TemperatureSensor{
						Name:          "temp1",
						Label:         "Core 0",
						Current:       40.0,
						High:          80.0,
						Critical:      100.0,
						CriticalAlarm: false,
					},
					&TemperatureSensor{
						Name:          "temp2",
						Label:         "Core 1",
						Current:       42.0,
						High:          80.0,
						Critical:      100.0,
						CriticalAlarm: false,
					},
				},
			}},
		},
		{
			name: "it8728 device",
			fs: &memoryFilesystem{
				matches: []string{"/sys/devices/platform/it87.2608/name"},
				files: map[string]memoryFile{
					"/sys/devices/platform/it87.2608": {
						info: &memoryFileInfo{
							isDir: true,
						},
					},
					"/sys/devices/platform/it87.2608/name": {
						contents: "it8728",
					},
					"/sys/devices/platform/it87.2608/fan1_alarm": {
						contents: "0",
					},
					"/sys/devices/platform/it87.2608/fan1_beep": {
						contents: "1",
					},
					"/sys/devices/platform/it87.2608/fan1_input": {
						contents: "1010",
					},
					"/sys/devices/platform/it87.2608/fan1_min": {
						contents: "10",
					},
					"/sys/devices/platform/it87.2608/intrusion0_alarm": {
						contents: "1",
					},
					"/sys/devices/platform/it87.2608/temp1_alarm": {
						contents: "0",
					},
					"/sys/devices/platform/it87.2608/temp1_beep": {
						contents: "1",
					},
					"/sys/devices/platform/it87.2608/temp1_input": {
						contents: "43000",
					},
					"/sys/devices/platform/it87.2608/temp1_max": {
						contents: "127000",
					},
					"/sys/devices/platform/it87.2608/temp1_type": {
						contents: "4",
					},
				},
			},
			devices: []*Device{{
				Name: "it8728",
				Sensors: []Sensor{
					&FanSensor{
						Name:    "fan1",
						Alarm:   false,
						Beep:    true,
						Current: 1010,
						Minimum: 10,
					},
					&IntrusionSensor{
						Name:  "intrusion0",
						Alarm: true,
					},
					&TemperatureSensor{
						Name:    "temp1",
						Alarm:   false,
						Beep:    true,
						Type:    TemperatureSensorTypeThermistor,
						Current: 43.0,
						High:    127.0,
					},
				},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scanner{fs: tt.fs}

			devices, err := s.Scan()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if want, got := tt.devices, devices; !reflect.DeepEqual(want, got) {
				t.Fatalf("unexpected Devices:\n- want:\n%v\n-  got:\n%v",
					devicesStr(want), devicesStr(got))
			}
		})
	}
}

func devicesStr(ds []*Device) string {
	var out string
	for _, d := range ds {
		out += fmt.Sprintf("device: %q [%d sensors]\n", d.Name, len(d.Sensors))

		for _, s := range d.Sensors {
			out += fmt.Sprintf("  - sensor: %#v\n", s)
		}
	}

	return out
}

var _ filesystem = &memoryFilesystem{}

// A memoryFilesystem is an in-memory implementation of filesystem, used for
// tests.
type memoryFilesystem struct {
	globCalled bool
	matches    []string
	files      map[string]memoryFile
}

func (fs *memoryFilesystem) Glob(_ string) ([]string, error) {
	if !fs.globCalled {
		fs.globCalled = true
		return fs.matches, nil
	}

	return nil, nil
}

func (fs *memoryFilesystem) ReadFile(filename string) (string, error) {
	if f, ok := fs.files[filename]; ok {
		return f.contents, nil
	}

	return "", fmt.Errorf("file %q not in memory", filename)
}

func (fs *memoryFilesystem) Walk(root string, walkFn filepath.WalkFunc) error {
	if _, ok := fs.files[root]; !ok {
		return fmt.Errorf("file %q not in memory", root)
	}

	for k, v := range fs.files {
		info := v.info
		if info == nil {
			info = &memoryFileInfo{}
		}

		if err := walkFn(k, info, nil); err != nil {
			return err
		}
	}

	return nil
}

// A memoryFile is an in-memory file used by memoryFilesystem.
type memoryFile struct {
	contents string
	info     os.FileInfo
}

var _ os.FileInfo = &memoryFileInfo{}

// A memoryFileInfo is an os.FileInfo used by memoryFiles.
type memoryFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
}

func (fi *memoryFileInfo) Name() string       { return fi.name }
func (fi *memoryFileInfo) Size() int64        { return fi.size }
func (fi *memoryFileInfo) Mode() os.FileMode  { return fi.mode }
func (fi *memoryFileInfo) ModTime() time.Time { return fi.modTime }
func (fi *memoryFileInfo) IsDir() bool        { return fi.isDir }
func (fi *memoryFileInfo) Sys() interface{}   { return nil }
