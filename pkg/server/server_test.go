package server

import (
	"log"
	"myapp/pkg/config"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	cfg, err := config.New("../../config.yaml")
	if err != nil {
		log.Println(err)
	}
	server := NewServer(cfg)

	assert.NotNil(t, server)
	assert.Equal(t, ":8080", server.Serv.Addr)
	assert.IsType(t, &http.ServeMux{}, server.Router)
}
