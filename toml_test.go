// Copyright (c) 2018 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package confi

import (
	"bytes"
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	c := newTestConfig()

	if err := Read(strings.NewReader(testConfigTOML), c); err != nil {
		t.Fatalf("%v", err)
	}

	testConfigValues(t, c)
}

func TestReadKeepDefaults(t *testing.T) {
	var c struct {
		Foo string
		Bar string
		Baz struct {
			A string
			B string
		}
	}
	c.Bar = "default value"
	c.Baz.A = "preserve me please"

	if err := Read(strings.NewReader(`foo = "hello"
[baz]
b = "goodbye"
`), &c); err != nil {
		t.Fatalf("%v", err)
	}

	if c.Foo != "hello" {
		t.Error(c.Foo)
	}
	if c.Bar != "default value" {
		t.Error(c.Bar)
	}
	if c.Baz.A != "preserve me please" {
		t.Error(c.Baz.A)
	}
	if c.Baz.B != "goodbye" {
		t.Error(c.Baz.B)
	}
}

func TestReadFileIfExists(t *testing.T) {
	if err := ReadFileIfExists("/nonexistent", nil); err != nil {
		t.Error(err)
	}

	if ReadFileIfExists("/etc/issue", nil) == nil {
		t.Fail()
	}
}

func TestWrite(t *testing.T) {
	c := newTestConfig()

	if err := Read(strings.NewReader(testConfigTOML), c); err != nil {
		t.Fatal(err)
	}

	b := new(bytes.Buffer)

	if err := Write(b, c); err != nil {
		t.Fatal(err)
	}

	if s := b.String(); s != testConfigTOML {
		t.Error(s)
	}
}
