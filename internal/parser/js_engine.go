// Package parser implements parsing and extraction functionality for the Pong0 application.
// This file specifically handles JavaScript execution and key generation logic
// required to authenticate with the Ping0.cc service.
package parser

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"

	"ping0/internal/constants"
)

// calculateHashStart uses crypto/sha256 to hash the input string
// and returns the beginning of the hash as a string.
// This is equivalent to the calculateHash function in JavaScript.
func calculateHashStart(input string, length int) (string, error) {
	hasher := sha256.New()
	hasher.Write([]byte(input))
	hashBytes := hasher.Sum(nil)
	hashHex := hex.EncodeToString(hashBytes)

	if length > len(hashHex) {
		length = len(hashHex)
	}

	return hashHex[:length], nil
}

// calculatePow calculates a proof-of-work value that produces a hash
// starting with the specified difficulty prefix.
//
// Parameters:
//   - x1: The base string (usually a hex string)
//   - difficulty: The prefix that the hash should start with
//
// Returns:
//   - int: The POW value
//   - error: If an error occurs during hash calculation
func calculatePow(x1, difficulty string) (int, error) {
	counter := 0
	difficultyLen := len(difficulty)

	for {
		input := fmt.Sprintf("%s%d", x1, counter)
		hash, err := calculateHashStart(input, difficultyLen)
		if err != nil {
			return 0, err
		}

		if hash == difficulty {
			return counter, nil
		}

		counter++
		// Add a reasonable limit to prevent infinite loops
		if counter > 100000 {
			return 0, fmt.Errorf("超过最大迭代次数，无法找到符合条件的POW值")
		}
	}
}

// calculateJs1Key calculates js1key by applying a series of operations to the input
//
// Parameters:
//   - x1: The hex string to use as input
//   - animated: Whether the page has animation
//
// Returns:
//   - int: The js1key value
func calculateJs1Key(x1 string, animated bool) int {
	result := 0

	// Process x1 in 4-character chunks
	for i := 0; i < 32; i += 4 {
		// Parse chunk as hex and add 12
		chunk := x1[i : i+4]
		val, _ := strconv.ParseInt(chunk, 16, 64)
		result += int(val) + 12
		result &= 16777215 // Keep within 24 bits

		// Apply series of transformations as per original algorithm
		result = result ^ 9592
		result &= 16777215
		result = result + 4856
		result &= 16777215

		if animated {
			result = result ^ 40996
		}

		result = result ^ 5007
		result &= 16777215

		if animated {
			result = result - 83957
		}

		result = result - 8842
		result &= 16777215

		// Additional transformations
		result = result ^ 4621
		result &= 16777215
		result = result ^ 5497
		result &= 16777215

		if animated {
			result = result + 26924
		}

		result = result + 5961
		result &= 16777215

		if animated {
			result = result - 12005
		}

		result = result - 6533
		result &= 16777215
		result = result + 1149
		result &= 16777215

		// More transformations
		result = result ^ 4784
		result &= 16777215
		result = result - 3624
		result &= 16777215
		result = result ^ 1855
		result &= 16777215
		result = result - 2903
		result &= 16777215
		result = result ^ 9651
		result &= 16777215

		// Final transformations
		result = result ^ 9740
		result &= 16777215
		result = result - 7250
		result &= 16777215
		result = result + 8334
		result &= 16777215
		result = result - 5332
		result &= 16777215
		result = result + 8264
		result &= 16777215
		result = result - 1840
		result &= 16777215
		result = result ^ 7994
		result &= 16777215
		result = result - 6564
		result &= 16777215
		result = result - 9319
		result &= 16777215
		result = result ^ 9276
		result &= 16777215
		result = result - 8188
		result &= 16777215
		result = result - 6630
		result &= 16777215
		result = result ^ 4756
		result &= 16777215
		result = result - 8429
		result &= 16777215
		result = result - 5819
		result &= 16777215

		if animated {
			result = result - 84724
		}

		result = result - 3288
		result &= 16777215
		result = result + 3350
		result &= 16777215
		result = result - 7509
		result &= 16777215
		result = result ^ 8297
		result &= 16777215
		result = result ^ 5024
		result &= 16777215
		result = result ^ 2855
		result &= 16777215
		result = result - 3995
		result &= 16777215
		result = result ^ 3949
		result &= 16777215
		result = result + 5215
		result &= 16777215
		result = result + 1856
		result &= 16777215
		result = result - 6845
		result &= 16777215
		result = result ^ 8122
		result &= 16777215
		result = result ^ 4941
		result &= 16777215
		result = result + 2276
		result &= 16777215
		result = result - 5399
		result &= 16777215
		result = result - 1237
		result &= 16777215
		result = result ^ 4935
	}

	return result
}

// Keys 表示生成的js1key和pow值
type Keys struct {
	Js1key string
	Pow    string
}

// GenerateKey 根据新的算法生成访问密钥
// 该函数会生成两个密钥：js1key和pow，这是访问Ping0.cc服务的必要凭证。
//
// 参数:
//   - jsPath: JavaScript文件路径
//   - x1Value: 从初始页面提取的x1值
//   - difficultyValue: 从初始页面提取的difficulty值
//
// 返回:
//   - *Keys: 包含js1key和pow值的结构体
//   - error: 如果生成过程中出现错误则返回对应错误信息
func GenerateKey(jsPath, x1Value, difficultyValue string) (*Keys, error) {
	if len(x1Value) != 32 {
		return nil, fmt.Errorf("无效的x1Value长度: 期望32, 实际%d", len(x1Value))
	}

	if constants.Verbose {
		fmt.Printf("开始生成密钥:\n")
		fmt.Printf("- x1Value: %s\n", x1Value)
		fmt.Printf("- difficultyValue: %s\n", difficultyValue)
		fmt.Printf("- jsPath: %s\n", jsPath)
		fmt.Printf("- BaseURL: %s\n", constants.BaseURL)
	}

	// 1. 计算js1key值
	animated := false // 页面动画状态固定为关闭
	js1key := calculateJs1Key(x1Value, animated)

	// 2. 计算pow值
	pow, err := calculatePow(x1Value, difficultyValue)
	if err != nil {
		return nil, fmt.Errorf("计算POW失败: %w", err)
	}

	if constants.Verbose {
		fmt.Printf("生成的js1key: %d\n", js1key)
		fmt.Printf("生成的pow: %d\n", pow)
	}

	return &Keys{
		Js1key: fmt.Sprintf("%d", js1key),
		Pow:    fmt.Sprintf("%d", pow),
	}, nil
}
