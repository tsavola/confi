// Copyright (c) 2018 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package confi

import (
	"errors"
	"flag"
)

// FileReader is equivalent to FlagReader(config, false).
func FileReader(config interface{}) flag.Value {
	return FlagReader(config, false)
}

// FlagReader makes a ``dynamic value'' which reads files into the
// configuration as it receives filenames.  Unknown keys are silently skipped
// if skipUnknown is true.
func FlagReader(config interface{}, skipUnknown bool) flag.Value {
	return fileReader{config, skipUnknown}
}

type fileReader struct {
	config      interface{}
	skipUnknown bool
}

func (fr fileReader) Set(filename string) error {
	return readFile(filename, fr.config, fr.skipUnknown)
}

func (fileReader) String() string {
	return ""
}

// Assigner is equivalent to FlagSetter(config, false).
func Assigner(config interface{}) flag.Value {
	return FlagSetter(config, false)
}

// FlagSetter makes a ``dynamic value'' which sets fields in the configuration
// as it receives assignment expressions.  An error is not returned for unknown
// keys if ignoreUnknown is true.
func FlagSetter(config interface{}, ignoreUnknown bool) flag.Value {
	return assigner{config, ignoreUnknown}
}

type assigner struct {
	config        interface{}
	ignoreUnknown bool
}

func (a assigner) Set(expr string) (err error) {
	err = Assign(a.config, expr)
	if err != nil && a.ignoreUnknown {
		if _, ok := err.(unknownKeyError); ok {
			err = nil
		}
	}
	return
}

func (assigner) String() string {
	return ""
}

// Buffer configuration files and assignments.
type Buffer struct {
	list []buffered
}

func NewBuffer(optionalDefaultFilename string) *Buffer {
	b := new(Buffer)
	if optionalDefaultFilename != "" {
		b.list = append(b.list, buffered{
			filename: optionalDefaultFilename,
			optional: true,
		})
	}
	return b
}

// FileReplacer makes a ``dynamic value'' which buffers file names to be read.
// It discards previously buffered values.
func (b *Buffer) FileReplacer() flag.Value {
	return bufferedFileReader{b, true}
}

// FileReader makes a ``dynamic value'' which buffers file names to be read.
func (b *Buffer) FileReader() flag.Value {
	return bufferedFileReader{b, false}
}

// Assigner makes a ``dynamic value'' which buffers assignment expressions to
// be applied.
func (b *Buffer) Assigner() flag.Value {
	return bufferedAssigner{b}
}

// Apply files and assignments to the configuration.
func (b Buffer) Apply(config interface{}) error {
	for _, entry := range b.list {
		if err := entry.apply(config); err != nil {
			return err
		}
	}
	return nil
}

type buffered struct {
	filename string
	expr     string
	optional bool
}

func (b buffered) apply(config interface{}) error {
	switch {
	case b.filename != "":
		if b.optional {
			return ReadFileIfExists(b.filename, config)
		} else {
			return ReadFile(b.filename, config)
		}

	case b.expr != "":
		return Assign(config, b.expr)
	}

	panic(b)
}

type bufferedFileReader struct {
	b       *Buffer
	replace bool
}

func (fr bufferedFileReader) Set(filename string) error {
	if filename == "" {
		return errors.New("configuration file name is empty")
	}
	if fr.replace {
		fr.b.list = nil
	}
	fr.b.list = append(fr.b.list, buffered{filename: filename})
	return nil
}

func (bufferedFileReader) String() string {
	return ""
}

type bufferedAssigner struct {
	b *Buffer
}

func (a bufferedAssigner) Set(expr string) error {
	if expr == "" {
		return errors.New("configuration expression is empty")
	}
	a.b.list = append(a.b.list, buffered{expr: expr})
	return nil
}

func (bufferedAssigner) String() string {
	return ""
}
