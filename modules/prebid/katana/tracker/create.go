package tracker

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/prebid/openrtb/v20/openrtb2"
	"github.com/prebid/prebid-server/v2/modules/prebid/katana/models"
	"github.com/prebid/prebid-server/v2/modules/prebid/katana/utils"
	"github.com/prebid/prebid-server/v2/openrtb_ext"
)

func CreateTracker(rctx models.RequestCtx, bid *openrtb2.Bid, seat string, currency string, responseExt openrtb_ext.ExtBidResponse) models.TrackerDetails {
	floorsDetails := utils.GetFloorsDetails(responseExt)

	tracker := models.Tracker{
		IID:               rctx.Id,
		PubID:             rctx.PublisherID,
		PageURL:           rctx.PageURL,
		Timestamp:         rctx.StartTime,
		ProfileID:         rctx.ProfileID,
		VersionID:         rctx.VersionID,
		Platform:          int(rctx.DeviceCtx.Platform),
		Origin:            rctx.Origin,
		FloorModelVersion: floorsDetails.FloorModelVersion,
		FloorType:         floorsDetails.FloorType,
		FloorSkippedFlag:  floorsDetails.Skipfloors,
		FloorSource:       floorsDetails.FloorSource,
		PartnerInfo: models.Partner{
			PartnerID:  seat,
			BidderCode: seat,
			GrossECPM:  bid.Price,
			NetECPM:    bid.Price,
			BidID:      bid.ID,
			OrigBidID:  bid.ID,
			AdSize:     fmt.Sprintf("%dx%d", bid.W, bid.H),
			ServerSide: 1,
			DealID:     "-1",
		},
	}

	if rctx.DeviceCtx.Ext != nil {
		tracker.ATTS, _ = rctx.DeviceCtx.Ext.GetAtts()
	}

	if impCtx, ok := rctx.ImpCtx[bid.ImpID]; ok {
		tracker.SlotID = impCtx.TagId
		tracker.Adunit = impCtx.TagId

		if impCtx.IsRewardInventory != nil {
			tracker.RewardedInventory = int(*impCtx.IsRewardInventory)
		}

	}

	var bidExt openrtb_ext.ExtBid
	var adformat string
	if err := json.Unmarshal(bid.Ext, &bidExt); err == nil {
		if bidExt.Prebid != nil {
			if bidExt.Prebid.Video != nil && bidExt.Prebid.Video.Duration > 0 {
				tracker.PartnerInfo.AdDuration = bidExt.Prebid.Video.Duration
			}

			if bidExt.Prebid.Meta != nil && len(bidExt.Prebid.Meta.AdapterCode) != 0 && seat != bidExt.Prebid.Meta.AdapterCode {
				tracker.PartnerInfo.PartnerID = bidExt.Prebid.Meta.AdapterCode
			}

			if len(bidExt.Prebid.Type) > 0 {
				adformat = string(bidExt.Prebid.Type)
			}
		}
	}

	tracker.PartnerInfo.Adformat = adformat

	if len(bid.ADomain) != 0 {
		if domain, err := ExtractDomain(bid.ADomain[0]); err == nil {
			tracker.PartnerInfo.Advertiser = domain
		}
	}

	if len(bid.DealID) > 0 {
		tracker.PartnerInfo.DealID = bid.DealID
	}

	return models.TrackerDetails{
		Tracker:       tracker,
		TrackerURL:    constructTrackerURL(rctx.TrackerEndpoint, &tracker),
		BidType:       adformat,
		Price:         bid.Price,
		PriceModel:    models.VideoPricingModelCPM,
		PriceCurrency: currency,
	}
}

func ExtractDomain(rawURL string) (string, error) {
	if !strings.HasPrefix(rawURL, "http") {
		rawURL = "http://" + rawURL
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	return u.Host, nil
}
