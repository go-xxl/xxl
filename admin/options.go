package admin

import (
	"time"
)

type Options struct {
	AdmAddresses  []string
	AccessToken   string
	Timeout       time.Duration
	RegistryGroup string
	RegistryKey   string
	RegistryValue string
}

type Option func(o *Options)

func SetAdmAddresses(addresses []string) Option {
	return func(o *Options) {
		o.AdmAddresses = addresses
	}
}

func SetAccessToken(accessToken string) Option {
	return func(o *Options) {
		o.AccessToken = accessToken
	}
}

func SetTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.Timeout = timeout
	}
}

func SetRegistryGroup(registryGroup string) Option {
	return func(o *Options) {
		o.RegistryGroup = registryGroup
	}
}

func SetRegistryKey(registryKey string) Option {
	return func(o *Options) {
		o.RegistryKey = registryKey
	}
}

func SetRegistryValue(registryValue string) Option {
	return func(o *Options) {
		o.RegistryValue = registryValue
	}
}
