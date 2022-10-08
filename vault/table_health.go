package vault

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

type SysHealth struct {
	Initialized                bool
	Sealed                     bool
	Standby                    bool
	PerformanceStandby         bool
	ReplicationPerformanceMode string
	ReplicationDrMode          string
	ServerTimeUtc              int64
	Version                    string
	ClusterName                string
	ClusterID                  string
}

func tableSysHealth() *plugin.Table {
	return &plugin.Table{
		Name:        "vault_sys_health",
		Description: "Health of Vault",
		List: &plugin.ListConfig{
			Hydrate: getSysHealth,
		},
		Columns: []*plugin.Column{
			{Name: "initialized", Type: proto.ColumnType_BOOL, Description: "Is Initialized"},
			{Name: "sealed", Type: proto.ColumnType_BOOL, Description: "Is sealed"},
			{Name: "standby", Type: proto.ColumnType_BOOL, Description: "Is in standby"},
			{Name: "performance_standby", Type: proto.ColumnType_BOOL, Description: "Is Performance Standby"},
			{Name: "replication_performance_mode", Type: proto.ColumnType_STRING, Description: "Replication Performance Mode"},
			{Name: "replication_dr_mode", Type: proto.ColumnType_STRING, Description: "Replication Disaster Recovery Mode"},
			{Name: "server_time_utc", Type: proto.ColumnType_TIMESTAMP, Description: "Server Time in UTC", Transform: transform.FromField("ServerTimeUtc").Transform(convertTimestamp)},
			{Name: "version", Type: proto.ColumnType_STRING, Description: "Hashicorp Vault Version"},
			{Name: "cluster_name", Type: proto.ColumnType_STRING, Description: "Name of Vault Cluster"},
			{Name: "cluster_id", Type: proto.ColumnType_STRING, Description: "Identity of Vault Cluster"},
		},
	}
}

func getSysHealth(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)

	if err != nil {
		return nil, err
	}

	data, err := conn.Sys().Health()

	if err != nil {
		return nil, err
	}

	d.StreamListItem(ctx, &SysHealth{
		Initialized:                data.Initialized,
		Sealed:                     data.Sealed,
		Standby:                    data.Standby,
		PerformanceStandby:         data.PerformanceStandby,
		ReplicationPerformanceMode: data.ReplicationPerformanceMode,
		ReplicationDrMode:          data.ReplicationDRMode,
		ServerTimeUtc:              data.ServerTimeUTC,
		Version:                    data.Version,
		ClusterName:                data.ClusterName,
		ClusterID:                  data.ClusterID,
	})

	return nil, nil
}
