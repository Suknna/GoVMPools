package vm

type VMInfo struct {
	UUID      string
	Name      string
	State     string
	OSType    string
	CPUs      int
	Mem       string
	AutoStart bool
}

type VMRealTimeData struct {
	CpuUsed float32
	MemUsed float32
}

type VMCreateOps struct {
}
