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
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/kahlys/sms"
)

var apiurl = "https://extranet.mobimel.com/envoi_form.asp"

var (
	errAPISms     = errors.New("api-sms error")
	errAPINumber  = errors.New("api-sms: bad phone number")
	errAPIAccount = errors.New("api-sms: bad account configuration")
	errAPIProcess = errors.New("api-sms: something went wrong during process")
	errAPIMsg     = errors.New("api-sms: bad message format")
	errAPIBalance = errors.New("api-sms: no more credit")
)

type mobimelDriver struct {
	login    string
	password string
	sender   string
}

var errBadParams = errors.New("mobimel: missing parameter (login, password, sender)")

// Init sms driver
func (md *mobimelDriver) Init(params map[string]string) error {
	if _, ok := params["login"]; !ok {
		return errBadParams
	}
	if _, ok := params["password"]; !ok {
		return errBadParams
	}
	if _, ok := params["sender"]; !ok {
		return errBadParams
	}
	md.login = params["login"]
	md.password = params["password"]
	md.sender = params["sender"]
	return nil
}

// Send a message
func (md *mobimelDriver) Send(msg, num string) error {
	form := url.Values{}
	form.Add("login", md.login)
	form.Add("passwd", md.password)
	form.Add("tel", num)
	form.Add("msg", msg)
	form.Add("emetteur", md.sender)
	clt := &http.Client{
		Timeout: time.Second * 10,
	}
	response, err := clt.PostForm(apiurl, form)
	if err, ok := err.(net.Error); ok && err.Timeout() {
		return nil
	}
	if err != nil {
		return errAPISms
	}
	defer response.Body.Close()
	content, _ := ioutil.ReadAll(response.Body)
	t := strings.Split(string(content), "/")
	if len(t) < 3 {
		return errAPISms
	}
	switch t[0] {
	case "0":
		break
	case "101", "102", "103", "104", "108", "118":
		return errAPIAccount
	case "105", "106", "111", "112", "4":
		return errAPIBalance
	case "100", "107", "1", "5", "6":
		return errAPIProcess
	case "117", "2", "3":
		return errAPINumber
	case "109", "110":
		return errAPIMsg
	default:
		return errAPISms
	}
	return nil
}

func init() {
	sms.Register("mobimel", &mobimelDriver{})
}
