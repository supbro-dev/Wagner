/*
* @Author: supbro
* @Date:   2025/6/5 21:05
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/5 21:05
 */
package reflect_util

import (
	"fmt"
	"reflect"
)

// HasField判断是否存在这个属性
func HasField(obj interface{}, fieldName string) (bool, error) {
	// 获取反射值对象，必须是指针类型，否则无法设置
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr {
		return false, fmt.Errorf("obj must be a pointer to struct")
	}
	// 获取字段
	field := v.FieldByName(fieldName)

	return field.IsValid(), nil
}

// SetField 通过反射设置结构体字段
func SetField(obj interface{}, fieldName string, value interface{}) error {
	// 获取反射值对象，必须是指针类型，否则无法设置
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("obj must be a pointer to struct")
	}

	// 获取指针指向的元素
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("obj must be a pointer to struct")
	}

	// 获取字段
	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return fmt.Errorf("field %s does not exist", fieldName)
	}

	// 检查字段是否可设置
	if !field.CanSet() {
		return fmt.Errorf("field %s is not settable", fieldName)
	}

	// 获取字段的类型
	fieldType := field.Type()
	// 传入值的反射对象
	val := reflect.ValueOf(value)

	// 如果传入值的类型和字段类型不一致，尝试转换
	if val.Type() != fieldType {
		// 检查是否可以转换
		if val.Type().ConvertibleTo(fieldType) {
			// 转换类型
			val = val.Convert(fieldType)
		} else {
			return fmt.Errorf("value type %s doesn't match field type %s and is not convertible", val.Type(), fieldType)
		}
	}

	// 设置字段值
	field.Set(val)
	return nil
}

// GetField 获取结构体单个字段的值（不支持嵌套）
// obj: 结构体实例或指针
// fieldName: 要获取的字段名称
func GetField(obj interface{}, fieldName string) (interface{}, error) {
	val := reflect.ValueOf(obj)

	// 处理指针
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil, fmt.Errorf("对象指针为nil")
		}
		val = val.Elem()
	}

	// 确保是结构体
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("需要结构体类型, 实际类型: %s", val.Kind())
	}

	// 获取字段
	field := val.FieldByName(fieldName)
	if !field.IsValid() {
		return nil, fmt.Errorf("字段不存在: %s", fieldName)
	}

	// 处理指针
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return nil, nil
		}
		field = field.Elem()
	}

	// 处理接口
	if field.Kind() == reflect.Interface && !field.IsNil() {
		field = field.Elem()
	}

	return field.Interface(), nil
}
