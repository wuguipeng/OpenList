package _139

import (
	"strconv"
	"testing"
)

func TestShareEntriesAndRefEncoding(t *testing.T) {
	d := &Yun139{Addition: Addition{
		Type:   MetaShare,
		LinkID: "share-a,share-b,share-c#pass",
	}}

	entries := d.shareEntries()
	if len(entries) != 3 {
		t.Fatalf("expected 3 share entries, got %d", len(entries))
	}
	if entries[2].LinkID != "share-c" || entries[2].Password != "pass" {
		t.Fatalf("unexpected password share entry: %+v", entries[2])
	}

	refs := []shareRef{
		{LinkID: entries[0].LinkID, Password: entries[0].Password, NodeID: "root-a"},
		{LinkID: entries[2].LinkID, Password: entries[2].Password, NodeID: "root-b"},
	}
	encoded := encodeShareRefs(refs)
	decoded, ok := decodeShareRefs(encoded)
	if !ok || len(decoded) != len(refs) {
		t.Fatalf("failed to decode merged share refs: %q", encoded)
	}
	for i := range refs {
		if decoded[i] != refs[i] {
			t.Fatalf("unexpected decoded ref at %d: %+v", i, decoded[i])
		}
	}
}

func TestShareEntriesUrlParsing(t *testing.T) {
	d := &Yun139{Addition: Addition{
		Type: MetaShare,
		LinkID: "https://yun.139.com/shareweb/#/w/i/1A52003541jdf\n" +
			"https://yun.139.com/shareweb/#/w/i/1B52003541abc#pwd123\n" +
			"https://yun.139.com/w/#/detail?linkID=1C52003541xyz\n" +
			"https://yun.139.com/w/#/detail?linkID=1D52003541def#pwd456\n" +
			"rawLink1#pass1\n" +
			"rawLink2",
	}}

	entries := d.shareEntries()
	if len(entries) != 6 {
		t.Fatalf("expected 6 entries, got %d: %+v", len(entries), entries)
	}

	expected := []struct{ LinkID, Password string }{
		{"1A52003541jdf", ""},
		{"1B52003541abc", "pwd123"},
		{"1C52003541xyz", ""},
		{"1D52003541def", "pwd456"},
		{"rawLink1", "pass1"},
		{"rawLink2", ""},
	}

	for i, exp := range expected {
		if entries[i].LinkID != exp.LinkID || entries[i].Password != exp.Password {
			t.Errorf("entry %d mismatch: got %+v, want %+v", i, entries[i], exp)
		}
	}
}

func TestLargeShareRefsCacheAndDeduplication(t *testing.T) {
	// Test deduplication
	duplicateRefs := []shareRef{
		{LinkID: "linkA", Password: "pwd", NodeID: "node1"},
		{LinkID: "linkA", Password: "pwd", NodeID: "node1"},
		{LinkID: "linkB", Password: "", NodeID: "node2"},
	}
	dedup := deduplicateShareRefs(duplicateRefs)
	if len(dedup) != 2 {
		t.Fatalf("expected 2 deduped refs, got %d", len(dedup))
	}

	// Test large number of refs (cache storage)
	refs := make([]shareRef, 100)
	for i := 0; i < 100; i++ {
		refs[i] = shareRef{
			LinkID:   "share-link-id-" + strconv.Itoa(i) + "-000000000000000000",
			Password: "password-123456",
			NodeID:   "catalog-node-id-" + strconv.Itoa(i),
		}
	}
	encoded := encodeShareRefs(refs)
	decoded, ok := decodeShareRefs(encoded)
	if !ok {
		t.Fatalf("failed to decode cached large share refs: %q", encoded)
	}
	if len(decoded) != len(refs) {
		t.Fatalf("expected %d decoded refs, got %d", len(refs), len(decoded))
	}
}

func TestShareHasAnyCredentials(t *testing.T) {
	tests := []struct {
		name string
		d    Yun139
		want bool
	}{
		{
			name: "empty",
			d:    Yun139{},
			want: false,
		},
		{
			name: "authorization only",
			d:    Yun139{Addition: Addition{Authorization: "auth"}},
			want: true,
		},
		{
			name: "mail cookies only",
			d:    Yun139{Addition: Addition{MailCookies: "RMKEY=rm; Os_SSo_Sid=sid"}},
			want: true,
		},
		{
			name: "username and password",
			d:    Yun139{Addition: Addition{Username: "u", Password: "p"}},
			want: true,
		},
		{
			name: "username only",
			d:    Yun139{Addition: Addition{Username: "u"}},
			want: false,
		},
		{
			name: "full login credentials",
			d: Yun139{Addition: Addition{
				MailCookies: "RMKEY=rm; Os_SSo_Sid=sid",
				Username:    "u",
				Password:    "p",
			}},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.hasAnyCredentials(); got != tt.want {
				t.Fatalf("hasAnyCredentials() = %v, want %v", got, tt.want)
			}
		})
	}
}
