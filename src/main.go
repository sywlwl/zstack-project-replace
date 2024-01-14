package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

var (
	projectPath string
	name        string
	replace     string
)

func isDirExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func replaceContent(filename string, oldStr, newStr string) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		//err
		return
	}
	content := string(buf)
	//替换
	newContent := strings.Replace(content, oldStr, newStr, -1)

	//重新写入
	ioutil.WriteFile(filename, []byte(newContent), 0)
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "zpr",
		Short: "Zigbee Zstatic Project name replace",
		Long:  "Zigbee Zstatic Project name replace",
		Run: func(cmd *cobra.Command, args []string) {
			if !isDirExists(projectPath) {
				fmt.Printf("指定的项目路径不存在 %s\r\n", projectPath)
				os.Exit(-1)
			}
			if !isDirExists(path.Join(projectPath, name)) {
				fmt.Printf("指定的项目路径下不存在指定的项目 %s\r\n", name)
				os.Exit(-1)
			}

			if isDirExists(path.Join(projectPath, replace)) {
				fmt.Printf("指定的项目路径下已存在要替换的项目 %s\r\n", replace)
				os.Exit(-1)
			}

			destPath := path.Join(projectPath, name)
			files := [][]string{
				[]string{
					path.Join(destPath, "Source", "OSAL_"+name+".c"),
					path.Join(destPath, "Source", "OSAL_"+replace+".c"),
				},
				[]string{
					path.Join(destPath, "Source", name+".c"),
					path.Join(destPath, "Source", replace+".c"),
				},
				[]string{
					path.Join(destPath, "Source", name+".h"),
					path.Join(destPath, "Source", replace+".h"),
				},
				[]string{
					path.Join(destPath, "Source", name+"Hw.h"),
					path.Join(destPath, "Source", replace+"Hw.h"),
				},

				[]string{
					path.Join(destPath, "CC2530DB", name+".eww"),
					path.Join(destPath, "CC2530DB", replace+".eww"),
				},
				[]string{
					path.Join(destPath, "CC2530DB", name+".ewp"),
					path.Join(destPath, "CC2530DB", replace+".ewp"),
				},
				[]string{
					path.Join(destPath, "CC2530DB", name+".ewd"),
					path.Join(destPath, "CC2530DB", replace+".ewd"),
				},
				[]string{
					path.Join(destPath, "CC2530DB", "Source", name+"Hw.c"),
					path.Join(destPath, "CC2530DB", "Source", replace+"Hw.c"),
				},
			}

			for _, p := range files {
				replaceContent(p[0], name, replace)
				// 修改文件名
				os.Rename(p[0], p[1])
			}
			os.Rename(path.Join(projectPath, name), path.Join(projectPath, replace))
			fmt.Println("项目替换完成")
		},
	}

	rootCmd.PersistentFlags().StringVarP(&projectPath, "path", "p", "", "指定不包含项目名的路径，即项目的上级目录")
	rootCmd.PersistentFlags().StringVarP(&name, "name", "n", "", "指定项目名")
	rootCmd.PersistentFlags().StringVarP(&replace, "replace", "r", "", "指定要替换的项目名")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
