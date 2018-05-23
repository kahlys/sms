//Package ovh is a Go driver for the github.com/kahlys/sms package.
//
// In most cases clients will use the github.com/kahlys/sms package instead of
// using this package directly.
//
//		import (
//			"github.com/kahlys/sms"
//			_ "github.com/kahlys/sms/driver/ovh"
//		)
//		func main() {
//			param := map[string]string{
//				"account":"sms-xx4242-7",
//				"login":"bruce",
//				"password": "91939",
//				"sender": "wayne",
//			}
//			sender, _ := sms.Init("ovh", param)
//			sender.Send("Meet me at the roof !", "+33666666666")
//		}
//
// Make sure that the phone number you use is international formated.
package ovh

import (
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/kahlys/sms"
)

var (
	errAPISms    = errors.New("api-sms error")
	errBadParams = errors.New("ovh: missing parameter (account, login, password, sender)")
)

var apiurl = "https://www.ovh.com/cgi-bin/sms/http2sms.cgi?"

type ovhDriver struct {
	account  string
	login    string
	password string
	sender   string
}

// Init ovh driver
func (md *ovhDriver) Init(params map[string]string) error {
	if _, ok := params["account"]; !ok {
		return errBadParams
	}
	if _, ok := params["login"]; !ok {
		return errBadParams
	}
	if _, ok := params["password"]; !ok {
		return errBadParams
	}
	if _, ok := params["sender"]; !ok {
		return errBadParams
	}
	md.account = params["account"]
	md.login = params["login"]
	md.password = params["password"]
	md.sender = params["sender"]
	return nil
}

// Send a message
func (md *ovhDriver) Send(msg, num string) error {
	if num[0] == '+' {
		num = strings.Replace(num, "+", "00", 1)
	}
	u, _ := url.Parse(apiurl)
	q := u.Query()
	q.Set("account", md.account)
	q.Set("login", md.login)
	q.Set("password", md.password)
	q.Set("from", md.sender)
	q.Set("to", num)
	q.Set("message", msg)
	q.Set("noStop", "1")
	u.RawQuery = q.Encode()
	// timeout don't necessarily means the request errored, so it doesn't return an error.
	clt := &http.Client{
		Timeout: time.Second * 10,
	}
	response, err := clt.Get(u.String())
	if err, ok := err.(net.Error); ok && err.Timeout() {
		return nil
	}
	if err != nil {
		return errAPISms
	}
	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errAPISms
	}
	if strings.Contains(string(content), "KO") {
		return errAPISms
	}
	return nil
}

func init() {
	sms.Register("ovh", &ovhDriver{})
}
