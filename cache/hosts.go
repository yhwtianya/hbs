package cache

import (
	"sync"

	"github.com/open-falcon/common/model"
	"github.com/open-falcon/hbs/db"
)

// 每次心跳的时候agent把hostname汇报上来，经常要知道这个机器的hostid，把此信息缓存
// key: hostname value: hostid
type SafeHostMap struct {
	sync.RWMutex
	M map[string]int
}

// 保存所有hostname和host id对应关系
var HostMap = &SafeHostMap{M: make(map[string]int)}

// 根据hostname获取host id
func (this *SafeHostMap) GetID(hostname string) (int, bool) {
	this.RLock()
	defer this.RUnlock()
	id, exists := this.M[hostname]
	return id, exists
}

// 从数据库获取所有hostname和host id对应关系
func (this *SafeHostMap) Init() {
	m, err := db.QueryHosts()
	if err != nil {
		return
	}

	this.Lock()
	defer this.Unlock()
	this.M = m
}

type SafeMonitoredHosts struct {
	sync.RWMutex
	M map[int]*model.Host
}

// 保存所有非维护状态的主机信息
var MonitoredHosts = &SafeMonitoredHosts{M: make(map[int]*model.Host)}

// 获取所有非维护状态的主机信息
func (this *SafeMonitoredHosts) Get() map[int]*model.Host {
	this.RLock()
	defer this.RUnlock()
	return this.M
}

// 从数据库获取所有非维护状态的主机信息
func (this *SafeMonitoredHosts) Init() {
	m, err := db.QueryMonitoredHosts()
	if err != nil {
		return
	}

	this.Lock()
	defer this.Unlock()
	this.M = m
}
