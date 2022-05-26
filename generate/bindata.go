package generate

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

func getNewSpiderFuncByFile(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	reg, err := regexp.Compile(`func\s+(\w+)\w?\s?\(\)\s?spider.Spider`)
	if err != nil {
		return "", err
	}

	regName := reg.FindAllStringSubmatch(string(content), 1)
	if len(regName) == 0 {
		msg := `"this spider does not have an initialization creation 
		function, please define or check the file. spider path:\n%s"`
		return "", fmt.Errorf(msg, file)
	}

	return regName[0][1], nil
}

func MakeRegisterSpider(spiders []string) error {
	spiderPath := path.Join(projectPath, SpiderDir)

	if len(spiders) == 0 {
		var walkFunc = func(spiderFile string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				spiders = append(spiders, path.Base(spiderFile))
			}
			return nil
		}

		if err := filepath.Walk(spiderPath, walkFunc); err != nil {
			return err
		}
	}

	// 创建 register
	projectName := strings.ToLower(path.Base(projectPath))

	// 获取 New 爬虫创建方法
	makeSpiderFunc := make([]string, 0)

	for _, sp := range spiders {
		funcName, err := getNewSpiderFuncByFile(path.Join(spiderPath, sp))
		if err != nil {
			return err
		}
		makeSpiderFunc = append(makeSpiderFunc, funcName)
	}

	data := map[string]interface{}{
		"projectName": projectName,
		"projectMod":  projectMod,
		"spiders":     makeSpiderFunc,
	}

	return generateTemplate(path.Join(projectPath, "register_gen.go"), data, registerTemp, "register")
}
