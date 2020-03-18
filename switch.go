package delcom_switch_go

import (
	"errors"
	"github.com/karalabe/hid"
)

var (
	ErrNoDelcomDevices = errors.New("no delcom devices")
	ErrNoDelcomSwitch  = errors.New("no delcom switch")
	ErrNoSwitchStatus  = errors.New("can't recognize switch status")
)

const (
	DelcomVendorID  uint16 = 0x0FC5
	DelcomProductID uint16 = 0xB080
	ProductType            = "USB FS IO"

	maxPacketLen = 16
	readDataFlag = 100 // mandatory

	buttonPressedFlag  = 254
	buttonReleasedFlag = 255
)

type DelcomSwitch struct {
	dev *hid.Device
}

func NewSwitch() (*DelcomSwitch, error) {
	var info hid.DeviceInfo

	if !hid.Supported() {
		return nil, hid.ErrUnsupportedPlatform
	}

	delkomInfos := hid.Enumerate(DelcomVendorID, DelcomProductID)
	if len(delkomInfos) == 0 {
		return nil, ErrNoDelcomDevices
	}
	for _, di := range delkomInfos {
		if di.Product == ProductType {
			info = di
			break
		}
	}

	if info == (hid.DeviceInfo{}) {
		return nil, ErrNoDelcomSwitch
	}

	dev, err := info.Open()
	if err != nil {
		return nil, err
	}

	return &DelcomSwitch{
		dev: dev,
	}, nil
}

func (btn *DelcomSwitch) IsPressed() (bool, error) {
	data := make([]byte, maxPacketLen)
	data[0] = readDataFlag

	_, err := btn.dev.GetFeatureReport(data)
	if err != nil {
		return false, err
	}

	return btn.ackReport(data[0])
}

func (btn *DelcomSwitch) ackReport(flag byte) (bool, error) {
	switch flag {
	case buttonReleasedFlag:
		return false, nil
	case buttonPressedFlag:
		return true, nil
	default:
		return false, ErrNoSwitchStatus
	}
}

func (btn *DelcomSwitch) Close() error {
	return btn.dev.Close()
}
