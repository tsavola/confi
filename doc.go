// Copyright (c) 2018 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*

Package confi is an ergonomic configuration parsing toolkit.  The schema is
declared using a struct type, and values can be read from TOML files or set via
command-line flags.

A pointer to a preallocated configuration object of a user-defined struct type
must be passed to all functions.  The type can have an arbitrary number of
nested structs (either embedded or through an initialized pointer).  Only
exported fields can be used.  The object can be initialized with default
values.

Slices of structs can be populated by appending TOML table arrays, or by
indexing on the command line.

Dynamically created subtrees are supported via map[string]interface{} nodes.
The map values must be struct pointers.

The field names are spelled in lower case in TOML files and on the
command-line.  The accessor functions and flag values use dotted paths to
identify the field, such as "audio.samplerate".

Supported field types are bool, int, int8, int16, int32, int64, uint, uint8,
uint16, uint32, uint64, float32, float64, string, []string, and time.Duration.

The Get method is provided for completeness; the intended way to access
configuration values is through direct struct field access.

Short example:

	c := &myConfig{}

	flag.Usage = confi.FlagUsage(nil, c)
	flag.Var(confi.FileReader(c), "f", "read config from TOML files")
	flag.Var(confi.Assigner(c), "c", "set config keys (path.to.key=value)")
	flag.Parse()

Longer example:

	package main

	import (
		"flag"
		"fmt"
		"log"

		"github.com/tsavola/confi"
	)

	type myConfig struct {
		Comment string

		Size struct {
			Width  uint32
			Height uint32
		}

		Audio struct {
			Enabled    bool
			SampleRate int
		}
	}

	func main() {
		c := new(myConfig)
		c.Size.Width = 640
		c.Size.Height = 480
		c.Audio.SampleRate = 44100

		if err := confi.ReadFileIfExists("defaults.toml", c); err != nil {
			log.Print(err)
		}

		if x, _ := confi.Get(c, "audio.samplerate"); x.(int) <= 0 {
			confi.MustSet(c, "audio.enabled", false)
		}

		dump := flag.Bool("dump", false, "create defaults.toml")
		flag.Var(confi.FileReader(c), "f", "read config from TOML files")
		flag.Var(confi.Assigner(c), "c", "set config keys (path.to.key=value)")
		flag.Usage = confi.FlagUsage(nil, c)
		flag.Parse()

		if *dump {
			if err := confi.WriteFile("defaults.toml", c); err != nil {
				log.Fatal(err)
			}
		}

		fmt.Printf("Comment is %q\n", c.Comment)
		fmt.Printf("Size is %dx%d\n", c.Size.Width, c.Size.Height)
		if c.Audio.Enabled {
			fmt.Printf("Sample rate is %d\n", c.Audio.SampleRate)
		}
	}

Example usage output:

	$ example -help
	Usage of example:
	  -c value
	    	set config keys (path.to.key=value)
	  -dump
	    	create defaults.toml
	  -f value
	    	read config from TOML files

	Configuration settings:
	  comment string
	  size.width uint32
	  size.height uint32
	  audio.enabled bool
	  audio.samplerate int

*/
package confi
