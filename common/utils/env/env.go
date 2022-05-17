package envUtils

import "os"

const (
	EnvKey = "ENV"
)

type TypEnv string

var (
	env TypEnv
)

const (
	envDev  TypEnv = "dev"
	envTest TypEnv = "test"
	envPre  TypEnv = "pre"
	envProd TypEnv = "prod"
)

func init() {
	getEnv := TypEnv(os.Getenv(EnvKey))
	if getEnv == "" {
		env = "local"
	} else {
		env = getEnv
	}
}

func IsDev() bool {
	return env == envDev
}

func IsTest() bool {
	return env == envTest
}

func IsPre() bool {
	return env == envPre
}

func IsProd() bool {
	return env == envProd
}

func IsCustom() bool {
	return !IsDev() && !IsTest() && !IsPre() && !IsProd()
}

func Env() TypEnv {
	return env
}

func (t TypEnv) String() string {
	return string(t)
}
