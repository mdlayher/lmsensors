package lmsensors

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// A filesystem is an interface to a filesystem, used for testing.
type filesystem interface {
	Glob(pattern string) ([]string, error)
	Walk(root string, walkFn filepath.WalkFunc) error
	ReadFile(filename string) (string, error)
}

// A Scanner scans for Devices, so data can be read from their Sensors.
type Scanner struct {
	fs filesystem
}

// New creates a new Scanner.
func New() *Scanner {
	return &Scanner{
		fs: &systemFilesystem{},
	}
}

// Scan scans for Devices and their Sensors.
func (s *Scanner) Scan() ([]*Device, error) {
	// Determine common device locations in Linux /sys filesystem.
	paths, err := s.detectDevicePaths()
	if err != nil {
		return nil, err
	}

	var devices []*Device
	for _, p := range paths {
		d := &Device{}
		raw := make(map[string]map[string]string, 0)

		// Walk filesystem paths to fetch devices and sensors
		err := s.fs.Walk(filepath.Dir(p), func(path string, info os.FileInfo, err error) error {
			// Skip directories and anything that isn't a regular file
			if info.IsDir() || !info.Mode().IsRegular() {
				return nil
			}

			// Skip some files that can't be read or don't provide useful
			// sensor information
			file := filepath.Base(path)
			if shouldSkip(file) {
				return nil
			}

			s, err := s.fs.ReadFile(path)
			if err != nil {
				return nil
			}

			switch file {
			// Found name of device
			case "name":
				d.Name = s
			}

			// Sensor names in format "sensor#_foo", e.g. "temp1_input"
			fs := strings.SplitN(file, "_", 2)
			if len(fs) != 2 {
				return nil
			}

			// Gather sensor data into map for later processing
			if _, ok := raw[fs[0]]; !ok {
				raw[fs[0]] = make(map[string]string, 0)
			}

			raw[fs[0]][fs[1]] = s
			return nil
		})
		if err != nil {
			return nil, err
		}

		// Parse all possible sensors from raw data
		sensors, err := parseSensors(raw)
		if err != nil {
			return nil, err
		}

		d.Sensors = sensors
		devices = append(devices, d)
	}

	return devices, nil
}

// detectDevicePaths uses globbing to detect filesystem paths where devices
// may reside on Linux.
func (s *Scanner) detectDevicePaths() ([]string, error) {
	// Locations where device sensors typically reside in /sys on Linux
	globs := []string{
		"/sys/devices/platform/*/name",
		"/sys/devices/platform/*/hwmon/hwmon*/name",
		"/sys/devices/virtual/hwmon/*/name",
	}

	var paths []string
	for _, g := range globs {
		matches, err := s.fs.Glob(g)
		if err != nil {
			return nil, err
		}

		paths = append(paths, matches...)
	}

	return paths, nil
}

// shouldSkip indicates if a given filename should be skipped during the
// filesystem walk operation.
func shouldSkip(file string) bool {
	if strings.HasPrefix(file, "runtime_") {
		return true
	}

	switch file {
	case "async":
	case "autosuspend_delay_ms":
	case "control":
	case "driver_override":
	case "modalias":
	case "uevent":
	default:
		return false
	}

	return true
}

var _ filesystem = &systemFilesystem{}

// A systemFilesystem is a filesystem which uses operations on the host
// filesystem.
type systemFilesystem struct{}

func (fs *systemFilesystem) Glob(pattern string) ([]string, error) {
	return filepath.Glob(pattern)
}

func (fs *systemFilesystem) ReadFile(filename string) (string, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(b)), nil
}

func (fs *systemFilesystem) Walk(root string, walkFn filepath.WalkFunc) error {
	return filepath.Walk(root, walkFn)
}
