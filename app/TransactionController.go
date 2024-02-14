package app

import (
	"context"
	"errors"
	"github.com/edimarlnx/rinha-de-backend-2024-q1-edimarlnx/types"
	"github.com/edimarlnx/rinha-de-backend-2024-q1-edimarlnx/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"

	"github.com/jackc/pgx/v5"
)

type TransactionController struct {
	pool *pgxpool.Pool
}

func New(dbUri string) *TransactionController {
	transactionController := TransactionController{}
	transactionController.initPgx(dbUri)
	return &transactionController
}

func (t *TransactionController) initPgx(dbUri string) {
	pool, err := pgxpool.New(context.Background(), dbUri)
	if err != nil {
		utils.Log.WithField("err", err).Panicln("Erro ao conectar ao banco de dados.")
	}
	t.pool = pool
}

func (t *TransactionController) Transacao(clientId int, transacao types.Transacao) (*types.Saldo, *types.ErrorWithCode) {
	errWithCode := t.checkCliente(clientId)
	if errWithCode != nil {
		return nil, errWithCode
	}
	err := transacao.Valid()
	if err != nil {
		return nil, t.ErrorWithCode(err.Error(), 422)
	}
	now := time.Now()
	tx, err := t.pool.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		tx.Rollback(context.Background())
		return nil, t.ErrorWithCode(err.Error(), 400)
	}
	var saldoAtualizado int
	valorTransacao := transacao.Valor
	if transacao.Tipo == "d" {
		valorTransacao *= -1
	}
	err = tx.QueryRow(context.Background(), "update clientes c set saldo = (c.saldo + $2), data_extrato = $3 WHERE c.id = $1 RETURNING c.limite + c.saldo", clientId, valorTransacao, now).Scan(&saldoAtualizado)
	if err != nil {
		tx.Rollback(context.Background())
		return nil, t.ErrorWithCode(err.Error(), 400)
	}
	if saldoAtualizado < 0 {
		tx.Rollback(context.Background())
		return nil, t.ErrorWithCode("Valor ultrapassa o limite do cliente", 422)
	}
	errWithCode = t.saveTransacao(tx, clientId, now, transacao)
	if errWithCode != nil {
		tx.Rollback(context.Background())
		return nil, errWithCode
	}
	saldo, errWithCode := t.loadSaldo(tx, clientId)
	if errWithCode != nil {
		tx.Rollback(context.Background())
		return nil, errWithCode
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, t.ErrorWithCode(err.Error(), 400)
	}
	return saldo, nil
}

func (t *TransactionController) saveTransacao(tx pgx.Tx, clienteId int, realizadaEm time.Time, transacao types.Transacao) *types.ErrorWithCode {
	rows, err := tx.Query(context.Background(), `INSERT INTO transacoes (cliente_id,valor,tipo,descricao,realizada_em) 
		VALUES ($1, $2, $3, $4, $5)`, clienteId, transacao.Valor, transacao.Tipo, transacao.Descricao, realizadaEm)
	defer rows.Close()
	if err != nil {
		return t.ErrorWithCode("Erro ao salvar a transação.", 400)
	}
	return nil
}

func (t *TransactionController) loadSaldo(tx pgx.Tx, clienteId int) (*types.Saldo, *types.ErrorWithCode) {
	var saldo types.Saldo
	rows, err := tx.Query(context.Background(), "SELECT saldo,limite,data_extrato FROM clientes WHERE id = $1", clienteId)
	if err != nil {
		return nil, t.ErrorWithCode(err.Error(), 400)
	}
	for rows.Next() {
		err = rows.Scan(&saldo.Saldo, &saldo.Limite, &saldo.DataExtrato)
		if err != nil {
			return nil, t.ErrorWithCode(err.Error(), 400)
		}
	}
	rows.Close()
	if err != nil {
		return nil, t.ErrorWithCode(err.Error(), 400)
	}
	return &saldo, nil
}

func (t *TransactionController) Extrato(clientId int) (*types.Extrato, *types.ErrorWithCode) {
	errWithCode := t.checkCliente(clientId)
	if errWithCode != nil {
		return nil, errWithCode
	}
	saldoExtrato := types.SaldoExtrato{
		DataExtrato: time.Now(),
	}

	rows, err := t.pool.Query(context.Background(), "SELECT saldo, limite  FROM clientes WHERE id = $1", clientId)
	if err != nil {
		return nil, t.ErrorWithCode(err.Error(), 400)
	}
	for rows.Next() {
		err = rows.Scan(&saldoExtrato.Total, &saldoExtrato.Limite)
		if err != nil {
			return nil, t.ErrorWithCode(err.Error(), 400)
		}
	}
	rows.Close()
	rows, err = t.pool.Query(context.Background(), "SELECT valor,tipo,descricao,realizada_em FROM transacoes WHERE cliente_id = $1 order by realizada_em desc limit 10", clientId)
	if err != nil {
		return nil, t.ErrorWithCode(err.Error(), 400)
	}
	transacoes := []types.Transacao{}
	for rows.Next() {
		var tr types.Transacao
		err = rows.Scan(&tr.Valor, &tr.Tipo, &tr.Descricao, &tr.RealizadaEm)
		transacoes = append(transacoes, tr)
	}
	rows.Close()
	if err != nil {
		return nil, t.ErrorWithCode(err.Error(), 400)
	}
	extrato := types.Extrato{
		Saldo:             saldoExtrato,
		UltimasTransacoes: transacoes,
	}
	return &extrato, nil
}

func (t *TransactionController) checkCliente(clienteId int) *types.ErrorWithCode {
	rows, err := t.pool.Query(context.Background(), "SELECT * FROM clientes WHERE id = $1", clienteId)
	defer rows.Close()
	if err != nil || !rows.Next() {
		if err != nil {
			utils.Log.Error(err.Error())
		}
		return t.ErrorWithCode("Cliente não encontrado.", 404)
	}
	return nil
}

func (t *TransactionController) ErrorWithCode(msg string, code int) *types.ErrorWithCode {
	return &types.ErrorWithCode{
		Error: errors.New(msg),
		Code:  code,
	}
}

func (t *TransactionController) ResetClient(clienteId int) {
	t.pool.Exec(context.Background(), "delete from transacoes WHERE cliente_id = $1", clienteId)
	t.pool.Exec(context.Background(), "update clientes set saldo = 0 WHERE id = $1", clienteId)
}
