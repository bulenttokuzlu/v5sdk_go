package main

import (
	"context"
	"fmt"
	. "github.com/bulenttokuzlu/v5sdk_go/rest"
	. "github.com/bulenttokuzlu/v5sdk_go/ws"
	"log"
	"time"
)

/*
	rest API请求
	更多示例请查看 rest/rest_test.go
*/
func REST() {
	// 设置您的APIKey
	apikey := APIKeyInfo{
		ApiKey:     "xxxx",
		SecKey:     "xxxx",
		PassPhrase: "xxxx",
	}

	// 第三个参数代表是否为模拟环境，更多信息查看接口说明
	cli := NewRESTClient("https://www.okex.win", &apikey, true)
	rsp, err := cli.Get(context.Background(), "/api/v5/account/balance", nil)
	if err != nil {
		return
	}

	fmt.Println("Response:")
	fmt.Println("\thttp code: ", rsp.Code)
	fmt.Println("\tTotal time consuming: ", rsp.TotalUsedTime)
	fmt.Println("\tRequest time consuming: ", rsp.ReqUsedTime)
	fmt.Println("\tReturn message: ", rsp.Body)
	fmt.Println("\terrCode: ", rsp.V5Response.Code)
	fmt.Println("\terrMsg: ", rsp.V5Response.Msg)
	fmt.Println("\tdata: ", rsp.V5Response.Data)

}

// 订阅私有频道
func wsPriv() {
	ep := "wss://ws.okex.com:8443/ws/v5/private?brokerId=9999"

	// 填写您自己的APIKey信息
	apikey := "xxxx"
	secretKey := "xxxxx"
	passphrase := "xxxxx"

	// 创建ws客户端
	r, err := NewWsClient(ep)
	if err != nil {
		log.Println(err)
		return
	}

	// 设置连接超时
	r.SetDailTimeout(time.Second * 2)
	err = r.Start()
	if err != nil {
		log.Println(err)
		return
	}
	defer r.Stop()
	var res bool

	res, _, err = r.Login(apikey, secretKey, passphrase)
	if res {
		fmt.Println("login successful！")
	} else {
		fmt.Println("Login failed！", err)
		return
	}

	// 订阅账户频道
	var args []map[string]string
	arg := make(map[string]string)
	arg["ccy"] = "BTC"
	args = append(args, arg)

	start := time.Now()
	res, _, err = r.PrivAccout(OP_SUBSCRIBE, args)
	if res {
		usedTime := time.Since(start)
		fmt.Println("Successfully subscribed! time consuming:", usedTime.String())
	} else {
		fmt.Println("Subscription failed！", err)
	}

	time.Sleep(100 * time.Second)
	start = time.Now()
	res, _, err = r.PrivAccout(OP_UNSUBSCRIBE, args)
	if res {
		usedTime := time.Since(start)
		fmt.Println("Unsubscribe successfully！", usedTime.String())
	} else {
		fmt.Println("Failed to unsubscribe！", err)
	}

}

// 订阅公共频道
func wsPub() {
	ep := "wss://ws.okex.com:8443/ws/v5/public?brokerId=9999"

	// 创建ws客户端
	r, err := NewWsClient(ep)
	if err != nil {
		log.Println(err)
		return
	}

	// 设置连接超时
	r.SetDailTimeout(time.Second * 2)
	err = r.Start()
	if err != nil {
		log.Println(err)
		return
	}
	defer r.Stop()
	// 订阅产品频道
	var args []map[string]string
	arg := make(map[string]string)
	arg["instType"] = FUTURES
	//arg["instType"] = OPTION
	args = append(args, arg)

	start := time.Now()
	res, _, err := r.PubInstruemnts(OP_SUBSCRIBE, args)
	if res {
		usedTime := time.Since(start)
		fmt.Println("Successfully subscribed！", usedTime.String())
	} else {
		fmt.Println("Subscription failed！", err)
	}

	time.Sleep(30 * time.Second)

	start = time.Now()
	res, _, err = r.PubInstruemnts(OP_UNSUBSCRIBE, args)
	if res {
		usedTime := time.Since(start)
		fmt.Println("Unsubscribe successfully！", usedTime.String())
	} else {
		fmt.Println("Failed to unsubscribe！", err)
	}
}

// websocket交易
func wsJrpc() {
	ep := "wss://ws.okex.com:8443/ws/v5/private?brokerId=9999"

	// 填写您自己的APIKey信息
	apikey := "xxxx"
	secretKey := "xxxxx"
	passphrase := "xxxxx"

	var res bool
	var req_id string

	// 创建ws客户端
	r, err := NewWsClient(ep)
	if err != nil {
		log.Println(err)
		return
	}

	// 设置连接超时
	r.SetDailTimeout(time.Second * 2)
	err = r.Start()
	if err != nil {
		log.Println(err)
		return
	}

	defer r.Stop()

	res, _, err = r.Login(apikey, secretKey, passphrase)
	if res {
		fmt.Println("login successful！")
	} else {
		fmt.Println("Login failed！", err)
		return
	}

	start := time.Now()
	param := map[string]interface{}{}
	param["instId"] = "BTC-USDT"
	param["tdMode"] = "cash"
	param["side"] = "buy"
	param["ordType"] = "market"
	param["sz"] = "200"
	req_id = "00001"

	res, _, err = r.PlaceOrder(req_id, param)
	if res {
		usedTime := time.Since(start)
		fmt.Println("successfully ordered！", usedTime.String())
	} else {
		usedTime := time.Since(start)
		fmt.Println("下单失败！", usedTime.String(), err)
	}
}

func main() {
	// 公共订阅
	wsPub()

	// 私有订阅
	wsPriv()

	// websocket交易
	wsJrpc()

	// rest请求
	REST()
}
