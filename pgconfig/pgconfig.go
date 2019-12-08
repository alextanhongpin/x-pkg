package pgconfig

import (
	"fmt"
	"sort"
	"strings"
)

type Option map[string]interface{}

func (o Option) String() string {
	var result []string
	for k, v := range o {
		result = append(result, fmt.Sprintf("%s='%v'", k, v))
	}
	sort.Strings(result)
	return strings.Join(result, " ")
}

type OptionFn func(Option)

func Name(name string) OptionFn {
	return func(opt Option) {
		opt["dbname"] = strings.TrimSpace(name)
	}
}

func User(user string) OptionFn {
	return func(opt Option) {
		opt["user"] = strings.TrimSpace(user)
	}
}

func Password(pwd string) OptionFn {
	return func(opt Option) {
		opt["password"] = strings.TrimSpace(pwd)
	}
}
func Host(host string) OptionFn {
	return func(opt Option) {
		opt["host"] = strings.TrimSpace(host)
	}
}
func Port(port int) OptionFn {
	return func(opt Option) {
		opt["port"] = port
	}
}
func SSLMode(sslmode string) OptionFn {
	return func(opt Option) {
		opt["sslmode"] = strings.TrimSpace(sslmode)
	}
}

func FallbackApplicationName(name string) OptionFn {
	return func(opt Option) {
		opt["fallback_application_name"] = strings.TrimSpace(name)
	}
}

func ConnectTimeout(seconds int) OptionFn {
	return func(opt Option) {
		opt["connect_timeout"] = seconds
	}
}
func SSLCert(cert string) OptionFn {
	return func(opt Option) {
		opt["sslcert"] = strings.TrimSpace(cert)
	}
}

func SSLKey(cert string) OptionFn {
	return func(opt Option) {
		opt["sslkey"] = strings.TrimSpace(cert)
	}
}
func SSLRootCert(cert string) OptionFn {
	return func(opt Option) {
		opt["sslrootcert"] = strings.TrimSpace(cert)
	}
}

func New(fns ...OptionFn) Option {
	opt := Option{
		"dbname":   "postgres",
		"host":     "localhost",
		"password": "postgres",
		"port":     5432,
		"sslmode":  "disable",
		"user":     "root",
	}
	for _, fn := range fns {
		fn(opt)
	}
	return opt
}
