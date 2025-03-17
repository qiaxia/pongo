// Package parser implements parsing and extraction functionality for the Pong0 application.
// It handles HTML parsing, data extraction, and JavaScript execution to retrieve
// IP information from the Ping0.cc service responses.
package parser

import (
	"fmt"
	"html"
	"regexp"
	"strings"
	"sync"

	"ping0/internal/constants"
	"ping0/internal/models"

	"github.com/PuerkitoBio/goquery"
)

// 正则表达式编译缓存，用于提高性能
var (
	extractRegexCache = make(map[string]*regexp.Regexp)
	regexCacheMutex   sync.RWMutex
)

// getOrCompileRegex 获取或编译正则表达式
// 该函数实现了正则表达式的缓存机制，避免重复编译相同的正则表达式。
// 它使用读写锁保证并发安全。
//
// 参数:
//   - pattern: 正则表达式模式字符串
//
// 返回:
//   - *regexp.Regexp: 编译好的正则表达式对象
func getOrCompileRegex(pattern string) *regexp.Regexp {
	regexCacheMutex.RLock()
	re, exists := extractRegexCache[pattern]
	regexCacheMutex.RUnlock()

	if exists {
		return re
	}

	re = regexp.MustCompile(pattern)

	regexCacheMutex.Lock()
	extractRegexCache[pattern] = re
	regexCacheMutex.Unlock()

	return re
}

// ParseIPInfo 从HTML内容中解析IP信息
// 该函数分析从Ping0.cc服务获取的HTML响应，提取所有相关的IP信息，
// 并将其组织到IPInfo结构体中。
//
// 参数:
//   - htmlContent: 包含IP信息的HTML内容
//
// 返回:
//   - *models.IPInfo: 解析出的IP信息结构体
//   - error: 如果解析失败则返回相应错误
func ParseIPInfo(htmlContent string) (*models.IPInfo, error) {
	// 检查输入参数
	if htmlContent == "" {
		return nil, fmt.Errorf("HTML内容为空")
	}

	// 检查是否包含错误信息
	if strings.Contains(htmlContent, "系统发生错误") {
		// 尝试提取更详细的错误信息
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
		if err == nil {
			errorMsg := doc.Find(".error-message").Text()
			if errorMsg != "" {
				return nil, fmt.Errorf("网站返回错误: %s", errorMsg)
			}
		}
		return nil, fmt.Errorf("网站返回错误: 系统发生错误")
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("解析HTML失败: %w", err)
	}

	// 检查网站错误页面
	title := doc.Find("title").Text()
	if strings.Contains(title, "系统发生错误") || strings.Contains(title, "Error") {
		return nil, fmt.Errorf("网站返回错误页面: %s", title)
	}

	ipInfo := models.NewIPInfo()

	// 从脚本标签中直接提取常用变量
	scriptValues := extractScriptVariables(doc)
	if constants.Verbose && len(scriptValues) > 0 {
		fmt.Println("从脚本中提取的变量:")
		for k, v := range scriptValues {
			fmt.Printf("- %s: %s\n", k, v)
		}
	}

	// 设置IP
	if ip, ok := scriptValues["window.ip"]; ok && ip != "" {
		ipInfo.IP = ip
		if constants.Verbose {
			fmt.Printf("从脚本中提取到IP: %s\n", ip)
		}
	} else {
		// 备选方法：从title中提取
		ipParts := strings.Split(title, "-")
		if len(ipParts) > 0 {
			ipInfo.IP = strings.TrimSpace(ipParts[0])
			if constants.Verbose {
				fmt.Printf("从标题中提取到IP: %s\n", ipInfo.IP)
			}
		}
	}

	// 如果无法提取到IP，页面可能是错误页面
	if ipInfo.IP == "" {
		// 打印HTML内容的前200个字符以便调试
		if constants.Verbose {
			preview := htmlContent
			if len(preview) > 200 {
				preview = preview[:200] + "..."
			}
			fmt.Printf("HTML内容预览: %s\n", preview)
		}
		return nil, fmt.Errorf("无法从页面提取IP信息，可能是错误页面")
	}

	// 设置IP位置
	if loc, ok := scriptValues["window.loc"]; ok && loc != "" {
		// 解码HTML实体
		ipInfo.IPLocation = decodeHTMLEntities(loc)
		if constants.Verbose {
			fmt.Printf("从脚本中提取到位置: %s\n", ipInfo.IPLocation)
		}
	} else {
		// 备选方法：从DOM中提取
		extractIPLocation(doc, ipInfo)
		if constants.Verbose && ipInfo.IPLocation != "" {
			fmt.Printf("从DOM中提取到位置: %s\n", ipInfo.IPLocation)
		}
	}

	// 提取国家旗帜
	doc.Find(".line.loc .content img").Each(func(i int, s *goquery.Selection) {
		flagSrc, exists := s.Attr("src")
		if exists {
			parts := strings.Split(flagSrc, "/")
			if len(parts) > 0 {
				flagFile := parts[len(parts)-1]
				ipInfo.CountryFlag = strings.TrimSuffix(flagFile, ".png")
				if constants.Verbose {
					fmt.Printf("提取到国家旗帜: %s\n", ipInfo.CountryFlag)
				}
			}
		}
	})

	// 提取ASN
	doc.Find(".line.asn .content a").Each(func(i int, s *goquery.Selection) {
		ipInfo.ASN = strings.TrimSpace(s.Text())
		if constants.Verbose && ipInfo.ASN != "" {
			fmt.Printf("提取到ASN: %s\n", ipInfo.ASN)
		}
	})

	// 提取ASN所有者和类型
	extractASNInfo(doc, scriptValues, ipInfo)
	if constants.Verbose {
		if ipInfo.ASNOwner != "" {
			fmt.Printf("提取到ASN所有者: %s\n", ipInfo.ASNOwner)
		}
		if ipInfo.ASNType != "" {
			fmt.Printf("提取到ASN类型: %s\n", ipInfo.ASNType)
		}
	}

	// 提取组织信息和类型
	extractOrgInfo(doc, scriptValues, ipInfo)
	if constants.Verbose {
		if ipInfo.Organization != "" {
			fmt.Printf("提取到组织: %s\n", ipInfo.Organization)
		}
		if ipInfo.OrgType != "" {
			fmt.Printf("提取到组织类型: %s\n", ipInfo.OrgType)
		}
	}

	// 提取经度
	if longitude, ok := scriptValues["window.longitude"]; ok && longitude != "" {
		ipInfo.Longitude = longitude
		if constants.Verbose {
			fmt.Printf("提取到经度: %s\n", longitude)
		}
	} else {
		doc.Find(".line").Each(func(i int, s *goquery.Selection) {
			name := strings.TrimSpace(s.Find(".name").Text())
			if name == "经度" {
				ipInfo.Longitude = strings.TrimSpace(s.Find(".content").Text())
				if constants.Verbose {
					fmt.Printf("从DOM中提取到经度: %s\n", ipInfo.Longitude)
				}
			}
		})
	}

	// 提取纬度
	if latitude, ok := scriptValues["window.latitude"]; ok && latitude != "" {
		ipInfo.Latitude = latitude
		if constants.Verbose {
			fmt.Printf("提取到纬度: %s\n", latitude)
		}
	} else {
		doc.Find(".line").Each(func(i int, s *goquery.Selection) {
			name := strings.TrimSpace(s.Find(".name").Text())
			if name == "纬度" {
				ipInfo.Latitude = strings.TrimSpace(s.Find(".content").Text())
				if constants.Verbose {
					fmt.Printf("从DOM中提取到纬度: %s\n", ipInfo.Latitude)
				}
			}
		})
	}

	// 提取IP类型 - 收集所有类型并用分号分隔
	extractIPTypes(doc, ipInfo)
	if constants.Verbose && ipInfo.IPType != "" {
		fmt.Printf("提取到IP类型: %s\n", ipInfo.IPType)
	}

	// 提取风控值
	doc.Find(".line.line-risk .content .riskbar .riskcurrent").Each(func(i int, s *goquery.Selection) {
		value := strings.TrimSpace(s.Find(".value").Text())
		lab := strings.TrimSpace(s.Find(".lab").Text())
		if value != "" && lab != "" {
			ipInfo.RiskValue = value + " " + lab
			if constants.Verbose {
				fmt.Printf("提取到风控值: %s\n", ipInfo.RiskValue)
			}
		}
	})

	// 提取原生IP
	doc.Find(".line.line-nativeip .content .label").Each(func(i int, s *goquery.Selection) {
		ipInfo.NativeIP = strings.TrimSpace(s.Text())
		if constants.Verbose {
			fmt.Printf("提取到原生IP: %s\n", ipInfo.NativeIP)
		}
	})

	// 验证结果
	if ipInfo.IP == "" {
		return nil, fmt.Errorf("未能提取到IP信息")
	}

	// 返回前确保Princess字段有值
	if ipInfo.Princess == "" {
		ipInfo.Princess = "https://linux.do/u/amna"
	}

	return ipInfo, nil
}

// extractScriptVariables 从脚本标签中提取变量
func extractScriptVariables(doc *goquery.Document) map[string]string {
	scriptValues := make(map[string]string)

	varNames := []string{
		"window.ip", "window.tar", "window.longitude", "window.latitude", "window.loc",
	}

	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		content := s.Text()
		for _, varName := range varNames {
			extractScriptVar(content, varName, &scriptValues)
		}
	})

	return scriptValues
}

// extractIPLocation 从DOM中提取IP位置信息
func extractIPLocation(doc *goquery.Document, ipInfo *models.IPInfo) {
	doc.Find(".line.loc .content").Each(func(i int, s *goquery.Selection) {
		// 获取原始HTML和文本
		html, _ := s.Html()

		// 使用更强的正则表达式清理HTML和无用文本
		text := extractTextBetweenTags(html)

		// 移除"错误提交"文本
		text = strings.Replace(text, "错误提交", "", -1)

		// 清理空格
		text = strings.TrimSpace(text)
		re := getOrCompileRegex(`\s+`)
		text = re.ReplaceAllString(text, " ")

		// 解码HTML实体
		if text != "" {
			ipInfo.IPLocation = decodeHTMLEntities(text)
		}
	})
}

// extractIPTypes 提取IP类型
func extractIPTypes(doc *goquery.Document, ipInfo *models.IPInfo) {
	var ipTypes []string
	doc.Find(".line.line-iptype .content .label").Each(func(i int, s *goquery.Selection) {
		ipType := strings.TrimSpace(s.Text())
		if ipType != "" {
			ipTypes = append(ipTypes, ipType)
		}
	})
	// 用分号连接所有IP类型
	if len(ipTypes) > 0 {
		ipInfo.IPType = strings.Join(ipTypes, "; ")
	}
}

// extractASNInfo 提取ASN所有者和类型
func extractASNInfo(doc *goquery.Document, scriptValues map[string]string, ipInfo *models.IPInfo) {
	doc.Find(".line.asnname .content").Each(func(i int, s *goquery.Selection) {
		// 保存原始选择器，用于后续提取标签
		original := s.Clone()

		// 跳过标签元素，直接获取纯文本内容
		// 移除掉标签元素
		s.Find(".label").Each(func(i int, label *goquery.Selection) {
			label.Remove()
		})

		// 获取剩余内容
		content := strings.TrimSpace(s.Text())

		// 移除连字符和后面的内容
		if dashIndex := strings.Index(content, "—"); dashIndex != -1 {
			content = strings.TrimSpace(content[:dashIndex])
		}

		// 应用HTML实体解码
		ipInfo.ASNOwner = decodeHTMLEntities(content)

		// 提取ASN类型 - 收集所有类型并用分号分隔
		var asnTypes []string
		original.Find(".label").Each(func(i int, label *goquery.Selection) {
			asnType := strings.TrimSpace(label.Text())
			if asnType != "" {
				asnTypes = append(asnTypes, asnType)
			}
		})
		// 用分号连接所有ASN类型
		if len(asnTypes) > 0 {
			ipInfo.ASNType = strings.Join(asnTypes, "; ")
		}
	})
}

// extractOrgInfo 提取组织信息和类型
func extractOrgInfo(doc *goquery.Document, scriptValues map[string]string, ipInfo *models.IPInfo) {
	doc.Find(".line.orgname .content").Each(func(i int, s *goquery.Selection) {
		// 保存原始选择器，用于后续提取标签
		original := s.Clone()

		// 跳过标签元素，直接获取纯文本内容
		// 移除掉标签元素
		s.Find(".label").Each(func(i int, label *goquery.Selection) {
			label.Remove()
		})

		// 获取剩余内容
		content := strings.TrimSpace(s.Text())

		// 移除连字符和后面的内容
		if dashIndex := strings.Index(content, "—"); dashIndex != -1 {
			content = strings.TrimSpace(content[:dashIndex])
		}

		// 应用HTML实体解码
		ipInfo.Organization = decodeHTMLEntities(content)

		// 提取组织类型 - 收集所有类型并用分号分隔
		var orgTypes []string
		original.Find(".label").Each(func(i int, label *goquery.Selection) {
			orgType := strings.TrimSpace(label.Text())
			if orgType != "" {
				orgTypes = append(orgTypes, orgType)
			}
		})
		// 用分号连接所有组织类型
		if len(orgTypes) > 0 {
			ipInfo.OrgType = strings.Join(orgTypes, "; ")
		}
	})
}

// extractScriptVar 从脚本内容中提取变量
func extractScriptVar(content, varName string, result *map[string]string) {
	if strings.Contains(content, varName) {
		// 构建正则表达式模式
		pattern := fmt.Sprintf(`%s\s*=\s*['"]([^'"]*)['"]\s*;`, regexp.QuoteMeta(varName))
		re := getOrCompileRegex(pattern)

		matches := re.FindStringSubmatch(content)
		if len(matches) > 1 {
			(*result)[varName] = matches[1]
		}
	}
}

// extractTextBetweenTags 提取HTML标签之间的文本内容
func extractTextBetweenTags(html string) string {
	// 移除所有HTML标签
	re := getOrCompileRegex("<[^>]*>")
	text := re.ReplaceAllString(html, " ")

	// 移除Vue模板表达式
	re = getOrCompileRegex("{{[^}]*}}")
	text = re.ReplaceAllString(text, "")

	// 清理多余空格
	re = getOrCompileRegex(`\s+`)
	text = re.ReplaceAllString(text, " ")

	return strings.TrimSpace(text)
}

// cleanIPInfo 清理IPInfo中的Vue模板表达式
func cleanIPInfo(info *models.IPInfo) {
	// 检查并清理IP
	if strings.Contains(info.IP, "{{") || strings.Contains(info.IP, "}}") {
		info.IP = ""
	}
}

// decodeHTMLEntities 解码HTML实体为正确的Unicode字符
func decodeHTMLEntities(text string) string {
	// 优先使用标准库解码
	decoded := html.UnescapeString(text)

	// 处理一些特殊情况
	customReplacements := map[string]string{
		"\\u0026mdash;": "—",
		"\\u0026#8212;": "—",
	}

	for entity, replacement := range customReplacements {
		if strings.Contains(decoded, entity) {
			decoded = strings.ReplaceAll(decoded, entity, replacement)
		}
	}

	return decoded
}
