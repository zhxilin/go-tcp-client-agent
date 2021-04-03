package model

type EventItem struct {
	Type     int
	UserData interface{}
}

type Config struct {
	Id   int
	Host string
	Port int
}
