/*
* @Author: supbro
* @Date:   2025/6/3 12:47
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/3 12:47
 */
package script_util

import (
	"fmt"
	"github.com/PaesslerAG/gval"
	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
	"reflect"
	"wagner/app/global/business_error"
)

type ScriptType string

var (
	LUA ScriptType = "LUA"
	// 可以使用tengo，暂时未实现
	GOLANG ScriptType = "GOLANG"
	EL     ScriptType = "EL"
)

func Run[P any, V any](scriptName, script string, scriptType ScriptType, input P) (V, error) {
	var zero V
	switch scriptType {
	// lua脚本不支持原生的自定义类型
	case LUA:
		return runLua[P, V](script, input)
		// Yaegi不支持module，只能使用tengo，但tengo实现自定义类型非常繁琐，需要实现tengo.Object接口，赋值和取值时还要进行类型转换
	case GOLANG:
		return runGolang[P, V](script, input)
	case EL:
		// 检查输入类型是否为 map[string]interface{}
		if isMapStringInterface(input) {
			// 类型断言确保输入符合要求
			if m, ok := any(input).(map[string]interface{}); ok {
				return runEl[map[string]interface{}, V](script, m)
			}
		}

		if reflect.TypeOf(input).Kind() == reflect.Map {
			// 尝试将结构体转换为 map[string]interface{}
			return runEl[map[string]interface{}, V](script, convertToMapStringInterface(input))
		}
		return zero, business_error.ElScriptMustUseMap()
	}
	return zero, business_error.ScriptWrongType()
}

func convertToMapStringInterface(input any) map[string]interface{} {
	newMap := make(map[string]interface{})
	val := reflect.ValueOf(input)
	for _, key := range val.MapKeys() {
		newMap[key.String()] = val.MapIndex(key).Interface()
	}
	return newMap
}

// 检查是否为 map[string]interface{} 类型
func isMapStringInterface(input interface{}) bool {
	t := reflect.TypeOf(input)
	if t == nil {
		return false
	}

	return t.Kind() == reflect.Map &&
		t.Key().Kind() == reflect.String &&
		t.Elem().Kind() == reflect.Interface
}

func runEl[P map[string]interface{}, V any](script string, input P) (V, error) {
	// 简单表达式
	value, err := gval.Evaluate(script, input)
	return value.(V), err
}

func runGolang[P any, V any](script string, input P) (V, error) {
	// 创建脚本
	spt := tengo.NewScript([]byte(script))

	fmtPrintln := tengo.UserFunction{
		Value: func(args ...tengo.Object) (ret tengo.Object, err error) {
			fmt.Println(args)
			return nil, nil
		},
	}

	spt.SetImports(stdlib.GetModuleMap("fmt", "time"))
	// 设置输入变量
	_ = spt.Add("input", input)

	// 将包装后的函数添加到脚本，命名为"fmtPrintln"
	err := spt.Add("fmtPrintln", &fmtPrintln)

	// 运行脚本
	compiled, err := spt.Run()
	if err != nil {
		var zero V
		return zero, err
	}

	// 获取结果
	c := compiled.Get("ctxResult")

	return c.Object().(V), nil
}

func runLua[P any, V any](script string, input P) (V, error) {
	L := lua.NewState()
	defer L.Close()

	obj := luar.New(L, input)

	L.SetGlobal("input", obj)

	if err := L.DoString(script); err != nil {
		var zero V
		return zero, err
	}

	ret := L.Get(-1)

	return ret.(V), nil
}
