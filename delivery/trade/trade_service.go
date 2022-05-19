package trade

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/h9896/bingo-pkg-protobuf/services/delivery/v1"
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
	endpoint := fmt.Sprintf("%s/%s", s.domain, EntryPointPositionMode)
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

// Get user's position mode (Hedge Mode or One-way Mode ) on EVERY symbol
func (s *deliveryTradeService) GetPositionMode(ctx context.Context, request *pb.Empty) (*pb.GetPositionModeResponse, error) {
	endpoint := fmt.Sprintf("%s/%s", s.domain, EntryPointPositionMode)

	req := s.httpclient.GetHttpRequest(rpc.SetEndpoint(endpoint), rpc.SetMethod("get"),
		rpc.SetPrivate(), rpc.SetTimestamp(), rpc.SetSignature(s.secret))

	resp, err := s.httpclient.ExecuteHttpOperation(ctx, req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Something wrong")
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	out := &pb.GetPositionModeResponse{}

	// err = s.m.Unmarshal(respBody, out)
	err = json.Unmarshal(respBody, out)

	if err != nil {
		return nil, err
	}
	return out, nil
}

// Send in a new order.
func (s *deliveryTradeService) NewOrder(ctx context.Context, request *pb.NewOrderRequest) (*pb.NewOrderResponse, error) {
	endpoint := fmt.Sprintf("%s/%s", s.domain, EntryPointOrder)
	body := []*rpc.HttpParameter{
		{Key: "symbol", Val: request.GetSymobl()},
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
	endpoint := fmt.Sprintf("%s/%s", s.domain, EntryPointOrder)
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
