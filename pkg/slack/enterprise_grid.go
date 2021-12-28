// =================================================================
//
// Work of the U.S. Department of Defense, Defense Digital Service.
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package slack

import (
	"fmt"
	"strings"
)

type EnterpriseGrid struct {
	Archive                   *Archive
	DirectMessages            []*DirectMessage
	Groups                    []*Group
	IntegrationLogMessages    []*IntegrationLogMessage
	MultiPartyInstantMessages []*MultiPartyInstantMessage
	OrganizationUsers         []*User
	Teams                     []*Team
}

func (e *EnterpriseGrid) UnmarshalFile(name string, v interface{}) error {
	err := e.Archive.UnmarshalFile(name, v)
	if err != nil {
		return fmt.Errorf("error unmarshaling file from archive: %w", err)
	}
	return nil
}

func (e *EnterpriseGrid) GetDirectMessages() []*DirectMessage {
	return e.DirectMessages
}

func (e *EnterpriseGrid) GetMultiPartyInstantMessages() []*MultiPartyInstantMessage {
	return e.MultiPartyInstantMessages
}

func (e *EnterpriseGrid) GetTeams() []*Team {
	return e.Teams
}

func (e *EnterpriseGrid) GetMessages(prefix string) ([]*Message, error) {
	messages := []*Message{}
	for _, f := range e.Archive.GetFiles(prefix) {
		if !strings.HasSuffix(f.Name, "/") {
			m := make([]*Message, 0)
			err := e.UnmarshalFile(f.Name, &m)
			if err != nil {
				return nil, fmt.Errorf("error unmarshaling message from file %q: %w", f.Name, err)
			}
			messages = append(messages, m...)
		}
	}
	return messages, nil
}
