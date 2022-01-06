package utils

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

func ResolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		return ":12345"
	case 1:
		return addr[0]
	default:
		panic("too many parameters")
	}
}

func GetPort(addr string) string {
	u, err := url.Parse(addr)
	if err != nil {
		return ""
	}
	return u.Port()
}

func GetHost(addr string) string {
	u, err := url.Parse(addr)
	if err != nil {
		return ""
	}
	return u.Host
}

func BuildEndPoint(host, port string) string {
	rawUrl := fmt.Sprintf("%s:%s", host, port)
	if !strings.HasPrefix(rawUrl, "http") {
		rawUrl = "http://" + rawUrl
	}

	u, err := url.Parse(rawUrl)
	if err != nil {
		return ""
	}
	return u.String()
}

func NameOfFunction(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func WatchSignal() {
	quit := make(chan os.Signal)
	signal.Notify(quit,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGILL,
		syscall.SIGABRT,
		syscall.SIGFPE,
		syscall.SIGKILL,
		syscall.SIGSEGV,
		//syscall.SIGPIPE,
		syscall.SIGALRM,
		syscall.SIGTERM)
	<-quit
}

func GetLocalIp() string {
	var ips []string
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, a := range addresses {
		if ipNet, ok := a.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			ips = append(ips, ipNet.IP.String())
		}
	}

	if len(ips) > 0 {
		return ips[0]
	}

	return ""
}

func ObjToBytes(obj interface{}) []byte {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return nil
	}
	return bytes
}

func ObjToStr(obj interface{}) string {
	return Bytes2Str(ObjToBytes(obj))
}

func Bytes2Str(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}

func Str2bytes(str string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&str))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func Int2Str(i int64) string {
	return strconv.FormatInt(i, 10)
}

var render = rand.Reader

func Uuid() string {
	var buf [36]byte
	type UUID [16]byte
	var uuid UUID
	_, _ = io.ReadFull(render, uuid[:])
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10

	hex.Encode(buf[:], uuid[:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], uuid[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], uuid[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], uuid[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], uuid[10:])
	return Bytes2Str(buf[:])
}
