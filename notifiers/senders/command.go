// Statping
// Copyright (C) 2018.  Hunter Long and the project contributors
// Written by Hunter Long <info@socialeck.com> and the project contributors
//
// https://github.com/statping/statping
//
// The licenses for most software and other practical works are designed
// to take away your freedom to share and change the works.  By contrast,
// the GNU General Public License is intended to guarantee your freedom to
// share and change all versions of a program--to make sure it remains free
// software for all its users.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package senders

import (
	"fmt"
	"github.com/statping/statping/notifiers"
	"github.com/statping/statping/types/failures"
	"github.com/statping/statping/types/services"
	"github.com/statping/statping/utils"
	"strings"
	"time"
)

var _ notifiers.Notifier = (*commandLine)(nil)

type commandLine struct {
	*notifiers.Notification
}

var Command = &commandLine{&notifiers.Notification{
	Method:      "command",
	Title:       "Shell Command",
	Description: "Shell Command allows you to run a customized shell/bash Command on the local machine it's running on.",
	Author:      "Hunter Long",
	AuthorUrl:   "https://github.com/hunterlong",
	Delay:       time.Duration(1 * time.Second),
	Icon:        "fas fa-terminal",
	Host:        "/bin/bash",
	Form: []notifiers.NotificationForm{{
		Type:        "text",
		Title:       "Shell or Bash",
		Placeholder: "/bin/bash",
		DbField:     "host",
		SmallText:   "You can use '/bin/sh', '/bin/bash' or even an absolute path for an application.",
	}, {
		Type:        "text",
		Title:       "Command to Run on OnSuccess",
		Placeholder: "curl google.com",
		DbField:     "var1",
		SmallText:   "This Command will run every time a service is receiving a Successful event.",
	}, {
		Type:        "text",
		Title:       "Command to Run on OnFailure",
		Placeholder: "curl offline.com",
		DbField:     "var2",
		SmallText:   "This Command will run every time a service is receiving a Failing event.",
	}}},
}

func runCommand(app string, cmd ...string) (string, string, error) {
	outStr, errStr, err := utils.Command(app, cmd...)
	return outStr, errStr, err
}

func (u *commandLine) Select() *notifiers.Notification {
	return u.Notification
}

// OnFailure for commandLine will trigger failing service
func (u *commandLine) OnFailure(s *services.Service, f *failures.Failure) {
	u.AddQueue(fmt.Sprintf("service_%v", s.Id), u.Var2)
}

// OnSuccess for commandLine will trigger successful service
func (u *commandLine) OnSuccess(s *services.Service) {
	if !s.Online {
		u.ResetUniqueQueue(fmt.Sprintf("service_%v", s.Id))
		u.AddQueue(fmt.Sprintf("service_%v", s.Id), u.Var1)
	}
}

// OnSave for commandLine triggers when this notifier has been saved
func (u *commandLine) OnSave() error {
	u.AddQueue("saved", u.Var1)
	u.AddQueue("saved", u.Var2)
	return nil
}

// OnTest for commandLine triggers when this notifier has been saved
func (u *commandLine) OnTest() error {
	cmds := strings.Split(u.Var1, " ")
	in, out, err := runCommand(u.Host, cmds...)
	utils.Log.Infoln(in)
	utils.Log.Infoln(out)
	return err
}

// Send for commandLine will send message to expo Command push notifications endpoint
func (u *commandLine) Send(msg interface{}) error {
	cmd := msg.(string)
	_, _, err := runCommand(u.Host, cmd)
	return err
}