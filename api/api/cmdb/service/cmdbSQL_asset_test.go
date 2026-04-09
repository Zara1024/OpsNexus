package service

import (
	"testing"

	"dodevops-api/api/cmdb/model"
)

func TestNormalizeCreateCmdbSQLRequestPreservesCompatibilityFields(t *testing.T) {
	dto := model.CreateCmdbSQLDto{
		Name:          " analytics-prod ",
		Address:       " 10.0.0.15 ",
		Platform:      " mysql 8 ",
		AccountID:     7,
		GroupID:       4,
		ProtocolGroup: " mysql ",
		Tags:          " prod,reporting ",
		Type:          1,
		Description:   " legacy remark ",
	}

	normalized := normalizeCreateCmdbSQLRequest(dto)

	if normalized.Name != "analytics-prod" {
		t.Fatalf("expected trimmed name, got %q", normalized.Name)
	}
	if normalized.Address != "10.0.0.15" {
		t.Fatalf("expected trimmed address, got %q", normalized.Address)
	}
	if normalized.Platform != "mysql 8" {
		t.Fatalf("expected trimmed platform, got %q", normalized.Platform)
	}
	if normalized.Type != 1 {
		t.Fatalf("expected type to be preserved, got %d", normalized.Type)
	}
	if normalized.AccountID != 7 {
		t.Fatalf("expected accountId to be preserved, got %d", normalized.AccountID)
	}
	if normalized.GroupID != 4 {
		t.Fatalf("expected groupId to be preserved, got %d", normalized.GroupID)
	}
	if normalized.DefaultDatabase != "analytics-prod" {
		t.Fatalf("expected defaultDatabase to fall back to name, got %q", normalized.DefaultDatabase)
	}
	if normalized.ProtocolGroup != "mysql" {
		t.Fatalf("expected normalized protocol group mysql, got %q", normalized.ProtocolGroup)
	}
	if normalized.Tags != "prod,reporting" {
		t.Fatalf("expected trimmed tags, got %q", normalized.Tags)
	}
	if normalized.Remark != "legacy remark" {
		t.Fatalf("expected description compatibility to populate remark, got %q", normalized.Remark)
	}
	if normalized.Description != "legacy remark" {
		t.Fatalf("expected legacy description alias to be preserved, got %q", normalized.Description)
	}
	if !normalized.IsActive {
		t.Fatal("expected isActive to default to true")
	}
}

func TestNormalizeUpdateCmdbSQLRequestPreservesDefaultDatabaseAndIsActive(t *testing.T) {
	existing := model.CmdbSQL{
		ID:              99,
		Name:            "legacy-instance",
		Address:         "10.0.0.1",
		Platform:        "mysql",
		Type:            1,
		AccountID:       5,
		GroupID:         6,
		DefaultDatabase: "legacy_schema",
		ProtocolGroup:   "mysql",
		Tags:            "old",
		IsActive:        false,
		Remark:          "old note",
		Description:     "old note",
	}

	dto := model.UpdateCmdbSQLDto{
		ID:        99,
		Name:      " legacy-instance ",
		Address:   " 10.0.0.2 ",
		Platform:  " mysql 8 ",
		Type:      1,
		AccountID: 8,
		GroupID:   10,
		Tags:      " new ",
	}

	normalized := normalizeUpdateCmdbSQLRequest(existing, dto)

	if normalized.ID != existing.ID {
		t.Fatalf("expected id %d, got %d", existing.ID, normalized.ID)
	}
	if normalized.DefaultDatabase != "legacy_schema" {
		t.Fatalf("expected existing defaultDatabase to be preserved, got %q", normalized.DefaultDatabase)
	}
	if normalized.IsActive {
		t.Fatal("expected existing isActive=false to be preserved when omitted")
	}
	if normalized.ProtocolGroup != "mysql" {
		t.Fatalf("expected existing protocolGroup to be preserved, got %q", normalized.ProtocolGroup)
	}
	if normalized.Remark != "old note" {
		t.Fatalf("expected existing remark to be preserved, got %q", normalized.Remark)
	}
	if normalized.Description != "old note" {
		t.Fatalf("expected legacy description alias to be preserved, got %q", normalized.Description)
	}
}

func TestResolveCmdbSQLSchemaNamePrefersExplicitDefaultThenName(t *testing.T) {
	asset := model.CmdbSQL{
		Name:            "instance-name",
		DefaultDatabase: "default_schema",
	}

	if got := resolveCmdbSQLSchemaName(asset, " tenant_schema "); got != "tenant_schema" {
		t.Fatalf("expected explicit schema to win, got %q", got)
	}
	if got := resolveCmdbSQLSchemaName(asset, " "); got != "default_schema" {
		t.Fatalf("expected defaultDatabase fallback, got %q", got)
	}

	asset.DefaultDatabase = ""
	if got := resolveCmdbSQLSchemaName(asset, " "); got != "instance-name" {
		t.Fatalf("expected name fallback when defaultDatabase empty, got %q", got)
	}
}

func TestResolveCmdbSQLAssetLookupSupportsDatabaseNameOnly(t *testing.T) {
	asset, err := resolveCmdbSQLAssetLookup(
		0,
		" analytics-prod ",
		func(id uint) (*model.CmdbSQL, error) {
			t.Fatalf("expected id lookup to be skipped, got id=%d", id)
			return nil, nil
		},
		func(name string) ([]model.CmdbSQL, error) {
			if name != "analytics-prod" {
				t.Fatalf("expected trimmed lookup name, got %q", name)
			}
			return []model.CmdbSQL{{
				ID:              7,
				Name:            "analytics-prod",
				Type:            1,
				AccountID:       22,
				DefaultDatabase: "analytics",
				IsActive:        true,
				Remark:          "prod asset",
			}}, nil
		},
	)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if asset == nil {
		t.Fatal("expected resolved asset, got nil")
	}
	if asset.AccountID != 22 {
		t.Fatalf("expected actual accountId to be preserved, got %d", asset.AccountID)
	}
	if asset.Type != 1 {
		t.Fatalf("expected actual type to be preserved, got %d", asset.Type)
	}
	if asset.Description != "prod asset" {
		t.Fatalf("expected legacy description compatibility on resolved asset, got %q", asset.Description)
	}
}

func TestResolveCmdbSQLTargetUsesAssetDefaultDatabaseForNameOnlyLookup(t *testing.T) {
	asset, schemaName, err := resolveCmdbSQLTarget(
		0,
		" reporting-asset ",
		func(id uint) (*model.CmdbSQL, error) {
			t.Fatalf("expected id lookup to be skipped, got id=%d", id)
			return nil, nil
		},
		func(name string) ([]model.CmdbSQL, error) {
			if name != "reporting-asset" {
				t.Fatalf("expected trimmed asset lookup name, got %q", name)
			}
			return []model.CmdbSQL{{
				ID:              8,
				Name:            "reporting-asset",
				DefaultDatabase: "reporting_schema",
				AccountID:       33,
				Type:            1,
				IsActive:        true,
			}}, nil
		},
	)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if asset == nil {
		t.Fatal("expected resolved asset, got nil")
	}
	if schemaName != "reporting_schema" {
		t.Fatalf("expected defaultDatabase schema fallback, got %q", schemaName)
	}
}

func TestResolveCmdbSQLAssetLookupRejectsAmbiguousDatabaseName(t *testing.T) {
	asset, err := resolveCmdbSQLAssetLookup(
		0,
		"shared-name",
		func(id uint) (*model.CmdbSQL, error) {
			t.Fatalf("expected id lookup to be skipped, got id=%d", id)
			return nil, nil
		},
		func(name string) ([]model.CmdbSQL, error) {
			return []model.CmdbSQL{
				{ID: 1, Name: "shared-name", AccountID: 11, Type: 1, IsActive: true},
				{ID: 2, Name: "shared-name", AccountID: 12, Type: 1, IsActive: true},
			}, nil
		},
	)
	if err == nil {
		t.Fatal("expected ambiguous name lookup to fail")
	}
	if asset != nil {
		t.Fatalf("expected no asset on ambiguity, got %+v", asset)
	}
	if err.Error() != `database asset "shared-name" is ambiguous` {
		t.Fatalf("expected ambiguity error, got %q", err.Error())
	}
}

func TestResolveCmdbSQLAssetLookupRejectsInactiveAsset(t *testing.T) {
	asset, err := resolveCmdbSQLAssetLookup(
		0,
		"archive-db",
		func(id uint) (*model.CmdbSQL, error) {
			t.Fatalf("expected id lookup to be skipped, got id=%d", id)
			return nil, nil
		},
		func(name string) ([]model.CmdbSQL, error) {
			return []model.CmdbSQL{{
				ID:              5,
				Name:            "archive-db",
				DefaultDatabase: "archive",
				AccountID:       19,
				Type:            1,
				IsActive:        false,
			}}, nil
		},
	)
	if err == nil {
		t.Fatal("expected inactive asset lookup to fail")
	}
	if asset != nil {
		t.Fatalf("expected no asset for inactive lookup, got %+v", asset)
	}
	if err.Error() != `database asset "archive-db" is inactive` {
		t.Fatalf("expected inactive error, got %q", err.Error())
	}
}

func TestNormalizeCmdbSQLPageSizeCapsAndDefaults(t *testing.T) {
	if got := normalizeCmdbSQLPage(0); got != 1 {
		t.Fatalf("expected page 0 to normalize to 1, got %d", got)
	}
	if got := normalizeCmdbSQLPageSize(0); got != 10 {
		t.Fatalf("expected pageSize 0 to normalize to 10, got %d", got)
	}
	if got := normalizeCmdbSQLPageSize(500); got != 100 {
		t.Fatalf("expected oversized pageSize to clamp to 100, got %d", got)
	}
}
