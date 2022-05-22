package userdata

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/h9896/bingo-pkg-protobuf/services/delivery/v1"
	"github.com/h9896/bingo/mocks"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/encoding/protojson"
)

func getMockDeliveryUserDataService() pb.DeliveryUserDataServiceServer {
	return NewDeliveryUserDataService(mocks.MockDomain, mocks.MockApiKey, mocks.MockSecret, true, &mocks.MockHTTPClient{})
}

var m = &runtime.JSONPb{
	UnmarshalOptions: protojson.UnmarshalOptions{
		DiscardUnknown: true,
	},
	MarshalOptions: protojson.MarshalOptions{
		UseProtoNames: true,
	},
}

func TestGetPositionMode(t *testing.T) {
	service := getMockDeliveryUserDataService()

	mocks.GetDoFunc = func(req *http.Request) (resp *http.Response, err error) {
		assert.EqualValues(t, http.MethodGet, req.Method)
		mocks.CheckHeader(t, req.Header)
		assert.EqualValues(t, mocks.MockDomain, req.URL.Host)
		assert.EqualValues(t, "/dapi/v1/positionSide/dual", req.URL.Path)
		params := req.URL.Query()
		mocks.CheckTimestampAndSignature(t, params)
		data := `{
			"dualSidePosition": true 
		}`
		resp = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString(data))}
		return
	}
	resp, _ := service.GetPositionMode(context.Background(), &pb.Empty{})
	assert.EqualValues(t, true, resp.DualSidePosition)
}

func TestGetOrderModifyHistory(t *testing.T) {
	service := getMockDeliveryUserDataService()
	request := &pb.GetOrderModifyHistoryRequest{
		Symbol:            "BTCUSD_PERP",
		OrderId:           123,
		OrigClientOrderId: "asdyyy",
		StartTime:         666622,
		EndTime:           32188,
		Limit:             277,
		RecvWindow:        5000,
	}
	data := `[
			{
				"amendmentId": 5363,  
				"symbol": "BTCUSD_PERP",
				"pair": "BTCUSD",
				"orderId": 20072994037,
				"clientOrderId": "LJ9R4QZDihCaS8UAOOLpgW",
				"time": 1629184560899, 
				"amendment": {
					"price": {
						"before": "30004",
						"after": "30003.2"
					},
					"origQty": {
						"before": "1",
						"after": "1"
					},
					"count": 3 
				}
			},
			{
				"amendmentId": 5361,
				"symbol": "BTCUSD_PERP",
				"pair": "BTCUSD",
				"orderId": 20072994037,
				"clientOrderId": "LJ9R4QZDihCaS8UAOOLpgW",
				"time": 1629184533946,
				"amendment": {
					"price": {
						"before": "30005",
						"after": "30004"
					},
					"origQty": {
						"before": "1",
						"after": "1"
					},
					"count": 2
				}
			},
			{
				"amendmentId": 5325,
				"symbol": "BTCUSD_PERP",
				"pair": "BTCUSD",
				"orderId": 20072994037,
				"clientOrderId": "LJ9R4QZDihCaS8UAOOLpgW",
				"time": 1629182711787,
				"amendment": {
					"price": {
						"before": "30002",
						"after": "30005"
					},
					"origQty": {
						"before": "1",
						"after": "1"
					},
					"count": 1
				}
			}
		]`
	mocks.GetDoFunc = func(req *http.Request) (resp *http.Response, err error) {
		assert.EqualValues(t, http.MethodGet, req.Method)
		mocks.CheckHeader(t, req.Header)
		assert.EqualValues(t, mocks.MockDomain, req.URL.Host)
		assert.EqualValues(t, "/dapi/v1/orderAmendment", req.URL.Path)
		params := req.URL.Query()
		mocks.CheckTimestampAndSignature(t, params)
		assert.Contains(t, params["recvWindow"], fmt.Sprintf("%v", request.GetRecvWindow()))
		assert.Contains(t, params["symbol"], fmt.Sprintf("%v", request.GetSymbol()))
		assert.Contains(t, params["startTime"], fmt.Sprintf("%v", request.GetStartTime()))
		assert.Contains(t, params["endTime"], fmt.Sprintf("%v", request.GetEndTime()))
		resp = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString(data))}
		return
	}
	resp, _ := service.GetOrderModifyHistory(context.Background(), request)

	except := []*pb.OrderModifyHistory{}
	m.Unmarshal([]byte(data), &except)
	assert.EqualValues(t, except, resp.OrderModifyHistory)
}

func TestQueryOrder(t *testing.T) {
	service := getMockDeliveryUserDataService()
	request := &pb.QueryOrderRequest{
		Symbol:            "BTCUSD_PERP",
		OrderId:           123,
		OrigClientOrderId: "asdyyy",
		RecvWindow:        5000,
	}

	mocks.GetDoFunc = func(req *http.Request) (resp *http.Response, err error) {
		assert.EqualValues(t, http.MethodGet, req.Method)
		mocks.CheckHeader(t, req.Header)
		assert.EqualValues(t, mocks.MockDomain, req.URL.Host)
		assert.EqualValues(t, "/dapi/v1/order", req.URL.Path)
		params := req.URL.Query()
		mocks.CheckTimestampAndSignature(t, params)
		assert.Contains(t, params["recvWindow"], fmt.Sprintf("%v", request.GetRecvWindow()))
		assert.Contains(t, params["symbol"], fmt.Sprintf("%v", request.GetSymbol()))
		assert.Contains(t, params["orderId"], fmt.Sprintf("%v", request.GetOrderId()))
		assert.Contains(t, params["origClientOrderId"], fmt.Sprintf("%v", request.GetOrigClientOrderId()))
		data := `{
			"avgPrice": "0.0",
			"clientOrderId": "abc",
			"cumBase": "0",
			"executedQty": "0",
			"orderId": 1917641,
			"origQty": "0.40",
			"origType": "TRAILING_STOP_MARKET",
			"price": "0",
			"reduceOnly": false,
			"side": "BUY",
			"status": "NEW",
			"stopPrice": "9300",                
			"closePosition": false,             
			"symbol": "BTCUSD_200925",
			"pair": "BTCUSD",
			"time": 1579276756075,              
			"timeInForce": "GTC",
			"type": "TRAILING_STOP_MARKET",
			"activatePrice": "9020",   
			"priceRate": "0.3",                 
			"updateTime": 1579276756075,        
			"workingType": "CONTRACT_PRICE",
			"priceProtect": false             
		}`
		resp = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString(data))}
		return
	}
	resp, _ := service.QueryOrder(context.Background(), request)
	assert.EqualValues(t, "abc", resp.ClientOrderId)
	assert.EqualValues(t, "TRAILING_STOP_MARKET", resp.Type.String())
	assert.EqualValues(t, "BTCUSD_200925", resp.Symbol)
}

func TestQueryCurrentOpenOrder(t *testing.T) {
	service := getMockDeliveryUserDataService()
	request := &pb.QueryCurrentOpenOrderRequest{
		Symbol:            "BTCUSD_PERP",
		OrderId:           123,
		OrigClientOrderId: "asdyyy",
		RecvWindow:        5000,
	}

	mocks.GetDoFunc = func(req *http.Request) (resp *http.Response, err error) {
		assert.EqualValues(t, http.MethodGet, req.Method)
		mocks.CheckHeader(t, req.Header)
		assert.EqualValues(t, mocks.MockDomain, req.URL.Host)
		assert.EqualValues(t, "/dapi/v1/openOrder", req.URL.Path)
		params := req.URL.Query()
		mocks.CheckTimestampAndSignature(t, params)
		assert.Contains(t, params["recvWindow"], fmt.Sprintf("%v", request.GetRecvWindow()))
		assert.Contains(t, params["symbol"], fmt.Sprintf("%v", request.GetSymbol()))
		assert.Contains(t, params["orderId"], fmt.Sprintf("%v", request.GetOrderId()))
		assert.Contains(t, params["origClientOrderId"], fmt.Sprintf("%v", request.GetOrigClientOrderId()))
		data := `{
			"avgPrice": "0.0",              
			"clientOrderId": "abc",             
			"cumBase": "0",                     
			"executedQty": "0",                 
			"orderId": 1917641,                 
			"origQty": "0.40",                      
			"origType": "TRAILING_STOP_MARKET",
			"price": "0",
			"reduceOnly": false,
			"side": "BUY",
			"positionSide": "SHORT",
			"status": "NEW",
			"stopPrice": "9300",               
			"closePosition": false,             
			"symbol": "BTCUSD_200925",
			"pair": "BTCUSD",
			"time": 1579276756075,             
			"timeInForce": "GTC",
			"type": "TRAILING_STOP_MARKET",
			"activatePrice": "9020",          
			"priceRate": "0.3",                                    
			"updateTime": 1579276756075,        
			"workingType": "CONTRACT_PRICE",
			"priceProtect": false              
		}`
		resp = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString(data))}
		return
	}
	resp, _ := service.QueryCurrentOpenOrder(context.Background(), request)
	assert.EqualValues(t, "abc", resp.ClientOrderId)
	assert.EqualValues(t, 1917641, resp.OrderId)
	assert.EqualValues(t, "TRAILING_STOP_MARKET", resp.Type.String())
	assert.EqualValues(t, "BTCUSD_200925", resp.Symbol)
}

func TestCurrentAllOpenOrders(t *testing.T) {
	service := getMockDeliveryUserDataService()
	request := &pb.CurrentAllOpenOrdersRequest{
		Symbol:     "BTCUSD_PERP",
		Pair:       "BTC",
		RecvWindow: 5000,
	}
	data := `[
		{
		  "avgPrice": "0.0",
		  "clientOrderId": "abc",
		  "cumBase": "0",
		  "executedQty": "0",
		  "orderId": 1917641,
		  "origQty": "0.40",
		  "origType": "TRAILING_STOP_MARKET",
		  "price": "0",
		  "reduceOnly": false,
		  "side": "BUY",
		  "positionSide": "SHORT",
		  "status": "NEW",
		  "stopPrice": "9300",               
		  "closePosition": false,            
		  "symbol": "BTCUSD_200925",
		  "time": 1579276756075,             
		  "timeInForce": "GTC",
		  "type": "TRAILING_STOP_MARKET",
		  "activatePrice": "9020",          
		  "priceRate": "0.3",             
		  "updateTime": 1579276756075,        
		  "workingType": "CONTRACT_PRICE",
		  "priceProtect": false               
		}
	  ]`
	mocks.GetDoFunc = func(req *http.Request) (resp *http.Response, err error) {
		assert.EqualValues(t, http.MethodGet, req.Method)
		mocks.CheckHeader(t, req.Header)
		assert.EqualValues(t, mocks.MockDomain, req.URL.Host)
		assert.EqualValues(t, "/dapi/v1/openOrders", req.URL.Path)
		params := req.URL.Query()
		mocks.CheckTimestampAndSignature(t, params)
		assert.Contains(t, params["recvWindow"], fmt.Sprintf("%v", request.GetRecvWindow()))
		assert.Contains(t, params["symbol"], fmt.Sprintf("%v", request.GetSymbol()))
		assert.Contains(t, params["pair"], fmt.Sprintf("%v", request.GetPair()))

		resp = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString(data))}
		return
	}
	resp, _ := service.CurrentAllOpenOrders(context.Background(), request)
	except := []*pb.QueryCurrentOpenOrderResponse{}
	m.Unmarshal([]byte(data), &except)
	assert.EqualValues(t, except, resp.CurrentAllOpenOrders)
}

func TestAllOrders(t *testing.T) {
	service := getMockDeliveryUserDataService()
	request := &pb.AllOrdersRequest{
		Symbol:     "BTCUSD_PERP",
		Pair:       "BTC",
		OrderId:    2344,
		StartTime:  7788,
		EndTime:    9087,
		Limit:      8876,
		RecvWindow: 5000,
	}
	data := `[
		{
		  "avgPrice": "0.0",
		  "clientOrderId": "abc",
		  "cumBase": "0",
		  "executedQty": "0",
		  "orderId": 1917641,
		  "origQty": "0.40",
		  "origType": "TRAILING_STOP_MARKET",
		  "price": "0",
		  "reduceOnly": false,
		  "side": "BUY",
		  "positionSide": "SHORT",
		  "status": "NEW",
		  "stopPrice": "9300",               
		  "closePosition": false,             
		  "symbol": "BTCUSD_200925",
		  "pair": "BTCUSD",
		  "time": 1579276756075,             
		  "timeInForce": "GTC",
		  "type": "TRAILING_STOP_MARKET",
		  "activatePrice": "9020",         
		  "priceRate": "0.3",                
		  "updateTime": 1579276756075,        
		  "workingType": "CONTRACT_PRICE",
		  "priceProtect": false              
		}
	  ]`
	mocks.GetDoFunc = func(req *http.Request) (resp *http.Response, err error) {
		assert.EqualValues(t, http.MethodGet, req.Method)
		mocks.CheckHeader(t, req.Header)
		assert.EqualValues(t, mocks.MockDomain, req.URL.Host)
		assert.EqualValues(t, "/dapi/v1/allOrders", req.URL.Path)
		params := req.URL.Query()
		mocks.CheckTimestampAndSignature(t, params)
		assert.Contains(t, params["recvWindow"], fmt.Sprintf("%v", request.GetRecvWindow()))
		assert.Contains(t, params["symbol"], fmt.Sprintf("%v", request.GetSymbol()))
		assert.Contains(t, params["pair"], fmt.Sprintf("%v", request.GetPair()))
		assert.Contains(t, params["startTime"], fmt.Sprintf("%v", request.GetStartTime()))
		assert.Contains(t, params["endTime"], fmt.Sprintf("%v", request.GetEndTime()))
		assert.Contains(t, params["limit"], fmt.Sprintf("%v", request.GetLimit()))
		assert.Contains(t, params["orderId"], fmt.Sprintf("%v", request.GetOrderId()))
		resp = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString(data))}
		return
	}
	resp, _ := service.AllOrders(context.Background(), request)
	except := []*pb.QueryCurrentOpenOrderResponse{}
	m.Unmarshal([]byte(data), &except)
	assert.EqualValues(t, except, resp.AllOrders)
}

func TestFuturesAccountBalance(t *testing.T) {
	service := getMockDeliveryUserDataService()
	data := `[
		{
			"accountAlias": "SgsR",    
			"asset": "BTC",
			"balance": "0.00250000",
			"withdrawAvailable": "0.00250000",
			"crossWalletBalance": "0.00241969",
			"crossUnPnl": "0.00000000",
			"availableBalance": "0.00241969",
			"updateTime": 1592468353979
		}
	]`
	mocks.GetDoFunc = func(req *http.Request) (resp *http.Response, err error) {
		assert.EqualValues(t, http.MethodGet, req.Method)
		mocks.CheckHeader(t, req.Header)
		assert.EqualValues(t, mocks.MockDomain, req.URL.Host)
		assert.EqualValues(t, "/dapi/v1/balance", req.URL.Path)
		params := req.URL.Query()
		mocks.CheckTimestampAndSignature(t, params)
		resp = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString(data))}
		return
	}
	resp, _ := service.FuturesAccountBalance(context.Background(), &pb.Empty{})
	except := []*pb.Balance{}
	m.Unmarshal([]byte(data), &except)
	assert.EqualValues(t, except, resp.FuturesAccountBalance)
}

func TestAccountInformation(t *testing.T) {
	service := getMockDeliveryUserDataService()
	data := `{
		"assets": [
			{
				"asset": "BTC",  
				"walletBalance": "0.00241969",  
				"unrealizedProfit": "0.00000000",  
				"marginBalance": "0.00241969",  
				"maintMargin": "0.00000000",    
				"initialMargin": "0.00000000", 
				"positionInitialMargin": "0.00000000", 
				"openOrderInitialMargin": "0.00000000",  
				"maxWithdrawAmount": "0.00241969", 
				"crossWalletBalance": "0.00241969", 
				"crossUnPnl": "0.00000000",
				"availableBalance": "0.00241969" 
			}
		 ],
		 "positions": [
			 {
				"symbol": "BTCUSD_201225",
				"positionAmt":"0", 
				"initialMargin": "0",
				"maintMargin": "0",
				"unrealizedProfit": "0.00000000",
				"positionInitialMargin": "0",
				"openOrderInitialMargin": "0",
				"leverage": "125",
				"isolated": false,
				"positionSide": "BOTH", 
				"entryPrice": "0.0",
				"maxQty": "50", 
				"updateTime": 0
			},
			{
				"symbol": "BTCUSD_201225",
				"positionAmt":"0",
				"initialMargin": "0",
				"maintMargin": "0",
				"unrealizedProfit": "0.00000000",
				"positionInitialMargin": "0",
				"openOrderInitialMargin": "0",
				"leverage": "125",
				"isolated": false,
				"positionSide": "LONG", 
				"entryPrice": "0.0",
				"maxQty": "50",
				"updateTime": 0
			},
			{
				"symbol": "BTCUSD_201225",
				"positionAmt":"0",
				"initialMargin": "0",
				"maintMargin": "0",
				"unrealizedProfit": "0.00000000",
				"positionInitialMargin": "0",
				"openOrderInitialMargin": "0",
				"leverage": "125",
				"isolated": false,
				"positionSide": "SHORT",  
				"entryPrice": "0.0",
				"maxQty": "50",
				"updateTime":1627026881327
			}
		 ],
		 "canDeposit": true,
		 "canTrade": true,
		 "canWithdraw": true,
		 "feeTier": 2,
		 "updateTime": 0
	 }`
	mocks.GetDoFunc = func(req *http.Request) (resp *http.Response, err error) {
		assert.EqualValues(t, http.MethodGet, req.Method)
		mocks.CheckHeader(t, req.Header)
		assert.EqualValues(t, mocks.MockDomain, req.URL.Host)
		assert.EqualValues(t, "/dapi/v1/account", req.URL.Path)
		params := req.URL.Query()
		mocks.CheckTimestampAndSignature(t, params)
		resp = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString(data))}
		return
	}
	resp, _ := service.AccountInformation(context.Background(), &pb.Empty{})
	except := &pb.AccountInformationResponse{}
	m.Unmarshal([]byte(data), except)
	assert.EqualValues(t, except, resp)
}

func TestPositionInformation(t *testing.T) {
	service := getMockDeliveryUserDataService()
	request := &pb.PositionInformationRequest{
		MarginAsset: "BTCUSD_PERP",
		Pair:        "BTC",
		RecvWindow:  5000,
	}
	data := `[
		{
			"symbol": "BTCUSD_201225",
			"positionAmt": "0",
			"entryPrice": "0.0",
			"markPrice": "0.00000000",
			"unRealizedProfit": "0.00000000",
			"liquidationPrice": "0",
			"leverage": "125",
			"maxQty": "50", 
			"marginType": "cross",
			"isolatedMargin": "0.00000000",
			"isAutoAddMargin": "false",
			"positionSide": "BOTH",
			"updateTime": 0
		},
		{
			"symbol": "BTCUSD_201225",
			"positionAmt": "1",
			"entryPrice": "11707.70000003",
			"markPrice": "11788.66626667",
			"unRealizedProfit": "0.00005866",
			"liquidationPrice": "11667.63509587",
			"leverage": "125",
			"maxQty": "50",
			"marginType": "cross",
			"isolatedMargin": "0.00012357",
			"isAutoAddMargin": "false",
			"positionSide": "LONG",
			"updateTime": 0
		},
		{
			"symbol": "BTCUSD_201225",
			"positionAmt": "0",
			"entryPrice": "0.0",
			"markPrice": "0.00000000",
			"unRealizedProfit": "0.00000000",
			"liquidationPrice": "0",
			"leverage": "125",
			"maxQty": "50",
			"marginType": "cross",
			"isolatedMargin": "0.00000000",
			"isAutoAddMargin": "false",
			"positionSide": "SHORT",
			"updateTime":1627026881327
	  }
	]`
	mocks.GetDoFunc = func(req *http.Request) (resp *http.Response, err error) {
		assert.EqualValues(t, http.MethodGet, req.Method)
		mocks.CheckHeader(t, req.Header)
		assert.EqualValues(t, mocks.MockDomain, req.URL.Host)
		assert.EqualValues(t, "/dapi/v1/positionRisk", req.URL.Path)
		params := req.URL.Query()
		mocks.CheckTimestampAndSignature(t, params)
		assert.Contains(t, params["recvWindow"], fmt.Sprintf("%v", request.GetRecvWindow()))
		assert.Contains(t, params["marginAsset"], fmt.Sprintf("%v", request.GetMarginAsset()))
		assert.Contains(t, params["pair"], fmt.Sprintf("%v", request.GetPair()))

		resp = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString(data))}
		return
	}
	resp, err := service.PositionInformation(context.Background(), request)
	log.Println(err)
	except := []*pb.PositionString{}
	m.Unmarshal([]byte(data), &except)
	assert.EqualValues(t, except, resp.Positions)
}

func TestGetPositionMarginChangeHistory(t *testing.T) {
	service := getMockDeliveryUserDataService()
	request := &pb.GetPositionMarginChangeHistoryRequest{
		Symbol:     "BTCUSD_PERP",
		Type:       1,
		StartTime:  12355,
		EndTime:    48832,
		Limit:      23,
		RecvWindow: 5000,
	}
	data := `[
		{
			"amount": "23.36332311",
			"asset": "BTC",
			"symbol": "BTCUSD_200925",
			"time": 1578047897183,
			"type": 1,
			"positionSide": "BOTH"
		},
		{
			"amount": "100",
			"asset": "BTC",
			"symbol": "BTCUSD_200925",
			"time": 1578047900425,
			"type": 1,
			"positionSide": "LONG"
		}
	]`
	mocks.GetDoFunc = func(req *http.Request) (resp *http.Response, err error) {
		assert.EqualValues(t, http.MethodGet, req.Method)
		mocks.CheckHeader(t, req.Header)
		assert.EqualValues(t, mocks.MockDomain, req.URL.Host)
		assert.EqualValues(t, "/dapi/v1/positionMargin/history", req.URL.Path)
		params := req.URL.Query()
		mocks.CheckTimestampAndSignature(t, params)
		assert.Contains(t, params["recvWindow"], fmt.Sprintf("%v", request.GetRecvWindow()))
		assert.Contains(t, params["symbol"], fmt.Sprintf("%v", request.GetSymbol()))
		assert.Contains(t, params["startTime"], fmt.Sprintf("%v", request.GetStartTime()))
		assert.Contains(t, params["endTime"], fmt.Sprintf("%v", request.GetEndTime()))
		assert.Contains(t, params["limit"], fmt.Sprintf("%v", request.GetLimit()))

		resp = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString(data))}
		return
	}
	resp, err := service.GetPositionMarginChangeHistory(context.Background(), request)
	log.Println(err)
	except := []*pb.PositionMargin{}
	m.Unmarshal([]byte(data), &except)
	assert.EqualValues(t, except, resp.PositionMargins)
}

func TestAccountTradeList(t *testing.T) {
	service := getMockDeliveryUserDataService()
	request := &pb.AccountTradeListRequest{
		Symbol:     "BTCUSD_PERP",
		Pair:       "BTC",
		FromId:     2345,
		StartTime:  12355,
		EndTime:    48832,
		Limit:      23,
		RecvWindow: 5000,
	}
	data := `[
		{
			"symbol": "BTCUSD_200626",
			"id": 6,
			"orderId": 28,
			"pair": "BTCUSD",
			"side": "SELL",
			"price": "8800",
			"qty": "1",
			"realizedPnl": "0",
			"marginAsset": "BTC",
			"baseQty": "0.01136364",
			"commission": "0.00000454",
			"commissionAsset": "BTC",
			"time": 1590743483586,
			"positionSide": "BOTH",
			"buyer": false,
			"maker": false
		}
	]`
	mocks.GetDoFunc = func(req *http.Request) (resp *http.Response, err error) {
		assert.EqualValues(t, http.MethodGet, req.Method)
		mocks.CheckHeader(t, req.Header)
		assert.EqualValues(t, mocks.MockDomain, req.URL.Host)
		assert.EqualValues(t, "/dapi/v1/userTrades", req.URL.Path)
		params := req.URL.Query()
		mocks.CheckTimestampAndSignature(t, params)
		assert.Contains(t, params["recvWindow"], fmt.Sprintf("%v", request.GetRecvWindow()))
		assert.Contains(t, params["symbol"], fmt.Sprintf("%v", request.GetSymbol()))
		assert.Contains(t, params["pair"], fmt.Sprintf("%v", request.GetPair()))
		assert.Contains(t, params["startTime"], fmt.Sprintf("%v", request.GetStartTime()))
		assert.Contains(t, params["endTime"], fmt.Sprintf("%v", request.GetEndTime()))
		assert.Contains(t, params["limit"], fmt.Sprintf("%v", request.GetLimit()))
		assert.Contains(t, params["fromId"], fmt.Sprintf("%v", request.GetFromId()))

		resp = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString(data))}
		return
	}
	resp, err := service.AccountTradeList(context.Background(), request)
	log.Println(err)
	except := []*pb.AccountTrade{}
	m.Unmarshal([]byte(data), &except)
	assert.EqualValues(t, except, resp.AccountTrades)
}

func TestGetIncomeHistory(t *testing.T) {
	service := getMockDeliveryUserDataService()
	request := &pb.GetIncomeHistoryRequest{
		Symbol:     "BTCUSD_PERP",
		IncomeType: pb.IncomeType_COMMISSION,
		StartTime:  12355,
		EndTime:    48832,
		Limit:      23,
		RecvWindow: 5000,
	}
	data := `[
		{
			"symbol": "",               
			"incomeType": "TRANSFER",   
			"income": "-0.37500000",    
			"asset": "BTC",             
			"info":"WITHDRAW",          
			"time": 1570608000000,
			"tranId":"9689322392",      
			"tradeId":""                
		},
		{
			"symbol": "BTCUSD_200925",
			"incomeType": "COMMISSION", 
			"income": "-0.01000000",
			"asset": "BTC",
			"info":"",
			"time": 1570636800000,
			"tranId":"9689322392",
			"tradeId":"2059192"
		}
	]`
	mocks.GetDoFunc = func(req *http.Request) (resp *http.Response, err error) {
		assert.EqualValues(t, http.MethodGet, req.Method)
		mocks.CheckHeader(t, req.Header)
		assert.EqualValues(t, mocks.MockDomain, req.URL.Host)
		assert.EqualValues(t, "/dapi/v1/income", req.URL.Path)
		params := req.URL.Query()
		mocks.CheckTimestampAndSignature(t, params)
		assert.Contains(t, params["recvWindow"], fmt.Sprintf("%v", request.GetRecvWindow()))
		assert.Contains(t, params["symbol"], fmt.Sprintf("%v", request.GetSymbol()))
		assert.Contains(t, params["incomeType"], fmt.Sprintf("%v", request.GetIncomeType().String()))
		assert.Contains(t, params["startTime"], fmt.Sprintf("%v", request.GetStartTime()))
		assert.Contains(t, params["endTime"], fmt.Sprintf("%v", request.GetEndTime()))
		assert.Contains(t, params["limit"], fmt.Sprintf("%v", request.GetLimit()))

		resp = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString(data))}
		return
	}
	resp, err := service.GetIncomeHistory(context.Background(), request)
	log.Println(err)
	except := []*pb.IncomeHistory{}
	m.Unmarshal([]byte(data), &except)
	assert.EqualValues(t, except, resp.Incomes)
}

func TestNotionalBracketForSymbol(t *testing.T) {
	service := getMockDeliveryUserDataService()
	request := &pb.NotionalBracketForSymbolRequest{
		Pair:       "BTCUSD",
		RecvWindow: 5000,
	}
	data := `[
		{
			"pair": "BTCUSD",
			"brackets": [
				{
					"bracket": 1,   
					"initialLeverage": 125,  
					"qtyCap": 50,  
					"qtylFloor": 0,  
					"maintMarginRatio": 0.004,
					"cum": 0.0 
				}
			]
		}
	]`
	mocks.GetDoFunc = func(req *http.Request) (resp *http.Response, err error) {
		assert.EqualValues(t, http.MethodGet, req.Method)
		mocks.CheckHeader(t, req.Header)
		assert.EqualValues(t, mocks.MockDomain, req.URL.Host)
		assert.EqualValues(t, "/dapi/v1/leverageBracket", req.URL.Path)
		params := req.URL.Query()
		mocks.CheckTimestampAndSignature(t, params)
		assert.Contains(t, params["recvWindow"], fmt.Sprintf("%v", request.GetRecvWindow()))
		assert.Contains(t, params["pair"], fmt.Sprintf("%v", request.GetPair()))

		resp = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString(data))}
		return
	}
	resp, err := service.NotionalBracketForSymbol(context.Background(), request)
	log.Println(err)
	except := []*pb.NotionalBracketForSymbol{}
	m.Unmarshal([]byte(data), &except)
	assert.EqualValues(t, except, resp.Brackets)
}

func TestNotionalBracketForPair(t *testing.T) {
	service := getMockDeliveryUserDataService()
	request := &pb.NotionalBracketForPairRequest{
		Symbol:     "BTCUSD_PERP",
		RecvWindow: 5000,
	}
	data := `[
		{
			"symbol": "BTCUSD_PERP",
			"brackets": [
				{
					"bracket": 1,   
					"initialLeverage": 125,  
					"qtyCap": 50,  
					"qtylFloor": 0, 
					"maintMarginRatio": 0.004,
					"cum": 0.0 
				}
			]
		}
	]`
	mocks.GetDoFunc = func(req *http.Request) (resp *http.Response, err error) {
		assert.EqualValues(t, http.MethodGet, req.Method)
		mocks.CheckHeader(t, req.Header)
		assert.EqualValues(t, mocks.MockDomain, req.URL.Host)
		assert.EqualValues(t, "/dapi/v2/leverageBracket", req.URL.Path)
		params := req.URL.Query()
		mocks.CheckTimestampAndSignature(t, params)
		assert.Contains(t, params["recvWindow"], fmt.Sprintf("%v", request.GetRecvWindow()))
		assert.Contains(t, params["symbol"], fmt.Sprintf("%v", request.GetSymbol()))

		resp = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString(data))}
		return
	}
	resp, err := service.NotionalBracketForPair(context.Background(), request)
	log.Println(err)
	except := []*pb.NotionalBracketForPair{}
	m.Unmarshal([]byte(data), &except)
	assert.EqualValues(t, except, resp.Brackets)
}
