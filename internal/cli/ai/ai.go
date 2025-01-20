// Copyright 2025 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ai

import (
	"fmt"
	"regexp"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/output"
	"github.com/tidbcloud/tidbcloud-cli/internal/ui"
	"github.com/tidbcloud/tidbcloud-cli/internal/util"
	"github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/pingchat"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pingcap/errors"
	"github.com/spf13/cobra"
)

const (
	regexPattern = `\[\^(\d+)\]`
)

var (
	re     = regexp.MustCompile(regexPattern)
	domain = []pingchat.PingchatChatInfoDomainInner{"tidbcloud"}
)

type AIOpts struct {
	interactive bool
}

func (o AIOpts) NonInteractiveFlags() []string {
	return []string{
		flag.Query,
	}
}

func AICmd(h *internal.Helper) *cobra.Command {
	opts := AIOpts{
		interactive: true,
	}

	cmd := &cobra.Command{
		Use:         "ai",
		Short:       "Chat with TiDB Bot",
		Args:        cobra.NoArgs,
		Annotations: make(map[string]string),
		Example: fmt.Sprintf(`  Chat with TiDB Bot in interactive mode:
  $ %[1]s ai

  Chat with TiDB Bot in non-interactive mode:
  $ %[1]s ai -q "How to create a cluster?"`,
			config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			flags := opts.NonInteractiveFlags()
			for _, fn := range flags {
				f := cmd.Flags().Lookup(fn)
				if f != nil && f.Changed {
					opts.interactive = false
				}
			}

			// mark required flags in non-interactive mode
			if !opts.interactive {
				for _, fn := range flags {
					err := cmd.MarkFlagRequired(fn)
					if err != nil {
						return errors.Trace(err)
					}
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := h.Client()
			if err != nil {
				return err
			}

			context := cmd.Context()
			if opts.interactive {
				task := func(messages []ui.ChatMessage) tea.Msg {
					msgs := make([]pingchat.PingchatChatMessage, 0, len(messages))
					for _, message := range messages {
						content := message.Content
						role, err := convertRole(message.Role)
						if err != nil {
							return ui.EndSendingMsg{
								Err: err,
							}
						}
						msg := pingchat.PingchatChatMessage{
							Content: content,
							Role:    role,
						}
						msgs = append(msgs, msg)
					}
					chat, err := client.Chat(context, &pingchat.PingchatChatInfo{
						Messages: msgs,
						Domain:   domain,
					})

					if err != nil {
						return ui.EndSendingMsg{
							Err: err,
						}
					}

					linkContent := "\n\n"
					for i, link := range chat.Links {
						linkContent = fmt.Sprintf("%s[%d] [%s](%s)\n", linkContent, i+1, *link.Title, *link.Link)
					}

					// Replace occurrences of [^\d+] with [\d+] for better user comprehension.
					content := re.ReplaceAllString(*chat.Content, "[$1]")

					return ui.EndSendingMsg{
						Msg: ui.ChatMessage{
							Role:        ui.RoleBot,
							Content:     content,
							LinkContent: linkContent,
						},
					}
				}

				model := ui.InitialChatBoxModel(task, "TiDB Bot")
				p := tea.NewProgram(model,
					tea.WithAltScreen(), // use the full size of the terminal in its "alternate screen buffer"
				)
				typeModel, err := p.Run()
				if err != nil {
					return err
				}
				if m, _ := typeModel.(ui.ChatBoxModel); m.Interrupted {
					return util.InterruptError
				}
				if m, _ := typeModel.(ui.ChatBoxModel); m.Err != nil {
					return m.Err
				}
			} else {
				query, err := cmd.Flags().GetString(flag.Query)
				if err != nil {
					return errors.Trace(err)
				}

				chat, err := client.Chat(context, &pingchat.PingchatChatInfo{
					Messages: []pingchat.PingchatChatMessage{
						{
							Content: query,
							Role:    pingchat.PINGCHATCHATMESSAGEROLE_USER,
						},
					},
					Domain: domain,
				})
				if err != nil {
					return err
				}

				err = output.PrintJson(h.IOStreams.Out, chat)
				return errors.Trace(err)
			}

			return nil
		},
	}

	cmd.Flags().StringP(flag.Query, flag.QueryShort, "", "The query to chat with TiDB Bot.")
	return cmd
}

func convertRole(role ui.Role) (pingchat.PingchatChatMessageRole, error) {
	switch role {
	case ui.RoleUser:
		return pingchat.PINGCHATCHATMESSAGEROLE_USER, nil
	case ui.RoleBot:
		return pingchat.PINGCHATCHATMESSAGEROLE_ASSISTANT, nil
	default:
		return "", errors.New("unknown chat role")
	}
}
