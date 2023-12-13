// Copyright 2022 PingCAP, Inc.
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
	"encoding/json"
	"fmt"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/ui"
	"tidbcloud-cli/internal/util"
	"tidbcloud-cli/pkg/tidbcloud/pingchat/client/operations"
	"tidbcloud-cli/pkg/tidbcloud/pingchat/models"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pingcap/errors"
	"github.com/spf13/cobra"
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
			param := operations.NewChatParams()

			if opts.interactive {
				task := func(messages []ui.ChatMessage) tea.Msg {
					msgs := make([]*models.PingchatChatMessage, 0, len(messages))
					for _, message := range messages {
						content := message.Content
						role, err := convertRole(message.Role)
						if err != nil {
							return ui.EndSendingMsg{
								Err: err,
							}
						}
						msg := models.PingchatChatMessage{
							Content: &content,
							Role:    &role,
						}
						msgs = append(msgs, &msg)
					}
					chat, err := client.Chat(param.WithChatInfo(&models.PingchatChatInfo{
						Messages: msgs,
					}))

					if err != nil {
						return ui.EndSendingMsg{
							Err: err,
						}
					}

					content := chat.Payload.Content + "\n\n"
					for _, link := range chat.Payload.Links {
						content = fmt.Sprintf("%s[%s](%s)\n", content, link.Title, link.Link)
					}

					return ui.EndSendingMsg{
						Msg: ui.ChatMessage{
							Role:    ui.RoleBot,
							Content: content,
						},
					}
				}

				model := ui.InitialChatBoxModel(task, "TiDB Bot")
				p := tea.NewProgram(model,
					tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
					tea.WithMouseCellMotion(), // turn on mouse support we can track the mouse wheel
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

				role := models.PingchatChatMessageRoleUser
				chat, err := client.Chat(param.WithChatInfo(&models.PingchatChatInfo{
					Messages: []*models.PingchatChatMessage{
						{
							Content: &query,
							Role:    &role,
						},
					},
				}))

				if err != nil {
					return err
				}

				marshal, err := json.MarshalIndent(chat, "", "  ")
				if err != nil {
					return err
				}
				fmt.Fprintln(h.IOStreams.Out, string(marshal))
			}

			return nil
		},
	}

	cmd.Flags().StringP(flag.Query, flag.QueryShort, "", "The query to chat with TiDB Bot")
	return cmd
}

func convertRole(role ui.Role) (string, error) {
	switch role {
	case ui.RoleUser:
		return models.PingchatChatMessageRoleUser, nil
	case ui.RoleBot:
		return models.PingchatChatMessageRoleAssistant, nil
	default:
		return "", errors.New("unknown chat role")
	}
}
