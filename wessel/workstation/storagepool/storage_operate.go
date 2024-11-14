package storagepool

import (
	"libvirt.org/go/libvirt"
)

type StoragePoolOperateType uint

const (
	STORAGE_POOL_CREATE         StoragePoolOperateType = iota // 0 创建
	STORAGE_POOL_DELETE_NORMAL                                // 1 只删除存储池的配置
	STORAGE_POOL_DELETE_ZEROED                                // 2 将存储池中所有数据都清零
	STORAGE_POOL_IN_ACTIVE                                    // 3 将存储池置为活跃的
	STORAGE_POOL_IS_INACTIVE                                  // 4 将存储池置为不活跃的
	STORAGE_POOL_SET_AUTO_START                               // 6 将存储池设置开机自启
)

func StoragePoolOperate(uuid string, types StoragePoolOperateType, xml string, c *libvirt.Connect) error {
	sp, err := c.LookupStoragePoolByUUIDString(uuid)
	if err != nil {
		return err
	}
	defer sp.Free()
	switch types {
	case STORAGE_POOL_CREATE:
		if err := create(xml, c); err != nil {
			return err
		}
	case STORAGE_POOL_DELETE_NORMAL:
		if err := sp.Delete(libvirt.STORAGE_POOL_DELETE_NORMAL); err != nil {
			return err
		}
	case STORAGE_POOL_DELETE_ZEROED:
		if err := sp.Delete(libvirt.STORAGE_POOL_DELETE_ZEROED); err != nil {
			return err
		}
	case STORAGE_POOL_IS_INACTIVE:
		if err := sp.Destroy(); err != nil {
			return err
		}
	case STORAGE_POOL_IN_ACTIVE:
		if err := sp.Create(libvirt.STORAGE_POOL_CREATE_NORMAL); err != nil {
			return err
		}
	case STORAGE_POOL_SET_AUTO_START:
		if err := sp.SetAutostart(true); err != nil {
			return err
		}
	}
	return nil
}

func create(xml string, c *libvirt.Connect) error {
	sp, err := c.StoragePoolDefineXML(xml, libvirt.STORAGE_POOL_DEFINE_VALIDATE)
	if err != nil {
		return err
	}
	defer sp.Free()
	if err := sp.Create(libvirt.STORAGE_POOL_CREATE_WITH_BUILD_OVERWRITE); err != nil {
		return err
	}
	return nil
}
