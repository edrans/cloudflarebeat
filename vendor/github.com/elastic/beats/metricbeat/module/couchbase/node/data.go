package node

import (
	"encoding/json"
	"io"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"strconv"
)

type NodeSystemStats struct {
	CPUUtilizationRate float32 `json:"cpu_utilization_rate"`
	SwapTotal          int64   `json:"swap_total"`
	SwapUsed           int64   `json:"swap_used"`
	MemTotal           int64   `json:"mem_total"`
	MemFree            int64   `json:"mem_free"`
}

type NodeInterestingStats struct {
	CmdGet                   int64 `json:"cmd_get"`
	CouchDocsActualDiskSize  int64 `json:"couch_docs_actual_disk_size"`
	CouchDocsDataSize        int64 `json:"couch_docs_data_size"`
	CouchSpatialDataSize     int64 `json:"couch_spatial_data_size"`
	CouchSpatialDiskSize     int64 `json:"couch_spatial_disk_size"`
	CouchViewsActualDiskSize int64 `json:"couch_views_actual_disk_size"`
	CouchViewsDataSize       int64 `json:"couch_views_data_size"`
	CurrItems                int64 `json:"curr_items"`
	CurrItemsTot             int64 `json:"curr_items_tot"`
	EpBgFetched              int64 `json:"ep_bg_fetched"`
	GetHits                  int64 `json:"get_hits"`
	MemUsed                  int64 `json:"mem_used"`
	Ops                      int64 `json:"ops"`
	VbReplicaCurrItems       int64 `json:"vb_replica_curr_items"`
}

type Node struct {
	SystemStats          NodeSystemStats      `json:"systemStats"`
	InterestingStats     NodeInterestingStats `json:"interestingStats"`
	Uptime               string               `json:"uptime"`
	MemoryTotal          int64                `json:"memoryTotal"`
	MemoryFree           int64                `json:"memoryFree"`
	McdMemoryReserved    int64                `json:"mcdMemoryReserved"`
	McdMemoryAllocated   int64                `json:"mcdMemoryAllocated"`
	ClusterMembership    string               `json:"clusterMembership"`
	RecoveryType         string               `json:"recoveryType"`
	Status               string               `json:"status"`
	ThisNode             bool                 `json:"thisNode"`
	Hostname             string               `json:"hostname"`
	ClusterCompatibility int64                `json:"clusterCompatibility"`
	Version              string               `json:"version"`
	Os                   string               `json:"os"`
}

type Data struct {
	Nodes []Node `json:"nodes"`
}

func eventsMapping(body io.Reader) []common.MapStr {

	var d Data
	err := json.NewDecoder(body).Decode(&d)
	if err != nil {
		logp.Err("Error: ", err)
	}

	events := []common.MapStr{}

	for _, NodeItem := range d.Nodes {
		uptime, _ := strconv.ParseInt(NodeItem.Uptime, 10, 64)

		event := common.MapStr{
			"cmd_get": NodeItem.InterestingStats.CmdGet,
			"couch": common.MapStr{
				"docs": common.MapStr{
					"actual_disk_size": common.MapStr{
						"bytes": NodeItem.InterestingStats.CouchDocsActualDiskSize,
					},
					"data_size": common.MapStr{
						"bytes": NodeItem.InterestingStats.CouchDocsDataSize,
					},
				},
				"spacial": common.MapStr{
					"data_size": common.MapStr{
						"bytes": NodeItem.InterestingStats.CouchSpatialDataSize,
					},
					"disk_size": common.MapStr{
						"bytes": NodeItem.InterestingStats.CouchSpatialDiskSize,
					},
				},
				"views": common.MapStr{
					"actual_disk_size": common.MapStr{
						"bytes": NodeItem.InterestingStats.CouchViewsActualDiskSize,
					},
					"data_size": common.MapStr{
						"bytes": NodeItem.InterestingStats.CouchViewsDataSize,
					},
				},
			},
			"cpu_utilization_rate": common.MapStr{
				"pct": NodeItem.SystemStats.CPUUtilizationRate,
			},
			"current_items": common.MapStr{
				"value": NodeItem.InterestingStats.CurrItems,
				"total": NodeItem.InterestingStats.CurrItemsTot,
			},
			"ep_bg_fetched": NodeItem.InterestingStats.EpBgFetched,
			"get_hits":      NodeItem.InterestingStats.GetHits,
			"hostname":      NodeItem.Hostname,
			"mcd_memory": common.MapStr{
				"reserved": common.MapStr{
					"bytes": NodeItem.McdMemoryReserved,
				},
				"allocated": common.MapStr{
					"bytes": NodeItem.McdMemoryAllocated,
				},
			},
			"memory": common.MapStr{
				"total": common.MapStr{
					"bytes": NodeItem.SystemStats.MemTotal,
				},
				"free": common.MapStr{
					"bytes": NodeItem.SystemStats.MemFree,
				},
				"used": common.MapStr{
					"bytes": NodeItem.InterestingStats.MemUsed,
				},
			},
			"ops": NodeItem.InterestingStats.Ops,
			"swap": common.MapStr{
				"total": common.MapStr{
					"bytes": NodeItem.SystemStats.SwapTotal,
				},
				"used": common.MapStr{
					"bytes": NodeItem.SystemStats.SwapUsed,
				},
			},
			"uptime": common.MapStr{
				"sec": uptime,
			},
			"vb_replica_curr_items": NodeItem.InterestingStats.VbReplicaCurrItems,
		}
		events = append(events, event)
	}

	return events

}
