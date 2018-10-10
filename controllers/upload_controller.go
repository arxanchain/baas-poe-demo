package controllers

import (
	"math/rand"
	"net/http"
	"time"

	pw "github.com/arxanchain/sdk-go-common/protos/wallet"
	restapi "github.com/arxanchain/sdk-go-common/rest/api"
	"github.com/arxanchain/sdk-go-common/structs"
	"github.com/arxanchain/sdk-go-common/structs/did"
	"github.com/arxanchain/sdk-go-common/structs/pki"
	"github.com/arxanchain/sdk-go-common/structs/wallet"
	walletapi "github.com/arxanchain/wallet-sdk-go/api"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type UploadController struct {
	beego.Controller
	walletClient wallet.IWalletClient
	walletID     did.Identifier
	keyPair      *wallet.KeyPair
}

const (
	mainPage = "index.html"
)

var (
	logger *logs.BeeLogger
)

func init() {
	logger = logs.NewLogger(100000)
	logger.SetLogger("console", "")
	logger.EnableFuncCallDepth(true)
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (c *UploadController) Post() {
	hash := c.GetString("hash")
	if hash == "" {
		c.Data["ErrMsg"] = "file hash is empty"
		//Goto error page
		c.TplName = mainPage
		return
	}
	logger.Debug("File hash: %v", hash)

	// New Wallet clinet
	if c.walletClient == nil {
		config := restapi.Config{
			Address:     beego.AppConfig.String("wallet::address"),
			ApiKey:      beego.AppConfig.String("wallet::apikey"),
			CallbackUrl: beego.AppConfig.String("wallet::callbackurl"),
		}
		logger.Debug("Wallet config: %+v", config)
		walletClient, err := walletapi.NewWalletClient(&config)
		if err != nil {
			logger.Error("New wallet client fail: %v", err)
			c.Data["ErrMsg"] = err.Error()
			//Goto error page
			c.TplName = mainPage
			return
		}
		c.walletClient = walletClient
	}

	header := http.Header{}
	header.Set(structs.InvokeModeHeader, structs.InvokeModeSync)

	// Query poe object via hash
	indexGetPayload := &wallet.IndexGetPayload{
		Indexs: &wallet.IndexTags{
			IndividualIndex: []string{
				hash,
			},
		},
	}

	IDs, err := c.walletClient.IndexGet(header, indexGetPayload)
	if err != nil {
		logger.Error("Index get fail: %v", err)
		c.Data["ErrMsg"] = err.Error()
		//Goto error page
		c.TplName = mainPage
		return
	}
	logger.Debug("Index get succ. IDs: %+v", IDs)

	// file hash already exist, query the existed POE infomation
	if len(IDs) > 0 {
		poePayload, err := c.walletClient.QueryPOE(header, did.Identifier(IDs[0]))
		if err != nil {
			logger.Error("Query POE fail: %v", err)
			c.Data["ErrMsg"] = err.Error()
			//Goto error page
			c.TplName = mainPage
			return
		}

		// Goto main page
		c.Data["notice"] = "该文件已经存证，详细信息如下"
		c.Data["poe_id"] = poePayload.Id
		c.Data["file_name"] = poePayload.OffchainMetadata.Filename
		c.Data["file_hash"] = poePayload.OffchainMetadata.ContentHash
		c.Data["file_size"] = poePayload.OffchainMetadata.Size
		c.Data["read_only"] = poePayload.OffchainMetadata.ReadOnly
		c.Data["created"] = time.Unix(poePayload.Created, 0).Format("2006-01-02 15:04:05")
		c.TplName = mainPage
		return
	}

	// Save poe file
	f, h, err := c.GetFile("poe_file")
	if err != nil {
		logger.Error("Get POE file fail: %v", err)
		c.Data["ErrMsg"] = err.Error()
		//Goto error page
		c.TplName = mainPage
		return
	}
	defer f.Close()
	// 保存位置在 static/upload, 没有文件夹要先创建
	poeFile := "static/upload/" + h.Filename
	err = c.SaveToFile("poe_file", poeFile)
	if err != nil {
		logger.Error("Save POE file fail: %v", err)
		c.Data["ErrMsg"] = err.Error()
		//Goto error page
		c.TplName = mainPage
		return
	}

	// New wallet account
	if c.walletID == "" {
		username := "username-" + RandStringRunes(5)
		registerBody := &wallet.RegisterWalletBody{
			Type:   pw.DidType_ORGANIZATION,
			Access: username,
			Secret: "Alice#123456",
		}
		resp, err := c.walletClient.Register(header, registerBody)
		if err != nil {
			logger.Error("Register wallet fail: %v", err)
			c.Data["ErrMsg"] = err.Error()
			//Goto error page
			c.TplName = mainPage
			return
		}
		logger.Debug("Register issuer wallet succ.\nResponse: %+v", resp)
		logger.Debug("Wallet public key: %s\n", resp.KeyPair.PublicKey)
		logger.Debug("Wallet private key: %s\n", resp.KeyPair.PrivateKey)

		c.walletID = resp.Id
		c.keyPair = resp.KeyPair
	}

	// Create POE
	poeBody := &wallet.POEBody{
		Name:  "TestPOE",
		Owner: c.walletID,
	}
	signParam := &pki.SignatureParam{
		Creator:    c.walletID,
		Nonce:      "nonce",
		PrivateKey: c.keyPair.PrivateKey,
	}
	poeResp, err := c.walletClient.CreatePOE(header, poeBody, signParam)
	if err != nil {
		logger.Error("Create POE fail: %v", err)
		c.Data["ErrMsg"] = err.Error()
		//Goto error page
		c.TplName = mainPage
		return
	}
	poeID := poeResp.Id

	uploadResp, err := c.walletClient.UploadPOEFile(header, string(poeID), poeFile, true)
	if err != nil {
		logger.Error("Upload POE file fail: %v", err)
		c.Data["ErrMsg"] = err.Error()
		//Goto error page
		c.TplName = mainPage
		return
	}
	logger.Debug("Upload POE file succ.\nResponse: %+v\nOffchainMetadata: %+v", uploadResp, *uploadResp.OffchainMetadata)

	// Goto main page
	c.Data["notice"] = "文件存证成功，详细信息如下"
	c.Data["poe_id"] = uploadResp.Id
	c.Data["file_name"] = uploadResp.OffchainMetadata.Filename
	c.Data["file_hash"] = uploadResp.OffchainMetadata.ContentHash
	c.Data["file_size"] = uploadResp.OffchainMetadata.Size
	c.Data["read_only"] = uploadResp.OffchainMetadata.ReadOnly
	c.Data["created"] = time.Now().Format("2006-02-01 15:04:05")

	if len(uploadResp.TransactionIds) > 0 {
		c.Data["transaction_id"] = uploadResp.TransactionIds[0]
	}
	c.TplName = mainPage
}

func (c *UploadController) Check() {
	hash := c.GetString("hash")
	result := make(map[string]interface{})
	if hash == "" {
		result["ErrMsg"] = "file hash is empty"
		//Goto error page
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	logger.Debug("File hash: %v", hash)

	// New Wallet clinet
	if c.walletClient == nil {
		config := restapi.Config{
			Address:     beego.AppConfig.String("wallet::address"),
			ApiKey:      beego.AppConfig.String("wallet::apikey"),
			CallbackUrl: beego.AppConfig.String("wallet::callbackurl"),
		}
		logger.Debug("Wallet config: %+v", config)
		walletClient, err := walletapi.NewWalletClient(&config)
		if err != nil {
			logger.Error("New wallet client fail: %v", err)
			result["ErrMsg"] = err.Error()
			//Goto error page
			c.Data["json"] = result
			c.ServeJSON()
			return
		}
		c.walletClient = walletClient
	}

	header := http.Header{}
	header.Set(structs.InvokeModeHeader, structs.InvokeModeSync)

	// Query poe object via hash
	indexGetPayload := &wallet.IndexGetPayload{
		Indexs: &wallet.IndexTags{
			IndividualIndex: []string{
				hash,
			},
		},
	}

	IDs, err := c.walletClient.IndexGet(header, indexGetPayload)
	if err != nil {
		logger.Error("Index get fail: %v", err)
		result["ErrMsg"] = err.Error()
		//Goto error page
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	logger.Debug("Index get succ. IDs: %+v", IDs)

	// file hash already exist, query the existed POE infomation
	if len(IDs) > 0 {
		poePayload, err := c.walletClient.QueryPOE(header, did.Identifier(IDs[0]))
		if err != nil {
			logger.Error("Query POE fail: %v", err)
			result["ErrMsg"] = err.Error()
			//Goto error page
			c.Data["json"] = result
			c.ServeJSON()
			return
		}

		// Goto main page
		result["notice"] = "该文件已经存证，详细信息如下"
		result["poe_id"] = poePayload.Id
		result["file_name"] = poePayload.OffchainMetadata.Filename
		result["file_hash"] = poePayload.OffchainMetadata.ContentHash
		result["file_size"] = poePayload.OffchainMetadata.Size
		result["read_only"] = poePayload.OffchainMetadata.ReadOnly
		result["created"] = poePayload.Created
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	// Save poe file
	f, h, err := c.GetFile("poe_file")
	if err != nil {
		logger.Error("Get POE file fail: %v", err)
		result["ErrMsg"] = err.Error()
		//Goto error page
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	defer f.Close()
	// 保存位置在 static/upload, 没有文件夹要先创建
	poeFile := "static/upload/" + h.Filename
	err = c.SaveToFile("poe_file", poeFile)
	if err != nil {
		logger.Error("Save POE file fail: %v", err)
		result["ErrMsg"] = err.Error()
		//Goto error page
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	// New wallet account
	if c.walletID == "" {
		username := "username-" + RandStringRunes(5)
		registerBody := &wallet.RegisterWalletBody{
			Type:   pw.DidType_ORGANIZATION,
			Access: username,
			Secret: "Alice#123456",
		}
		resp, err := c.walletClient.Register(header, registerBody)
		if err != nil {
			logger.Error("Register wallet fail: %v", err)
			result["ErrMsg"] = err.Error()
			//Goto error page
			c.Data["json"] = result
			c.ServeJSON()
			return
		}
		logger.Debug("Register issuer wallet succ.\nResponse: %+v", resp)
		logger.Debug("Wallet public key: %s\n", resp.KeyPair.PublicKey)
		logger.Debug("Wallet private key: %s\n", resp.KeyPair.PrivateKey)

		c.walletID = resp.Id
		c.keyPair = resp.KeyPair
	}

	// Create POE
	poeBody := &wallet.POEBody{
		Name:  "TestPOE",
		Owner: c.walletID,
	}
	signParam := &pki.SignatureParam{
		Creator:    c.walletID,
		Nonce:      "nonce",
		PrivateKey: c.keyPair.PrivateKey,
	}
	poeResp, err := c.walletClient.CreatePOE(header, poeBody, signParam)
	if err != nil {
		logger.Error("Create POE fail: %v", err)
		result["ErrMsg"] = err.Error()
		//Goto error page
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	poeID := poeResp.Id

	uploadResp, err := c.walletClient.UploadPOEFile(header, string(poeID), poeFile, true)
	if err != nil {
		logger.Error("Upload POE file fail: %v", err)
		result["ErrMsg"] = err.Error()
		//Goto error page
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	logger.Debug("Upload POE file succ.\nResponse: %+v\nOffchainMetadata: %+v", uploadResp, *uploadResp.OffchainMetadata)

	// Goto main page
	result["notice"] = "文件存证成功，详细信息如下"
	result["poe_id"] = uploadResp.Id
	result["file_name"] = uploadResp.OffchainMetadata.Filename
	result["file_hash"] = uploadResp.OffchainMetadata.ContentHash
	result["file_size"] = uploadResp.OffchainMetadata.Size
	result["read_only"] = uploadResp.OffchainMetadata.ReadOnly
	result["created"] = time.Now().Unix()
	if len(uploadResp.TransactionIds) > 0 {
		result["transaction_id"] = uploadResp.TransactionIds[0]
	}
	c.Data["json"] = result
	c.ServeJSON()
	return
}
