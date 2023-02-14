package vault

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/api"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type PkiCert struct {
	Path          string
	Serial        string
	RequestID     string
	LeaseID       string
	LeaseDuration int64
	Renewable     bool
}

// Table Function
func tablePkiCert() *plugin.Table {
	return &plugin.Table{
		Name:        "vault_pki_cert",
		Description: "Vault PKI Certificates",
		List: &plugin.ListConfig{
			Hydrate: listCerts,
		},
		Columns: []*plugin.Column{
			{Name: "path", Type: proto.ColumnType_STRING, Description: "The path (mount point) of the engine containing the PKI certificate"},
			{Name: "serial", Type: proto.ColumnType_STRING, Description: "The serial identifier of the certificate"},
			{Name: "request_id", Type: proto.ColumnType_STRING, Description: "Request Identifier"},
			{Name: "lease_id", Type: proto.ColumnType_STRING, Description: "Lease Identifier", Transform: transform.FromGo()},
			{Name: "lease_duration", Type: proto.ColumnType_INT, Description: "Duration of Lease in seconds (0 [infinite] if not set)", Transform: transform.FromGo()},
			{Name: "renewable", Type: proto.ColumnType_BOOL, Description: "Indication if the certificate is renewable"},
		},
	}
}

// Hydrate Functions
func listCerts(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}
	allMounts, err := conn.Sys().ListMounts()
	if err != nil {
		return nil, err
	}

	mounts := filterMounts(ctx, allMounts, "pki", d.Quals)
	for mount := range mounts {
		certs, err := getCertDetails(ctx, conn, mount)
		if err != nil {
			return nil, err
		}
		for _, cert := range certs {
			d.StreamListItem(ctx, cert)
		}
	}
	return nil, nil
}

// Data Obtaining Functions
func getCertDetails(ctx context.Context, client *api.Client, engine string) ([]*PkiCert, error) {
	data, err := client.Logical().List(replaceDoubleSlash(fmt.Sprintf("/%s/certs", engine)))
	if err != nil {
		return nil, err
	}

	out := []*PkiCert{}

	certs := getSecretAsStrings(data)
	for _, cert := range certs {
		out = append(out, &PkiCert{
			Path:          engine,
			Serial:        cert,
			RequestID:     data.RequestID,
			LeaseID:       data.LeaseID,
			LeaseDuration: int64(data.LeaseDuration),
			Renewable:     data.Renewable,
		})
	}
	return out, nil
}
