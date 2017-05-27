/**

Copyright (C) 2017  ZhiQiang Huang (email: ngxpkg@gmail.com)

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.

**/

package main

import "fmt"

var (
	version = "0.0.1"
	osarch  string // set by ldflags

	cmdVersion = &command{
		run:       runVersion,
		UsageLine: "version",
		Short:     "display ngx version",
		Long:      "display ngx version and build info.\n",
	}
)

func init() {
	commands = append([]*command{cmdVersion}, commands...)
}

func runVersion(args []string) {
	fmt.Println(version, osarch)
}
