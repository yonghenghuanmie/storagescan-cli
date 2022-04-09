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
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/yonghenghuanmie/storagescan"
)

type QueryData struct {
	Address          string   `json:"address"`
	Layout_file_path string   `json:"layout_file_path"`
	Name             []string `json:"name"`
	Layout           string   `json:"layout"`
}
type QueryArray struct {
	Contracts []QueryData `json:"contracts"`
}

func main() {
	rpc_node := "https://ropsten.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161"
	reference_rpc_node := &rpc_node
	var file_name string
	reference_file_name := &file_name
	query_array := QueryArray{}
	reference_query_array := &query_array

	ReadJsonData := func(file_name string, json_object any) error {
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
		Args: func(cmd *cobra.Command, args []string) error {
			//check parameter number
			if len(*reference_file_name) == 0 && len(args) < 3 {
				return errors.New("Not enough parameter, at least 3 parameters or 1 file path.")
			}

			//construct query data
			if len(*reference_file_name) != 0 {
				err := ReadJsonData(*reference_file_name, reference_query_array)
				if err != nil {
					return err
				}
			}
			if len(args) != 0 {
				(*reference_query_array).Contracts = append((*reference_query_array).Contracts, QueryData{args[0], args[1], args[2:], ""})
			}

			for i := 0; i < len((*reference_query_array).Contracts); i++ {
				query_array := (*reference_query_array).Contracts[i]
				//check data format
				if strings.Compare(query_array.Address[0:2], "0x") != 0 {
					return errors.New("Only support hex format.")
				}
				if len(query_array.Address) != 42 {
					return errors.New("Not valid hex data.")
				}

				//fill layout data if layout_file_path is specified
				if len(query_array.Layout_file_path) != 0 {
					layout_data, err := ioutil.ReadFile(query_array.Layout_file_path)
					if err != nil {
						return errors.New("Failed to open file:" + file_name)
					}
					if !json.Valid(layout_data) {
						return errors.New("Invalid json file:" + file_name)
					}
					(*reference_query_array).Contracts[i].Layout = string(layout_data)
				}
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			for i := 0; i < len((*reference_query_array).Contracts); i++ {
				query_array := (*reference_query_array).Contracts[i]
				contract := storagescan.NewContract(common.HexToAddress(query_array.Address), *reference_rpc_node)
				err := contract.ParseByStorageLayout(query_array.Layout)
				if err != nil {
					fmt.Println("Parse went wrong. " + err.Error())
					return
				}

				for j := 0; j < len(query_array.Name); j++ {
					fmt.Printf("%v:%v\n", query_array.Name[j], contract.GetVariableValue(query_array.Name[j]))
				}
			}
		},
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
			*reference_rpc_node = args[1]
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
			fmt.Println("rpc node:" + *reference_rpc_node)
		},
	}

	var root_cmd = &cobra.Command{Use: "storagescan-cli"}
	root_cmd.AddCommand(cmd_ls, cmd_get, cmd_set)
	root_cmd.Execute()
}
