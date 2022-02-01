package ws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/bulenttokuzlu/v5sdk_go/config"
	. "github.com/bulenttokuzlu/v5sdk_go/utils"
	. "github.com/bulenttokuzlu/v5sdk_go/ws/wImpl"
	"log"
	"regexp"
	"runtime/debug"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Global callback function
type ReceivedDataCallback func(*Msg) error

// Ordinary subscription push data callback function
type ReceivedMsgDataCallback func(time.Time, MsgData) error

// Deep subscription push data callback function
type ReceivedDepthDataCallback func(time.Time, DepthData) error

// websocket
type WsClient struct {
	WsEndPoint string
	WsApi      *ApiInfo
	conn       *websocket.Conn
	sendCh     chan string //Message queue
	resCh      chan *Msg   //Receive message queue

	errCh chan *Msg
	regCh map[Event]chan *Msg //Request response queue

	quitCh chan struct{}
	lock   sync.RWMutex

	onMessageHook ReceivedDataCallback      //Global message callback function
	onBookMsgHook ReceivedMsgDataCallback   //Ordinary subscription message callback function
	onDepthHook   ReceivedDepthDataCallback //Deep subscription message callback function
	OnErrorHook   ReceivedDataCallback      //Error handling callback function

	// Record depth information
	DepthDataList map[string]DepthDetail
	autoDepthMgr  bool // In-depth data management (checksum, etc.)
	DepthDataLock sync.RWMutex

	isStarted   bool //Prevent repeated startup and shutdown
	dailTimeout time.Duration
}

/*
Server response details
Timestamp: The time when the message was received
Info: The received message string
*/
type Msg struct {
	Timestamp time.Time   `json:"timestamp"`
	Info      interface{} `json:"info"`
}

func (this *Msg) Print() {
	fmt.Println("【Message time】", this.Timestamp.Format("2006-01-02 15:04:05.000"))
	str, _ := json.Marshal(this.Info)
	fmt.Println("【Message content】", string(str))
}

/*
	Message structure after subscription result encapsulation
*/
type ProcessDetail struct {
	EndPoint string        `json:"endPoint"`
	ReqInfo  string        `json:"ReqInfo"`  //订阅请求
	SendTime time.Time     `json:"sendTime"` //发送订阅请求的时间
	RecvTime time.Time     `json:"recvTime"` //接受到订阅结果的时间
	UsedTime time.Duration `json:"UsedTime"` //耗时
	Data     []*Msg        `json:"data"`     //订阅结果数据
}

func (p *ProcessDetail) String() string {
	data, _ := json.Marshal(p)
	return string(data)
}

// 创建ws对象
func NewWsClient(ep string) (r *WsClient, err error) {
	if ep == "" {
		err = errors.New("websocket endpoint cannot be null")
		return
	}

	r = &WsClient{
		WsEndPoint: ep,
		sendCh:     make(chan string),
		resCh:      make(chan *Msg),
		errCh:      make(chan *Msg),
		regCh:      make(map[Event]chan *Msg),
		//cbs:        make(map[Event]ReceivedDataCallback),
		quitCh:        make(chan struct{}),
		DepthDataList: make(map[string]DepthDetail),
		dailTimeout:   time.Second * 5,
		// 自动深度校验默认开启
		autoDepthMgr: true,
	}

	return
}

/*
	新增记录深度信息
*/
func (a *WsClient) addDepthDataList(key string, dd DepthDetail) error {
	a.DepthDataLock.Lock()
	defer a.DepthDataLock.Unlock()
	a.DepthDataList[key] = dd
	return nil
}

/*
	更新记录深度信息（如果没有记录不会更新成功）
*/
func (a *WsClient) updateDepthDataList(key string, dd DepthDetail) error {
	a.DepthDataLock.Lock()
	defer a.DepthDataLock.Unlock()
	if _, ok := a.DepthDataList[key]; !ok {
		return errors.New("Update failed! No record found" + key)
	}

	a.DepthDataList[key] = dd
	return nil
}

/*
	删除记录深度信息
*/
func (a *WsClient) deleteDepthDataList(key string) error {
	a.DepthDataLock.Lock()
	defer a.DepthDataLock.Unlock()
	delete(a.DepthDataList, key)
	return nil
}

/*
	设置是否自动深度管理，开启 true，关闭 false
*/
func (a *WsClient) EnableAutoDepthMgr(b bool) error {
	a.DepthDataLock.Lock()
	defer a.DepthDataLock.Unlock()

	if len(a.DepthDataList) != 0 {
		err := errors.New("In-depth data is currently subscribed")
		return err
	}

	a.autoDepthMgr = b
	return nil
}

/*
	获取当前的深度快照信息(合并后的)
*/
func (a *WsClient) GetSnapshotByChannel(data DepthData) (snapshot *DepthDetail, err error) {
	key, err := json.Marshal(data.Arg)
	if err != nil {
		return
	}
	a.DepthDataLock.Lock()
	defer a.DepthDataLock.Unlock()
	val, ok := a.DepthDataList[string(key)]
	if !ok {
		return
	}
	snapshot = new(DepthDetail)
	raw, err := json.Marshal(val)
	if err != nil {
		return
	}
	err = json.Unmarshal(raw, &snapshot)
	if err != nil {
		return
	}
	return
}

// 设置dial超时时间
func (a *WsClient) SetDailTimeout(tm time.Duration) {
	a.dailTimeout = tm
}

// Non-blocking start
func (a *WsClient) Start(keepalive bool) error {
	a.lock.RLock()
	if a.isStarted {
		a.lock.RUnlock()
		fmt.Println("ws has started")
		return nil
	} else {
		a.lock.RUnlock()
		a.lock.Lock()
		defer a.lock.Unlock()
		// Increase timeout handling
		done := make(chan struct{})
		ctx, cancel := context.WithTimeout(context.Background(), a.dailTimeout)
		defer cancel()
		go func(ctx context.Context) {
			defer func() {
				close(done)
			}()
			var c *websocket.Conn
			c, _, err := websocket.DefaultDialer.Dial(a.WsEndPoint, nil)
			if err != nil {
				err = errors.New("dial error:" + err.Error())
				return
			}
			a.conn = c

		}(ctx)
		select {
		case <-ctx.Done():
			err := errors.New("Connection timeout exit！")
			return err
		case <-done:

		}

		go a.receive()
		if keepalive {
			go a.work()
		}
		a.isStarted = true
		log.Println("Client started!", a.WsEndPoint)
		return nil
	}
}

// Client exit message channel
func (a *WsClient) IsQuit() <-chan struct{} {
	return a.quitCh
}

func (a *WsClient) work() {
	defer func() {
		a.Stop()
		err := recover()
		if err != nil {
			log.Printf("work End. Recover msg: %+v", a)
			debug.PrintStack()
		}

	}()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C: // Keep heartbeat
			// go a.Ping(1000)
			go func() {
				_, _, err := a.Ping(1000)
				if err != nil {
					fmt.Println("Heartbeat detection failed！", err)
					a.Stop()
					return
				}

			}()

		case <-a.quitCh: // Keep heartbeat
			return
		case data, ok := <-a.resCh: //Receive a message from the server
			if !ok {
				return
			}
			//log.Println("Received a message from resCh:", data)
			if a.onMessageHook != nil {
				err := a.onMessageHook(data)
				if err != nil {
					log.Println("Error executing onMessageHook function！", err)
				}
			}
		case errMsg, ok := <-a.errCh: //Error handling
			if !ok {
				return
			}
			if a.OnErrorHook != nil {
				err := a.OnErrorHook(errMsg)
				if err != nil {
					log.Println("Execute OnErrorHook function error！", err)
				}
			}
		case req, ok := <-a.sendCh: //Remove the data from the sending queue and send it to the server
			if !ok {
				return
			}
			//log.Println("Received a message from req:", req)
			err := a.conn.WriteMessage(websocket.TextMessage, []byte(req))
			if err != nil {
				log.Printf("Failed to send request: %s\n", err)
				return
			}
			log.Printf("[send request] %v\n", req)
		}
	}

}

/*
	Process the received message
*/
func (a *WsClient) receive() {
	defer func() {
		a.Stop()
		err := recover()
		if err != nil {
			log.Printf("Receive End. Recover msg: %+v", a)
			debug.PrintStack()
		}

	}()

	for {
		messageType, message, err := a.conn.ReadMessage()
		if err != nil {
			if a.isStarted {
				log.Println("receive message error!" + err.Error())
			}

			break
		}

		txtMsg := message
		switch messageType {
		case websocket.TextMessage:
		case websocket.BinaryMessage:
			txtMsg, err = GzipDecode(message)
			if err != nil {
				log.Println("Unzip failed！")
				continue
			}
		}

		log.Println("[Received the news]", string(txtMsg))

		//Send the result to the default message processing channel

		timestamp := time.Now()
		msg := &Msg{Timestamp: timestamp, Info: string(txtMsg)}

		a.resCh <- msg

		evt, data, err := a.parseMessage(txtMsg)
		if err != nil {
			log.Println("Failed to parse the message！", err)
			continue
		}

		//log.Println("Message parsing succeeded! Message type =", evt)

		a.lock.RLock()
		ch, ok := a.regCh[evt]
		a.lock.RUnlock()
		if !ok {
			//Only push messages will actively create channels and consumption queues
			if evt == EVENT_BOOKED_DATA || evt == EVENT_DEPTH_DATA {
				//log.Println("channel does not exist！event:", evt)
				//a.lock.RUnlock()
				a.lock.Lock()
				a.regCh[evt] = make(chan *Msg)
				ch = a.regCh[evt]
				a.lock.Unlock()

				//log.Println("create", evt, "aisle")

				// Create consumption queue
				go func(evt Event) {
					//log.Println("Create goroutine  evt:", evt)

					for msg := range a.regCh[evt] {
						//log.Println(msg)
						// msg.Print()
						switch evt {
						// Processing normal push data
						case EVENT_BOOKED_DATA:
							fn := a.onBookMsgHook
							if fn != nil {
								err = fn(msg.Timestamp, msg.Info.(MsgData))
								if err != nil {
									log.Println("The subscribe data callback function failed to execute！", err)
								}
								//log.Println("Function executed successfully！", err)
							}
						// Processing deep push data
						case EVENT_DEPTH_DATA:
							fn := a.onDepthHook

							depData := msg.Info.(DepthData)

							// If the deep data management function is turned on, deep data will be merged
							if a.autoDepthMgr {
								a.MergeDepth(depData)
							}

							// Run user-defined callback function
							if fn != nil {
								err = fn(msg.Timestamp, msg.Info.(DepthData))
								if err != nil {
									log.Println("Deep callback function execution failed！", err)
								}

							}
						}

					}
					//log.Println("Exit goroutine  evt:", evt)
				}(evt)

				//continue
			} else {
				//log.Println("The program is abnormal! Channel closed", evt)
				continue
			}

		}

		//log.Println(evt,"Event registered",ch)

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Millisecond*1000)
		defer cancel()
		select {
		/*
			Discarding messages can easily lead to data processing errors
		*/
		// case <-ctx.Done():
		// 	log.Println("Waiting timeout, message discarded - ", data)
		case ch <- &Msg{Timestamp: timestamp, Info: data}:
		}
	}
}

/*
	After the in-depth data management function is turned on, the system will automatically merge in-depth information
*/
func (a *WsClient) MergeDepth(depData DepthData) (err error) {
	if !a.autoDepthMgr {
		return
	}

	key, err := json.Marshal(depData.Arg)
	if err != nil {
		err = errors.New("data error")
		return
	}

	// books5 No need to do checksum
	if depData.Arg["channel"] == "books5" {
		a.addDepthDataList(string(key), depData.Data[0])
		return
	}

	if depData.Action == "snapshot" {

		_, err = depData.CheckSum(nil)
		if err != nil {
			log.Println("Verification failed", err)
			return
		}

		a.addDepthDataList(string(key), depData.Data[0])

	} else {

		var newSnapshot *DepthDetail
		a.DepthDataLock.RLock()
		oldSnapshot, ok := a.DepthDataList[string(key)]
		if !ok {
			log.Println("Depth data error, not found in all data！")
			err = errors.New("data error")
			return
		}
		a.DepthDataLock.RUnlock()
		newSnapshot, err = depData.CheckSum(&oldSnapshot)
		if err != nil {
			log.Println("Depth verification failed", err)
			err = errors.New("Verification failed")
			return
		}

		a.updateDepthDataList(string(key), *newSnapshot)
	}
	return
}

/*
	Judge the event type by ErrorCode
*/
func GetInfoFromErrCode(data ErrData) Event {
	switch data.Code {
	case "60001":
		return EVENT_LOGIN
	case "60002":
		return EVENT_LOGIN
	case "60003":
		return EVENT_LOGIN
	case "60004":
		return EVENT_LOGIN
	case "60005":
		return EVENT_LOGIN
	case "60006":
		return EVENT_LOGIN
	case "60007":
		return EVENT_LOGIN
	case "60008":
		return EVENT_LOGIN
	case "60009":
		return EVENT_LOGIN
	case "60010":
		return EVENT_LOGIN
	case "60011":
		return EVENT_LOGIN
	}

	return EVENT_UNKNOWN
}

/*
Resolve the corresponding channel from the error return
    Sample error message
 {"event":"error","msg":"channel:index-tickers,instId:BTC-USDT1 doesn't exist","code":"60018"}
*/
func GetInfoFromErrMsg(raw string) (channel string) {
	reg := regexp.MustCompile(`channel:(.*?),`)
	if reg == nil {
		fmt.Println("MustCompile err")
		return
	}
	//Extract key information
	result := reg.FindAllStringSubmatch(raw, -1)
	for _, text := range result {
		channel = text[1]
	}
	return
}

/*
	Parse the message type
*/
func (a *WsClient) parseMessage(raw []byte) (evt Event, data interface{}, err error) {
	evt = EVENT_UNKNOWN
	//log.Println("Parse the message")
	//log.Println("Message content:", string(raw))
	if string(raw) == "pong" {
		evt = EVENT_PING
		data = raw
		return
	}
	//log.Println(0, evt)
	var rspData = RspData{}
	err = json.Unmarshal(raw, &rspData)
	if err == nil {
		op := rspData.Event
		if op == OP_SUBSCRIBE || op == OP_UNSUBSCRIBE {
			channel := rspData.Arg["channel"]
			evt = GetEventId(channel)
			data = rspData
			return
		}
	}

	//log.Println("ErrData")
	var errData = ErrData{}
	err = json.Unmarshal(raw, &errData)
	if err == nil {
		op := errData.Event
		switch op {
		case OP_LOGIN:
			evt = EVENT_LOGIN
			data = errData
			//log.Println(3, evt)
			return
		case OP_ERROR:
			data = errData
			// TODO:Refine the event judgment corresponding to the error report

			//Try to parse the corresponding event type from the msg field
			evt = GetInfoFromErrCode(errData)
			if evt != EVENT_UNKNOWN {
				return
			}
			evt = GetEventId(GetInfoFromErrMsg(errData.Msg))
			if evt == EVENT_UNKNOWN {
				evt = EVENT_ERROR
				return
			}
			return
		}
		//log.Println(5, evt)
	}

	//log.Println("JRPCRsp")
	var jRPCRsp = JRPCRsp{}
	err = json.Unmarshal(raw, &jRPCRsp)
	if err == nil {
		data = jRPCRsp
		evt = GetEventId(jRPCRsp.Op)
		if evt != EVENT_UNKNOWN {
			return
		}
	}

	var depthData = DepthData{}
	err = json.Unmarshal(raw, &depthData)
	if err == nil {
		evt = EVENT_DEPTH_DATA
		data = depthData
		//log.Println("-->>EVENT_DEPTH_DATA", evt)
		//log.Println(evt, data)
		//log.Println(6)
		switch depthData.Arg["channel"] {
		case "books":
			return
		case "books-l2-tbt":
			return
		case "books50-l2-tbt":
			return
		case "books5":
			return
		default:

		}
	}

	//log.Println("MsgData")
	var msgData = MsgData{}
	err = json.Unmarshal(raw, &msgData)
	if err == nil {
		evt = EVENT_BOOKED_DATA
		data = msgData
		//log.Println("-->>EVENT_BOOK_DATA", evt)
		//log.Println(evt, data)
		//log.Println(6)
		return
	}

	evt = EVENT_UNKNOWN
	err = errors.New("message unknown")
	return
}

func (a *WsClient) Stop() error {

	a.lock.Lock()
	defer a.lock.Unlock()
	if !a.isStarted {
		return nil
	}

	a.isStarted = false

	if a.conn != nil {
		a.conn.Close()
	}
	close(a.errCh)
	close(a.sendCh)
	close(a.resCh)
	close(a.quitCh)

	for _, ch := range a.regCh {
		close(ch)
	}

	log.Println("ws client exit!")
	return nil
}

/*
	Add callback function for global message processing
*/
func (a *WsClient) AddMessageHook(fn ReceivedDataCallback) error {
	a.onMessageHook = fn
	return nil
}

/*
	Add callback function for subscription message processing
*/
func (a *WsClient) AddBookMsgHook(fn ReceivedMsgDataCallback) error {
	a.onBookMsgHook = fn
	return nil
}

/*
	Add callback function for in-depth message processing
	For example:
	cli.AddDepthHook(func(ts time.Time, data DepthData) error { return nil })
*/
func (a *WsClient) AddDepthHook(fn ReceivedDepthDataCallback) error {
	a.onDepthHook = fn
	return nil
}

/*
	Add callback function for error message processing
*/
func (a *WsClient) AddErrMsgHook(fn ReceivedDataCallback) error {
	a.OnErrorHook = fn
	return nil
}

/*
	Determine if the connection is alive
*/
func (a *WsClient) IsAlive() bool {
	res := false
	if a.conn == nil {
		return res
	}
	res, _, _ = a.Ping(500)
	return res
}
