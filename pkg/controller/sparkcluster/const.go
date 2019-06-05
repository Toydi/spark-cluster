package sparkcluster

// Const Variable
const (
	Master    = "master"
	MasterPvc = "namenode-pvc"
	// MasterImage      = "registry.njuics.cn/wdongyu/hive_on_kube:0.2"
	// "registry.njuics.cn/wdongyu/spark_master_on_kube:1.0.2"
	// SlaveImage       = "registry.njuics.cn/wdongyu/spark_slave_on_kube:0.2"
	// MasterImage      = "registry.njuics.cn/qr/spark-master:0.1"
	// SlaveImage       = "registry.njuics.cn/qr/spark-slave:0.1"
	// VscodeImage      = "registry.njuics.cn/dlkit/code-server"
	MasterImage      = "registry.njuics.cn/wdongyu/spark_master_on_kube:1.0.3"
	SlaveImage       = "registry.njuics.cn/wdongyu/spark_slave_on_kube:1.0.3"
	VscodeImage      = "registry.njuics.cn/qr/code-server:1.0"
	Slave            = "slave"
	SlavePvc         = "datanode-pvc"
	ShareServer      = "114.212.189.141"
	StorageClassName = "cephfs"
	DefaultGitRepo   = "https://github.com/Toydi/Test"
)
