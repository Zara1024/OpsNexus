package service

import (
	"testing"

	"dodevops-api/api/cmdb/model"
)

func TestResolveHostConnectivityTarget(t *testing.T) {
	tests := []struct {
		name       string
		host       model.CmdbHost
		wantProto  string
		wantAddr   string
		wantReason string
	}{
		{
			name: "linux host prefers ssh endpoint",
			host: model.CmdbHost{
				ID:         718,
				HostName:   "opsnexus-local-verify",
				DeviceType: "linux",
				SSHIP:      "10.0.0.200",
				SSHPort:    22,
				SSHName:    "root",
				SSHKeyID:   33,
			},
			wantProto: "ssh",
			wantAddr:  "10.0.0.200:22",
		},
		{
			name: "windows host prefers rdp endpoint",
			host: model.CmdbHost{
				ID:           719,
				HostName:     "opsnexus-win",
				DeviceType:   "windows",
				RemoteDomain: "10.0.0.201",
				RDPPort:      3389,
			},
			wantProto: "rdp",
			wantAddr:  "10.0.0.201:3389",
		},
		{
			name: "windows host falls back to ssh ip for rdp target",
			host: model.CmdbHost{
				ID:         720,
				HostName:   "opsnexus-win-fallback",
				DeviceType: "windows",
				SSHIP:      "10.0.0.202",
				RDPPort:    3390,
			},
			wantProto: "rdp",
			wantAddr:  "10.0.0.202:3390",
		},
		{
			name: "missing linux ssh config returns reason",
			host: model.CmdbHost{
				ID:         721,
				HostName:   "opsnexus-linux-missing-ssh",
				DeviceType: "linux",
			},
			wantReason: "主机未配置 SSH 连接信息",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target, err := resolveHostConnectivityTarget(tt.host)
			if tt.wantReason != "" {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tt.wantReason)
				}
				if err.Error() != tt.wantReason {
					t.Fatalf("expected error %q, got %q", tt.wantReason, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
			if target.Protocol != tt.wantProto {
				t.Fatalf("expected protocol %q, got %q", tt.wantProto, target.Protocol)
			}
			if target.Address != tt.wantAddr {
				t.Fatalf("expected address %q, got %q", tt.wantAddr, target.Address)
			}
		})
	}
}
