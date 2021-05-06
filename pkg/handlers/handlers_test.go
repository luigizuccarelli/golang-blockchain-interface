package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/microlib/simple"
	"lmzsoftware.com/lzuccarelli/golang-blockchain-interface/pkg/connectors"
)

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("Injected error")
}

// Fake all connections
type FakeConnections struct {
	log  *simple.Logger
	Flag string
}

func (r *FakeConnections) Error(msg string, val ...interface{}) {
	r.log.Error(fmt.Sprintf(msg, val...))
}

func (r *FakeConnections) Info(msg string, val ...interface{}) {
	r.log.Info(fmt.Sprintf(msg, val...))
}

func (r *FakeConnections) Debug(msg string, val ...interface{}) {
	r.log.Debug(fmt.Sprintf(msg, val...))
}

func (r *FakeConnections) Trace(msg string, val ...interface{}) {
	r.log.Trace(fmt.Sprintf(msg, val...))
}

// NewTestConnections - create all mock connections
func NewTestConnections(file string, code int, logger *simple.Logger) connectors.Clients {
	conns := &FakeConnections{log: logger}
	return conns
}

func TestHandlers(t *testing.T) {

	logger := &simple.Logger{Level: "info"}

	t.Run("IsAlive : should pass", func(t *testing.T) {
		var STATUS int = 200
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v2/sys/info/isalive", nil)
		NewTestConnections("../../tests/payload.json", STATUS, logger)
		handler := http.HandlerFunc(IsAlive)
		handler.ServeHTTP(rr, req)

		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "IsAlive", rr.Code, STATUS))
		}
	})

	t.Run("Init : GET should pass", func(t *testing.T) {
		var STATUS int = 200
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/genesis", nil)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Init(w, r)
		})
		handler.ServeHTTP(rr, req)

		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "Init", rr.Code, STATUS))
		}
	})

	t.Run("GetBlockchain : GET should pass", func(t *testing.T) {
		var STATUS int = 200
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/blockchain/1", nil)
		conn := NewTestConnections("../../tests/payload.json", STATUS, logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			GetBlockChain(w, r, conn)
		})

		// inject var for gorilla mux
		vars := map[string]string{
			"index": "test",
		}
		req = mux.SetURLVars(req, vars)

		handler.ServeHTTP(rr, req)

		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "GetBlockChain", rr.Code, STATUS))
		}
	})

	t.Run("GetBlockchainList : GET should pass", func(t *testing.T) {
		var STATUS int = 200
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/blockchain/list", nil)
		conn := NewTestConnections("../../tests/payload.json", STATUS, logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			GetBlockChainList(w, r, conn)
		})
		handler.ServeHTTP(rr, req)

		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "GetBlockChainList", rr.Code, STATUS))
		}
	})

	t.Run("WriteBlockChain : POST should pass", func(t *testing.T) {
		var STATUS int = 200
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		data, _ := ioutil.ReadFile("../../tests/input-to-encrypt.json")
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/blockchain", bytes.NewBuffer(data))
		conn := NewTestConnections("../../tests/payload.json", STATUS, logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			WriteBlockChain(w, r, conn)
		})
		handler.ServeHTTP(rr, req)

		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "GetBlockChainList", rr.Code, STATUS))
		}
	})

	t.Run("WriteBlockChain : POST should fail (force readall error)", func(t *testing.T) {
		var STATUS int = 500
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/blockchain", errReader(0))
		conn := NewTestConnections("../../tests/payload.json", STATUS, logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			WriteBlockChain(w, r, conn)
		})
		handler.ServeHTTP(rr, req)

		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "GetBlockChainList", rr.Code, STATUS))
		}
	})

	t.Run("WriteBlockChain : POST should fail (json error)", func(t *testing.T) {
		var STATUS int = 500
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/blockchain", bytes.NewBuffer([]byte("{ \"test\":")))
		conn := NewTestConnections("../../tests/payload.json", STATUS, logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			WriteBlockChain(w, r, conn)
		})
		handler.ServeHTTP(rr, req)

		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "GetBlockChainList", rr.Code, STATUS))
		}
	})

	t.Run("WriteBlockChain : POST should fail (empty object blocks)", func(t *testing.T) {
		var STATUS int = 500
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/blockchain", bytes.NewBuffer([]byte("{ \"objectA\":\"\" }")))
		conn := NewTestConnections("../../tests/payload.json", STATUS, logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			WriteBlockChain(w, r, conn)
		})
		handler.ServeHTTP(rr, req)

		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "GetBlockChainList", rr.Code, STATUS))
		}
	})

}
