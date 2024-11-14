package vm

import (
	"os"

	"libvirt.org/go/libvirt"
)

// domain 列表返回值
type VirtualMachineInfo struct {
	VirtualMachineID        string
	VirtualMachineName      string
	HypervisorHostName      string
	VirtualMachineStatus    string
	VirtualMachineImage     string
	VirtualMachineVcpus     uint
	VirtualMachineMemoryUse uint64
}

/*
获取domain 的列表包括如下字段
1. uuid
2. 名称
3. 所属节点
4. 当前状态
6. 镜像版本
*/
func GetVMList(c *libvirt.Connect) ([]VirtualMachineInfo, error) {
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
	// 获取当前主机名
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	// 存储domain列表
	vs := make([]VirtualMachineInfo, len(ld))
	for num, domain := range ld {
		// 声明单一的vm对象
		var v VirtualMachineInfo
		// 获取uuid
		v.VirtualMachineID, err = domain.GetUUIDString()
		if err != nil {
			return nil, err
		}
		// 获取名称
		v.VirtualMachineName, err = domain.GetName()
		if err != nil {
			return nil, err
		}
		// 获取节点名称
		v.HypervisorHostName = hostname
		// 获取当前状态
		domainStateInt, _, err := domain.GetState()
		if err != nil {
			return nil, err
		}
		v.VirtualMachineStatus = ParserVMState(int(domainStateInt))
		// 获取当前虚拟机的镜像类型
		v.VirtualMachineImage, err = domain.GetOSType()
		if err != nil {
			return nil, err
		}
		// 获取虚拟机分配的最大cpu
		v.VirtualMachineVcpus, err = domain.GetMaxVcpus()
		if err != nil {
			return nil, err
		}
		// 获取虚拟机分配的最大内存
		v.VirtualMachineMemoryUse, err = domain.GetMaxMemory()
		if err != nil {
			return nil, err
		}
		vs[num] = v
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
