package envexist

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

const (
	K1 = "KEY_1"
	K2 = "KEY_2"
	M1 = "M1"
	M2 = "M2"
)

func s(t *testing.T, m, k, v string) {
	k = tokey(m, k)
	if err := os.Setenv(k, v); err != nil {
		t.Fatalf("cannot update env var [%s] to [%s]: %s", k, v, err)
	}

}

func reset(t *testing.T) {
	s(t, M1, K1, "")
	s(t, M1, K2, "")
	s(t, M2, K1, "")
	s(t, M2, K2, "")
	Release()
}

func TestAllSet(t *testing.T) {
	reset(t)
	s(t, M1, K1, "1")
	s(t, M1, K2, "2")

	m := New(M1, nil)
	m.Need(K1, "", "")
	m.Want(K2, "", "")
	if !Parse() {
		t.Fatal("should be true as M1_K1 has value")
	}
}

func TestOnlyNeeded(t *testing.T) {
	reset(t)
	s(t, M1, K1, "1")
	m := New(M1, nil)
	m.Need(K1, "", "")
	m.Want(K2, "", "")
	if !Parse() {
		t.Fatal("should be true as M1_K1 has value")
	}
}

func TestLackNeeded(t *testing.T) {
	reset(t)
	s(t, M1, K2, "1")
	m := New(M1, nil)
	m.Need(K1, "", "")
	m.Want(K2, "", "")
	if Parse() {
		t.Fatal("should be false as M1_K1 is empty")
	}
}

func TestAllEmpty(t *testing.T) {
	reset(t)
	m := New(M1, nil)
	m.Need(K1, "", "")
	m.Want(K2, "", "")
	if Parse() {
		t.Fatal("should be false as M1_K1 is empty")
	}
}

func TestOneModuleFailed(t *testing.T) {
	reset(t)
	s(t, M1, K1, "1")
	m := New(M1, nil)
	m.Need(K1, "", "")
	m.Want(K2, "", "")
	m = New(M2, nil)
	m.Need(K1, "", "")
	m.Want(K2, "", "")
	if Parse() {
		t.Fatal("should be false as M2_K1 is empty")
	}
}

type datacase struct {
	keys   map[string]bool
	vals   map[string]string
	expect map[string]string
}

var cases = []datacase{
	{
		keys:   map[string]bool{K1: true},
		vals:   map[string]string{K1: "1"},
		expect: nil,
	},
	{
		keys:   map[string]bool{K1: true, K2: false},
		vals:   map[string]string{K1: "1"},
		expect: nil,
	},
	{
		keys:   map[string]bool{K1: true},
		vals:   map[string]string{K1: "1", K2: "2"},
		expect: map[string]string{K1: "1"},
	},
	{
		keys:   map[string]bool{K1: true, K2: false},
		vals:   map[string]string{K1: "1", K2: "2"},
		expect: nil,
	},
}

func initdatacase(t *testing.T, c datacase) (exp map[string]string) {
	reset(t)
	for k, v := range c.vals {
		s(t, M1, k, v)
	}
	exp = c.expect
	if exp == nil {
		exp = c.vals
	}
	return
}

func TestDataInNew(t *testing.T) {
	for idx, c := range cases {
		t.Run(fmt.Sprintf("#%d", idx), func(t *testing.T) {
			exp := initdatacase(t, c)

			var actual map[string]string
			m := New(M1, func(data map[string]string) {
				actual = data
			})
			for k, need := range c.keys {
				f := m.Need
				if !need {
					f = m.Want
				}
				f(k, "", "")
			}

			if !Parse() {
				t.Fatal("should be true as all required vars are set")
			}
			if !reflect.DeepEqual(exp, actual) {
				t.Fatalf("expected %+v, got %+v", exp, actual)
			}
		})
	}
}

func TestDataInMain(t *testing.T) {
	for idx, c := range cases {
		t.Run(fmt.Sprintf("#%d", idx), func(t *testing.T) {
			exp := initdatacase(t, c)

			m, ch := Main(M1)
			for k, need := range c.keys {
				f := m.Need
				if !need {
					f = m.Want
				}
				f(k, "", "")
			}

			if !Parse() {
				t.Fatal("should be true as all required vars are set")
			}
			actual := <-ch
			if !reflect.DeepEqual(exp, actual) {
				t.Fatalf("expected %+v, got %+v", exp, actual)
			}
		})
	}
}
