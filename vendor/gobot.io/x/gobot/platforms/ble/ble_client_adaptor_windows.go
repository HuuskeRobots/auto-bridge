package ble

import (
	"errors"

	blelib "github.com/currantlabs/ble"
)

func defaultDevice(impl string) (d blelib.Device, err error) {
	return nil, errors.New("Not supported on windows")
}
