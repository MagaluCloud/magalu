package cmd

import (
	"fmt"
	"os/exec"
	"path"
	"strconv"

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
		var result result

		for _, cmmd := range currentCommands {
			// fmt.Printf(`Executando comando: "%v - %s" -> $%s`, cmd.ID, cmd.Name, cmd.Command)

			output, err := exec.Command("sh", "-c", cmmd.Command+" --raw").CombinedOutput()

			if rewriteSnap { // TODO: normalizar o nome do arquivo, utilizando o proprio comando.
				_ = writeSnapshot(output, SNAP_DIR, cmmd.Module)
			}

			if err == nil {
				err = compareSnapshot(output, SNAP_DIR, cmmd.Module)
			}

			if err != nil {
				result.errors = append(result.errors, resultError{
					commandsList: cmmd,
					Error:        err.Error(),
				})
				continue
			}

			result.success = append(result.success, cmmd)
		}

		//TODO: Fazer um table bonitinho =)
		if len(result.errors) == 0 {
			fmt.Println("Sucesso! Todos os comandos executados sem alterações.")
			return
		}

		fmt.Println("Erros encontrados:")
		for _, er := range result.errors {

			fmt.Println("Command: ", er.Command)
			fmt.Println(er.Error)
		}
	},
}
