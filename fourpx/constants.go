package fourpx

// Environments.
const (
	EnvSandbox    = "sandbox"
	EnvTest       = "test"
	EnvProduction = "production"
)

// API version fixed by 4PX open platform.
const APIVersion = "1.0"

// Method names used by this SDK surface.
const (
	MethodOrderCreate             = "ds.xms.order.create"
	MethodOrderCancel             = "ds.xms.order.cancel"
	MethodLabelGet                = "ds.xms.label.get"
	MethodEstimatedCostGet        = "ds.xms.estimated_cost.get"
	MethodLogisticsProductGetList = "ds.xms.logistics_product.getlist"
	MethodTrackingGet             = "tr.order.tracking.get"
)

// Duty types (税费承担方式).
const (
	DutyDDU         = "U" // 收件人付关税
	DutyDDP         = "P" // 寄件方付关税
	DutyDDPPlatform = "F" // 平台代收代缴 VAT
)

// Cargo types.
const (
	CargoGift     = "1"
	CargoDocument = "2"
	CargoSample   = "3"
	CargoOther    = "5"
)

// Yes/No flags used across 4PX payloads.
const (
	Yes = "Y"
	No  = "N"
)

// Label formats.
const (
	LabelPDF = "PDF"
	LabelIMG = "IMG"
)

// Transport modes for logistics product list.
const (
	TransportAll           = "1"
	TransportIntlExpress   = "2"
	TransportIntlParcel    = "3"
	TransportDedicatedLine = "4"
	TransportLianYouTong   = "5"
	TransportOther         = "6"
)

// Cargo codes for estimated cost.
const (
	CargoCodeParcel   = "P"
	CargoCodeDocument = "D"
)
