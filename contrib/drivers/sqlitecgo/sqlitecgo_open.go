// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package sqlitecgo implements gdb.Driver, which supports operations for database SQLite.
//
// Note:
//  1. Using sqlitecgo is for building a 32-bit Windows operating system
//  2. You need to set the environment variable CGO_ENABLED=1 and make sure that GCC is installed
//     on your path. windows gcc: https://jmeubank.github.io/tdm-gcc/
package sqlitecgo

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gurl"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

// Open creates and returns an underlying sql.DB object for sqlite.
// https://github.com/mattn/go-sglite3
func (d *Driver) Open(config *gdb.ConfigNode) (db *sql.DB, err error) {
	var (
		source               string
		underlyingDriverName = "sqlite3"
	)
	source = config.Name
	// It searches the source file to locate its absolute path..
	if absolutePath, _ := gfile.Search(source); absolutePath != "" {
		source = absolutePath
	}

	// Multiple PRAGMAs can be specified, e.g.:
	// path/to/some.db?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)
	if config.Extra != "" {
		var (
			options  string
			extraMap map[string]interface{}
		)
		if extraMap, err = gstr.Parse(config.Extra); err != nil {
			return nil, err
		}
		for k, v := range extraMap {
			if options != "" {
				options += "&"
			}
			options += fmt.Sprintf(`_pragma=%s(%s)`, k, gurl.Encode(gconv.String(v)))
		}
		if len(options) > 1 {
			source += "?" + options
		}
	}

	if db, err = sql.Open(underlyingDriverName, source); err != nil {
		err = gerror.WrapCodef(
			gcode.CodeDbOperationError, err,
			`sql.Open failed for driver "%s" by source "%s"`, underlyingDriverName, source,
		)
		return nil, err
	}
	return
}
