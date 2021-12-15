package ws

import (
	"errors"
	. "github.com/bulenttokuzlu/v5sdk_go/ws/wImpl"
	. "github.com/bulenttokuzlu/v5sdk_go/ws/wInterface"
	"log"
	"runtime/debug"
)

// Judge whether the returned result is successful or unsuccessful
func checkResult(wsReq WSReqData, wsRsps []*Msg) (res bool, err error) {
	defer func() {
		a := recover()
		if a != nil {
			log.Printf("Receive End. Recover msg: %+v", a)
			debug.PrintStack()
		}
		return
	}()

	res = false
	if len(wsRsps) == 0 {
		return
	}

	for _, v := range wsRsps {
		switch v.Info.(type) {
		case ErrData:
			return
		}
		if wsReq.GetType() != v.Info.(WSRspData).MsgType() {
			err = errors.New("Inconsistent message types")
			return
		}
	}

	//Check whether all channels are updated successfully
	if wsReq.GetType() == MSG_NORMAL {
		req, ok := wsReq.(ReqData)
		if !ok {
			log.Println("Type conversion failed", req)
			err = errors.New("Type conversion failed")
			return
		}

		for idx, _ := range req.Args {
			ok := false
			i_req := req.Args[idx]
			//fmt.Println("an examination",i_req)
			for i, _ := range wsRsps {
				info, _ := wsRsps[i].Info.(RspData)
				//fmt.Println("<<",info)
				if info.Event == req.Op && info.Arg["channel"] == i_req["channel"] && info.Arg["instType"] == i_req["instType"] {
					ok = true
					continue
				}
			}
			if !ok {
				err = errors.New("Did not get all the expected return results")
				return
			}
		}
	} else {
		for i, _ := range wsRsps {
			info, _ := wsRsps[i].Info.(JRPCRsp)
			if info.Code != "0" {
				return
			}
		}
	}

	res = true
	return
}
