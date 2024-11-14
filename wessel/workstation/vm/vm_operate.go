package vm

import (
	"fmt"
	"os"

	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

type VirtualMachineOperateType uint

// 操作相关常量
const (
	VIRTUAL_MACHINE_START           VirtualMachineOperateType = iota // 0 启动
	VIRTUAL_MACHINE_SAFETY_SHUTDOWN                                  // 1 安全关机
	VIRTUAL_MACHINE_FORCE_SHUTDOWN                                   // 2 强制关机
	VIRTUAL_MACHINE_RESTART                                          // 3 重启
	VIRTUAL_MACHINE_UNDEFINE                                         // 4 临时删除，只删除了虚拟机的配置信息。数据依旧保留在磁盘中
	VIRTUAL_MACHINE_DELETE                                           // 5 永久删除虚拟机
	VIRTUAL_MACHINE_SUSPEND                                          // 6 将虚拟机置为暂停状态
	VIRTUAL_MACHINE_DEFINE                                           // 7 创建虚拟机
	VIRTUAL_MACHINE_CREATE                                           // 8 临时创建虚拟机
	VIRTUAL_MACHINE_ATTACH_DEVICE                                    // 9 附加驱动器
)

// 虚拟机的常规运维操作，接收一个虚拟机的uuid，操作类型，一个libvirt.Connect指针用于和libvirt进行连接
func VirtualMachineOperate(uuid string, types VirtualMachineOperateType, xml string, c *libvirt.Connect) error {
	// 通过libvirt.Connect和uuid获取domain对象的指针
	VirtualMachine, err := c.LookupDomainByUUIDString(uuid)
	if err != nil {
		return err
	}
	// 函数退出时释放域对象
	defer VirtualMachine.Free()
	// 根据传入的类型，对domain进行相关操作
	switch types {
	case VIRTUAL_MACHINE_START:
		if err := VirtualMachine.Create(); err != nil {
			return err
		}
	case VIRTUAL_MACHINE_SAFETY_SHUTDOWN:
		if err := VirtualMachine.Shutdown(); err != nil {
			return err
		}
	case VIRTUAL_MACHINE_FORCE_SHUTDOWN:
		if err := VirtualMachine.Destroy(); err != nil {
			return err
		}
	case VIRTUAL_MACHINE_RESTART:
		if err := VirtualMachine.Reboot(libvirt.DOMAIN_REBOOT_DEFAULT); err != nil {
			return err
		}
	case VIRTUAL_MACHINE_UNDEFINE:
		if err := VirtualMachine.Undefine(); err != nil {
			return err
		}
	case VIRTUAL_MACHINE_SUSPEND:
		if err := VirtualMachine.Suspend(); err != nil {
			return err
		}
	case VIRTUAL_MACHINE_DELETE:
		if err := delete(VirtualMachine); err != nil {
			return err
		}
	case VIRTUAL_MACHINE_DEFINE:
		if err := define(xml, c); err != nil {
			return err
		}
	case VIRTUAL_MACHINE_CREATE:
		if err := create(xml, c); err != nil {
			return err
		}
	case VIRTUAL_MACHINE_ATTACH_DEVICE:
		// VIR_DOMAIN_AFFECT_LIVE：表示设备更改仅在当前活跃的域实例上进行，并立即生效。
		// VIR_DOMAIN_AFFECT_CONFIG：表示设备更改将持久化到域的配置文件中，以便在下一次启动时仍然有效。
		if err := VirtualMachine.AttachDeviceFlags(xml, libvirt.DOMAIN_DEVICE_MODIFY_LIVE|libvirt.DOMAIN_DEVICE_MODIFY_CONFIG); err != nil {
			return err
		}
	}
	return fmt.Errorf("unknown operation type")
}

// 从磁盘删除虚拟机，该删除将永久删除虚拟机
func delete(VirtualMachine *libvirt.Domain) error {
	// 获取虚拟机的快照列表,按照快照树进行返回
	snapshots, err := VirtualMachine.ListAllSnapshots(libvirt.DOMAIN_SNAPSHOT_LIST_TOPOLOGICAL)
	if err != nil {
		return nil
	}
	defer func() {
		for _, snapshot := range snapshots {
			snapshot.Free()
		}
	}()
	// 判断虚拟机是否存在快照
	if len(snapshots) != 0 {
		// 按照顺序删除快照
		for _, snapshot := range snapshots {
			err := snapshot.Delete(libvirt.DOMAIN_SNAPSHOT_DELETE_CHILDREN)
			if err != nil {
				// 释放内存
				snapshot.Free()
				return err
			}
		}
	}
	// 获取虚拟机的xml文件
	domainXMLStr, err := VirtualMachine.GetXMLDesc(libvirt.DOMAIN_XML_INACTIVE)
	if err != nil {
		return err
	}
	// 调用libvirt-xml库解析xml文件
	var d libvirtxml.Domain
	if err := d.Unmarshal(domainXMLStr); err != nil {
		return err
	}
	diskPaths := make([]string, len(d.Devices.Disks))
	// 获取磁盘相关配置
	for i, diskinfo := range d.Devices.Disks {
		if diskinfo.Device == "cdrom" {
			continue
		}
		diskPaths[i] = diskinfo.Source.File.File
	}
	// 对虚拟机进行关机
	if err := VirtualMachine.Destroy(); err != nil {
		return err
	}
	// 删除虚拟机配置
	if err := VirtualMachine.Undefine(); err != nil {
		return err
	}
	// 根据路径删除虚拟机磁盘
	for _, v := range diskPaths {
		if err := os.Remove(v); err != nil {
			return err
		}
	}
	return nil
}

// 永久创建虚拟机
func define(xml string, c *libvirt.Connect) error {
	// 通过xml文件生成domain指针
	d, err := c.DomainDefineXML(xml)
	if err != nil {
		return err
	}
	defer d.Free()
	// 执行创建操作
	if err := d.Create(); err != nil {
		return err
	}
	return nil
}

func create(xml string, c *libvirt.Connect) error {
	// 通过xml文件生成domain指针
	d, err := c.DomainCreateXML(xml, libvirt.DOMAIN_NONE)
	defer d.Free()
	if err != nil {
		return err
	}
	// 执行创建操作
	if err := d.Create(); err != nil {
		return err
	}
	return nil
}
