// =================================================================
//
// Work of the U.S. Department of Defense, Defense Digital Service.
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package slack

import (
	"encoding/json"
)

type MessageProfile struct {
	AvatarHash        string `json:"avatar_hash"`
	DisplayName       string `json:"display_name"`
	Name              string `json:"name"`
	FirstName         string `json:"first_name"`
	RealName          string `json:"real_name"`
	Image72           string `json:"image_72"`
	IsRestricted      bool   `json:"is_restricted"`
	IsUltraRestricted bool   `json:"is_ultra_restricted"`
	Team              string `json:"team"`
}

/*
type MessageText struct {
	Emoji bool `json:"emoji"`
	Type string `json:"type"`
	Text string `json:"text,omitempty"` // if type==text, then text is used
	Verbatim bool `json:"verbatim,omitempty"`
}
*/

type MessageBlockElement struct {
	ActionID string                 `json:"action_id"` // if type==button, then action_id is used
	Type     string                 `json:"type"`
	Text     json.RawMessage        `json:"text,omitempty"` // could be string or objecty
	Name     string                 `json:"name,omitempty"` // if type==emoji, then name is used
	Elements []*MessageBlockElement `json:"elements,omitempty"`
}

type MessageBlock struct {
	BlockID  string                 `json:"block_id"`
	Type     string                 `json:"type"`
	Elements []*MessageBlockElement `json:"elements"`
}

type MessageReaction struct {
	Name  string   `json:"name"`
	Users []string `json:"users"`
	Count int      `json:"count"`
}

type MessageFile struct {
	ID                 string `json:"id"`
	Created            int    `json:"created,omitempty"`
	Timestamp          int    `json:"timestamp,omitempty"`
	Name               string `json:"name,omitempty"`
	Title              string `json:"title,omitempty"`
	MimeType           string `json:"mimetype,omitempty"`
	FileType           string `json:"filetype,omitempty"`
	PrettyType         string `json:"pretty_type,omitempty"`
	User               string `json:"user,omitempty"`
	Editable           bool   `json:"editable,omitempty"`
	Size               int    `json:"size,omitempty"`
	Mode               string `json:"mode"`
	IsExternal         bool   `json:"is_external,omitempty"`
	ExternalType       string `json:"external_type,omitempty"`
	IsPublic           bool   `json:"is_public,omitempty"`
	PublicURLShared    bool   `json:"public_url_shared,omitempty"`
	DisplayAsBot       bool   `json:"display_as_bot,omitempty"`
	Username           string `json:"username,omitempty"`
	URLPrivate         string `json:"url_private,omitempty"`
	URLPrivateDownload string `json:"url_private_download,omitempty"`
	MediaDisplayType   string `json:"media_display_type,omitempty"`
	Thumb64            string `json:"thumb_64,omitempty"`
	Thumb80            string `json:"thumb_80,omitempty"`
	Thumb160           string `json:"thumb_160,omitempty"`
	Thumb360           string `json:"thumb_360,omitempty"`
	Thumb360Width      int    `json:"thumb_360_w,omitempty"`
	Thumb360Height     int    `json:"thumb_360_h,omitempty"`
	Thumb480           string `json:"thumb_480,omitempty"`
	Thumb480Width      int    `json:"thumb_480_w,omitempty"`
	Thumb480Height     int    `json:"thumb_480_h,omitempty"`
	Thumb720           string `json:"thumb_720,omitempty"`
	Thumb720Width      int    `json:"thumb_720_w,omitempty"`
	Thumb720Height     int    `json:"thumb_720_h,omitempty"`
	Thumb800           string `json:"thumb_800,omitempty"`
	Thumb800Width      int    `json:"thumb_800_w,omitempty"`
	Thumb800Height     int    `json:"thumb_800_h,omitempty"`
	Thumb960           string `json:"thumb_960,omitempty"`
	Thumb960Width      int    `json:"thumb_960_w,omitempty"`
	Thumb960Height     int    `json:"thumb_960_h,omitempty"`
	Thumb1024          string `json:"thumb_1024,omitempty"`
	Thumb1024Width     int    `json:"thumb_1024_w,omitempty"`
	Thumb1024Height    int    `json:"thumb_1024_h,omitempty"`
	ThumbTiny          string `json:"thumb_tiny,omitempty"`
	ImageEXIFRotation  int    `json:"image_exif_rotation,omitempty"`
	OriginalWidth      int    `json:"original_width,omitempty"`
	OriginalHeight     int    `json:"original_height,omitempty"`
	Permalink          string `json:"permalink,omitempty"`
	PermalinkPublic    string `json:"permalink_public,omitempty"`
	IsStarred          bool   `json:"is_starred,omitempty"`
	HasRichPreview     bool   `json:"has_rich_preview,omitempty"`
	FileAccess         string `json:"file_access,omitempty"`
}

func (f MessageFile) IsHosted() bool {
	return f.Mode == "hosted"
}

func (f MessageFile) IsTombstone() bool {
	return f.Mode == "tombstone"
}

type MessageReply struct {
	User      string `json:"user"`
	Timestamp string `json:"ts"`
}

type Message struct {
	Blocks          []MessageBlock    `json:"blocks"`
	Files           []MessageFile     `json:"files,omitempty"`
	ClientMessageID string            `json:"client_msg_id"`
	IsLocked        bool              `json:"is_locked"`
	LatestReply     string            `json:"latest_reply,omitempty"`
	Reactions       []MessageReaction `json:"reactions,omitempty"`
	Replies         []MessageReply    `json:"replies,omitempty"`
	ReplyCount      int               `json:"reply_count,omitempty"`
	ReplyUsersCount int               `json:"reply_users_count,omitempty"`
	ReplyUsers      []string          `json:"reply_users,omitempty"`
	SourceTeam      string            `json:"source_team"`
	Subscribed      bool              `json:"subscribed"`
	Text            string            `json:"text"`
	ThreadTimestamp string            `json:"thread_ts,omitempty"`
	Timestamp       string            `json:"ts"`
	Type            string            `json:"type"`
	Team            string            `json:"team"`
	Upload          bool              `json:"upload,omitempty"`
	User            string            `json:"user"`
	UserTeam        string            `json:"user_team"`
	UserProfile     MessageProfile    `json:"user_profile"`
}
