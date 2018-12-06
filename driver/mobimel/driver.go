//Package mobimel is a Go driver for the github/kahlys/sms package.
//
// In most cases clients will use the github/kahlys/sms package instead of
// using this package directly.
//
//    import (
//    	"github/kahlys/sms"
//    	_ "github/kahlys/sms/driver/mobimel"
//    )
//    func main() {
//    	param := map[string]string{
//    		"login":"bruce",
//    		"password": "91939",
//    		"sender": "wayne",
//    	}
//    	sender, _ := sms.Init("mobimel", param)
//    	sender.Send("Meet me at the roof !", "+33666666666")
//    }
//
// Make sure that the phone number you use is international formated.
package mobimel

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/kahlys/sms"
)

var apiurl = "https://portail.mobimel.com/api/envoyer/sms"

type mobimelDriver struct {
	key     string
	sender  string
	typesms string
}

// Init mobimel driver
func (md *mobimelDriver) Init(params map[string]string) error {
	if _, ok := params["key"]; !ok {
		return fmt.Errorf("mobimel: missing parameter key")
	}
	if _, ok := params["typesms"]; !ok {
		return fmt.Errorf("mobimel: missing parameter typesms")
	}
	if _, ok := params["sender"]; !ok {
		return fmt.Errorf("mobimel: missing parameter sender")
	}
	md.key = params["key"]
	md.typesms = params["typesms"]
	md.sender = params["sender"]
	return nil
}

// Send a message
func (md *mobimelDriver) Send(msg, num string) error {
	req, err := http.NewRequest("GET", apiurl, nil)
	if err != nil {
		return fmt.Errorf("mobimel: unable to create api request")
	}
	q := req.URL.Query()
	q.Add("key", md.key)
	q.Add("type", md.typesms)
	q.Add("message", msg)
	q.Add("destinataires", num)
	q.Add("expediteur", md.sender)
	req.URL.RawQuery = q.Encode()

	client := http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err, ok := err.(net.Error); ok && err.Timeout() {
		// timeout are infrequent and don't necessarily means the request errored, so the function doesn't return an error.
		return nil
	}
	if err != nil {
		return fmt.Errorf("mobimel: unable to send request: %v", err)
	}
	var body struct {
		Result int    `json:"resultat"`
		Errors string `json:"erreurs"`
	}
	json.NewDecoder(resp.Body).Decode(&body)
	if body.Result == 0 {
		return fmt.Errorf("mobimel api error: %v (see mobimel api documentation)", body.Errors)
	}
	return nil
}

func init() {
	sms.Register("mobimel", &mobimelDriver{})
}
