/*
* @Author: supbro
* @Date:   2025/6/6 13:02
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/6 13:02
 */
package olap_dao

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"reflect"
	"strconv"
	"strings"
	"wagner/app/domain"
	"wagner/app/service/calc_dynamic_param"
	"wagner/infrastructure/persistence/common"
	"wagner/infrastructure/persistence/entity"
	"wagner/infrastructure/persistence/query"
)

type HourSummaryResultDao struct {
	olapDb *gorm.DB
}

func CreateHourSummaryResultDao(olapClient *gorm.DB) *HourSummaryResultDao {
	return &HourSummaryResultDao{olapClient}
}

const batchSize = 500

func (dao *HourSummaryResultDao) BatchInsertOrUpdateByUnqKey(resultList []*entity.HourSummaryResultEntity) {
	// todo 如果没有任何字段更新，gmt_modified即便设置了CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP，也不会自动更新，看这里是否需要手动更新该字段
	dao.olapDb.Omit("gmt_create", "gmt_modified").Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "unique_key"}},                                                                                                                                                                                                                                                                                                                                                                                            // 冲突检测列（唯一索引或主键）
		DoUpdates: clause.AssignmentColumns([]string{"operate_time", "operate_day", "process_code", "position_code", "workplace_code", "workplace_name", "employee_number", "employee_name", "employee_position_code", "work_group_code", "region_code", "industry_code", "sub_industry_code", "work_load", "direct_work_time", "indirect_work_time", "idle_time", "rest_time", "attendance_time", "process_property", "properties", "is_deleted"}), // 更新字段
	}).CreateInBatches(resultList, batchSize)
}

func buildDynamicWorkLoadSelect(workLoadUnit []calc_dynamic_param.WorkLoadUnit) string {
	selects := make([]string, 0)
	for _, field := range workLoadUnit {
		selects = append(selects,
			fmt.Sprintf("JSON_UNQUOTE(JSON_EXTRACT(work_load, '$.%s')) as `%s`", field.Code, field.Code),
		)
	}
	return strings.Join(selects, ", ")
}

func (dao *HourSummaryResultDao) QueryEmployeeEfficiency(query query.HourSummaryResultQuery) []*entity.WorkLoadWithSummaryEntity {
	fixedSelect := "employee_number, employee_name, operate_day, workplace_code, workplace_name, process_code, position_code, " +
		"employee_position_code, work_group_code, region_code, industry_code, sub_industry_code, process_property, " +
		" direct_work_time, indirect_work_time, idle_time, rest_time, attendance_time"
	workLoadSelect := buildDynamicWorkLoadSelect(query.WorkLoadUnit)

	where := "operate_day >= ? and operate_day <= ? and workplace_code = ?"
	if query.EmployeeNumber != "" {
		where += " and employee_number = " + query.EmployeeNumber
	}
	if query.IsCrossPosition == domain.Cross {
		where += " and position_code = employee_position_code"
	} else if query.IsCrossPosition == domain.NoCross {
		where += " and position_code != employee_position_code"
	}

	subQuery := dao.olapDb.Table("hour_summary_result").
		Select(fixedSelect+","+workLoadSelect).
		Where(where, query.DateRange[0], query.DateRange[1], query.WorkplaceCode)

	mainSelect := "employee_number, employee_name, operate_day, workplace_code, workplace_name, " +
		" max(region_code) region_code, max(industry_code) industry_code, max(sub_industry_code) sub_industry_code, " +
		" max(employee_position_code) employee_position_code, max(work_group_code) work_group_code, max(process_property) process_property, " +
		" sum(direct_work_time) direct_work_time, sum(indirect_work_time) indirect_work_time, sum(idle_time) idle_time, sum(rest_time) rest_time, sum(attendance_time) attendance_time"
	groupBy := "employee_number, employee_name, operate_day, workplace_code, workplace_name"
	orderBy := "operate_day, employee_number"

	for _, workLoadUnit := range query.WorkLoadUnit {
		mainSelect += fmt.Sprintf(", sum(%s) %s", workLoadUnit.Code, workLoadUnit.Code)
	}

	if query.AggregateDimension == domain.Process {
		groupBy += " ,process_code, position_code"
		mainSelect += ", process_code, max(position_code) position_code"
		orderBy += ", process_code"
	} else if query.AggregateDimension == domain.Position {
		groupBy += " ,position_code"
		mainSelect += ", position_code"
		orderBy += ", position_code"
	}

	var rawResult []map[string]interface{}
	//var result []entity.EmployeeSummaryEntity
	dao.olapDb.Table("(?) as summary", subQuery).Select(mainSelect).
		Order(orderBy).
		Group(groupBy).
		Find(&rawResult)

	return dao.convertRaw2Entity(rawResult)
}

func (dao *HourSummaryResultDao) convertRaw2Entity(resultList []map[string]interface{}) []*entity.WorkLoadWithSummaryEntity {
	// 创建命名策略（默认使用蛇形命名）
	namer := schema.NamingStrategy{
		SingularTable: true, // 可选：单数表名
	}

	// 解析模型 Schema（不再需要 CacheStore）
	sch, err := schema.Parse(entity.EmployeeSummaryEntity{}, common.SchemaCache, &namer)
	if err != nil {
		panic(err)
	}

	result := make([]*entity.WorkLoadWithSummaryEntity, 0)
	for _, rawResult := range resultList {
		// 获取反射值并确保可设置
		e := entity.EmployeeSummaryEntity{}
		v := reflect.ValueOf(&e).Elem()

		for _, field := range sch.Fields {
			value := rawResult[field.DBName]
			delete(rawResult, field.DBName)
			entityName := field.Name

			entityField := v.FieldByName(entityName)
			if entityField.CanSet() && entityField.IsValid() {
				err := dao.setEntityValue(&entityField, value)
				if err != nil {
					panic(err)
				}
			}
		}

		var workLoad map[string]float64
		if len(rawResult) > 0 {
			workLoad = make(map[string]float64, len(rawResult))
			for key, value := range rawResult {
				workLoad[key] = value.(float64)
			}
		}

		result = append(result, &entity.WorkLoadWithSummaryEntity{
			EmployeeSummary: &e, WorkLoad: workLoad,
		})
	}

	return result

}

// setValue 处理具体类型转换
func (dao *HourSummaryResultDao) setEntityValue(field *reflect.Value, value interface{}) error {
	if value == nil {
		return nil // 忽略 NULL 值（假设字段允许零值）
	}

	val := reflect.ValueOf(value)
	fieldType := field.Type()

	// 类型完全匹配时直接赋值
	if val.Type().AssignableTo(fieldType) {
		field.Set(val)
		return nil
	}

	// 类型不匹配时尝试转换（例如 string -> time.Time）
	if val.Type().ConvertibleTo(fieldType) {
		convertedVal := val.Convert(fieldType)
		field.Set(convertedVal)
		return nil
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(value.(string))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if i, err := strconv.ParseInt(value.(string), 10, 64); err == nil {
			field.SetInt(i)
		} else {
			panic(err)
		}
	case reflect.Float32, reflect.Float64:
		if f, err := strconv.ParseFloat(value.(string), 64); err == nil {
			field.SetFloat(f)
		} else {
			panic(err)
		}
	case reflect.Bool:
		if b, err := strconv.ParseBool(value.(string)); err == nil {
			field.SetBool(b)
		} else {
			panic(err)
		}
	default:
		return fmt.Errorf("unsupported field type: %s", field.Type())
	}
	return nil
}

func (dao *HourSummaryResultDao) UpdateDeletedByUniqueKeyList(delete *query.HourSummaryResultDelete) {
	dao.olapDb.Model(entity.HourSummaryResultEntity{}).
		Where("unique_key not in (?)", delete.UniqueKeyList).
		Where("employee_number = ? and workplace_code = ? and operate_day = ? ", delete.EmployeeNumber, delete.WorkplaceCode, delete.OperateDay).
		Update("is_deleted", 1)

}
