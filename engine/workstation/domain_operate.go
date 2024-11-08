package workstation

import "libvirt.org/go/libvirt"

// 创建虚拟机所需的参数
type DomainOpts struct {
}

// 启动domain
func DomainStart(id string, c *libvirt.Connect) error {
	return nil
}

// 关闭domain
func DomainStop(id string, c *libvirt.Connect) error {
	return nil
}

// 重启domain
func DomainRestart(id string, c *libvirt.Connect) error {
	return nil
}

// 删除domain
func DomainDelete(id string, c *libvirt.Connect) error {
	return nil
}

// 暂停domain
func DomainPaused(id string, c *libvirt.Connect) error {
	return nil
}

// 创建domain
func DomainCreate(opts *DomainOpts, c *libvirt.Connect) error {
	return nil
}
