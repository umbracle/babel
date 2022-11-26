package babel

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type MockHttpServer struct {
	srv  *http.Server
	addr string
	mux  *http.ServeMux
}

func NewMockHttpServer(port uint64) *MockHttpServer {
	mux := http.NewServeMux()

	addr := fmt.Sprintf("127.0.0.1:%d", port)
	s := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		if err := s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()
	time.Sleep(500 * time.Millisecond)

	srv := &MockHttpServer{
		srv:  s,
		addr: addr,
		mux:  mux,
	}
	return srv
}

func (m *MockHttpServer) Mux() *http.ServeMux {
	return m.mux
}

func (m *MockHttpServer) Http() string {
	return "http://" + m.addr
}

func (m *MockHttpServer) Stop() {
	if err := m.srv.Shutdown(context.Background()); err != nil {
		log.Fatal(err)
	}
}
