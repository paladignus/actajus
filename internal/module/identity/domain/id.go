// Package domain
package domain

import "strconv"

type (
	IDUser    int64
	IDSession int64
)

func (i IDUser) Value() int64   { return int64(i) }
func (i IDUser) String() string { return strconv.FormatInt(int64(i), 10) }

func (i IDSession) Value() int64   { return int64(i) }
func (i IDSession) String() string { return strconv.FormatInt(int64(i), 10) }
