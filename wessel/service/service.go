package service

import (
	"log"

	"libvirt.org/go/libvirt"
)

func Run() {
	// 获取libvirt的指针
	libvirtConnect, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		log.Fatalln(err)
	}
	_ = libvirtConnect
}
