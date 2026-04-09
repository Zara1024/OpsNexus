package service

import (
	"testing"
)

func TestBuildAgentMetricsSnapshotParsesCurrentValues(t *testing.T) {
	snapshot, ok := buildAgentMetricsSnapshot(`system_cpu_usage_percent{cpu="cpu0",instance="ubuntu2204.wang.org"} 0.5
system_cpu_usage_percent{cpu="cpu1",instance="ubuntu2204.wang.org"} 0.7
system_memory_usage_percent{instance="ubuntu2204.wang.org"} 8.37
system_disk_usage_percent{instance="ubuntu2204.wang.org",mountpoint="/"} 25.02
system_network_receive_kb_per_second{instance="ubuntu2204.wang.org"} 0.01
system_network_send_kb_per_second{instance="ubuntu2204.wang.org"} 0.02
system_load_average{instance="ubuntu2204.wang.org",period="1min"} 0.19
system_load_average{instance="ubuntu2204.wang.org",period="5min"} 0.14
system_load_average{instance="ubuntu2204.wang.org",period="15min"} 0.10
system_disk_read_kb_per_second{device="vda",instance="ubuntu2204.wang.org"} 0
system_disk_write_kb_per_second{device="vda",instance="ubuntu2204.wang.org"} 560.4
system_total_processes{instance="ubuntu2204.wang.org"} 542`)

	if !ok {
		t.Fatalf("expected snapshot to parse successfully")
	}
	if snapshot.CPU <= 0 {
		t.Fatalf("expected parsed cpu value, got %v", snapshot.CPU)
	}
	if snapshot.Memory != 8.37 {
		t.Fatalf("expected parsed memory value 8.37, got %v", snapshot.Memory)
	}
	if snapshot.Disk != 25.02 {
		t.Fatalf("expected parsed disk value 25.02, got %v", snapshot.Disk)
	}
	if snapshot.TotalProcesses != 542 {
		t.Fatalf("expected parsed process count 542, got %v", snapshot.TotalProcesses)
	}
}
