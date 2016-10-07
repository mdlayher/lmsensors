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
				symlinks: map[string]string{
					"/sys/class/hwmon/hwmon1":                              "../../devices/platform/coretemp.0/hwmon/hwmon1",
					"/sys/devices/platform/coretemp.0/hwmon/hwmon1/device": "../../../coretemp.0",
				},
				files: []memoryFile{
					{
						name: "/sys/class/hwmon",
						info: &memoryFileInfo{
							isDir: true,
						},
					},
					{
						name: "/sys/class/hwmon/hwmon1",
						info: &memoryFileInfo{
							mode: os.ModeSymlink,
						},
					},
					{
						name: "/sys/devices/platform/coretemp.0",
						info: &memoryFileInfo{
							isDir: true,
						},
					},
					{
						name: "/sys/devices/platform/coretemp.0/hwmon/hwmon1/name",
						err:  os.ErrNotExist,
					},
					{
						name: "/sys/devices/platform/coretemp.0/hwmon/hwmon1/device",
						info: &memoryFileInfo{
						// TODO(mdlayher): why does this only work if this isn't a symlink,
						// even though it is in the actual filesystem (and the actual filesystem
						// exhibits the same behavior)?
						// mode: os.ModeSymlink,
						},
					},
					{
						name:     "/sys/devices/platform/coretemp.0/name",
						contents: "coretemp",
					},
					{
						name:     "/sys/devices/platform/coretemp.0/temp1_crit",
						contents: "100000",
					},
					{
						name:     "/sys/devices/platform/coretemp.0/temp1_crit_alarm",
						contents: "0",
					},
					{
						name:     "/sys/devices/platform/coretemp.0/temp1_input",
						contents: "40000",
					},
					{
						name:     "/sys/devices/platform/coretemp.0/temp1_label",
						contents: "Core 0",
					},
					{
						name:     "/sys/devices/platform/coretemp.0/temp1_max",
						contents: "80000",
					},
					{
						name:     "/sys/devices/platform/coretemp.0/temp2_crit",
						contents: "100000",
					},
					{
						name:     "/sys/devices/platform/coretemp.0/temp2_crit_alarm",
						contents: "0",
					},
					{
						name:     "/sys/devices/platform/coretemp.0/temp2_input",
						contents: "42000",
					},
					{
						name:     "/sys/devices/platform/coretemp.0/temp2_label",
						contents: "Core 1",
					},
					{
						name:     "/sys/devices/platform/coretemp.0/temp2_max",
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
				symlinks: map[string]string{
					"/sys/class/hwmon/hwmon2":                             "../../devices/platform/it87.2608/hwmon/hwmon2",
					"/sys/devices/platform/it87.2608/hwmon/hwmon2/device": "../../../it87.2608",
				},
				files: []memoryFile{
					{
						name: "/sys/class/hwmon",
						info: &memoryFileInfo{
							isDir: true,
						},
					},
					{
						name: "/sys/class/hwmon/hwmon2",
						info: &memoryFileInfo{
							mode: os.ModeSymlink,
						},
					},
					{
						name: "/sys/devices/platform/it87.2608",
						info: &memoryFileInfo{
							isDir: true,
						},
					},
					{
						name: "/sys/devices/platform/it87.2608/hwmon/hwmon2/name",
						err:  os.ErrNotExist,
					},
					{
						name: "/sys/devices/platform/it87.2608/hwmon/hwmon2/device",
						info: &memoryFileInfo{
						// TODO(mdlayher): why does this only work if this isn't a symlink,
						// even though it is in the actual filesystem (and the actual filesystem
						// exhibits the same behavior)?
						// mode: os.ModeSymlink,
						},
					},
					{
						name:     "/sys/devices/platform/it87.2608/name",
						contents: "it8728",
					},
					{
						name:     "/sys/devices/platform/it87.2608/fan1_alarm",
						contents: "0",
					},
					{
						name:     "/sys/devices/platform/it87.2608/fan1_beep",
						contents: "1",
					},
					{
						name:     "/sys/devices/platform/it87.2608/fan1_input",
						contents: "1010",
					},
					{
						name:     "/sys/devices/platform/it87.2608/fan1_min",
						contents: "10",
					},
					{
						name:     "/sys/devices/platform/it87.2608/in0_alarm",
						contents: "0",
					},
					{
						name:     "/sys/devices/platform/it87.2608/in0_beep",
						contents: "0",
					},
					{
						name:     "/sys/devices/platform/it87.2608/in0_input",
						contents: "1056",
					},
					{
						name:     "/sys/devices/platform/it87.2608/in0_max",
						contents: "3060",
					},
					{
						name:     "/sys/devices/platform/it87.2608/in1_alarm",
						contents: "0",
					},
					{
						name:     "/sys/devices/platform/it87.2608/in1_beep",
						contents: "0",
					},
					{
						name:     "/sys/devices/platform/it87.2608/in1_input",
						contents: "3384",
					},
					{
						name:     "/sys/devices/platform/it87.2608/in1_label",
						contents: "3VSB",
					},
					{
						name:     "/sys/devices/platform/it87.2608/in1_max",
						contents: "6120",
					},
					{
						name:     "/sys/devices/platform/it87.2608/intrusion0_alarm",
						contents: "1",
					},
					{
						name:     "/sys/devices/platform/it87.2608/temp1_alarm",
						contents: "0",
					},
					{
						name:     "/sys/devices/platform/it87.2608/temp1_beep",
						contents: "1",
					},
					{
						name:     "/sys/devices/platform/it87.2608/temp1_input",
						contents: "43000",
					},
					{
						name:     "/sys/devices/platform/it87.2608/temp1_max",
						contents: "127000",
					},
					{
						name:     "/sys/devices/platform/it87.2608/temp1_type",
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
					&VoltageSensor{
						Name:    "in0",
						Alarm:   false,
						Beep:    false,
						Current: 1.056,
						Maximum: 3.060,
					},
					&VoltageSensor{
						Name:    "in1",
						Label:   "3VSB",
						Alarm:   false,
						Beep:    false,
						Current: 3.384,
						Maximum: 6.120,
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
	symlinks map[string]string
	files    []memoryFile
}

func (fs *memoryFilesystem) ReadFile(filename string) (string, error) {
	for _, f := range fs.files {
		if f.name == filename {
			return f.contents, nil
		}
	}

	return "", fmt.Errorf("readfile: file %q not in memory", filename)
}

func (fs *memoryFilesystem) Readlink(name string) (string, error) {
	if l, ok := fs.symlinks[name]; ok {
		return l, nil
	}

	return "", fmt.Errorf("readlink: symlink %q not in memory", name)
}

func (fs *memoryFilesystem) Stat(name string) (os.FileInfo, error) {
	for _, f := range fs.files {
		if f.name == name {
			info := f.info
			if info == nil {
				info = &memoryFileInfo{}
			}

			return info, f.err
		}
	}

	return nil, fmt.Errorf("stat: file %q not in memory", name)
}

func (fs *memoryFilesystem) Walk(root string, walkFn filepath.WalkFunc) error {
	if _, err := fs.Stat(root); err != nil {
		return err
	}

	for _, f := range fs.files {
		info := f.info
		if info == nil {
			info = &memoryFileInfo{}
		}

		if err := walkFn(f.name, info, nil); err != nil {
			return err
		}
	}

	return nil
}

// A memoryFile is an in-memory file used by memoryFilesystem.
type memoryFile struct {
	name     string
	contents string
	info     os.FileInfo
	err      error
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
