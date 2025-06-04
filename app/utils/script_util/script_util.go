/*
* @Author: supbro
* @Date:   2025/6/3 12:47
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/3 12:47
 */
package script_util

import (
	"fmt"
	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
	"wagner/app/global/my_error"
	"wagner/infrastructure/persistence/entity"
)

func Run[P any, V any](script string, scriptType entity.ScriptType, input P, inputName string) (V, error) {
	var zero V
	switch scriptType {
	// lua脚本不支持原生的自定义类型
	case entity.LUA:
		return runLua[P, V](script, input, inputName)
		// Yaegi不支持module，只能使用tengo，但tengo实现自定义类型非常繁琐，需要实现tengo.Object接口，赋值和取值时还要进行类型转换
	case entity.GOLANG:
		return runGolang[P, V](script, input, inputName)
	case entity.REFLECT:
		return zero, nil
	case entity.EL:
		return zero, nil
	}
	return zero, fmt.Errorf("脚本解析失败: %v, %v", my_error.ScriptWrongTypeCode, my_error.ScriptWrongTypeMsg)
}

func runGolang[P any, V any](script string, input P, inputName string) (V, error) {
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
	_ = spt.Add(inputName, input)

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

func runLua[P any, V any](script string, input P, inputName string) (V, error) {
	L := lua.NewState()
	defer L.Close()

	obj := luar.New(L, input)

	L.SetGlobal(inputName, obj)

	if err := L.DoString(script); err != nil {
		var zero V
		return zero, err
	}

	ret := L.Get(-1)

	return ret.(V), nil
}
