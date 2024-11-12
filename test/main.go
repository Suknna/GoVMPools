package main

import (
	"fmt"
	"log"

	"libvirt.org/go/libvirt"
)

func main() {
	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer conn.Close()
	storages, err := conn.ListAllStoragePools(libvirt.CONNECT_LIST_STORAGE_POOLS_ACTIVE | libvirt.CONNECT_LIST_STORAGE_POOLS_AUTOSTART)
	if err != nil {
		log.Fatalf(err.Error())
	}
	for _, v := range storages {
		xml, err := v.GetXMLDesc(libvirt.STORAGE_XML_INACTIVE)
		if err != nil {
			log.Fatalf(err.Error())
		}
		fmt.Println(xml)
	}
}
