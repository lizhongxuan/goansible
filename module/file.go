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

//// Group 组
//type group struct {
//	Name      string
//	Hosts     map[string]*model.Host
//	Parents   []string
//	GroupVars map[string]interface{}
//	Depth     int
//}
//
//const (
//	varsOpt int = iota
//	childrenOpt
//	groupOpt
//)
//
//type fileInv struct {
//	path string
//}
//
//func NewFile(path string) (model.Inventory, error) {
//	return &fileInv{path: path}, nil
//}
//
//func (g *fileInv) Vars() (map[string]interface{}, error) {
//	return make(map[string]interface{}), nil
//}
//
//// GetGroup get group struct by inventory file
//func (g *fileInv) Groups() (map[string]*model.Group, error) {
//	groupMap := make(map[string]*group)
//	resMap := make(map[string]*model.Group)
//	groupName := "all"
//	groupMap[groupName] = &group{
//		Name:      groupName,
//		Hosts:     make(map[string]*model.Host),
//		Parents:   make([]string, 0),
//		GroupVars: make(map[string]interface{}),
//	}
//	opt := groupOpt
//	err := readLine(g.path, func(s string) {
//		if strings.Contains(s, "#") || s == "" {
//			return
//		}
//
//		// group vars
//		gvarRegexp := regexp.MustCompile(`\[(.*):vars\]`)
//		gvarParams := gvarRegexp.FindStringSubmatch(s)
//		if len(gvarParams) == 2 {
//			opt = varsOpt
//			groupName = gvarParams[1]
//			if _, ok := groupMap[groupName]; !ok {
//				groupMap[groupName] = &group{
//					Name:      groupName,
//					Hosts:     make(map[string]*model.Host),
//					Parents:   make([]string, 0),
//					GroupVars: make(map[string]interface{}),
//				}
//			}
//			return
//		}
//
//		// group children
//		gchildRegexp := regexp.MustCompile(`\[(.*):children\]`)
//		gchildParams := gchildRegexp.FindStringSubmatch(s)
//		if len(gchildParams) == 2 {
//			opt = childrenOpt
//			groupName = gchildParams[1]
//			if _, ok := groupMap[groupName]; !ok {
//				groupMap[groupName] = &group{
//					Name:      groupName,
//					Hosts:     make(map[string]*model.Host),
//					Parents:   make([]string, 0),
//					GroupVars: make(map[string]interface{}),
//				}
//			}
//			return
//		}
//
//		// group
//		groupRegexp := regexp.MustCompile(`\[(.*)\]`)
//		groupParams := groupRegexp.FindStringSubmatch(s)
//		if len(groupParams) == 2 {
//			opt = groupOpt
//			groupName = groupParams[1]
//			if _, ok := groupMap[groupName]; !ok {
//				groupMap[groupName] = &group{
//					Name:      groupName,
//					Hosts:     make(map[string]*model.Host),
//					Parents:   make([]string, 0),
//					GroupVars: make(map[string]interface{}),
//				}
//			}
//			return
//		}
//
//		switch opt {
//		case childrenOpt:
//			ps := groupMap[groupName].Parents
//			ps = append(ps, strings.TrimSpace(s))
//			groupMap[groupName].Parents = ps
//			return
//		case varsOpt:
//			fileds := strings.Split(s, " ")
//			for _, filed := range fileds {
//				if filed == "" {
//					continue
//				}
//				fRegexp := regexp.MustCompile(`(.*)=(.*)`)
//				fParams := fRegexp.FindStringSubmatch(filed)
//				if len(fParams) == 3 {
//					groupMap[groupName].GroupVars[fParams[1]] = convertVar(fParams[2])
//				}
//			}
//		case groupOpt:
//			fileds := strings.Split(s, " ")
//			hostname := fileds[0]
//			if _, ok := groupMap[groupName].Hosts[hostname]; !ok {
//				groupMap[groupName].Hosts[hostname] = &model.Host{
//					Name:     hostname,
//					HostVars: make(map[string]interface{}),
//				}
//			}
//			for _, filed := range fileds {
//				if filed == "" {
//					continue
//				}
//				fRegexp := regexp.MustCompile(`(.*)=(.*)`)
//				fParams := fRegexp.FindStringSubmatch(filed)
//				if len(fParams) == 3 {
//					groupMap[groupName].Hosts[hostname].HostVars[fParams[1]] = convertVar(fParams[2])
//				}
//			}
//		}
//
//	})
//
//	if err != nil {
//		return resMap, err
//	}
//
//	groupSlice := make([]*group, 0)
//	for name, g := range groupMap {
//		groupMap[name].Depth = depth(0, g.Parents, groupMap)
//		groupSlice = append(groupSlice, groupMap[name])
//	}
//
//	sort.Slice(groupSlice, func(i, j int) bool {
//		return groupSlice[i].Depth < groupSlice[j].Depth
//	})
//
//	for _, g := range groupSlice {
//		name := g.Name
//		for _, p := range g.Parents {
//			for k, v := range groupMap[p].Hosts {
//				groupMap[name].Hosts[k] = v
//			}
//		}
//		for hn := range groupMap[name].Hosts {
//			for alln, allh := range groupMap["all"].Hosts {
//				if alln != hn {
//					continue
//				}
//				for allk, allv := range allh.HostVars {
//					groupMap[name].Hosts[hn].HostVars[allk] = allv
//				}
//			}
//			groupMap[name].Hosts[hn].HostVars["ansible_hostname"] = hn
//			for gk, gv := range g.GroupVars {
//				groupMap[name].Hosts[hn].HostVars[gk] = gv
//			}
//
//		}
//		resMap[name] = &model.Group{
//			Name:  name,
//			Hosts: groupMap[name].Hosts,
//		}
//	}
//	return resMap, nil
//
//}
//
//func depth(n int, ps []string, groupMap map[string]*group) int {
//	if len(ps) > 0 {
//		n += 1
//	}
//	ds := []int{}
//	for _, p := range ps {
//		d := depth(n, groupMap[p].Parents, groupMap)
//		ds = append(ds, d)
//	}
//	max := n
//	for _, d := range ds {
//		if d > max {
//			max = d
//		}
//	}
//	return max
//}
//
//func readLine(fileName string, handler func(string)) error {
//	f, err := os.Open(fileName)
//	if err != nil {
//		return err
//	}
//	buf := bufio.NewReader(f)
//	for {
//		line, err := buf.ReadString('\n')
//		line = strings.TrimSpace(line)
//		handler(line)
//		if err != nil {
//			if err == io.EOF {
//				return nil
//			}
//			return err
//		}
//	}
//}
//
//func convertVar(src string) interface{} {
//	if src == "yes" || src == "true" {
//		return true
//	}
//	if src == "no" || src == "false" {
//		return false
//	}
//	return src
//}
