// Package valueobject
package valueobject

import "strconv"

type ID int64

func (i ID) Value() int64 {
	return int64(i)
}

func (i ID) String() string {
	return strconv.FormatInt(int64(i), 10)
}

func (i ID) Int() int {
	return int(i)
}
