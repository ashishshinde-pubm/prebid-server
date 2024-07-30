package models

type RequestCtx struct {
	Id             string
	StartTime      int64
	PublisherIDStr string
	PublisherID    int
	ProfileID      string
	VersionID      string
	PageURL        string
	Platform       string
	SSAI           string
	TestGroup      int
	ATTS           *float64
	UA             string
	DeviceCtx      DeviceCtx
	Source, Origin string
	Country        string
	IP             string
	ImpCtx         map[string]ImpCtx

	//trackers per bid
	Trackers        map[string]TrackerDetails
	TrackerEndpoint string
}

type ImpCtx struct {
	TagId             string
	Secure            *int8
	IsRewardInventory *int8
	BidFloor          float64
	BidFloorCur       string
}
