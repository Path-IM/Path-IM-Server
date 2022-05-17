package dc

func defaultListOption() *FirstOption {
	return &FirstOption{
		where:     "",
		args:      []interface{}{},
		keySuffix: "",
		fieldId:   "id",
	}
}
