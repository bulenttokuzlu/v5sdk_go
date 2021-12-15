package ws

import (
	"context"
	"errors"
	. "github.com/bulenttokuzlu/v5sdk_go/config"
	"github.com/bulenttokuzlu/v5sdk_go/rest"
	. "github.com/bulenttokuzlu/v5sdk_go/utils"
	. "github.com/bulenttokuzlu/v5sdk_go/ws/wImpl"
	. "github.com/bulenttokuzlu/v5sdk_go/ws/wInterface"
	"log"
	"sync"
	"time"
)

/*
	Ping服务端保持心跳。
	timeOut:超时时间(毫秒)，如果不填默认为5000ms
*/
func (a *WsClient) Ping(timeOut ...int) (res bool, detail *ProcessDetail, err error) {
	tm := 5000
	if len(timeOut) != 0 {
		tm = timeOut[0]
	}
	res = true

	detail = &ProcessDetail{
		EndPoint: a.WsEndPoint,
	}

	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, time.Duration(tm)*time.Millisecond)
	ctx = context.WithValue(ctx, "detail", detail)
	msg, err := a.process(ctx, EVENT_PING, nil)
	if err != nil {
		res = false
		log.Println("Failed to process request!", err)
		return
	}
	detail.Data = msg

	if len(msg) == 0 {
		res = false
		return
	}

	str := string(msg[0].Info.([]byte))
	if str != "pong" {
		res = false
		return
	}

	return
}

/*
	Sign in to private channel
*/
func (a *WsClient) Login(apiKey, secKey, passPhrase string, timeOut ...int) (res bool, detail *ProcessDetail, err error) {

	if apiKey == "" {
		err = errors.New("ApiKey cannot be null")
		return
	}

	if secKey == "" {
		err = errors.New("SecretKey cannot be null")
		return
	}

	if passPhrase == "" {
		err = errors.New("Passphrase cannot be null")
		return
	}

	a.WsApi = &ApiInfo{
		ApiKey:     apiKey,
		SecretKey:  secKey,
		Passphrase: passPhrase,
	}

	tm := 5000
	if len(timeOut) != 0 {
		tm = timeOut[0]
	}
	res = true

	timestamp := EpochTime()

	preHash := PreHashString(timestamp, rest.GET, "/users/self/verify", "")
	//fmt.Println("preHash:", preHash)
	var sign string
	if sign, err = HmacSha256Base64Signer(preHash, secKey); err != nil {
		log.Println("Failed to process signature！", err)
		return
	}

	args := map[string]string{}
	args["apiKey"] = apiKey
	args["passphrase"] = passPhrase
	args["timestamp"] = timestamp
	args["sign"] = sign
	req := &ReqData{
		Op:   OP_LOGIN,
		Args: []map[string]string{args},
	}

	detail = &ProcessDetail{
		EndPoint: a.WsEndPoint,
	}

	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, time.Duration(tm)*time.Millisecond)
	ctx = context.WithValue(ctx, "detail", detail)

	msg, err := a.process(ctx, EVENT_LOGIN, req)
	if err != nil {
		res = false
		log.Println("Failed to process request!", req, err)
		return
	}
	detail.Data = msg

	if len(msg) == 0 {
		res = false
		return
	}

	info, _ := msg[0].Info.(ErrData)

	if info.Code == "0" && info.Event == OP_LOGIN {
		log.Println("login successful!")
	} else {
		log.Println("Login failed!")
		res = false
		return
	}

	return
}

/*
	Wait for result response
*/
func (a *WsClient) waitForResult(e Event, timeOut int) (data interface{}, err error) {

	if _, ok := a.regCh[e]; !ok {
		a.lock.Lock()
		a.regCh[e] = make(chan *Msg)
		a.lock.Unlock()
		//log.Println("Registered ", e, "event succeeded")
	}

	a.lock.RLock()
	defer a.lock.RUnlock()
	ch := a.regCh[e]
	//log.Println(e, "Waiting for response！")
	select {
	case <-time.After(time.Duration(timeOut) * time.Millisecond):
		log.Println(e, "Timeout not responding！")
		err = errors.New(e.String() + "Timeout not responding！")
		return
	case data = <-ch:
		//log.Println(data)
	}

	return
}

/*
	Send a message to the server
*/
func (a *WsClient) Send(ctx context.Context, op WSReqData) (err error) {
	select {
	case <-ctx.Done():
		log.Println("Exit on failure！")
		err = errors.New("Send timeout exit！")
	case a.sendCh <- op.ToString():
	}

	return
}

func (a *WsClient) process(ctx context.Context, e Event, op WSReqData) (data []*Msg, err error) {
	defer func() {
		_ = recover()
	}()

	var detail *ProcessDetail
	if val := ctx.Value("detail"); val != nil {
		detail = val.(*ProcessDetail)
	} else {
		detail = &ProcessDetail{
			EndPoint: a.WsEndPoint,
		}
	}
	defer func() {
		//fmt.Println("Processing complete,", e.String())
		detail.UsedTime = detail.RecvTime.Sub(detail.SendTime)
	}()

	//Check if the event is registered
	if _, ok := a.regCh[e]; !ok {
		a.lock.Lock()
		a.regCh[e] = make(chan *Msg)
		a.lock.Unlock()
		//log.Println("Register", e, "Event successful")
	} else {
		//log.Println("Event ", e, "registered! ")
		err = errors.New("event" + e.String() + "Not yet finished")
		return
	}

	//The number of expected request responses
	expectCnt := 1
	if op != nil {
		expectCnt = op.Len()
	}
	recvCnt := 0

	//Wait for completion notification
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(ctx context.Context) {
		defer func() {
			a.lock.Lock()
			delete(a.regCh, e)
			//log.Println("Event has been logged out!",e)
			a.lock.Unlock()
			wg.Done()
		}()

		a.lock.RLock()
		ch := a.regCh[e]
		a.lock.RUnlock()

		//log.Println(e, "Waiting for response！")
		done := false
		ok := true
		for {
			var item *Msg
			select {
			case <-ctx.Done():
				log.Println(e, "Timeout not responding！")
				err = errors.New(e.String() + "Timeout not responding！")
				return
			case item, ok = <-ch:
				if !ok {
					return
				}
				detail.RecvTime = time.Now()
				//log.Println(e, "Data received", item)
				data = append(data, item)
				recvCnt++
				//log.Println(data)
				if recvCnt == expectCnt {
					done = true
					break
				}
			}
			if done {
				break
			}
		}
		if ok {
			close(ch)
		}

	}(ctx)

	switch e {
	case EVENT_PING:
		msg := "ping"
		detail.ReqInfo = msg
		a.sendCh <- msg
		detail.SendTime = time.Now()
	default:
		detail.ReqInfo = op.ToString()
		err = a.Send(ctx, op)
		if err != nil {
			log.Println("send[", e, "]Message failed！", err)
			return
		}
		detail.SendTime = time.Now()
	}

	wg.Wait()
	return
}

/*
Determine the request type according to the args request parameter
For example: {"channel": "account","ccy": "BTC"} The type is EVENT_BOOK_ACCOUNT
*/
func GetEventByParam(param map[string]string) (evtId Event) {
	evtId = EVENT_UNKNOWN
	channel, ok := param["channel"]
	if !ok {
		return
	}

	evtId = GetEventId(channel)
	return
}

/*
Subscribe to the channel.
req: request json string
*/
func (a *WsClient) Subscribe(param map[string]string, timeOut ...int) (res bool, detail *ProcessDetail, err error) {
	res = true
	tm := 5000
	if len(timeOut) != 0 {
		tm = timeOut[0]
	}

	evtid := GetEventByParam(param)
	if evtid == EVENT_UNKNOWN {
		err = errors.New("Illegal request parameter！")
		return
	}

	var args []map[string]string
	args = append(args, param)

	req := ReqData{
		Op:   OP_SUBSCRIBE,
		Args: args,
	}

	detail = &ProcessDetail{
		EndPoint: a.WsEndPoint,
	}

	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, time.Duration(tm)*time.Millisecond)
	ctx = context.WithValue(ctx, "detail", detail)

	msg, err := a.process(ctx, evtid, req)
	if err != nil {
		res = false
		log.Println("Failed to process request!", req, err)
		return
	}
	detail.Data = msg

	//Check whether all channels are updated successfully
	res, err = checkResult(req, msg)
	if err != nil {
		res = false
		return
	}

	return
}

/*
Unsubscribe from the channel.
req: request json string
*/
func (a *WsClient) UnSubscribe(param map[string]string, timeOut ...int) (res bool, detail *ProcessDetail, err error) {
	res = true
	tm := 5000
	if len(timeOut) != 0 {
		tm = timeOut[0]
	}

	evtid := GetEventByParam(param)
	if evtid == EVENT_UNKNOWN {
		err = errors.New("Illegal request parameter！")
		return
	}

	var args []map[string]string
	args = append(args, param)

	req := ReqData{
		Op:   OP_UNSUBSCRIBE,
		Args: args,
	}

	detail = &ProcessDetail{
		EndPoint: a.WsEndPoint,
	}

	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, time.Duration(tm)*time.Millisecond)
	ctx = context.WithValue(ctx, "detail", detail)
	msg, err := a.process(ctx, evtid, req)
	if err != nil {
		res = false
		log.Println("Failed to process request!", req, err)
		return
	}
	detail.Data = msg
	//Check whether all channels are updated successfully
	res, err = checkResult(req, msg)
	if err != nil {
		res = false
		return
	}

	return
}

/*
	jrpc request
*/
func (a *WsClient) Jrpc(id, op string, params []map[string]interface{}, timeOut ...int) (res bool, detail *ProcessDetail, err error) {
	res = true
	tm := 5000
	if len(timeOut) != 0 {
		tm = timeOut[0]
	}

	evtid := GetEventId(op)
	if evtid == EVENT_UNKNOWN {
		err = errors.New("Illegal request parameter！")
		return
	}

	req := JRPCReq{
		Id:   id,
		Op:   op,
		Args: params,
	}
	detail = &ProcessDetail{
		EndPoint: a.WsEndPoint,
	}

	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, time.Duration(tm)*time.Millisecond)
	ctx = context.WithValue(ctx, "detail", detail)
	msg, err := a.process(ctx, evtid, req)
	if err != nil {
		res = false
		log.Println("Failed to process request!", req, err)
		return
	}
	detail.Data = msg

	//Check whether all channels are updated successfully
	res, err = checkResult(req, msg)
	if err != nil {
		res = false
		return
	}

	return
}

func (a *WsClient) PubChannel(evtId Event, op string, params []map[string]string, pd Period, timeOut ...int) (res bool, msg []*Msg, err error) {

	// Parameter verification
	pa, err := checkParams(evtId, params, pd)
	if err != nil {
		return
	}

	res = true
	tm := 5000
	if len(timeOut) != 0 {
		tm = timeOut[0]
	}

	req := ReqData{
		Op:   op,
		Args: pa,
	}

	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, time.Duration(tm)*time.Millisecond)
	msg, err = a.process(ctx, evtId, req)
	if err != nil {
		res = false
		log.Println("Failed to process request!", req, err)
		return
	}

	//Check whether all channels are updated successfully

	res, err = checkResult(req, msg)
	if err != nil {
		res = false
		return
	}

	return
}

// Parameter verification
func checkParams(evtId Event, params []map[string]string, pd Period) (res []map[string]string, err error) {

	channel := evtId.GetChannel(pd)
	if channel == "" {
		err = errors.New("Parameter verification failed! Unknown type:" + evtId.String())
		return
	}
	log.Println(channel)
	if params == nil {
		tmp := make(map[string]string)
		tmp["channel"] = channel
		res = append(res, tmp)
	} else {
		//log.Println(params)
		for _, param := range params {

			tmp := make(map[string]string)
			for k, v := range param {
				tmp[k] = v
			}

			val, ok := tmp["channel"]
			if !ok {
				tmp["channel"] = channel
			} else {
				if val != channel {
					err = errors.New("Parameter verification failed! channel should be" + channel + val)
					return
				}
			}

			res = append(res, tmp)
		}
	}

	return
}
