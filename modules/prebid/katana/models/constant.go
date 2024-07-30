package models

const (
	PLATFORM_DISPLAY = "display"
	PLATFORM_AMP     = "amp"
	PLATFORM_APP     = "in-app"
	PLATFORM_VIDEO   = "video"
)

// impression tracker url parameters
const (
	// constants for query parameter names for tracker call
	TRKPubID                 = "pubid"
	TRKPageURL               = "purl"
	TRKTimestamp             = "tst"
	TRKIID                   = "iid"
	TRKProfileID             = "pid"
	TRKVersionID             = "pdvid"
	TRKIP                    = "ip"
	TRKUserAgent             = "ua"
	TRKSlotID                = "slot"
	TRKAdunit                = "au"
	TRKRewardedInventory     = "rwrd"
	TRKPartnerID             = "pn"
	TRKBidderCode            = "bc"
	TRKKGPV                  = "kgpv"
	TRKGrossECPM             = "eg"
	TRKNetECPM               = "en"
	TRKBidID                 = "bidid"
	TRKOrigBidID             = "origbidid"
	TRKQMARK                 = "?"
	TRKAmpersand             = "&"
	TRKSSAI                  = "ssai"
	TRKPlatform              = "plt"
	TRKAdSize                = "psz"
	TRKTestGroup             = "tgid"
	TRKAdvertiser            = "adv"
	TRKPubDomain             = "orig"
	TRKServerSide            = "ss"
	TRKAdformat              = "af"
	TRKAdDuration            = "dur"
	TRKAdPodExist            = "aps"
	TRKFloorType             = "ft"
	TRKFloorModelVersion     = "fmv"
	TRKFloorSkippedFlag      = "fskp"
	TRKFloorSource           = "fsrc"
	TRKFloorValue            = "fv"
	TRKFloorRuleValue        = "frv"
	TRKServerLogger          = "sl"
	TRKDealID                = "di"
	TRKCustomDimensions      = "cds"
	TRKATTS                  = "atts"
	TRKProfileType           = "pt"
	TRKProfileTypePlatform   = "ptp"
	TRKAppPlatform           = "ap"
	TRKAppIntegrationPath    = "aip"
	TRKAppSubIntegrationPath = "asip"
	TRKPriceBucket           = "pb"
)

const (
	//constant for adformat
	Banner = "banner"
	Video  = "video"
	Native = "native"
)

var (
	//EmptyVASTResponse Empty VAST Response
	EmptyVASTResponse = []byte(`<VAST version="2.0"/>`)
	//EmptyString to check for empty value
	EmptyString = ""
	//ContentType HTTP Response Header Content Type
	ContentType = `Content-Type`
	//ContentTypeApplicationJSON HTTP Header Content-Type Value
	ContentTypeApplicationJSON = `application/json`
	//ContentTypeApplicationXML HTTP Header Content-Type Value
	ContentTypeApplicationXML = `application/xml`
	//EmptyJSONResponse Empty JSON Response
	EmptyJSONResponse = []byte{}
	//VASTErrorResponse VAST Error Response
	VASTErrorResponse = `<VAST version="2.0"><Ad><InLine><Extensions><Extension><OWStatus><Error code="%v">%v</Error></OWStatus></Extension></Extensions></InLine></Ad></VAST>`
	//TrackerCallWrap
	TrackerCallWrap = `<div style="position:absolute;left:0px;top:0px;visibility:hidden;"><img src="${escapedUrl}"></div>`
	//Tracker Format for Native
	NativeTrackerMacro = `{"event":1,"method":1,"url":"${trackerUrl}"}`
	//TrackerCallWrapOMActive for Open Measurement in In-App Banner
	TrackerCallWrapOMActive = `<script id="OWPubOMVerification" data-owurl="${escapedUrl}" src="${OMScript}"></script>`
	//Universal Pixel Macro
	UniversalPixelMacroForUrl = `<div style="position:absolute;left:0px;top:0px;visibility:hidden;"><img src="${pixelUrl}"></div>`
)

const (
	//constants for video
	VIDEO_CACHE_PATH          = "/cache"
	VideoSizeSuffix           = "v"
	PartnerURLPlaceholder     = "$PARTNER_URL_PLACEHOLDER"
	TrackerPlaceholder        = "$TRACKER_PLACEHOLDER"
	ErrorPlaceholder          = "$ERROR_PLACEHOLDER"
	ImpressionElement         = "Impression"
	ErrorElement              = "Error"
	VASTAdElement             = ".//VAST/Ad"
	AdWrapperElement          = "./Wrapper"
	AdInlineElement           = "./InLine"
	VASTAdWrapperElement      = ".//VAST/Ad/Wrapper"
	VASTAdInlineElement       = ".//VAST/Ad/InLine"
	CdataPrefix               = "<![CDATA["
	CdataSuffix               = "]]>"
	HTTPProtocol              = "http"
	HTTPSProtocol             = "https"
	VASTImpressionURLTemplate = `<Impression><![CDATA[` + TrackerPlaceholder + `]]></Impression>`
	VASTErrorURLTemplate      = `<Error><![CDATA[` + ErrorPlaceholder + `]]></Error>`
	VastWrapper               = `<VAST version="3.0"><Ad id="1"><Wrapper><AdSystem>PubMatic Wrapper</AdSystem><VASTAdTagURI><![CDATA[$PARTNER_URL_PLACEHOLDER]]></VASTAdTagURI>` + VASTImpressionURLTemplate + `</Wrapper></Ad></VAST>`
)

const (
	//VideoVASTTag video VAST parameter constant
	VideoVASTTag = "./VAST"
	//VideoVASTVersion video version parameter constant
	VideoVASTVersion = "version"
	//VideoVASTVersion2_0 video version 2.0 parameter constant
	VideoVASTVersion2_0 = "2.0"
	//VideoVASTVersion3_0 video version 3.0 parameter constant
	VideoVASTVersion3_0 = "3.0"
	//VideoVASTAdWrapperTag video ad/wrapper element constant
	VideoVASTAdWrapperTag = "./Ad/Wrapper"
	//VideoVASTAdInLineTag video ad/inline element constant
	VideoVASTAdInLineTag = "./Ad/InLine"
	//VideoExtensionsTag video extensions element constant
	VideoExtensionsTag = "Extensions"
	//VideoExtensionTag video extension element constant
	VideoExtensionTag = "Extension"
	//VideoPricingTag video pricing element constant
	VideoPricingTag = "Pricing"
	//VideoPricingModel video model attribute constant
	VideoPricingModel = "model"
	//VideoPricingModelCPM video cpm attribute value constant
	VideoPricingModelCPM = "CPM"
	//VideoPricingCurrencyUSD video USD default currency constant
	VideoPricingCurrencyUSD = "USD"
	//VideoPricingCurrency video currency constant
	VideoPricingCurrency = "currency"
	//VideoTagLookupStart video xpath constant
	VideoTagLookupStart = "./"
	//VideoTagForwardSlash video forward slash for xpath constant
	VideoTagForwardSlash = "/"
	//VideoVAST2ExtensionPriceElement video parameter constant
	VideoVAST2ExtensionPriceElement = VideoTagLookupStart + VideoExtensionTag + VideoTagForwardSlash + VideoPricingTag
)
