package network

import "libvirt.org/go/libvirt"

type NetworkOperateType uint

const (
	NETWORK_CREATE NetworkOperateType = iota // 0 创建
	NETWORK_DELETE                           // 1 删除
	NETWORK_EDIT                             // 2 修改, 暂不实现
)

func NetworkOperate(uuid string, types NetworkOperateType, opts NetOpts, c *libvirt.Connect) error {
	n, err := c.LookupNetworkByUUIDString(uuid)
	if err != nil {
		return err
	}
	switch types {
	case NETWORK_CREATE:
		if err := networkCreate(&opts, c); err != nil {
			return err
		}
	case NETWORK_DELETE:
		if err := n.Undefine(); err != nil {
			return err
		}
	}
	return nil
}
