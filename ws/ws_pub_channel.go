package ws

import (
	"errors"
	. "github.com/bulenttokuzlu/v5sdk_go/ws/wImpl"
)

/*
	Product Channel
*/
func (a *WsClient) PubInstruemnts(op string, params []map[string]string, timeOut ...int) (res bool, msg []*Msg, err error) {

	return a.PubChannel(EVENT_BOOK_INSTRUMENTS, op, params, PERIOD_NONE, timeOut...)
}

func (a *WsClient) PubStatus(op string, timeOut ...int) (res bool, msg []*Msg, err error) {
	return a.PubChannel(EVENT_STATUS, op, nil, PERIOD_NONE, timeOut...)
}

/*
	Quotes Channel
*/
func (a *WsClient) PubTickers(op string, params []map[string]string, timeOut ...int) (res bool, msg []*Msg, err error) {

	return a.PubChannel(EVENT_BOOK_TICKERS, op, params, PERIOD_NONE, timeOut...)
}

/*
	Total position channel
*/
func (a *WsClient) PubOpenInsterest(op string, params []map[string]string, timeOut ...int) (res bool, msg []*Msg, err error) {
	return a.PubChannel(EVENT_BOOK_OPEN_INTEREST, op, params, PERIOD_NONE, timeOut...)
}

/*
	K line channel
*/
func (a *WsClient) PubKLine(op string, period Period, params []map[string]string, timeOut ...int) (res bool, msg []*Msg, err error) {

	return a.PubChannel(EVENT_BOOK_KLINE, op, params, period, timeOut...)
}

/*
	Trading channel
*/
func (a *WsClient) PubTrade(op string, params []map[string]string, timeOut ...int) (res bool, msg []*Msg, err error) {

	return a.PubChannel(EVENT_BOOK_TRADE, op, params, PERIOD_NONE, timeOut...)
}

/*
	Estimated delivery/exercise price channel
*/
func (a *WsClient) PubEstDePrice(op string, params []map[string]string, timeOut ...int) (res bool, msg []*Msg, err error) {

	return a.PubChannel(EVENT_BOOK_ESTIMATE_PRICE, op, params, PERIOD_NONE, timeOut...)

}

/*
	Mark price channel
*/
func (a *WsClient) PubMarkPrice(op string, params []map[string]string, timeOut ...int) (res bool, msg []*Msg, err error) {

	return a.PubChannel(EVENT_BOOK_MARK_PRICE, op, params, PERIOD_NONE, timeOut...)
}

/*
	Mark price K-line channel
*/
func (a *WsClient) PubMarkPriceCandle(op string, pd Period, params []map[string]string, timeOut ...int) (res bool, msg []*Msg, err error) {

	return a.PubChannel(EVENT_BOOK_MARK_PRICE_CANDLE_CHART, op, params, pd, timeOut...)
}

/*
	Limit channel
*/
func (a *WsClient) PubLimitPrice(op string, params []map[string]string, timeOut ...int) (res bool, msg []*Msg, err error) {

	return a.PubChannel(EVENT_BOOK_LIMIT_PRICE, op, params, PERIOD_NONE, timeOut...)
}

/*
	Deep channel
*/
func (a *WsClient) PubOrderBooks(op string, channel string, params []map[string]string, timeOut ...int) (res bool, msg []*Msg, err error) {

	switch channel {
	// 400 file snapshots
	case "books":
		return a.PubChannel(EVENT_BOOK_ORDER_BOOK, op, params, PERIOD_NONE, timeOut...)
	// 5 file snapshots
	case "books5":
		return a.PubChannel(EVENT_BOOK_ORDER_BOOK5, op, params, PERIOD_NONE, timeOut...)
	// 400 tbt
	case "books-l2-tbt":
		return a.PubChannel(EVENT_BOOK_ORDER_BOOK_TBT, op, params, PERIOD_NONE, timeOut...)
	// 50 tbt
	case "books50-l2-tbt":
		return a.PubChannel(EVENT_BOOK_ORDER_BOOK50_TBT, op, params, PERIOD_NONE, timeOut...)

	default:
		err = errors.New("Unknown channel")
		return
	}

}

/*
	Option Pricing Channel
*/
func (a *WsClient) PubOptionSummary(op string, params []map[string]string, timeOut ...int) (res bool, msg []*Msg, err error) {

	return a.PubChannel(EVENT_BOOK_OPTION_SUMMARY, op, params, PERIOD_NONE, timeOut...)
}

/*
	Funding rate channel
*/
func (a *WsClient) PubFundRate(op string, params []map[string]string, timeOut ...int) (res bool, msg []*Msg, err error) {

	return a.PubChannel(EVENT_BOOK_FUND_RATE, op, params, PERIOD_NONE, timeOut...)
}

/*
	Index K-Line Channel
*/
func (a *WsClient) PubKLineIndex(op string, pd Period, params []map[string]string, timeOut ...int) (res bool, msg []*Msg, err error) {

	return a.PubChannel(EVENT_BOOK_KLINE_INDEX, op, params, pd, timeOut...)
}

/*
	Index market channel
*/
func (a *WsClient) PubIndexTickers(op string, params []map[string]string, timeOut ...int) (res bool, msg []*Msg, err error) {

	return a.PubChannel(EVENT_BOOK_INDEX_TICKERS, op, params, PERIOD_NONE, timeOut...)
}
