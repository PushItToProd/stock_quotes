package alphavantage

import "testing"

// Quick test to validate that the query doesn't fail.
func TestGetClosingData(t *testing.T) {
	t.Logf("Getting closing data")
	data, err := GetClosingData("MSFT", 7)
	t.Logf("Got closing data")
	if err != nil {
		t.Fatalf("Got error: %v", err)
	}
	if len(data) != 7 {
		t.Fatalf("Wrong number of days of data returned - got data: %v", data)
	}
}
