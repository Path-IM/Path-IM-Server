package rc

type Option func(*relationOption)

type relationOption struct {
	order string // 排序字段
	size  int    // 存储大小 <=0 则会全部存储
}

func defaultRelationOption() *relationOption {
	return &relationOption{
		order: "created_at desc",
		size:  0,
	}
}

func Order(field string) Option {
	return func(r *relationOption) {
		r.order = field
	}
}

func Size(limit int) Option {
	return func(r *relationOption) {
		r.size = limit
	}
}
