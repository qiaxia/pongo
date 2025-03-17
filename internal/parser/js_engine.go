// Package parser implements parsing and extraction functionality for the Pong0 application.
// This file specifically handles JavaScript execution and key generation logic
// required to authenticate with the Ping0.cc service.
package parser

import (
	"fmt"
	"strconv"

	"ping0/internal/constants"
)

// transformValue 辅助函数，执行特定序列的数学运算
// 该函数对输入值应用一系列预定义的数学运算，模拟Ping0.cc服务的JavaScript算法。
// 这些运算包括加法、减法和异或操作，顺序和参数值是固定的。
//
// 参数:
//   - value: 要转换的初始整数值
//
// 返回:
//   - int: 经过一系列运算后的结果值
func transformValue(value int) int {
	operations := []struct {
		op    string
		value int
	}{
		{"-", 1381},
		{"-", 5990},
		{"^", 8629},
		{"+", 7102},
		{"+", 1130},
		{"+", 6993},
		{"+", 1630},
		{"-", 9155},
		{"^", 6178},
		{"-", 3193},
		{"^", 63238},
	}

	for _, operation := range operations {
		switch operation.op {
		case "+":
			value += operation.value
		case "-":
			value -= operation.value
		case "^":
			value ^= operation.value
		}
	}

	return value
}

// generateJs1Key 生成js1key值
// 该函数根据x1值和当前URL生成访问密钥，这是访问Ping0.cc服务的必要凭证。
// 它模拟了Ping0.cc前端JavaScript中的密钥生成算法。
//
// 参数:
//   - x1Value: 从初始页面提取的x1值
//   - currentURL: 当前请求的URL路径
//
// 返回:
//   - int: 生成的密钥值
func generateJs1Key(x1Value, currentURL string) int {
	if len(x1Value) != 32 {
		return 0
	}

	result := 0
	// 动画状态固定为关闭
	hasAnimation := false

	// 循环处理x1Value中的每4个字符
	for i := 0; i < len(x1Value); i += 4 {
		if i+4 > len(x1Value) {
			break // 防止越界
		}

		// 解析4个字符为16进制数值并加12
		val, err := strconv.ParseInt(x1Value[i:i+4], 16, 64)
		if err != nil {
			continue // 解析失败则跳过
		}

		result += int(val) + 12
		// 限制为24位
		result &= 16777215

		// 执行一系列数学运算
		result = applyMathSequence(result, hasAnimation, len(currentURL), len(x1Value))
	}

	return result
}

// applyMathSequence 应用数学运算序列
func applyMathSequence(result int, hasAnimation bool, urlLength, x1Length int) int {
	// 24位掩码
	const mask24 = 16777215

	// 基本运算
	result = (result + 3464) & mask24
	result = (result - 9490) & mask24
	result = (result ^ 3351) & mask24

	// 检查页面动画状态
	if hasAnimation {
		result = (result + 12493) & mask24
	}

	result = (result + 5064) & mask24
	result = (result + 3508) & mask24
	result = (result - 3539) & mask24

	// URL长度相关条件
	if transformValue(63750)+urlLength > transformValue(63932) {
		result = (result - 1373) & mask24
	}

	result = (result + 4812) & mask24
	result = (result + 8552) & mask24
	result = (result ^ 5344) & mask24

	// 更多URL长度相关条件
	if transformValue(63945)+urlLength > transformValue(63945) {
		result = (result + 1064) & mask24
	}

	// x1长度相关条件
	if x1Length+transformValue(63750) < transformValue(63777) {
		result = (result + 6416) & mask24
	}

	if x1Length+transformValue(61368) > transformValue(63728) {
		result = (result + 8914) & mask24
	}

	if hasAnimation {
		result = (result ^ 65589) & mask24
	}

	result = (result ^ 8309) & mask24
	result = (result + 2767) & mask24
	result = (result - 8772) & mask24

	if transformValue(66758)+transformValue(68531) < transformValue(29831) {
		result = (result - 5595) & mask24
	}

	if hasAnimation {
		result = (result ^ 81826) & mask24
	}

	result = (result ^ 6722) & mask24
	result = (result ^ 7352) & mask24

	if x1Length+transformValue(63931) < transformValue(63753) {
		result = (result ^ 8071) & mask24
	}

	if hasAnimation {
		result = (result + 96879) & mask24
	}

	result = (result + 3828) & mask24
	result = (result ^ 4591) & mask24

	if hasAnimation {
		result = (result + 28190) & mask24
	}

	result = (result + 3148) & mask24
	result = (result ^ 8553) & mask24

	if transformValue(63946)+urlLength > transformValue(63933) {
		result = (result - 9610) & mask24
	}

	result = (result ^ 6118) & mask24
	result = (result ^ 2005) & mask24

	if x1Length+transformValue(71367) > transformValue(63738) {
		result = (result + 1963) & mask24
	}

	result = (result - 5915) & mask24
	result = (result + 2651) & mask24

	if hasAnimation {
		result = (result ^ 12414) & mask24
	}

	result = (result ^ 6886) & mask24
	result = (result ^ 5011) & mask24

	if hasAnimation {
		result = (result ^ 79397) & mask24
	}

	result = (result ^ 4167) & mask24

	if hasAnimation {
		result = (result - 67531) & mask24
	}

	result = (result - 3290) & mask24
	result = (result ^ 6567) & mask24

	return result
}

// GenerateKey 生成密钥
func GenerateKey(jsPath, x1Value string) (string, error) {
	if len(x1Value) != 32 {
		return "", fmt.Errorf("无效的x1Value长度: 期望32, 实际%d", len(x1Value))
	}

	if constants.Verbose {
		fmt.Printf("开始生成密钥:\n")
		fmt.Printf("- x1Value: %s\n", x1Value)
		fmt.Printf("- jsPath: %s\n", jsPath)
		fmt.Printf("- BaseURL: %s\n", constants.BaseURL)
	}

	// 直接生成js1key，不需要下载JS文件
	key := generateJs1Key(x1Value, constants.BaseURL)
	if constants.Verbose {
		fmt.Printf("生成的密钥: %d\n", key)
	}

	return fmt.Sprintf("%d", key), nil
}
