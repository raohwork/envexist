package envexist

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"golang.org/x/text/width"
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
		// sort keys
		keys := make([]string, 0, len(m.params))
		for n, _ := range m.params {
			keys = append(keys, n)
		}
		sort.Sort(sort.StringSlice(keys))

		for _, n := range keys {
			s := m.params[n]
			dumpSpec(m.name, n, s)
		}
	}
}

func toarr(orig string, sz int) (ret []string) {
	x := strings.Split(orig, "\n")
	ret = make([]string, 0, len(x))
	for _, s := range x {
		arr := []rune(s)
		l := 0
		buf := []rune{}
		for _, a := range arr {
			p := width.LookupRune(a)
			w := p.Kind()
			cur := 1
			if w == width.EastAsianWide || w == width.EastAsianFullwidth {
				cur += 1
			}
			if l+cur > sz {
				// CJK double width
				delta := sz - l
				if delta > 0 {
					tmp := make([]rune, delta)
					for idx, _ := range tmp {
						tmp[idx] = ' '
					}
					buf = append(buf, tmp...)
				}
				ret = append(ret, string(buf))
				l = cur
				buf = []rune{a}
				continue
			}

			l += cur
			buf = append(buf, a)
			if l == sz {
				ret = append(ret, string(buf))
				l = 0
				buf = []rune{}
				continue
			}
		}

		if l != 0 {
			delta := sz - l
			if delta > 0 {
				tmp := make([]rune, delta)
				for idx, _ := range tmp {
					tmp[idx] = ' '
				}
				buf = append(buf, tmp...)
			}
			ret = append(ret, string(buf))
		}
	}

	return
}

func max(arr []string, l int) (ret int) {
	if x := len(arr); x > l {
		return x
	}
	return l
}

func pad(arr []string, l, sz int) (ret []string) {
	delta := l - len(arr)
	if delta < 1 {
		return arr
	}

	empty := strings.Repeat(" ", sz)
	x := make([]string, delta)
	for idx, _ := range x {
		x[idx] = empty
	}
	ret = append(arr, x...)
	return
}

func dumpSpec(m, n string, s *spec) {
	defer fmt.Print("+----------------------+----------------------+---------------------------+-----------------+\n")

	l := 0
	name := toarr(tokey(m, n), 20)
	l = max(name, l)
	val := toarr(s.val, 20)
	l = max(val, l)
	desc := toarr(s.desc, 25)
	l = max(desc, l)
	ex := toarr(s.example, 15)
	l = max(ex, l)

	name = pad(name, l, 20)
	val = pad(val, l, 20)
	desc = pad(desc, l, 25)
	ex = pad(ex, l, 15)

	req := " "
	if s.required {
		req = "*"
	}

	for x := 0; x < l; x++ {
		fmt.Printf(
			"|%s%s | %s | %s | %s |\n",
			req,
			name[x],
			val[x],
			desc[x],
			ex[x],
		)
		req = " "
	}
}
