package rc

const MagicKey = "$"

var (
	gtKey    = MagicKey + "gt"
	gteKey   = MagicKey + "gte"
	ltKey    = MagicKey + "lt"
	lteKey   = MagicKey + "lte"
	whereKey = MagicKey + "where"
)

func KeyGt(key string) string {
	return key + gtKey
}

func KeyGte(key string) string {
	return key + gteKey
}

func KeyLt(key string) string {
	return key + ltKey
}

func KeyLte(key string) string {
	return key + lteKey
}

func KeyWhere(key string) string {
	return key + whereKey
}
