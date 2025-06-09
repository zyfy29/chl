package feishu

import (
	"fmt"
)

type RangeData struct {
	Revision         int    `json:"revision"`
	SpreadsheetToken string `json:"spreadsheetToken"`
	ValueRange       struct {
		MajorDimension string  `json:"majorDimension"`
		Range          string  `json:"range"`
		Revision       int     `json:"revision"`
		Values         [][]any `json:"values"` // deprecated, use Data instead
		Data           [][]string
	} `json:"valueRange"`
}

func any2string(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case float64:
		return fmt.Sprintf("%.0f", v)
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", v)
	}
}

// ReadRangeData reads data from a specified range in a spreadsheet.
// The `border` parameter can be used to specify a range like "A1:B2".
// If `border` is empty, it reads the entire sheet specified by `sheetId`.
func (c *Client) ReadRangeData(spreadsheetId, sheetId, border string) (*RangeData, error) {
	range_ := sheetId
	if border != "" {
		range_ = fmt.Sprintf("%s!%s", sheetId, border)
	}
	resp, err := c.r.R().
		SetResult(&Response[RangeData]{}).
		SetPathParams(map[string]string{
			"spreadsheet_token": spreadsheetId,
			"range":             range_,
		}).
		Get("/sheets/v2/spreadsheets/{spreadsheet_token}/values/{range}")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error reading range data\nstatus: %d\nbody: %s", resp.StatusCode(), resp.String())
	}

	data := resp.Result().(*Response[RangeData])
	data.Data.ValueRange.Data = make([][]string, len(data.Data.ValueRange.Values))
	for i := range data.Data.ValueRange.Values {
		data.Data.ValueRange.Data[i] = make([]string, len(data.Data.ValueRange.Values[i]))
		for j := range data.Data.ValueRange.Values[i] {
			data.Data.ValueRange.Data[i][j] = any2string(data.Data.ValueRange.Values[i][j])
		}
	}
	return &resp.Result().(*Response[RangeData]).Data, nil
}

// WriteCellData writes data to a specific cell in a spreadsheet.
// The `cellIndex` parameter should be in the format "A1", "B2", etc.
func (c *Client) WriteCellData(spreadsheetId, sheetId, cellIndex string, data string) error {
	range_ := fmt.Sprintf("%s!%s", sheetId, cellIndex+":"+cellIndex)
	resp, err := c.r.R().
		SetBody(map[string]any{
			"valueRange": map[string]any{
				"range":  range_,
				"values": [][]string{{data}},
			},
		}).
		SetPathParams(map[string]string{
			"spreadsheet_token": spreadsheetId,
		}).
		Put("/sheets/v2/spreadsheets/{spreadsheet_token}/values")

	if err != nil {
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("error writing range data\nstatus: %d\nbody: %s", resp.StatusCode(), resp.String())
	}
	return nil
}
