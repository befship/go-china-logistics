package yanwen

import "context"

// CreateOrder creates a waybill (express.order.create).
func (c *Client) CreateOrder(ctx context.Context, req *OrderCreateRequest) (*OrderCreateData, error) {
	env, err := c.Call(ctx, MethodOrderCreate, req)
	if err != nil {
		return nil, err
	}
	var data OrderCreateData
	if err := decodeData(env, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// CancelOrder cancels a waybill (express.order.cancel).
func (c *Client) CancelOrder(ctx context.Context, req *OrderCancelRequest) error {
	_, err := c.Call(ctx, MethodOrderCancel, req)
	return err
}

// GetLabel fetches a shipping label PDF as base64 (express.order.label.get).
func (c *Client) GetLabel(ctx context.Context, req *LabelGetRequest) (*LabelGetData, error) {
	env, err := c.Call(ctx, MethodOrderLabelGet, req)
	if err != nil {
		return nil, err
	}
	var data LabelGetData
	if err := decodeData(env, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// GetOrders queries order details by waybill / order numbers (express.order.getlist).
func (c *Client) GetOrders(ctx context.Context, numbers []string) ([]OrderDetail, error) {
	env, err := c.Call(ctx, MethodOrderGetList, &OrderGetListRequest{ListNumber: numbers})
	if err != nil {
		return nil, err
	}
	var items []OrderDetail
	if err := decodeData(env, &items); err != nil {
		return nil, err
	}
	return items, nil
}

// CalcList runs freight calculation (calc.list). Requires sales activation.
func (c *Client) CalcList(ctx context.Context, req *CalcListRequest) ([]CalcListItem, error) {
	env, err := c.Call(ctx, MethodCalcList, req)
	if err != nil {
		return nil, err
	}
	var items []CalcListItem
	if err := decodeData(env, &items); err != nil {
		return nil, err
	}
	return items, nil
}

// VerifyUSAddress validates a US destination address (common.verify.us.address).
func (c *Client) VerifyUSAddress(ctx context.Context, req *VerifyUSAddressRequest) (*VerifyUSAddressData, error) {
	env, err := c.Call(ctx, MethodVerifyUSAddress, req)
	if err != nil {
		return nil, err
	}
	var data VerifyUSAddressData
	if err := decodeData(env, &data); err != nil {
		return nil, err
	}
	return &data, nil
}
