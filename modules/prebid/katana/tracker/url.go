package tracker

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/prebid/prebid-server/v2/modules/prebid/katana/models"
)

func constructTrackerURL(trackerEndpoint string, tracker *models.Tracker) string {
	trackerURL, err := url.Parse(trackerEndpoint)
	if err != nil {
		return ""
	}

	v := url.Values{}
	v.Set(models.TRKPubID, strconv.Itoa(tracker.PubID))
	v.Set(models.TRKPageURL, tracker.PageURL)
	v.Set(models.TRKTimestamp, strconv.FormatInt(tracker.Timestamp, 10))
	v.Set(models.TRKIID, tracker.IID)
	v.Set(models.TRKProfileID, tracker.ProfileID)
	v.Set(models.TRKVersionID, tracker.VersionID)
	v.Set(models.TRKSlotID, tracker.SlotID)
	v.Set(models.TRKAdunit, tracker.Adunit)
	if tracker.RewardedInventory == 1 {
		v.Set(models.TRKRewardedInventory, strconv.Itoa(tracker.RewardedInventory))
	}
	v.Set(models.TRKPlatform, strconv.Itoa(tracker.Platform))
	v.Set(models.TRKPubDomain, tracker.Origin)

	partner := tracker.PartnerInfo
	v.Set(models.TRKPartnerID, partner.PartnerID)
	v.Set(models.TRKBidderCode, partner.BidderCode)
	v.Set(models.TRKGrossECPM, fmt.Sprint(partner.GrossECPM))
	v.Set(models.TRKNetECPM, fmt.Sprint(partner.NetECPM))
	v.Set(models.TRKBidID, partner.BidID)

	v.Set(models.TRKOrigBidID, partner.OrigBidID)
	v.Set(models.TRKAdSize, partner.AdSize)
	if partner.AdDuration > 0 {
		v.Set(models.TRKAdDuration, strconv.Itoa(partner.AdDuration))
	}
	v.Set(models.TRKAdformat, partner.Adformat)
	v.Set(models.TRKServerSide, strconv.Itoa(partner.ServerSide))
	v.Set(models.TRKAdvertiser, partner.Advertiser)

	v.Set(models.TRKFloorType, strconv.Itoa(tracker.FloorType))
	if tracker.FloorSkippedFlag != nil {
		v.Set(models.TRKFloorSkippedFlag, strconv.Itoa(*tracker.FloorSkippedFlag))
	}
	if len(tracker.FloorModelVersion) > 0 {
		v.Set(models.TRKFloorModelVersion, tracker.FloorModelVersion)
	}
	if tracker.FloorSource != nil {
		v.Set(models.TRKFloorSource, strconv.Itoa(*tracker.FloorSource))
	}
	if partner.FloorValue > 0 {
		v.Set(models.TRKFloorValue, fmt.Sprint(partner.FloorValue))
	}
	if partner.FloorRuleValue > 0 {
		v.Set(models.TRKFloorRuleValue, fmt.Sprint(partner.FloorRuleValue))
	}
	v.Set(models.TRKServerLogger, "1")
	v.Set(models.TRKDealID, partner.DealID)

	if tracker.ATTS != nil {
		v.Set(models.TRKATTS, strconv.Itoa(int(*tracker.ATTS)))
	}

	queryString := v.Encode()

	finalTrackerEndpoint := trackerURL.String() + models.TRKQMARK + queryString

	fmt.Println(finalTrackerEndpoint)
	return finalTrackerEndpoint
}
