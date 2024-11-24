package vm

import (
	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

func GetVMList(c *libvirt.Connect) ([]libvirtxml.Domain, error) {
	// 获取全部domain列表，处于在线的和离线的。返回一个libvirt.Domain列表
	ld, err := c.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE | libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
	if err != nil {
		return nil, err
	}
	// 如果下面任意一行代码出现错误，释放域对象。正在运行的实例保持活动状态。数据结构已释放，此后不应再使用。
	defer func() {
		for _, d := range ld {
			d.Free()
		}
	}()
	// 存储domain列表
	vs := make([]libvirtxml.Domain, len(ld))
	for num, domain := range ld {
		xml, err := domain.GetXMLDesc(libvirt.DOMAIN_XML_SECURE | libvirt.DOMAIN_XML_INACTIVE | libvirt.DOMAIN_XML_MIGRATABLE | libvirt.DOMAIN_XML_UPDATE_CPU)
		if err != nil {
			return nil, err
		}
		domainInfo := libvirtxml.Domain{}
		if err := domainInfo.Unmarshal(xml); err != nil {
			return nil, err
		}
		vs[num] = domainInfo
	}
	return vs, nil
}

// 将传入的domainstate int转换为string
func ParserVMState(VirtualMachineInfoState int) string {
	switch VirtualMachineInfoState {
	case int(libvirt.DOMAIN_NOSTATE):
		return "NOSTATE"
	case int(libvirt.DOMAIN_RUNNING):
		return "RUNNING"
	case int(libvirt.DOMAIN_BLOCKED):
		return "BLOCKED"
	case int(libvirt.DOMAIN_PAUSED):
		return "PAUSED"
	case int(libvirt.DOMAIN_SHUTDOWN):
		return "SHUTDOWN"
	case int(libvirt.DOMAIN_CRASHED):
		return "CRASHED"
	case int(libvirt.DOMAIN_PMSUSPENDED):
		return "PMSUSPENDED"
	case int(libvirt.DOMAIN_SHUTOFF):
		return "SHUTOFF"
	default:
		return "UNKNOWN"
	}
}
