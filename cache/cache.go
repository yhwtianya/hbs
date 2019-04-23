package cache

import (
	"log"
	"time"
)

// 周期性用数据库数据更新缓存数据
func Init() {
	log.Println("cache begin")

	// 查询数据库，获得group id和plugin_dir关系
	log.Println("#1 GroupPlugins...")
	GroupPlugins.Init()

	// 查询数据库，获得group id和template ids关系
	log.Println("#2 GroupTemplates...")
	GroupTemplates.Init()

	// 查询数据库，获取host id和group ids的关系
	log.Println("#3 HostGroupsMap...")
	HostGroupsMap.Init()

	// 从数据库获取所有hostname和host id对应关系
	log.Println("#4 HostMap...")
	HostMap.Init()

	// 从数据库获取所有template数据
	log.Println("#5 TemplateCache...")
	TemplateCache.Init()

	// 从数据库获取所有生效中的strategy
	log.Println("#6 Strategies...")
	Strategies.Init(TemplateCache.GetMap())

	// 从数据库获取所有host id对应的group所关联的template id
	log.Println("#7 HostTemplateIds...")
	HostTemplateIds.Init()

	// 从数据库获取所有生效中的Expressions
	log.Println("#8 ExpressionCache...")
	ExpressionCache.Init()

	// 从数据库获取所有非维护状态的主机信息
	log.Println("#9 MonitoredHosts...")
	MonitoredHosts.Init()

	log.Println("cache done")

	// 周期性用数据库数据更新缓存数据
	go LoopInit()

}

// 周期性用数据库数据更新缓存数据
func LoopInit() {
	for {
		time.Sleep(time.Minute)
		GroupPlugins.Init()
		GroupTemplates.Init()
		HostGroupsMap.Init()
		HostMap.Init()
		TemplateCache.Init()
		Strategies.Init(TemplateCache.GetMap())
		HostTemplateIds.Init()
		ExpressionCache.Init()
		MonitoredHosts.Init()
	}
}
