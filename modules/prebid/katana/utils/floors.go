package utils

import (
	"math"

	"github.com/prebid/prebid-server/v2/modules/prebid/katana/models"
	"github.com/prebid/prebid-server/v2/openrtb_ext"
	"github.com/prebid/prebid-server/v2/util/ptrutil"
)

const (
	NotSet = -1

	// USD denotes currency USD
	USD = "USD"

	//floor types
	SoftFloor = 0
	HardFloor = 1
)

var FloorSourceMap = map[string]int{
	openrtb_ext.NoDataLocation:  0,
	openrtb_ext.RequestLocation: 1,
	openrtb_ext.FetchLocation:   2,
}

// FetchStatusMap maps floor fetch status with integer codes
var FetchStatusMap = map[string]int{
	openrtb_ext.FetchNone:       0,
	openrtb_ext.FetchSuccess:    1,
	openrtb_ext.FetchError:      2,
	openrtb_ext.FetchInprogress: 3,
	openrtb_ext.FetchTimeout:    4,
}

// GetBidLevelFloorsDetails return floorvalue and floorrulevalue
func GetBidLevelFloorsDetails(bidExt openrtb_ext.ExtBid, impCtx models.ImpCtx,
	currencyConversion func(from, to string, value float64) (float64, error)) (fv, frv float64) {
	var floorCurrency string
	frv = NotSet

	if bidExt.Prebid != nil && bidExt.Prebid.Floors != nil {
		floorCurrency = bidExt.Prebid.Floors.FloorCurrency
		fv = RoundToTwoDigit(bidExt.Prebid.Floors.FloorValue)
		frv = fv
		if bidExt.Prebid.Floors.FloorRuleValue > 0 {
			frv = RoundToTwoDigit(bidExt.Prebid.Floors.FloorRuleValue)
		}
	}

	// if floor values are not set from bid.ext then fall back to imp.bidfloor
	if frv == NotSet && impCtx.BidFloor != 0 {
		fv = RoundToTwoDigit(impCtx.BidFloor)
		frv = fv
		floorCurrency = impCtx.BidFloorCur
	}

	// convert the floor values in USD currency
	if floorCurrency != "" && floorCurrency != USD {
		value, _ := currencyConversion(floorCurrency, USD, fv)
		fv = RoundToTwoDigit(value)
		value, _ = currencyConversion(floorCurrency, USD, frv)
		frv = RoundToTwoDigit(value)
	}

	if frv == NotSet {
		frv = 0 // set it back to 0
	}

	return
}

func RoundToTwoDigit(value float64) float64 {
	output := math.Pow(10, float64(2))
	return float64(math.Round(value*output)) / output
}

// GetFloorsDetails returns floors details from response.ext.prebid
func GetFloorsDetails(responseExt openrtb_ext.ExtBidResponse) (floorDetails models.FloorsDetails) {
	if responseExt.Prebid != nil && responseExt.Prebid.Floors != nil {
		floors := responseExt.Prebid.Floors
		if floors.Skipped != nil {
			floorDetails.Skipfloors = ptrutil.ToPtr(0)
			if *floors.Skipped {
				floorDetails.Skipfloors = ptrutil.ToPtr(1)
			}
		}
		if floors.Data != nil && len(floors.Data.ModelGroups) > 0 {
			floorDetails.FloorModelVersion = floors.Data.ModelGroups[0].ModelVersion
		}
		if len(floors.PriceFloorLocation) > 0 {
			if source, ok := FloorSourceMap[floors.PriceFloorLocation]; ok {
				floorDetails.FloorSource = &source
			}
		}
		if status, ok := FetchStatusMap[floors.FetchStatus]; ok {
			floorDetails.FloorFetchStatus = &status
		}
		floorDetails.FloorProvider = floors.FloorProvider
		if floors.Data != nil && len(floors.Data.FloorProvider) > 0 {
			floorDetails.FloorProvider = floors.Data.FloorProvider
		}
		if floors.Enforcement != nil && floors.Enforcement.EnforcePBS != nil && *floors.Enforcement.EnforcePBS {
			floorDetails.FloorType = HardFloor
		}
	}
	return floorDetails
}
