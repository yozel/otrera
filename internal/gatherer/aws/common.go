package aws

import (
	"log"

	"github.com/pkg/errors"
	"gopkg.in/ini.v1"
)

func ListProfiles(filename string) ([]string, error) {
	var r []string
	cfg, err := ini.Load(filename)
	err = errors.Wrapf(err, "Failed to read file")
	if err != nil {
		return nil, err
	}
	sections := cfg.Sections()
	for _, section := range sections {
		var pn string
		sn := section.Name()
		if sn == "DEFAULT" {
			continue
		} else if sn == "default" {
			pn = sn
		} else {
			if sn[0:8] != "profile " {
				log.Printf("Can't parse section: %s\n", sn)
				continue
			}
			pn = sn[8:]
		}
		r = append(r, pn)
	}
	return r, nil
}
