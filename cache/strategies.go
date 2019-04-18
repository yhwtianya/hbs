package cache

import (
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/open-falcon/common/model"
	"github.com/open-falcon/hbs/db"
	"github.com/toolkits/container/set"
)

type SafeStrategies struct {
	sync.RWMutex
	M map[int]*model.Strategy
}

// 保存所有生效中的strategy信息，一个template对应多个strategy
var Strategies = &SafeStrategies{M: make(map[int]*model.Strategy)}

// 获取所有生效中的strategy
func (this *SafeStrategies) GetMap() map[int]*model.Strategy {
	this.RLock()
	defer this.RUnlock()
	return this.M
}

// 从数据库获取所有生效中的strategy
func (this *SafeStrategies) Init(tpls map[int]*model.Template) {
	m, err := db.QueryStrategies(tpls)
	if err != nil {
		return
	}

	this.Lock()
	defer this.Unlock()
	this.M = m
}

// 获取主机关联的模板、父模板策略里设置的内置指标
func GetBuiltinMetrics(hostname string) ([]*model.BuiltinMetric, error) {
	ret := []*model.BuiltinMetric{}
	hid, exists := HostMap.GetID(hostname)
	if !exists {
		return ret, nil
	}

	gids, exists := HostGroupsMap.GetGroupIds(hid)
	if !exists {
		return ret, nil
	}

	// 根据gids，获取绑定的所有tids
	tidSet := set.NewIntSet()
	for _, gid := range gids {
		tids, exists := GroupTemplates.GetTemplateIds(gid)
		if !exists {
			continue
		}

		for _, tid := range tids {
			// 去重
			tidSet.Add(tid)
		}
	}

	tidSlice := tidSet.ToSlice()
	if len(tidSlice) == 0 {
		return ret, nil
	}

	// 继续寻找这些tid的ParentId
	allTpls := TemplateCache.GetMap()
	for _, tid := range tidSlice {
		pids := ParentIds(allTpls, tid)
		for _, pid := range pids {
			tidSet.Add(pid)
		}
	}

	// 终于得到了最终的tid列表
	tidSlice = tidSet.ToSlice()

	// 把tid列表用逗号拼接在一起
	count := len(tidSlice)
	tidStrArr := make([]string, count)
	for i := 0; i < count; i++ {
		tidStrArr[i] = strconv.Itoa(tidSlice[i])
	}

	// 内置指标包括必监控指标和需配置策略才监控的指标，这里获取策略配置里涉及到的内置指标
	return db.QueryBuiltinMetrics(strings.Join(tidStrArr, ","))
}

// 根据template id获取其所有的父template ids,结果包含自身id，且按id升序排序
func ParentIds(allTpls map[int]*model.Template, tid int) (ret []int) {
	depth := 0
	for {
		if tid <= 0 {
			break
		}

		if t, exists := allTpls[tid]; exists {
			// partend id插入到列表尾部
			ret = append(ret, tid)
			tid = t.ParentId
		} else {
			break
		}

		// 模板最多嵌套10层??
		depth++
		if depth == 10 {
			log.Println("[ERROR] template inherit cycle. id:", tid)
			return []int{}
		}
	}

	sz := len(ret)
	if sz <= 1 {
		return
	}

	desc := make([]int, sz)
	for i, item := range ret {
		j := sz - i - 1
		// 按id升序排序，
		desc[j] = item
	}

	return desc
}
