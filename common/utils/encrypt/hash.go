package encrypt

import "github.com/speps/go-hashids"

func EnCodeInt64(id int64, salt string) (string, error) {
	hd := hashids.NewData()
	hd.Salt = salt //盐值
	hd.MinLength = 10
	h, _ := hashids.NewWithData(hd)
	return h.EncodeInt64([]int64{id})
}

func DeCodeInt64(code string, salt string) (int64, error) {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = 10
	h, _ := hashids.NewWithData(hd)
	d, err := h.DecodeInt64WithError(code)
	if err == nil {
		return d[0], nil
	}
	return 0, err
}
