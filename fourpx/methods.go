package fourpx

import "context"

// CreateOrder creates a direct-shipping consignment (ds.xms.order.create).
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

// CancelOrder cancels a consignment (ds.xms.order.cancel).
func (c *Client) CancelOrder(ctx context.Context, req *OrderCancelRequest) error {
	_, err := c.Call(ctx, MethodOrderCancel, req)
	return err
}

// GetLabel fetches a shipping label (ds.xms.label.get).
func (c *Client) GetLabel(ctx context.Context, req *LabelGetRequest) (*LabelGetData, error) {
	env, err := c.Call(ctx, MethodLabelGet, req)
	if err != nil {
		return nil, err
	}
	var data LabelGetData
	if err := decodeData(env, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// EstimatedCost runs freight estimation (ds.xms.estimated_cost.get).
func (c *Client) EstimatedCost(ctx context.Context, req *EstimatedCostRequest) ([]EstimatedCostItem, error) {
	env, err := c.Call(ctx, MethodEstimatedCostGet, req)
	if err != nil {
		return nil, err
	}
	var items []EstimatedCostItem
	if err := decodeData(env, &items); err != nil {
		return nil, err
	}
	return items, nil
}

// ListLogisticsProducts lists logistics products (ds.xms.logistics_product.getlist).
func (c *Client) ListLogisticsProducts(ctx context.Context, req *LogisticsProductGetListRequest) ([]LogisticsProductItem, error) {
	env, err := c.Call(ctx, MethodLogisticsProductGetList, req)
	if err != nil {
		return nil, err
	}
	var items []LogisticsProductItem
	if err := decodeData(env, &items); err != nil {
		return nil, err
	}
	return items, nil
}

// Track queries tracking by delivery order number (tr.order.tracking.get).
func (c *Client) Track(ctx context.Context, deliveryOrderNo string) (*TrackingData, error) {
	env, err := c.Call(ctx, MethodTrackingGet, &TrackingGetRequest{DeliveryOrderNo: deliveryOrderNo})
	if err != nil {
		return nil, err
	}
	var data TrackingData
	if err := decodeData(env, &data); err != nil {
		return nil, err
	}
	return &data, nil
}
