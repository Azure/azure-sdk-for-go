package easykey

import (
	"context"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
)

var pfx []byte

func init() {
	var err error
	pfx, err = ioutil.ReadFile("./testdata/server.pxf")
	if err != nil {
		panic(err)
	}
}

func TestPrivateKey(t *testing.T) {
	ctx := context.Background()
	mock, err := NewMock("vault")
	if err != nil {
		panic(err)
	}

	if err := mock.AddPKCS12("private", LatestVersion, pfx, "", 1); err != nil {
		panic(err)
	}

	kv, err := New(ctx, "vault", Config{}, MockService(mock))
	if err != nil {
		panic(err)
	}

	tests := []struct {
		desc string
		name string
		err  bool
		want []byte
	}{
		{
			desc: "secret doesn't exist",
			err:  true,
		},
		{
			desc: "secret is not base64 encoded",
			err:  true,
		},
		{
			desc: "success",
			name: "private",
		},
	}

	for _, test := range tests {
		_, err := kv.PrivateKey(ctx, test.name, LatestVersion)

		switch {
		case test.err && err == nil:
			t.Errorf("TestPrivateKey(%s): got err == nil, want err != nil", test.desc)
			continue
		case !test.err && err != nil:
			t.Errorf("TestPrivateKey(%s): got err == %s, want err == nil", test.desc, err)
			continue
		case err != nil:
			continue
		}
	}
}

func TestTLSCert(t *testing.T) {
	ctx := context.Background()
	mock, err := NewMock("vault")
	if err != nil {
		panic(err)
	}
	pfx, err := ioutil.ReadFile("./testdata/server.pxf")
	if err != nil {
		panic(err)
	}
	if err := mock.AddPKCS12("private", LatestVersion, pfx, "", 0); err != nil {
		panic(err)
	}

	kv, err := New(ctx, "vault", Config{}, MockService(mock))
	if err != nil {
		panic(err)
	}

	cert, err := kv.Certificate(ctx, "private", LatestVersion)
	if err != nil {
		panic(err)
	}
	_, err = cert.X509()
	if err != nil {
		panic(err)
	}

	tests := []struct {
		desc string
		name string
		err  bool
	}{
		{
			desc: "secret doesn't exist",
			name: "nope",
			err:  true,
		},
		{
			desc: "success",
			name: "private",
		},
	}

	for _, test := range tests {
		tlsCert, err := kv.TLSCert(ctx, test.name, LatestVersion, 0, 1)
		switch {
		case test.err && err == nil:
			t.Errorf("TestPrivateKey(%s): got err == nil, want err != nil", test.desc)
			continue
		case !test.err && err != nil:
			t.Errorf("TestPrivateKey(%s): got err == %s, want err == nil", test.desc, err)
			continue
		case err != nil:
			continue
		}

		cfg := &tls.Config{Certificates: []tls.Certificate{tlsCert}}
		srv := &http.Server{
			Addr:         ":8080",
			TLSConfig:    cfg,
			ReadTimeout:  time.Minute,
			WriteTimeout: time.Minute,
		}
		go func() {
			if err := srv.ListenAndServeTLS("", ""); err != nil {
				log.Println("server stopped listening, this may just be the server shutting down: ", err)
			}
		}()
		defer srv.Close()

		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

		stop := time.Now().Add(10 * time.Second)
		success := false
		for time.Now().Before(stop) {
			_, err := http.Get("https://127.0.0.1:8080")
			if err != nil {
				time.Sleep(1 * time.Second)
				continue
			}

			success = true
			break
		}
		if !success {
			t.Errorf("TestPrivateKey(%s): server did not come up", test.desc)
		}
	}
}
