// Copyright (c) 2018 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package confi

import (
	"errors"
	"flag"
	"path"
	"path/filepath"
	"sort"
)

// FileReader is equivalent to FlagReader(config, false).
func FileReader(config interface{}) flag.Value {
	return FlagReader(config, false)
}

// FlagReader makes a ``dynamic value'' which reads files into the
// configuration as it receives filenames.  Unknown keys are silently skipped
// if ignoreUnknown is true.
func FlagReader(config interface{}, ignoreUnknown bool) flag.Value {
	return fileReader{config, ignoreUnknown}
}

type fileReader struct {
	config        interface{}
	ignoreUnknown bool
}

func (fr fileReader) Set(filename string) error {
	return readFile(filename, fr.config, fr.ignoreUnknown)
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

func (a assigner) Set(expr string) error {
	return assign(a.config, expr, a.ignoreUnknown)
}

func (assigner) String() string {
	return ""
}

// Buffer configuration files and assignments.
type Buffer struct {
	list []buffered
}

// NewBuffer with optional default absolute configuration filenames or glob patterns.
func NewBuffer(defaults ...string) *Buffer {
	b := new(Buffer)
	for _, pattern := range defaults {
		if pattern != "" {
			b.list = append(b.list, buffered{
				pattern: pattern,
			})
		}
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

// DirReader makes a ``dynamic value'' which buffers directories to read files
// from.
func (b *Buffer) DirReader(pattern string) flag.Value {
	if pattern == "" {
		panic("a glob pattern must be specified")
	}
	return bufferedDirReader{b, pattern}
}

// Assigner makes a ``dynamic value'' which buffers assignment expressions to
// be applied.
func (b *Buffer) Assigner() flag.Value {
	return bufferedAssigner{b}
}

// Apply is equivalent to Flush(config, false).
func (b Buffer) Apply(config interface{}) error {
	return b.Flush(config, false)
}

// Flush files and assignments to the configuration.  Unknown keys are silently
// skipped if ignoreUnknown is true.
func (b Buffer) Flush(config interface{}, ignoreUnknown bool) error {
	for _, entry := range b.list {
		if err := entry.flush(config, ignoreUnknown); err != nil {
			return err
		}
	}
	return nil
}

type buffered struct {
	filename string
	pattern  string
	expr     string
}

func (b buffered) flush(config interface{}, ignoreUnknown bool) error {
	switch {
	case b.filename != "":
		return readFile(b.filename, config, ignoreUnknown)

	case b.pattern != "":
		names, err := filepath.Glob(b.pattern)
		if err != nil {
			return err
		}
		sort.Strings(names)
		for _, name := range names {
			if err := readFileIfExists(name, config, ignoreUnknown); err != nil {
				return err
			}
		}
		return nil

	case b.expr != "":
		return assign(config, b.expr, ignoreUnknown)

	default:
		panic(b)
	}
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

type bufferedDirReader struct {
	b       *Buffer
	pattern string
}

func (dr bufferedDirReader) Set(dirname string) error {
	if dirname == "" {
		return errors.New("configuration directory name is empty")
	}
	dr.b.list = append(dr.b.list, buffered{pattern: path.Join(dirname, dr.pattern)})
	return nil
}

func (bufferedDirReader) String() string {
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
