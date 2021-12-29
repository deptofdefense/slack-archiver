// =================================================================
//
// Work of the U.S. Department of Defense, Defense Digital Service.
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/deptofdefense/slack-archiver/pkg/slack"
)

const (
	SlackArchiverVersion = "1.0.0"
)

const (
	FlagSource      = "src"
	FlagDestination = "dest"
	FlagOverwrite   = "overwrite"
	FlagVersion     = "version"
)

func initListFlags(flag *pflag.FlagSet) {
	flag.StringP(FlagSource, "s", "", "path to Slack zip file")
	flag.BoolP(FlagVersion, "v", false, "show version")
}

func initDownloadFilesFlags(flag *pflag.FlagSet) {
	flag.String(FlagSource, "", "path to Slack zip file")
	flag.String(FlagDestination, "", "path to where to download files")
	flag.Bool(FlagOverwrite, false, "overwrite existing files")
	flag.BoolP(FlagVersion, "v", false, "show version")
}

func initViper(cmd *cobra.Command) (*viper.Viper, error) {
	v := viper.New()
	err := v.BindPFlags(cmd.Flags())
	if err != nil {
		return v, fmt.Errorf("error binding flag set to viper: %w", err)
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv() // set environment variables to overwrite config
	return v, nil
}

func checkConfig(v *viper.Viper) error {
	src := v.GetString(FlagSource)
	if len(src) == 0 {
		return fmt.Errorf("src is missing")
	}
	return nil
}

func checkDownloadFilesConfig(v *viper.Viper) error {
	src := v.GetString(FlagSource)
	if len(src) == 0 {
		return fmt.Errorf("src is missing")
	}
	dest := v.GetString(FlagDestination)
	if len(dest) == 0 {
		return fmt.Errorf("dest is missing")
	}
	return nil
}

func main() {

	rootCommand := &cobra.Command{
		Use:                   `slack-archiver [flags]`,
		DisableFlagsInUseLine: true,
		Short:                 "slack-archiver is a tool to archive a Slack enterprise grid.",
		Long:                  "slack-archiver is a tool to archive a Slack enterprise grid.",
	}

	listCommand := &cobra.Command{
		Use:                   `list`,
		DisableFlagsInUseLine: true,
		Short:                 "list data",
		SilenceErrors:         true,
		SilenceUsage:          true,
	}

	listFilesCommand := &cobra.Command{
		Use:                   `files [flags]`,
		DisableFlagsInUseLine: true,
		Short:                 "list files",
		Long:                  "list files",
		SilenceErrors:         true,
		SilenceUsage:          true,
		RunE: func(cmd *cobra.Command, args []string) error {
			v, err := initViper(cmd)
			if err != nil {
				return fmt.Errorf("error initializing viper: %w", err)
			}

			if len(args) > 0 {
				return cmd.Usage()
			}

			if v.GetBool(FlagVersion) {
				fmt.Println(SlackArchiverVersion)
				return nil
			}

			if errConfig := checkConfig(v); errConfig != nil {
				return errConfig
			}

			src := v.GetString(FlagSource)

			archive, err := slack.OpenArchive(src)
			if err != nil {
				return fmt.Errorf("error reading source %q: %w", src, err)
			}

			enterpriseGrid, err := archive.GetEnterpriseGrid()
			if err != nil {
				return fmt.Errorf("error reading enterprise grid from %q: %w", src, err)
			}

			encoder := json.NewEncoder(os.Stdout)

			// Files from multiparty instant messages

			for i, mpim := range enterpriseGrid.GetMultiPartyInstantMessages() {
				messages, getMessagesError := enterpriseGrid.GetMessages(fmt.Sprintf("%s/", mpim.Name))
				if getMessagesError != nil {
					return fmt.Errorf(
						"error reading mulitparty instant messages %d from %q for mpim %q : %w",
						i,
						src,
						mpim.Name,
						getMessagesError,
					)
				}
				for _, msg := range messages {
					if files := msg.Files; len(files) > 0 {
						for k, file := range files {
							encodeError := encoder.Encode(file)
							if encodeError != nil {
								return fmt.Errorf(
									"error encoding file from %q: %w",
									fmt.Sprintf("%s/%s/%d", mpim.Name, file.Name, k),
									encodeError,
								)
							}
						}
					}
				}
			}

			// Files from direct messages

			for i, dm := range enterpriseGrid.GetDirectMessages() {
				messages, getMessagesError := enterpriseGrid.GetMessages(fmt.Sprintf("%s/", dm.ID))
				if getMessagesError != nil {
					return fmt.Errorf(
						"error reading direct message %d from %q for mpim %q : %w",
						i,
						src,
						dm.ID,
						getMessagesError,
					)
				}
				for _, msg := range messages {
					if files := msg.Files; len(files) > 0 {
						for k, file := range files {
							encodeError := encoder.Encode(file)
							if encodeError != nil {
								return fmt.Errorf(
									"error encoding file from %q: %w",
									fmt.Sprintf("%s/%s/%d", dm.ID, file.Name, k),
									encodeError,
								)
							}
						}
					}
				}
			}

			// Files from teams

			for _, t := range enterpriseGrid.GetTeams() {

				// Files from channels

				for _, c := range t.Channels {
					messages, getMessagesError := enterpriseGrid.GetMessages(fmt.Sprintf("teams/%s/%s/", t.Name, c.Name))
					if getMessagesError != nil {
						return fmt.Errorf(
							"error reading messages for channel %q in team %q from %q: %w",
							c.Name,
							t.Name,
							src,
							getMessagesError,
						)
					}
					for _, msg := range messages {
						if files := msg.Files; len(files) > 0 {
							for k, file := range files {
								encodeError := encoder.Encode(file)
								if encodeError != nil {
									return fmt.Errorf(
										"error encoding file from %q: %w",
										fmt.Sprintf("teams/%s/%s/%s/%d", t.Name, c.Name, file.Name, k),
										encodeError,
									)
								}
							}
						}
					}
				}

				// Files from groups

				for _, g := range t.Groups {
					messages, getMessagesError := enterpriseGrid.GetMessages(fmt.Sprintf("teams/%s/%s/", t.Name, g.Name))
					if getMessagesError != nil {
						return fmt.Errorf(
							"error reading messages for group %q in team %q from %q: %w",
							g.Name,
							t.Name,
							src,
							getMessagesError,
						)
					}
					for _, msg := range messages {
						if files := msg.Files; len(files) > 0 {
							for k, file := range files {
								encodeError := encoder.Encode(file)
								if encodeError != nil {
									return fmt.Errorf(
										"error encoding file from %q: %w",
										fmt.Sprintf("teams/%s/%s/%s/%d", t.Name, g.Name, file.Name, k),
										encodeError,
									)
								}
							}
						}
					}
				}
			}

			err = archive.Close()
			if err != nil {
				return fmt.Errorf("error closing file for source %q: %w", src, err)
			}
			return nil
		},
	}
	initListFlags(listFilesCommand.Flags())

	listTeamsCommand := &cobra.Command{
		Use:                   `teams [flags]`,
		DisableFlagsInUseLine: true,
		Short:                 "list teams",
		Long:                  "list teams",
		SilenceErrors:         true,
		SilenceUsage:          true,
		RunE: func(cmd *cobra.Command, args []string) error {
			v, err := initViper(cmd)
			if err != nil {
				return fmt.Errorf("error initializing viper: %w", err)
			}

			if len(args) > 0 {
				return cmd.Usage()
			}

			if v.GetBool(FlagVersion) {
				fmt.Println(SlackArchiverVersion)
				return nil
			}

			if errConfig := checkConfig(v); errConfig != nil {
				return errConfig
			}

			src := v.GetString(FlagSource)

			archive, err := slack.OpenArchive(src)
			if err != nil {
				return fmt.Errorf("error reading source %q: %w", src, err)
			}

			enterpriseGrid, err := archive.GetEnterpriseGrid()
			if err != nil {
				return fmt.Errorf("error reading enterprise grid from %q: %w", src, err)
			}

			teams := enterpriseGrid.GetTeams()

			encoder := json.NewEncoder(os.Stdout)
			for _, team := range teams {
				encodeError := encoder.Encode(team)
				if encodeError != nil {
					return fmt.Errorf(
						"error encoding team %q from %q: %w",
						team.Name,
						src,
						encodeError,
					)
				}
			}

			err = archive.Close()
			if err != nil {
				return fmt.Errorf("error closing file for source %q: %w", src, err)
			}
			return nil
		},
	}
	initListFlags(listTeamsCommand.Flags())

	listCommand.AddCommand(listFilesCommand, listTeamsCommand)

	downloadCommand := &cobra.Command{
		Use:                   `download`,
		DisableFlagsInUseLine: true,
		Short:                 "download data",
		SilenceErrors:         true,
		SilenceUsage:          true,
	}

	downloadFilesCommand := &cobra.Command{
		Use:                   `files [flags]`,
		DisableFlagsInUseLine: true,
		Short:                 "download files",
		Long:                  "download files",
		SilenceErrors:         true,
		SilenceUsage:          true,
		RunE: func(cmd *cobra.Command, args []string) error {
			v, err := initViper(cmd)
			if err != nil {
				return fmt.Errorf("error initializing viper: %w", err)
			}

			if len(args) > 0 {
				return cmd.Usage()
			}

			if v.GetBool(FlagVersion) {
				fmt.Println(SlackArchiverVersion)
				return nil
			}

			if errConfig := checkDownloadFilesConfig(v); errConfig != nil {
				return errConfig
			}

			src := v.GetString(FlagSource)
			dest := v.GetString(FlagDestination)
			overwrite := v.GetBool(FlagOverwrite)

			archive, err := slack.OpenArchive(src)
			if err != nil {
				return fmt.Errorf("error reading source %q: %w", src, err)
			}

			enterpriseGrid, err := archive.GetEnterpriseGrid()
			if err != nil {
				return fmt.Errorf("error reading enterprise grid from %q: %w", src, err)
			}

			allFiles := make([]slack.MessageFile, 0)

			// Files from multiparty instant messages

			for i, mpim := range enterpriseGrid.GetMultiPartyInstantMessages() {
				messages, getMessagesError := enterpriseGrid.GetMessages(fmt.Sprintf("%s/", mpim.Name))
				if getMessagesError != nil {
					return fmt.Errorf(
						"error reading mulitparty instant messages %d from %q for mpim %q : %w",
						i,
						src,
						mpim.Name,
						getMessagesError,
					)
				}
				for _, msg := range messages {
					if files := msg.Files; len(files) > 0 {
						allFiles = append(allFiles, files...)
					}
				}
			}

			// Files from direct messages

			for i, dm := range enterpriseGrid.GetDirectMessages() {
				messages, getMessagesError := enterpriseGrid.GetMessages(fmt.Sprintf("%s/", dm.ID))
				if getMessagesError != nil {
					return fmt.Errorf(
						"error reading direct message %d from %q for mpim %q : %w",
						i,
						src,
						dm.ID,
						getMessagesError,
					)
				}
				for _, msg := range messages {
					if files := msg.Files; len(files) > 0 {
						allFiles = append(allFiles, files...)
					}
				}
			}

			// Files from teams

			for _, t := range enterpriseGrid.GetTeams() {

				// Files from channels

				for _, c := range t.Channels {
					messages, getMessagesError := enterpriseGrid.GetMessages(fmt.Sprintf("teams/%s/%s/", t.Name, c.Name))
					if getMessagesError != nil {
						return fmt.Errorf(
							"error reading messages for channel %q in team %q from %q: %w",
							c.Name,
							t.Name,
							src,
							getMessagesError,
						)
					}
					for _, msg := range messages {
						if files := msg.Files; len(files) > 0 {
							allFiles = append(allFiles, files...)
						}
					}
				}

				// Files from groups

				for _, g := range t.Groups {
					messages, getMessagesError := enterpriseGrid.GetMessages(fmt.Sprintf("teams/%s/%s/", t.Name, g.Name))
					if getMessagesError != nil {
						return fmt.Errorf(
							"error reading messages for group %q in team %q from %q: %w",
							g.Name,
							t.Name,
							src,
							getMessagesError,
						)
					}
					for _, msg := range messages {
						if files := msg.Files; len(files) > 0 {
							allFiles = append(allFiles, files...)
						}
					}
				}
			}

			client := http.Client{
				CheckRedirect: func(r *http.Request, via []*http.Request) error {
					r.URL.Opaque = r.URL.Path
					return nil
				},
			}

			for _, f := range allFiles {
				if !f.IsTombstone() {

					u, parseError := url.Parse(f.URLPrivateDownload)
					if parseError != nil {
						return fmt.Errorf(
							"error parsing url for file %q: %w",
							f.ID,
							parseError,
						)
					}

					filename := u.Path[strings.LastIndex(u.Path, "/")+1 : len(u.Path)]

					created := time.Unix(int64(f.Created), 0)
					createdYear, createdMonth, createdDay := created.Date()

					fullpath := filepath.Join(
						dest,
						fmt.Sprintf("year=%d", createdYear),
						fmt.Sprintf("month=%d", int(createdMonth)),
						fmt.Sprintf("day=%d", createdDay),
						fmt.Sprintf("user=%s", f.User),
						fmt.Sprintf("filetype=%s", f.FileType),
						fmt.Sprintf("id=%s", f.ID),
						filename,
					)

					mkdirAllError := os.MkdirAll(filepath.Dir(fullpath), 0775)
					if err != nil {
						return fmt.Errorf("error creating directory for file %q from url %q: %w", f.ID, u.String(), mkdirAllError)
					}

					if fi, err := os.Stat(fullpath); err == nil {
						if !overwrite {
							if fi.Size() == f.Size {
								// if not overwriting and the sizes match, then skip
								continue
							}
						}
					}

					resp, getError := client.Get(u.String())
					if getError != nil {
						return fmt.Errorf("error downloading file %q from url %q: %w", f.ID, u.String(), getError)
					}

					downloadedFile, createError := os.Create(fullpath)
					if createError != nil {
						return fmt.Errorf("error creating file %q for url %q: %w", fullpath, u.String(), createError)
					}

					_, copyError := io.Copy(downloadedFile, resp.Body)
					if copyError != nil {
						return fmt.Errorf("error copying file %q for url %q: %w", fullpath, u.String(), copyError)
					}

					_ = downloadedFile.Close()

					_ = resp.Body.Close()

				}
			}

			err = archive.Close()
			if err != nil {
				return fmt.Errorf("error closing file for source %q: %w", src, err)
			}
			return nil
		},
	}
	initDownloadFilesFlags(downloadFilesCommand.Flags())

	downloadCommand.AddCommand(downloadFilesCommand)

	versionCommand := &cobra.Command{
		Use:                   `version`,
		DisableFlagsInUseLine: true,
		Short:                 "show version",
		SilenceErrors:         true,
		SilenceUsage:          true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(SlackArchiverVersion)
			return nil
		},
	}

	rootCommand.AddCommand(listCommand, downloadCommand, versionCommand)

	if err := rootCommand.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "slack-archiver: "+err.Error())
		_, _ = fmt.Fprintln(os.Stderr, "Try slack-archiver --help for more information.")
		os.Exit(1)
	}
}
