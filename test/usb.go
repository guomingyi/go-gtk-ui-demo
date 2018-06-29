package main

import (
	"fmt"
	"usbhid"

)

var mHidDevice *usbhid.HidDevice

func Usb_init() int {
	var ret int = 0

	ret, mHidDevice = Usb_connect()
	if ret == 0 {
		fmt.Printf("Usb_init() success:", mHidDevice)
		return 0;
	}
	fmt.Printf("Usb_init() failed.\n")
	return -1;
}

func Usb_connect() (int, *usbhid.HidDevice) {
	devs := usbhid.HidEnumerate(0, 0)
	for _, dev := range devs {
		if match(&dev) {
			hid_dev, err := dev.Open()
			if err != nil {
				fmt.Printf("hidapi Open failed:\n", dev)
				break
			}
			fmt.Printf("hidapi - Usb_connect -done\n")
			return 0, hid_dev
		}
	}
	fmt.Printf("hidapi Usb_connect failed.\n")

	return -1, nil
}

func match(d *usbhid.HidDeviceInfo) bool {
	vid := d.VendorID
	pid := d.ProductID
	ret := (vid == JS_VID) && (pid == JS_PID)
	return ret
}


func Usb_write(msg string) int {
	var ret int = 0
	
	if msg != "" && len(msg) <= 64 {
		var buffer []byte
		buffer = []byte(msg)
		
		ret,_ = mHidDevice.Write(buffer, false)
		if ret < 0 {
			fmt.Printf("Usb_write: failed\n")
		} else {
			fmt.Printf("Usb_write: %s success\n", string(buffer[:]))
		}
	}
	return ret
}

func Usb_read() (int, string) {
	buffer := make([]byte, 64, 64)
	
	ret, err := mHidDevice.Read(buffer, 1000*1000)
	if err != nil {
		fmt.Printf("Usb_read: err\n")
		return ret, ""
	}
	return ret, string(buffer[:])
}



