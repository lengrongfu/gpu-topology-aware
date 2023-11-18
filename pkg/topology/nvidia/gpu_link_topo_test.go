package nvidia

import "testing"

func Test_printTopology(t *testing.T) {
	inputMatrix := [][]LinkType{
		{P2PLinkSameBoard, P2PLinkSingleSwitch, P2PLinkHostBridge, P2PLinkHostBridge, P2PLinkCrossCPU},
		{P2PLinkSingleSwitch, P2PLinkSameBoard, P2PLinkHostBridge, P2PLinkHostBridge, P2PLinkCrossCPU},
		{P2PLinkHostBridge, P2PLinkHostBridge, P2PLinkSameBoard, P2PLinkSingleSwitch, P2PLinkCrossCPU},
		{P2PLinkHostBridge, P2PLinkHostBridge, P2PLinkSingleSwitch, P2PLinkSameBoard, P2PLinkCrossCPU},
		{P2PLinkCrossCPU, P2PLinkCrossCPU, P2PLinkCrossCPU, P2PLinkCrossCPU, P2PLinkSameBoard},
	}
	g := gpuTopologyLinkCollector{}
	g.printTopology(inputMatrix)
}
