package cache

import (
	"sync"

	"github.com/open-falcon/common/model"
	"github.com/open-falcon/hbs/db"
)

type SafeExpressionCache struct {
	sync.RWMutex
	L []*model.Expression
}

// 保存所有生效中的Expressions
var ExpressionCache = &SafeExpressionCache{}

// 获取所有生效中的Expressions
func (this *SafeExpressionCache) Get() []*model.Expression {
	this.RLock()
	defer this.RUnlock()
	return this.L
}

// 从数据库获取所有生效中的Expressions
func (this *SafeExpressionCache) Init() {
	es, err := db.QueryExpressions()
	if err != nil {
		return
	}

	this.Lock()
	defer this.Unlock()
	this.L = es
}
