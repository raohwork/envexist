package envexist

import "os"

func Example() {
	// reset envexist as we're in test environment, you should not need this
	Release()
	// simulate
	os.Setenv("MYMODULE_PARAM_2", "中文")

	m, ch := Main("MYMODULE")
	m.Need("PARAM_1", "a super detailed and descriptive decription which introduces how this variable should be and what it should do", "example value")
	m.Want("PARAM_2", "desc", "example")

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
}
