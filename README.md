# auto_ssl_push

1Panel 面板自动推送 SSL 证书至 CDN 提供商
目前仅支持 UCLOUD CDN

## 配置文件

配置文件为项目根目录下的 **config.yml** 具体属性解释如下
|属性|描述|类型|
|-|-|-|
|base_url|ucloudApi 地址|string|
|public_key|Api 公钥|string|
|private_key|Api 密钥|string|
|region|地域节点|string|
|notify_type|通知模式|string|
|email_host|发件邮箱 SMTP 地址|string|
|email_port|SMTP 端口|string|
|email_username|发件邮箱|string|
|email_pwd|邮箱密码|string|
|email_receiver|接收邮箱|string|
|wx_webhook_url|企业微信 WebHook 地址|string|

## 使用方法

### 为 CDN 单域名推送新证书1

```bash
./autoUCDNcert -domainID="域名资源id" -certName="证书名" -certPath="证书文件路径" -keyPath="密钥文件路径"
```

### 为 CDN 多域名推送新证书

```bash
./autoUCDNcert -certName="证书名" -certPath="证书文件路径" -keyPath="密钥文件路径"
```

### 1Panel 面板配置简介

- 开启"推送证书到本地目录"
- 开启"申请证书之后执行脚本"

```bash
#脚本内容中路径仅供参考，应以自己设置的实际路径为准
cd 本项目的目录 && sudo 本项目的目录/打包后的项目文件 -certName="证书名" -certPath="证书文件路径" -keyPath="密钥文件路径"
#如果与长亭雷池WAF配合使用，则可以参考如下进行配置
sudo cp -f 证书文件路径  /data/safeline/resources/nginx/certs/cert_1.crt # cert_1.crt中的cert_1是根据在雷池WAF中添加证书时从1向上递增的
sudo cp -f 密钥文件路径  /data/safeline/resources/nginx/certs/cert_1.key # 同如上解释
sudo docker restart safeline-tengine
sudo docker restart safeline-mgt
```
