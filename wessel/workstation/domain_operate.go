package workstation

import (
	"fmt"
	"os"

	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

type VMOperateType uint

// 操作相关常量
const (
	VM_START           VMOperateType = iota // 0 启动
	VM_SAFETY_SHUTDOWN                      // 1 安全关机
	VM_FORCE_SHUTDOWN                       // 2 强制关机
	VM_RESTART                              // 3 重启
	VM_UNDEFINE                             // 4 临时删除，只删除了虚拟机的配置信息。数据依旧保留在磁盘中
	VM_DELETE                               // 5 永久删除虚拟机
	VM_SUSPEND                              // 6 将虚拟机置为暂停状态
	VM_DEFINE                               // 7 创建虚拟机
	VM_CREATE                               // 8 临时创建虚拟机
	VM_ATTACH_DEVICE                        // 9 附加驱动器
)

// 虚拟机的常规运维操作，接收一个虚拟机的uuid，操作类型，一个libvirt.Connect指针用于和libvirt进行连接
func VMOperate(uuid string, types VMOperateType, xml string, c *libvirt.Connect) error {
	// 通过libvirt.Connect和uuid获取domain对象的指针
	vm, err := c.LookupDomainByUUIDString(uuid)
	if err != nil {
		return err
	}
	// 函数退出时释放域对象
	defer vm.Free()
	// 根据传入的类型，对domain进行相关操作
	switch types {
	case VM_START:
		if err := vm.Create(); err != nil {
			return err
		}
	case VM_SAFETY_SHUTDOWN:
		if err := vm.Shutdown(); err != nil {
			return err
		}
	case VM_FORCE_SHUTDOWN:
		if err := vm.Destroy(); err != nil {
			return err
		}
	case VM_RESTART:
		if err := vm.Reboot(libvirt.DOMAIN_REBOOT_DEFAULT); err != nil {
			return err
		}
	case VM_UNDEFINE:
		if err := vm.Undefine(); err != nil {
			return err
		}
	case VM_SUSPEND:
		if err := vm.Suspend(); err != nil {
			return err
		}
	case VM_DELETE:
		if err := delete(vm); err != nil {
			return err
		}
	case VM_DEFINE:
		if err := define(xml, c); err != nil {
			return err
		}
	case VM_CREATE:
		if err := create(xml, c); err != nil {
			return err
		}
	case VM_ATTACH_DEVICE:
		// VIR_DOMAIN_AFFECT_LIVE：表示设备更改仅在当前活跃的域实例上进行，并立即生效。
		// VIR_DOMAIN_AFFECT_CONFIG：表示设备更改将持久化到域的配置文件中，以便在下一次启动时仍然有效。
		if err := vm.AttachDeviceFlags(xml, libvirt.DOMAIN_DEVICE_MODIFY_LIVE|libvirt.DOMAIN_DEVICE_MODIFY_CONFIG); err != nil {
			return err
		}
	}
	return fmt.Errorf("unknown operation type")
}

// 从磁盘删除虚拟机，该删除将永久删除虚拟机
func delete(vm *libvirt.Domain) error {
	// 获取虚拟机的快照列表,按照快照树进行返回
	snapshots, err := vm.ListAllSnapshots(libvirt.DOMAIN_SNAPSHOT_LIST_TOPOLOGICAL)
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
	domainXMLStr, err := vm.GetXMLDesc(libvirt.DOMAIN_XML_INACTIVE)
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
	if err := vm.Destroy(); err != nil {
		return err
	}
	// 删除虚拟机配置
	if err := vm.Undefine(); err != nil {
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
