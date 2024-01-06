package exchange

import (
	"github.com/prebid/openrtb/v19/openrtb2"
	"github.com/prebid/prebid-server/v2/openrtb_ext"
)

// AuctionResponse contains OpenRTB Bid Response object and its extension (un-marshalled) object
type AuctionResponse struct {
	*openrtb2.BidResponse
	ExtBidResponse *openrtb_ext.ExtBidResponse
	SeatNonBid     openrtb_ext.NonBidsWrapper
}
