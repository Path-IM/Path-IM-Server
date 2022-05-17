package dc

type FuncFirstOption func(*FirstOption)

type FirstOption struct {
	where     string        // 缓存查不到 去数据库查询的时候 用 where语句 查询
	args      []interface{} // 缓存查不到 去数据库查询的时候 用 where语句+args 查询
	keySuffix string        // redis key的后缀
	fieldId   string        // 替换id的字段
}

func WithWhere(
	where string,
	args ...interface{},
) FuncFirstOption {
	return func(option *FirstOption) {
		option.where = where
		option.args = args
	}
}

func WithKeySuffix(
	suffix string,
) FuncFirstOption {
	return func(option *FirstOption) {
		option.keySuffix = suffix
	}
}

func WithFieldId(
	fieldId string,
) FuncFirstOption {
	return func(option *FirstOption) {
		option.fieldId = fieldId
		option.where = fieldId + " = ?"
	}
}

func defaultOption(id string) *FirstOption {
	return &FirstOption{
		where: "id = ?",
		args:  []interface{}{id},
	}
}
