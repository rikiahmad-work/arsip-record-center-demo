package seeders

import (
	"testing"

	"github.com/kemenpppa/siarsip-api/internal/model"
)

func TestMasterKodeKlasifikasi_Phase0Criteria(t *testing.T) {
	if len(masterKodeKlasifikasi) != 435 {
		t.Fatalf("expected 435 kode klasifikasi, got %d", len(masterKodeKlasifikasi))
	}

	seen := make(map[string]bool, 435)
	for i, item := range masterKodeKlasifikasi {
		if item.Kode == "" {
			t.Fatalf("row %d has empty kode", i)
		}
		if seen[item.Kode] {
			t.Fatalf("duplicate kode found: %s", item.Kode)
		}
		seen[item.Kode] = true

		if !model.SKKAD(item.SKKAD).IsValid() {
			t.Fatalf("kode %s has invalid SKKAD: %s", item.Kode, item.SKKAD)
		}

		if item.NasibAkhir != "" && !model.NasibAkhir(item.NasibAkhir).IsValid() {
			t.Fatalf("kode %s has invalid NasibAkhir: %s", item.Kode, item.NasibAkhir)
		}
	}
}

func TestCanonicalEnums_Phase0Criteria(t *testing.T) {
	roles := model.AllRoles()
	if len(roles) != 4 {
		t.Fatalf("expected 4 canonical roles, got %d", len(roles))
	}
	for _, r := range roles {
		if !r.IsValid() {
			t.Fatalf("role %s should be valid", r)
		}
	}

	if model.NormalizeRole("super-admin") != model.RoleAdministrator {
		t.Fatalf("super-admin should map to administrator")
	}
	if model.NormalizeRole("admin-arsip") != model.RoleOperator {
		t.Fatalf("admin-arsip should map to operator")
	}
}
