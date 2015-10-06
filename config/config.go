package config

import (
	"github.com/samalba/dockerclient"
)

type Service struct {
	Name              string            `json:"name"`
	Version           string            `json:"version"`
	Instances         int64             `json:"instances"`
	Cpus              int64             `json:"cpus"`
	Memory            int64             `json:"memory"`
	Volumes           []Volume          `json:"volumes"`
	Networks          []Network         `json:"networks"`
	Pods              []Pod             `json:"pods"`
	UpdateConfig      UpdateConfig      `json:"updateConfig"`
	ReplicationConfig ReplicationConfig `json:"replicationConfig"`
	HealthCheck       HealthCheck       `json:"healthCheck"`
	Labels            []string          `json:"labels"`
}

type UpdateConfig struct {
	Delay                  string `json:"delay"`
	ReplacementConcurrency int64  `json:"replacementConcurrency"`
}

type ReplicationConfig struct {
	Strategy string `json:"strategy"`
}

type HealthCheck struct {
	Type    string `json:"type"`
	URL     string `json:"url"`
	Delay   string `json:"delay"`
	Retries int64  `json:"retries"`
}

type Network struct {
	Name    string            `json:"name"`
	Type    string            `json:"type"`
	Options map[string]string `json:"options"`
	Labels  []string          `json:"labels"`
}

type Volume struct {
	Name       string            `json:"name"`
	Driver     string            `json:"driver"`
	DriverOpts map[string]string `json:"driverOpts"`
	ReadOnly   bool              `json:"readOnly"`
	Labels     []string          `json:"labels"`
}

type Pod struct {
	Name       string      `json:"name"`
	Volumes    []string    `json:"volumes"`
	Networks   []string    `json:"networks"`
	Images     []Image     `json:"images"`
	Containers []Container `json:"containers"`
	Labels     []string    `json:"labels"`
}

type Image struct {
	Id           string      `json:"id"`
	Name         string      `json:"name"`
	Architecture string      `json:"architecture"`
	OS           string      `json:"os"`
	Size         int64       `json:"size"`
	VirtualSize  int64       `json:"virtualSize"`
	GraphDriver  GraphDriver `json:"graphDriver"`
	Labels       []string    `json:"labels"`
}

type GraphDriver struct {
	Name string            `json:"name"`
	Data map[string]string `json:"data"`
}

type Container struct {
	Config *dockerclient.ContainerConfig `json:"config"`
	Labels []string                      `json:"labels"`
}
