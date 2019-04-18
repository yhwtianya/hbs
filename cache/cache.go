package cache

import (
	"log"
	"time"
)

// 周期性用数据库数据更新缓存数据
func Init() {
	log.Println("cache begin")

	log.Println("#1 GroupPlugins...")
	GroupPlugins.Init()

	log.Println("#2 GroupTemplates...")
	GroupTemplates.Init()

	log.Println("#3 HostGroupsMap...")
	HostGroupsMap.Init()

	log.Println("#4 HostMap...")
	HostMap.Init()

	log.Println("#5 TemplateCache...")
	TemplateCache.Init()

	log.Println("#6 Strategies...")
	Strategies.Init(TemplateCache.GetMap())

	log.Println("#7 HostTemplateIds...")
	HostTemplateIds.Init()

	log.Println("#8 ExpressionCache...")
	ExpressionCache.Init()

	log.Println("#9 MonitoredHosts...")
	MonitoredHosts.Init()

	log.Println("cache done")

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
