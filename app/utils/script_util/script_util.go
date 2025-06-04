/*
* @Author: supbro
* @Date:   2025/6/3 12:47
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/3 12:47
 */
package script_util

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
	"wagner/app/global/my_error"
	"wagner/infrastructure/persistence/entity"
)

func Run[P any, V any](script string, input P, scriptType entity.ScriptType) (V, error) {
	var zero V
	switch scriptType {
	case entity.LUA:
		return runLua[P, V](script, input)
	case entity.GOLANG:
		return zero, nil
	case entity.REFLECT:
		return zero, nil
	case entity.EL:
		return zero, nil
	}
	return zero, fmt.Errorf("脚本解析失败: %v, %v", my_error.ScriptWrongTypeCode, my_error.ScriptWrongTypeMsg)
}

func runLua[P any, V any](script string, input P) (V, error) {
	L := lua.NewState()
	defer L.Close()

	L.SetGlobal("u", luar.New(L, input))

	if err := L.DoString(script); err != nil {
		var zero V
		return zero, err
	}

	ret := L.Get(-1)

	return ret.(V), nil
}
