# storagescan-cli
A tool for read data(include private data) from blockchain(theoretically any blockchain if only EVM compatible) through command-line.
## Compile
1. Install and config go development tool suit. For detail see https://github.com/golang/go.  
2. Use `go get -u github.com/spf13/cobra@latest` to install cobra. For detail see https://github.com/spf13/cobra.  
3. Use `go get -u github.com/MetaplasiaTeam/storagescan@latest` to install storagescan. For detail see https://github.com/MetaplasiaTeam/storagescan.  
4. Use `go build -o storagescan-cli.exe` to generate executable file.
## Usage
`storagescan-cli ls 0x24302f327764f94c15d930f5Ac70D362B4a156F9 storage_layout.json int1 string1`  
`storagescan-cli ls query_data.json`  
`storagescan-cli get rpc`  
`storagescan-cli set rpc https://ropsten.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161`  
