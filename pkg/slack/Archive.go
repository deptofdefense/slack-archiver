// =================================================================
//
// Work of the U.S. Department of Defense, Defense Digital Service.
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package slack

import (
	"archive/zip"
	"fmt"
	"os"
	"strings"

	"github.com/deptofdefense/slack-archiver/pkg/ziputil"
)

type Archive struct {
	name   string
	file   *os.File
	reader *zip.Reader
}

func (a *Archive) init() error {
	f, err := os.Open(a.name)
	if err != nil {
		return fmt.Errorf("error reading source %q: %w", a.name, err)
	}

	fi, err := f.Stat()
	if err != nil {
		_ = f.Close()
		return fmt.Errorf("error stating source %q: %w", a.name, err)
	}

	zr, err := zip.NewReader(f, fi.Size())
	if err != nil {
		_ = f.Close()
		return fmt.Errorf("error opening zip reader for source %q: %w", a.name, err)
	}
	a.file = f
	a.reader = zr
	return nil
}

func (a *Archive) UnmarshalFile(name string, v interface{}) error {
	err := ziputil.UnmarshalFile(a.reader, name, v)
	if err != nil {
		return fmt.Errorf("error unmarshaling file from %q: %w", a.name, err)
	}
	return nil
}

func (a *Archive) Close() error {
	err := a.file.Close()
	if err != nil {
		return fmt.Errorf("error closing source %q: %w", a.name, err)
	}
	return nil
}

func (a *Archive) GetFiles(prefix string) []*zip.File {
	if len(prefix) == 0 {
		return a.reader.File
	}
	files := make([]*zip.File, 0)
	for _, f := range a.reader.File {
		if strings.HasPrefix(f.Name, prefix) {
			files = append(files, f)
		}
	}
	return files
}

func (a *Archive) GetEnterpriseGrid() (*EnterpriseGrid, error) {

	enterpriseGrid := &EnterpriseGrid{
		Archive: a,
	}

	// Direct Messages

	directMessages := []*DirectMessage{}
	err := a.UnmarshalFile("dms.json", &directMessages)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling direct messages from %q: %w", a.name, err)
	}
	enterpriseGrid.DirectMessages = directMessages

	// Organization Users

	organizationUsers := []*User{}
	err = a.UnmarshalFile("org_users.json", &organizationUsers)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling organization users from %q: %w", a.name, err)
	}
	enterpriseGrid.OrganizationUsers = organizationUsers

	// Multiparty Instant Messages

	multiPartyInstantMessages := []*MultiPartyInstantMessage{}
	err = a.UnmarshalFile("mpims.json", &multiPartyInstantMessages)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling multiparty instant messages from %q: %w", a.name, err)
	}
	enterpriseGrid.MultiPartyInstantMessages = multiPartyInstantMessages

	// Groups

	groups := []*Group{}
	err = a.UnmarshalFile("groups.json", &groups)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling groups from %q: %w", a.name, err)
	}
	enterpriseGrid.Groups = groups

	// Integration Logs

	integrationLogMessages := []*IntegrationLogMessage{}
	err = a.UnmarshalFile("integration_logs.json", &integrationLogMessages)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling integration log messages from %q: %w", a.name, err)
	}
	enterpriseGrid.IntegrationLogMessages = integrationLogMessages

	// Teams

	teams := []*Team{}

	for _, f := range a.GetFiles("teams/") {
		if f.Name != "teams/" && strings.HasSuffix(f.Name, "/") && strings.Count(f.Name, "/") == 2 {
			name := f.Name[len("teams/") : len("teams/")+strings.Index(f.Name[len("teams/"):], "/")]
			teamChannels := []*Channel{}
			err = a.UnmarshalFile(fmt.Sprintf("teams/%s/channels.json", name), &teamChannels)
			if err != nil {
				return nil, fmt.Errorf("error unmarshaling channels for team %q from %q: %w", name, a.name, err)
			}
			teamGroups := []*Group{}
			err = a.UnmarshalFile(fmt.Sprintf("teams/%s/groups.json", name), &teamGroups)
			if err != nil {
				return nil, fmt.Errorf("error unmarshaling groups for team %q from %q: %w", name, a.name, err)
			}
			teamUsers := []*User{}
			err = a.UnmarshalFile(fmt.Sprintf("teams/%s/users.json", name), &teamUsers)
			if err != nil {
				return nil, fmt.Errorf("error unmarshaling users for team %q from %q: %w", name, a.name, err)
			}
			teams = append(teams, &Team{
				Name:     name,
				Channels: teamChannels,
				Groups:   teamGroups,
				Users:    teamUsers,
			})
		}
	}
	enterpriseGrid.Teams = teams

	return enterpriseGrid, nil
}

func OpenArchive(name string) (*Archive, error) {
	a := &Archive{
		name: name,
	}
	err := a.init()
	if err != nil {
		return nil, fmt.Errorf("error initializing archive: %w", err)
	}
	return a, nil
}
