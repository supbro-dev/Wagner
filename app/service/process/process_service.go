/*
* @Author: supbro
* @Date:   2025/7/1 13:14
* @Last Modified by:   supbro
* @Last Modified time: 2025/7/1 13:14
 */
package process

import (
	"math"
	"wagner/app/domain"
	"wagner/app/utils/json_util"
	"wagner/infrastructure/persistence/dao"
	"wagner/infrastructure/persistence/entity"
)

type ProcessService interface {
	FindProcessPositionFirstProcess(positionCode string, industryCode, subIndustryCode string) *domain.ProcessPosition
	FindProcessPositionByWorkplace(workplaceCode string) []*domain.ProcessPosition
	FindProcessPositionByIndustry(industryCode, subIndustryCode string) []*domain.ProcessPosition
	FindProcessPositionListByIndustry(industryCode, subIndustryCode string) []*domain.ProcessPosition
}

type ProcessServiceImpl struct {
	processPositionDao  *dao.ProcessPositionDao
	processImplementDao *dao.ProcessImplementationDao
	workplaceDao        *dao.WorkplaceDao
}

func CreateProcessServiceImpl(processPositionDao *dao.ProcessPositionDao, processImplementDao *dao.ProcessImplementationDao, workplaceDao *dao.WorkplaceDao) ProcessService {
	return &ProcessServiceImpl{processPositionDao, processImplementDao, workplaceDao}
}

var OtherProcess = &domain.ProcessPosition{
	Name: "其他",
	Code: "Others",
}

func (service *ProcessServiceImpl) FindProcessPositionFirstProcess(positionCode string, industryCode, subIndustryCode string) *domain.ProcessPosition {
	implementationEntity := service.processImplementDao.FindByIndustry(industryCode, subIndustryCode)

	if implementationEntity == nil {
		return nil
	}
	version := implementationEntity.ProcessPositionRootId

	positionList := service.processPositionDao.FindByIndustry(industryCode, subIndustryCode, version)

	minOrder := math.MaxInt
	var minProcess *entity.ProcessPositionEntity
	for _, positionEntity := range positionList {
		if positionEntity.ParentCode == positionCode && positionEntity.Order < minOrder {
			minOrder = positionEntity.Order
			minProcess = positionEntity
		}
	}
	if minOrder == math.MaxInt {
		return nil
	}

	return service.convertEntity2Domain(minProcess)
}

// 根据工作点编码获取标准岗位模型
func (service *ProcessServiceImpl) FindProcessPositionByWorkplace(workplaceCode string) []*domain.ProcessPosition {
	// todo 这里应该查找工序实施配置，在没有实施流程时，直接根据工作点查找行业的标准模型
	positions := make([]*domain.ProcessPosition, 0)

	workplace := service.workplaceDao.FindByCode(workplaceCode)

	if workplace == nil {
		return positions
	}

	positionList := service.FindProcessPositionByIndustry(workplace.IndustryCode, workplace.SubIndustryCode)
	return positionList
}

// 按从根至叶子节点的顺序查出来
func (service *ProcessServiceImpl) FindProcessPositionListByIndustry(industryCode, subIndustryCode string) []*domain.ProcessPosition {
	implementationEntity := service.processImplementDao.FindByIndustry(industryCode, subIndustryCode)

	if implementationEntity == nil {
		return nil
	}
	version := implementationEntity.ProcessPositionRootId

	positionList := service.processPositionDao.FindByIndustry(industryCode, subIndustryCode, version)

	domainList := make([]*domain.ProcessPosition, 0)
	for _, positionEntity := range positionList {
		domain := service.convertEntity2Domain(positionEntity)
		domainList = append(domainList, domain)
	}

	return domainList
}

func (service *ProcessServiceImpl) FindProcessPositionByIndustry(industryCode, subIndustryCode string) []*domain.ProcessPosition {
	implementationEntity := service.processImplementDao.FindByIndustry(industryCode, subIndustryCode)

	if implementationEntity == nil {
		return nil
	}
	version := implementationEntity.ProcessPositionRootId
	positionList := service.processPositionDao.FindByIndustry(industryCode, subIndustryCode, version)

	return service.buildLeafNodePaths(positionList)
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
		d := service.convertEntity2Domain(leaf)
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

func (service *ProcessServiceImpl) convertEntity2Domain(e *entity.ProcessPositionEntity) *domain.ProcessPosition {

	d := domain.ProcessPosition{
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
		path = append(path, service.convertEntity2Domain(current))

		// 移动到上一级父节点
		current = parentMap[current.Code]
	}

	// 返回从父节点到根节点的路径
	return path
}
