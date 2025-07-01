/*
* @Author: supbro
* @Date:   2025/6/11 22:43
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/11 22:43
 */
package calc_node

import (
	"wagner/app/domain"
)

var registry = make(map[string]func(*domain.ComputeContext) *domain.ComputeContext)

// 注册函数
func Register(name string, fn func(*domain.ComputeContext) *domain.ComputeContext) {
	registry[name] = fn
}

func GetFunction(name string) (func(*domain.ComputeContext) *domain.ComputeContext, bool) {
	fn, ok := registry[name]
	return fn, ok
}
