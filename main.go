package main

import (
	"fmt"

	"github.com/ramonamorim/go-stress-test/stress"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "go-stress-test",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {

			sUrl, err := cmd.PersistentFlags().GetString("url")

			if err != nil {
				fmt.Printf("Erro ao obter valor da Url: %s", err.Error())
				return
			} else if sUrl == "" {
				fmt.Printf("Falta a URL a ser processada. Verifique!!!")
				return
			}

			nMaxRequest, err := cmd.PersistentFlags().GetInt64("requests")

			if err != nil {
				fmt.Printf("Erro ao obter valor das Request: %s", err.Error())
				return
			}

			nConcurrency, err := cmd.PersistentFlags().GetInt64("concurrency")

			if err != nil {
				fmt.Errorf("Erro ao obter valor da Concorrencia: %s", err)
				return
			}

			stress.StressTest(sUrl, nMaxRequest, nConcurrency)
		},
	}

	rootCmd.PersistentFlags().String("url", "", "Digite a Url que deve ser testada")
	rootCmd.PersistentFlags().Int64("requests", 100, "")
	rootCmd.PersistentFlags().Int64("concurrency", 2, "")

	rootCmd.Execute()
}
