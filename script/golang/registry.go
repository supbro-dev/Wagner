/*
* @Author: supbro
* @Date:   2025/6/4 13:15
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/4 13:15
 */
package golang

import "reflect"

var registry = map[string]reflect.Value{}

// 注册函数（自动调用）
func register(name string, fn any) {
	registry[name] = reflect.ValueOf(fn)
}

func GetFunction(name string) (reflect.Value, bool) {
	fn, ok := registry[name]
	return fn, ok
}

// 注册 Run 函数
func init() {
	register("RunTest", RunTest)
}
