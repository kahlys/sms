package ovh

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInit(t *testing.T) {
	param := map[string]string{
		"account":  "sms-xx4242-7",
		"login":    "bruce",
		"password": "wayne",
		"sender":   "batman",
	}
	sender := &ovhDriver{}

	err := sender.Init(param)
	if err != nil {
		t.Fatalf("wanted nil, got %s.", err)
	}
	if sender.account != param["account"] {
		t.Fatalf("wanted %s, got %s.", param["account"], sender.account)
	}
	if sender.login != param["login"] {
		t.Fatalf("wanted %s, got %s.", param["login"], sender.login)
	}
	if sender.password != param["password"] {
		t.Fatalf("wanted %s, got %s.", param["password"], sender.password)
	}
	if sender.sender != param["sender"] {
		t.Fatalf("wanted %s, got %s.", param["sender"], sender.sender)
	}
}

func TestInitErr(t *testing.T) {
	cases := []map[string]string{
		{
			"password": "wayne",
			"sender":   "batman",
		},
		{
			"login":  "bruce",
			"sender": "batman",
		},
		{
			"login":    "bruce",
			"password": "wayne",
		},
	}

	for _, c := range cases {
		sender := &ovhDriver{}
		if err := sender.Init(c); err != errBadParams {
			t.Fatalf("wanted %s, got %s.", errBadParams, err)
		}
	}
}

func TestSend(t *testing.T) {
	sender := &ovhDriver{}

	tests := []struct {
		responder *httptest.Server
		number    string
		message   string
		wantedErr error
	}{
		{
			responder: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "OK")
			})),
			number:    "+33605040302",
			message:   "test",
			wantedErr: nil,
		},
		{
			responder: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "KO")
			})),
			number:    "+33605040302",
			message:   "test",
			wantedErr: errAPISms,
		},
	}
	for _, tc := range tests {
		apiurl = tc.responder.URL
		if err := sender.Send(tc.message, tc.number); err != tc.wantedErr {
			t.Fatalf("wanted %s, got %s.", tc.wantedErr, err)
		}
	}
}
