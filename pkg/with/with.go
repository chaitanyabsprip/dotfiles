package with

import "os"

type PopFunc func() error

func Pwd(dir string) (pop PopFunc, err error) {
	pwd := ``
	if pwd, err = os.Getwd(); err != nil {
		return nil, err
	}
	if err = os.Chdir(dir); err != nil {
		return nil, err
	}
	return func() error { return os.Chdir(pwd) }, nil
}

func Env(name, value string) (PopFunc, error) {
	old, exists := os.LookupEnv(name)
	if err := os.Setenv(name, value); err != nil {
		return nil, err
	}
	if exists {
		return func() error { return os.Setenv(name, old) }, nil
	} else {
		return func() error { return os.Unsetenv(name) }, nil
	}
}
