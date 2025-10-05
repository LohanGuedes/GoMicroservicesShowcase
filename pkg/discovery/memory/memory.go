package memory

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/lohanguedes/movie-microservices/pkg/discovery"
)

type Registry struct {
	sync.RWMutex
	serviceAddrs map[string]map[string]*serviceInstance
}

type serviceInstance struct {
	hostPort   string
	lastActive time.Time
}

func NewRegistry() *Registry {
	return &Registry{
		serviceAddrs: map[string]map[string]*serviceInstance{},
	}
}

func (r *Registry) Register(
	ctx context.Context,
	instanceID string,
	serviceName string,
	hostPort string,
) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddrs[serviceName]; !ok {
		r.serviceAddrs[serviceName] = map[string]*serviceInstance{}
	}
	r.serviceAddrs[serviceName][instanceID] = &serviceInstance{hostPort: hostPort, lastActive: time.Now()}
	return nil
}

func (r *Registry) Deregister(
	ctx context.Context,
	instanceID string,
	serviceName string,
) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddrs[serviceName]; !ok {
		return nil
	}
	delete(r.serviceAddrs[serviceName], instanceID)
	return nil
}

func (r *Registry) ReportHealhyState(instanceID string, serviceName string) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddrs[serviceName]; !ok {
		return errors.New("service not registered")
	}

	if _, ok := r.serviceAddrs[serviceName][instanceID]; !ok {
		return errors.New("instance" + instanceID + " of service" + serviceName + " is not registered")
	}

	r.serviceAddrs[serviceName][instanceID].lastActive = time.Now()
	return nil
}

func (r *Registry) ServiceAddresses(
	ctx context.Context,
	serviceID string,
) ([]string, error) {
	r.RLock()
	defer r.RUnlock()
	if _, ok := r.serviceAddrs[serviceID]; !ok {
		return nil, discovery.ErrNotFound
	}

	res := []string{}
	for instanceID, instance := range r.serviceAddrs[serviceID] {
		if instance.lastActive.Before(time.Now().Add(-5 * time.Second)) {
			log.Println("Instance " + instanceID + "of service" + serviceID + " is not active, skipping")
			continue
		}
		res = append(res, instance.hostPort)
	}
	return res, nil
}
