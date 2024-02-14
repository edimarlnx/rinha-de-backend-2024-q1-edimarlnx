package types

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Validation interface {
	Valid() error
}

type Response struct {
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
	Status  int    `json:"status"`
}

type ErrorWithCode struct {
	Error error
	Code  int
}

type Transacao struct {
	Valor       float32   `json:"valor"`
	Tipo        string    `json:"tipo"`
	Descricao   string    `json:"descricao"`
	RealizadaEm time.Time `json:"realizada_em,omitempty"`
}

type Saldo struct {
	Limite      int       `json:"limite"`
	Saldo       int       `json:"saldo"`
	DataExtrato time.Time `json:"data_extrato,omitempty"`
}

type Cliente struct {
	Id string `json:"id"`
	Saldo
}

type SaldoExtrato struct {
	Limite      int       `json:"limite"`
	Total       int       `json:"total"`
	DataExtrato time.Time `json:"data_extrato,omitempty"`
}

type Extrato struct {
	Saldo             SaldoExtrato `json:"saldo"`
	UltimasTransacoes []Transacao  `json:"ultimas_transacoes"`
}

func (t Transacao) Valid() error {
	tipo := strings.ToLower(t.Tipo)
	if tipo != "c" && tipo != "d" {
		return errors.New(fmt.Sprintf("Tipo da transação aceitas 'd' e 'd'. '[%s]' não é um valor válido.", tipo))
	}
	if t.Valor < 0 {
		return errors.New(fmt.Sprintf("Valor da transação deve ser maior um número positivo [%f].", t.Valor))
	}
	intValue := float32(int32(t.Valor))
	if (t.Valor - intValue) > 0 {
		return errors.New(fmt.Sprintf("Valor da transação deve ser um número inteiro [%f].", t.Valor))
	}
	descSize := len(strings.TrimSpace(t.Descricao))
	if descSize == 0 || descSize > 10 {
		return errors.New(fmt.Sprintf("Descrição é obrigatória e deve ter no máximo 10 caracteres. '%s' contém '%d'.", t.Descricao, len(t.Descricao)))
	}
	return nil
}
