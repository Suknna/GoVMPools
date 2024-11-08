package main

import (
	"fmt"
	"govmpools/engine/workstation"
	"log"

	"libvirt.org/go/libvirt"
)

func main() {
	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer conn.Close()
	vms, err := workstation.GetVMList(conn)
	if err != nil {
		log.Fatalf(err.Error())
	}
	for _, v := range vms {
		fmt.Println(v)
	}
}
