package trade

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/h9896/bingo-pkg-protobuf/services/delivery/v1"
	"github.com/h9896/bingo/delivery"
	"github.com/h9896/bingo/rpc"
	"google.golang.org/protobuf/encoding/protojson"
)

type deliveryTradeService struct {
	httpclient rpc.GenericHttpClient
	domain     string
	secret     string
	m          *runtime.JSONPb
}

func NewDeliveryTradeService(domain, apikey, secret string, useSSL bool, client rpc.HTTPClient) pb.DeliveryTradeServiceServer {
	service := &deliveryTradeService{
		domain: domain,
		secret: secret,
		m: &runtime.JSONPb{
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
		},
	}
	if client == nil {
		service.httpclient = rpc.NewGenericHttpClient(apikey, useSSL, nil)
	} else {
		service.httpclient = rpc.NewGenericHttpClient(apikey, useSSL, client)
	}

	return service
}

// Change user's position mode (Hedge Mode or One-way Mode ) on EVERY symbol
func (s *deliveryTradeService) ChangePositionMode(ctx context.Context, request *pb.ChangePositionModeRequest) (*pb.ChangePositionModeResponse, error) {
	endpoint := fmt.Sprintf("%s/%s", s.domain, delivery.EntryPointPositionMode)
	body := []*rpc.HttpParameter{
		{Key: "dualSidePosition", Val: request.GetDualSidePosition()},
	}
	if request.GetRecvWindow() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "recvWindow", Val: fmt.Sprintf("%v", request.GetRecvWindow())})
	}
	req := s.httpclient.GetHttpRequest(rpc.SetEndpoint(endpoint), rpc.SetMethod("post"),
		rpc.SetParams(body...), rpc.SetPrivate(), rpc.SetTimestamp(), rpc.SetSignature(s.secret))

	resp, err := s.httpclient.ExecuteHttpOperation(ctx, req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Something wrong")
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	out := &pb.ChangePositionModeResponse{}

	err = s.m.Unmarshal(respBody, out)

	if err != nil {
		return nil, err
	}
	return out, nil
}

// Send in a new order.
func (s *deliveryTradeService) NewOrder(ctx context.Context, request *pb.NewOrderRequest) (*pb.NewOrderResponse, error) {
	endpoint := fmt.Sprintf("%s/%s", s.domain, delivery.EntryPointOrder)
	body := []*rpc.HttpParameter{
		{Key: "symbol", Val: request.GetSymbol()},
		{Key: "side", Val: request.GetSide().String()},
		{Key: "type", Val: request.GetType().String()},
	}

	if request.GetPositionSide() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "positionSide", Val: request.GetPositionSide().String()})
	}

	if request.GetQuantity() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "quantity", Val: fmt.Sprintf("%v", request.GetQuantity())})
	}

	if request.GetReduceOnly() != "" {
		body = append(body, &rpc.HttpParameter{Key: "reduceOnly", Val: request.GetReduceOnly()})
	}

	if request.GetPrice() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "price", Val: fmt.Sprintf("%v", request.GetPrice())})
	}

	if request.GetNewClientOrderId() != "" {
		body = append(body, &rpc.HttpParameter{Key: "newClientOrderId", Val: request.GetNewClientOrderId()})
	}

	if request.GetStopPrice() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "stopPrice", Val: fmt.Sprintf("%v", request.GetStopPrice())})
	}

	if request.GetClosePosition() != "" {
		body = append(body, &rpc.HttpParameter{Key: "closePosition", Val: request.GetClosePosition()})
	}

	if request.GetActivationPrice() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "activationPrice", Val: fmt.Sprintf("%v", request.GetActivationPrice())})
	}

	if request.GetCallbackRate() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "callbackRate", Val: fmt.Sprintf("%v", request.GetCallbackRate())})
	}

	if request.GetWorkingType() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "workingType", Val: request.GetWorkingType().String()})
	}

	if request.GetPriceProtect() != "" {
		body = append(body, &rpc.HttpParameter{Key: "priceProtect", Val: request.GetPriceProtect()})
	}

	if request.GetNewOrderRespType() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "newOrderRespType", Val: request.GetNewOrderRespType().String()})
	}

	if request.GetTimeInForce() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "timeInForce", Val: request.GetTimeInForce().String()})
	}

	if request.GetRecvWindow() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "recvWindow", Val: fmt.Sprintf("%v", request.GetRecvWindow())})
	}

	req := s.httpclient.GetHttpRequest(rpc.SetEndpoint(endpoint), rpc.SetMethod("post"),
		rpc.SetParams(body...), rpc.SetPrivate(), rpc.SetTimestamp(), rpc.SetSignature(s.secret))

	resp, err := s.httpclient.ExecuteHttpOperation(ctx, req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Something wrong")
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	out := &pb.NewOrderResponse{}

	err = s.m.Unmarshal(respBody, out)

	if err != nil {
		return nil, err
	}
	return out, nil
}

// Cancel an active order.
func (s *deliveryTradeService) CancelOrder(ctx context.Context, request *pb.CancelOrderRequest) (*pb.CancelOrderResponse, error) {
	endpoint := fmt.Sprintf("%s/%s", s.domain, delivery.EntryPointOrder)
	body := []*rpc.HttpParameter{
		{Key: "symbol", Val: request.GetSymbol()},
	}

	if request.GetRecvWindow() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "recvWindow", Val: fmt.Sprintf("%v", request.GetRecvWindow())})
	}

	if request.GetOrderId() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "orderId", Val: fmt.Sprintf("%v", request.GetOrderId())})
	}

	if request.GetOrigClientOrderId() != "" {
		body = append(body, &rpc.HttpParameter{Key: "origClientOrderId", Val: request.GetOrigClientOrderId()})
	}

	req := s.httpclient.GetHttpRequest(rpc.SetEndpoint(endpoint), rpc.SetMethod("delete"),
		rpc.SetParams(body...), rpc.SetPrivate(), rpc.SetTimestamp(), rpc.SetSignature(s.secret))

	resp, err := s.httpclient.ExecuteHttpOperation(ctx, req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Something wrong")
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	out := &pb.CancelOrderResponse{}

	err = s.m.Unmarshal(respBody, out)

	if err != nil {
		return nil, err
	}
	return out, nil
}

// Order modify function, currently only LIMIT order modification is supported,
// modified orders will be reordered in the match queue
func (s *deliveryTradeService) ModifyOrder(ctx context.Context, request *pb.ModifyOrderRequest) (*pb.ModifyOrderResponse, error) {
	endpoint := fmt.Sprintf("%s/%s", s.domain, delivery.EntryPointOrder)
	body := []*rpc.HttpParameter{
		{Key: "symbol", Val: request.GetSymbol()},
		{Key: "side", Val: request.GetSide().String()},
	}

	if request.GetRecvWindow() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "recvWindow", Val: fmt.Sprintf("%v", request.GetRecvWindow())})
	}

	if request.GetOrderId() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "orderId", Val: fmt.Sprintf("%v", request.GetOrderId())})
	}

	if request.GetOrigClientOrderId() != "" {
		body = append(body, &rpc.HttpParameter{Key: "origClientOrderId", Val: request.GetOrigClientOrderId()})
	}

	if request.GetQuantity() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "quantity", Val: fmt.Sprintf("%v", request.GetQuantity())})
	}

	if request.GetPrice() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "price", Val: fmt.Sprintf("%v", request.GetPrice())})
	}

	req := s.httpclient.GetHttpRequest(rpc.SetEndpoint(endpoint), rpc.SetMethod("put"),
		rpc.SetParams(body...), rpc.SetPrivate(), rpc.SetTimestamp(), rpc.SetSignature(s.secret))

	resp, err := s.httpclient.ExecuteHttpOperation(ctx, req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Something wrong")
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	out := &pb.ModifyOrderResponse{}

	err = s.m.Unmarshal(respBody, out)

	if err != nil {
		return nil, err
	}
	return out, nil
}

// Place Multiple Orders
func (s *deliveryTradeService) PlaceMultipleOrders(ctx context.Context, request *pb.PlaceMultipleOrdersRequest) (*pb.PlaceMultipleOrdersResponse, error) {
	endpoint := fmt.Sprintf("%s/%s", s.domain, delivery.EntryPointMultipleOrders)
	batch := ""
	for _, val := range request.GetBatchOrders() {
		subBody, err := s.m.Marshal(val)
		if err != nil {
			return nil, err
		}
		if batch != "" {
			batch += fmt.Sprintf(",%s", string(subBody[:]))
		} else {
			batch += string(subBody[:])
		}
	}
	body := []*rpc.HttpParameter{
		{Key: "batchOrders", Val: fmt.Sprintf("[%s]", batch)},
	}

	if request.GetRecvWindow() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "recvWindow", Val: fmt.Sprintf("%v", request.GetRecvWindow())})
	}

	req := s.httpclient.GetHttpRequest(rpc.SetEndpoint(endpoint), rpc.SetMethod("POST"),
		rpc.SetParams(body...), rpc.SetPrivate(), rpc.SetTimestamp(), rpc.SetSignature(s.secret))

	resp, err := s.httpclient.ExecuteHttpOperation(ctx, req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Something wrong")
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(respBody[:]))
	fmt.Println(req)
	out := &pb.PlaceMultipleOrdersResponse{}

	batchOrders := &[]*pb.NewOrderResponse{}

	err = s.m.Unmarshal(respBody, batchOrders)

	if err != nil {
		return nil, err
	}

	out.BatchOrders = *batchOrders

	return out, nil
}

// // Modify Multiple Orders
// func (s *deliveryTradeService) ModifyMultipleOrders(ctx context.Context, request *pb.ModifyMultipleOrdersRequest) (*pb.ModifyMultipleOrdersResponse, error) {
// 	endpoint := fmt.Sprintf("%s/%s", s.domain, EntryMultipleOrders)
// 	body := []*rpc.HttpParameter{
// 		{Key: "batchOrders", Val: fmt.Sprintf("%v", request.GetBatchOrders())},
// 	}

// 	if request.GetRecvWindow() != 0 {
// 		body = append(body, &rpc.HttpParameter{Key: "recvWindow", Val: fmt.Sprintf("%v", request.GetRecvWindow())})
// 	}

// 	req := s.httpclient.GetHttpRequest(rpc.SetEndpoint(endpoint), rpc.SetMethod("PUT"),
// 		rpc.SetParams(body...), rpc.SetPrivate(), rpc.SetTimestamp(), rpc.SetSignature(s.secret))

// 	resp, err := s.httpclient.ExecuteHttpOperation(ctx, req)

// 	if err != nil {
// 		return nil, err
// 	}

// 	if resp.StatusCode != 200 {
// 		return nil, fmt.Errorf("Something wrong")
// 	}

// 	respBody, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	out := &pb.ModifyMultipleOrdersResponse{}

// 	err = s.m.Unmarshal(respBody, out.BatchOrders)

// 	if err != nil {
// 		return nil, err
// 	}
// 	return out, nil
// }

// Cancel All Open Orders
func (s *deliveryTradeService) CancelAllOpenOrders(ctx context.Context, request *pb.CancelAllOpenOrdersRequest) (*pb.CancelAllOpenOrdersResponse, error) {
	endpoint := fmt.Sprintf("%s/%s", s.domain, delivery.EntryPointAllOpenOrders)
	body := []*rpc.HttpParameter{
		{Key: "symbol", Val: request.GetSymbol()},
	}

	if request.GetRecvWindow() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "recvWindow", Val: fmt.Sprintf("%v", request.GetRecvWindow())})
	}

	req := s.httpclient.GetHttpRequest(rpc.SetEndpoint(endpoint), rpc.SetMethod("DELETE"),
		rpc.SetParams(body...), rpc.SetPrivate(), rpc.SetTimestamp(), rpc.SetSignature(s.secret))

	resp, err := s.httpclient.ExecuteHttpOperation(ctx, req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Something wrong")
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	out := &pb.CancelAllOpenOrdersResponse{}

	err = s.m.Unmarshal(respBody, out)

	if err != nil {
		return nil, err
	}
	return out, nil
}

// Cancel all open orders of the specified symbol at the end of the specified countdown
func (s *deliveryTradeService) AutoCancelAllOpenOrder(ctx context.Context, request *pb.AutoCancelAllOpenOrdersRequest) (*pb.AutoCancelAllOpenOrdersResponse, error) {
	endpoint := fmt.Sprintf("%s/%s", s.domain, delivery.EntryPointCountdownCancelAll)
	body := []*rpc.HttpParameter{
		{Key: "symbol", Val: request.GetSymbol()},
		{Key: "countdownTime", Val: fmt.Sprintf("%v", request.GetCountdownTime())},
	}

	if request.GetRecvWindow() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "recvWindow", Val: fmt.Sprintf("%v", request.GetRecvWindow())})
	}

	req := s.httpclient.GetHttpRequest(rpc.SetEndpoint(endpoint), rpc.SetMethod("POST"),
		rpc.SetParams(body...), rpc.SetPrivate(), rpc.SetTimestamp(), rpc.SetSignature(s.secret))

	resp, err := s.httpclient.ExecuteHttpOperation(ctx, req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Something wrong")
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	out := &pb.AutoCancelAllOpenOrdersResponse{}

	err = s.m.Unmarshal(respBody, out)

	if err != nil {
		return nil, err
	}
	return out, nil
}

// Change user's initial leverage in the specific symbol market.
// For Hedge Mode, LONG and SHORT positions of one symbol use
// the same initial leverage and share a total notional value.
func (s *deliveryTradeService) ChangeInitialLeverage(ctx context.Context, request *pb.ChangeInitialLeverageRequest) (*pb.ChangeInitialLeverageResponse, error) {
	endpoint := fmt.Sprintf("%s/%s", s.domain, delivery.EntryPointLeverage)
	body := []*rpc.HttpParameter{
		{Key: "symbol", Val: request.GetSymbol()},
		{Key: "leverage", Val: fmt.Sprintf("%v", request.GetLeverage())},
	}

	if request.GetRecvWindow() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "recvWindow", Val: fmt.Sprintf("%v", request.GetRecvWindow())})
	}

	req := s.httpclient.GetHttpRequest(rpc.SetEndpoint(endpoint), rpc.SetMethod("POST"),
		rpc.SetParams(body...), rpc.SetPrivate(), rpc.SetTimestamp(), rpc.SetSignature(s.secret))

	resp, err := s.httpclient.ExecuteHttpOperation(ctx, req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Something wrong")
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	out := &pb.ChangeInitialLeverageResponse{}

	err = s.m.Unmarshal(respBody, out)

	if err != nil {
		return nil, err
	}
	return out, nil
}

// Change user's margin type in the specific symbol market.
/// For Hedge Mode, LONG and SHORT positions of
// one symbol use the same margin type.
// With ISOLATED margin type, margins of
// the LONG and SHORT positions are isolated from each other.
func (s *deliveryTradeService) ChangeMarginType(ctx context.Context, request *pb.ChangeMarginTypeRequest) (*pb.ChangeMarginTypeResponse, error) {
	endpoint := fmt.Sprintf("%s/%s", s.domain, delivery.EntryPointMarginType)
	body := []*rpc.HttpParameter{
		{Key: "symbol", Val: request.GetSymbol()},
		{Key: "marginType", Val: request.GetMarginType().String()},
	}

	if request.GetRecvWindow() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "recvWindow", Val: fmt.Sprintf("%v", request.GetRecvWindow())})
	}

	req := s.httpclient.GetHttpRequest(rpc.SetEndpoint(endpoint), rpc.SetMethod("POST"),
		rpc.SetParams(body...), rpc.SetPrivate(), rpc.SetTimestamp(), rpc.SetSignature(s.secret))

	resp, err := s.httpclient.ExecuteHttpOperation(ctx, req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Something wrong")
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	out := &pb.ChangeMarginTypeResponse{}

	err = s.m.Unmarshal(respBody, out)

	if err != nil {
		return nil, err
	}
	return out, nil
}

// Modify Isolated Position Margin
func (s *deliveryTradeService) ModifyIsolatedPositionMargin(ctx context.Context, request *pb.ModifyIsolatedPositionMarginRequest) (*pb.ModifyIsolatedPositionMarginResponse, error) {
	endpoint := fmt.Sprintf("%s/%s", s.domain, delivery.EntryPointPositionMargin)
	body := []*rpc.HttpParameter{
		{Key: "symbol", Val: request.GetSymbol()},
		{Key: "amount", Val: fmt.Sprintf("%v", request.GetAmount())},
		{Key: "type", Val: fmt.Sprintf("%v", request.GetType())},
	}

	if request.GetRecvWindow() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "recvWindow", Val: fmt.Sprintf("%v", request.GetRecvWindow())})
	}

	if request.GetPositionSide() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "positionSide", Val: request.GetPositionSide().String()})
	}

	req := s.httpclient.GetHttpRequest(rpc.SetEndpoint(endpoint), rpc.SetMethod("POST"),
		rpc.SetParams(body...), rpc.SetPrivate(), rpc.SetTimestamp(), rpc.SetSignature(s.secret))

	resp, err := s.httpclient.ExecuteHttpOperation(ctx, req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Something wrong")
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	out := &pb.ModifyIsolatedPositionMarginResponse{}

	err = s.m.Unmarshal(respBody, out)

	if err != nil {
		return nil, err
	}
	return out, nil
}
