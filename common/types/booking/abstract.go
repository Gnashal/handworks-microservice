package types

import "time"

type CleaningEquipment struct {
	ID       string
	Name     string
	Type     string
	PhotoURL string
}

type CleaningResources struct {
	ID       string
	Name     string
	Type     string
	PhotoURL string
}

type CleanerAssigned struct {
	ID               string
	CleanerFirstName string
	CleanerLastName  string
	PFPUrl           string
}
type AddonCleaningPrice struct {
	AddonName  string
	AddonPrice float32
}
type CleaningPrices struct {
	MainServicePrice float32
	AddonPrices      []AddonCleaningPrice
}

type ServiceDetail struct {
	General  *GeneralCleaningDetails
	Couch    *CouchCleaningDetails
	Mattress *MattressCleaningDetails
	Car      *CarCleaningDetails
	Post     *PostConstructionDetails
}

// because of course this is fucking different
// totally refactoring our detail logic because gubot jud kaayu siya
type ServiceDetails struct {
	ID          string
	ServiceType string
	Details     any
}

// detail factory types
type DetailType string

const (
	ServiceGeneral  DetailType = "GENERAL_CLEANING"
	ServiceCouch    DetailType = "COUCH"
	ServiceMattress DetailType = "MATTRESS"
	ServiceCar      DetailType = "CAR"
	ServicePost     DetailType = "POST"
)

// Used for unmarshaling dynamically
var DetailFactories = map[DetailType]func() any{
	ServiceGeneral:  func() any { return &GeneralCleaningDetails{} },
	ServiceCouch:    func() any { return &CouchCleaningDetails{} },
	ServiceMattress: func() any { return &MattressCleaningDetails{} },
	ServiceCar:      func() any { return &CarCleaningDetails{} },
	ServicePost:     func() any { return &PostConstructionDetails{} },
}

type GeneralCleaningDetails struct {
	HomeType string `json:"home_type"`
	SQM      int32  `json:"sqm"`
}

// Couch cleaning
type CouchCleaningSpecifications struct {
	CouchType string `json:"couch_type"`
	WidthCM   int32  `json:"width_cm"`
	DepthCM   int32  `json:"depth_cm"`
	HeightCM  int32  `json:"height_cm"`
	Quantity  int32  `json:"quantity"`
}

type CouchCleaningDetails struct {
	CleaningSpecs []CouchCleaningSpecifications `json:"cleaning_specs"`
	BedPillows    int32                         `json:"bed_pillows"`
}

// Mattress cleaning
type MattressCleaningSpecifications struct {
	BedType  string `json:"bed_type"`
	WidthCM  int32  `json:"width_cm"`
	DepthCM  int32  `json:"depth_cm"`
	HeightCM int32  `json:"height_cm"`
	Quantity int32  `json:"quantity"`
}

type MattressCleaningDetails struct {
	CleaningSpecs []MattressCleaningSpecifications `json:"cleaning_specs"`
}

// Car cleaning
type CarCleaningSpecifications struct {
	CarType  string `json:"car_type"`
	Quantity int32  `json:"quantity"`
}

type CarCleaningDetails struct {
	CleaningSpecs []CarCleaningSpecifications `json:"cleaning_specs"`
	ChildSeats    int32                       `json:"child_seats"`
}
type PostConstructionDetails struct {
	SQM int32 `json:"sqm"`
}
type BaseBookingDetails struct {
	ID                string
	CustID            string
	CustomerFirstName string
	CustomerLastName  string
	Address           Address
	StartSched        time.Time
	EndSched          time.Time
	DirtyScale        int32
	PaymentStatus     string
	ReviewStatus      string
	Photos            []string
	CreatedAt         time.Time
	UpdatedAt         *time.Time
	QuoteId           string
}
type Address struct {
	AddressHuman string
	AddressLat   float64
	AddressLng   float64
}
type BookingReply struct {
	Source     string              `json:"source"`
	Equipments []CleaningEquipment `json:"equipments,omitempty"`
	Resources  []CleaningResources `json:"resources,omitempty"`
	Cleaners   []CleanerAssigned   `json:"cleaners,omitempty"`
	Prices     CleaningPrices      `json:"prices,omitempty"`
	Error      string              `json:"error,omitempty"`
}

type MainServiceType string

const (
	ServiceTypeUnspecified MainServiceType = "SERVICE_TYPE_UNSPECIFIED"
	GeneralCleaning        MainServiceType = "GENERAL_CLEANING"
	CouchCleaning          MainServiceType = "COUCH"
	MattressCleaning       MainServiceType = "MATTRESS"
	CarCleaning            MainServiceType = "CAR"
	PostCleaning           MainServiceType = "POST"
)

type ServicesRequest struct {
	ServiceType MainServiceType `json:"service_type"`
	Details     ServiceDetail   `json:"details"`
}
type AddOnRequest struct {
	ServiceDetail ServicesRequest `json:"service_detail"`
}

type CreateBookingEvent struct {
	Base        BaseBookingDetails `json:"base"`
	MainService ServicesRequest    `json:"main_service"`
	Addons      []AddOnRequest     `json:"addons"`
}
type AddOns struct {
	ID            string
	ServiceDetail ServiceDetails
	Price         float32
}
type Booking struct {
	ID          string
	Base        BaseBookingDetails
	MainService ServiceDetails
	Addons      []AddOns
	Equipments  []CleaningEquipment
	Resources   []CleaningResources
	Cleaners    []CleanerAssigned
	TotalPrice  float32
}
