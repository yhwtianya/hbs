package http

import (
	"fmt"
	"net/http"

	"github.com/open-falcon/common/model"
	"github.com/open-falcon/hbs/cache"
)

func configProcRoutes() {
	// 获取所有生效中的Expressions
	http.HandleFunc("/expressions", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, cache.ExpressionCache.Get())
	})

	// 获取所有上报自身信息agent的主机名
	http.HandleFunc("/agents", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, cache.Agents.Keys())
	})

	// 获取所有非维护状态的主机信息
	http.HandleFunc("/hosts", func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]*model.Host, len(cache.MonitoredHosts.Get()))
		for k, v := range cache.MonitoredHosts.Get() {
			data[fmt.Sprint(k)] = v
		}
		RenderDataJson(w, data)
	})

	// 获取所有生效中的strategy
	http.HandleFunc("/strategies", func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]*model.Strategy, len(cache.Strategies.GetMap()))
		for k, v := range cache.Strategies.GetMap() {
			data[fmt.Sprint(k)] = v
		}
		RenderDataJson(w, data)
	})

	// 获取所有template数据
	http.HandleFunc("/templates", func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]*model.Template, len(cache.TemplateCache.GetMap()))
		for k, v := range cache.TemplateCache.GetMap() {
			data[fmt.Sprint(k)] = v
		}
		RenderDataJson(w, data)
	})

	// 根据hostname获取关联的插件
	http.HandleFunc("/plugins/", func(w http.ResponseWriter, r *http.Request) {
		hostname := r.URL.Path[len("/plugins/"):]
		RenderDataJson(w, cache.GetPlugins(hostname))
	})

}
