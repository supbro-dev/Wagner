/*
* @Author: supbro
* @Date:   2025/7/1 13:14
* @Last Modified by:   supbro
* @Last Modified time: 2025/7/1 13:14
 */
package process

import (
	"math"
	"strconv"
	"wagner/app/domain"
	"wagner/app/global/business_error"
	"wagner/app/utils/json_util"
	"wagner/app/utils/pinyin_util"
	"wagner/infrastructure/persistence/dao"
	"wagner/infrastructure/persistence/entity"
	"wagner/infrastructure/persistence/query"
)

type ProcessService interface {
	FindFirstProcess(positionCode string, workplace *domain.Workplace) *domain.ProcessPosition
	FindProcessList(workplace *domain.Workplace) []*domain.ProcessPosition
	FindProcessPositionList(workplace *domain.Workplace) []*domain.ProcessPosition
	FindProcessImplementationListByPage(targetType string, workplaceCode string, industryCode string, subIndustryCode string, currentPage int, pageSize int) ([]*domain.ProcessImplementation, int)
	Save(processImpl *domain.ProcessImplementation) (int64, *business_error.BusinessError)
	GetImplementationById(id int64) *domain.ProcessImplementation
	GetProcessPositionTree(id int64) *domain.ProcessPositionTreeNode
	FindProcessByParentCode(parentCode string, version int64) []*domain.ProcessPosition
	GenerateProcessCode(processName string, version int64) string
	SaveProcessPosition(processPosition *domain.ProcessPosition)
	FindProcessByCode(code string, version int64) *domain.ProcessPosition
}

type ProcessServiceImpl struct {
	processPositionDao       *dao.ProcessPositionDao
	processImplementationDao *dao.ProcessImplementationDao
	workplaceDao             *dao.WorkplaceDao
}

func CreateProcessServiceImpl(processPositionDao *dao.ProcessPositionDao, processImplementationDao *dao.ProcessImplementationDao, workplaceDao *dao.WorkplaceDao) ProcessService {
	return &ProcessServiceImpl{processPositionDao, processImplementationDao, workplaceDao}
}

var OtherProcess = &domain.ProcessPosition{
	Name: "其他",
	Code: "Others",
}

func (service *ProcessServiceImpl) FindProcessByCode(code string, version int64) *domain.ProcessPosition {
	e := service.processPositionDao.FindByCode(code, version)
	return service.convertPositionEntity2Domain(e)
}

func (service *ProcessServiceImpl) SaveProcessPosition(processPosition *domain.ProcessPosition) {
	e := service.convertPositionDomain2Entity(processPosition)
	service.processPositionDao.Insert(e)
}

func (service *ProcessServiceImpl) GenerateProcessCode(processName string, version int64) string {
	processCode := pinyin_util.ConvertMixedString(processName)

	// 查找是否有同名code
	var existedPosition *entity.ProcessPositionEntity
	for i := 1; i < 100; i++ {
		existedPosition = service.processPositionDao.FindByCode(processCode, version)
		if existedPosition == nil {
			break
		} else {
			processCode = processCode + strconv.Itoa(i)
		}
	}

	return processCode
}

func (service *ProcessServiceImpl) FindProcessByParentCode(parentCode string, version int64) []*domain.ProcessPosition {
	result := service.iterateFindProcessByParentCode(parentCode, version)

	domainList := make([]*domain.ProcessPosition, 0)
	for _, positionEntity := range result {
		d := service.convertPositionEntity2Domain(positionEntity)
		domainList = append(domainList, d)
	}

	return domainList
}

// 数据量可控，采用递归查询
func (service *ProcessServiceImpl) iterateFindProcessByParentCode(parentCode string, version int64) []*entity.ProcessPositionEntity {
	list := service.processPositionDao.FindByParentCodeAndVersion(parentCode, version)
	result := make([]*entity.ProcessPositionEntity, 0)
	for _, positionEntity := range list {
		if positionEntity.Type == entity.INDIRECT_PROCESS || positionEntity.Type == entity.DIRECT_PROCESS {
			result = append(result, positionEntity)
		}
	}
	if len(list) > 0 {
		for _, child := range list {
			if child.Type != entity.INDIRECT_PROCESS && child.Type != entity.DIRECT_PROCESS {
				childResult := service.iterateFindProcessByParentCode(child.Code, version)
				result = append(result, childResult...)
			}
		}
	}

	return result
}

func (service *ProcessServiceImpl) GetProcessPositionTree(id int64) *domain.ProcessPositionTreeNode {
	implementationEntity := service.processImplementationDao.FindById(id)

	var positionList []*entity.ProcessPositionEntity

	version := implementationEntity.ProcessPositionRootId
	switch implementationEntity.TargetType {
	case entity.Workplace:
		workplace := service.workplaceDao.FindByCode(implementationEntity.TargetCode)
		positionList = service.processPositionDao.FindByIndustry(workplace.IndustryCode, workplace.SubIndustryCode, version)
	case entity.Industry:
		positionList = service.processPositionDao.FindByIndustry(implementationEntity.TargetCode, "", version)
	case entity.SubIndustry:
		subIndustryCode := implementationEntity.TargetCode
		industryCode := service.workplaceDao.FindSubIndustryBySubindustryCode(subIndustryCode)
		positionList = service.processPositionDao.FindByIndustry(industryCode, subIndustryCode, version)
	}

	code2Node := make(map[string]*domain.ProcessPositionTreeNode)

	for _, position := range positionList {
		if position.Type == entity.DIRECT_PROCESS || position.Type == entity.INDIRECT_PROCESS {
			continue
		}
		node := domain.ProcessPositionTreeNode{
			Name: position.Name,
			Code: position.Code,
			Type: position.Type,
		}

		parentCode := position.ParentCode
		if parentCode != "-1" {
			parentNode := code2Node[parentCode]
			parentNode.Children = append(parentNode.Children, &node)
		}

		code2Node[node.Code] = &node
	}

	return code2Node[implementationEntity.Code]
}

// 根据id查找环节实施信息
func (service *ProcessServiceImpl) GetImplementationById(id int64) *domain.ProcessImplementation {
	e := service.processImplementationDao.FindById(id)
	impl := service.convertImplEntity2Domain(e)

	return impl
}

// 查找岗位下的第一个环节
func (service *ProcessServiceImpl) FindFirstProcess(positionCode string, workplace *domain.Workplace) *domain.ProcessPosition {
	implementationEntity := service.processImplementationDao.FindByWorkplaceCode(workplace.Code)
	if implementationEntity == nil {
		implementationEntity = service.processImplementationDao.FindByIndustry(workplace.IndustryCode, workplace.SubIndustryCode)
	}

	if implementationEntity == nil {
		return nil
	}
	version := implementationEntity.ProcessPositionRootId

	positionList := service.processPositionDao.FindByIndustry(workplace.IndustryCode, workplace.SubIndustryCode, version)

	minOrder := math.MaxInt
	var minProcess *entity.ProcessPositionEntity
	for _, positionEntity := range positionList {
		if positionEntity.ParentCode == positionCode && positionEntity.SortIndex < minOrder {
			minOrder = positionEntity.SortIndex
			minProcess = positionEntity
		}
	}
	if minOrder == math.MaxInt {
		return nil
	}

	return service.convertPositionEntity2Domain(minProcess)
}

// 根据工作点查找所有环节
func (service *ProcessServiceImpl) FindProcessList(workplace *domain.Workplace) []*domain.ProcessPosition {
	implementationEntity := service.processImplementationDao.FindByWorkplaceCode(workplace.Code)
	if implementationEntity == nil {
		implementationEntity = service.processImplementationDao.FindByIndustry(workplace.IndustryCode, workplace.SubIndustryCode)
	}

	if implementationEntity == nil {
		return nil
	}
	version := implementationEntity.ProcessPositionRootId
	positionList := service.processPositionDao.FindByIndustry(workplace.IndustryCode, workplace.SubIndustryCode, version)

	return service.buildLeafNodePaths(positionList)
}

// 按从根至叶子节点的顺序查出来
func (service *ProcessServiceImpl) FindProcessPositionList(workplace *domain.Workplace) []*domain.ProcessPosition {
	implementationEntity := service.processImplementationDao.FindByWorkplaceCode(workplace.Code)
	if implementationEntity == nil {
		implementationEntity = service.processImplementationDao.FindByIndustry(workplace.IndustryCode, workplace.SubIndustryCode)
	}

	if implementationEntity == nil {
		return nil
	}
	version := implementationEntity.ProcessPositionRootId

	positionList := service.processPositionDao.FindByIndustry(workplace.IndustryCode, workplace.SubIndustryCode, version)

	domainList := make([]*domain.ProcessPosition, 0)
	for _, positionEntity := range positionList {
		domain := service.convertPositionEntity2Domain(positionEntity)
		domainList = append(domainList, domain)
	}

	return domainList
}

func (service *ProcessServiceImpl) FindProcessImplementationListByPage(targetType string, workplaceCode string, industryCode string, subIndustryCode string, currentPage int, pageSize int) ([]*domain.ProcessImplementation, int) {
	var targetCode string
	switch entity.TargetType(targetType) {
	case entity.Workplace:
		targetCode = workplaceCode
	case entity.Industry:
		targetCode = industryCode
	case entity.SubIndustry:
		targetCode = subIndustryCode
	}
	processImplementationQuery := query.ProcessImplementationQuery{
		TargetType:  targetType,
		TargetCode:  targetCode,
		CurrentPage: currentPage,
		PageSize:    pageSize,
	}

	implementationEntities := service.processImplementationDao.QueryProcessImplementation(processImplementationQuery)

	implementationList := make([]*domain.ProcessImplementation, 0)
	for _, e := range implementationEntities {
		impl := service.convertImplEntity2Domain(e)
		implementationList = append(implementationList, impl)
	}

	total := service.processImplementationDao.CountProcessImplementation(processImplementationQuery)

	return implementationList, total
}

func (service *ProcessServiceImpl) convertImplEntity2Domain(e *entity.ProcessImplementationEntity) *domain.ProcessImplementation {
	impl := &domain.ProcessImplementation{
		Id:                    e.Id,
		Code:                  e.Code,
		Name:                  e.Name,
		TargetType:            e.TargetType,
		TargetCode:            e.TargetCode,
		Status:                e.Status,
		ProcessPositionRootId: e.ProcessPositionRootId,
	}
	switch e.TargetType {
	case entity.Workplace:
		workplace := service.workplaceDao.FindByCode(e.TargetCode)
		impl.TargetName = workplace.Name
	default:
		impl.TargetName = e.TargetCode
	}

	return impl
}

func (service *ProcessServiceImpl) Save(processImpl *domain.ProcessImplementation) (int64, *business_error.BusinessError) {
	existed := service.processImplementationDao.FindOne(&query.ProcessImplementationQuery{
		TargetType: string(processImpl.TargetType),
		TargetCode: processImpl.TargetCode,
		Code:       processImpl.Code,
	})

	if existed != nil && existed.Id != processImpl.Id {
		return 0, business_error.ExistSameCodeProcessImpl(processImpl.Code)
	}

	e := entity.ProcessImplementationEntity{
		Code:       processImpl.Code,
		Name:       processImpl.Name,
		TargetType: processImpl.TargetType,
		TargetCode: processImpl.TargetCode,
		Status:     processImpl.Status,
	}

	service.processImplementationDao.Save(&e)

	// 如果没有进行更新，gorm不会回填id
	if e.Id == 0 {
		return processImpl.Id, nil
	} else {
		return e.Id, nil
	}

}

func (service *ProcessServiceImpl) buildLeafNodePaths(positionEntities []*entity.ProcessPositionEntity) []*domain.ProcessPosition {
	// 创建三个核心映射
	entityMap := make(map[string]*entity.ProcessPositionEntity)     // code -> 实体指针
	childrenMap := make(map[string][]*entity.ProcessPositionEntity) // parentCode -> 子节点列表
	parentMap := make(map[string]*entity.ProcessPositionEntity)     // code -> 父节点指针

	// 最大部门层级
	maxDeptLevel := 0

	// 构建映射关系
	for i := range positionEntities {
		e := positionEntities[i]

		maxDeptLevel = max(maxDeptLevel, e.Level)

		code := (*e).Code
		parentCode := e.ParentCode

		// 添加到实体映射
		entityMap[code] = e

		// 添加到父节点映射
		if parentCode != "" {
			parentMap[code] = entityMap[parentCode]
		}

		// 添加到子节点映射
		childrenMap[parentCode] = append(childrenMap[parentCode], e)
	}

	// 部门层级排除最后两级的环节和岗位
	maxDeptLevel = maxDeptLevel - 2

	// 收集所有叶子节点（没有子节点的节点）
	var leafNodes []*entity.ProcessPositionEntity
	for code, e := range entityMap {
		// 没有子节点即为叶子节点
		if len(childrenMap[code]) == 0 {
			leafNodes = append(leafNodes, e)
		}
	}

	// 为每个叶子节点构建路径
	result := make([]*domain.ProcessPosition, 0, len(leafNodes))
	for _, leaf := range leafNodes {
		path := service.buildParentPath(leaf, parentMap)
		d := service.convertPositionEntity2Domain(leaf)
		d.Path = path
		// 记录最大部门层级
		d.MaxDeptLevel = maxDeptLevel
		// 只收集直接/间接环节
		if leaf.Type == entity.DIRECT_PROCESS || leaf.Type == entity.INDIRECT_PROCESS {
			result = append(result, d)
		}
	}

	return result
}

func (service *ProcessServiceImpl) convertPositionEntity2Domain(e *entity.ProcessPositionEntity) *domain.ProcessPosition {

	d := domain.ProcessPosition{
		Id:         e.Id,
		Name:       e.Name,
		Code:       e.Code,
		Level:      e.Level,
		Script:     e.Script,
		Version:    e.Version,
		ParentCode: e.ParentCode,
		Type:       e.Type,
	}

	if e.Properties != "" {
		if propertyMap, err := json_util.Parse2Map(e.Properties); err == nil {
			d.Properties = propertyMap
		}
	}
	return &d
}

// 递归构建从叶子节点到根节点的路径
func (service *ProcessServiceImpl) buildParentPath(node *entity.ProcessPositionEntity, parentMap map[string]*entity.ProcessPositionEntity) []*domain.ProcessPosition {
	var path []*domain.ProcessPosition

	// 从直接父节点开始
	current := parentMap[node.Code]

	// 递归向上遍历父节点
	for current != nil {
		path = append(path, service.convertPositionEntity2Domain(current))

		// 移动到上一级父节点
		current = parentMap[current.Code]
	}

	// 返回从父节点到根节点的路径
	return path
}

func (service *ProcessServiceImpl) convertPositionDomain2Entity(position *domain.ProcessPosition) *entity.ProcessPositionEntity {
	parent := service.processPositionDao.FindByCode(position.ParentCode, int64(position.Version))

	return &entity.ProcessPositionEntity{
		Code:            position.Code,
		Name:            position.Name,
		ParentCode:      position.ParentCode,
		Type:            position.Type,
		Level:           parent.Level + 1,
		Version:         position.Version,
		IndustryCode:    parent.IndustryCode,
		SubIndustryCode: parent.SubIndustryCode,
		Properties:      json_util.ToJsonString(position.Properties),
	}
}
