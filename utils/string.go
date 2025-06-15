/*
 * @Description:
 * @Author: 安知鱼
 * @Date: 2025-06-15 12:20:47
 * @LastEditTime: 2025-06-15 12:20:55
 * @LastEditors: 安知鱼
 */
package utils

import "strings"

func JoinTags(tags []string) string {
	return strings.Join(tags, ",")
}
