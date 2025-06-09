package feishu

import (
	"chl/config"
	"testing"
)

func TestReadRangeData(t *testing.T) {

	data, err := Api.ReadRangeData(config.Conf.Table.TableToken, config.Conf.Table.SheetID, "")
	if err != nil {
		t.Fatalf("Failed to read range data: %v", err)
	}

	if len(data.ValueRange.Values) == 0 {
		t.Fatal("No data returned in the specified range")
	}

	t.Log("Data read successfully")
	for i, row := range data.ValueRange.Values {
		t.Logf("Row %d: %v", i+1, row)
	}
}
