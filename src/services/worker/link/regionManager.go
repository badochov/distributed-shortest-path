package link

import (
	"context"
	"fmt"
	"math/rand"
	"sync"

	"github.com/badochov/distributed-shortest-path/src/services/worker/discoverer"
	"github.com/hashicorp/go-multierror"
	"golang.org/x/exp/slices"
)

type Address string

type RegionManager interface {
	UpdateInstances(instances []discoverer.WorkerInstance) error
	GetLink() (Address, Link, error)
}

type regionManger struct {
	rwlock    sync.RWMutex
	links     map[Address]Link
	addresses []Address
	port      int
}

func (r *regionManger) UpdateInstances(instances []discoverer.WorkerInstance) error {
	r.rwlock.Lock()
	defer r.rwlock.Unlock()

	r.addresses = make([]Address, 0, len(instances))
	for _, i := range instances {
		r.addresses = append(r.addresses, r.toAddr(i.Ip))
	}

	newLinks := make(map[Address]Link, len(instances))
	for _, i := range instances {
		a := r.toAddr(i.Ip)
		l, ok := r.links[a]
		if ok {
			delete(r.links, a)
		} else {
		}
		newLinks[a] = l
	}

	var err error
	// Close no longer existing links
	for a, l := range r.links {
		if closeErr := l.Close(); closeErr != nil {
			err = multierror.Append(err, fmt.Errorf("error closing link to %s, %w", a, err))
		}
	}
	r.links = newLinks

	return err
}

func (r *regionManger) UpdateInstanceStatus(ctx context.Context, s discoverer.WorkerInstanceStatus) (err error) {
	r.rwlock.Lock()
	defer r.rwlock.Unlock()

	addr := r.toAddr(s.Ip)
	if !slices.Contains(r.addresses, addr) {
		r.addresses = append(r.addresses, addr)
	}
	if s.IsRunning() {
		r.links[addr], err = New(ctx, string(addr))
		if err != nil {
			return
		}
	} else {
		l, ok := r.links[addr]
		if ok {
			delete(r.links, addr)
			err = l.Close()
			if err != nil {
				return
			}
		}
	}

	return
}

func (r *regionManger) toAddr(ip string) Address {
	return Address(fmt.Sprintf("%s:%d", ip, r.port))
}

func (r *regionManger) GetLink() (Address, Link, error) {
	r.rwlock.RLock()
	defer r.rwlock.RUnlock()

	if len(r.addresses) == 0 {
		return "", nil, fmt.Errorf("no links to region")
	}
	idx := rand.Intn(len(r.addresses))
	addr := r.addresses[idx]
	return addr, r.links[addr], nil
}

func NewRegionManager(port int) RegionManager {
	return &regionManger{
		port: port,
	}
}
