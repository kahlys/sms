package mobimel

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInit(t *testing.T) {
	param := map[string]string{
		"login":    "bruce",
		"password": "wayne",
		"sender":   "batman",
	}
	sender := &mobimelDriver{}

	err := sender.Init(param)
	if err != nil {
		t.Fatalf("wanted nil, got %s.", err)
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
		sender := &mobimelDriver{}
		if err := sender.Init(c); err != errBadParams {
			t.Fatalf("wanted %s, got %s.", errBadParams, err)
		}
	}
}

func TestSend(t *testing.T) {
	sender := &mobimelDriver{}

	tests := []struct {
		responder *httptest.Server
		number    string
		message   string
		wantedErr error
	}{
		{
			responder: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("0/+33605040302/Le message a été envoyé."))
			})),
			number:    "+33605040302",
			message:   "test",
			wantedErr: nil,
		},
		{
			responder: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("47/"))
			})),
			number:    "+33605040302",
			message:   "test",
			wantedErr: errAPISms,
		},
		{
			responder: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("47/+33605040302/This is a new millenium."))
			})),
			number:    "+33605040302",
			message:   "test",
			wantedErr: errAPISms,
		},
		{
			responder: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("2/0102030405/Le numéro n'est pas valide."))
			})),
			number:    "+33605040302",
			message:   "test",
			wantedErr: errAPINumber,
		},
		{
			responder: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("101/+33605040302/Vous n'êtes pas autorisé à envoyer des messages."))
			})),
			number:    "+33605040302",
			message:   "test",
			wantedErr: errAPIAccount,
		},
		{
			responder: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("107/+33605040302/Service momentanément indisponible (maintenance)."))
			})),
			number:    "+33605040302",
			message:   "test",
			wantedErr: errAPIProcess,
		},
		{
			responder: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("109/+33605040302/Message absent ou vide."))
			})),
			number:    "+33605040302",
			message:   "test",
			wantedErr: errAPIMsg,
		},
		{
			responder: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("105/+33605040302/Erreur : votre solde est épuisé."))
			})),
			number:    "+33605040302",
			message:   "test",
			wantedErr: errAPIBalance,
		},
	}

	for _, tc := range tests {
		apiurl = tc.responder.URL
		if err := sender.Send(tc.number, tc.number); err != tc.wantedErr {
			t.Fatalf("wanted %s, got %s.", tc.wantedErr, err)
		}
	}
}
