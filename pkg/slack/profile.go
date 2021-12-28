// =================================================================
//
// Work of the U.S. Department of Defense, Defense Digital Service.
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package slack

type FieldValue struct {
	Value     string `json:"value"`
	Alternate string `json:"alt"`
}

type Profile struct {
	AvatarHash             string                `json:"avatar_hash"`
	DisplayName            string                `json:"display_name"`
	DisplayNameNormalized  string                `json:"display_name_normalized"`
	Email                  string                `json:"email"`
	Fields                 map[string]FieldValue `json:"fields"`
	FirstName              string                `json:"first_name"`
	GuestInvitedBy         string                `json:"guest_invited_by,omitempty"`
	IsCustomImage          bool                  `json:"is_custom_image"`
	Image24                string                `json:"image_24"`
	Image32                string                `json:"image_32"`
	Image48                string                `json:"image_48"`
	Image72                string                `json:"image_72"`
	Image192               string                `json:"image_192"`
	Image512               string                `json:"image_512"`
	Image1024              string                `json:"image_1024"`
	ImageOriginal          string                `json:"image_original"`
	LastName               string                `json:"last_name"`
	Phone                  string                `json:"phone"`
	RealName               string                `json:"real_name"`
	RealNameNormalized     string                `json:"real_name_normalized"`
	Skype                  string                `json:"skype"`
	StatusEmoji            string                `json:"status_emoji"`
	StatusEmojiDisplayInfo []string              `json:"status_emoji_display_info"`
	StatusExpiration       int                   `json:"status_expiration"`
	StatusText             string                `json:"status_text"`
	StatusTextCanonical    string                `json:"status_text_canonical"`
	Team                   string                `json:"team"`
	Title                  string                `json:"title"`
}
