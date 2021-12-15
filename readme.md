# Introduction
The v5sdk of OKEX go version is only for learning and communication.
(The documentation continues to be improved)
# project instruction

## REST call
``` go
    // Set your APIKey
	apikey := APIKeyInfo{
		ApiKey:     "xxxx",
		SecKey:     "xxxx",
		PassPhrase: "xxxx",
	}

	// The third parameter represents whether it is a simulated environment, see the interface description for more information
	cli := NewRESTClient("https://www.okex.win", &apikey, true)
	rsp, err := cli.Get(context.Background(), "/api/v5/account/balance", nil)
	if err != nil {
		return
	}

	fmt.Println("Response:")
	fmt.Println("\thttp code: ", rsp.Code)
	fmt.Println("\t total time: ", rsp.TotalUsedTime)
	fmt.Println("\t request time consuming: ", rsp.ReqUsedTime)
	fmt.Println("\t return message: ", rsp.Body)
	fmt.Println("\terrCode: ", rsp.V5Response.Code)
	fmt.Println("\terrMsg: ", rsp.V5Response.Msg)
	fmt.Println("\tdata: ", rsp.V5Response.Data)
 ```
See more examples rest/rest_test.go  

## websocket subscription

### Private channel
```go
    ep := "wss://ws.okex.com:8443/ws/v5/private?brokerId=9999"

	// Fill in your own APIKey information
	apikey := "xxxx"
	secretKey := "xxxxx"
	passphrase := "xxxxx"

	// Create ws client
	r, err := NewWsClient(ep)
	if err != nil {
		log.Println(err)
		return
	}

	// Set connection timeout
	r.SetDailTimeout(time.Second * 2)
	err = r.Start()
	if err != nil {
		log.Println(err)
		return
	}
	defer r.Stop()

	var res bool
	// Private channel requires login
	res, _, err = r.Login(apikey, secretKey, passphrase)
	if res {
		fmt.Println("login successful！")
	} else {
		fmt.Println("Login failed！", err)
		return
	}

	
	var args []map[string]string
	arg := make(map[string]string)
	arg["ccy"] = "BTC"
	args = append(args, arg)

	start := time.Now()
	// Subscribe to account channel
	res, _, err = r.PrivAccout(OP_SUBSCRIBE, args)
	if res {
		usedTime := time.Since(start)
		fmt.Println("Successfully subscribed! time consuming:", usedTime.String())
	} else {
		fmt.Println("Subscription failed！", err)
	}

	time.Sleep(100 * time.Second)
	start = time.Now()
	// Unsubscribe account channel
	res, _, err = r.PrivAccout(OP_UNSUBSCRIBE, args)
	if res {
		usedTime := time.Since(start)
		fmt.Println("Unsubscribe successfully！", usedTime.String())
	} else {
		fmt.Println("Failed to unsubscribe！", err)
	}
```
See more examples ws/ws_priv_channel_test.go  

### Public channel
```go
    ep := "wss://ws.okex.com:8443/ws/v5/public?brokerId=9999"

	// Create ws client
	r, err := NewWsClient(ep)
	if err != nil {
		log.Println(err)
		return
	}

	
	// Set connection timeout
	r.SetDailTimeout(time.Second * 2)
	err = r.Start()
	if err != nil {
		log.Println(err)
		return
	}

	defer r.Stop()

	
	var args []map[string]string
	arg := make(map[string]string)
	arg["instType"] = FUTURES
	//arg["instType"] = OPTION
	args = append(args, arg)

	start := time.Now()

	// Subscribe to product channel
	res, _, err := r.PubInstruemnts(OP_SUBSCRIBE, args)
	if res {
		usedTime := time.Since(start)
		fmt.Println("Successfully subscribed！", usedTime.String())
	} else {
		fmt.Println("Subscription failed！", err)
	}

	time.Sleep(30 * time.Second)

	start = time.Now()

	// Unsubscribe product channel
	res, _, err = r.PubInstruemnts(OP_UNSUBSCRIBE, args)
	if res {
		usedTime := time.Since(start)
		fmt.Println("Unsubscribe successfully！", usedTime.String())
	} else {
		fmt.Println("Failed to unsubscribe！", err)
	}
```
See more examples ws/ws_pub_channel_test.go  

## websocket trade
```go
    ep := "wss://ws.okex.com:8443/ws/v5/private?brokerId=9999"

	// Fill in your own APIKey information
	apikey := "xxxx"
	secretKey := "xxxxx"
	passphrase := "xxxxx"

	var res bool
	var req_id string

	// Create ws client
	r, err := NewWsClient(ep)
	if err != nil {
		log.Println(err)
		return
	}

	// Set connection timeout
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

	// Single order
	res, _, err = r.PlaceOrder(req_id, param)
	if res {
		usedTime := time.Since(start)
		fmt.Println("successfully ordered！", usedTime.String())
	} else {
		usedTime := time.Since(start)
		fmt.Println("Order failed！", usedTime.String(), err)
	}

```
See more examples ws/ws_jrpc_test.go  

## wesocket push
Websocket push data is divided into two types of data: `normal push data` and `deep type data`.

```go
ws/wImpl/BookData.go

// Ordinary push
type MsgData struct {
	Arg  map[string]string `json:"arg"`
	Data []interface{}     `json:"data"`
}

// Depth data
type DepthData struct {
	Arg    map[string]string `json:"arg"`
	Action string            `json:"action"`
	Data   []DepthDetail     `json:"data"`
}
```
If you need to process the push data, users can customize the callback function:
1. Callback function for global message processing
   The callback function will process all the data received from the server.
```go
/*
	Add callback function for global message processing
*/
func (a *WsClient) AddMessageHook(fn ReceivedDataCallback) error {
	a.onMessageHook = fn
	return nil
}
```
For usage, please refer to the test case TestAddMessageHook in ws/ws_test.go.

2. Subscription message processing callback function
   It can handle all non-deep types of data, including subscription/unsubscription, and general push data.
```go
/*
Add callback function for subscription message processing
*/
 */
func (a *WsClient) AddBookMsgHook(fn ReceivedMsgDataCallback) error {
	a.onBookMsgHook = fn
	return nil
}
```
For usage, please refer to the test case TestAddBookedDataHook in ws/ws_test.go.


3. In-depth message processing callback function
   What needs to be explained here is that Wsclient provides in-depth data management and automatic checksum functions. If users need to turn off this function, they only need to call the EnableAutoDepthMgr method.
```go
/*
Add callback function for in-depth message processing
*/
 */
func (a *WsClient) AddDepthHook(fn ReceivedDepthDataCallback) error {
	a.onDepthHook = fn
	return nil
}
```
For how to use it, please refer to the test case TestOrderBooks in ws/ws_pub_channel_test.go.

4. Error message type callback function
```go
func (a *WsClient) AddErrMsgHook(fn ReceivedDataCallback) error {
	a.OnErrorHook = fn
	return nil
}
```

# Contact information
Email: caron_co@163.com
WeChat: caron_co