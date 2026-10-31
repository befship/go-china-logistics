package yanwen

// Environments.
const (
	EnvSandbox    = "sandbox"
	EnvProduction = "production"
)

// API version fixed by Yanwen open platform.
const APIVersion = "V1.0"

// Method names used by this SDK surface.
const (
	MethodOrderCreate     = "express.order.create"
	MethodOrderCancel     = "express.order.cancel"
	MethodOrderLabelGet   = "express.order.label.get"
	MethodOrderGetList    = "express.order.getlist"
	MethodCalcList        = "calc.list"
	MethodVerifyUSAddress = "common.verify.us.address"
)

// Product types (calc.list productTypeList).
const (
	ProductTypeEconomy     = "1"  // 燕文经济
	ProductTypeRegistered  = "2"  // 燕文挂号
	ProductTypeExpress     = "3"  // 燕文专线
	ProductTypeCommercial  = "17" // 商业快递
	ProductTypeForeignPost = "19" // 外邮产品
	ProductTypeChinaPost   = "26" // 中邮挂号
)

// Product attributes (calc.list productAttributes).
const (
	ProductAttrGeneral   = "1" // 普货
	ProductAttrSpecial   = "2" // 特货
	ProductAttrSensitive = "3" // 敏感货
	ProductAttrPremium   = "4" // 特品
)

// US overseas company / warehouse codes (create companyCode).
const (
	CompanyLAX01 = "LAX01"
	CompanyNJC01 = "NJC01"
	CompanyORD01 = "ORD01"
	CompanyMIA01 = "MIA01"
	CompanyDFW01 = "DFW01"
	CompanySFO01 = "SFO01"
	CompanyATL01 = "ATL01"
)

// Order status codes from express.order.getlist.
const (
	OrderStatusCreated       = 0  // 已制单
	OrderStatusConfirmed     = 1  // 已确认发货
	OrderStatusReceived      = 2  // 已收货
	OrderStatusInTransit     = 3  // 运输中
	OrderStatusDelivered     = 4  // 已妥投
	OrderStatusCancelled     = 5  // 已取消
	OrderStatusHeld          = 6  // 已截留
	OrderStatusDeliveryFail  = 7  // 投递失败
	OrderStatusWarehouseAbn  = 8  // 仓内异常
	OrderStatusWarehouseRtn  = 9  // 仓内退件
	OrderStatusReturnSigned  = 10 // 退件签收
	OrderStatusTransferAbn   = 11 // 转运异常
	OrderStatusOutForDeliv   = 12 // 派送中
	OrderStatusPendingPickup = 13 // 待提取
	OrderStatusTrackEnded    = 15 // 追踪结束
)
