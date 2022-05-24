package generate

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"
)

const (
	SpiderDir = "spiders"
)

func generateTemplate(savePath string, data interface{}, temp string, name string) error {
	f, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer f.Close()

	t, err := template.New(name).Parse(temp)
	if err != nil {
		return err
	}

	return t.Execute(f, data)
}

// MkdirProjectDir generate project dir.
func MakeDirProjectDir(proPath, name string) error {
	dir := path.Join(proPath, strings.ToLower(name))
	spiderDir := path.Join(dir, SpiderDir)
	settingDir := path.Join(dir, strings.ToLower(name))

	if err := os.Mkdir(dir, 0777); err != nil {
		return err
	}

	// os.Chmod(dir, 0777)
	if err := os.Mkdir(spiderDir, 0777); err != nil {
		return err
	}

	if err := os.Mkdir(settingDir, 0777); err != nil {
		return err
	}

	return nil
}

// generate simple spider
func MakeSpider(name string, fqdn string, uid string) error {
	nameCap := strings.ToUpper(name[:1]) + strings.ToLower(name[1:])

	data := map[string]string{
		"nameCap": nameCap,
		"uid":     uid,
		"fqdn":    fqdn,
	}

	// open name
	spiderPath := path.Join(projectPath, SpiderDir, strings.ToLower(name)+".go")
	_, exist := os.Stat(spiderPath)
	if !os.IsNotExist(exist) {
		return fmt.Errorf("spider already exists, spider path: %s", spiderPath)
	}

	return generateTemplate(spiderPath, data, spiderTemp, "spider")
}

func makeSetting() error {
	proSetPath := path.Join(projectPath, path.Base(projectPath))
	if _, err := os.Stat(proSetPath); os.IsNotExist(err) {
		return fmt.Errorf("not a legitimate slub scaffolding, path: %s", projectPath)
	}

	name := path.Base(proSetPath)

	data := map[string]string{
		"nameCap": strings.ToUpper(name[:1]) + strings.ToLower(name[1:]),
		"pkgName": strings.ToLower(name),
	}

	return generateTemplate(path.Join(proSetPath, "settings.go"), data, settingTemp, "setting")
}

func makeItem() error {
	proSetPath := path.Join(projectPath, path.Base(projectPath))
	if _, err := os.Stat(proSetPath); os.IsNotExist(err) {
		return fmt.Errorf("not a legitimate slub scaffolding, path: %s", projectPath)
	}

	name := path.Base(proSetPath)

	data := map[string]string{
		"nameCap": strings.ToUpper(name[:1]) + strings.ToLower(name[1:]),
		"pkgName": strings.ToLower(name),
	}

	return generateTemplate(path.Join(proSetPath, "items.go"), data, itemTemp, "item")
}

func makePipline() error {
	proSetPath := path.Join(projectPath, path.Base(projectPath))
	if _, err := os.Stat(proSetPath); os.IsNotExist(err) {
		return fmt.Errorf("not a legitimate slub scaffolding, path: %s", projectPath)
	}

	name := path.Base(proSetPath)

	data := map[string]string{
		"nameCap": strings.ToUpper(name[:1]) + strings.ToLower(name[1:]),
		"pkgName": strings.ToLower(name),
	}

	return generateTemplate(path.Join(proSetPath, "piplines.go"), data, piplineTemp, "pipline")
}

// generate pipline
func makeMiddleware() error {
	proSetPath := path.Join(projectPath, path.Base(projectPath))
	if _, err := os.Stat(proSetPath); os.IsNotExist(err) {
		return fmt.Errorf("not a legitimate slub scaffolding, path: %s", projectPath)
	}

	name := path.Base(proSetPath)

	data := map[string]string{
		"nameCap": strings.ToUpper(name[:1]) + strings.ToLower(name[1:]),
		"pkgName": strings.ToLower(name),
	}

	return generateTemplate(path.Join(proSetPath, "middleware.go"), data, middlewareTemp, "middle")
}

func makeConfig() error {
	return generateTemplate(path.Join(projectPath, "conf.toml"), nil, confTem, "conf")
}

func makeMain() error {

	name := path.Base(projectPath)
	data := map[string]string{
		"confFile": path.Join(projectPath, "conf.toml"),
		"mod":      projectMod,
		"pkgName":  strings.ToLower(name),
	}

	return generateTemplate(path.Join(projectPath, "main.go"), data, mainTemp, "main")
}

func MakeDefaultSettings() error {
	exit := make([]func() error, 0)
	exit = append(
		exit,
		makeSetting,
		makeItem,
		makeMiddleware,
		makePipline,
		makeConfig,
		makeMain,
	)

	for _, e := range exit {
		if err := e(); err != nil {
			return err
		}
	}

	// go mod init

	dos := exec.Command("go", "mod", "init", projectMod)
	dos.Dir = projectPath

	fmt.Println(dos.String())
	return dos.Run()
}
