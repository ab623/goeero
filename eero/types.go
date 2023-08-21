package eero

import (
	"fmt"
	"time"
)

type MetaResponse struct {
	Meta struct {
		Code         int       `json:"code"`
		ServerTime   time.Time `json:"server_time"`
		ErrorMessage string    `json:"error"`
	} `json:"meta"`
}

func (e *MetaResponse) Error() string {
	return fmt.Sprintf("[Eero API Error] Http Status %d: %s", e.Meta.Code, e.Meta.ErrorMessage)
}

type LoginRequest struct {
	Identifier string `json:"login"`
}

type LoginVerifyRequest struct {
	Code string `json:"code"`
}

type LoginResponse struct {
	MetaResponse
	Data LoginData
}

type LoginData struct {
	UserToken string `json:"user_token"`
}

type AccountsResponse struct {
	MetaResponse
	Data AccountsData `json:"data"`
}

type DevicesResponse struct {
	MetaResponse
	Data []DeviceData `json:"data"`
}

type AccountsData struct {
	Name  string `json:"name"`
	Phone struct {
		Value          string `json:"value"`
		CountryCode    string `json:"country_code"`
		NationalNumber string `json:"national_number"`
		Verified       bool   `json:"verified"`
	} `json:"phone"`
	Email struct {
		Value    string `json:"value"`
		Verified bool   `json:"verified"`
	} `json:"email"`
	LogID string `json:"log_id"`
	// OrganizationID string `json:"organization_id"`
	// ImageAssets    string `json:"image_assets"`
	Networks struct {
		Count int           `json:"count"`
		Data  []NetworkData `json:"data"`
	} `json:"networks"`
	// Auth struct {
	// 	Type       string `json:"type"`
	// 	ProviderID string `json:"provider_id"`
	// 	ServiceID  string `json:"service_id"`
	// } `json:"auth"`
	// Role                      string `json:"role"`
	// IsBetaBugReporterEligible bool `json:"is_beta_bug_reporter_eligible"`
	// CanTransfer               string `json:"can_transfer"`
	// IsPremiumCapable          bool `json:"is_premium_capable"`
	// PaymentFailed             string `json:"payment_failed"`
	// PremiumStatus             string `json:"premium_status"`
	// PremiumDetails            struct {
	// 	TrialEnds            string `json:"trial_ends"`
	// 	HasPaymentInfo       string `json:"has_payment_info"`
	// 	Tier                 string `json:"tier"`
	// 	IsIapCustomer        string `json:"is_iap_customer"`
	// 	PaymentMethod        string `json:"payment_method"`
	// 	Interval             string `json:"interval"`
	// 	NextBillingEventDate string `json:"next_billing_event_date"`
	// } `json:"premium_details"`
	// PushSettings struct {
	// 	NetworkOffline string `json:"networkOffline"`
	// 	NodeOffline    string `json:"nodeOffline"`
	// } `json:"push_settings"`
	// TrustCertificatesEtag string `json:"trust_certificates_etag"`
	// Consents              struct {
	// 	MarketingEmails struct {
	// 		Consented string `json:"consented"`
	// 	} `json:"marketing_emails"`
	// } `json:"consents"`
	// CanMigrateToAmazonLogin string `json:"can_migrate_to_amazon_login"`
	// EeroForBusiness         string `json:"eero_for_business"`
	// MduProgram              string `json:"mdu_program"`
}

type NetworkData struct {
	URL     string    `json:"url"`
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
	// AccessExpiresOn  time.Time `json:"access_expires_on"`
	// AmazonDirectedID string    `json:"amazon_directed_id"`
}

type DeviceData struct {
	URL            string    `json:"url"`
	Mac            string    `json:"mac"`
	Manufacturer   string    `json:"manufacturer"`
	IP             string    `json:"ip"`
	Ips            []string  `json:"ips"`
	Nickname       string    `json:"nickname"`
	Connected      bool      `json:"connected"`
	Wireless       bool      `json:"wireless"`
	ConnectionType string    `json:"connection_type"`
	LastActive     time.Time `json:"last_active"`
	FirstActive    time.Time `json:"first_active"`
	Interface      struct {
		Frequency     string `json:"frequency"`
		FrequencyUnit string `json:"frequency_unit"`
	} `json:"interface"`
	DeviceType  string `json:"device_type"`
	Blacklisted bool   `json:"blacklisted"`
	IsGuest     bool   `json:"is_guest"`
	Paused      bool   `json:"paused"`
	SSID        string `json:"ssid"`
	DisplayName string `json:"display_name"`
}
