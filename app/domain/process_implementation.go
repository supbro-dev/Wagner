/*
* @Author: supbro
* @Date:   2025/7/2 16:29
* @Last Modified by:   supbro
* @Last Modified time: 2025/7/2 16:29
 */
package domain

import "wagner/infrastructure/persistence/entity"

type ProcessImplementation struct {
	Id         int64
	Name       string
	TargetType entity.TargetType
	TargetCode string
	TargetName string
	Status     entity.ImplementationStatus
}
