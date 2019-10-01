// Copyright (c) 2018 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package confi

import (
	"reflect"
	"testing"
	"time"
)

type testConfig struct {
	Foo struct {
		Key1  bool
		Key2  int
		Key2b int8
		Key3a int16
		Key3  int32
		Key4  int64
		Key5  uint
		Key5b uint8
		Key6a uint16
		Key6  uint32
		Key7  uint64
		Key8  float32
		Key9  float64
		Key10 string
		Key11 []string
	}

	Bar int

	Baz struct {
		Quux     testConfigQuux
		Interval time.Duration

		TestConfigEmbed
		Embed1 struct {
			TestConfigEmbed
		}
		Embed2 struct {
			*TestConfigEmbed
		}
	}

	Dummy struct{}
	_     struct{}

	Ignore struct {
		A *int
		B *struct{}
		C func()
		*int
		_ *struct{}
		d func()
	}

	Ext map[string]interface{}
}

type testConfigQuux struct {
	Key_a string
	Key_b bool
}

type TestConfigEmbed struct {
	EmBedded bool
}

type testExtA struct{ Average int }
type testExtB struct{ Beverage int }

var testConfigTOML = `bar = 12345

[baz]
embedded = false
interval = "10h9m8.007006005s"

[baz.embed1]
embedded = false

[baz.embed2]
embedded = false

[baz.quux]
key_a = "true"
key_b = true

[ext.a]
average = 123

[ext.b]
beverage = 456

[foo]
key1 = true
key10 = "hello, world"
key11 = ["hello", "world"]
key2 = -10
key2b = -128
key3 = -11
key3a = -32768
key4 = -100000000000000
key5 = 10
key5b = 255
key6 = 11
key6a = 65535
key7 = 100000000000000
key8 = 1.5e+00
key9 = 1.0000000000005e+00
`

func newTestConfig() *testConfig {
	c := new(testConfig)
	c.Bar = 67890
	c.Baz.Embed2.TestConfigEmbed = new(TestConfigEmbed)
	c.Ext = map[string]interface{}{
		"a": new(testExtA),
		"b": new(testExtB),
	}
	return c
}

func testConfigValues(t *testing.T, c *testConfig) {
	if !c.Foo.Key1 {
		t.Error(c.Foo.Key1)
	}
	if c.Foo.Key2 != -10 {
		t.Error(c.Foo.Key2)
	}
	if c.Foo.Key2b != -128 {
		t.Error(c.Foo.Key2b)
	}
	if c.Foo.Key3a != -32768 {
		t.Error(c.Foo.Key3a)
	}
	if c.Foo.Key3 != -11 {
		t.Error(c.Foo.Key3)
	}
	if c.Foo.Key4 != -100000000000000 {
		t.Error(c.Foo.Key4)
	}
	if c.Foo.Key5 != 10 {
		t.Error(c.Foo.Key5)
	}
	if c.Foo.Key5b != 255 {
		t.Error(c.Foo.Key5b)
	}
	if c.Foo.Key6a != 65535 {
		t.Error(c.Foo.Key6a)
	}
	if c.Foo.Key6 != 11 {
		t.Error(c.Foo.Key6)
	}
	if c.Foo.Key7 != 100000000000000 {
		t.Error(c.Foo.Key7)
	}
	if c.Foo.Key8 != 1.5 {
		t.Error(c.Foo.Key8)
	}
	if c.Foo.Key9 != 1.0000000000005 {
		t.Error(c.Foo.Key9)
	}
	if c.Foo.Key10 != "hello, world" {
		t.Error(c.Foo.Key10)
	}
	if !reflect.DeepEqual(c.Foo.Key11, []string{"hello", "world"}) {
		t.Error()
	}
	if c.Bar != 12345 {
		t.Error(c.Bar != 12345)
	}
	if c.Baz.Quux.Key_a != "true" {
		t.Error(c.Baz.Quux.Key_a)
	}
	if !c.Baz.Quux.Key_b {
		t.Error(c.Baz.Quux.Key_b)
	}
	if c.Baz.EmBedded {
		t.Error(c.Baz.EmBedded)
	}
	if c.Baz.Embed1.EmBedded {
		t.Error(c.Baz.Embed1.EmBedded)
	}
	if c.Baz.Embed2.EmBedded {
		t.Error(c.Baz.Embed2.EmBedded)
	}
	if c.Baz.Interval != 10*time.Hour+9*time.Minute+8*time.Second+7*time.Millisecond+6*time.Microsecond+5*time.Nanosecond {
		t.Error(c.Baz.Interval)
	}
	if a, ok := c.Ext["a"].(*testExtA); ok {
		if a.Average != 123 {
			t.Error(a.Average)
		}
	} else {
		t.Error(c.Ext["a"])
	}
	if b, ok := c.Ext["b"].(*testExtB); ok {
		if b.Beverage != 456 {
			t.Error(b.Beverage)
		}
	} else {
		t.Error(c.Ext["b"])
	}
}
