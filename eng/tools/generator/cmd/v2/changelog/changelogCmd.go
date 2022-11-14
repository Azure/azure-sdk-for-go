package changelog

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/autorest"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/autorest/model"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/exports"
	"github.com/spf13/cobra"
	"log"
	"regexp"
	"strings"
)

type filterOperation func(changelog *model.Changelog)

func Filter(changelog *model.Changelog, opts ...filterOperation) {
	if changelog.Modified != nil {
		for _, opt := range opts {
			opt(changelog)
		}
	}

	fmt.Println("ChangeLog:", changelog.String())
}

func EnumFilter(changelog *model.Changelog) {
	if changelog.Modified.HasBreakingChanges() {
		for typeAliases := range changelog.Modified.BreakingChanges.TypeAliases {
			constKeys, constExist := searchKey(changelog.Modified.AdditiveChanges.Consts, typeAliases, "")
			funcKeys, funcExist := searchKey(changelog.Modified.AdditiveChanges.Funcs, typeAliases, "Possible")

			if constExist && funcExist && len(funcKeys) == 1 {
				fmt.Println(constKeys, funcKeys)
				for _, c := range constKeys {
					delete(changelog.Modified.BreakingChanges.Consts, c)
				}
				for _, f := range funcKeys {
					delete(changelog.Modified.BreakingChanges.Funcs, f)
				}
			}
		}
	}

	if changelog.Modified.HasAdditiveChanges() {
		for typeAliases := range changelog.Modified.AdditiveChanges.TypeAliases {
			constKeys, constExist := searchKey(changelog.Modified.AdditiveChanges.Consts, typeAliases, "")
			funcKeys, funcExist := searchKey(changelog.Modified.AdditiveChanges.Funcs, typeAliases, "Possible")

			if constExist && funcExist && len(funcKeys) == 1 {
				for _, c := range constKeys {
					delete(changelog.Modified.AdditiveChanges.Consts, c)
				}
				for _, f := range funcKeys {
					delete(changelog.Modified.AdditiveChanges.Funcs, f)
				}
			}
		}
	}
}

type FilterInterface interface {
	exports.Const | exports.Func
}

func searchKey[T FilterInterface](m map[string]T, key1, prefix string) ([]string, bool) {
	keys := make([]string, 0)
	for k := range m {
		if regexp.MustCompile(fmt.Sprintf("^%s%s\\w*", prefix, key1)).MatchString(k) {
			keys = append(keys, k)
		}
	}
	if len(keys) != 0 {
		return keys, true
	}
	return nil, false
}

func OperationFiler(changelog *model.Changelog) {
	if changelog.Modified.HasAdditiveChanges() {
		for funcName, funcValue := range changelog.Modified.AdditiveChanges.Funcs {
			clientFunc := strings.Split(funcName, ".")
			if len(clientFunc) == 1 || clientFunc[1] == "MarshalJSON" || clientFunc[1] == "UnmarshalJSON" {
				continue
			}
			// 获取最后一个参数
			ps := strings.Split(*funcValue.Params, ",")
			clientFuncOptions := ps[len(ps)-1]
			clientFuncOptions = strings.TrimLeft(strings.TrimSpace(clientFuncOptions), "*")
			// 获取第一个返回值
			rs := strings.Split(*funcValue.Returns, ",")
			clientFuncResponse := rs[0]
			if strings.Contains(clientFunc[1], "Begin") {
				re := regexp.MustCompile("\\[(?P<response>.*)\\]")
				clientFuncResponse = re.FindString(clientFuncResponse)
				clientFuncResponse = re.ReplaceAllString(clientFuncResponse, "${response}")
			} else {
				clientFuncResponse = strings.TrimLeft(clientFuncResponse, "*")
			}
			// remove
			if clientFuncOptions != "" && clientFuncResponse != "" {
				delete(changelog.Modified.AdditiveChanges.Structs, clientFuncOptions)
				delete(changelog.Modified.AdditiveChanges.Structs, clientFuncResponse)
				for i, v := range changelog.Modified.AdditiveChanges.CompleteStructs {
					if v == clientFuncOptions {
						changelog.Modified.AdditiveChanges.CompleteStructs = append(changelog.Modified.AdditiveChanges.CompleteStructs[:i-1], changelog.Modified.AdditiveChanges.CompleteStructs[i+2:]...)
					}
				}
			}
			fmt.Println(clientFuncResponse)
		}
	}
}

// LROFilter 在 OperationFilter 之后
func LROFilter(config *model.Changelog) {
	if config.Modified.HasBreakingChanges() && config.Modified.HasAdditiveChanges() {
		for bFunc := range config.Modified.BreakingChanges.Funcs {
			clientFunc := strings.Split(bFunc, ".")
			if len(clientFunc) == 1 || clientFunc[1] == "MarshalJSON" || clientFunc[1] == "UnmarshalJSON" {
				continue
			}

			beginFunc := fmt.Sprintf("%s.Begin%s", clientFunc[0], clientFunc[1])
			if _, ok := config.Modified.AdditiveChanges.Funcs[beginFunc]; ok {
				structName := fmt.Sprintf("%s%sOpetions", strings.TrimLeft(strings.TrimSpace(clientFunc[0]), "*"), clientFunc[1])
				if _, structOk := config.Modified.BreakingChanges.Structs[structName]; structOk {
					delete(config.Modified.BreakingChanges.Structs, structName)
					delete(config.Modified.AdditiveChanges.Funcs, beginFunc)
				}
			}
		}
	}
}

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "changelog",
		Short: "Generate a go readme file or add go track2 part to go readme file according to base swagger readme file",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			log.SetFlags(0) // remove the time stamp prefix
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			oriPackagePath := "D:\\Go\\src\\github.com\\Azure\\azure-sdk-for-go\\sdk\\resourcemanager\\elastic\\armelastic"
			newPackagePath := "D:\\Go\\src\\github.com\\Azure\\dev\\azure-sdk-for-go\\sdk\\resourcemanager\\elastic\\armelastic"

			oriExports, err := exports.Get(oriPackagePath)
			newExports, err := exports.Get(newPackagePath)

			changelog, err := autorest.GetChangelogForPackage(&oriExports, &newExports)
			if err != nil {
				log.Fatalln(err)
			}

			_ = changelog
			Filter(changelog, OperationFiler, LROFilter)
			return nil
		},
		SilenceUsage: true, // this command is used for a pipeline, the usage should never show
	}

	return cmd
}
