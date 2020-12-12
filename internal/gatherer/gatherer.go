package gatherer

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Gatherer is an struct to gather Description
type Gatherer struct {
	cachePath   string
	descriptors map[string]func(options map[string]string) ([]RawObjectInterface, error)
}

func (g *Gatherer) getCachePath(cacheKey string, options map[string]string) string {
	keys := make([]string, 0, len(options))
	for k := range options {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	h := sha256.New()
	for _, k := range keys {
		k = url.QueryEscape(k)
		v := url.QueryEscape(options[k])
		_, err := h.Write([]byte(fmt.Sprintf("%s:%s;", k, v)))
		if err != nil {
			panic(err)
		}
	}

	key := fmt.Sprintf("%s-%s", strings.Replace(cacheKey, "/", ".", -1), fmt.Sprintf("%x", h.Sum(nil)))
	return path.Join(g.cachePath, key)
}

func (g *Gatherer) setCache(cacheFilePath string, objects []RawObjectInterface, ttl time.Duration) error {
	rawObjects := []RawObject{}
	for _, obj := range objects {
		rawObjects = append(rawObjects, obj.Copy())
	}
	b, err := json.Marshal(rawObjects)
	if err != nil {
		return err // TODO: wrap error
	}
	err = ioutil.WriteFile(fmt.Sprintf("%s-%d", cacheFilePath, time.Now().Unix()+int64(ttl.Seconds())), b, 0644)
	if err != nil {
		return err // TODO: wrap error
	}
	return nil
}

func (g *Gatherer) getCache(cacheFilePath string) ([]RawObject, error) {

	files, err := filepath.Glob(fmt.Sprintf("%s-*", cacheFilePath))
	if err != nil {
		return nil, err // TODO: wrap error
	}
	if len(files) == 0 {
		return nil, nil
	}
	sort.Strings(files)
	cacheFilePath = files[len(files)-1]

	cs := strings.Split(cacheFilePath, "-")

	timestamp, err := strconv.Atoi(cs[len(cs)-1])
	if err != nil {
		return nil, err // TODO: wrap error
	}

	if int64(timestamp) < time.Now().Unix() {
		return nil, nil
	}

	b, err := ioutil.ReadFile(cacheFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err // TODO: wrap error
	}

	var d []RawObject
	err = json.Unmarshal(b, &d)
	if err != nil {
		return nil, err // TODO: wrap error
	}
	return d, nil
}

func (g *Gatherer) UpdateCache(key string, options map[string]string, ttl time.Duration) error {
	cp := g.getCachePath(key, options)
	r, err := g.descriptors[key](options)
	if err != nil {
		return err // TODO: wrap error
	}
	err = g.setCache(cp, r, ttl)
	if err != nil {
		return err // TODO: wrap error
	}
	return nil
}

// Gather returns Description for given name and options with cache
func (g *Gatherer) Gather(key string, options map[string]string, ttl time.Duration) ([]RawObject, error) {
	cp := g.getCachePath(key, options)
	r, err := g.getCache(cp)
	if err != nil {
		return nil, err // TODO: wrap error
	}
	if r == nil {
		err = g.UpdateCache(key, options, ttl)
		if err != nil {
			log.Fatal(err)
		}
	}

	r, err = g.getCache(cp)
	if err != nil {
		return nil, err // TODO: wrap error
	}
	if r == nil {
		log.Fatalf("Something wrong, cache is still empty after UpdateCache")
	}
	return r, nil
}

// New creates a new Gatherer
func New(cachePath string, descriptors map[string]func(options map[string]string) ([]RawObjectInterface, error)) *Gatherer {
	g := &Gatherer{
		cachePath:   cachePath,
		descriptors: descriptors,
	}
	return g
}
