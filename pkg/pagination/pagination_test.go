package pagination

import "testing"

func TestNormalizeAppliesDefaultsAndCaps(t *testing.T) {
	cases := []struct {
		name     string
		in       Query
		wantPage int
		wantSize int
	}{
		{"zero values", Query{}, 1, defaultPageSize},
		{"negative page", Query{Page: -3, PageSize: 20}, 1, 20},
		{"oversized page size", Query{Page: 2, PageSize: maxPageSize + 100}, 2, maxPageSize},
		{"valid input untouched", Query{Page: 5, PageSize: 30}, 5, 30},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			q := tc.in
			q.Normalize()

			if q.Page != tc.wantPage {
				t.Errorf("page = %d, want %d", q.Page, tc.wantPage)
			}
			if q.PageSize != tc.wantSize {
				t.Errorf("page size = %d, want %d", q.PageSize, tc.wantSize)
			}
		})
	}
}

func TestOffsetAndLimit(t *testing.T) {
	q := Query{Page: 3, PageSize: 20}
	q.Normalize()

	if got := q.Offset(); got != 40 {
		t.Errorf("offset = %d, want 40", got)
	}
	if got := q.Limit(); got != 20 {
		t.Errorf("limit = %d, want 20", got)
	}
}
