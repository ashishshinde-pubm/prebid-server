package katana

import (
	"context"
	"net/url"
	"strconv"

	"github.com/prebid/prebid-server/v2/hooks/hookstage"
	"github.com/prebid/prebid-server/v2/modules/prebid/katana/models"
	"github.com/prebid/prebid-server/v2/modules/prebid/katana/utils"
	"github.com/prebid/prebid-server/v2/openrtb_ext"
)

func (k Katana) handleProcessedAuctionHook(
	ctx context.Context,
	moduleCtx hookstage.ModuleInvocationContext,
	payload hookstage.ProcessedAuctionRequestPayload,
) (result hookstage.HookResult[hookstage.ProcessedAuctionRequestPayload], err error) {
	if payload.Request == nil {
		return result, nil
	}

	if len(moduleCtx.ModuleContext) == 0 {
		result.DebugMessages = append(result.DebugMessages, "error: module-ctx not found in handleProcessedAuctionHook()")
		result.Reject = true
		return result, nil
	}

	rctx, ok := moduleCtx.ModuleContext["rctx"].(models.RequestCtx)
	if !ok {
		result.DebugMessages = append(result.DebugMessages, "error: request-ctx not found in handleProcessedAuctionHook()")
		result.Reject = true
		return result, nil
	}

	defer func() {
		result.ModuleContext = make(hookstage.ModuleContext)
		result.ModuleContext["rctx"] = rctx
	}()

	rctx.PublisherIDStr = getPublisherId(payload.Request)
	rctx.PublisherID, err = strconv.Atoi(rctx.PublisherIDStr)
	if err != nil {
		result.DebugMessages = append(result.DebugMessages, "error: failed to convert publisher id to int")
		result.Reject = true
		return result, nil
	}

	// profile id from configs
	rctx.ProfileID = k.cfg.ProfileId
	rctx.TrackerEndpoint = k.cfg.TrackerUrl
	rctx.Id = payload.Request.ID

	rctx.Source, rctx.Origin = getSourceAndOrigin(payload.Request)
	rctx.PageURL = getPageURL(payload.Request)
	rctx.Platform = getPlatform(payload.Request)
	rctx.UA = getUserAgent(payload.Request, rctx.UA)
	rctx.IP = getIP(payload.Request, rctx.IP)
	rctx.Country = getCountry(payload.Request)
	rctx.DeviceCtx.Platform = utils.GetDevicePlatform(rctx, payload.Request)
	utils.PopulateDeviceContext(&rctx.DeviceCtx, payload.Request.Device)

	for _, imp := range payload.Request.GetImp() {
		impCtx := models.ImpCtx{
			TagId:       imp.TagID,
			Secure:      imp.Secure,
			BidFloor:    imp.BidFloor,
			BidFloorCur: imp.BidFloorCur,
		}

		if ext, err := imp.GetImpExt(); err == nil {
			prebidExt := ext.GetPrebid()
			impCtx.IsRewardInventory = prebidExt.IsRewardedInventory
		}

		rctx.ImpCtx[imp.ID] = impCtx
	}

	return result, nil
}

func getPublisherId(req *openrtb_ext.RequestWrapper) string {
	if req.Site != nil && req.Site.Publisher != nil {
		return req.Site.Publisher.ID
	}

	if req.App != nil && req.App.Publisher != nil {
		return req.App.Publisher.ID
	}
	return ""
}

func getPageURL(req *openrtb_ext.RequestWrapper) string {
	if req.App != nil && req.App.StoreURL != "" {
		return req.App.StoreURL
	} else if req.Site != nil && req.Site.Page != "" {
		return req.Site.Page
	}
	return ""
}

func getPlatform(request *openrtb_ext.RequestWrapper) string {
	var platform string
	if request == nil {
		return platform
	}
	if request.Site != nil {
		return models.PLATFORM_DISPLAY
	}
	if request.App != nil {
		return models.PLATFORM_APP
	}
	return platform
}

func getUserAgent(req *openrtb_ext.RequestWrapper, defaultUA string) string {
	userAgent := defaultUA
	if req != nil && req.Device != nil && len(req.Device.UA) > 0 {
		userAgent = req.Device.UA
	}
	return userAgent
}

func getSourceAndOrigin(bidRequest *openrtb_ext.RequestWrapper) (string, string) {
	var source, origin string
	if bidRequest.Site != nil {
		if len(bidRequest.Site.Domain) != 0 {
			source = bidRequest.Site.Domain
			origin = source
		} else if len(bidRequest.Site.Page) != 0 {
			source = getDomainFromUrl(bidRequest.Site.Page)
			origin = source

		}
	} else if bidRequest.App != nil {
		source = bidRequest.App.Bundle
		origin = source
	}
	return source, origin
}

func getDomainFromUrl(pageUrl string) string {
	u, err := url.Parse(pageUrl)
	if err != nil {
		return ""
	}

	return u.Host
}

func getIP(bidRequest *openrtb_ext.RequestWrapper, defaultIP string) string {
	ip := defaultIP
	if bidRequest != nil && bidRequest.Device != nil {
		if len(bidRequest.Device.IP) > 0 {
			ip = bidRequest.Device.IP
		} else if len(bidRequest.Device.IPv6) > 0 {
			ip = bidRequest.Device.IPv6
		}
	}
	return ip
}
func getCountry(bidRequest *openrtb_ext.RequestWrapper) string {
	if bidRequest.Device != nil && bidRequest.Device.Geo != nil && bidRequest.Device.Geo.Country != "" {
		return bidRequest.Device.Geo.Country
	}
	if bidRequest.User != nil && bidRequest.User.Geo != nil && bidRequest.User.Geo.Country != "" {
		return bidRequest.User.Geo.Country
	}
	return ""
}
