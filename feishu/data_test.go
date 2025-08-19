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

func TestWriteCellImage(t *testing.T) {
	err := Api.WriteCellImage(
		config.Conf.Table.TableToken, config.Conf.Table.SheetID, "K2",
		[]byte{137, 80, 78, 71, 13, 10, 26, 10, 0, 0, 0, 13, 73, 72, 68, 82, 0, 0, 0, 2, 0, 0, 0, 1, 1, 0, 0, 0, 0, 220, 89, 66, 39, 0, 0, 0, 2, 116, 82, 78, 83, 0, 0, 118, 147, 205, 56, 0, 0, 0, 10, 73, 68, 65, 84, 8, 215, 99, 104, 0, 0, 0, 130, 0, 129, 221, 67, 106, 244, 0, 0, 0, 0, 73, 69, 78, 68, 174, 66, 96, 130},
		"test.png",
	)
	if err != nil {
		t.Fatalf("Failed to write cell data: %v", err)
	}
	t.Log("Data written successfully")
}
