/*
 * @Description:
 * @Author: Allen
 * @Date: 2021-05-28 20:38:38
 * @LastEditors: Allen
 * @LastEditTime: 2021-05-28 20:38:40
 */
package version

import (
	"fmt"
	"runtime"
)

var (
	Version string
	BuildDate string
	BuildUser string
	Branch string
)

func PrintCLIVersion() {
	println(fmt.Sprintf(
		"version %s,build user %s,branch %s,build on %s,%s",
		Version,
		BuildUser,
		Branch,
		BuildDate,
		runtime.Version(),
	)) 
}