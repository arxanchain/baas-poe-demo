# 数字资产存证Demo

这是一个用来演示如何使用阿尔山区块链BaaS平台进行数字资产存证(POE)的网页应用Demo。

首先，clone下来Demo的代码：

```
git clone https://github.com/arxanchain/baas-poe-demo.git
```

注意，clone下来的`baas-poe-demo`工程位置必须在`$GOPATH/github.com/arxanchain`目录下，否则不能正常编译运行。

## 修改Demo配置

在编译运行Demo之前，需要根据实际的环境来配置Demo的一些参数，打开`conf/app.conf`文件，可以看到Demo运行的配置参数。

主要的配置参数有：


```
[wallet]
## 阿尔山区块链BaaS平台对外API网关地址和端口，请根据实际环境进行配置
address = http://API-Gateway-IP:Port
## 开发者在Chain-Console管理页面上申请的API-Key，请根据实际申请的API-Key进行配置
apikey = Your-API-Key

[crypto]
## 是否启用API加密通道，默认需要启动
enable=true
## 启用API加密通道后，需要配置相关证书，请根据Chain-Console管理页面上实际申请的证书进行配置，提供证书所在的路径
certspath=path/to/certs/path

[signature]
## 企业开发者的钱包DID，如果是对接阿尔山公有链，需要使用该账户进行签名
creator=your-enterprise-wallet-did
## 签名用的随机值，可以保持默认值，也可以进行自定义设置
nonce=nonce
## 企业开发者的交易签名凭证(钱包私钥），如果是对接阿尔山公有链，需要使用该交易签名凭证进行签名
privatekey=your-enterprise-wallet-private-key
```

对于`certspath`的证书路径配置，请参考[wallet-sdk-go README](https://github.com/arxanchain/wallet-sdk-go/tree/v2.1.2)中的`Organize your client certificates`一小节。

其它不需要修改的配置参数，保持默认值即可。


## 如何运行Demo

在`baas-poe-demo`工程根目录下，运行下面命令编译Demo：

```
make
```

会在`baas-poe-demo`工程根目录下，生成一个命名为`baas-poe-demo`的二进制可执行文件。

然后运行该二进制文件，启动Demo：

```
./baas-poe-demo
```

最后提示在浏览器访问地址： http://localhost:8081 

Enjoy!

