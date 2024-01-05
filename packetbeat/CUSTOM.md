# 自定义文档

## 初始化环境

### go 环境

安装 go 环境(版本`>=1.20`)

```shell
go mod tidy
```

### linux 环境(debian/ubuntu)

```shell
# 安装pcap开发依赖库
apt-get install libpcap-dev
```

## 授权

授权文件位置：`packetbeat/license.dat`

- expired_at 为`null`时，表示不过期
- issued_at 为授权时间，必须早于当前时间
- 目前对所有命令访问进行授权验证拦截，如果授权验证不通过，拒绝除`-h 帮助`之外的所有命令

## 构建

- 编辑`libbeat/cmd/license.go`添加`license key`用于解密

```shell
# 在项目根目录下构建，即：beats目录下
go build -o build/packet-audit  packetbeat/main.go
```

## 运行

```shell
# 在项目根目录下运行，即：beats目录下
build/packet-audit -c packetbeat/packetbeat.yml --license packetbeat/license.dat --strict.perms=false
```

> 需要 `root` 身份运行，否则会报权限错误
