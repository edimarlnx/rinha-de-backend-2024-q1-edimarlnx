package webhook

import (
	"bytes"
	"encoding/json"
	"github.com/edimarlnx/rinha-de-backend-2024-q1-edimarlnx/app"
	"github.com/edimarlnx/rinha-de-backend-2024-q1-edimarlnx/types"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingRoute(t *testing.T) {
	dbUri := "postgres://postgres:postgres@localhost:55432/rinha-backend?sslmode=disable"
	transactionController := app.New(dbUri)
	server, _ := CreateWebhook(transactionController)
	transactionController.ResetClient(1)
	type args struct {
		builder string
		m       map[string]string
	}
	tests := []struct {
		name   string
		method string
		path   string
		body   interface{}
		assert func(w *httptest.ResponseRecorder, t *testing.T)
		args   args
	}{
		{
			name:   "extrato-inicial",
			method: "GET",
			path:   "/clientes/1/extrato",
			assert: func(w *httptest.ResponseRecorder, t *testing.T) {
				var extrato types.Extrato
				json.Unmarshal(w.Body.Bytes(), &extrato)
				assert.Equal(t, 200, w.Code)
				assert.Equal(t, 0, extrato.Saldo.Total)
				assert.Equal(t, 100000, extrato.Saldo.Limite)
				assert.Equal(t, 0, len(extrato.UltimasTransacoes))
			},
		},
		{
			name:   "transacao-cliente-not-found",
			method: "POST",
			path:   "/clientes/6/transacoes",
			body:   types.Transacao{Tipo: "d", Valor: 50000},
			assert: func(w *httptest.ResponseRecorder, t *testing.T) {
				assert.Equal(t, 404, w.Code)
			},
		},
		{
			name:   "transacao-1-debito",
			method: "POST",
			path:   "/clientes/1/transacoes",
			body:   types.Transacao{Tipo: "d", Valor: 50000, Descricao: "debito"},
			assert: func(w *httptest.ResponseRecorder, t *testing.T) {
				var saldo types.Saldo
				json.Unmarshal(w.Body.Bytes(), &saldo)
				assert.Equal(t, 200, w.Code)
				assert.Equal(t, -50000, saldo.Saldo)
				assert.Equal(t, 100000, saldo.Limite)
			},
		},
		{
			name:   "transacao-2-debito",
			method: "POST",
			path:   "/clientes/1/transacoes",
			body:   types.Transacao{Tipo: "d", Valor: 50000, Descricao: "debito"},
			assert: func(w *httptest.ResponseRecorder, t *testing.T) {
				var saldo types.Saldo
				json.Unmarshal(w.Body.Bytes(), &saldo)
				assert.Equal(t, 200, w.Code)
				assert.Equal(t, -100000, saldo.Saldo)
				assert.Equal(t, 100000, saldo.Limite)
			},
		},
		{
			name:   "transacao-3-falha-sem-limite",
			method: "POST",
			path:   "/clientes/1/transacoes",
			body:   types.Transacao{Tipo: "d", Valor: 50000, Descricao: "debito"},
			assert: func(w *httptest.ResponseRecorder, t *testing.T) {
				assert.Equal(t, 422, w.Code)
			},
		},
		{
			name:   "transacao-4-credito",
			method: "POST",
			path:   "/clientes/1/transacoes",
			body:   types.Transacao{Tipo: "c", Valor: 2000, Descricao: "credito"},
			assert: func(w *httptest.ResponseRecorder, t *testing.T) {
				var saldo types.Saldo
				json.Unmarshal(w.Body.Bytes(), &saldo)
				assert.Equal(t, 200, w.Code)
				assert.Equal(t, -98000, saldo.Saldo)
				assert.Equal(t, 100000, saldo.Limite)
			},
		},
		{
			name:   "transacao-5-falha-float",
			method: "POST",
			path:   "/clientes/1/transacoes",
			body:   types.Transacao{Tipo: "d", Valor: 2.2, Descricao: "debito"},
			assert: func(w *httptest.ResponseRecorder, t *testing.T) {
				assert.Equal(t, 422, w.Code)
			},
		},
		{
			name:   "extrato-final",
			method: "GET",
			path:   "/clientes/1/extrato",
			assert: func(w *httptest.ResponseRecorder, t *testing.T) {
				var extrato types.Extrato
				json.Unmarshal(w.Body.Bytes(), &extrato)
				assert.Equal(t, 200, w.Code)
				assert.Equal(t, -98000, extrato.Saldo.Total)
				assert.Equal(t, 100000, extrato.Saldo.Limite)
				assert.Equal(t, 3, len(extrato.UltimasTransacoes))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			var body io.Reader
			if tt.body != nil {
				data, _ := json.Marshal(tt.body)
				body = bytes.NewReader(data)
			}
			req, _ := http.NewRequest(tt.method, tt.path, body)
			server.ServeHTTP(w, req)
			tt.assert(w, t)
		})
	}
}
