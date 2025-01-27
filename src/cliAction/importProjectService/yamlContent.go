package importProjectService

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zeropsio/zcli/src/constants"
	"github.com/zeropsio/zcli/src/i18n"
)

func getImportYamlContent(config RunConfig) ([]byte, error) {
	fmt.Println(i18n.YamlCheck)

	importYamlPath := config.ImportYamlPath
	if !filepath.IsAbs(importYamlPath) {

		workingDir, err := filepath.Abs(config.WorkingDir)
		if err != nil {
			return nil, err
		}

		importYamlPath = filepath.Join(workingDir, importYamlPath)
	}

	fileInfo, err := os.Stat(importYamlPath)
	if err != nil {
		return nil, err
	}

	if fileInfo.IsDir() {
		return nil, errors.New(i18n.ImportYamlNotFound)
	}

	fmt.Printf("%s: %s\n", i18n.ImportYamlFound, importYamlPath)

	if fileInfo.Size() == 0 {
		return nil, errors.New(i18n.ImportYamlEmpty)
	}

	if fileInfo.Size() > 100*1024 {
		return nil, errors.New(i18n.ImportYamlTooLarge)
	}

	yamlContent, err := os.ReadFile(importYamlPath)
	if err != nil {
		return nil, err
	}

	if len(yamlContent) == 0 {
		return nil, errors.New(i18n.ImportYamlCorrupted)
	}

	fmt.Println(constants.Success + i18n.ImportYamlOk)
	return yamlContent, nil
}
