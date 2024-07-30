package utils

import (
	"regexp"
	"strings"

	"github.com/prebid/openrtb/v20/adcom1"
	"github.com/prebid/openrtb/v20/openrtb2"
	"github.com/prebid/prebid-server/v2/modules/prebid/katana/models"
	"github.com/prebid/prebid-server/v2/openrtb_ext"
	"github.com/prebid/prebid-server/v2/util/ptrutil"
)

const (
	OpenRTBDeviceOsIosRegexPattern     string = `(ios).*`
	OpenRTBDeviceOsAndroidRegexPattern string = `(android).*`
	IosUARegexPattern                  string = `(iphone|ipad|darwin).*`
	AndroidUARegexPattern              string = `android.*`
	MobileDeviceUARegexPattern         string = `(mobi|tablet|ios).*`
	ConnectedDeviceUARegexPattern      string = `Roku|SMART-TV|SmartTV|AndroidTV|Android TV|AppleTV|Apple TV|VIZIO|PHILIPS|BRAVIA|PlayStation|Chromecast|ExoPlayerLib|MIBOX3|Xbox|ComcastAppPlatform|AFT|HiSmart|BeyondTV|D.*ATV|PlexTV|Xstream|MiTV|AI PONT`
)

var (
	openRTBDeviceOsAndroidRegex *regexp.Regexp
	androidUARegex              *regexp.Regexp
	iosUARegex                  *regexp.Regexp
	openRTBDeviceOsIosRegex     *regexp.Regexp
	mobileDeviceUARegex         *regexp.Regexp
	ctvRegex                    *regexp.Regexp
)

func init() {
	openRTBDeviceOsAndroidRegex = regexp.MustCompile(OpenRTBDeviceOsAndroidRegexPattern)
	androidUARegex = regexp.MustCompile(AndroidUARegexPattern)
	iosUARegex = regexp.MustCompile(IosUARegexPattern)
	openRTBDeviceOsIosRegex = regexp.MustCompile(OpenRTBDeviceOsIosRegexPattern)
	mobileDeviceUARegex = regexp.MustCompile(MobileDeviceUARegexPattern)
	ctvRegex = regexp.MustCompile(ConnectedDeviceUARegexPattern)
}

func isMobile(deviceType adcom1.DeviceType, userAgentString string) bool {
	if deviceType != 0 {
		return deviceType == adcom1.DeviceMobile || deviceType == adcom1.DeviceTablet || deviceType == adcom1.DevicePhone
	}

	if mobileDeviceUARegex.Match([]byte(strings.ToLower(userAgentString))) {
		return true
	}
	return false
}

func isIos(os string, userAgentString string) bool {
	if openRTBDeviceOsIosRegex.Match([]byte(strings.ToLower(os))) || iosUARegex.Match([]byte(strings.ToLower(userAgentString))) {
		return true
	}
	return false
}

func isAndroid(os string, userAgentString string) bool {
	if openRTBDeviceOsAndroidRegex.Match([]byte(strings.ToLower(os))) || androidUARegex.Match([]byte(strings.ToLower(userAgentString))) {
		return true
	}
	return false
}

func isCTV(userAgent string) bool {
	return ctvRegex.Match([]byte(userAgent))
}

// getDevicePlatform determines the device from which request has been generated
func GetDevicePlatform(rCtx models.RequestCtx, bidRequest *openrtb_ext.RequestWrapper) models.DevicePlatform {
	userAgentString := rCtx.UA

	switch rCtx.Platform {
	case models.PLATFORM_AMP:
		return models.DevicePlatformMobileWeb

	case models.PLATFORM_APP:
		//Its mobile; now determine ios or android
		var os = ""
		if bidRequest != nil && bidRequest.Device != nil && len(bidRequest.Device.OS) != 0 {
			os = bidRequest.Device.OS
		}
		if isIos(os, userAgentString) {
			return models.DevicePlatformMobileAppIos
		} else if isAndroid(os, userAgentString) {
			return models.DevicePlatformMobileAppAndroid
		}

	case models.PLATFORM_DISPLAY:
		//Its web; now determine mobile or desktop
		var deviceType adcom1.DeviceType
		if bidRequest != nil && bidRequest.Device != nil && bidRequest.Device.DeviceType != 0 {
			deviceType = bidRequest.Device.DeviceType
		}
		if isMobile(deviceType, userAgentString) {
			return models.DevicePlatformMobileWeb
		}
		return models.DevicePlatformDesktop

	case models.PLATFORM_VIDEO:
		var deviceType adcom1.DeviceType
		if bidRequest != nil && bidRequest.Device != nil && bidRequest.Device.DeviceType != 0 {
			deviceType = bidRequest.Device.DeviceType
		}
		isCtv := isCTV(userAgentString)

		if deviceType != 0 {
			if deviceType == adcom1.DeviceTV || deviceType == adcom1.DeviceConnected || deviceType == adcom1.DeviceSetTopBox {
				return models.DevicePlatformConnectedTv
			}
		}

		if deviceType == 0 && isCtv {
			return models.DevicePlatformConnectedTv
		}

		if bidRequest != nil && bidRequest.Site != nil {
			//Its web; now determine mobile or desktop
			var deviceType adcom1.DeviceType
			if bidRequest.Device != nil {
				deviceType = bidRequest.Device.DeviceType
			}
			if isMobile(deviceType, userAgentString) {
				return models.DevicePlatformMobileWeb
			}
			return models.DevicePlatformDesktop
		}

		if bidRequest != nil && bidRequest.App != nil {
			//Its mobile; now determine ios or android
			var os = ""
			if bidRequest.Device != nil && len(bidRequest.Device.OS) != 0 {
				os = bidRequest.Device.OS
			}

			if isIos(os, userAgentString) {
				return models.DevicePlatformMobileAppIos
			} else if isAndroid(os, userAgentString) {
				return models.DevicePlatformMobileAppAndroid
			}
		}

	default:
		return models.DevicePlatformNotDefined

	}

	return models.DevicePlatformNotDefined
}

func PopulateDeviceContext(dvc *models.DeviceCtx, device *openrtb2.Device) {
	if device == nil {
		return
	}
	//this is needed in determine ifa_type parameter
	dvc.DeviceIFA = strings.TrimSpace(device.IFA)

	if device.Ext == nil {
		return
	}

	//unmarshal device ext
	var deviceExt models.ExtDevice
	if err := deviceExt.UnmarshalJSON(device.Ext); err != nil {
		return
	}
	dvc.Ext = &deviceExt

	//update device IFA Details
	updateDeviceIFADetails(dvc)
}

func updateDeviceIFADetails(dvc *models.DeviceCtx) {
	if dvc == nil || dvc.Ext == nil {
		return
	}

	deviceExt := dvc.Ext
	extIFAType, ifaTypeFound := deviceExt.GetIFAType()
	extSessionID, _ := deviceExt.GetSessionID()

	if ifaTypeFound {
		if dvc.DeviceIFA != "" {
			if ifaTypeID, ok := models.DeviceIFATypeID[strings.ToLower(extIFAType)]; !ok {
				deviceExt.DeleteIFAType()
			} else {
				dvc.IFATypeID = &ifaTypeID
				deviceExt.SetIFAType(extIFAType)
			}
		} else if extSessionID != "" {
			dvc.DeviceIFA = extSessionID
			dvc.IFATypeID = ptrutil.ToPtr(models.DeviceIfaTypeIdSessionId)
			deviceExt.SetIFAType(models.DeviceIFATypeSESSIONID)
		} else {
			deviceExt.DeleteIFAType()
		}
	} else if extSessionID != "" {
		dvc.DeviceIFA = extSessionID
		dvc.IFATypeID = ptrutil.ToPtr(models.DeviceIfaTypeIdSessionId)
		deviceExt.SetIFAType(models.DeviceIFATypeSESSIONID)
	}
}
