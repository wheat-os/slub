package generate

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMkdirProjectDir(t *testing.T) {
	dir, err := os.Getwd()
	require.NoError(t, err)

	projectName := "testProject"

	MakeDirProjectDir(dir, projectName)

	spiderDir := path.Join(dir, strings.ToLower(projectName), SpiderDir)
	setDir := path.Join(dir, strings.ToLower(projectName))

	_, err = os.Stat(spiderDir)
	require.NoError(t, err)

	_, err = os.Stat(setDir)
	require.NoError(t, err)

	dir = path.Join(dir, strings.ToLower(projectName))
	os.RemoveAll(dir)

}

const testPro = "temps"

func TestMkdirSpider(t *testing.T) {
	dir, _ := os.Getwd()
	MakeDirProjectDir(dir, testPro)
	p := path.Join(dir, testPro)

	SetProjectPath(p)
	require.NoError(t, MakeSpider("temps", "www.baidu.com", "test"))
	os.RemoveAll(p)
}

func TestMakeSetting(t *testing.T) {
	dir, _ := os.Getwd()
	MakeDirProjectDir(dir, testPro)

	p := path.Join(dir, testPro)
	SetProjectPath(p)
	require.NoError(t, makeSetting())

	os.RemoveAll(p)
}

func TestMakeDefaultSettings(t *testing.T) {
	dir, _ := os.Getwd()
	p := path.Join(dir, testPro)

	SetProjectPath(p)
	MakeDirProjectDir(dir, testPro)

	require.NoError(t, MakeDefaultSettings())

	os.RemoveAll(p)
}
