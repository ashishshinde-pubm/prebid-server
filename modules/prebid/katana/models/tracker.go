package models

type Tracker struct {
	PubID             int      //done
	PageURL           string   //done
	Timestamp         int64    //done
	IID               string   //done
	ProfileID         string   //done
	VersionID         string   //done
	SlotID            string   //done
	Adunit            string   //done
	PartnerInfo       Partner  //done
	RewardedInventory int      //done
	SURL              string   // Ignored
	Platform          int      // done
	Origin            string   // done
	ATTS              *float64 //done
	FloorSkippedFlag  *int     //done
	FloorModelVersion string   //done
	FloorSource       *int     //done
	FloorType         int      //done
}

type Partner struct {
	PartnerID      string  //done
	BidderCode     string  //done
	GrossECPM      float64 //done
	NetECPM        float64 //done
	BidID          string  //done
	OrigBidID      string  //done
	AdSize         string  //done
	AdDuration     int     //done
	Adformat       string  //done
	ServerSide     int     //done
	Advertiser     string  //done
	FloorValue     float64 //done
	FloorRuleValue float64 //done
	DealID         string  //done
}

type FloorsDetails struct {
	FloorType         int
	FloorModelVersion string
	FloorProvider     string
	Skipfloors        *int
	FloorFetchStatus  *int
	FloorSource       *int
}

type TrackerDetails struct {
	Tracker       Tracker
	TrackerURL    string
	BidType       string
	Price         float64
	PriceModel    string
	PriceCurrency string
}
