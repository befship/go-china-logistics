package fourpx

import (
	"encoding/json"
	"strconv"
)

// Person is a sender / recipient address block.
type Person struct {
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	Company     string `json:"company,omitempty"`
	Phone       string `json:"phone,omitempty"`
	Phone2      string `json:"phone2,omitempty"`
	Email       string `json:"email,omitempty"`
	PostCode    string `json:"post_code,omitempty"`
	Country     string `json:"country,omitempty"`
	State       string `json:"state,omitempty"`
	City        string `json:"city,omitempty"`
	District    string `json:"district,omitempty"`
	Street      string `json:"street,omitempty"`
	HouseNumber string `json:"house_number,omitempty"`
}

// DeclareProductInfo is a customs declaration line.
type DeclareProductInfo struct {
	DeclareProductNameCN   string  `json:"declare_product_name_cn,omitempty"`
	DeclareProductNameEN   string  `json:"declare_product_name_en,omitempty"`
	DeclareProductCodeQty  float64 `json:"declare_product_code_qty"`
	DeclareUnitPriceExport float64 `json:"declare_unit_price_export"`
	CurrencyExport         string  `json:"currency_export"`
	DeclareUnitPriceImport float64 `json:"declare_unit_price_import"`
	CurrencyImport         string  `json:"currency_import"`
	DeclareProductCode     string  `json:"declare_product_code,omitempty"`
	BrandExport            string  `json:"brand_export"`
	BrandImport            string  `json:"brand_import"`
	HSCodeExport           string  `json:"hscode_export,omitempty"`
	HSCodeImport           string  `json:"hscode_import,omitempty"`
	UnitNetWeight          float64 `json:"unit_net_weight,omitempty"`
	UnitGrossWeight        float64 `json:"unit_gross_weight,omitempty"`
}

// Parcel is one package in an order.
type Parcel struct {
	Weight             float64              `json:"weight"`
	Length             float64              `json:"length,omitempty"`
	Width              float64              `json:"width,omitempty"`
	Height             float64              `json:"height,omitempty"`
	ParcelValue        float64              `json:"parcel_value"`
	IncludeBattery     string               `json:"include_battery"` // Y/N
	Currency           string               `json:"currency"`
	DeclareProductInfo []DeclareProductInfo `json:"declare_product_info"`
}

// LogisticsServiceInfo selects the logistics product.
type LogisticsServiceInfo struct {
	LogisticsProductCode string `json:"logistics_product_code"`
	CustomsService       string `json:"customs_service,omitempty"`   // Y/N
	SignatureService     string `json:"signature_service,omitempty"` // Y/N
}

// ReturnInfo controls domestic / overseas return behaviour.
type ReturnInfo struct {
	IsReturnOnDomestic string `json:"is_return_on_domestic"`
	IsReturnOnOversea  string `json:"is_return_on_oversea"`
	IsReturnOnOverdue  string `json:"is_return_on_overdue,omitempty"`
}

// DeliverTypeInfo is optional delivery-type metadata.
type DeliverTypeInfo struct {
	DeliverType string `json:"deliver_type,omitempty"`
}

// OrderCreateRequest is the body for ds.xms.order.create.
type OrderCreateRequest struct {
	RefNo                string               `json:"ref_no"`
	BusinessType         string               `json:"business_type"`
	DutyType             string               `json:"duty_type"`
	CargoType            string               `json:"cargo_type,omitempty"`
	LogisticsServiceInfo LogisticsServiceInfo `json:"logistics_service_info"`
	ReturnInfo           ReturnInfo           `json:"return_info"`
	IsInsure             string               `json:"is_insure"`
	ParcelList           []Parcel             `json:"parcel_list"`
	Sender               Person               `json:"sender"`
	RecipientInfo        Person               `json:"recipient_info"`
	DeliverTypeInfo      *DeliverTypeInfo     `json:"deliver_type_info,omitempty"`
	VATNo                string               `json:"vat_no,omitempty"`
	EORINo               string               `json:"eori_no,omitempty"`
	IOSSNo               string               `json:"ioss_no,omitempty"`
}

// OrderCreateData is data from a successful create.
type OrderCreateData struct {
	TrackingNo4PX        string `json:"4px_tracking_no,omitempty"`
	DSConsignmentNo      string `json:"ds_consignment_no,omitempty"`
	RefNo                string `json:"ref_no,omitempty"`
	LabelBarcode         string `json:"label_barcode,omitempty"`
	LabelURL             string `json:"label_url,omitempty"`
	LogisticsChannelNo   string `json:"logistics_channel_no,omitempty"`
	LogisticsChannelName string `json:"logistics_channel_name,omitempty"`
	ODAResultSign        string `json:"oda_result_sign,omitempty"`
}

// OrderGetRequest is the body for ds.xms.order.get.
// Provide request_no for a single lookup, or a create-time window (max 7 days).
type OrderGetRequest struct {
	RequestNo                    string `json:"request_no,omitempty"`
	StartTimeOfCreateConsignment *int64 `json:"start_time_of_create_consignment,omitempty"` // unix ms
	EndTimeOfCreateConsignment   *int64 `json:"end_time_of_create_consignment,omitempty"`   // unix ms
	ConsignmentStatus            string `json:"consignment_status,omitempty"`               // P/V/H/C/F/R/X/ALL
}

// ConsignmentInfo is consignment_info from ds.xms.order.get.
type ConsignmentInfo struct {
	DSConsignmentNo       string `json:"ds_consignment_no,omitempty"`
	TrackingNo4PX         string `json:"4px_tracking_no,omitempty"`
	RefNo                 string `json:"ref_no,omitempty"`
	GetNoMode             string `json:"get_no_mode,omitempty"` // C|U
	LogisticsChannelNo    string `json:"logistics_channel_no,omitempty"`
	GetNoExMsg            string `json:"get_no_exmsg,omitempty"`
	LogisticsProductCode  string `json:"logistics_product_code,omitempty"`
	LogisticsProductName  string `json:"logistics_product_name,omitempty"`
	ConsignmentStatus     string `json:"consignment_status,omitempty"`
	InsureStatus          string `json:"insure_status,omitempty"`
	InsureType            string `json:"insure_type,omitempty"`
	HasCheckODA           string `json:"has_check_oda,omitempty"`
	ODAResultSign         string `json:"oda_result_sign,omitempty"`
	IsHoldSign            string `json:"is_hold_sign,omitempty"`
	ConsignmentCreateDate int64  `json:"consignment_create_date,omitempty"` // unix ms
	InboundDate4PX        int64  `json:"4px_inbound_date,omitempty"`        // unix ms
	OutboundDate4PX       int64  `json:"4px_outbound_date,omitempty"`       // unix ms
}

// ParcelConfirmInfo is parcel_confirm_info from ds.xms.order.get.
type ParcelConfirmInfo struct {
	ConfirmParcelQty          string              `json:"confirm_parcel_qty,omitempty"`
	ConfirmParcelWeight       float64             `json:"confirm_parcel_weight,omitempty"`
	ConfirmParcelVolumeWeight float64             `json:"confirm_parcel_volume_weight,omitempty"`
	ConfirmParcelChargeWeight float64             `json:"confirm_parcel_charge_weight,omitempty"`
	ParcelListConfirmInfo     []ParcelConfirmItem `json:"parcel_list_confirm_info,omitempty"`
}

// ParcelConfirmItem is one row in parcel_list_confirm_info.
type ParcelConfirmItem struct {
	ConfirmWeight           float64 `json:"confirm_weight,omitempty"`
	ConfirmVolumeWeight     float64 `json:"confirm_volume_weight,omitempty"`
	ConfirmLength           float64 `json:"confirm_length,omitempty"` // mm
	ConfirmWidth            float64 `json:"confirm_width,omitempty"`  // mm
	ConfirmHigh             float64 `json:"confirm_high,omitempty"`   // mm
	ConfirmChargeWeight     float64 `json:"confirm_charge_weight,omitempty"`
	ConfirmIncludeBattery   string  `json:"confirm_include_battery,omitempty"`
	ConfirmBatteryType      string  `json:"confirm_battery_type,omitempty"`
	ParcelTotalValueConfirm float64 `json:"parcel_total_value_confirm,omitempty"`
	CurrencyCode            string  `json:"currency_code,omitempty"`
}

// OrderGetItem is one element of the data array from ds.xms.order.get.
type OrderGetItem struct {
	ConsignmentInfo       *ConsignmentInfo    `json:"consignment_info,omitempty"`
	ParcelConfirmInfo     *ParcelConfirmInfo  `json:"parcel_confirm_info,omitempty"`
	ParcelListConfirmInfo []ParcelConfirmItem `json:"parcel_list_confirm_info,omitempty"`
}

// OrderCancelRequest is the body for ds.xms.order.cancel.
type OrderCancelRequest struct {
	RequestNo    string `json:"request_no"`
	CancelReason string `json:"cancel_reason"`
}

// LabelGetRequest is the body for ds.xms.label.get.
type LabelGetRequest struct {
	RequestNo              string `json:"request_no"`
	ResponseLabelFormat    string `json:"response_label_format,omitempty"`
	IsPrintTime            string `json:"is_print_time,omitempty"`
	IsPrintBuyerID         string `json:"is_print_buyer_id,omitempty"`
	IsPrintPickInfo        string `json:"is_print_pick_info,omitempty"`
	IsPrintDeclarationList string `json:"is_print_declaration_list,omitempty"`
	IsPrintCustomerWeight  string `json:"is_print_customer_weight,omitempty"`
	CreatePackageLabel     string `json:"create_package_label,omitempty"`
	IsPrintPickBarcode     string `json:"is_print_pick_barcode,omitempty"`
	IsPrintMerge           string `json:"is_print_merge,omitempty"`
	BarCodeOrderType       string `json:"bar_code_order_type,omitempty"`
}

// LabelURLInfo holds label URLs returned by label.get.
type LabelURLInfo struct {
	LogisticsLabel string `json:"logistics_label,omitempty"`
	CustomLabel    string `json:"custom_label,omitempty"`
	PackageLabel   string `json:"package_label,omitempty"`
	InvoiceLabel   string `json:"invoice_label,omitempty"`
}

// LabelGetData is data from label.get.
type LabelGetData struct {
	LabelBarcode      string          `json:"label_barcode,omitempty"`
	ChildLabelBarcode []string        `json:"child_label_barcode,omitempty"`
	LabelURLInfo      *LabelURLInfo   `json:"label_url_info,omitempty"`
	LabelURL          string          `json:"label_url,omitempty"`
	LabelInfo         string          `json:"label_info,omitempty"`
	LabelInfoList     json.RawMessage `json:"label_info_list,omitempty"`
}

// EstimatedCostRequest is the body for ds.xms.estimated_cost.get.
// Weight / dimensions are strings per official docs (grams / cm).
type EstimatedCostRequest struct {
	RequestNo            string   `json:"request_no,omitempty"`
	CountryCode          string   `json:"country_code,omitempty"`
	Weight               string   `json:"weight,omitempty"`
	Length               string   `json:"length,omitempty"`
	Width                string   `json:"width,omitempty"`
	Height               string   `json:"height,omitempty"`
	CargoCode            string   `json:"cargocode,omitempty"`
	LogisticsProductCode []string `json:"logistics_product_code,omitempty"`
	RecipientPostCode    string   `json:"recipient_post_code,omitempty"`
}

// EstimatedCostItem is one rate quote row.
type EstimatedCostItem struct {
	TrackingNo4PX        string  `json:"4px_tracking_no,omitempty"`
	RefNo                string  `json:"ref_no,omitempty"`
	LogisticsChannelNo   string  `json:"logistics_channel_no,omitempty"`
	LogisticsProductCode string  `json:"logistics_product_code,omitempty"`
	LumpSumFee           float64 `json:"lump_sum_fee,omitempty"` // CNY; API may send number or string
	IsVolumeCargo        string  `json:"is_volume_cargo,omitempty"`
	ChargeWeight         string  `json:"charge_weight,omitempty"`
	EstimatedTime        string  `json:"estimated_time,omitempty"`
	IsShowTrack          string  `json:"is_show_track,omitempty"`
	Remarks              string  `json:"remarks,omitempty"`
}

// UnmarshalJSON accepts lump_sum_fee as either a JSON number or numeric string
// (test env returns strings inside a stringified data payload).
func (e *EstimatedCostItem) UnmarshalJSON(b []byte) error {
	type alias EstimatedCostItem
	aux := struct {
		LumpSumFee json.RawMessage `json:"lump_sum_fee"`
		*alias
	}{alias: (*alias)(e)}
	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}
	if len(aux.LumpSumFee) == 0 || string(aux.LumpSumFee) == "null" {
		return nil
	}
	if aux.LumpSumFee[0] == '"' {
		var s string
		if err := json.Unmarshal(aux.LumpSumFee, &s); err != nil {
			return err
		}
		if s == "" {
			return nil
		}
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		e.LumpSumFee = v
		return nil
	}
	return json.Unmarshal(aux.LumpSumFee, &e.LumpSumFee)
}

// LogisticsProductGetListRequest is the body for ds.xms.logistics_product.getlist.
type LogisticsProductGetListRequest struct {
	TransportMode       string `json:"transport_mode"`
	SourceCountryCode   string `json:"source_country_code,omitempty"`
	SourceWarehouseCode string `json:"source_warehouse_code,omitempty"`
	DestCountryCode     string `json:"dest_country_code,omitempty"`
}

// LogisticsProductItem is one product catalog row.
type LogisticsProductItem struct {
	LogisticsProductCode   string `json:"logistics_product_code,omitempty"`
	LogisticsProductNameCN string `json:"logistics_product_name_cn,omitempty"`
	LogisticsProductNameEN string `json:"logistics_product_name_en,omitempty"`
	TransportMode          string `json:"transport_mode,omitempty"`
	WithBattery            string `json:"with_battery,omitempty"`
	OrderTrack             string `json:"order_track,omitempty"`
	BillingVolumeWeight    string `json:"billing_volume_weight,omitempty"`
	SurfaceMail            string `json:"surface_mail,omitempty"`
}

// TrackingGetRequest is the body for tr.order.tracking.get.
type TrackingGetRequest struct {
	DeliveryOrderNo string `json:"deliveryOrderNo"`
}

// TrackingEvent is one tracking timeline entry.
type TrackingEvent struct {
	BusinessLinkCode string `json:"businessLinkCode,omitempty"`
	OccurDatetime    string `json:"occurDatetime,omitempty"`
	TrackingContent  string `json:"trackingContent,omitempty"`
	TimeZone         string `json:"timeZone,omitempty"`
	OccurLocation    string `json:"occurLocation,omitempty"`
	Country          string `json:"country,omitempty"`
	City             string `json:"city,omitempty"`
}

// TrackingData is tracking payload.
type TrackingData struct {
	DeliveryOrderNo    string          `json:"deliveryOrderNo,omitempty"`
	DestinationCountry string          `json:"destinationCountry,omitempty"`
	ServerName         string          `json:"serverName,omitempty"`
	ServerNum          string          `json:"serverNum,omitempty"`
	TrackingList       []TrackingEvent `json:"trackingList,omitempty"`
	Status             string          `json:"status,omitempty"`
	POD                string          `json:"pod,omitempty"`
	ServerHawbcode     string          `json:"server_hawbcode,omitempty"`
}

// APIErrorItem is one entry in the router errors array.
type APIErrorItem struct {
	ErrorCode     string `json:"error_code,omitempty"`
	ErrorMsg      string `json:"error_msg,omitempty"`
	ReferenceCode string `json:"reference_code,omitempty"`
}

// Response is the 4PX router envelope.
// Success when Result is "1" or "2"; failure when "0".
type Response struct {
	Result string          `json:"result"`
	Msg    string          `json:"msg"`
	Data   json.RawMessage `json:"data"`
	Errors []APIErrorItem  `json:"errors,omitempty"`
}
