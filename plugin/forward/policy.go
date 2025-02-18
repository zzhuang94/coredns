package forward

import (
	"sync/atomic"
	"time"

	"github.com/coredns/coredns/plugin/pkg/rand"
)

// Policy defines a policy we use for selecting upstreams.
type Policy interface {
	List([]*Proxy) []*Proxy
	String() string
}

// random is a policy that implements random upstream selection.
type random struct{}

func (r *random) String() string { return "random" }

func (r *random) List(p []*Proxy) []*Proxy {
	switch len(p) {
	case 1:
		return p
	case 2:
		if rn.Int()%2 == 0 {
			return []*Proxy{p[1], p[0]} // swap
		}
		return p
	}

	perms := rn.Perm(len(p))
	rnd := make([]*Proxy, len(p))

	for i, p1 := range perms {
		rnd[i] = p[p1]
	}
	return rnd
}

// roundRobin is a policy that selects hosts based on round robin ordering.
type roundRobin struct {
	robin uint32
}

func (r *roundRobin) String() string { return "round_robin" }

func (r *roundRobin) List(p []*Proxy) []*Proxy {
	poolLen := uint32(len(p))
	i := atomic.AddUint32(&r.robin, 1) % poolLen

	robin := []*Proxy{p[i]}
	robin = append(robin, p[:i]...)
	robin = append(robin, p[i+1:]...)

	return robin
}

// sequential is a policy that selects hosts based on sequential ordering.
type sequential struct{}

func (r *sequential) String() string { return "sequential" }

func (r *sequential) List(p []*Proxy) []*Proxy {
	return p
}

// race is a policy that try all hosts at once and pick the fastest result.
type race struct{}

func (r *race) String() string { return "race" }

func (r *race) List(p []*Proxy) []*Proxy {
	return p
}

var rn = rand.New(time.Now().UnixNano())
