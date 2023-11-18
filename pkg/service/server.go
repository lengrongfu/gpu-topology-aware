/*
Copyright 2022 The Kangaroo Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/DaoCloud-OpenSource/gpu-topology-aware/pkg/topology"
	"github.com/DaoCloud-OpenSource/gpu-topology-aware/pkg/topology/nvidia"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	DefaultInterval          = time.Second * 60
	GpuTopologyAnnotationKey = "gpu.topology/value"
)

type Config struct {
	// interval is collect gpu topology info use ticker triage
	Interval string `json:"interval,omitempty"`
	// Watch is collect gpu topology use watch node resource change to triage
	Watch bool `json:"watch,omitempty"`
	// NodeName is current kubernetes cluster node name
	NodeName string `json:"nodeName"`
}

type Service struct {
	c         Config
	clientset *kubernetes.Clientset
	topoColl  topology.Collector
}

func NewService(c Config) Service {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	nvidiaTopoCollect := nvidia.NewLinkCollector()
	return Service{
		clientset: clientset,
		c:         c,
		topoColl:  nvidiaTopoCollect,
	}
}

func (s *Service) Run(ctx context.Context) error {
	if s.c.Watch {
		return s.watchCollect(ctx)
	}
	return s.intervalCollect(ctx)
}

func (s *Service) intervalCollect(ctx context.Context) error {
	var (
		interval time.Duration
		err      error
	)
	if s.c.Interval != "" {
		interval, err = time.ParseDuration(s.c.Interval)
		if err != nil {
			log.Printf("parse interval value invalid, config error")
			interval = DefaultInterval
		}
	}
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			matrix, err := s.topoColl.Collect(ctx)
			if err != nil {
				log.Printf("gpu topology collect error: %v", err)
				break
			}
			matrixBytes, err := json.Marshal(matrix)
			if err != nil {
				log.Printf("marshal gpu topology matrix error: %v", err)
				break
			}
			patchData := []byte(fmt.Sprintf(`{"metadata":{"annotation":{%s:%s}}}`, GpuTopologyAnnotationKey, matrixBytes))
			_, err = s.clientset.CoreV1().Nodes().Patch(ctx, s.c.NodeName, types.MergePatchType, patchData, metav1.PatchOptions{})
			if err != nil {
				log.Printf("patch node error: %v", err)
			}
		case <-ctx.Done():
			log.Printf("stop interval collect gpu topology info")
		}
	}
	return nil
}

func (s *Service) watchCollect(ctx context.Context) error {
	return nil
}
