package trade

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	pb "github.com/h9896/bingo-pkg-protobuf/services/delivery/v1"
	"github.com/h9896/bingo/mocks"
	"github.com/stretchr/testify/assert"
)

func getMockDeliveryTradeService() pb.DeliveryTradeServiceServer {
	return NewDeliveryTradeService(mocks.MockDomain, mocks.MockApiKey, mocks.MockSecret, true, &mocks.MockHTTPClient{})
}

func TestChangePositionMode(t *testing.T) {
	service := getMockDeliveryTradeService()
	request := &pb.ChangePositionModeRequest{
		DualSidePosition: "YES",
		RecvWindow:       5000,
	}

	mocks.GetDoFunc = func(req *http.Request) (resp *http.Response, err error) {
		assert.EqualValues(t, http.MethodPost, req.Method)
		mocks.CheckHeader(t, req.Header)
		assert.EqualValues(t, mocks.MockDomain, req.URL.Host)
		assert.EqualValues(t, "/dapi/v1/positionSide/dual", req.URL.Path)
		params := req.URL.Query()
		mocks.CheckTimestampAndSignature(t, params)
		assert.Contains(t, params["dualSidePosition"], request.GetDualSidePosition())
		assert.Contains(t, params["recvWindow"], fmt.Sprintf("%v", request.GetRecvWindow()))
		data := `{
			"code": 200,
			"msg": "success"
		}`
		resp = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString(data))}
		return
	}
	resp, _ := service.ChangePositionMode(context.Background(), request)
	assert.EqualValues(t, 200, resp.Code)
	assert.EqualValues(t, "success", resp.GetMsg())
}

func TestNewOrder(t *testing.T) {
	service := getMockDeliveryTradeService()
	request := &pb.NewOrderRequest{
		Symbol:           "BTCUSD_200925",
		Side:             pb.OrderSide_BUY,
		PositionSide:     pb.PositionSide_LONG,
		Type:             pb.OrderType_LIMIT,
		TimeInForce:      pb.TimeInForce_FOK,
		ReduceOnly:       "true",
		Quantity:         1.11,
		Price:            598.2,
		NewClientOrderId: "123",
		StopPrice:        0.444,
		ClosePosition:    "true",
		ActivationPrice:  31234.2,
		CallbackRate:     123.1,
		PriceProtect:     "TRUE",
		WorkingType:      pb.WorkingType_MARK_PRICE,
		NewOrderRespType: pb.ResponseType_ACK,
		RecvWindow:       5000,
	}

	mocks.GetDoFunc = func(req *http.Request) (resp *http.Response, err error) {
		assert.EqualValues(t, http.MethodPost, req.Method)
		mocks.CheckHeader(t, req.Header)
		assert.EqualValues(t, mocks.MockDomain, req.URL.Host)
		assert.EqualValues(t, "/dapi/v1/order", req.URL.Path)
		params := req.URL.Query()
		mocks.CheckTimestampAndSignature(t, params)
		assert.Contains(t, params["symbol"], request.GetSymbol())
		assert.Contains(t, params["side"], request.GetSide().String())
		assert.Contains(t, params["type"], request.GetType().String())
		assert.Contains(t, params["timeInForce"], request.GetTimeInForce().String())
		assert.Contains(t, params["workingType"], request.GetWorkingType().String())
		assert.Contains(t, params["newOrderRespType"], request.GetNewOrderRespType().String())
		assert.Contains(t, params["quantity"], fmt.Sprintf("%v", request.GetQuantity()))
		assert.Contains(t, params["price"], fmt.Sprintf("%v", request.GetPrice()))
		assert.Contains(t, params["recvWindow"], fmt.Sprintf("%v", request.GetRecvWindow()))
		data := `{
			"clientOrderId": "testOrder",
			"cumQty": "0",
			"cumBase": "0",
			"executedQty": "0",
			"orderId": 22542179,
			"avgPrice": "0.0",
			"origQty": "10",
			"price": "0",
			"reduceOnly": false,
			"side": "BUY",
			"positionSide": "SHORT", 
			"status": "NEW",
			"stopPrice": "9300",             
			"closePosition": false,         
			"symbol": "BTCUSD_200925",
			"pair": "BTCUSD",
			"timeInForce": "GTC",
			"type": "TRAILING_STOP_MARKET",
			"origType": "TRAILING_STOP_MARKET",
			"activatePrice": "9020",          
			"priceRate": "0.3",                 
			"updateTime": 1566818724722,
			"workingType": "CONTRACT_PRICE",
			"priceProtect": false               
		}`
		resp = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString(data))}
		return
	}
	resp, _ := service.NewOrder(context.Background(), request)
	assert.EqualValues(t, "testOrder", resp.ClientOrderId)
	assert.EqualValues(t, "BTCUSD", resp.Pair)
	assert.EqualValues(t, "GTC", resp.TimeInForce.String())
	assert.EqualValues(t, "TRAILING_STOP_MARKET", resp.Type.String())
	assert.EqualValues(t, "NEW", resp.Status.String())
	assert.EqualValues(t, false, resp.ReduceOnly)
}

func TestCancelOrder(t *testing.T) {
	service := getMockDeliveryTradeService()
	request := &pb.CancelOrderRequest{
		Symbol:            "BTCUSD_PERP",
		OrderId:           124566,
		OrigClientOrderId: "asduuu",
		RecvWindow:        5000,
	}

	mocks.GetDoFunc = func(req *http.Request) (resp *http.Response, err error) {
		assert.EqualValues(t, http.MethodDelete, req.Method)
		mocks.CheckHeader(t, req.Header)
		assert.EqualValues(t, mocks.MockDomain, req.URL.Host)
		assert.EqualValues(t, "/dapi/v1/order", req.URL.Path)
		params := req.URL.Query()
		mocks.CheckTimestampAndSignature(t, params)
		assert.Contains(t, params["symbol"], request.GetSymbol())
		assert.Contains(t, params["origClientOrderId"], request.GetOrigClientOrderId())
		assert.Contains(t, params["orderId"], fmt.Sprintf("%v", request.GetOrderId()))
		assert.Contains(t, params["recvWindow"], fmt.Sprintf("%v", request.GetRecvWindow()))
		data := `{
			"avgPrice": "0.0",
			"clientOrderId": "myOrder1",
			"cumQty": "0",
			"cumBase": "0",
			"executedQty": "0",
			"orderId": 283194212,
			"origQty": "11",
			"origType": "TRAILING_STOP_MARKET",
			"price": "0",
			"reduceOnly": false,
			"side": "BUY",
			"positionSide": "SHORT",            
			"status": "CANCELED",
			"stopPrice": "9300",                
			"closePosition": false,             
			"symbol": "BTCUSD_200925",
			"pair": "BTCUSD",
			"timeInForce": "GTC",
			"type": "TRAILING_STOP_MARKET",
			"activatePrice": "9020",            
			"priceRate": "0.3",                 
			"updateTime": 1571110484038,
			"workingType": "CONTRACT_PRICE",
			"priceProtect": false             
		}`
		resp = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString(data))}
		return
	}
	resp, _ := service.CancelOrder(context.Background(), request)
	assert.EqualValues(t, "myOrder1", resp.ClientOrderId)
	assert.EqualValues(t, "BUY", resp.Side.String())
	assert.EqualValues(t, float64(9020), resp.ActivatePrice)
}

func TestModifyOrder(t *testing.T) {
	service := getMockDeliveryTradeService()
	request := &pb.ModifyOrderRequest{
		Symbol:            "BTCUSD_200925",
		Side:              pb.OrderSide_BUY,
		OrderId:           2988,
		OrigClientOrderId: "123ff",
		Quantity:          1.11,
		Price:             598.2,
		RecvWindow:        5000,
	}

	mocks.GetDoFunc = func(req *http.Request) (resp *http.Response, err error) {
		assert.EqualValues(t, http.MethodPut, req.Method)
		mocks.CheckHeader(t, req.Header)
		assert.EqualValues(t, mocks.MockDomain, req.URL.Host)
		assert.EqualValues(t, "/dapi/v1/order", req.URL.Path)
		params := req.URL.Query()
		mocks.CheckTimestampAndSignature(t, params)
		assert.Contains(t, params["symbol"], request.GetSymbol())
		assert.Contains(t, params["side"], request.GetSide().String())
		assert.Contains(t, params["origClientOrderId"], request.GetOrigClientOrderId())
		assert.Contains(t, params["orderId"], fmt.Sprintf("%v", request.GetOrderId()))
		assert.Contains(t, params["quantity"], fmt.Sprintf("%v", request.GetQuantity()))
		assert.Contains(t, params["price"], fmt.Sprintf("%v", request.GetPrice()))
		assert.Contains(t, params["recvWindow"], fmt.Sprintf("%v", request.GetRecvWindow()))
		data := `{
			"orderId": 20072994037,
			"symbol": "BTCUSD_PERP",
			"pair": "BTCUSD",
			"status": "NEW",
			"clientOrderId": "LJ9R4QZDihCaS8UAOOLpgW",
			"price": "30005",
			"avgPrice": "0.0",
			"origQty": "1",
			"executedQty": "0",
			"cumQty": "0",
			"cumBase": "0",
			"timeInForce": "GTC",
			"type": "LIMIT",
			"reduceOnly": false,
			"closePosition": false,
			"side": "BUY",
			"positionSide": "LONG",
			"stopPrice": "0",
			"workingType": "CONTRACT_PRICE",
			"priceProtect": false,
			"origType": "LIMIT",
			"updateTime": 1629182711600
		}`
		resp = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString(data))}
		return
	}
	resp, _ := service.ModifyOrder(context.Background(), request)
	assert.EqualValues(t, "LJ9R4QZDihCaS8UAOOLpgW", resp.ClientOrderId)
	assert.EqualValues(t, "BTCUSD", resp.Pair)
	assert.EqualValues(t, "GTC", resp.TimeInForce.String())
	assert.EqualValues(t, "LIMIT", resp.Type.String())
	assert.EqualValues(t, 1629182711600, resp.UpdateTime)
	assert.EqualValues(t, false, resp.ReduceOnly)
}

func TestCancelAllOpenOrders(t *testing.T) {
	service := getMockDeliveryTradeService()
	request := &pb.CancelAllOpenOrdersRequest{
		Symbol:     "BTCUSD_PERP",
		RecvWindow: 5000,
	}

	mocks.GetDoFunc = func(req *http.Request) (resp *http.Response, err error) {
		assert.EqualValues(t, http.MethodDelete, req.Method)
		mocks.CheckHeader(t, req.Header)
		assert.EqualValues(t, mocks.MockDomain, req.URL.Host)
		assert.EqualValues(t, "/dapi/v1/allOpenOrders", req.URL.Path)
		params := req.URL.Query()
		mocks.CheckTimestampAndSignature(t, params)
		assert.Contains(t, params["symbol"], request.GetSymbol())
		assert.Contains(t, params["recvWindow"], fmt.Sprintf("%v", request.GetRecvWindow()))
		data := `{
			"code": "200", 
			"msg": "The operation of cancel all open order is done."
		}`
		resp = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString(data))}
		return
	}
	resp, _ := service.CancelAllOpenOrders(context.Background(), request)
	assert.EqualValues(t, 200, resp.Code)
	assert.EqualValues(t, "The operation of cancel all open order is done.", resp.Msg)
}

func TestAutoCancelAllOpenOrder(t *testing.T) {
	service := getMockDeliveryTradeService()
	request := &pb.AutoCancelAllOpenOrdersRequest{
		Symbol:        "BTCUSD_PERP",
		CountdownTime: 21399,
		RecvWindow:    5000,
	}

	mocks.GetDoFunc = func(req *http.Request) (resp *http.Response, err error) {
		assert.EqualValues(t, http.MethodPost, req.Method)
		mocks.CheckHeader(t, req.Header)
		assert.EqualValues(t, mocks.MockDomain, req.URL.Host)
		assert.EqualValues(t, "/dapi/v1/countdownCancelAll", req.URL.Path)
		params := req.URL.Query()
		mocks.CheckTimestampAndSignature(t, params)
		assert.Contains(t, params["symbol"], request.GetSymbol())
		assert.Contains(t, params["countdownTime"], fmt.Sprintf("%v", request.GetCountdownTime()))
		assert.Contains(t, params["recvWindow"], fmt.Sprintf("%v", request.GetRecvWindow()))
		data := `{
			"symbol": "BTCUSD_200925", 
			"countdownTime": "100000"
		}`
		resp = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString(data))}
		return
	}
	resp, _ := service.AutoCancelAllOpenOrder(context.Background(), request)
	assert.EqualValues(t, "BTCUSD_200925", resp.Symbol)
	assert.EqualValues(t, 100000, resp.CountdownTime)
}

func TestChangeInitialLeverage(t *testing.T) {
	service := getMockDeliveryTradeService()
	request := &pb.ChangeInitialLeverageRequest{
		Symbol:     "BTCUSD_PERP",
		Leverage:   21399,
		RecvWindow: 5000,
	}

	mocks.GetDoFunc = func(req *http.Request) (resp *http.Response, err error) {
		assert.EqualValues(t, http.MethodPost, req.Method)
		mocks.CheckHeader(t, req.Header)
		assert.EqualValues(t, mocks.MockDomain, req.URL.Host)
		assert.EqualValues(t, "/dapi/v1/leverage", req.URL.Path)
		params := req.URL.Query()
		mocks.CheckTimestampAndSignature(t, params)
		assert.Contains(t, params["symbol"], request.GetSymbol())
		assert.Contains(t, params["leverage"], fmt.Sprintf("%v", request.GetLeverage()))
		assert.Contains(t, params["recvWindow"], fmt.Sprintf("%v", request.GetRecvWindow()))
		data := `{
			"leverage": 21,
			"maxQty": "1000",  
			"symbol": "BTCUSD_200925"
		}`
		resp = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString(data))}
		return
	}
	resp, _ := service.ChangeInitialLeverage(context.Background(), request)
	assert.EqualValues(t, "BTCUSD_200925", resp.Symbol)
	assert.EqualValues(t, 21, resp.Leverage)
	assert.EqualValues(t, 1000, resp.MaxQty)
}

func TestChangeMarginType(t *testing.T) {
	service := getMockDeliveryTradeService()
	request := &pb.ChangeMarginTypeRequest{
		Symbol:     "BTCUSD_PERP",
		MarginType: pb.MarginType_CROSSED,
		RecvWindow: 5000,
	}

	mocks.GetDoFunc = func(req *http.Request) (resp *http.Response, err error) {
		assert.EqualValues(t, http.MethodPost, req.Method)
		mocks.CheckHeader(t, req.Header)
		assert.EqualValues(t, mocks.MockDomain, req.URL.Host)
		assert.EqualValues(t, "/dapi/v1/marginType", req.URL.Path)
		params := req.URL.Query()
		mocks.CheckTimestampAndSignature(t, params)
		assert.Contains(t, params["symbol"], request.GetSymbol())
		assert.Contains(t, params["marginType"], request.GetMarginType().String())
		assert.Contains(t, params["recvWindow"], fmt.Sprintf("%v", request.GetRecvWindow()))
		data := `{
			"code": 200,
			"msg": "success"
		}`
		resp = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString(data))}
		return
	}
	resp, _ := service.ChangeMarginType(context.Background(), request)
	assert.EqualValues(t, 200, resp.Code)
	assert.EqualValues(t, "success", resp.GetMsg())
}

func TestModifyIsolatedPositionMargin(t *testing.T) {
	service := getMockDeliveryTradeService()
	request := &pb.ModifyIsolatedPositionMarginRequest{
		Symbol:       "BTCUSD_PERP",
		PositionSide: pb.PositionSide_SHORT,
		Amount:       567.2,
		Type:         2,
		RecvWindow:   5000,
	}

	mocks.GetDoFunc = func(req *http.Request) (resp *http.Response, err error) {
		assert.EqualValues(t, http.MethodPost, req.Method)
		mocks.CheckHeader(t, req.Header)
		assert.EqualValues(t, mocks.MockDomain, req.URL.Host)
		assert.EqualValues(t, "/dapi/v1/positionMargin", req.URL.Path)
		params := req.URL.Query()
		mocks.CheckTimestampAndSignature(t, params)
		assert.Contains(t, params["symbol"], request.GetSymbol())
		assert.Contains(t, params["positionSide"], request.GetPositionSide().String())
		assert.Contains(t, params["amount"], fmt.Sprintf("%v", request.GetAmount()))
		assert.Contains(t, params["type"], fmt.Sprintf("%v", request.GetType()))
		assert.Contains(t, params["recvWindow"], fmt.Sprintf("%v", request.GetRecvWindow()))
		data := `{
			"amount": 100.0,
			"code": 200,
			"msg": "Successfully modify position margin.",
			"type": 1
		}`
		resp = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString(data))}
		return
	}
	resp, _ := service.ModifyIsolatedPositionMargin(context.Background(), request)
	assert.EqualValues(t, 200, resp.Code)
	assert.EqualValues(t, 1, resp.Type)
	assert.EqualValues(t, 100.0, resp.Amount)
	assert.EqualValues(t, "Successfully modify position margin.", resp.GetMsg())
}
