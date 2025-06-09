package feishu

import (
	"chl/config"
	"testing"
)

func TestGetSheets(t *testing.T) {
	sheets, err := Api.GetSheets(config.Conf.Table.TableToken)
	if err != nil {
		t.Fatalf("GetSheets failed: %v", err)
	}
	if len(sheets) == 0 {
		t.Error("Expected at least one sheet, got none")
	} else {
		t.Logf("Retrieved %d sheets", len(sheets))
		for _, sheet := range sheets {
			t.Logf("Sheet ID: %s, Title: %s", sheet.SheetId, sheet.Title)
		}
	}
}
