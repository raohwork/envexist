/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package envexist

import "os"

func Example() {
	// reset envexist as we're in test environment, you should not need this
	Release()
	// simulate
	os.Setenv("MYMODULE_PARAM_2", "中文")

	m, ch := Main("MYMODULE")
	// there will be an asterisk before the Name column
	m.Need("PARAM_1", "a super detailed and descriptive decription which introduces how this variable should be and what it should do", "example value")

	// nothing special
	m.Want("PARAM_2", "desc", "example")

	// there will be an asterisk before the Example column
	m.May("PARAM_3", "the value in example works as default value", "default_Value")

	if !Parse() {
		PrintEnvList()
		return
	}

	// consume your data here
	env := <-ch
	param1 := env["PARAM_1"]
	_ = param1

	// output: +----------------------+----------------------+---------------------------+-----------------+
	// | Name                 | Value                | Description               | Example         |
	// +----------------------+----------------------+---------------------------+-----------------+
	// |*MYMODULE_PARAM_1     |                      | a super detailed and desc | example value   |
	// |                      |                      | riptive decription which  |                 |
	// |                      |                      | introduces how this varia |                 |
	// |                      |                      | ble should be and what it |                 |
	// |                      |                      |  should do                |                 |
	// +----------------------+----------------------+---------------------------+-----------------+
	// | MYMODULE_PARAM_2     | 中文                 | desc                      | example         |
	// +----------------------+----------------------+---------------------------+-----------------+
	// | MYMODULE_PARAM_3     |                      | the value in example work |*default_Value   |
	// |                      |                      | s as default value        |                 |
	// +----------------------+----------------------+---------------------------+-----------------+
}
