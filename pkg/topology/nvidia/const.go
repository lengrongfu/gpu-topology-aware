package nvidia

import (
	"errors"
	"log"
	"strconv"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
)

type LinkType uint

const (
	P2PLinkUnknown LinkType = iota
	P2PLinkCrossCPU
	P2PLinkSameCPU
	P2PLinkHostBridge
	P2PLinkMultiSwitch
	P2PLinkSingleSwitch
	P2PLinkSameBoard
	SingleNVLINKLink
	TwoNVLINKLinks
	ThreeNVLINKLinks
	FourNVLINKLinks
	FiveNVLINKLinks
	SixNVLINKLinks
	SevenNVLINKLinks
	EightNVLINKLinks
	NineNVLINKLinks
	TenNVLINKLinks
	ElevenNVLINKLinks
	TwelveNVLINKLinks
)

var (
	ErrUnsupportedP2PLink = errors.New("unsupported P2P link type")
	LinkTypeMap           = make(map[LinkType]string)
)

func init() {
	LinkTypeMap[P2PLinkUnknown] = "N/A"
	LinkTypeMap[P2PLinkCrossCPU] = "SYS"
	LinkTypeMap[P2PLinkSameCPU] = "NODE"
	LinkTypeMap[P2PLinkHostBridge] = "PHB"
	LinkTypeMap[P2PLinkMultiSwitch] = "PXB"
	LinkTypeMap[P2PLinkSingleSwitch] = "PIX"
	LinkTypeMap[P2PLinkSameBoard] = "X"
	LinkTypeMap[SingleNVLINKLink] = "NV1"
	LinkTypeMap[TwoNVLINKLinks] = "NV2"
	LinkTypeMap[ThreeNVLINKLinks] = "NV3"
	LinkTypeMap[FourNVLINKLinks] = "NV4"
	LinkTypeMap[FiveNVLINKLinks] = "NV5"
	LinkTypeMap[SixNVLINKLinks] = "NV6"
	LinkTypeMap[SevenNVLINKLinks] = "NV7"
	LinkTypeMap[EightNVLINKLinks] = "NV8"
	LinkTypeMap[NineNVLINKLinks] = "NV9"
	LinkTypeMap[TenNVLINKLinks] = "NV10"
	LinkTypeMap[ElevenNVLINKLinks] = "NV11"
	LinkTypeMap[TwelveNVLINKLinks] = "NV12"
}

func getP2PLink(dev1, dev2 nvml.Device) (link LinkType, err nvml.Return) {
	level, ret := nvml.DeviceGetTopologyCommonAncestor(dev1, dev2)
	if ret != nvml.SUCCESS {
		return P2PLinkUnknown, ret
	}
	switch level {
	case nvml.TOPOLOGY_INTERNAL:
		link = P2PLinkSameBoard
	case nvml.TOPOLOGY_SINGLE:
		link = P2PLinkSingleSwitch
	case nvml.TOPOLOGY_MULTIPLE:
		link = P2PLinkMultiSwitch
	case nvml.TOPOLOGY_HOSTBRIDGE:
		link = P2PLinkHostBridge
	case nvml.TOPOLOGY_NODE:
		link = P2PLinkSameCPU
	case nvml.TOPOLOGY_SYSTEM:
		link = P2PLinkCrossCPU
	default:
		err = nvml.ERROR_NOT_SUPPORTED
	}
	return
}

func deviceGetAllNvLinkRemotePciInfo(dev nvml.Device) ([]string, nvml.Return) {
	var busIds []string
	for i := 0; i < nvml.NVLINK_MAX_LINKS; i++ {
		nvLinkState, ret := nvml.DeviceGetNvLinkState(dev, i)
		if ret != nvml.SUCCESS {
			log.Printf("nvml call DeviceGetNvLinkState ret is: %d", ret)
			return nil, ret
		}
		if nvLinkState == nvml.FEATURE_ENABLED {
			remotePciInfo, ret := nvml.DeviceGetNvLinkRemotePciInfo(dev, i)
			if ret != nvml.SUCCESS {
				log.Printf("nvml call DeviceGetNvLinkRemotePciInfo ret is: %d", ret)
				return nil, ret
			}
			busId := strconv.Itoa(int(remotePciInfo.BusId[0]))
			busIds = append(busIds, busId)
		}
	}
	return busIds, nvml.ERROR_NOT_SUPPORTED
}

func getNVLink(dev1, dev2 nvml.Device) (link LinkType, err nvml.Return) {
	nvBusIds1, err := deviceGetAllNvLinkRemotePciInfo(dev1)
	if err != nvml.SUCCESS || nvBusIds1 == nil {
		return P2PLinkUnknown, err
	}
	pciInfo, _ := dev2.GetPciInfo()
	busId := strconv.Itoa(int(pciInfo.BusId[0]))
	nvLink := P2PLinkUnknown
	for _, nvBusID := range nvBusIds1 {
		if nvBusID == busId {
			switch nvLink {
			case P2PLinkUnknown:
				nvLink = SingleNVLINKLink
			case SingleNVLINKLink:
				nvLink = TwoNVLINKLinks
			case TwoNVLINKLinks:
				nvLink = ThreeNVLINKLinks
			case ThreeNVLINKLinks:
				nvLink = FourNVLINKLinks
			case FourNVLINKLinks:
				nvLink = FiveNVLINKLinks
			case FiveNVLINKLinks:
				nvLink = SixNVLINKLinks
			case SixNVLINKLinks:
				nvLink = SevenNVLINKLinks
			case SevenNVLINKLinks:
				nvLink = EightNVLINKLinks
			case EightNVLINKLinks:
				nvLink = NineNVLINKLinks
			case NineNVLINKLinks:
				nvLink = TenNVLINKLinks
			case TenNVLINKLinks:
				nvLink = ElevenNVLINKLinks
			case ElevenNVLINKLinks:
				nvLink = TwelveNVLINKLinks
			}
		}
	}
	// TODO: Handle NVSwitch semantics
	return nvLink, nvml.SUCCESS
}
