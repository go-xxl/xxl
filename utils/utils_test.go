package utils

import "testing"

func TestBuildEndPoint(t *testing.T) {
	t.Log(BuildEndPoint("127.0.0.1", "8080"))
}

func TestUuid(t *testing.T) {
	t.Log(Uuid())
}

func BenchmarkUuid(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Uuid()
	}
}

func BenchmarkBytes2Str(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Bytes2Str([]byte{99,99,99,99})
	}
}