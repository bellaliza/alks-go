package alks

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// SessionRequest is used to represent a new STS session request.
type SessionRequest struct {
	Account         string `json:"account"`
	Role            string `json:"role"`
	SessionDuration int    `json:"sessionTime"`
}

// SessionResponse is used to represent a new STS session.
type SessionResponse struct {
	AccessKey       string    `json:"accessKey"`
	SecretKey       string    `json:"secretKey"`
	SessionToken    string    `json:"sessionToken"`
	SessionDuration int       `json:"sessionDuration"`
	Expires         time.Time `json:"expires"`
}


// CreateSession will create a new STS session on AWS. If no error is
// returned then you will receive a SessionResponse object representing
// your STS session.
func (c *Client) CreateSession(account string, role string, sessionDuration int) (*SessionResponse, error) {
	log.Printf("[INFO] Creating %v hr session", sessionDuration)

	var found bool = false
	for _, v := range c.Durations() {
		log.Printf("compare %v to %v", sessionDuration, v)
		if sessionDuration == v {
			found = true
		}
	}

	log.Printf("did we find? %v", found)
	if !found {
		log.Printf("oops")
		return nil, fmt.Errorf("Unsupported session duration")
	}

	session := SessionRequest{
		account,
		role,
		sessionDuration,
	}

	b, err := json.Marshal(struct {
		SessionRequest
		AlksAccount
	}{session, c.Account})

	if err != nil {
		return nil, fmt.Errorf("Error encoding session create JSON: %s", err)
	}

	req, err := c.NewRequest(b, "POST", "/getKeys/")
	if err != nil {
		return nil, err
	}

	resp, err := checkResp(c.Http.Do(req))
	if err != nil {
		return nil, err
	}

	sr := new(SessionResponse)
	err = decodeBody(resp, &sr)

	if err != nil {
		return nil, fmt.Errorf("Error parsing session create response: %s", err)
	}

	sr.Expires = time.Now().Local().Add(time.Hour * time.Duration(sessionDuration))
	sr.SessionDuration = sessionDuration

	return sr, nil
}