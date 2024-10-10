package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"sync"

	"github.com/spf13/cobra"
)

type resultError struct {
	commandsList
	Error string
}

type result struct {
	errors  []resultError
	success []commandsList
}

var runTestsCmd = &cobra.Command{
	Use:    "test",
	Short:  "Run all available tests",
	Hidden: false,
	Run: func(cmd *cobra.Command, args []string) {

		if os.Getenv("MGC_API_KEY") == "" {
			// TODO: Validar object storage também
			fmt.Println("Se faz necessário a variavel de ambiente MGC_API_KEY")
			return
		}

		rewriteSnap := false
		if len(args) > 0 {
			rewriteSnap, _ = strconv.ParseBool(args[0])
		}

		_ = ensureDirectoryExists(path.Join(currentDir(), SNAP_DIR))

		currentCommands, err := loadList()

		if err != nil {
			fmt.Println(err)
			return

		}
		result := result{}
		var wg sync.WaitGroup

		for _, cmmd := range currentCommands {
			wg.Add(1)
			go func() {
				defer wg.Done()
				output, err := exec.Command("sh", "-c", cmmd.Command+" --raw").CombinedOutput()
				if err != nil {
					result.errors = append(result.errors, resultError{
						commandsList: cmmd,
						Error:        err.Error(),
					})
					return
				}

				snapshotFile := normalizeCommandToFile(cmmd.Command)
				if rewriteSnap { // TODO: normalizar o nome do arquivo, utilizando o proprio comando.
					_ = writeSnapshot(output, SNAP_DIR, snapshotFile)
				}

				err = compareSnapshot(output, SNAP_DIR, snapshotFile)

				if err != nil {
					result.errors = append(result.errors, resultError{
						commandsList: cmmd,
						Error:        err.Error(),
					})
					return
				}

				result.success = append(result.success, cmmd)
			}()
		}

		wg.Wait()

		//TODO: Fazer um table bonitinho =)
		if len(result.errors) == 0 {
			fmt.Println("Sucesso! Todos os comandos executados sem alterações.")
			return
		}

		fmt.Print("\nErros encontrados:\n\n")
		for _, er := range result.errors {

			fmt.Println("Command: ", er.Command)
			fmt.Println("Error: ", er.Error)
		}
	},
}
