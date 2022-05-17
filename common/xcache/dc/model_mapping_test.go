package dc

import (
	"testing"
)

func TestDbMapping_FirstByID(t *testing.T) {
	db := GetDbMapping(nil, nil)
	db.FirstByID(nil, WithFieldId(""))
}
