package nag

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dylantientcheu/nbacli/nag/params"
)

// CommonPlayerInfo wraps request to and response from commonplayerinfo endpoint.
type CommonPlayerInfo struct {
	*Client
	Response *Response
	PlayerID string
	LeagueID string
}

// NewCommonPlayerInfo creates a default CommonPlayerInfo instance.
func NewCommonPlayerInfo(id string) *CommonPlayerInfo {
	return &CommonPlayerInfo{
		Client:   NewDefaultClient(),
		PlayerID: id,
		LeagueID: params.LeagueID.Default(),
	}
}

// Get sends a GET request to commonplayerinfo endpoint.
func (c *CommonPlayerInfo) Get() error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/commonplayerinfo", c.BaseURL.String()), nil)
	// fmt.Println(fmt.Sprintf("%s/commonplayerinfo", c.BaseURL.String()))
	if err != nil {
		return err
	}

	req.Header = DefaultStatsHeader

	q := req.URL.Query()
	q.Add("PlayerID", c.PlayerID)
	q.Add("LeagueID", c.LeagueID)
	req.URL.RawQuery = q.Encode()
	println(req.URL.String())

	b, err := c.Do(req)
	if err != nil {
		return err
	}

	var res Response
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	c.Response = &res
	return nil
}
