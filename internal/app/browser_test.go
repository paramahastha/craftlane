package app

import "testing"

func TestBrowserFailureIsReportedToBestEffortCaller(t *testing.T) {
	t.Setenv("BROWSER", t.TempDir()+"/missing-browser")
	if err := OpenBrowser("http://127.0.0.1:7331"); err == nil {
		t.Fatal("expected browser launch failure")
	}
}
