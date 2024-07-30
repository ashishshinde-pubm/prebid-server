package katana

import (
	"context"
	"encoding/json"

	"github.com/prebid/prebid-server/v2/hooks/hookstage"
	"github.com/prebid/prebid-server/v2/modules/prebid/katana/models"
	"github.com/prebid/prebid-server/v2/modules/prebid/katana/tracker"
	"github.com/prebid/prebid-server/v2/openrtb_ext"
)

func (k Katana) handleAuctionResponseHook(
	ctx context.Context,
	moduleCtx hookstage.ModuleInvocationContext,
	payload hookstage.AuctionResponsePayload,
) (hookstage.HookResult[hookstage.AuctionResponsePayload], error) {
	result := hookstage.HookResult[hookstage.AuctionResponsePayload]{}
	result.ChangeSet = hookstage.ChangeSet[hookstage.AuctionResponsePayload]{} // Initialize the change set

	if payload.BidResponse == nil || len(payload.BidResponse.SeatBid) == 0 {
		return result, nil
	}

	if len(moduleCtx.ModuleContext) == 0 {
		result.DebugMessages = append(result.DebugMessages, "[katana] error: module-ctx not found in handleAuctionResponseHook()")
		result.Reject = true
		return result, nil
	}

	rctx, ok := moduleCtx.ModuleContext["rctx"].(models.RequestCtx)
	if !ok {
		result.DebugMessages = append(result.DebugMessages, "[katana] error: request-ctx not found in handleAuctionResponseHook()")
		result.Reject = true
		return result, nil
	}

	var responseExt openrtb_ext.ExtBidResponse
	err := json.Unmarshal(payload.BidResponse.Ext, &responseExt)
	if err != nil {
		result.DebugMessages = append(result.DebugMessages, "[katana] error: failed to unmarshal response ext")
		result.Reject = true
		return result, nil
	}

	if rctx.Trackers == nil {
		rctx.Trackers = make(map[string]models.TrackerDetails)
	}

	for _, seatBid := range payload.BidResponse.SeatBid {
		for _, bid := range seatBid.Bid {
			tracker := tracker.CreateTracker(rctx, &bid, seatBid.Seat, payload.BidResponse.Cur, responseExt)
			rctx.Trackers[bid.ID] = tracker
		}
	}

	result.ChangeSet.AddMutation(func(arp hookstage.AuctionResponsePayload) (hookstage.AuctionResponsePayload, error) {
		arp.BidResponse, err = tracker.InjectTrackers(rctx, arp.BidResponse)
		return arp, err
	}, hookstage.MutationUpdate, "Injecting trackers")

	return result, nil
}
