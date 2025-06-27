package deadshotgolib

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

type LogModel struct {
	ID              int       `json:"id" db:"id"`
	Method          string    `json:"method" db:"method"`
	URL             string    `json:"url" db:"url"`
	Headers         string    `json:"headers" db:"headers"`
	QueryParams     string    `json:"query_params" db:"query_params"`
	Body            string    `json:"body" db:"body"`
	ReceivedAt      time.Time `json:"received_at" db:"received_at"`
	StatusCode      int       `json:"status_code" db:"status_code"`
	ResponseHeaders string    `json:"response_headers" db:"response_headers"`
	ResponseBody    string    `json:"response_body" db:"response_body"`
	Tags            string    `json:"tags" db:"tags"`
	Source          string    `json:"source" db:"source"`
	Replayed        bool      `json:"replayed" db:"replayed"`
	Error           string    `json:"error" db:"error"`
}

type DeadShot struct {
	EndPoint string
}

func (d DeadShot) Send(log LogModel) error {
	hClient := http.Client{}

	serializedBody, serErr := json.Marshal(log)
	if serErr != nil {
		return serErr
	}

	body := strings.NewReader(string(serializedBody))
	req, reqErr := http.NewRequest(http.MethodPost, d.EndPoint, body)
	if reqErr != nil {
		return reqErr
	}

	res, resErr := hClient.Do(req)
	if resErr != nil {
		return resErr
	}

	if res.StatusCode == http.StatusOK {
		return nil
	} else {
		defer res.Body.Close()
		bodyBytes, ioErr := io.ReadAll(res.Body)
		if ioErr != nil {
			return ioErr
		}
		return errors.New(string(bodyBytes))
	}
}
