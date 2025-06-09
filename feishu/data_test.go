package feishu

import (
	"github.com/zyfy29/chl/config"
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

func TestWriteCellData(t *testing.T) {
	err := Api.WriteCellData(config.Conf.Table.TableToken, config.Conf.Table.SheetID, "I2", "Hello from golang!")
	if err != nil {
		t.Fatalf("Failed to write cell data: %v", err)
	}
	t.Log("Data written successfully")
}
