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

### 为 CDN 单域名推送新证书

```bash
./autoUCDNcert -domainID="域名资源id" -certName="证书名" -certPath="证书文件路径" -keyPath="密钥文件路径"
```

### 为 CDN 多域名推送新证书

```bash
./autoUCDNcert -certName="证书名" -certPath="证书文件路径" -keyPath="密钥文件路径"
```
