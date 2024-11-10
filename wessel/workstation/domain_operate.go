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
	VM_DELETE
	VM_SUSPEND // 5 将虚拟机置为暂停状态
)

// 虚拟机的常规运维操作，接收一个虚拟机的uuid，操作类型，一个libvirt.Connect指针用于和libvirt进行连接
func VMOperate(uuid string, types VMOperateType, c *libvirt.Connect) error {
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
		if err := vmDelete(vm); err != nil {
			return err
		}
	}
	return fmt.Errorf("unknown operation type")
}

// 从磁盘删除虚拟机，该删除将永久删除虚拟机
func vmDelete(vm *libvirt.Domain) error {
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
