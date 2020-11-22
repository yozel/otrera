package gatherer

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/yozel/otrera/gatherer/aws"
	"github.com/yozel/otrera/types"
)

// Gatherer is an struct to gather Description
type Gatherer struct {
	cachePath   string
	descriptors map[string]func(options map[string]string) (*types.Description, error)
}

func (g *Gatherer) getCachePath(name string, options map[string]string) string {
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

	key := fmt.Sprintf("%s-%s", name, fmt.Sprintf("%x", h.Sum(nil)))
	return path.Join(g.cachePath, key)
}

func (g *Gatherer) setCache(cacheFilePath string, desc *types.Description, ttl time.Duration) error {
	b, err := json.Marshal(desc)
	if err != nil {
		return err // TODO: wrap error
	}
	err = ioutil.WriteFile(fmt.Sprintf("%s-%d", cacheFilePath, time.Now().Unix()+int64(ttl.Seconds())), b, 0644)
	if err != nil {
		return err // TODO: wrap error
	}
	return nil
}

func (g *Gatherer) getCache(cacheFilePath string) (*types.Description, error) {

	files, err := filepath.Glob(fmt.Sprintf("%s-*", cacheFilePath))
	if err != nil {
		return nil, err // TODO: wrap error
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

	var desc types.Description
	err = json.Unmarshal(b, &desc)
	if err != nil {
		return nil, err // TODO: wrap error
	}
	return &desc, nil
}

// Gather returns Description for given name and options with cache
func (g *Gatherer) Gather(name string, options map[string]string, ttl time.Duration) (*types.Description, error) {
	cp := g.getCachePath(name, options)
	r, err := g.getCache(cp)
	if err != nil {
		return nil, err // TODO: wrap error
	}
	if r != nil {
		return r, nil
	}

	r, err = g.descriptors[name](options)
	if err != nil {
		return nil, err // TODO: wrap error
	}
	err = g.setCache(cp, r, ttl)
	if err != nil {
		return nil, err // TODO: wrap error
	}
	return r, nil
}

// New creates a new Gatherer
func New(cachePath string) *Gatherer {
	g := &Gatherer{
		cachePath: cachePath,
		descriptors: map[string]func(options map[string]string) (*types.Description, error){
			"EC2Instances": aws.DescribeEC2Instances,
		},
	}
	return g
}
