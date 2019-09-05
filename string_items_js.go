// Lute - A structured markdown engine.
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//     http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.

// +build js

package lute

// []byte~string 之间的快速转换优化会导致生成 JavaScript 端代码问题，所以此处还是使用内存拷贝。

// fromItems 快速转换 items 为 string。
func fromItems(items items) string {
	//return *(*string)(unsafe.Pointer(&items))
	return string(items)
}

// toItems 快速转换 str 为 items。
func toItems(str string) items {
	//x := (*[2]uintptr)(unsafe.Pointer(&str))
	//h := [3]uintptr{x[0], x[1], x[1]}
	//return *(*items)(unsafe.Pointer(&h))
	return items(str)
}