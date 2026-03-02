package util

import "testing"

func TestSanitizeFilename(t *testing.T) {
	t.Parallel()

	cases := []struct {
		in     string
		maxLen int
		want   string
	}{
		{in: "my file.pdf", maxLen: 100, want: "my_file.pdf"},
		{in: " ../etc/passwd ", maxLen: 100, want: "etc_passwd"},
		{in: "a/b\\c?.txt", maxLen: 100, want: "a_b_c.txt"},
		{in: "----", maxLen: 100, want: "file"},
		{in: "report.final.v2.PDF", maxLen: 100, want: "report_final_v2.PDF"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.in, func(t *testing.T) {
			t.Parallel()
			got := SanitizeFilename(tc.in, tc.maxLen)
			if got != tc.want {
				t.Fatalf("SanitizeFilename(%q, %d) = %q, want %q", tc.in, tc.maxLen, got, tc.want)
			}
		})
	}
}

func TestSanitizeFilename_MaxLen(t *testing.T) {
	t.Parallel()

	long := ""
	for i := 0; i < 200; i++ {
		long += "a"
	}

	got := SanitizeFilename(long+".pdf", 20)
	if len(got) != 20 {
		t.Fatalf("expected len=20, got=%d (%q)", len(got), got)
	}
	if got[len(got)-4:] != ".pdf" {
		t.Fatalf("expected to keep extension, got=%q", got)
	}
}
