package hredis

type UnlockFunc func() error

type Option struct {
	Addr     string
	Password string
	DB       int
	Timeout  int64
}
