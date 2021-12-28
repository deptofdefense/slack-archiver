// =================================================================
//
// Work of the U.S. Department of Defense, Defense Digital Service.
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package slack

type IntegrationLogMessage struct {
	UserID     int    `json:"user_id"`
	UserName   string `json:"user_name"`
	Date       string `json:"date"`
	ChangeType string `json:"change_type"`
	AdminAppID string `json:"admin_app_id"`
	Resolution string `json:"resolution"`
	AppID      string `json:"app_id"`
	AppType    string `json:"app_type"`
}
