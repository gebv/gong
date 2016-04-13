package store

import (
	"encoding/binary"
	"github.com/golang/glog"
	"strconv"
)

func NewIDFromString(v string) ID {
	id := ID(0)
	if err := id.SetString(v); err != nil {
		glog.Errorf("id: invalid id=%v, err=%v", v, err)
	}
	return id
}

type ID uint64

func (i ID) IsNil() bool {
	return i == 0
}

func (i *ID) SetString(v string) error {
	_id, err := strconv.ParseUint(v, 10, 64)

	if err != nil {
		return err
	}

	*i = ID(_id)

	return err
}

func (i *ID) SetUint64(v uint64) {
	*i = ID(v)
}

func (i ID) String() string {

	return strconv.FormatUint(uint64(i), 10)
}

func (i ID) Bytes() []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(i))
	return b
}
