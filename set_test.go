// Copyright (c) 2018 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package confi

import (
	"reflect"
	"testing"
	"time"
)

func TestSet(t *testing.T) {
	c := newTestConfig()
	c.Foo.Key2 = 67890

	if err := Set(c, "foo.key1", true); err != nil {
		t.Error(err)
	}
	if !c.Foo.Key1 {
		t.Error(c.Foo.Key1)
	}

	if err := Set(c, "foo.key2", true); err == nil {
		t.Error("foo.key2")
	}
	if c.Foo.Key2 != 67890 {
		t.Error(c.Foo.Key2)
	}

	if err := Set(c, "foo.key2b", 10); err == nil {
		t.Error("foo.key2b")
	}
	if c.Foo.Key2b != 0 {
		t.Error(c.Foo.Key2b)
	}

	if err := Set(c, "foo.key3", 10); err == nil {
		t.Error("foo.key3")
	}
	if c.Foo.Key3 != 0 {
		t.Error(c.Foo.Key3)
	}

	if err := Set(c, "foo.key3", int32(10)); err != nil {
		t.Error(err)
	}
	if c.Foo.Key3 != 10 {
		t.Error(c.Foo.Key3)
	}

	if err := Set(c, "foo.key4", 10); err == nil {
		t.Error("foo.key4")
	}
	if c.Foo.Key4 != 0 {
		t.Error(c.Foo.Key4)
	}

	if err := Set(c, "foo.key4", int64(10)); err != nil {
		t.Error(err)
	}
	if c.Foo.Key4 != 10 {
		t.Error(c.Foo.Key4)
	}

	if err := Set(c, "foo.key9", "Hello, World"); err == nil {
		t.Error("foo.key9")
	}
	if c.Foo.Key9 != 0 {
		t.Error(c.Foo.Key9)
	}

	if err := Set(c, "foo.key10", "Hello, World"); err != nil {
		t.Error(err)
	}
	if c.Foo.Key10 != "Hello, World" {
		t.Error(c.Foo.Key10)
	}

	if err := Set(c, "foo.key11", []string{"Hello", "World"}); err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(c.Foo.Key11, []string{"Hello", "World"}) {
		t.Error(c.Foo.Key11)
	}

	if err := Set(c, "baz.interval", int64(time.Second)); err == nil {
		t.Error("baz.interval")
	}
	if c.Baz.Interval != 0 {
		t.Error(c.Baz.Interval)
	}

	if err := Set(c, "baz.interval", time.Second); err != nil {
		t.Error(err)
	}
	if c.Baz.Interval != time.Second {
		t.Error(c.Baz.Interval)
	}

	if err := Set(c, "ext.a.average", 123); err != nil {
		t.Error(err)
	}
	if a, ok := c.Ext["a"].(*testExtA); ok {
		if a.Average != 123 {
			t.Error(a.Average)
		}
	} else {
		t.Error(c.Ext["a"])
	}

	if err := Set(c, "ext.b.beverage", 456); err != nil {
		t.Error(err)
	}
	if b, ok := c.Ext["b"].(*testExtB); ok {
		if b.Beverage != 456 {
			t.Error(b.Beverage)
		}
	} else {
		t.Error(c.Ext["b"])
	}
}

func TestGet(t *testing.T) {
	c := newTestConfig()

	if x, err := Get(c, "bar"); err != nil {
		t.Error(err)
	} else if x.(int) != 67890 {
		t.Error(x)
	}
}
