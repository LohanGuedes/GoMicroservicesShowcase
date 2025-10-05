package consul

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	consul "github.com/hashicorp/consul/api"
	"github.com/lohanguedes/movie-microservices/pkg/discovery"
)

type Registry struct {
	client *consul.Client
}

func NewResgistry(addr string) (*Registry, error) {
	config := consul.DefaultConfig()
	config.Address = addr
	client, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &Registry{client}, nil
}

func (r *Registry) Register(
	ctx context.Context,
	instanceID string,
	serviceName string,
	hostPort string,
) error {
	parts := strings.Split(hostPort, ":")
	if len(parts) != 2 {
		return errors.New("hostport  must be in a form <host>:<port>, example -- localhost:8080")
	}
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}
	return r.client.Agent().ServiceRegister(&consul.AgentServiceRegistration{
		Address: parts[0],
		ID:      instanceID,
		Name:    serviceName,
		Port:    port,
		Check:   &consul.AgentServiceCheck{CheckID: instanceID, TTL: "5s"},
	})
}

func (r *Registry) Deregister(
	ctx context.Context,
	instanceID string,
	serviceName string,
) error {
	return r.client.Agent().CheckDeregister(instanceID)
}

func (r *Registry) ServiceAddresses(
	ctx context.Context,
	serviceID string,
) ([]string, error) {
	entries, _, err := r.client.Health().Service(serviceID, "", true, nil)
	if err != nil {
		return nil, err
	} else if len(entries) == 0 {
		return nil, discovery.ErrNotFound
	}

	res := []string{}
	for _, e := range entries {
		res = append(res, fmt.Sprintf("%s:%d", e.Service.Address, e.Service.Port))
	}
	return res, nil
}

func (r *Registry) ReportHealhyState(instanceID string, serviceName string) error {
	return r.client.Agent().PassTTL(instanceID, "")
}
