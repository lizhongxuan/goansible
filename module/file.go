package module

// FileModule file模块的配置参数
type FileModule struct {
	Path          string `json:"path"`           // 文件或目录的路径
	State         string `json:"state"`          // 期望的状态
	Mode          string `json:"mode"`           // 权限模式，八进制数
	Owner         string `json:"owner"`          // 文件或目录的所有者
	Group         string `json:"group"`          // 文件或目录的组
	Setype        string `json:"setype"`         // SELinux的类型
	Selevel       string `json:"selevel"`        // SELinux的级别
	Recurse       bool   `json:"recurse"`        // 是否递归应用设置
	Follow        bool   `json:"follow"`         // 是否跟随符号链接
	Force         bool   `json:"force"`          // 是否强制创建链接
	Src           string `json:"src"`            // 硬链接或符号链接的源路径
	Dest          string `json:"dest"`           // 硬链接或符号链接的目标路径
	Attributes    string `json:"attributes"`     // 文件属性
	Ctime         string `json:"ctime"`          // 创建时间
	Mtime         string `json:"mtime"`          // 修改时间
	Atime         string `json:"atime"`          // 访问时间
	Checksum      string `json:"checksum"`       // 文件校验和
	Backup        bool   `json:"backup"`         // 是否备份原始文件
	DirectoryMode string `json:"directory_mode"` // 子目录的权限模式
	Umask         string `json:"umask"`          // 文件和目录创建时的默认权限掩码
}

type FileModuleState string

const (
	AbsentFileModuleState    FileModuleState = "absent"    // 确保文件或目录不存在
	DirectoryFileModuleState FileModuleState = "directory" // 确保路径是一个目录
	FileFileModuleState      FileModuleState = "file"      // 确保路径是一个文件
	HardFileModuleState      FileModuleState = "hard"      // 创建硬链接
	LinkFileModuleState      FileModuleState = "link"      // 创建符号链接（默认）
	MountedFileModuleState   FileModuleState = "mounted"   // 确保文件系统已挂载
	TouchFileModuleState     FileModuleState = "touch"     // 创建文件，如果文件已存在，则更新其时间戳
)

func (m *FileModule) Show() string {
	return "File Module"
}
