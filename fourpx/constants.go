package fourpx

// Environments.
const (
	EnvTest       = "test"
	EnvProduction = "production"
)

// API version fixed by 4PX open platform.
const APIVersion = "1.0"

// Method names used by this SDK surface.
const (
	MethodOrderCreate             = "ds.xms.order.create"
	MethodOrderGet                = "ds.xms.order.get"
	MethodOrderCancel             = "ds.xms.order.cancel"
	MethodLabelGet                = "ds.xms.label.get"
	MethodEstimatedCostGet        = "ds.xms.estimated_cost.get"
	MethodLogisticsProductGetList = "ds.xms.logistics_product.getlist"
	MethodTrackingGet             = "tr.order.tracking.get"
)

// Get-number modes from ds.xms.order.get (consignment_info.get_no_mode).
const (
	GetNoModeCreate    = "C" // 创建订单时取号
	GetNoModeWarehouse = "U" // 仓库作业时取号
)

// Consignment statuses for ds.xms.order.get request/response.
const (
	ConsignmentDraft       = "D" // 草稿（仅响应）
	ConsignmentPredicted   = "P" // 已预报
	ConsignmentHandedOver  = "V" // 已交接/已交货
	ConsignmentInWarehouse = "H" // 库内作业中/已入库
	ConsignmentOutbound    = "C" // 已出库
	ConsignmentFinished    = "F" // 已完成（仅请求筛选）
	ConsignmentReturned    = "R" // 已退件（仅请求筛选）
	ConsignmentClosed      = "X" // 已关闭
	ConsignmentAll         = "ALL"
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
