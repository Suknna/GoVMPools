package network

import (
	"fmt"

	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

type Supplier string

const (
	SupplierLocal Supplier = "local"
	SupplierFlat  Supplier = "flat"
	// 暂未实现
	SupplierVlan Supplier = "vlan"
)

type NetOpts struct {
	Name     string
	Supplier Supplier
	NicName  string
	Address  string
	Netmask  string
}

func networkCreate(opts *NetOpts, c *libvirt.Connect) error {
	switch opts.Supplier {
	case SupplierLocal:
		if err := local(c, opts.Name, opts.Address, opts.Netmask); err != nil {
			return fmt.Errorf("local network create failed, err: %s", err.Error())
		}
	case SupplierFlat:
		if err := flat(c, opts.Name, opts.NicName); err != nil {
			return fmt.Errorf("flat network create failed, err: %s", err.Error())
		}
	}
	return nil
}

// 本地网络
func local(c *libvirt.Connect, name string, address string, Netmask string) error {
	n := libvirtxml.Network{
		Name: name,
		Forward: &libvirtxml.NetworkForward{
			Mode: "nat",
		},
		Bridge: &libvirtxml.NetworkBridge{
			Name: name,
		},
		IPs: []libvirtxml.NetworkIP{
			{
				Address: address,
				Netmask: Netmask,
			},
		},
	}
	xml, err := n.Marshal()
	if err != nil {
		return err
	}
	ln, err := c.NetworkDefineXML(xml)
	if err != nil {
		return err
	}
	if err := ln.Create(); err != nil {
		return err
	}
	return nil
}

// 直连物理网卡
func flat(c *libvirt.Connect, name string, nicName string) error {
	n := libvirtxml.Network{
		Name: name,
		Forward: &libvirtxml.NetworkForward{
			Mode: "bridge",
			Interfaces: []libvirtxml.NetworkForwardInterface{
				{
					Dev: nicName,
				},
			},
		},
	}
	xml, err := n.Marshal()
	if err != nil {
		return err
	}
	ln, err := c.NetworkDefineXML(xml)
	if err != nil {
		return err
	}
	if err := ln.Create(); err != nil {
		return err
	}
	return nil
}

func saveNetworkCfg(path string)
