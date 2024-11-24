package network

import (
	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

func GetNetworkList(c *libvirt.Connect) ([]libvirtxml.Network, error) {
	nets, err := c.ListAllNetworks(libvirt.CONNECT_LIST_NETWORKS_ACTIVE | libvirt.CONNECT_LIST_NETWORKS_INACTIVE)
	if err != nil {
		return nil, err
	}
	defer func() {
		for _, n := range nets {
			n.Free()
		}
	}()
	netInfos := make([]libvirtxml.Network, len(nets))
	for num, n := range nets {
		xml, err := n.GetXMLDesc(libvirt.NETWORK_XML_INACTIVE)
		if err != nil {
			return nil, err
		}
		netInfo := libvirtxml.Network{}
		if err := netInfo.Unmarshal(xml); err != nil {
			return nil, err
		}
		netInfos[num] = netInfo
	}
	return netInfos, nil
}
