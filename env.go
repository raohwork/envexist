package envexist

import (
	"fmt"
	"os"
	"strings"
)

func tokey(m, k string) (ret string) {
	return m + "_" + k
}

type spec struct {
	desc     string
	example  string
	val      string
	required bool
	ok       bool
}

// Module represents a set of env vars
type Module struct {
	name   string
	params map[string]*spec
	cb     func(data map[string]string)
}

func (m *Module) notify() {
	data := map[string]string{}
	for n, s := range m.params {
		if s.val == "" {
			continue
		}
		data[n] = s.val
	}

	if m.cb != nil {
		m.cb(data)
	}
}

// Need registers a required env var
//
// It ALWAYS converts name to uppercase.
func (m *Module) Need(name string, desc, example string) (ret *Module) {
	name = strings.ToUpper(name)
	v := os.Getenv(tokey(m.name, name))
	x := &spec{
		desc:     desc,
		example:  example,
		val:      v,
		required: true,
		ok:       v != "",
	}
	m.params[name] = x
	return m
}

// Want registers an optional env var
//
// It ALWAYS converts name to uppercase.
func (m *Module) Want(name string, desc, example string) (ret *Module) {
	name = strings.ToUpper(name)
	x := &spec{
		desc:     desc,
		example:  example,
		val:      os.Getenv(tokey(m.name, name)),
		required: false,
		ok:       true,
	}
	m.params[name] = x
	return m
}

var (
	modules = []*Module{}
)

// New registers a new module suitable for libraries
//
// The env var will passed through data after calling Parse()
//
// It ALWAYS converts name to uppercase.
func New(name string, cb func(map[string]string)) (ret *Module) {
	ret = &Module{
		name:   strings.ToUpper(name),
		cb:     cb,
		params: map[string]*spec{},
	}

	modules = append(modules, ret)
	return
}

// Main registers a new module suitable for main application entry
//
// The env var will passed through data after calling Parse()
//
// It calls New() with a callback, which pushes data into ch. As ch is a buffered
// channel, Parse() can return before you retrieve the result from ch.
func Main(name string) (ret *Module, ch chan map[string]string) {
	ch = make(chan map[string]string, 1)
	ret = New(name, func(data map[string]string) {
		ch <- data
		close(ch)
	})

	return
}

// Parse checks all required env vars are set, and pass final data through channel
//
// It blocks until first error or all callbacks executed, that's why you should not
// use New() in application entry.
//
// Env vars are passed only if this returns true.
//
// YOU MUST NOT CALL New() or Main() AFTER THIS.
func Parse() (ok bool) {
	// check if there's something not set
	for _, m := range modules {
		for _, s := range m.params {
			if !s.ok {
				return
			}
		}
	}

	for _, m := range modules {
		m.notify()
	}

	return true
}

// Release releases used resources
//
// YOU SHOULD NOT CALL Parse() or PrintEnvList() AFTER THIS.
func Release() {
	modules = nil
}

// PrintEnvList lists all registered env vars with fmt.Print
//
// TODO: customizable formatting and output writer
func PrintEnvList() {
	// header
	fmt.Print("+----------------------+----------------------+---------------------------+-----------------+\n")
	fmt.Printf("| %-20s | %-20s | %-25s | %-15s |\n", "Name", "Value", "Description", "Example")
	fmt.Print("+----------------------+----------------------+---------------------------+-----------------+\n")

	for _, m := range modules {
		for n, s := range m.params {
			dumpSpec(m.name, n, s)
		}
	}
}

func stripNewline(orig string, l int) (ret string) {
	idx := strings.Index(orig, "\n")
	for idx != -1 {
		pos := idx % l
		pos = l - pos
		orig = strings.Replace(orig, "\n", strings.Repeat(" ", pos), 1)
		idx = strings.Index(orig, "\n")
	}

	return orig
}

func cut(orig string, l int, c bool) (res, rest string, cutted bool) {
	if len(orig) <= l {
		return orig, "", c
	}
	return orig[:l], orig[l:], true
}

func dumpSpec(m, n string, s *spec) {
	defer fmt.Print("+----------------------+----------------------+---------------------------+-----------------+\n")

	name := stripNewline(tokey(m, n), 20)
	val := stripNewline(s.val, 20)
	desc := stripNewline(s.desc, 25)
	ex := stripNewline(s.example, 15)
	cutted := true
	req := " "
	if s.required {
		req = "*"
	}

	for cutted {
		var n, v, d, e string
		c := false
		n, name, c = cut(name, 20, c)
		v, val, c = cut(val, 20, c)
		d, desc, c = cut(desc, 25, c)
		e, ex, c = cut(ex, 15, c)
		cutted = c

		fmt.Printf("|%s%-20s | %-20s | %-25s | %-15s |\n", req, n, v, d, e)
		req = " "
	}
}
