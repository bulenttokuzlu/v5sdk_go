package ws

import (
	"encoding/json"
	"fmt"
	. "github.com/bulenttokuzlu/v5sdk_go/ws/wImpl"
	"log"
	"strings"
	"testing"
	"time"
)

func prework() *WsClient {
	ep := "wss://ws.okex.com:8443/ws/v5/public?brokerId=9999"
	r, err := NewWsClient(ep)
	if err != nil {
		log.Fatal(err)
	}

	err = r.Start()
	if err != nil {
		log.Fatal(err, ep)
	}
	return r
}

// 产品频道测试
func TestInstruemnts(t *testing.T) {
	r := prework()
	var err error
	var res bool

	var args []map[string]string
	arg := make(map[string]string)
	arg["instType"] = FUTURES
	//arg["instType"] = OPTION
	args = append(args, arg)

	start := time.Now()
	res, _, err = r.PubInstruemnts(OP_SUBSCRIBE, args)
	if res {
		usedTime := time.Since(start)
		fmt.Println("Successfully subscribed！", usedTime.String())
	} else {
		fmt.Println("Subscription failed！", err)
		t.Fatal("Subscription failed！", err)
		//return
	}

	time.Sleep(3 * time.Second)
	//等待推送

	start = time.Now()
	res, _, err = r.PubInstruemnts(OP_UNSUBSCRIBE, args)
	if res {
		usedTime := time.Since(start)
		fmt.Println("Unsubscribe successfully！", usedTime.String())
	} else {
		fmt.Println("Failed to unsubscribe！", err)
		t.Fatal("Failed to unsubscribe！", err)
	}

}

// status频道测试
func TestStatus(t *testing.T) {
	r := prework()
	var err error
	var res bool

	start := time.Now()
	res, _, err = r.PubStatus(OP_SUBSCRIBE)
	if res {
		usedTime := time.Since(start)
		fmt.Println("Successfully subscribed！", usedTime.String())
	} else {
		fmt.Println("Subscription failed！", err)
		t.Fatal("Subscription failed！", err)
		//return
	}

	time.Sleep(10000 * time.Second)
	//等待推送

	start = time.Now()
	res, _, err = r.PubStatus(OP_UNSUBSCRIBE)
	if res {
		usedTime := time.Since(start)
		fmt.Println("Unsubscribe successfully！", usedTime.String())
	} else {
		fmt.Println("Cancel Subscription failed！", err)
		t.Fatal("Cancel Subscription failed！", err)
	}

}

// 行情频道测试
func TestTickers(t *testing.T) {
	r := prework()
	var err error
	var res bool

	var args []map[string]string
	arg := make(map[string]string)
	arg["instId"] = "BTC-USDT"

	args = append(args, arg)

	start := time.Now()
	res, _, err = r.PubTickers(OP_SUBSCRIBE, args)
	if res {
		usedTime := time.Since(start)
		fmt.Println("Successfully subscribed！", usedTime.String())
	} else {
		fmt.Println("Subscription failed！", err)
		t.Fatal("Subscription failed！", err)
		//return
	}

	time.Sleep(600 * time.Second)
	//等待推送

	start = time.Now()
	res, _, err = r.PubTickers(OP_UNSUBSCRIBE, args)
	if res {
		usedTime := time.Since(start)
		fmt.Println("Unsubscribe successfully！", usedTime.String())
	} else {
		fmt.Println("Cancel Subscription failed！", err)
		t.Fatal("Cancel Subscription failed！", err)
	}

}

// 持仓总量频道 测试
func TestOpenInsterest(t *testing.T) {
	r := prework()
	var err error
	var res bool

	var args []map[string]string
	arg := make(map[string]string)
	arg["instId"] = "LTC-USD-SWAP"

	args = append(args, arg)

	start := time.Now()
	res, _, err = r.PubOpenInsterest(OP_SUBSCRIBE, args)
	if res {
		usedTime := time.Since(start)
		fmt.Println("Successfully subscribed！", usedTime.String())
	} else {
		fmt.Println("Subscription failed！", err)
		t.Fatal("Subscription failed！", err)
		//return
	}

	time.Sleep(60 * time.Second)
	//等待推送

	start = time.Now()
	res, _, err = r.PubOpenInsterest(OP_UNSUBSCRIBE, args)
	if res {
		usedTime := time.Since(start)
		fmt.Println("Unsubscribe successfully！", usedTime.String())
	} else {
		fmt.Println("Cancel Subscription failed！", err)
		t.Fatal("Cancel Subscription failed！", err)
	}

}

// K线频道测试
func TestKLine(t *testing.T) {
	r := prework()
	var err error
	var res bool

	var args []map[string]string
	arg := make(map[string]string)
	arg["instId"] = "BTC-USDT"
	args = append(args, arg)

	// 1分钟K
	period := PERIOD_1MIN

	start := time.Now()
	res, _, err = r.PubKLine(OP_SUBSCRIBE, period, args)
	if res {
		usedTime := time.Since(start)
		fmt.Println("Successfully subscribed！", usedTime.String())
	} else {
		fmt.Println("Subscription failed！", err)
		t.Fatal("Subscription failed！", err)
	}

	time.Sleep(60 * time.Second)
	//等待推送

	start = time.Now()
	res, _, err = r.PubKLine(OP_UNSUBSCRIBE, period, args)
	if res {
		usedTime := time.Since(start)
		fmt.Println("Unsubscribe successfully！", usedTime.String())
	} else {
		fmt.Println("Cancel Subscription failed！", err)
		t.Fatal("Cancel Subscription failed！", err)
	}

}

// 交易频道测试
func TestTrade(t *testing.T) {
	r := prework()
	var err error
	var res bool

	var args []map[string]string
	arg := make(map[string]string)
	arg["instId"] = "BTC-USDT"
	args = append(args, arg)

	start := time.Now()
	res, _, err = r.PubTrade(OP_SUBSCRIBE, args)
	if res {
		usedTime := time.Since(start)
		fmt.Println("Successfully subscribed！", usedTime.String())
	} else {
		fmt.Println("Subscription failed！", err)
		t.Fatal("Subscription failed！", err)
		//return
	}

	time.Sleep(60 * time.Second)
	//等待推送

	start = time.Now()
	res, _, err = r.PubTrade(OP_UNSUBSCRIBE, args)
	if res {
		usedTime := time.Since(start)
		fmt.Println("Unsubscribe successfully！", usedTime.String())
	} else {
		fmt.Println("Cancel Subscription failed！", err)
		t.Fatal("Cancel Subscription failed！", err)
	}

}

// 预估交割/行权价格频道 测试
func TestEstDePrice(t *testing.T) {
	r := prework()
	var err error
	var res bool

	var args []map[string]string
	arg := make(map[string]string)
	arg["instType"] = FUTURES
	arg["uly"] = "BTC-USD"
	args = append(args, arg)

	start := time.Now()
	res, _, err = r.PubEstDePrice(OP_SUBSCRIBE, args)
	if res {
		usedTime := time.Since(start)
		fmt.Println("Successfully subscribed！", usedTime.String())
	} else {
		fmt.Println("Subscription failed！", err)
		t.Fatal("Subscription failed！", err)
		//return
	}

	time.Sleep(60 * time.Second)
	//等待推送

	start = time.Now()
	res, _, err = r.PubEstDePrice(OP_UNSUBSCRIBE, args)
	if res {
		usedTime := time.Since(start)
		fmt.Println("Unsubscribe successfully！", usedTime.String())
	} else {
		fmt.Println("Cancel Subscription failed！", err)
		t.Fatal("Cancel Subscription failed！", err)
	}

}

// 标记价格频道 测试
func TestMarkPrice(t *testing.T) {
	r := prework()
	var err error
	var res bool

	var args []map[string]string
	arg := make(map[string]string)
	arg["instId"] = "BTC-USDT"

	args = append(args, arg)

	start := time.Now()
	res, _, err = r.PubMarkPrice(OP_SUBSCRIBE, args)
	if res {
		usedTime := time.Since(start)
		fmt.Println("Successfully subscribed！", usedTime.String())
	} else {
		fmt.Println("Subscription failed！", err)
		t.Fatal("Subscription failed！", err)
		//return
	}

	time.Sleep(60 * time.Second)
	//等待推送

	start = time.Now()
	res, _, err = r.PubMarkPrice(OP_UNSUBSCRIBE, args)
	if res {
		usedTime := time.Since(start)
		fmt.Println("Unsubscribe successfully！", usedTime.String())
	} else {
		fmt.Println("Cancel Subscription failed！", err)
		t.Fatal("Cancel Subscription failed！", err)
	}

}

// 标记价格K线频道 测试s
func TestMarkPriceCandle(t *testing.T) {
	r := prework()
	var err error
	var res bool

	var args []map[string]string
	arg := make(map[string]string)
	arg["instId"] = "BTC-USDT"
	args = append(args, arg)

	period := PERIOD_1MIN

	start := time.Now()
	res, _, err = r.PubMarkPriceCandle(OP_SUBSCRIBE, period, args)
	if res {
		usedTime := time.Since(start)
		fmt.Println("Successfully subscribed！", usedTime.String())
	} else {
		fmt.Println("Subscription failed！", err)
		t.Fatal("Subscription failed！", err)
		//return
	}

	time.Sleep(60 * time.Second)
	//等待推送

	start = time.Now()
	res, _, err = r.PubMarkPriceCandle(OP_UNSUBSCRIBE, period, args)
	if res {
		usedTime := time.Since(start)
		fmt.Println("Unsubscribe successfully！", usedTime.String())
	} else {
		fmt.Println("Cancel Subscription failed！", err)
		t.Fatal("Cancel Subscription failed！", err)
	}

}

// 限价频道 测试
func TestLimitPrice(t *testing.T) {
	r := prework()
	var err error
	var res bool

	var args []map[string]string
	arg := make(map[string]string)
	arg["instId"] = "BTC-USDT"
	args = append(args, arg)

	start := time.Now()
	res, _, err = r.PubLimitPrice(OP_SUBSCRIBE, args)
	if res {
		usedTime := time.Since(start)
		fmt.Println("Successfully subscribed！", usedTime.String())
	} else {
		fmt.Println("Subscription failed！", err)
		t.Fatal("Subscription failed！", err)
		//return
	}

	time.Sleep(60 * time.Second)
	//等待推送

	start = time.Now()
	res, _, err = r.PubLimitPrice(OP_UNSUBSCRIBE, args)
	if res {
		usedTime := time.Since(start)
		fmt.Println("Unsubscribe successfully！", usedTime.String())
	} else {
		fmt.Println("Cancel Subscription failed！", err)
		t.Fatal("Cancel Subscription failed！", err)
	}

}

// 深度频道 测试
func TestOrderBooks(t *testing.T) {
	r := prework()
	var err error
	var res bool

	/*
		设置关闭深度数据管理
	*/
	// err = r.EnableAutoDepthMgr(false)
	// if err != nil {
	// 	fmt.Println("关闭自动校验失败！")
	// }

	end := make(chan struct{})

	r.AddDepthHook(func(ts time.Time, data DepthData) error {
		// 对于深度类型数据处理的用户可以自定义

		// 检测深度数据是否正常
		key, _ := json.Marshal(data.Arg)
		fmt.Println("个数：", len(data.Data[0].Asks))
		checksum := data.Data[0].Checksum
		fmt.Println("[自定义方法] ", string(key), ", checksum = ", checksum)

		for _, ask := range data.Data[0].Asks {

			arr := strings.Split(ask[0], ".")
			//fmt.Println(arr)
			if len(arr) > 1 && len(arr[1]) > 2 {
				fmt.Println("ask数据异常,", checksum, "ask:", ask)
				t.Fatal()
				end <- struct{}{}
				return nil
			} else {
				fmt.Println("bid数据正常,", checksum, "ask:", ask)
			}

		}

		for _, bid := range data.Data[0].Bids {

			arr := strings.Split(bid[0], ".")
			//fmt.Println(arr)
			if len(arr) > 1 && len(arr[1]) > 2 {
				fmt.Println("bid数据异常,", checksum, "bid:", bid)
				t.Fatal()
				end <- struct{}{}
				return nil
			} else {
				fmt.Println("ask数据正常,", checksum, "bid:", bid)
			}

		}

		// // View the current combined full depth data
		// snapshot, err := r.GetSnapshotByChannel(data)
		// if err != nil {
		// 	t.Fatal("Depth data does not exist！")
		// }
		// // Display the top 5 data of ask/bid
		// fmt.Println(" Ask 5 File data >> ")
		// for _, v := range snapshot.Asks[:5] {
		// 	fmt.Println(" price:", v[0], " amount:", v[1])
		// }
		// fmt.Println(" Bid 5 File data >> ")
		// for _, v := range snapshot.Bids[:5] {
		// 	fmt.Println(" price:", v[0], " amount:", v[1])
		// }
		return nil
	})

	// 可选类型：books books5 books-l2-tbt
	channel := "books50-l2-tbt"

	instIds := []string{"BTC-USDT"}
	for _, instId := range instIds {
		var args []map[string]string
		arg := make(map[string]string)
		arg["instId"] = instId
		args = append(args, arg)

		start := time.Now()
		res, _, err = r.PubOrderBooks(OP_SUBSCRIBE, channel, args)
		if res {
			usedTime := time.Since(start)
			fmt.Println("Successfully subscribed！", usedTime.String())
		} else {
			fmt.Println("Subscription failed！", err)
			t.Fatal("Subscription failed！", err)
		}
	}

	select {
	case <-end:

	}
	//等待推送
	for _, instId := range instIds {
		var args []map[string]string
		arg := make(map[string]string)
		arg["instId"] = instId
		args = append(args, arg)

		start := time.Now()
		res, _, err = r.PubOrderBooks(OP_UNSUBSCRIBE, channel, args)
		if res {
			usedTime := time.Since(start)
			fmt.Println("Unsubscribe successfully！", usedTime.String())
		} else {
			fmt.Println("Cancel Subscription failed！", err)
			t.Fatal("Cancel Subscription failed！", err)
		}
	}

}

// 期权定价频道 测试
func TestOptionSummary(t *testing.T) {
	r := prework()
	var err error
	var res bool

	var args []map[string]string
	arg := make(map[string]string)
	arg["uly"] = "BTC-USD"
	args = append(args, arg)

	start := time.Now()
	res, _, err = r.PubOptionSummary(OP_SUBSCRIBE, args)
	if res {
		usedTime := time.Since(start)
		fmt.Println("Successfully subscribed！", usedTime.String())
	} else {
		fmt.Println("Subscription failed！", err)
		t.Fatal("Subscription failed！", err)
		//return
	}

	time.Sleep(60 * time.Second)
	//等待推送

	start = time.Now()
	res, _, err = r.PubOptionSummary(OP_UNSUBSCRIBE, args)
	if res {
		usedTime := time.Since(start)
		fmt.Println("Unsubscribe successfully！", usedTime.String())
	} else {
		fmt.Println("Cancel Subscription failed！", err)
		t.Fatal("Cancel Subscription failed！", err)
	}

}

// 资金费率 测试
func TestFundRate(t *testing.T) {
	r := prework()
	var err error
	var res bool

	var args []map[string]string
	arg := make(map[string]string)
	arg["instId"] = "BTC-USD-SWAP"
	args = append(args, arg)

	start := time.Now()
	res, _, err = r.PubFundRate(OP_SUBSCRIBE, args)
	if res {
		usedTime := time.Since(start)
		fmt.Println("Successfully subscribed！", usedTime.String())
	} else {
		fmt.Println("Subscription failed！", err)
		t.Fatal("Subscription failed！", err)
		//return
	}

	time.Sleep(600 * time.Second)
	//等待推送

	start = time.Now()
	res, _, err = r.PubFundRate(OP_UNSUBSCRIBE, args)
	if res {
		usedTime := time.Since(start)
		fmt.Println("Unsubscribe successfully！", usedTime.String())
	} else {
		fmt.Println("Cancel Subscription failed！", err)
		t.Fatal("Cancel Subscription failed！", err)
	}

}

// 指数K线频道 测试
func TestKLineIndex(t *testing.T) {
	r := prework()
	var err error
	var res bool

	var args []map[string]string
	arg := make(map[string]string)

	arg["instId"] = "BTC-USDT"
	args = append(args, arg)
	period := PERIOD_1MIN

	start := time.Now()
	res, _, err = r.PubKLineIndex(OP_SUBSCRIBE, period, args)
	if res {
		usedTime := time.Since(start)
		fmt.Println("Successfully subscribed！", usedTime.String())
	} else {
		fmt.Println("Subscription failed！", err)
		t.Fatal("Subscription failed！", err)
		//return
	}

	time.Sleep(60 * time.Second)
	//等待推送

	start = time.Now()
	res, _, err = r.PubKLineIndex(OP_UNSUBSCRIBE, period, args)
	if res {
		usedTime := time.Since(start)
		fmt.Println("Unsubscribe successfully！", usedTime.String())
	} else {
		fmt.Println("Cancel Subscription failed！", err)
		t.Fatal("Cancel Subscription failed！", err)
	}

}

// 指数行情频道 测试
func TestIndexMarket(t *testing.T) {
	r := prework()
	var err error
	var res bool

	var args []map[string]string
	arg := make(map[string]string)
	arg["instId"] = "BTC-USDT"
	args = append(args, arg)

	start := time.Now()
	res, _, err = r.PubIndexTickers(OP_SUBSCRIBE, args)
	if err != nil {
		fmt.Println("Subscription failed！", err)
	}
	usedTime := time.Since(start)
	if res {
		fmt.Println("Successfully subscribed！", usedTime.String())
	} else {
		fmt.Println("Subscription failed！", usedTime.String())
		//return
	}

	time.Sleep(600 * time.Second)
	//等待推送

	start = time.Now()
	res, _, err = r.PubIndexTickers(OP_UNSUBSCRIBE, args)
	if res {
		usedTime := time.Since(start)
		fmt.Println("Unsubscribe successfully！", usedTime.String())
	} else {
		fmt.Println("Cancel Subscription failed！", err)
		t.Fatal("Cancel Subscription failed！", err)
	}

}
