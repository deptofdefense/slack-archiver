// =================================================================
//
// Work of the U.S. Department of Defense, Defense Digital Service.
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package slack

type EnterpriseUser struct {
	ID             string   `json:"id"`
	EnterpriseID   string   `json:"enterprise_id"`
	EnterpriseName string   `json:"enterprise_name"`
	IsAdmin        bool     `json:"is_admin"`
	IsOwner        bool     `json:"is_owner"`
	Teams          []string `json:"teams"`
}
