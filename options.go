package redis_pagination

type Options struct {
	Addr     string
	Password string
	DB       int
	Key      string
	Field    string
	Page     int
	PageSize int
}
