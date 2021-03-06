/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

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
	K3 = "KEY_3"
	M1 = "M1"
	M2 = "M2"
	M3 = "M3"
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
	s(t, M1, K3, "")
	s(t, M2, K1, "")
	s(t, M2, K2, "")
	s(t, M2, K3, "")
	s(t, M3, K1, "")
	s(t, M3, K2, "")
	s(t, M3, K3, "")
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

func TestLackDefaultValue(t *testing.T) {
	reset(t)
	m, ch := Main(M1)
	m.May(K1, "", "1")
	if !Parse() {
		t.Fatal("should be safe to omit with May")
	}
	x := <-ch
	res, ok := x[K1]
	if !ok {
		t.Fatal("there's no needed data in result")
	}
	if res != "1" {
		t.Fatal("unexpected result:", res)
	}
}

func TestSetDefaultValue(t *testing.T) {
	reset(t)
	s(t, M1, K1, "2")
	m, ch := Main(M1)
	m.May(K1, "", "1")
	if !Parse() {
		t.Fatal("should be safe to overwrite with May")
	}
	x := <-ch
	res, ok := x[K1]
	if !ok {
		t.Fatal("there's no needed data in result")
	}
	if res != "2" {
		t.Fatal("unexpected result:", res)
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

func TestToArr(t *testing.T) {
	cases := []struct {
		str    string
		expect []string
	}{
		{
			str:    "abc",
			expect: []string{"abc  "},
		},
		{
			str:    "abc\n123",
			expect: []string{"abc  ", "123  "},
		},
		{
			str:    "abc\n1234567",
			expect: []string{"abc  ", "12345", "67   "},
		},
		{
			str:    "abc\n12三4567",
			expect: []string{"abc  ", "12三4", "567  "},
		},
		{
			str:    "abc\n123四567",
			expect: []string{"abc  ", "123四", "567  "},
		},
		{
			str:    "abc\n1234五67",
			expect: []string{"abc  ", "1234 ", "五67 "},
		},
	}

	for idx, c := range cases {
		t.Run(fmt.Sprintf("#%d", idx), func(t *testing.T) {
			actual := toarr(c.str, 5)
			if !reflect.DeepEqual(actual, c.expect) {
				t.Fatalf("expected %+v, got %+v", c.expect, actual)
			}
		})
	}
}
