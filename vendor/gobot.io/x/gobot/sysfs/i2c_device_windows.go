package sysfs

import "errors"

func (d *i2cDevice) queryFunctionality() (err error) {
	return errors.New("Not implemented on windows")
}

func (d *i2cDevice) SetAddress(address int) (err error) {
	return errors.New("Not implemented on windows")
}

func (d *i2cDevice) smbusAccess(readWrite byte, command byte, size uint32, data uintptr) error {
	return errors.New("Not implemented on windows")
}
