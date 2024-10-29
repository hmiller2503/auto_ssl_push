package ucloud

import (
	"autoUCDNcert/config"
	"autoUCDNcert/utils"
	"fmt"
	"os"
	"time"

	"github.com/ucloud/ucloud-sdk-go/services/ucdn"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
	"github.com/ucloud/ucloud-sdk-go/ucloud/auth"
)

type UCloudClient struct {
	BaseURL    string
	PublicKey  string
	PrivateKey string
	UcdnClient *ucdn.UCDNClient
	LogFile    *os.File
}

func NewClient(logFile *os.File) *UCloudClient {
	cfg := ucloud.NewConfig()
	cfg.BaseUrl = config.AppConfig.BaseURL

	cred := auth.NewCredential()
	cred.PublicKey = config.AppConfig.PublicKey
	cred.PrivateKey = config.AppConfig.PrivateKey

	ucdnClient := ucdn.NewClient(&cfg, &cred)

	return &UCloudClient{
		BaseURL:    config.AppConfig.BaseURL,
		PublicKey:  cred.PublicKey,
		PrivateKey: cred.PrivateKey,
		UcdnClient: ucdnClient,
		LogFile:    logFile,
	}
}

func (client *UCloudClient) LogOperation(message string) {
	utils.LogOperation(client.LogFile, message)
}

func (client *UCloudClient) HandleFailure(operation string, err error) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("操作 %s 失败: %v\n", operation, err)
	utils.SendNotify(fmt.Sprint("操作: ", operation, " 执行失败\n 原因: ", err, "\n 服务器时间:", timestamp, " 请注意！"))
}
func StartProgramNotify() {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	utils.SendNotify(fmt.Sprint("自动同步SSL至UCLOUD程序完美结束运行 服务器时间: ", timestamp))
}
func EndProgramNotify() {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	utils.SendNotify(fmt.Sprint("自动同步SSL至UCLOUD程序完美结束运行 服务器时间: ", timestamp))
}
