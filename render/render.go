package render

import (
	"fmt"
	"strings"

	"github.com/better0fdead/plugin_demo_tgv3/context"
	"github.com/better0fdead/plugin_demo_tgv3/parser"
)

func Render(pluginCtx context.PluginCtx, parsedPackage parser.PackageInfo) (err error) {
	for _, service := range parsedPackage.Services {
		fmt.Printf("Сервис: %s\n", service.Name)
		for _, method := range service.Methods {
			fmt.Printf("  Метод: %s\n", method.Name)

			if len(method.Parameters) > 0 {
				argStrings := make([]string, len(method.Parameters))
				for i, arg := range method.Parameters {
					argStrings[i] = fmt.Sprintf("%s: %s", arg.Name, arg.Kind)
				}
				fmt.Printf("    Аргументы: %s\n", strings.Join(argStrings, ", "))
			} else {
				fmt.Println("    Аргументы: нет")
			}

			if len(method.Returns) > 0 {
				returnStrings := make([]string, len(method.Returns))
				for i, ret := range method.Returns {
					returnStrings[i] = ret.Kind
				}
				fmt.Printf("    Возвращаемые значения: %s\n", strings.Join(returnStrings, ", "))
			} else {
				fmt.Println("    Возвращаемые значения: нет")
			}
		}
		fmt.Println()
	}

	return err
}
