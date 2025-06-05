/*
* @Author: supbro
* @Date:   2025/5/31 23:06
* @Last Modified by:   supbro
* @Last Modified time: 2025/5/31 23:06
 */
package lock

func Lock(employeeNumber string) (bool, error) {
	// todo redis分布式锁
	return true, nil
}

func Unlock(employeeNumber string) (bool, error) {
	return true, nil
}
