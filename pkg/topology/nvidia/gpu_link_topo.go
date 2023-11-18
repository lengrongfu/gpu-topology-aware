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

package nvidia

import (
	"context"
	"fmt"
	"github.com/DaoCloud-OpenSource/gpu-topology-aware/pkg/topology"
	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"log"
)

const (
	// TopologyNv is gpy link type ls nvlink
	TopologyNv nvml.GpuTopologyLevel = 60
)

func NewLinkCollector() topology.Collector {
	return &gpuTopologyLinkCollector{}
}

var _ topology.Collector = &gpuTopologyLinkCollector{}

type gpuTopologyLinkCollector struct {
}

func (g *gpuTopologyLinkCollector) Collect(ctx context.Context) (topology.Matrix, error) {
	result := nvml.Init()
	if result != nvml.SUCCESS {
		return nil, fmt.Errorf("unable to initialize NVML: %v", nvml.ErrorString(result))
	}
	defer func() {
		result = nvml.Shutdown()
		if result != nvml.SUCCESS {
			log.Printf("unable to shutdown NVML: %v", nvml.ErrorString(result))
		}
	}()

	count, ret := nvml.DeviceGetCount()
	if ret != nvml.SUCCESS {
		log.Printf("unable to get device count: %v", nvml.ErrorString(ret))
		return nil, fmt.Errorf("unable to get device count: %v", nvml.ErrorString(ret))
	}
	matrix := make([][]LinkType, count)
	for i := 0; i < count; i++ {
		matrix[i] = make([]LinkType, count)
		for j := 0; j < count; j++ {
			if i == j {
				matrix[i][j] = 0
				continue
			}
			deviceI, ret := nvml.DeviceGetHandleByIndex(i)
			if ret != nvml.SUCCESS {
				log.Printf("unable to get device at index %d: %v", i, nvml.ErrorString(ret))
				return nil, fmt.Errorf("unable to get device at index %d: %v", i, nvml.ErrorString(ret))
			}
			deviceJ, ret := nvml.DeviceGetHandleByIndex(j)
			if ret != nvml.SUCCESS {
				log.Printf("unable to get device at index %d: %v", j, nvml.ErrorString(ret))
				return nil, fmt.Errorf("unable to get device at index %d: %v", j, nvml.ErrorString(ret))
			}
			deviceLinkType, err := g.collectDeviceToDeviceLinkType(deviceI, deviceJ)
			if err != nvml.SUCCESS {
				matrix[i][j] = P2PLinkUnknown
			} else {
				matrix[i][j] = deviceLinkType
			}
		}
	}
	return g.convertToMatrix(matrix), nil
}

// collectDeviceToDeviceLinkType is get gpu to gpu link type
func (g *gpuTopologyLinkCollector) collectDeviceToDeviceLinkType(src, dest nvml.Device) (LinkType, nvml.Return) {
	srcUUID, ret := src.GetUUID()
	if ret != nvml.SUCCESS {
		return P2PLinkUnknown, ret
	}
	destUUID, ret := dest.GetUUID()
	if ret != nvml.SUCCESS {
		return P2PLinkUnknown, ret
	}
	if srcUUID == destUUID {
		return P2PLinkSameBoard, nvml.SUCCESS
	}
	linkType, err := getP2PLink(src, dest)
	if err != nvml.SUCCESS {
		return linkType, err
	}
	nvLink, err := getNVLink(src, dest)
	if err != nvml.SUCCESS {
		return linkType, err
	}
	return nvLink, nvml.SUCCESS
}

func (g *gpuTopologyLinkCollector) printTopology(matrix [][]LinkType) {
	log.Println("gpu topology link type: ")
	// 打印表头
	fmt.Print("\t")
	for i := 0; i < len(matrix); i++ {
		fmt.Printf("GPU%d\t", i)
	}
	fmt.Println()

	// 打印内容
	for i := 0; i < len(matrix); i++ {
		fmt.Printf("GPU%d\t", i)
		for j := 0; j < len(matrix[i]); j++ {
			fmt.Printf("%s\t", LinkTypeMap[matrix[i][j]])
		}
		fmt.Println()
	}
}

func (g *gpuTopologyLinkCollector) convertToMatrix(matrix [][]LinkType) topology.Matrix {
	g.printTopology(matrix)
	resultMatrix := make(topology.Matrix, len(matrix))
	for i := range matrix {
		resultMatrix[i] = make([]uint64, 0)
		for j := range matrix[i] {
			resultMatrix[i][j] = uint64(matrix[i][j])
		}
	}
	return resultMatrix
}
