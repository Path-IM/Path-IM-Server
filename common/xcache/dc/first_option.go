package dc

type FuncFirstOption func(*FirstOption)

type FirstOption struct {
	order string
}

func WithOrder(
	order string,
) FuncFirstOption {
	return func(option *FirstOption) {
		option.order = order
	}
}

func defaultOption() *FirstOption {
	return &FirstOption{
		order: "",
	}
}
