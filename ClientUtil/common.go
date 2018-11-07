package ClientUtil

import (
	"github.com/beego/bee/logger/colors"
	"fmt"
	"strings"
	"go-min-chat/ClientApp"
)

func GetPre() string {
	cliSing := ClientApp.GetCli()
	var pre string
	if (cliSing.RoomId != 0) {
		room := colors.Red(fmt.Sprintf("%s", cliSing.RoomName))
		pre = fmt.Sprintf("%s%s 我: ", strings.TrimSuffix(pre, " "), room)
	}
	return pre
}
