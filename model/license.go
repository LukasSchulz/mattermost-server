// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package model

import (
	"encoding/json"
	"io"
	"net/http"
)

const (
	EXPIRED_LICENSE_ERROR = "api.license.add_license.expired.app_error"
	INVALID_LICENSE_ERROR = "api.license.add_license.invalid.app_error"
)

type LicenseRecord struct {
	Id       string `json:"id"`
	CreateAt int64  `json:"create_at"`
	Bytes    string `json:"-"`
}

type License struct {
	Id           string    `json:"id"`
	IssuedAt     int64     `json:"issued_at"`
	StartsAt     int64     `json:"starts_at"`
	ExpiresAt    int64     `json:"expires_at"`
	Customer     *Customer `json:"customer"`
	Features     *Features `json:"features"`
	SkuName      string    `json:"sku_name"`
	SkuShortName string    `json:"sku_short_name"`
}

type Customer struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Company     string `json:"company"`
	PhoneNumber string `json:"phone_number"`
}

type Features struct {
	Users                     *int  `json:"users"`
	LDAP                      *bool `json:"ldap"`
	LDAPGroups                *bool `json:"ldap_groups"`
	MFA                       *bool `json:"mfa"`
	GoogleOAuth               *bool `json:"google_oauth"`
	Office365OAuth            *bool `json:"office365_oauth"`
	Compliance                *bool `json:"compliance"`
	Cluster                   *bool `json:"cluster"`
	Metrics                   *bool `json:"metrics"`
	MHPNS                     *bool `json:"mhpns"`
	SAML                      *bool `json:"saml"`
	Elasticsearch             *bool `json:"elastic_search"`
	Announcement              *bool `json:"announcement"`
	ThemeManagement           *bool `json:"theme_management"`
	EmailNotificationContents *bool `json:"email_notification_contents"`
	DataRetention             *bool `json:"data_retention"`
	MessageExport             *bool `json:"message_export"`
	CustomPermissionsSchemes  *bool `json:"custom_permissions_schemes"`
	CustomTermsOfService      *bool `json:"custom_terms_of_service"`
	GuestAccounts             *bool `json:"guest_accounts"`
	GuestAccountsPermissions  *bool `json:"guest_accounts_permissions"`
	IDLoadedPushNotifications *bool `json:"id_loaded"`
	LockTeammateNameDisplay   *bool `json:"lock_teammate_name_display"`

	// after we enabled more features we'll need to control them with this
	FutureFeatures *bool `json:"future_features"`
}

func (f *Features) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"ldap":                        *f.LDAP,
		"ldap_groups":                 *f.LDAPGroups,
		"mfa":                         *f.MFA,
		"google":                      *f.GoogleOAuth,
		"office365":                   *f.Office365OAuth,
		"compliance":                  *f.Compliance,
		"cluster":                     *f.Cluster,
		"metrics":                     *f.Metrics,
		"mhpns":                       *f.MHPNS,
		"saml":                        *f.SAML,
		"elastic_search":              *f.Elasticsearch,
		"email_notification_contents": *f.EmailNotificationContents,
		"data_retention":              *f.DataRetention,
		"message_export":              *f.MessageExport,
		"custom_permissions_schemes":  *f.CustomPermissionsSchemes,
		"guest_accounts":              *f.GuestAccounts,
		"guest_accounts_permissions":  *f.GuestAccountsPermissions,
		"id_loaded":                   *f.IDLoadedPushNotifications,
		"lock_teammate_name_display":  *f.LockTeammateNameDisplay,
		"future":                      *f.FutureFeatures,
	}
}

func (f *Features) SetDefaults() {
	if f.FutureFeatures == nil {
		f.FutureFeatures = NewBool(true)
	}

	if f.Users == nil {
		f.Users = NewInt(0)
	}

	if f.LDAP == nil {
		f.LDAP = NewBool(true)
	}

	if f.LDAPGroups == nil {
		f.LDAPGroups = NewBool(true)
	}

	if f.MFA == nil {
		f.MFA = NewBool(true)
	}

	if f.GoogleOAuth == nil {
		f.GoogleOAuth = NewBool(*f.FutureFeatures)
	}

	if f.Office365OAuth == nil {
		f.Office365OAuth = NewBool(*f.FutureFeatures)
	}

	if f.Compliance == nil {
		f.Compliance = NewBool(true)
	}

	if f.Cluster == nil {
		f.Cluster = NewBool(true)
	}

	if f.Metrics == nil {
		f.Metrics = NewBool(true)
	}

	if f.MHPNS == nil {
		f.MHPNS = NewBool(true)
	}

	if f.SAML == nil {
		f.SAML = NewBool(true)
	}

	if f.Elasticsearch == nil {
		f.Elasticsearch = NewBool(true)
	}

	if f.Announcement == nil {
		f.Announcement = NewBool(true)
	}

	if f.ThemeManagement == nil {
		f.ThemeManagement = NewBool(true)
	}

	if f.EmailNotificationContents == nil {
		f.EmailNotificationContents = NewBool(true)
	}

	if f.DataRetention == nil {
		f.DataRetention = NewBool(true)
	}

	if f.MessageExport == nil {
		f.MessageExport = NewBool(true)
	}

	if f.CustomPermissionsSchemes == nil {
		f.CustomPermissionsSchemes = NewBool(true)
	}

	if f.GuestAccounts == nil {
		f.GuestAccounts = NewBool(true)
	}

	if f.GuestAccountsPermissions == nil {
		f.GuestAccountsPermissions = NewBool(true)
	}

	if f.CustomTermsOfService == nil {
		f.CustomTermsOfService = NewBool(true)
	}

	if f.IDLoadedPushNotifications == nil {
		f.IDLoadedPushNotifications = NewBool(true)
	}

	if f.LockTeammateNameDisplay == nil {
		f.LockTeammateNameDisplay = NewBool(true)
	}
}

func (l *License) IsExpired() bool {
	return l.ExpiresAt < GetMillis()
}

func (l *License) IsStarted() bool {
	return l.StartsAt < GetMillis()
}

func (l *License) ToJson() string {
	b, _ := json.Marshal(l)
	return string(b)
}

// NewTestLicense returns a license that expires in the future and has the given features.
func NewTestLicense(features ...string) *License {
	ret := &License{
		ExpiresAt: GetMillis() + 90*24*60*60*1000,
		Customer:  &Customer{},
		Features:  &Features{},
	}
	ret.Features.SetDefaults()

	featureMap := map[string]bool{}
	for _, feature := range features {
		featureMap[feature] = true
	}
	featureJson, _ := json.Marshal(featureMap)
	json.Unmarshal(featureJson, &ret.Features)

	return ret
}

func LicenseFromJson(data io.Reader) *License {
	var o *License
	json.NewDecoder(data).Decode(&o)
	return o
}

func (lr *LicenseRecord) IsValid() *AppError {
	if len(lr.Id) != 26 {
		return NewAppError("LicenseRecord.IsValid", "model.license_record.is_valid.id.app_error", nil, "", http.StatusBadRequest)
	}

	if lr.CreateAt == 0 {
		return NewAppError("LicenseRecord.IsValid", "model.license_record.is_valid.create_at.app_error", nil, "", http.StatusBadRequest)
	}

	if len(lr.Bytes) == 0 || len(lr.Bytes) > 10000 {
		return NewAppError("LicenseRecord.IsValid", "model.license_record.is_valid.create_at.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

func (lr *LicenseRecord) PreSave() {
	lr.CreateAt = GetMillis()
}
