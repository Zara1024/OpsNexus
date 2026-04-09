package dao

import "testing"

func TestTruncateSysLoginInfoSQLUsesCanonicalTableName(t *testing.T) {
	got := truncateSysLoginInfoSQL()
	want := "truncate table sys_login_info"
	if got != want {
		t.Fatalf("truncate sql mismatch: got %q want %q", got, want)
	}
}
