package config

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"gopkg.in/yaml.v2"
)

type UCloudConfig struct {
	BaseURL       string `yaml:"base_url"`
	PublicKey     string `yaml:"public_key"`
	PrivateKey    string `yaml:"private_key"`
	Region        string `yaml:"region"`
	NotifyType    string `yaml:"notify_type"` // 通知类型: "email" 或 "wx_webhook"
	EmailHost     string `yaml:"email_host"`  // Email 配置
	EmailPort     string `yaml:"email_port"`
	EmailUsername string `yaml:"email_username"`
	EmailPwd      string `yaml:"email_pwd"`
	EmailReceiver string `yaml:"email_receiver"`
	WxWebhookURL  string `yaml:"wx_webhook_url"` // 微信 Webhook URL
}

var AppConfig *UCloudConfig

// LoadConfig 读取并解析 YAML 配置文件
func LoadConfig(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("无法打开配置文件: %v", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	config := &UCloudConfig{}
	if err := decoder.Decode(config); err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 检查必要配置项是否缺失
	missingFields := checkRequiredFields(config)
	if len(missingFields) > 0 {
		return fmt.Errorf("配置文件缺少以下必要项: %v", missingFields)
	}

	AppConfig = config
	return nil
}

// checkRequiredFields 使用反射和自定义规则检查配置结构体中的必需字段是否为空
func checkRequiredFields(cfg *UCloudConfig) []string {
	var missingFields []string

	// 检查基本字段
	val := reflect.ValueOf(cfg).Elem()
	for _, fieldName := range []string{"BaseURL", "PublicKey", "PrivateKey", "Region"} {
		field := val.FieldByName(fieldName)
		if isZeroValue(field) {
			missingFields = append(missingFields, fieldName)
		}
	}

	// 检查通知相关字段
	switch cfg.NotifyType {
	case "email":
		for _, fieldName := range []string{"EmailHost", "EmailPort", "EmailUsername", "EmailPwd", "EmailReceiver"} {
			field := val.FieldByName(fieldName)
			if isZeroValue(field) {
				missingFields = append(missingFields, fieldName)
			}
		}
	case "wx_webhook":
		if isZeroValue(val.FieldByName("WxWebhookURL")) {
			missingFields = append(missingFields, "WxWebhookURL")
		}
	case "": // 不指定 NotifyType 则跳过通知配置的检查
		log.Println("未指定通知类型，通知配置将被忽略。")
	default: // 同时启用 email 和 wx_webhook 时，检查两者的字段
		for _, fieldName := range []string{"EmailHost", "EmailPort", "EmailUsername", "EmailPwd", "EmailReceiver", "WxWebhookURL"} {
			field := val.FieldByName(fieldName)
			if isZeroValue(field) {
				missingFields = append(missingFields, fieldName)
			}
		}
	}

	return missingFields
}

// isZeroValue 检查字段是否为零值
func isZeroValue(field reflect.Value) bool {
	return reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface())
}

// InitConfig 初始化并检查配置文件
func InitConfig(filePath string) {
	err := LoadConfig(filePath)
	if err != nil {
		log.Fatalf("配置文件错误: %v", err) // 记录日志并终止程序
	}
}
