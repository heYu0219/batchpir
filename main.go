package main

import (
	"batchPIR/utils"
	"fmt"
	"math/big"
	"strconv"
)

func main() {
	config := utils.InitConfig("config")

	// 元素大小
	num_bits, _ := strconv.Atoi(config["num_bits"])
	//文件行列数
	n, _ := strconv.Atoi(config["n"])

	alpha, _ := strconv.Atoi(config["alpha"])
	beta, _ := strconv.Atoi(config["beta"])

	//原始数据库构造
	db, _ := utils.GenerateRawDB(n, num_bits)

	//构造随机向量
	vector := utils.GenerateRandomVector(n)

	//生成提示
	hint := utils.GenerateHint(db, vector)

	//构造查询
	indices := utils.GetIndices(n)

	//索引转换为坐标集合
	coordinates := utils.GetCoordinates(indices, n)

	//根据坐标集合生成查询
	queries := utils.GetQueries(coordinates, vector, alpha, beta)

	//服务器计算响应
	ans := utils.GetAnswer(db, queries)

	//构造错误响应索引
	errorIdx1 := []int{1} // 错误索引1
	errorIdx2 := []int{3} // 错误索引2

	//构造错误响应
	error_ans1 := utils.MakeErrorAns(ans, errorIdx1)
	error_ans2 := utils.MakeErrorAns(ans, errorIdx2)

	//解码响应
	decode_ans1 := utils.DecodeAns(error_ans1, hint, alpha, beta)
	decode_ans2 := utils.DecodeAns(error_ans2, hint, alpha, beta)

	// 记录结果，初始化为其中一个响应的
	result := make(map[int]*big.Int)

	// 将结果填充到 map 中
	for i := 0; i < len(decode_ans1); i++ {
		result[indices[i]] = decode_ans1[i]
	}

	//返回两个服务器不同响应的元素集合 每个服务器的响应保存在一个列表中
	diffAns, idxs := utils.CompareAns(decode_ans1, decode_ans2)

	// 遍历所有不相同的索引
	for _, id := range idxs {
		// 将对应的索引位置的值设置为 nil
		result[indices[id]] = nil
	}

	//验证：返回两个服务器不同响应的元素的索引集合
	verifyIndex := utils.GetAllAnsIndex(diffAns, db)

	// 提取第一组和第二组数据
	serverIndex1 := verifyIndex[0] // 获取第一组错误响应的索引
	serverIndex2 := verifyIndex[1] // 获取第二组错误响应的索引

	//
	// 遍历 idxs，获取对应的索引并检查是否匹配
	for i := 0; i < len(idxs); i++ {
		id := idxs[i]           // 客户端请求的索引下标
		queryIdx := indices[id] // 获取查询的索引值

		indexList := serverIndex1[i] // 服务器1返回的索引集合
		if len(indexList) == 0 || indexList == nil {
			continue // 如果服务器1返回的索引为空，则跳过
		}

		// 查找 queryIdx 是否在 indexList 中
		found := false
		for _, idx := range indexList {
			if idx == queryIdx {
				found = true
				break
			}
		}

		if found {
			result[queryIdx] = decode_ans1[id] // 如果找到，更新结果
		}
	}

	for i := 0; i < len(idxs); i++ {
		id := idxs[i]           // 客户端请求的索引下标
		queryIdx := indices[id] // 获取查询的索引值

		indexList := serverIndex2[i] // 服务器1返回的索引集合
		if len(indexList) == 0 || indexList == nil {
			continue // 如果服务器1返回的索引为空，则跳过
		}

		// 查找 queryIdx 是否在 indexList 中
		found := false
		for _, idx := range indexList {
			if idx == queryIdx {
				found = true
				break
			}
		}

		if found {
			result[queryIdx] = decode_ans2[id] // 如果找到，更新结果
		}
	}

	fmt.Println(indices)
	fmt.Println(result)

}
