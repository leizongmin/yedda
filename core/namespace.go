package core

import (
	"time"
)

type Namespace struct {
	DbNum uint32
	Data  map[time.Duration]map[string]*Table `json:"data"`
	Queue []NamespaceResetItem                `json:"queue"`
}

type NamespaceResetItem struct {
	Name     string
	Duration time.Duration
	Expires  time.Time
}

func NewNamespace(db uint32) *Namespace {
	return &Namespace{
		DbNum: db,
		Data:  make(map[time.Duration]map[string]*Table),
		Queue: make([]NamespaceResetItem, 0),
	}
}

func (n *Namespace) addToExpiresQueue(name string, duration time.Duration) {
	v := NamespaceResetItem{
		Name:     name,
		Duration: duration,
		Expires:  time.Now().Add(duration),
	}
	n.Queue = append(n.Queue, v)
	//fmt.Println("addToExpiresQueue", n.DbNum, name, duration, v.Expires)
}

/// 获取指定 name & duration 的数据表，如果不存在则先初始化
func (n *Namespace) Get(name string, duration time.Duration) *Table {
	if _, exists := n.Data[duration]; !exists {
		// 如果指定 duration 不存在任何数据，先初始化
		n.Data[duration] = make(map[string]*Table)
	}
	if _, exists := n.Data[duration][name]; !exists {
		// 如果指定 duration 和 name 不存在任何数据，先初始化
		// 并且添加到队列，在 duration 过期后会有专门的程序去将整个表删除
		n.Data[duration][name] = NewTable(name)
		n.addToExpiresQueue(name, duration)
	}
	return n.Data[duration][name]
}

/// 删除指定 name & duration 的数据表
func (n *Namespace) Delete(name string, duration time.Duration) {
	if _, exists := n.Data[duration]; exists {
		delete(n.Data[duration], name)
		// 如果该 duration 已没有任何 namespace，则直接将其删除
		if len(n.Data[duration]) < 1 {
			delete(n.Data, duration)
		}
	}
	//fmt.Printf("%x #%d Delete %s %s %v\n", &n, n.DbNum, name, duration, n.Data)
}

/// 删除已经过期的 namespace
func (n *Namespace) DeleteExpired(t time.Time) {
	queue := make([]NamespaceResetItem, 0)
	// 处理已经过期的数据
	for _, v := range n.Queue {
		if v.Expires.Sub(t) > 0 {
			queue = append(queue, v)
		} else {
			n.Delete(v.Name, v.Duration)
		}
	}
	// 清理为空的 duration namespace 数据表
	for d := range n.Data {
		if len(n.Data[d]) < 1 {
			delete(n.Data, d)
		}
	}
	n.Queue = queue
	//fmt.Printf("%x #%d DeleteExpired %v\n", &n, n.DbNum, n.Data)
}

/// 销毁
func (n *Namespace) Destroy() {
	for d := range n.Data {
		for k := range n.Data[d] {
			n.Data[d][k].Destroy()
		}
	}
	n.Data = nil
}
