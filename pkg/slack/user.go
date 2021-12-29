// =================================================================
//
// Work of the U.S. Department of Defense, Defense Digital Service.
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package slack

type User struct {
	Color                  string          `json:"color"`
	Deleted                bool            `json:"deleted"`
	EnterpriseName         string          `json:"enterprise_name"`
	EnterpriseUser         *EnterpriseUser `json:"enterprise_user"`
	EnterpriseID           string          `json:"enterprise_id"`
	ID                     string          `json:"id"`
	IsAdmin                bool            `json:"is_admin"`
	IsAppUser              bool            `json:"is_app_user"`
	IsBot                  bool            `json:"is_bot"`
	IsEmailConfirmed       bool            `json:"is_email_confirmed"`
	IsInvitedUser          bool            `json:"is_invited_user"`
	IsOwner                bool            `json:"is_owner"`
	IsPrimaryOwner         bool            `json:"is_primary_owner"`
	IsRestricted           bool            `json:"is_restricted"`
	IsUltraRestricted      bool            `json:"is_ultra_restricted"`
	Name                   string          `json:"name"`
	Profile                *Profile        `json:"profile"`
	RealName               string          `json:"real_name"`
	Teams                  []string        `json:"teams"`
	TimeZone               string          `json:"tz"`
	TimeZoneOffset         int             `json:"tz_offset"`
	TimeZoneLabel          string          `json:"tz_label"`
	Updated                int             `json:"updated"`
	WhoCanShareContactCard string          `json:"who_can_share_contact_card"`
}
