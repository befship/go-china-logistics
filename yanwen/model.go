package yanwen

import "encoding/json"

// AddressInfo is a sender / receiver address block.
type AddressInfo struct {
	Name         string `json:"name,omitempty"`
	Phone        string `json:"phone,omitempty"`
	Email        string `json:"email,omitempty"`
	Company      string `json:"company,omitempty"`
	Country      string `json:"country,omitempty"`
	State        string `json:"state,omitempty"`
	City         string `json:"city,omitempty"`
	ZipCode      string `json:"zipCode,omitempty"`
	HouseNumber  string `json:"houseNumber,omitempty"`
	Address      string `json:"address,omitempty"`
	TaxNumber    string `json:"taxNumber,omitempty"`
	FINFLNumber  string `json:"finflNumber,omitempty"`
	ShortAddress string `json:"shortAddress,omitempty"`
}

// ProductItem is a declaration line inside a parcel.
type ProductItem struct {
	GoodsNameCh string  `json:"goodsNameCh"`
	GoodsNameEn string  `json:"goodsNameEn"`
	Price       float64 `json:"price"`
	PriceExport float64 `json:"priceExport,omitempty"`
	HSCode      string  `json:"hscode,omitempty"`
	URL         string  `json:"url,omitempty"`
	Material    string  `json:"material,omitempty"`
	Quantity    int     `json:"quantity"`
	Weight      int     `json:"weight"` // grams per unit
	SKU         string  `json:"sku,omitempty"`
	IMEI        string  `json:"imei,omitempty"`
}

// ParcelInfo is parcel / declaration block for create.
type ParcelInfo struct {
	HasBattery      int           `json:"hasBattery"` // 1 yes / 0 no
	Currency        string        `json:"currency"`
	TotalQuantity   int           `json:"totalQuantity"`
	TotalWeight     int           `json:"totalWeight"` // grams
	TotalPrice      float64       `json:"totalPrice,omitempty"`
	Height          int           `json:"height,omitempty"` // cm
	Width           int           `json:"width,omitempty"`
	Length          int           `json:"length,omitempty"`
	IOSS            string        `json:"ioss,omitempty"`
	FreightCharges  float64       `json:"freightCharges,omitempty"`
	FreightCurrency string        `json:"freightCurrency,omitempty"`
	ProductList     []ProductItem `json:"productList,omitempty"`
}

// PopStation is an optional pickup-point block.
type PopStation struct {
	PointID string `json:"pointId,omitempty"`
}

// OrderCreateRequest is the body for express.order.create.
type OrderCreateRequest struct {
	ChannelID         string       `json:"channelId"`
	OrderSource       string       `json:"orderSource,omitempty"`
	OrderNumber       string       `json:"orderNumber"`
	TransactionNumber string       `json:"transactionNumber,omitempty"`
	DateOfReceipt     string       `json:"dateOfReceipt,omitempty"`
	CompanyCode       string       `json:"companyCode,omitempty"`
	Remark            string       `json:"remark,omitempty"`
	SalesPlatform     string       `json:"salesPlatform,omitempty"`
	ReceiverInfo      AddressInfo  `json:"receiverInfo"`
	SenderInfo        *AddressInfo `json:"senderInfo,omitempty"`
	ParcelInfo        ParcelInfo   `json:"parcelInfo"`
	PoPStation        *PopStation  `json:"poPStation,omitempty"`
}

// OrderCreateData is data from a successful create.
type OrderCreateData struct {
	WaybillNumber     string `json:"waybillNumber,omitempty"`
	OrderNumber       string `json:"orderNumber,omitempty"`
	YanwenOrderNumber string `json:"yanwenOrderNumber,omitempty"`
}

// OrderCancelRequest is the body for express.order.cancel.
type OrderCancelRequest struct {
	WaybillNumber string `json:"waybillNumber"`
	Note          string `json:"note"`
}

// LabelGetRequest is the body for express.order.label.get.
type LabelGetRequest struct {
	WaybillNumber string `json:"waybillNumber"`
	PrintRemark   int    `json:"printRemark,omitempty"`
}

// LabelGetData is data from label.get.
type LabelGetData struct {
	WaybillNumber string `json:"waybillNumber,omitempty"`
	IsSuccess     bool   `json:"isSuccess,omitempty"`
	ErrorMsg      string `json:"errorMsg,omitempty"`
	Base64String  string `json:"base64String,omitempty"`
}

// OrderGetListRequest is the body for express.order.getlist.
type OrderGetListRequest struct {
	ListNumber []string `json:"listNumber"`
}

// OrderDetail is one express.order.getlist data item.
type OrderDetail struct {
	WaybillNumber   string          `json:"waybillNumber,omitempty"`
	OrderNumber     string          `json:"orderNumber,omitempty"`
	ReferenceNumber string          `json:"referenceNumber,omitempty"`
	YanwenNumber    string          `json:"yanwenNumber,omitempty"`
	ChannelID       string          `json:"channelId,omitempty"`
	ChannelName     string          `json:"channelName,omitempty"`
	OrderSource     string          `json:"orderSource,omitempty"`
	UserID          string          `json:"userId,omitempty"`
	DateOfReceipt   string          `json:"dateOfReceipt,omitempty"`
	CreateTime      string          `json:"createTime,omitempty"`
	Status          int             `json:"status,omitempty"`
	CompanyCode     string          `json:"companyCode,omitempty"`
	IsPrint         int             `json:"isPrint,omitempty"`
	Remark          string          `json:"remark,omitempty"`
	ReceiverInfo    json.RawMessage `json:"receiverInfo,omitempty"`
	SenderInfo      json.RawMessage `json:"senderInfo,omitempty"`
	ParcelInfo      json.RawMessage `json:"parcelInfo,omitempty"`
}

// CalcListRequest is the body for calc.list.
type CalcListRequest struct {
	CityID            string   `json:"cityId"`
	CountryID         string   `json:"countryId"`
	ProductAttributes string   `json:"productAttributes,omitempty"`
	ProductTypeList   []string `json:"productTypeList,omitempty"`
	Weight            float64  `json:"weight"`
	High              float64  `json:"high,omitempty"`
	Length            float64  `json:"length,omitempty"`
	Width             float64  `json:"width,omitempty"`
	PostCode          string   `json:"postCode,omitempty"`
}

// CalcListItem is one rate quote row.
type CalcListItem struct {
	BasicFreight      float64 `json:"basicFreight,omitempty"`
	TrunkFreight      float64 `json:"trunkFreight,omitempty"`
	FuelCost          float64 `json:"fuelCost,omitempty"`
	AdditionalCharge  float64 `json:"additionalCharge,omitempty"`
	ExceedLengthMoney float64 `json:"exceedLengthMoney,omitempty"`
	RMBTotalMoney     float64 `json:"rmbTotalMoney,omitempty"`
	TotalMoney        float64 `json:"totalMoney,omitempty"`
	Currency          string  `json:"currency,omitempty"`
	ProductNumber     string  `json:"productNumber,omitempty"` // = channelId / service code
	ProductName       string  `json:"productName,omitempty"`
	ReferAging        string  `json:"referAging,omitempty"`
	Zone              string  `json:"zone,omitempty"`
}

// USAddressFields is the nested address for US address verification.
type USAddressFields struct {
	Address string `json:"address"`
	ZipCode string `json:"zipCode"`
	City    string `json:"city"`
	State   string `json:"state"`
}

// VerifyUSAddressRequest is the body for common.verify.us.address.
type VerifyUSAddressRequest struct {
	ReceiverInfo USAddressFields `json:"receiverInfo"`
}

// VerifyUSAddressResult is the corrected address payload.
type VerifyUSAddressResult struct {
	Address  string `json:"address,omitempty"`
	City     string `json:"city,omitempty"`
	State    string `json:"state,omitempty"`
	ZipCode4 string `json:"zipCode4,omitempty"`
	ZipCode5 string `json:"zipCode5,omitempty"`
}

// VerifyUSAddressData wraps receiverInfo from verify response.
type VerifyUSAddressData struct {
	ReceiverInfo VerifyUSAddressResult `json:"receiverInfo"`
}

// Response is the Yanwen order-API envelope.
type Response struct {
	Success bool            `json:"success"`
	Code    string          `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

// TrackingResponse is the public tracking API payload.
type TrackingResponse struct {
	Data []TrackingShipment `json:"data"`
}

// TrackingShipment is one tracking result.
type TrackingShipment struct {
	Status         string          `json:"status,omitempty"`
	POD            string          `json:"pod,omitempty"`
	ExchangeNumber string          `json:"exchange_number,omitempty"`
	Extra          json.RawMessage `json:"-"`
}
