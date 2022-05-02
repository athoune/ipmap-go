package ipmap

import (
	"github.com/athoune/ipmap-go/csv"
	"github.com/athoune/iptree/tree"
)

type Ranges struct {
	Tree tree.Trunk
	cpt  int
}

func New(c *csv.CVS) (*Ranges, error) {
	r := &Ranges{
		Tree: tree.NewTrunk(3),
	}
	for c.Next() {
		line, err := c.Value()
		if err != nil {
			return nil, err
		}
		n := line.Network()
		r.Tree.Append(&n, line)
		r.cpt++
	}

	return r, nil
}

func (r *Ranges) Length() int {
	return r.cpt
}
