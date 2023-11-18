package nvidia

import (
	"context"
	"github.com/DaoCloud-OpenSource/gpu-topology-aware/pkg/topology"
	"github.com/NVIDIA/go-nvml/pkg/nvml"
)

var _ topology.Collector = &gpuTopologyBandwidthCollector{}

type gpuTopologyBandwidthCollector struct {
}

func (g *gpuTopologyBandwidthCollector) Collect(ctx context.Context) (topology.Matrix, error) {
	//TODO implement me
	panic("implement me")
}

// collectDeviceToDeviceBandwidth is get gpu to gpu communication bandwidth
func (g *gpuTopologyBandwidthCollector) collectDeviceToDeviceBandwidth(src, dest nvml.Device) uint64 {
	return 0
}
