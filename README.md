# domain-admin-ssl-deploy

支持[Domain Admin](https://github.com/mouday/domain-admin)项目证书自定义部署,用python写了一版，不过不同主机部署的话有些麻烦，go的话部署管理方便些

- 通过控制3个header字段进行自定义部署
- 配置文件进行命令映射防止危险命令
- 通过请求头选择执行哪个命令

## 使用

打包及运行

```sh
CGO_ENABLED=0  GOOS=linux  GOARCH=amd64  go build main.go
chmod a+x main

./main
```

## 测试

```sh
curl --location 'http://localhost:55433/issueCertificate' \
--header 'Token: vqQ2uVaAFj2DXDbbicw7' \
--header 'Key-Save-Path: .' \
--header 'Deploy-Cmd: cmd_simple_nginx' \
--header 'Content-Type: text/plain' \
--data '{
  "domains": [
      "www.baidu.com",
      "zhidao.baidu.com"
  ],
  "ssl_certificate":"-----BEGIN CERTIFICATE-----\nMIIGdTCCBN2gAwI\n-----END CERTIFICATE-----",
  "ssl_certificate_key":"-----BEGIN PRIVATE KEY-----\nMIIEvH+bpTwI=\n-----END PRIVATE KEY-----",
  "start_time": "2023-01-04 14:33:39",
  "expire_time": "2023-04-04 14:33:39"
}
'
```

```sh
http://localhost:55433/issueCertificate


{
  "Token": "vqQ2uVaAFj2DXDbbicw7",
  "Key-Save-Path": "/data/safeline/resources/nginx/certs/",
  "Deploy-Cmd": "cmd_waf_base_02"
}
```
