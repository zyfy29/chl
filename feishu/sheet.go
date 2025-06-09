package feishu

import (
	"fmt"
)

type SheetItem struct {
	SheetId        string `json:"sheet_id"`
	Title          string `json:"title"`
	Index          int    `json:"index"`
	Hidden         bool   `json:"hidden"`
	GridProperties struct {
		FrozenRowCount    int `json:"frozen_row_count"`
		FrozenColumnCount int `json:"frozen_column_count"`
		RowCount          int `json:"row_count"`
		ColumnCount       int `json:"column_count"`
	} `json:"grid_properties"`
	ResourceType string `json:"resource_type"`
	Merges       []struct {
		StartRowIndex    int `json:"start_row_index"`
		EndRowIndex      int `json:"end_row_index"`
		StartColumnIndex int `json:"start_column_index"`
		EndColumnIndex   int `json:"end_column_index"`
	} `json:"merges"`
}

func (c *Client) GetSheets(spreadsheetId string) ([]SheetItem, error) {
	type responseData Response[struct {
		Sheets []SheetItem `json:"sheets"`
	}]
	resp, err := c.r.R().
		SetResult(&responseData{}).
		SetPathParams(map[string]string{
			"spreadsheet_token": spreadsheetId,
		}).
		Get("/sheets/v3/spreadsheets/{spreadsheet_token}/sheets/query")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error fetching sheets\nstatus: %d\nbody: %s", resp.StatusCode(), resp.String())
	}

	result := resp.Result().(*responseData)
	return result.Data.Sheets, nil
}
