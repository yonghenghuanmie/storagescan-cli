# storagescan-cli
A tool for read data(include private data) from blockchain(theoretically any blockchain if only EVM compatible) through command-line.
## Compile
1. Install and config go development tool suit. For detail see https://github.com/golang/go.  
2. Use `go get -u github.com/spf13/cobra@latest` to install cobra. For detail see https://github.com/spf13/cobra.  
3. Use `go get -u github.com/MetaplasiaTeam/storagescan@latest` to install storagescan. For detail see https://github.com/MetaplasiaTeam/storagescan.  
4. Use `go build -o storagescan-cli.exe` to generate executable file.  
## Usage
1. Deploy contracts to blockchain.  
2. Generate storage_layout json strings by solc compiler.    
`solc --storage-layout storage_scan_examples.sol`
3. Using deployed contracts address and generated storage layout file and variable names as parameter.
## Example
Contract example: https://github.com/MetaplasiaTeam/storagescan/blob/main/README.md  
`storagescan-cli ls 0x24302f327764f94c15d930f5Ac70D362B4a156F9 storage_layout.json int1 string1`  
`storagescan-cli ls 0x24302f327764f94c15d930f5Ac70D362B4a156F9 storage_layout.json slice[1] mapping[1] i.id`  
`storagescan-cli ls 0x24302f327764f94c15d930f5Ac70D362B4a156F9 storage_layout.json array5[0].id`  
`storagescan-cli ls query_data.json`  
`storagescan-cli get rpc`  
`storagescan-cli set rpc https://ropsten.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161`  
