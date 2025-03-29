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

// obf replicates the obfuscation function (_0x34ab46) from the updated newjs1keypow.js
// It applies a series of arithmetic and bitwise operations to the input number.
//
// Parameters:
//   - n: The input number to transform
//
// Returns:
//   - int: The transformed number
func obf(n int) int {
	n = n + 4566
	n = n ^ 2258
	n = n ^ 7675
	n = n ^ 7668
	n = n + 4032
	n = n - 7637
	n = n - 6257
	n = n - 7159
	n = n ^ 8417
	n = n + 4296
	n = n ^ 42851
	return n
}

// calculateJs1Key calculates js1key by applying a series of operations to the input
// Updated to match the algorithm in newjs1keypow.js
//
// Parameters:
//   - x1: The hex string to use as input
//   - locationHref: The current location.href value
//   - animated: Whether the page has animation
//
// Returns:
//   - int: The js1key value
func calculateJs1Key(x1 string, locationHref string, animated bool) int {
	// Start with an initial value derived from the first 4 hex digits of x1
	hexVal, _ := strconv.ParseInt(x1[:4], 16, 64)
	js1key := (int(hexVal) + 12) & 0xFFFFFF

	// Apply a conditional modification using constants derived from obf
	if obf(46441)+obf(35365) < obf(28159) {
		js1key = (js1key - 1608) & 0xFFFFFF
	}

	if animated {
		js1key = js1key ^ 53400
	}

	js1key = (js1key ^ 3084) & 0xFFFFFF
	js1key = (js1key - 6341) & 0xFFFFFF

	if animated {
		js1key = (js1key + 90994) & 0xFFFFFF
	}

	js1key = (js1key + 3769) & 0xFFFFFF
	js1key = (js1key - 8999) & 0xFFFFFF

	// Add additional value
	js1key = (js1key + 6920) & 0xFFFFFF
	if animated {
		js1key = (js1key + 78439) & 0xFFFFFF
	}
	js1key = (js1key + 1264) & 0xFFFFFF

	if animated {
		js1key = js1key ^ 39852
	}
	js1key = js1key ^ 7249
	js1key = (js1key - 3808) & 0xFFFFFF
	js1key = (js1key + 2625) & 0xFFFFFF
	js1key = (js1key + 9192) & 0xFFFFFF

	// New condition based on x1.length
	if len(x1)+obf(56907) < obf(57347) {
		js1key = (js1key - 3116) & 0xFFFFFF
	}

	js1key = js1key ^ 8580
	js1key = (js1key - 5467) & 0xFFFFFF

	if animated {
		js1key = js1key ^ 24814
	}

	js1key = js1key ^ 7636
	js1key = (js1key & 0xFFFFFF)

	// Additional conditional operations
	if len(x1)+obf(56910) < obf(57268) {
		js1key = (js1key + 6090) & 0xFFFFFF
	}

	if len(x1)+obf(63126) > obf(57338) {
		js1key = (js1key - 4085) & 0xFFFFFF
	}

	if animated {
		js1key = (js1key - 85694) & 0xFFFFFF
	}

	js1key = (js1key - 3211) & 0xFFFFFF
	js1key = (js1key + 6561) & 0xFFFFFF

	if animated {
		js1key = (js1key - 64921) & 0xFFFFFF
	}

	js1key = (js1key - 3687) & 0xFFFFFF
	js1key = (js1key + 6366) & 0xFFFFFF

	// Process second chunk of x1 if available
	if len(x1) >= 8 {
		hexVal, _ := strconv.ParseInt(x1[4:8], 16, 64)
		js1key += int(hexVal) + 12
		js1key &= 0xFFFFFF

		if animated {
			js1key = (js1key - 23486) & 0xFFFFFF
		}

		js1key = (js1key - 1608) & 0xFFFFFF
		js1key = (js1key ^ 3084) & 0xFFFFFF
		js1key = (js1key - 6341) & 0xFFFFFF

		if obf(35569)+obf(47549) < obf(25553) {
			js1key = (js1key + 3769) & 0xFFFFFF
		}

		js1key = (js1key - 8999) & 0xFFFFFF

		if obf(56917)+len(locationHref) > obf(56913) {
			js1key = (js1key + 6920) & 0xFFFFFF
		}

		if animated {
			js1key = (js1key + 70398) & 0xFFFFFF
		}

		js1key = (js1key + 1264) & 0xFFFFFF
		js1key = (js1key ^ 7249) & 0xFFFFFF

		if obf(56912)+len(locationHref) > obf(56914) {
			js1key = (js1key - 3808) & 0xFFFFFF
		}

		js1key = (js1key + 2625) & 0xFFFFFF
		js1key = (js1key + 9192) & 0xFFFFFF
		js1key = (js1key - 3116) & 0xFFFFFF
		js1key = (js1key ^ 8580) & 0xFFFFFF

		if animated {
			js1key = (js1key - 58679) & 0xFFFFFF
		}

		js1key = (js1key - 5467) & 0xFFFFFF

		if len(x1)+obf(58748) > obf(57264) {
			js1key = (js1key ^ 7636) & 0xFFFFFF
		}

		js1key = (js1key + 6090) & 0xFFFFFF

		if obf(56972)+len(locationHref) > obf(56911) {
			js1key = (js1key - 4085) & 0xFFFFFF
		}

		js1key = (js1key - 3211) & 0xFFFFFF
		js1key = (js1key + 6561) & 0xFFFFFF

		if obf(56911)+len(locationHref) > obf(56912) {
			js1key = (js1key - 3687) & 0xFFFFFF
		}

		js1key = (js1key + 6366) & 0xFFFFFF
	}

	// Process third chunk of x1 if available
	if len(x1) >= 12 {
		hexVal, _ := strconv.ParseInt(x1[8:12], 16, 64)
		js1key += int(hexVal) + 12
		js1key &= 0xFFFFFF

		if animated {
			js1key = (js1key - 33986) & 0xFFFFFF
		}

		js1key = (js1key - 1608) & 0xFFFFFF

		if animated {
			js1key = (js1key ^ 16876) & 0xFFFFFF
		}

		js1key = (js1key ^ 3084) & 0xFFFFFF
		js1key = (js1key - 6341) & 0xFFFFFF
		js1key = (js1key + 3769) & 0xFFFFFF

		if len(x1)+obf(56911) < obf(56926) {
			js1key = (js1key - 8999) & 0xFFFFFF
		}

		if obf(35250)+obf(58115) < obf(24814) {
			js1key = (js1key + 6920) & 0xFFFFFF
		}

		js1key = (js1key + 1264) & 0xFFFFFF
		js1key = (js1key ^ 7249) & 0xFFFFFF

		if animated {
			js1key = (js1key - 19308) & 0xFFFFFF
		}

		js1key = (js1key - 3808) & 0xFFFFFF
		js1key = (js1key + 2625) & 0xFFFFFF
		js1key = (js1key + 9192) & 0xFFFFFF
		js1key = (js1key - 3116) & 0xFFFFFF

		if obf(56972)+len(locationHref) > obf(56912) {
			js1key = (js1key ^ 8580) & 0xFFFFFF
		}

		if animated {
			js1key = (js1key - 29698) & 0xFFFFFF
		}

		js1key = (js1key - 5467) & 0xFFFFFF
		js1key = (js1key ^ 7636) & 0xFFFFFF

		if obf(58053)+obf(32799) < obf(72110) {
			js1key = (js1key + 6090) & 0xFFFFFF
		}

		js1key = (js1key - 4085) & 0xFFFFFF

		if obf(56913)+len(locationHref) > obf(56911) {
			js1key = (js1key - 3211) & 0xFFFFFF
		}

		if obf(45344)+obf(44944) < obf(24687) {
			js1key = (js1key + 6561) & 0xFFFFFF
		}

		js1key = (js1key - 3687) & 0xFFFFFF
		js1key = (js1key + 6366) & 0xFFFFFF
	}

	// Process fourth chunk of x1 if available
	if len(x1) >= 16 {
		hexVal, _ := strconv.ParseInt(x1[12:16], 16, 64)
		js1key += int(hexVal) + 12
		js1key &= 0xFFFFFF
		js1key = (js1key - 1608) & 0xFFFFFF

		if obf(47787)+obf(58712) < obf(69976) {
			js1key = (js1key ^ 3084) & 0xFFFFFF
		}

		if animated {
			js1key = (js1key - 69093) & 0xFFFFFF
		}

		js1key = (js1key - 6341) & 0xFFFFFF
		js1key = (js1key + 3769) & 0xFFFFFF
		js1key = (js1key - 8999) & 0xFFFFFF

		if obf(56972)+len(locationHref) > obf(56972) {
			js1key = (js1key + 6920) & 0xFFFFFF
		}

		if animated {
			js1key = (js1key + 65053) & 0xFFFFFF
		}

		js1key = (js1key + 1264) & 0xFFFFFF
		js1key = (js1key ^ 7249) & 0xFFFFFF
		js1key = (js1key - 3808) & 0xFFFFFF

		if animated {
			js1key = (js1key + 91021) & 0xFFFFFF
		}

		js1key = (js1key + 2625) & 0xFFFFFF

		if animated {
			js1key = (js1key + 51460) & 0xFFFFFF
		}

		js1key = (js1key + 9192) & 0xFFFFFF
		js1key = (js1key - 3116) & 0xFFFFFF
		js1key = (js1key ^ 8580) & 0xFFFFFF
		js1key = (js1key - 5467) & 0xFFFFFF
		js1key = (js1key ^ 7636) & 0xFFFFFF
		js1key = (js1key + 6090) & 0xFFFFFF

		if len(x1)+obf(57982) > obf(57364) {
			js1key = (js1key - 4085) & 0xFFFFFF
		}

		js1key = (js1key - 3211) & 0xFFFFFF
		js1key = (js1key + 6561) & 0xFFFFFF
		js1key = (js1key - 3687) & 0xFFFFFF
		js1key = (js1key + 6366) & 0xFFFFFF
	}

	// Process fifth chunk of x1 if available
	if len(x1) >= 20 {
		hexVal, _ := strconv.ParseInt(x1[16:20], 16, 64)
		js1key += int(hexVal) + 12
		js1key &= 0xFFFFFF

		if obf(58011)+obf(60648) < obf(24265) {
			js1key = (js1key - 1608) & 0xFFFFFF
		}

		js1key = (js1key ^ 3084) & 0xFFFFFF
		js1key = (js1key - 6341) & 0xFFFFFF

		if obf(48375)+obf(51432) < obf(15265) {
			js1key = (js1key + 3769) & 0xFFFFFF
		}

		if animated {
			js1key = (js1key - 99151) & 0xFFFFFF
		}

		js1key = (js1key - 8999) & 0xFFFFFF

		if len(x1)+obf(56907) < obf(57336) {
			js1key = (js1key + 6920) & 0xFFFFFF
		}

		js1key = (js1key + 1264) & 0xFFFFFF
		js1key = (js1key ^ 7249) & 0xFFFFFF

		if animated {
			js1key = (js1key - 21490) & 0xFFFFFF
		}

		js1key = (js1key - 3808) & 0xFFFFFF
		js1key = (js1key + 2625) & 0xFFFFFF
		js1key = (js1key + 9192) & 0xFFFFFF

		if len(x1)+obf(56972) < obf(57344) {
			js1key = (js1key - 3116) & 0xFFFFFF
		}

		if len(x1)+obf(56972) < obf(57264) {
			js1key = (js1key ^ 8580) & 0xFFFFFF
		}

		if obf(44918)+obf(45865) < obf(25134) {
			js1key = (js1key - 5467) & 0xFFFFFF
		}

		if len(x1)+obf(48463) > obf(56892) {
			js1key = (js1key ^ 7636) & 0xFFFFFF
		}

		js1key = (js1key + 6090) & 0xFFFFFF

		if animated {
			js1key = (js1key - 43662) & 0xFFFFFF
		}

		js1key = (js1key - 4085) & 0xFFFFFF

		if animated {
			js1key = (js1key - 40853) & 0xFFFFFF
		}

		js1key = (js1key - 3211) & 0xFFFFFF
		js1key = (js1key + 6561) & 0xFFFFFF

		if obf(58656)+obf(36353) < obf(15184) {
			js1key = (js1key - 3687) & 0xFFFFFF
		}

		if obf(58323)+obf(63494) < obf(70537) {
			js1key = (js1key + 6366) & 0xFFFFFF
		}
	}

	// Process sixth chunk of x1 if available
	if len(x1) >= 24 {
		hexVal, _ := strconv.ParseInt(x1[20:24], 16, 64)
		js1key += int(hexVal) + 12
		js1key &= 0xFFFFFF

		if animated {
			js1key = (js1key - 69932) & 0xFFFFFF
		}

		js1key = (js1key - 1608) & 0xFFFFFF

		if animated {
			js1key = (js1key ^ 51824) & 0xFFFFFF
		}

		js1key = (js1key ^ 3084) & 0xFFFFFF
		js1key = (js1key - 6341) & 0xFFFFFF

		if animated {
			js1key = (js1key + 44364) & 0xFFFFFF
		}

		js1key = (js1key + 3769) & 0xFFFFFF

		if len(x1)+obf(35762) > obf(56936) {
			js1key = (js1key - 8999) & 0xFFFFFF
		}

		if animated {
			js1key = (js1key + 45930) & 0xFFFFFF
		}

		js1key = (js1key + 6920) & 0xFFFFFF
		js1key = (js1key + 1264) & 0xFFFFFF

		if len(x1)+obf(45131) > obf(57326) {
			js1key = (js1key ^ 7249) & 0xFFFFFF
		}

		if animated {
			js1key = (js1key - 21919) & 0xFFFFFF
		}

		js1key = (js1key - 3808) & 0xFFFFFF
		js1key = (js1key + 2625) & 0xFFFFFF

		if len(x1)+obf(45629) > obf(57358) {
			js1key = (js1key + 9192) & 0xFFFFFF
		}

		js1key = (js1key - 3116) & 0xFFFFFF
		js1key = (js1key ^ 8580) & 0xFFFFFF
		js1key = (js1key - 5467) & 0xFFFFFF

		if obf(56914)+len(locationHref) > obf(56913) {
			js1key = (js1key ^ 7636) & 0xFFFFFF
		}

		js1key = (js1key + 6090) & 0xFFFFFF
		js1key = (js1key - 4085) & 0xFFFFFF
		js1key = (js1key - 3211) & 0xFFFFFF
		js1key = (js1key + 6561) & 0xFFFFFF

		if len(x1)+obf(36030) > obf(56902) {
			js1key = (js1key - 3687) & 0xFFFFFF
		}

		if animated {
			js1key = (js1key + 69043) & 0xFFFFFF
		}

		js1key = (js1key + 6366) & 0xFFFFFF
	}

	// Process seventh chunk of x1 if available
	if len(x1) >= 28 {
		hexVal, _ := strconv.ParseInt(x1[24:28], 16, 64)
		js1key += int(hexVal) + 12
		js1key &= 0xFFFFFF

		if len(x1)+obf(57832) > obf(57332) {
			js1key = (js1key - 1608) & 0xFFFFFF
		}

		if obf(48352)+obf(45592) < obf(17965) {
			js1key = (js1key ^ 3084) & 0xFFFFFF
		}

		js1key = (js1key - 6341) & 0xFFFFFF

		if obf(56913)+len(locationHref) > obf(56972) {
			js1key = (js1key + 3769) & 0xFFFFFF
		}

		js1key = (js1key - 8999) & 0xFFFFFF

		if len(x1)+obf(56914) < obf(56935) {
			js1key = (js1key + 6920) & 0xFFFFFF
		}

		if obf(34892)+obf(48570) < obf(28130) {
			js1key = (js1key + 1264) & 0xFFFFFF
		}

		if animated {
			js1key = (js1key ^ 87291) & 0xFFFFFF
		}

		js1key = (js1key ^ 7249) & 0xFFFFFF

		if obf(60069)+obf(59777) < obf(27136) {
			js1key = (js1key - 3808) & 0xFFFFFF
		}

		js1key = (js1key + 2625) & 0xFFFFFF

		if animated {
			js1key = (js1key + 86379) & 0xFFFFFF
		}

		js1key = (js1key + 9192) & 0xFFFFFF
		js1key = (js1key - 3116) & 0xFFFFFF
		js1key = (js1key ^ 8580) & 0xFFFFFF

		if len(x1)+obf(46288) > obf(57344) {
			js1key = (js1key - 5467) & 0xFFFFFF
		}

		js1key = (js1key ^ 7636) & 0xFFFFFF
		js1key = (js1key + 6090) & 0xFFFFFF
		js1key = (js1key - 4085) & 0xFFFFFF

		if obf(35609)+obf(32327) < obf(25583) {
			js1key = (js1key - 3211) & 0xFFFFFF
		}

		js1key = (js1key + 6561) & 0xFFFFFF
		js1key = (js1key - 3687) & 0xFFFFFF

		if obf(56910)+len(locationHref) > obf(56972) {
			js1key = (js1key + 6366) & 0xFFFFFF
		}
	}

	// Process eighth chunk of x1 if available
	if len(x1) >= 32 {
		hexVal, _ := strconv.ParseInt(x1[28:32], 16, 64)
		js1key += int(hexVal) + 12
		js1key &= 0xFFFFFF
		js1key = (js1key - 1608) & 0xFFFFFF
		js1key = (js1key ^ 3084) & 0xFFFFFF

		if animated {
			js1key = (js1key - 89270) & 0xFFFFFF
		}

		js1key = (js1key - 6341) & 0xFFFFFF
		js1key = (js1key + 3769) & 0xFFFFFF

		if obf(56913)+len(locationHref) > obf(56913) {
			js1key = (js1key - 8999) & 0xFFFFFF
		}

		js1key = (js1key + 6920) & 0xFFFFFF
		js1key = (js1key + 1264) & 0xFFFFFF
		js1key = (js1key ^ 7249) & 0xFFFFFF

		if animated {
			js1key = (js1key - 51206) & 0xFFFFFF
		}

		js1key = (js1key - 3808) & 0xFFFFFF
		js1key = (js1key + 2625) & 0xFFFFFF

		if animated {
			js1key = (js1key + 39367) & 0xFFFFFF
		}

		js1key = (js1key + 9192) & 0xFFFFFF

		if obf(46156)+obf(45907) < obf(13117) {
			js1key = (js1key - 3116) & 0xFFFFFF
		}

		if obf(56913)+len(locationHref) > obf(56911) {
			js1key = (js1key ^ 8580) & 0xFFFFFF
		}

		js1key = (js1key - 5467) & 0xFFFFFF

		if animated {
			js1key = (js1key ^ 25650) & 0xFFFFFF
		}

		js1key = (js1key ^ 7636) & 0xFFFFFF
		js1key = (js1key + 6090) & 0xFFFFFF
		js1key = (js1key - 4085) & 0xFFFFFF
		js1key = (js1key - 3211) & 0xFFFFFF
		js1key = (js1key + 6561) & 0xFFFFFF
		js1key = (js1key - 3687) & 0xFFFFFF
		js1key = (js1key + 6366) & 0xFFFFFF
	}

	return js1key
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
	animated := false                 // 页面动画状态固定为关闭
	locationHref := constants.BaseURL // 使用基础URL作为locationHref参数
	js1key := calculateJs1Key(x1Value, locationHref, animated)

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
