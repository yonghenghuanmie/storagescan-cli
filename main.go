/*
	@author	huanmie<yonghenghuanmie@gmail.com>
	@date	2022.4.10
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/yonghenghuanmie/storagescan"
)

var rpc_node = "https://ropsten.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161"

const (
	NamePattern uint8 = iota
	BeginPattern
)

type Resolver struct {
	set []*regexp.Regexp
}

func ResolverConstructor() *Resolver {
	return &Resolver{[]*regexp.Regexp{
		regexp.MustCompile(`(.+?)([\.|\[].+)`),
		regexp.MustCompile(`\.([^\.\[]+)(.*)|\[(.+?)\](.*)`),
	}}
}

func (this *Resolver) GetValueName(s string) (value_name, substring string) {
	if match_string := this.set[NamePattern].FindStringSubmatch(s); match_string != nil {
		return match_string[1], match_string[2]
	}
	return s, ""
}

func (this *Resolver) GetFirstParameter(s string) (parameter, substring string) {
	if match_string := this.set[BeginPattern].FindStringSubmatch(s); match_string != nil {
		if match_string[1] == "" {
			return match_string[3], match_string[4]
		} else {
			return match_string[1], match_string[2]
		}
	}
	return "", s
}

type QueryData struct {
	Address          string   `json:"address"`
	Layout_file_path string   `json:"layout_file_path"`
	Name             []string `json:"name"`
	Layout           string   `json:"layout"`
}
type QueryArray struct {
	Contracts []QueryData `json:"contracts"`
}

var file_name string
var query_array = QueryArray{}
var resolver Resolver = *ResolverConstructor()

func ReadJsonData(file_name string, json_object any) error {
	data, err := ioutil.ReadFile(file_name)
	if err != nil {
		return errors.New("Failed to open file:" + file_name)
	}
	if !json.Valid(data) {
		return errors.New("Invalid json file:" + file_name)
	}
	err = json.Unmarshal(data, json_object)
	if err != nil {
		return errors.New("Json unmarshal failed. File:" + file_name)
	}
	return nil
}

func CheckListArgument(cmd *cobra.Command, args []string) error {
	//check parameter number
	if len(file_name) == 0 && len(args) < 3 {
		return errors.New("Not enough parameter, at least 3 parameters or 1 file path.")
	}

	//construct query data
	if len(file_name) != 0 {
		err := ReadJsonData(file_name, &query_array)
		if err != nil {
			return err
		}
	}
	if len(args) != 0 && len(args) >= 3 {
		query_array.Contracts = append(query_array.Contracts, QueryData{args[0], args[1], args[2:], ""})
	}

	for i := 0; i < len(query_array.Contracts); i++ {
		query_data := query_array.Contracts[i]
		//check data format
		if strings.Compare(query_data.Address[0:2], "0x") != 0 {
			return errors.New("Only support hex format.")
		}
		if len(query_data.Address) != 42 {
			return errors.New("Not valid hex data.")
		}

		//fill layout data if layout_file_path is specified
		if len(query_data.Layout_file_path) != 0 {
			layout_data, err := ioutil.ReadFile(query_data.Layout_file_path)
			if err != nil {
				return errors.New("Failed to open file:" + file_name)
			}
			if !json.Valid(layout_data) {
				return errors.New("Invalid json file:" + file_name)
			}
			query_array.Contracts[i].Layout = string(layout_data)
		}
	}
	return nil
}

func RunList(cmd *cobra.Command, args []string) {
	for i := 0; i < len(query_array.Contracts); i++ {
		query_array := query_array.Contracts[i]
		contract := storagescan.NewContract(common.HexToAddress(query_array.Address), rpc_node)
		err := contract.ParseByStorageLayout(query_array.Layout)
		if err != nil {
			fmt.Println("Parse contract went wrong. " + err.Error())
			return
		}

		for j := 0; j < len(query_array.Name); j++ {
			value_name, substring := resolver.GetValueName(query_array.Name[j])
			if _, ok := contract.Variables[value_name]; !ok {
				fmt.Printf("Do not find any value use this name, please check your input or layout file. %v\n", query_array.Name[j])
				continue
			}
			value := contract.GetVariableValue(value_name)

			if substring != "" {
				for {
					var parameter string
					var index uint64
					switch interface_ := value.(type) {
					case storagescan.StructValueI:
						parameter, substring = resolver.GetFirstParameter(substring)
						value = interface_.Field(parameter)

					case storagescan.SliceArrayValueI:
						parameter, substring = resolver.GetFirstParameter(substring)
						index, err = strconv.ParseUint(parameter, 10, 64)
						if err != nil {
							fmt.Println("Parse index went wrong. " + err.Error())
							goto failed
						}
						value = interface_.Index(index)

					case storagescan.MappingValueI:
						parameter, substring = resolver.GetFirstParameter(substring)
						value = interface_.Key(parameter)

					default:
						if substring == "" {
							goto success
						} else {
							fmt.Println("Input format error. " + query_array.Name[j])
							goto failed
						}
					}
				}
			}
		success:
			fmt.Printf("%v:%v\n", query_array.Name[j], value)
		failed:
		}
	}
}

func main() {
	//command ls
	var cmd_ls = &cobra.Command{
		Use:   "ls [<contract Address> <json file> <variable Name [...]>]",
		Short: "List variables",
		Long: `List variables. Json file just need provide Layout_file_path or Layout. File format following the follow example:
	{
		"Contracts":
			[
				{"Address":"0x1","Layout_file_path":"layout1.json","Name":["a","b"],"Layout":"{...}"},
				{"Address":"0x2","Layout_file_path":"layout2.json","Name":["c","d"],"Layout":"{...}"}
			]
	}`,
		Args: CheckListArgument,
		Run:  RunList,
	}
	cmd_ls.Flags().StringVarP(&file_name, "file", "f", "", "Specify list file path.")

	var cmd_set = &cobra.Command{
		Use:   "set <rpc>",
		Short: "Set parameter.",
		Long:  `rpc:Use specified rpc node to get data.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New("Not enough parameter, at least 2 parameters.")
			}
			if args[0] == "rpc" {
				return nil
			}
			return errors.New("Command not found. command:" + args[0])
		},
		Run: func(cmd *cobra.Command, args []string) {
			rpc_node = args[1]
		},
	}

	var cmd_get = &cobra.Command{
		Use:   "get <rpc>",
		Short: "Get parameter.",
		Long:  `rpc:Get specified rpc node which used to get data.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("Not enough parameter, at least 1 parameters.")
			}
			if args[0] == "rpc" {
				return nil
			}
			return errors.New("Command not found. command:" + args[0])
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("rpc node:" + rpc_node)
		},
	}

	var root_cmd = &cobra.Command{Use: "storagescan-cli"}
	root_cmd.AddCommand(cmd_ls, cmd_get, cmd_set)
	root_cmd.Execute()
}
