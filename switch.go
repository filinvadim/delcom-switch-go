package delcom_switch_go

import (
	"errors"
	"github.com/karalabe/hid"
	"runtime"
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
	var (
		info hid.DeviceInfo
		product  = ProductType
	)

	// FIXME karalabe/hid lib marshalling problem: linux hid devices doesn't have 'Product' field because it's called 'iProduct'
	if runtime.GOOS == "linux" {
		product = ""
	}

	// also will cause if CGO isn't enabled or app wasn't properly cross-compiled
	if !hid.Supported() {
		return nil, hid.ErrUnsupportedPlatform
	}

	delkomInfos := hid.Enumerate(DelcomVendorID, DelcomProductID)
	if len(delkomInfos) == 0 {
		return nil, ErrNoDelcomDevices
	}

	for _, di := range delkomInfos {
		if di.Product == product {
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
	if btn.dev == nil {
		return false, ErrNoDelcomSwitch
	}
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
	if btn.dev == nil {
		return ErrNoDelcomSwitch
	}
	return btn.dev.Close()
}
