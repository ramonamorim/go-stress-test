package stress

import "time"

type Resposta struct {
	Ini  time.Time
	Fim  time.Time
	Code int
	Erro error
}

type Relatorio struct {
	RequisicoesFeitas  int64
	RequisicoesOk      int64
	Concorrencia       int64
	TempoGasto         string
	InfTotais          int64
	RedirecoesTotais   int64
	ProblemaaCliTotais int64
	PorblemasSrvTotais int64
	Erros              int64
}
