package ipmap

import (
	"github.com/athoune/ipmap-go/csv"
	"github.com/yl2chen/cidranger"
)

type Ranges struct {
	Ranger cidranger.Ranger
}

func New(c *csv.CVS) (*Ranges, error) {
	r := &Ranges{
		Ranger: cidranger.NewPCTrieRanger(),
	}
	for c.Next() {
		line, err := c.Value()
		if err != nil {
			return nil, err
		}
		r.Ranger.Insert(line)
	}

	return r, nil
}
