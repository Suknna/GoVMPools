package main

import (
	"fmt"
	"log"

	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

func main() {
	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer conn.Close()
	vm, err := conn.LookupDomainByUUIDString("d5de9a0c-170b-47e4-b867-5d94713f1c99")
	if err != nil {
		log.Fatalf(err.Error())
	}
	domainXMLStr, err := vm.GetXMLDesc(libvirt.DOMAIN_XML_INACTIVE)
	if err != nil {
		log.Fatalf(err.Error())
	}
	// 调用libvirt-xml库解析xml文件
	var d libvirtxml.Domain
	if err := d.Unmarshal(domainXMLStr); err != nil {
		log.Fatalf(err.Error())
	}
	diskPaths := make([]string, len(d.Devices.Disks))
	// 获取磁盘相关配置
	for i, diskinfo := range d.Devices.Disks {
		if diskinfo.Device == "cdrom" {
			continue
		}
		diskPaths[i] = diskinfo.Source.File.File
	}
	fmt.Println(diskPaths)
}
