/*
* @Author: supbro
* @Date:   2025/6/3 12:47
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/3 12:47
 */
package script_util

import (
	"github.com/traefik/yaegi/interp"
	"wagner/app/domain"
	"wagner/app/global/my_error"
	"wagner/infrastructure/persistence/entity"
)

func Parse[P any, V any](scriptName, script string, scriptType entity.ScriptType) (*func(P) V, error) {
	switch scriptType {
	case entity.GOLANG:
		return parseGolang[P, V](scriptName, script)
	case entity.REFLECT:
		return nil, nil
	case entity.EL:
		return nil, nil
	}
	return nil, my_error.NewError(my_error.ScriptWrongTypeCode, my_error.ScriptWrongTypeMsg)
}

func parseGolang[P any, V any](scriptName string, script string) (*func(P) V, error) {
	i := interp.New(interp.Options{})

	_, err := i.Eval(script)
	if err != nil {
		return nil, err
	}

	// 然后调用函数
	v, err := i.Eval(scriptName + "." + scriptName)
	if err != nil {
		return nil, err
	}

	f := v.Interface().(func(p P) V)
	return &f, nil
}

func Run(scriptName, script string, scriptType entity.ScriptType, ctx domain.ComputeContext) (interface{}, error) {
	switch scriptType {
	case entity.GOLANG:
		return runGolang(scriptName, script, ctx)
	case entity.REFLECT:
		return runReflect(script, ctx)
	case entity.EL:
		return runEl(script, ctx)
	}
	return nil, my_error.NewError(my_error.ScriptWrongTypeCode, my_error.ScriptWrongTypeMsg)
}

func runGolang(name string, script string, ctx domain.ComputeContext) (interface{}, error) {
	return nil, nil
}

func runEl(script string, ctx domain.ComputeContext) (interface{}, error) {
	return nil, nil
}

func runReflect(script string, ctx domain.ComputeContext) (interface{}, error) {
	return nil, nil
}
