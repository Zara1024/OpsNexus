package service

import (
	"testing"

	"dodevops-api/api/cmdb/model"
)

func TestNormalizeCreateCmdbDeviceRequestAppliesProtocolDefaults(t *testing.T) {
	dto := model.CreateCmdbDeviceDto{
		Name:          " core-switch-01 ",
		Address:       " 10.10.10.10 ",
		Platform:      " ios ",
		GroupID:       7,
		AccountID:     9,
		ProtocolGroup: " ssh ",
		Tags:          " core,edge ",
		DeviceType:    " switch ",
	}

	normalized := normalizeCreateCmdbDeviceRequest(dto)

	if normalized.Name != "core-switch-01" {
		t.Fatalf("expected trimmed name, got %q", normalized.Name)
	}
	if normalized.Address != "10.10.10.10" {
		t.Fatalf("expected trimmed address, got %q", normalized.Address)
	}
	if normalized.Platform != "ios" {
		t.Fatalf("expected trimmed platform, got %q", normalized.Platform)
	}
	if normalized.ProtocolGroup != "ssh" {
		t.Fatalf("expected normalized protocol group ssh, got %q", normalized.ProtocolGroup)
	}
	if normalized.DeviceType != "switch" {
		t.Fatalf("expected normalized device type switch, got %q", normalized.DeviceType)
	}
	if normalized.SSHPort != 22 {
		t.Fatalf("expected default ssh port 22, got %d", normalized.SSHPort)
	}
	if normalized.TelnetPort != 23 {
		t.Fatalf("expected default telnet port 23, got %d", normalized.TelnetPort)
	}
	if !normalized.IsActive {
		t.Fatalf("expected isActive to default to true")
	}
}

func TestResolveCmdbDeviceConnectivityTargetPrefersSSHThenTelnetThenWeb(t *testing.T) {
	tests := []struct {
		name       string
		device     model.CmdbDevice
		wantProto  string
		wantAddr   string
		wantReason string
	}{
		{
			name: "ssh is preferred when configured",
			device: model.CmdbDevice{
				Name:       "router-01",
				Address:    "10.0.0.10",
				SSHPort:    22,
				TelnetPort: 23,
				WebURL:     "https://10.0.0.10",
			},
			wantProto: "ssh",
			wantAddr:  "10.0.0.10:22",
		},
		{
			name: "telnet is used when ssh is unavailable",
			device: model.CmdbDevice{
				Name:       "switch-02",
				Address:    "10.0.0.11",
				TelnetPort: 2323,
				WebURL:     "http://10.0.0.11",
			},
			wantProto: "telnet",
			wantAddr:  "10.0.0.11:2323",
		},
		{
			name: "web is used when only web url is available",
			device: model.CmdbDevice{
				Name:    "firewall-03",
				WebURL:  "https://fw.example.com",
				Address: "10.0.0.12",
			},
			wantProto: "web",
			wantAddr:  "fw.example.com:443",
		},
		{
			name: "missing endpoints returns reason",
			device: model.CmdbDevice{
				Name: "device-04",
			},
			wantReason: "device has no reachable connectivity target configured",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target, err := resolveCmdbDeviceConnectivityTarget(tt.device)
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

func TestNormalizeCmdbDeviceProtocolGroupCanonicalizesWhitespaceAndOrder(t *testing.T) {
	got := normalizeCmdbDeviceProtocolGroup(" web, SSH , telnet,ssh ,unknown ")

	if got != "ssh,telnet,web" {
		t.Fatalf("expected canonical protocol group %q, got %q", "ssh,telnet,web", got)
	}
}

func TestValidateCmdbDeviceRejectsOutOfRangePorts(t *testing.T) {
	device := model.CmdbDevice{
		Name:       "dist-sw-01",
		Address:    "10.0.0.20",
		GroupID:    1,
		AccountID:  2,
		SSHPort:    65536,
		TelnetPort: 23,
	}

	err := validateCmdbDevice(device)
	if err == nil {
		t.Fatal("expected out-of-range ssh port to be rejected")
	}
}

func TestNormalizeCmdbDevicePageSizeCapsAtReasonableMaximum(t *testing.T) {
	if got := normalizeCmdbDevicePageSize(0); got != 10 {
		t.Fatalf("expected default page size 10, got %d", got)
	}
	if got := normalizeCmdbDevicePageSize(500); got != 100 {
		t.Fatalf("expected capped page size 100, got %d", got)
	}
}

func TestValidateCmdbDeviceConnectivityBatchSizeRejectsOversizedRequests(t *testing.T) {
	deviceIDs := make([]uint, cmdbDeviceConnectivityBatchLimit+1)

	err := validateCmdbDeviceConnectivityBatchSize(deviceIDs)
	if err == nil {
		t.Fatal("expected oversized connectivity batch to be rejected")
	}
}

func TestNormalizeCmdbDeviceConnectivityIDsRemovesZeroAndPreservesOrder(t *testing.T) {
	got := normalizeCmdbDeviceConnectivityIDs([]uint{0, 7, 7, 0, 3, 9, 3, 9, 11})
	want := []uint{7, 3, 9, 11}

	if len(got) != len(want) {
		t.Fatalf("expected %d ids, got %d: %v", len(want), len(got), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected ids %v, got %v", want, got)
		}
	}
}

func TestValidateCmdbDeviceConnectivityBatchSizeUsesEffectiveDedupedCount(t *testing.T) {
	repeated := make([]uint, cmdbDeviceConnectivityBatchLimit+50)
	for i := range repeated {
		repeated[i] = 42
	}

	if err := validateCmdbDeviceConnectivityBatchSize(normalizeCmdbDeviceConnectivityIDs(repeated)); err != nil {
		t.Fatalf("expected repeated ids to pass effective batch validation, got %v", err)
	}

	unique := make([]uint, cmdbDeviceConnectivityBatchLimit+1)
	for i := range unique {
		unique[i] = uint(i + 1)
	}
	if err := validateCmdbDeviceConnectivityBatchSize(normalizeCmdbDeviceConnectivityIDs(unique)); err == nil {
		t.Fatal("expected unique oversized batch to be rejected")
	}
}
