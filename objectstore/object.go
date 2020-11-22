package objectstore

import "time"

type Object struct {
	Key               string
	Labels            map[string]string
	CreationTimestamp time.Time
	Data              interface{}
}
