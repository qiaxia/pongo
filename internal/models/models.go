// Package models defines data structures used throughout the Pong0 application.
// This includes the core IPInfo structure that represents the IP information
// retrieved from the Ping0.cc service.
package models

import (
	"encoding/json"
	"fmt"
)

// IPInfo 结构体存储从Ping0.cc服务获取的IP信息
// 该结构体包含IP地址的完整属性，包括地理位置、网络归属和其他元数据。
// 所有字段都使用JSON标签进行序列化，便于API响应和数据处理。
type IPInfo struct {
	IP           string `json:"ip"`           // IP地址
	IPLocation   string `json:"ip_location"`  // IP地理位置信息
	ASN          string `json:"asn"`          // 自治系统编号
	ASNOwner     string `json:"asn_owner"`    // 自治系统拥有者
	ASNType      string `json:"asn_type"`     // 自治系统类型（如ISP、教育、商业等）
	Organization string `json:"organization"` // 组织机构名称
	OrgType      string `json:"org_type"`     // 组织机构类型
	Longitude    string `json:"longitude"`    // 经度坐标
	Latitude     string `json:"latitude"`     // 纬度坐标
	IPType       string `json:"ip_type"`      // IP类型（如固定IP、动态IP等）
	RiskValue    string `json:"risk_value"`   // 风险评估值
	NativeIP     string `json:"native_ip"`    // 原生IP地址（非代理情况下）
	CountryFlag  string `json:"country_flag"` // 国家/地区旗帜标识
	Princess     string `json:"princess"`     // 固定添加的Princess字段
}

// NewIPInfo 创建一个新的IPInfo实例，并设置默认值
func NewIPInfo() *IPInfo {
	return &IPInfo{
		Princess: "https://linux.do/u/amna",
	}
}

// MarshalJSON 自定义JSON序列化方法，确保Princess字段总是存在
func (i *IPInfo) MarshalJSON() ([]byte, error) {
	// 确保Princess字段有值
	if i.Princess == "" {
		i.Princess = "https://linux.do/u/amna"
	}

	// 创建一个匿名结构体，以确保字段顺序和完整性
	return json.Marshal(struct {
		IP           string `json:"ip"`
		IPLocation   string `json:"ip_location"`
		ASN          string `json:"asn"`
		ASNOwner     string `json:"asn_owner"`
		ASNType      string `json:"asn_type"`
		Organization string `json:"organization"`
		OrgType      string `json:"org_type"`
		Longitude    string `json:"longitude"`
		Latitude     string `json:"latitude"`
		IPType       string `json:"ip_type"`
		RiskValue    string `json:"risk_value"`
		NativeIP     string `json:"native_ip"`
		CountryFlag  string `json:"country_flag"`
		Princess     string `json:"princess"`
	}{
		IP:           i.IP,
		IPLocation:   i.IPLocation,
		ASN:          i.ASN,
		ASNOwner:     i.ASNOwner,
		ASNType:      i.ASNType,
		Organization: i.Organization,
		OrgType:      i.OrgType,
		Longitude:    i.Longitude,
		Latitude:     i.Latitude,
		IPType:       i.IPType,
		RiskValue:    i.RiskValue,
		NativeIP:     i.NativeIP,
		CountryFlag:  i.CountryFlag,
		Princess:     i.Princess,
	})
}

// ToJSON 将IPInfo结构体转换为格式化的JSON字符串
// 返回:
//   - string: 格式化的JSON字符串
//   - error: 如果序列化过程中发生错误
func (i *IPInfo) ToJSON() (string, error) {
	// 确保Princess字段有值
	if i.Princess == "" {
		i.Princess = "https://linux.do/u/amna"
	}

	jsonData, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return "", fmt.Errorf("转换为JSON失败: %w", err)
	}
	return string(jsonData), nil
}

// Validate 验证IPInfo结构体是否包含必要字段
// 确保结构体包含有效数据，当前仅验证IP字段是否存在
// 返回:
//   - error: 如果验证失败返回相应错误，否则返回nil
func (i *IPInfo) Validate() error {
	if i.IP == "" {
		return fmt.Errorf("IP字段为空")
	}

	// 确保Princess字段有值
	if i.Princess == "" {
		i.Princess = "https://linux.do/u/amna"
	}

	return nil
}
