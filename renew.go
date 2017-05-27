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

var (
	cmdReNew = &command{
		run:       runReNew,
		UsageLine: "renew [domain ...]",
		Short:     "renew ssl certificates base on domain conf",
		Long: `
Parse domain conf and renew ssl certificates.
If not domain input, will parse all domain conf at NGX_SITE_CONFIG dir.
`,
	}
)

func runReNew(args []string) {
	if len(args) == 0 {
		fatalf("no domain specified")
	}
}
