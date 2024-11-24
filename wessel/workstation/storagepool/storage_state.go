package storagepool

import (
	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

// 实现存储相关操作
// 获取当前系统存储池的状态

/*
存储池类型：
VIR_CONNECT_LIST_STORAGE_POOLS_DIR			目录
VIR_CONNECT_LIST_STORAGE_POOLS_FS			文件系统
VIR_CONNECT_LIST_STORAGE_POOLS_NETFS		网络文件系统
VIR_CONNECT_LIST_STORAGE_POOLS_LOGICAL		逻辑卷
VIR_CONNECT_LIST_STORAGE_POOLS_DISK			磁盘
VIR_CONNECT_LIST_STORAGE_POOLS_ISCSI		iSCSI
VIR_CONNECT_LIST_STORAGE_POOLS_SCSI			SCSI
VIR_CONNECT_LIST_STORAGE_POOLS_MPATH		多路径
VIR_CONNECT_LIST_STORAGE_POOLS_RBDRADOS		块设备
VIR_CONNECT_LIST_STORAGE_POOLS_SHEEPDOG		Sheepdog
VIR_CONNECT_LIST_STORAGE_POOLS_GLUSTER		GlusterFS
VIR_CONNECT_LIST_STORAGE_POOLS_ZFS			ZFS
VIR_CONNECT_LIST_STORAGE_POOLS_VSTORAGE		虚拟存储
VIR_CONNECT_LIST_STORAGE_POOLS_ISCSI_DIRECT	直接iSCSI
*/

func GetStorageList(c *libvirt.Connect) ([]libvirtxml.StoragePool, error) {
	storagePools, err := c.ListAllStoragePools(libvirt.CONNECT_LIST_STORAGE_POOLS_ACTIVE | libvirt.CONNECT_LIST_STORAGE_POOLS_AUTOSTART)
	if err != nil {
		return nil, err
	}
	defer func() {
		for _, s := range storagePools {
			s.Free()
		}
	}()
	spis := make([]libvirtxml.StoragePool, len(storagePools))
	// 获取storage的配置信息
	for num, storagePool := range storagePools {
		// 获取xml内容
		xml, err := storagePool.GetXMLDesc(libvirt.STORAGE_XML_INACTIVE)
		if err != nil {
			return nil, err
		}
		// 解析xml内容
		storagepool := libvirtxml.StoragePool{}
		if err := storagepool.Unmarshal(xml); err != nil {
			return nil, err
		}
		spis[num] = storagepool
	}
	return spis, nil
}
