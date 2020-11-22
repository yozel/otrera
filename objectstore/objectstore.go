package objectstore

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/yozel/otrera/gatherer"
	"github.com/yozel/otrera/gatherer/aws"
)

type ObjectStore struct {
	store    map[string]Object
	gatherer *gatherer.Gatherer
}

func NewObjectStore() (*ObjectStore, error) {
	return &ObjectStore{
		store: make(map[string]Object),
		gatherer: gatherer.New(
			"/tmp/.otrera.cache",
			map[string]func(options map[string]string) ([]gatherer.RawObjectInterface, error){
				"AWS/EC2Instances": aws.DescribeEC2Instances,
			}), // TODO: get this from parameter
	}, nil
}

func (s *ObjectStore) Keys() []string {
	keys := make([]string, 0, len(s.store))
	for k := range s.store {
		keys = append(keys, k)
	}
	return keys
}

func (s *ObjectStore) Set(key string, l map[string]string, c time.Time, d interface{}) error {
	labels := make(map[string]string, len(l))
	for k, v := range l {
		labels[k] = v
	}
	r, err := DeepCopy(&d)
	if err != nil {
		return err // TODO: wrap error
	}
	s.store[key] = Object{Key: key, Labels: l, CreationTimestamp: c, Data: *r}
	return nil
}

func (s *ObjectStore) Get(key string) (Object, error) {
	return s.store[key], nil
}

func (s *ObjectStore) GetAll(key string, l map[string]string) ([]Object, error) {
	if l == nil {
		l = map[string]string{}
	}
	r := []Object{}
	keyParts := strings.Split(key, "*")
	for i, part := range keyParts {
		keyParts[i] = regexp.QuoteMeta(part)
	}
	if len(keyParts) > 0 {
		key = strings.Join(keyParts, ".*")
	} else {
		key = ".*"
	}
	keyR, err := regexp.Compile(key)
	if err != nil {
		return nil, err
	}

eachobject:
	for k, v := range s.store {
		if !keyR.MatchString(k) {
			continue
		}

		for lk, lv := range l {
			if v.Labels[lk] != lv {
				continue eachobject
			}
		}
		r = append(r, v)

	}
	return r, nil
}

func (s *ObjectStore) Gather(key string, o map[string]string, l map[string]string, ttl time.Duration) error {
	c, err := s.gatherer.Gather(key, o, ttl)
	if err != nil {
		return err // TODO: wrap error
	}
	for _, object := range c {
		s.Set(
			fmt.Sprintf("%s/%s", key, object.Name()),
			l,
			time.Now().UTC(),
			object.Content(),
		)
	}
	return nil
}
