package alphavantage

import "testing"

func TestGetClosingData(t *testing.T) {
	data := GetClosingData("MSFT", 7)
	if len(data) != 7 {
		t.Fatalf("Wrong number of days returned")
	}
}
