package mobimel

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInit(t *testing.T) {
	sender := &mobimelDriver{}
	param := map[string]string{
		"key":     "bruce",
		"typesms": "wayne",
		"sender":  "batman",
	}
	// init ok
	err := sender.Init(param)
	if err != nil {
		t.Fatalf("wanted nil, got %s.", err)
	}
	if sender.key != param["key"] {
		t.Fatalf("wanted %s, got %s.", param["login"], sender.key)
	}
	if sender.typesms != param["typesms"] {
		t.Fatalf("wanted %s, got %s.", param["password"], sender.typesms)
	}
	if sender.sender != param["sender"] {
		t.Fatalf("wanted %s, got %s.", param["sender"], sender.sender)
	}
	// init error missing parameter
	delete(param, "sender")
	err = sender.Init(param)
	if err == nil {
		t.Fatalf("expected an error")
	}
	delete(param, "typesms")
	err = sender.Init(param)
	if err == nil {
		t.Fatalf("expected an error")
	}
	delete(param, "key")
	err = sender.Init(param)
	if err == nil {
		t.Fatalf("expected an error")
	}

}

func TestSend(t *testing.T) {
	sender := &mobimelDriver{}

	cases := []struct {
		responder   *httptest.Server
		number      string
		message     string
		expectedErr bool
	}{
		{
			responder: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{"resultat": 1})
			})),
			number:      "+33605040302",
			message:     "test",
			expectedErr: false,
		},
		{
			responder: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{"resultat": 0, "erreurs": "2"})
			})),
			number:      "+33605040302",
			message:     "test",
			expectedErr: true,
		},
	}

	for _, c := range cases {
		apiurl = c.responder.URL
		err := sender.Send(c.number, c.number)
		if err == nil && c.expectedErr {
			t.Errorf("expected an error")
		}
	}
}
