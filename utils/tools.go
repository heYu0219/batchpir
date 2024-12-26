package utils

import (
	"bufio"
	cRand "crypto/rand"
	"fmt"
	"io"
	"math/big"
	mRand "math/rand"
	"os"
	"strings"
	"time"
)

// GenerateInteger 生成指定字节大小的整数
func GenerateInteger(bitSize int) *big.Int {
	// Generate a random big.Int with the specified bit size
	max := new(big.Int).Lsh(big.NewInt(1), uint(bitSize)) // 2^bitSize
	n, err := cRand.Int(cRand.Reader, max)
	if err != nil {
		panic(fmt.Sprintf("Error generating random big.Int: %v", err))
	}
	return n
}

// GenerateRawDB 生成 m x n 的随机整数数据库
func GenerateRawDB(n, numBytes int) ([][]*big.Int, error) {
	// 创建 m x n 的二维切片
	rawDB := make([][]*big.Int, n)
	for i := 0; i < n; i++ {
		rawDB[i] = make([]*big.Int, n)
		for j := 0; j < n; j++ {
			element := GenerateInteger(numBytes)

			rawDB[i][j] = element
		}
	}
	return rawDB, nil
}

// 生成随机向量
func GenerateRandomVector(length int) []*big.Int {
	// 设置随机种子，确保每次运行生成的随机数不同
	mRand.Seed(time.Now().UnixNano())

	// 创建一个长度为 length 的大整数切片
	vector := make([]*big.Int, length)

	// 填充数组，生成 1 到 100 之间的随机整数并转换为 *big.Int
	for i := 0; i < length; i++ {
		randomValue := mRand.Intn(100-1) + 1 // 随机整数范围为 [1, 100)
		vector[i] = big.NewInt(int64(randomValue))
	}

	return vector
}

// GenerateHint 计算提示，按照随机向量与每一行计算内积的方式，返回一个提示向量
func GenerateHint(DB [][]*big.Int, vector []*big.Int) []*big.Int {
	n := len(DB)
	ans := make([]*big.Int, n)

	for i := 0; i < n; i++ {
		tmp := big.NewInt(0)
		for j := 0; j < n; j++ {
			product := new(big.Int).Mul(DB[i][j], vector[j]) // 计算 DB[i][j] * vector[j]
			tmp.Add(tmp, product)                            // 累加到 tmp
		}
		ans[i] = tmp
	}
	return ans
}

// GetIndices 生成 n*n 数据库对角线上的索引集合
func GetIndices(n int) []int {
	indices := make([]int, n)
	for i := 0; i < n; i++ {
		indices[i] = i*n + i
	}
	return indices
}

// IndexToCoordinate 将查询索引转换为数据库对应的坐标
func IndexToCoordinate(index, n int) [2]int {
	var coordinate [2]int
	count := 0

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if count == index {
				coordinate[0] = i
				coordinate[1] = j
				return coordinate
			}
			count++
		}
	}
	return coordinate
}

// GetCoordinates 根据索引集合获得坐标集合
func GetCoordinates(indices []int, n int) [][2]int {
	var coordinates [][2]int

	for _, index := range indices {
		coordinate := IndexToCoordinate(index, n)
		coordinates = append(coordinates, coordinate)
	}

	return coordinates
}

// ConstructQuery 生成单行查询
func ConstructQuery(vector []*big.Int, col, alpha, beta int) []int {
	query := make([]int, len(vector))
	for i := 0; i < len(vector); i++ {
		// 使用 big.Int 的 SetInt64 将 alpha 与 vector[i] 相乘
		tmp := new(big.Int).Mul(vector[i], big.NewInt(int64(alpha)))
		query[i] = int(tmp.Int64()) // 转换为 int
		if i == col {
			// 对应列加上 beta
			query[i] += beta
		}
	}
	return query
}

// GetQueries 根据坐标生成最终查询
func GetQueries(coordinate [][2]int, vector []*big.Int, alpha, beta int) [][]int {
	queries := make([][]int, len(vector))
	for i := range queries {
		queries[i] = make([]int, len(vector)) // 初始化每一行的查询数组
	}

	// 遍历坐标，生成查询
	for _, cor := range coordinate {
		row := cor[0] // 行
		col := cor[1] // 列
		queries[row] = ConstructQuery(vector, col, alpha, beta)
	}

	return queries
}

// GetAnswer 服务器计算响应
func GetAnswer(DB [][]*big.Int, queries [][]int) []*big.Int {
	// 初始化结果数组，长度为数据库的行数
	ans := make([]*big.Int, len(DB))
	for i := 0; i < len(DB); i++ {
		tmp := big.NewInt(0) // 使用 big.NewInt 初始化 tmp
		for j := 0; j < len(DB[i]); j++ {
			// 使用 big.Int 的 Mul 方法进行乘法计算
			tmp.Add(tmp, new(big.Int).Mul(DB[i][j], big.NewInt(int64(queries[i][j]))))
		}
		// 将计算的结果赋值到 ans[i]
		ans[i] = tmp
	}
	return ans
}

// MakeErrorAns 基于正确的应答构造错误应答
func MakeErrorAns1(ans []*big.Int, idx []int) []*big.Int {
	// 创建一个新的错误应答数组，初始为 ans 的副本
	errorAns := make([]*big.Int, len(ans))
	for i, v := range ans {
		errorAns[i] = new(big.Int).Set(v)
	}

	// 随机数生成器，确保每次运行有不同的结果
	mRand.Seed(time.Now().UnixNano())

	// 修改指定索引的值
	for _, i := range idx {
		// 获取当前值，并加上一个随机数（0到19）
		// randValue := mRand.Intn(20)
		errorAns[i].Add(errorAns[i], big.NewInt(5))
	}

	return errorAns
}
func MakeErrorAns2(ans []*big.Int, idx []int) []*big.Int {
	// 创建一个新的错误应答数组，初始为 ans 的副本
	errorAns := make([]*big.Int, len(ans))
	for i, v := range ans {
		errorAns[i] = new(big.Int).Set(v)
	}

	// 随机数生成器，确保每次运行有不同的结果
	mRand.Seed(time.Now().UnixNano())

	// 修改指定索引的值
	for _, i := range idx {
		// 获取当前值，并加上一个随机数（0到19）
		// randValue := mRand.Intn(20)
		errorAns[i].Sub(errorAns[i], big.NewInt(7))
	}

	return errorAns
}

// DecodeAns 解码服务器的响应
func DecodeAns(ans []*big.Int, hint []*big.Int, alpha int, beta int) []*big.Int {
	// 创建一个文件数组来存储解码后的结果
	file := make([]*big.Int, len(ans))

	// 对每个元素进行解码
	for i := 0; i < len(ans); i++ {
		// 计算 alpha * hint[i]
		alphaHint := new(big.Int).Mul(big.NewInt(int64(alpha)), hint[i])

		// ans[i] - alpha * hint[i]
		delta := new(big.Int).Sub(ans[i], alphaHint)

		// (ans[i] - alpha * hint[i]) / beta
		file[i] = new(big.Int).Div(delta, big.NewInt(int64(beta)))
	}

	return file
}

// CompareAns 比较两个响应对应位置是否相等，不等则记录
func CompareAns(ans1, ans2 []*big.Int) ([][]*big.Int, []int) {
	// 用来存储差异的值
	var idx []int
	var diffAns1, diffAns2 []*big.Int

	// 遍历 ans1 和 ans2，比较对应位置的值
	for i := 0; i < len(ans1); i++ {
		if ans1[i].Cmp(ans2[i]) != 0 { // 使用 Cmp 比较两个 BigInt 是否相等
			idx = append(idx, i)                                   // 记录不相等的位置
			diffAns1 = append(diffAns1, new(big.Int).Set(ans1[i])) // 记录 ans1 中的值
			diffAns2 = append(diffAns2, new(big.Int).Set(ans2[i])) // 记录 ans2 中的值
		}
	}

	// 返回两个差异值数组和差异位置的索引
	return [][]*big.Int{diffAns1, diffAns2}, idx
}

// searchEntry 查找元素在数据库中的索引，并返回这些索引的集合
func SearchEntry(entry *big.Int, db [][]*big.Int) []int {
	var indexSet []int
	count := -1

	// 遍历数据库，查找与 entry 匹配的所有元素
	for i := 0; i < len(db); i++ {
		for j := 0; j < len(db[i]); j++ {
			count++
			// 使用 big.Int 的 SetString 比较两个大整数是否相等
			if db[i][j].Cmp(entry) == 0 {
				// 如果找到匹配的元素，添加当前的索引位置到 indexSet
				indexSet = append(indexSet, count)
			}
		}
	}
	return indexSet
}

// getAllAnsIndex 返回错误响应元组对应的索引集合
func GetAllAnsIndex(anss [][]*big.Int, db [][]*big.Int) [][][]int {
	var indexList [][][]int
	for _, ans := range anss {
		var serverIndex [][]int
		for _, ansEntry := range ans {
			// 查找每个响应元素在数据库中的索引
			entryIndexList := SearchEntry(ansEntry, db)
			// 将每个服务器的索引集合加入到当前的索引集合中
			serverIndex = append(serverIndex, entryIndexList)
		}
		// 记录所有服务器的索引集合
		indexList = append(indexList, serverIndex)
	}
	return indexList
}

// 读取key=value类型的配置文件
func InitConfig(path string) map[string]string {
	config := make(map[string]string)

	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	r := bufio.NewReader(f)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		s := strings.TrimSpace(string(b))
		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}
		key := strings.TrimSpace(s[:index])
		if len(key) == 0 {
			continue
		}
		value := strings.TrimSpace(s[index+1:])
		if len(value) == 0 {
			continue
		}
		config[key] = value
	}
	return config
}
