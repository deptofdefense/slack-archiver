// =================================================================
//
// Work of the U.S. Department of Defense, Defense Digital Service.
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package slack

type DirectMessage struct {
	ID      string   `json:"id"`
	Created int      `json:"created"`
	Members []string `json:"members"`
}
