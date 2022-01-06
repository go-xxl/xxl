package xxl

import (
	"time"
)

type Options struct {
	AdmAddresses  []string
	AccessToken   string
	Timeout       time.Duration
	ExecutorIp    string
	ExecutorPort  string
	RegistryKey   string
	RegistryValue string
}

type Option func(o *Options)

func AdmAdmAddresses(addr ...string) Option {
	return func(o *Options) {
		o.AdmAddresses = addr
	}
}

func AccessToken(token string) Option {
	return func(o *Options) {
		o.AccessToken = token
	}
}

func ExecutorIp(ip string) Option {
	return func(o *Options) {
		o.ExecutorIp = ip
	}
}

func ExecutorPort(port string) Option {
	return func(o *Options) {
		o.ExecutorPort = port
	}
}

func RegistryKey(registryKey string) Option {
	return func(o *Options) {
		o.RegistryKey = registryKey
	}
}
