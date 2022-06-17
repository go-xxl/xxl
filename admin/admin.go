package admin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-xxl/xxl/log"
	"github.com/go-xxl/xxl/utils"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var (
	client *AdmApi
	once   sync.Once
)

type AdmApi struct {
	opt    Options
	isStop bool
}

func NewAdmApi() *AdmApi {
	once.Do(func() {
		client = &AdmApi{}
	})
	return client
}

func GetClient() *AdmApi {
	if client == nil {
		NewAdmApi()
	}
	return client
}

func (adm *AdmApi) SetOpt(opts ...Option) {
	for _, opt := range opts {
		opt(&adm.opt)
	}
}

func (adm *AdmApi) Register() {
	opt := adm.opt

	if opt.RegistryKey == EmptyString ||
		opt.RegistryValue == EmptyString {
		panic("Illegal Argument.")
	}

	registerParams := Registry{
		RegistryGroup: opt.RegistryGroup,
		RegistryKey:   opt.RegistryKey,
		RegistryValue: opt.RegistryValue,
	}

	go func() {
		t := time.NewTimer(LoopFrequency)

		for !adm.isStop {
			for _, address := range adm.opt.AdmAddresses {
				go func(address string) {
					resp := adm.postBody(address+UrlRegister, opt.AccessToken, opt.Timeout, registerParams)
					if resp.Code == SuccessCode {
						log.Debug(">>>>>>>>>>> xxl-job registry success",
							log.Field("registerParams", utils.ObjToStr(registerParams)),
							log.Field("resp", utils.ObjToStr(resp)),
							log.Field("url", address+UrlRegister),
						)
					} else {
						log.Info(">>>>>>>>>>> xxl-job registry fail ",
							log.Field("registryParam", utils.ObjToStr(registerParams)),
							log.Field("registryResult", utils.ObjToStr(resp)),
							log.Field("url", address+UrlRegister),
						)
					}
				}(address)
			}
			<-t.C
			t.Reset(LoopFrequency)
		}

		if adm.isStop {
			t.Stop()
		}
	}()
}

func (adm *AdmApi) RegistryRemove() {

	adm.isStop = true

	opt := adm.opt

	registerParams := Registry{
		RegistryGroup: opt.RegistryGroup,
		RegistryKey:   opt.RegistryKey,
		RegistryValue: opt.RegistryValue,
	}

	for _, address := range adm.opt.AdmAddresses {

		resp := adm.postBody(address+UrlRegistryRemove, opt.AccessToken, opt.Timeout, registerParams)

		if resp.Code == SuccessCode {
			log.Debug(">>>>>>>>>>> xxl-job registry-remove success",
				log.Field("registryParam", utils.ObjToStr(registerParams)),
				log.Field("registryResult", utils.ObjToStr(resp)),
				log.Field("url", address+UrlRegister),
			)
		} else {
			log.Warn(">>>>>>>>>>> xxl-job registry-remove fail",
				log.Field("registryParam", utils.ObjToStr(registerParams)),
				log.Field("registryResult", utils.ObjToStr(resp)),
				log.Field("url", address+UrlRegister),
			)
		}
	}
}

func (adm *AdmApi) Callback(call []HandleCallbackParams) {

	opt := adm.opt

	for _, address := range adm.opt.AdmAddresses {
		resp := adm.postBody(address+UrlCallback, opt.AccessToken, opt.Timeout, call)

		if resp.Code == SuccessCode {
			log.Debug(">>>>>>>>>>> xxl-job callback success",
				log.Field("call", utils.ObjToStr(call)),
				log.Field("respResult", utils.ObjToStr(resp)),
				log.Field("url", address+UrlCallback),
			)

		} else {
			log.Warn(">>>>>>>>>>> xxl-job callback fail",
				log.Field("call", utils.ObjToStr(call)),
				log.Field("respResult", utils.ObjToStr(resp)),
				log.Field("url", address+UrlCallback),
			)
		}
	}

}

func (adm *AdmApi) postBody(requestUrl, accessToken string, timeout time.Duration, requestObj interface{}) Resp {
	resp := Resp{
		Code:    FailCode,
		Msg:     "",
		Content: nil,
	}

	u, err := url.Parse(requestUrl)
	if err != nil {
		resp.Msg = err.Error()
		return resp
	}

	header := http.Header{}
	header.Set("Content-Type", "application/json")
	if accessToken != EmptyString {
		header.Set(XxlJobAccessToken, accessToken)
	}

	data := utils.ObjToBytes(requestObj)

	req := &http.Request{
		Method: http.MethodPost,
		URL:    u,
		Host:   u.Host,
		Header: header,
		Body:   ioutil.NopCloser(bytes.NewBuffer(data)),
	}
	c := &http.Client{
		Timeout: timeout,
	}
	respHandler, err := c.Do(req)
	if err != nil {
		resp.Msg = err.Error()
		return resp
	}
	defer respHandler.Body.Close()

	if respHandler.StatusCode != http.StatusOK {
		resp.Code = FailCode
		resp.Msg = fmt.Sprintf("xxl-job remoting fail, StatusCode(%d) invalid. for url : %s",
			respHandler.StatusCode, requestUrl)
		return resp
	}

	all, err := ioutil.ReadAll(respHandler.Body)
	if err != nil {
		resp.Msg = err.Error()
		return resp
	}

	_ = json.Unmarshal(all, &resp)
	return resp

}
