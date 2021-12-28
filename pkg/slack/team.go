// =================================================================
//
// Work of the U.S. Department of Defense, Defense Digital Service.
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package slack

type Team struct {
	Name     string     `json:"name"`     // The name of the team
	Channels []*Channel `json:"channels"` // the public channels for the team
	Groups   []*Group   `json:"groups"`   // the private groups for the team
	Users    []*User    `json:"users"`    // the list of users for the team
}
