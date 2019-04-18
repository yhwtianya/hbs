package cache

import (
	"sync"

	"github.com/open-falcon/common/model"
	"github.com/open-falcon/hbs/db"
)

// 一个HostGroup对应多个Template
type SafeGroupTemplates struct {
	sync.RWMutex
	M map[int][]int
}

// 保存所有group id和template id对应关系
var GroupTemplates = &SafeGroupTemplates{M: make(map[int][]int)}

// 根据group id获取template id
func (this *SafeGroupTemplates) GetTemplateIds(gid int) ([]int, bool) {
	this.RLock()
	defer this.RUnlock()
	templateIds, exists := this.M[gid]
	return templateIds, exists
}

// 从数据库获取所有group id和template id对应关系
func (this *SafeGroupTemplates) Init() {
	m, err := db.QueryGroupTemplates()
	if err != nil {
		return
	}

	this.Lock()
	defer this.Unlock()
	this.M = m
}

type SafeTemplateCache struct {
	sync.RWMutex
	M map[int]*model.Template
}

// 保存所有template数据,template关联的strategy在Strategies中保存，一个template多个strategy
var TemplateCache = &SafeTemplateCache{M: make(map[int]*model.Template)}

// 获取所有template数据
func (this *SafeTemplateCache) GetMap() map[int]*model.Template {
	this.RLock()
	defer this.RUnlock()
	return this.M
}

// 从数据库获取所有template数据
func (this *SafeTemplateCache) Init() {
	ts, err := db.QueryTemplates()
	if err != nil {
		return
	}

	this.Lock()
	defer this.Unlock()
	this.M = ts
}

type SafeHostTemplateIds struct {
	sync.RWMutex
	M map[int][]int
}

// 保存所有host id对应的group所关联的template id
var HostTemplateIds = &SafeHostTemplateIds{M: make(map[int][]int)}

// 获取所有host id对应的group所关联的template id
func (this *SafeHostTemplateIds) GetMap() map[int][]int {
	this.RLock()
	defer this.RUnlock()
	return this.M
}

// 从数据库获取所有host id对应的group所关联的template id
func (this *SafeHostTemplateIds) Init() {
	m, err := db.QueryHostTemplateIds()
	if err != nil {
		return
	}

	this.Lock()
	defer this.Unlock()
	this.M = m
}
