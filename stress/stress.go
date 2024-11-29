package stress

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"text/template"
	"time"
)

func StressTest(sUrl string, nMaxRequest int64, nConcurrency int64) {
	var wg *sync.WaitGroup

	respostasChan := make(chan Resposta, nMaxRequest)
	respostasConcur := make(chan struct{}, nConcurrency)

	InicioTeste := time.Now()

	wg = &sync.WaitGroup{}

	for nCount := int64(0); nCount < nMaxRequest; nCount++ {
		wg.Add(1)
		go testeUrl(sUrl, respostasChan, wg, respostasConcur)

	}
	wg.Wait()
	close(respostasChan)
	close(respostasConcur)
	relatorio(sUrl, respostasChan, nMaxRequest, nConcurrency, time.Since(InicioTeste), "")

}

func testeUrl(sUrl string, Respostas chan<- Resposta, wg *sync.WaitGroup, Semaforo chan struct{}) {
	var Resp Resposta
	Semaforo <- struct{}{}

	Incio := time.Now()
	resp, err := http.Get(sUrl)
	Final := time.Now()

	if err != nil {
		Resp = Resposta{Ini: Incio, Fim: Final, Code: -1, Erro: err}
	} else {
		Resp = Resposta{Ini: Incio, Fim: Final, Code: resp.StatusCode, Erro: err}
	}

	Respostas <- Resp
	wg.Done()
	<-Semaforo

}

func relatorio(sUrl string, Respostas <-chan Resposta, nMaxRequest int64, nConcurrency int64, TempoGasto time.Duration, sFormato string) {
	var RelDados Relatorio

	RelDados.TempoGasto = TempoGasto.String()
	RelDados.RequisicoesFeitas = nMaxRequest
	RelDados.Concorrencia = nConcurrency

	for respostas := range Respostas {
		switch respostas.Code / 100 {
		case 1:
			RelDados.InfTotais++
		case 2:
			RelDados.RequisicoesOk++
		case 3:
			RelDados.RedirecoesTotais++
		case 4:
			RelDados.ProblemaaCliTotais++
		case 5:
			RelDados.PorblemasSrvTotais++
		}
		if respostas.Erro != nil {
			RelDados.Erros++
		}
	}

	switch sFormato {
	case "csv":
	case "html":
	default:
		relDefault(sUrl, RelDados)
	}
}

func relDefault(sUrl string, RelDados Relatorio) {
	var Tpl string = `

---- Relatório de Teste ----
Dados Gerais:
    Url: ` + sUrl + `
    Concorrencia utilizada: {{.Concorrencia}}
    Tempo Gasto: {{.TempoGasto}}
    Requisições Feitas: {{.RequisicoesFeitas}}
    Requisições completadas (http 200): {{.RequisicoesOk}}


Dados por Códigos http:
   1xx Informativos       ---> {{.InfTotais}}
   3xx Redireccionamentos ---> {{.RedirecoesTotais}}
   4xx Erro no Cliente    ---> {{.ProblemaaCliTotais}}
   5xx Erro no Servidor   ---> {{.PorblemasSrvTotais}}

Requisições com Erro : {{.Erros }}
`
	t := template.Must(template.New("report").Parse(Tpl))

	err := t.Execute(os.Stdout, RelDados)

	if err != nil {
		fmt.Errorf("Erro ao criar report: %s", err)
	}
}
