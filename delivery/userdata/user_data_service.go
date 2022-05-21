package userdata

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

type deliveryUserDataService struct {
	httpclient rpc.GenericHttpClient
	domain     string
	secret     string
	m          *runtime.JSONPb
}

func NewDeliveryUserDataService(domain, apikey, secret string, useSSL bool, client rpc.HTTPClient) pb.DeliveryUserDataServiceServer {
	service := &deliveryUserDataService{
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

// Get user's position mode (Hedge Mode or One-way Mode ) on EVERY symbol
func (s *deliveryUserDataService) GetPositionMode(ctx context.Context, request *pb.Empty) (*pb.GetPositionModeResponse, error) {
	endpoint := fmt.Sprintf("%s/%s", s.domain, delivery.EntryPointPositionMode)

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

	err = s.m.Unmarshal(respBody, out)

	if err != nil {
		return nil, err
	}
	return out, nil
}

// Get order modification history
func (s *deliveryUserDataService) GetOrderModifyHistory(ctx context.Context, request *pb.GetOrderModifyHistoryRequest) (*pb.GetOrderModifyHistoryResponse, error) {
	endpoint := fmt.Sprintf("%s/%s", s.domain, delivery.EntryPointOrderAmendment)

	body := []*rpc.HttpParameter{
		{Key: "symbol", Val: request.GetSymbol()},
	}

	if request.GetOrderId() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "orderId", Val: fmt.Sprintf("%v", request.GetOrderId())})
	}

	if request.GetOrigClientOrderId() != "" {
		body = append(body, &rpc.HttpParameter{Key: "origClientOrderId", Val: request.GetOrigClientOrderId()})
	}

	if request.GetStartTime() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "startTime", Val: fmt.Sprintf("%v", request.GetStartTime())})
	}

	if request.GetEndTime() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "endTime", Val: fmt.Sprintf("%v", request.GetEndTime())})
	}

	if request.GetLimit() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "limit", Val: fmt.Sprintf("%v", request.GetLimit())})
	}

	if request.GetRecvWindow() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "recvWindow", Val: fmt.Sprintf("%v", request.GetRecvWindow())})
	}

	req := s.httpclient.GetHttpRequest(rpc.SetEndpoint(endpoint), rpc.SetMethod("get"),
		rpc.SetParams(body...), rpc.SetPrivate(),
		rpc.SetTimestamp(), rpc.SetSignature(s.secret))

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

	out := &pb.GetOrderModifyHistoryResponse{}

	listHistory := &[]*pb.OrderModifyHistory{}

	err = s.m.Unmarshal(respBody, listHistory)

	if err != nil {
		return nil, err
	}

	out.OrderModifyHistory = *listHistory

	return out, nil
}

// Check an order's status
func (s *deliveryUserDataService) QueryOrder(ctx context.Context, request *pb.QueryOrderRequest) (*pb.QueryOrderResponse, error) {
	endpoint := fmt.Sprintf("%s/%s", s.domain, delivery.EntryPointOrder)

	body := []*rpc.HttpParameter{
		{Key: "symbol", Val: request.GetSymbol()},
	}

	if request.GetOrderId() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "orderId", Val: fmt.Sprintf("%v", request.GetOrderId())})
	}

	if request.GetOrigClientOrderId() != "" {
		body = append(body, &rpc.HttpParameter{Key: "origClientOrderId", Val: request.GetOrigClientOrderId()})
	}

	if request.GetRecvWindow() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "recvWindow", Val: fmt.Sprintf("%v", request.GetRecvWindow())})
	}

	req := s.httpclient.GetHttpRequest(rpc.SetEndpoint(endpoint), rpc.SetMethod("get"),
		rpc.SetParams(body...), rpc.SetPrivate(),
		rpc.SetTimestamp(), rpc.SetSignature(s.secret))

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

	out := &pb.QueryOrderResponse{}

	err = s.m.Unmarshal(respBody, out)

	if err != nil {
		return nil, err
	}

	return out, nil
}

// Query Current Open Order
func (s *deliveryUserDataService) QueryCurrentOpenOrder(ctx context.Context, request *pb.QueryCurrentOpenOrderRequest) (*pb.QueryCurrentOpenOrderResponse, error) {
	endpoint := fmt.Sprintf("%s/%s", s.domain, delivery.EntryPointOpenOrder)

	body := []*rpc.HttpParameter{
		{Key: "symbol", Val: request.GetSymbol()},
	}

	if request.GetOrderId() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "orderId", Val: fmt.Sprintf("%v", request.GetOrderId())})
	}

	if request.GetOrigClientOrderId() != "" {
		body = append(body, &rpc.HttpParameter{Key: "origClientOrderId", Val: request.GetOrigClientOrderId()})
	}

	if request.GetRecvWindow() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "recvWindow", Val: fmt.Sprintf("%v", request.GetRecvWindow())})
	}

	req := s.httpclient.GetHttpRequest(rpc.SetEndpoint(endpoint), rpc.SetMethod("get"),
		rpc.SetParams(body...), rpc.SetPrivate(),
		rpc.SetTimestamp(), rpc.SetSignature(s.secret))

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

	out := &pb.QueryCurrentOpenOrderResponse{}

	err = s.m.Unmarshal(respBody, out)

	if err != nil {
		return nil, err
	}

	return out, nil
}

// Get all open orders on a symbol.
// Careful when accessing this with no symbol.
func (s *deliveryUserDataService) CurrentAllOpenOrders(ctx context.Context, request *pb.CurrentAllOpenOrdersRequest) (*pb.CurrentAllOpenOrdersResponse, error) {
	endpoint := fmt.Sprintf("%s/%s", s.domain, delivery.EntryPointOpenOrders)

	body := []*rpc.HttpParameter{}

	if request.GetSymbol() != "" {
		body = append(body, &rpc.HttpParameter{Key: "symbol", Val: request.GetSymbol()})
	}

	if request.GetPair() != "" {
		body = append(body, &rpc.HttpParameter{Key: "pair", Val: request.GetPair()})
	}

	if request.GetRecvWindow() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "recvWindow", Val: fmt.Sprintf("%v", request.GetRecvWindow())})
	}

	req := s.httpclient.GetHttpRequest(rpc.SetEndpoint(endpoint), rpc.SetMethod("get"),
		rpc.SetParams(body...), rpc.SetPrivate(),
		rpc.SetTimestamp(), rpc.SetSignature(s.secret))

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

	out := &pb.CurrentAllOpenOrdersResponse{}

	openOrders := &[]*pb.QueryCurrentOpenOrderResponse{}

	err = s.m.Unmarshal(respBody, openOrders)

	if err != nil {
		return nil, err
	}

	out.CurrentAllOpenOrders = *openOrders

	return out, nil
}

// Get all account orders; active, canceled, or filled.{
func (s *deliveryUserDataService) AllOrders(ctx context.Context, request *pb.AllOrdersRequest) (*pb.AllOrdersResponse, error) {
	endpoint := fmt.Sprintf("%s/%s", s.domain, delivery.EntryPointAllOrders)

	body := []*rpc.HttpParameter{}

	if request.GetSymbol() != "" {
		body = append(body, &rpc.HttpParameter{Key: "symbol", Val: request.GetSymbol()})
	}

	if request.GetPair() != "" {
		body = append(body, &rpc.HttpParameter{Key: "pair", Val: request.GetPair()})
	}

	if request.GetOrderId() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "orderId", Val: fmt.Sprintf("%v", request.GetOrderId())})
	}

	if request.GetStartTime() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "startTime", Val: fmt.Sprintf("%v", request.GetStartTime())})
	}

	if request.GetEndTime() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "endTime", Val: fmt.Sprintf("%v", request.GetEndTime())})
	}

	if request.GetRecvWindow() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "recvWindow", Val: fmt.Sprintf("%v", request.GetRecvWindow())})
	}

	req := s.httpclient.GetHttpRequest(rpc.SetEndpoint(endpoint), rpc.SetMethod("get"),
		rpc.SetParams(body...), rpc.SetPrivate(),
		rpc.SetTimestamp(), rpc.SetSignature(s.secret))

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

	out := &pb.AllOrdersResponse{}

	allOrders := &[]*pb.QueryCurrentOpenOrderResponse{}

	err = s.m.Unmarshal(respBody, allOrders)

	if err != nil {
		return nil, err
	}

	out.AllOrders = *allOrders

	return out, nil
}

// Futures Account Balance
func (s *deliveryUserDataService) FuturesAccountBalance(ctx context.Context, request *pb.Empty) (*pb.FuturesAccountBalanceResponse, error) {
	endpoint := fmt.Sprintf("%s/%s", s.domain, delivery.EntryPointBalance)

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

	out := &pb.FuturesAccountBalanceResponse{}

	err = s.m.Unmarshal(respBody, out)

	if err != nil {
		return nil, err
	}

	return out, nil
}

// Position Information
func (s *deliveryUserDataService) PositionInformation(ctx context.Context, request *pb.PositionInformationRequest) (*pb.PositionInformationResponse, error) {
	endpoint := fmt.Sprintf("%s/%s", s.domain, delivery.EntryPointPositionRisk)

	body := []*rpc.HttpParameter{}

	if request.GetMarginAsset() != "" {
		body = append(body, &rpc.HttpParameter{Key: "marginAsset", Val: request.GetMarginAsset()})
	}

	if request.GetPair() != "" {
		body = append(body, &rpc.HttpParameter{Key: "pair", Val: request.GetPair()})
	}

	if request.GetRecvWindow() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "recvWindow", Val: fmt.Sprintf("%v", request.GetRecvWindow())})
	}

	req := s.httpclient.GetHttpRequest(rpc.SetEndpoint(endpoint), rpc.SetMethod("get"),
		rpc.SetParams(body...), rpc.SetPrivate(),
		rpc.SetTimestamp(), rpc.SetSignature(s.secret))

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

	out := &pb.PositionInformationResponse{}

	positions := &[]*pb.Position{}

	err = s.m.Unmarshal(respBody, positions)

	if err != nil {
		return nil, err
	}

	out.Positions = *positions

	return out, nil
}

// Get Position Margin Change History
func (s *deliveryUserDataService) GetPositionMarginChangeHistory(ctx context.Context, request *pb.GetPositionMarginChangeHistoryRequest) (*pb.GetPositionMarginChangeHistoryResponse, error) {
	endpoint := fmt.Sprintf("%s/%s/%s", s.domain, delivery.EntryPointPositionMargin, delivery.History)
	body := []*rpc.HttpParameter{
		{Key: "symbol", Val: request.GetSymbol()},
	}

	if request.GetRecvWindow() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "recvWindow", Val: fmt.Sprintf("%v", request.GetRecvWindow())})
	}

	if request.GetType() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "type", Val: fmt.Sprintf("%v", request.GetType())})
	}

	if request.GetStartTime() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "startTime", Val: fmt.Sprintf("%v", request.GetStartTime())})
	}

	if request.GetEndTime() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "endTime", Val: fmt.Sprintf("%v", request.GetEndTime())})
	}

	if request.GetLimit() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "limit", Val: fmt.Sprintf("%v", request.GetLimit())})
	}

	req := s.httpclient.GetHttpRequest(rpc.SetEndpoint(endpoint), rpc.SetMethod("GET"),
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

	out := &pb.GetPositionMarginChangeHistoryResponse{}

	listPositionMargin := &[]*pb.PositionMargin{}

	err = s.m.Unmarshal(respBody, listPositionMargin)

	out.PositionMargins = *listPositionMargin

	if err != nil {
		return nil, err
	}
	return out, nil
}

// Get current account information.
func (s *deliveryUserDataService) AccountInformation(ctx context.Context, request *pb.Empty) (*pb.AccountInformationResponse, error) {
	endpoint := fmt.Sprintf("%s/%s", s.domain, delivery.EntryPointAccount)

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

	out := &pb.AccountInformationResponse{}

	err = s.m.Unmarshal(respBody, out)

	if err != nil {
		return nil, err
	}

	return out, nil
}

// Get trades for a specific account and symbol.
func (s *deliveryUserDataService) AccountTradeList(ctx context.Context, request *pb.AccountTradeListRequest) (*pb.AccountTradeListResponse, error) {
	endpoint := fmt.Sprintf("%s/%s", s.domain, delivery.EntryPointUserTrades)
	body := []*rpc.HttpParameter{}

	if request.GetRecvWindow() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "recvWindow", Val: fmt.Sprintf("%v", request.GetRecvWindow())})
	}

	if request.GetSymbol() != "" {
		body = append(body, &rpc.HttpParameter{Key: "symbol", Val: request.GetSymbol()})
	}

	if request.GetPair() != "" {
		body = append(body, &rpc.HttpParameter{Key: "pair", Val: request.GetPair()})
	}

	if request.GetStartTime() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "startTime", Val: fmt.Sprintf("%v", request.GetStartTime())})
	}

	if request.GetEndTime() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "endTime", Val: fmt.Sprintf("%v", request.GetEndTime())})
	}

	if request.GetFromId() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "fromId", Val: fmt.Sprintf("%v", request.GetFromId())})
	}

	if request.GetLimit() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "limit", Val: fmt.Sprintf("%v", request.GetLimit())})
	}

	req := s.httpclient.GetHttpRequest(rpc.SetEndpoint(endpoint), rpc.SetMethod("GET"),
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

	out := &pb.AccountTradeListResponse{}

	listAccountTrade := &[]*pb.AccountTrade{}

	err = s.m.Unmarshal(respBody, listAccountTrade)

	out.AccountTrades = *listAccountTrade

	if err != nil {
		return nil, err
	}
	return out, nil
}

// Get Income History
func (s *deliveryUserDataService) GetIncomeHistory(ctx context.Context, request *pb.GetIncomeHistoryRequest) (*pb.GetIncomeHistoryResponse, error) {
	endpoint := fmt.Sprintf("%s/%s", s.domain, delivery.EntryPointIncome)
	body := []*rpc.HttpParameter{}

	if request.GetRecvWindow() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "recvWindow", Val: fmt.Sprintf("%v", request.GetRecvWindow())})
	}

	if request.GetSymbol() != "" {
		body = append(body, &rpc.HttpParameter{Key: "symbol", Val: request.GetSymbol()})
	}

	if request.GetIncomeType() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "incomeType", Val: fmt.Sprintf("%v", request.GetIncomeType())})
	}

	if request.GetStartTime() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "startTime", Val: fmt.Sprintf("%v", request.GetStartTime())})
	}

	if request.GetEndTime() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "endTime", Val: fmt.Sprintf("%v", request.GetEndTime())})
	}

	if request.GetLimit() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "limit", Val: fmt.Sprintf("%v", request.GetLimit())})
	}

	req := s.httpclient.GetHttpRequest(rpc.SetEndpoint(endpoint), rpc.SetMethod("GET"),
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

	out := &pb.GetIncomeHistoryResponse{}

	incomes := &[]*pb.IncomeHistory{}

	err = s.m.Unmarshal(respBody, incomes)

	out.Incomes = *incomes

	if err != nil {
		return nil, err
	}
	return out, nil
}

// Get the pair's default notional bracket list.
func (s *deliveryUserDataService) NotionalBracketForSymbol(ctx context.Context, request *pb.NotionalBracketForSymbolRequest) (*pb.NotionalBracketForSymbolResponse, error) {
	endpoint := fmt.Sprintf("%s/%s", s.domain, delivery.EntryPointLeverageBracket)
	body := []*rpc.HttpParameter{}

	if request.GetRecvWindow() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "recvWindow", Val: fmt.Sprintf("%v", request.GetRecvWindow())})
	}

	if request.GetPair() != "" {
		body = append(body, &rpc.HttpParameter{Key: "pair", Val: request.GetPair()})
	}

	req := s.httpclient.GetHttpRequest(rpc.SetEndpoint(endpoint), rpc.SetMethod("GET"),
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

	out := &pb.NotionalBracketForSymbolResponse{}

	brackets := &[]*pb.NotionalBracketForSymbol{}

	err = s.m.Unmarshal(respBody, brackets)

	out.Brackets = *brackets

	if err != nil {
		return nil, err
	}
	return out, nil
}

// Get the symbol's notional bracket list.
func (s *deliveryUserDataService) NotionalBracketForPair(ctx context.Context, request *pb.NotionalBracketForPairRequest) (*pb.NotionalBracketForPairResponse, error) {
	endpoint := fmt.Sprintf("%s/%s", s.domain, delivery.EntryPointLeverageBracketV2)
	body := []*rpc.HttpParameter{}

	if request.GetRecvWindow() != 0 {
		body = append(body, &rpc.HttpParameter{Key: "recvWindow", Val: fmt.Sprintf("%v", request.GetRecvWindow())})
	}

	if request.GetSymbol() != "" {
		body = append(body, &rpc.HttpParameter{Key: "symbol", Val: request.GetSymbol()})
	}

	req := s.httpclient.GetHttpRequest(rpc.SetEndpoint(endpoint), rpc.SetMethod("GET"),
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

	out := &pb.NotionalBracketForPairResponse{}

	brackets := &[]*pb.NotionalBracketForPair{}

	err = s.m.Unmarshal(respBody, brackets)

	out.Brackets = *brackets

	if err != nil {
		return nil, err
	}
	return out, nil
}
