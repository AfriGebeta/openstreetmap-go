package etrade

import "time"

type ClientFor string

var ClientForValue = struct {
	RegistrationInfo      ClientFor
	BusinessLicenseNumber ClientFor
}{
	RegistrationInfo:      "REGISTRATION_INFO",
	BusinessLicenseNumber: "BUSINESS_LICENSE_NUMBER",
}

var ClientForValues = []ClientFor{
	ClientForValue.RegistrationInfo,
	ClientForValue.BusinessLicenseNumber,
}

type BusinessInfo struct {
	MainGuid                   string                   `json:"MainGuid"`
	OwnerTIN                   string                   `json:"OwnerTIN"`
	DateRegistered             string                   `json:"DateRegistered"` // use time.Time if you parse it
	TradeName                  string                   `json:"TradeName"`
	LicenceNumber              string                   `json:"LicenceNumber"`
	Status                     int                      `json:"Status"`
	Capital                    int                      `json:"Capital"`
	AssociateShortInfos        []AssociateShortInfo     `json:"AssociateShortInfos"`
	AddressInfo                AddressInfo              `json:"BusinessAddressInfo"`
	BusinessLicensingGroupMain []BusinessLicensingGroup `json:"BusinessLicensingGroupMain"`
	RenewedTo                  time.Time                `json:"RenewedTo"`
	RenewedToDateString        string                   `json:"RenewedToDateString"`
	RenewalDate                string                   `json:"RenewalDate"`
	RenewedFrom                string                   `json:"RenewedFrom"`
}

type AddressInfo struct {
	Region       string `json:"Region"`
	Zone         string `json:"Zone"`
	Woreda       string `json:"Woreda"`
	Kebele       string `json:"Kebele"`
	HouseNo      string `json:"HouseNo"`
	MobilePhone  string `json:"MobilePhone"`
	RegularPhone string `json:"RegularPhone"`
}

type BusinessLicensingGroup struct {
	MainGuid         string  `json:"MainGuid"`
	BusinessMainGuid string  `json:"BusinessMainGuid"`
	MajorDivision    int     `json:"MajorDivision"`
	Division         int     `json:"Division"`
	MajorGroup       int     `json:"MajorGroup"`
	BGroup           int     `json:"BGroup"`
	SubGroup         int     `json:"BusinessSubGroup"`
	Tinc             *string `json:"Tinc"`         // nullable
	BusinessMain     *string `json:"BusinessMain"` // nullable
}

type BusinessRegistration struct {
	Tin                 string               `json:"Tin"`
	LegalCondtion       string               `json:"LegalCondtion"` // assuming string type based on quotes
	RegNo               string               `json:"RegNo"`
	RegDate             string               `json:"RegDate"` // can be time.Time with parsing
	BusinessName        string               `json:"BusinessName"`
	BusinessNameAmh     string               `json:"BusinessNameAmh"`
	PaidUpCapital       int                  `json:"PaidUpCapital"`
	AssociateShortInfos []AssociateShortInfo `json:"AssociateShortInfos"`
	Businesses          []Business           `json:"Businesses"`
}

type AssociateShortInfo struct {
	Position       *string `json:"Position"`
	ManagerName    string  `json:"ManagerName"`
	ManagerNameEng string  `json:"ManagerNameEng"`
	Photo          string  `json:"Photo"`
	MobilePhone    *string `json:"MobilePhone"`
	RegularPhone   *string `json:"RegularPhone"`
}

type Business struct {
	MainGuid                   string     `json:"MainGuid"`
	OwnerTIN                   string     `json:"OwnerTIN"`
	DateRegistered             string     `json:"DateRegistered"`
	TradeNameAmh               string     `json:"TradeNameAmh"`
	TradesName                 string     `json:"TradesName"`
	LicenceNumber              string     `json:"LicenceNumber"`
	RenewalDate                string     `json:"RenewalDate"`
	RenewedFrom                string     `json:"RenewedFrom"`
	RenewedTo                  string     `json:"RenewedTo"`
	BusinessLicensingGroupMain *string    `json:"BusinessLicensingGroupMain"` // nullable
	SubGroups                  []SubGroup `json:"SubGroups"`
}

type SubGroup struct {
	Code        int    `json:"Code"`
	Description string `json:"Description"`
}
