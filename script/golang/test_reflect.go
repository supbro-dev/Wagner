/*
* @Author: supbro
* @Date:   2025/6/4 13:07
* @Last Modified by:   supbro
* @Last Modified time: 2025/6/4 13:07
 */
package golang

import (
	"fmt"
	"wagner/app/domain"
)

func RunTest(ctx *domain.ComputeContext) *domain.ComputeContext {
	fmt.Println("runTest")
	return &domain.ComputeContext{}
}
