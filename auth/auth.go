package auth

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Login returns an access_token using the Headr Test Client's client_id and client secret
func Login() string {
	url := "https://headr.auth0.com/oauth/token"
	payload := strings.NewReader("{\"client_id\":\"diOgwfmhPiX3lcsb7JkRUr10HZNBwgWr\",\"client_secret\":\"6_To1KTBxGJ6ep7OZb0XGoxtvkL1CSIohIlMmQGHUq3Z1y6CqssKYKfktDhtpEN4\",\"audience\":\"https://api.headr.io\",\"grant_type\":\"client_credentials\"}")
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("content-type", "application/json")
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	var response struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   string `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}
	json.NewDecoder(res.Body).Decode(&response)
	return response.AccessToken
}
