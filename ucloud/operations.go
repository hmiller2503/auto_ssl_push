package ucloud

import (
	"autoUCDNcert/config"
	"autoUCDNcert/utils"

	"fmt"
	"net/http"

	"github.com/ucloud/ucloud-sdk-go/ucloud"
)

func (client *UCloudClient) DeleteCert(certName string) error {
	req := client.UcdnClient.NewDeleteCertificateRequest()
	req.ProjectId = ucloud.String("org-ixkwze")
	req.CertName = ucloud.String(certName)

	resp, err := client.UcdnClient.DeleteCertificate(req)
	if err != nil {
		client.LogOperation(fmt.Sprintf("删除证书错误: %v", err))
		client.HandleFailure("DeleteCert", err)
		return err
	}

	client.LogOperation(fmt.Sprintf("删除证书请求参数: CertName=%s", certName))
	client.LogOperation(fmt.Sprintf("删除证书成功: %v", resp))

	if resp.RetCode != 0 {
		client.LogOperation(fmt.Sprintf("DeleteCert 操作失败, RetCode: %d", resp.RetCode))
		client.HandleFailure("DeleteCert", fmt.Errorf("RetCode: %d", resp.RetCode))
		return fmt.Errorf("RetCode: %d", resp.RetCode)
	}
	return nil
}

func (client *UCloudClient) GetDomainIDsForCert(certName string) ([]string, error) {
	certReq := client.UcdnClient.NewGetCertificateV2Request()
	certReq.ProjectId = ucloud.String("org-ixkwze")

	certResp, err := client.UcdnClient.GetCertificateV2(certReq)
	if err != nil {
		return nil, fmt.Errorf("获取证书列表失败: %v", err)
	}

	var domains []string
	for _, cert := range certResp.CertList {
		if cert.CertName == certName {
			domains = cert.Domains
			break
		}
	}

	// 如果未找到匹配的证书名称，返回空的域名列表
	if len(domains) == 0 {
		client.LogOperation(fmt.Sprintf("未找到匹配的证书名称: %s，返回空域名列表", certName))
		return []string{}, nil
	}

	domainReq := client.UcdnClient.NewGetUcdnDomainInfoListRequest()
	domainReq.ProjectId = ucloud.String("org-ixkwze")

	domainResp, err := client.UcdnClient.GetUcdnDomainInfoList(domainReq)
	if err != nil {
		return nil, fmt.Errorf("获取域名列表失败: %v", err)
	}

	var domainIDs []string
	for _, domain := range domainResp.DomainInfoList {
		for _, d := range domains {
			if domain.Domain == d {
				domainIDs = append(domainIDs, domain.DomainId)
			}
		}
	}

	// 如果域名 ID 列表为空，返回空列表
	return domainIDs, nil
}

func (client *UCloudClient) ChangeHttps(domainID, httpsStatus, certName string) error {
	client.LogOperation(fmt.Sprint("开始请求变更HTTPS状态 domainID: ", domainID, " httpsStatus: ", httpsStatus, " certName: ", certName))
	params := map[string]string{
		"Action":        "UpdateUcdnDomainHttpsConfigV2",
		"Region":        config.AppConfig.Region,
		"PublicKey":     client.PublicKey,
		"DomainId":      domainID,
		"HttpsStatusCn": httpsStatus,
	}

	if httpsStatus != "disable" {
		params["CertName"] = certName
	}

	signature := utils.GenerateSignature(params, client.PrivateKey)
	params["Signature"] = signature

	req, err := http.NewRequest("GET", client.BaseURL, nil)
	if err != nil {
		client.LogOperation(fmt.Sprintf("创建请求失败: %v", err))
		client.HandleFailure("ChangeHttps", err)
		return err
	}

	query := req.URL.Query()
	for k, v := range params {
		query.Add(k, v)
	}
	req.URL.RawQuery = query.Encode()

	clientHTTP := &http.Client{}
	resp, err := clientHTTP.Do(req)
	if err != nil {
		client.LogOperation(fmt.Sprintf("请求失败: %v", err))
		client.HandleFailure("ChangeHttps", err)
		return err
	}
	defer resp.Body.Close()

	client.LogOperation(fmt.Sprintf("请求 URL: %s", req.URL.String()))
	client.LogOperation(fmt.Sprintf("请求响应状态: %s", resp.Status))
	client.LogOperation(fmt.Sprint("请求变更HTTPS状态结束 domainID: ", domainID, " httpsStatus: ", httpsStatus, " certName: ", certName))
	return nil
}
func (client *UCloudClient) UploadCert(certName, userCert, userKey string) error {
	req := client.UcdnClient.NewAddCertificateRequest()
	req.ProjectId = ucloud.String("org-ixkwze")
	req.CertName = ucloud.String(certName)
	req.UserCert = ucloud.String(userCert)
	req.PrivateKey = ucloud.String(userKey)

	resp, err := client.UcdnClient.AddCertificate(req)
	if err != nil {
		client.LogOperation(fmt.Sprintf("上传证书错误: %v", err))
		client.HandleFailure("UploadCert", err)
		return err
	}

	client.LogOperation(fmt.Sprintf("上传证书请求参数: CertName=%s", certName))
	client.LogOperation(fmt.Sprintf("上传证书成功: %v", resp))

	if resp.RetCode != 0 {
		client.LogOperation(fmt.Sprintf("UploadCert 操作失败, RetCode: %d", resp.RetCode))
		client.HandleFailure("UploadCert", fmt.Errorf("RetCode: %d", resp.RetCode))
		return fmt.Errorf("RetCode: %d", resp.RetCode)
	}
	return nil
}
