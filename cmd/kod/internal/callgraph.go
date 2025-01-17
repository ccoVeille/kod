package internal

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dominikbraun/graph/draw"
	"github.com/go-kod/kod/internal/callgraph"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var callgraphCmd = &cobra.Command{
	Use:   "callgraph",
	Short: "generate kod callgraph for your kod application.",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.PrintErr("please input the binary filepath")
			return
		}

		g := lo.Must(callgraph.ReadComponentGraph(args[0]))
		o := lo.Must(cmd.Flags().GetString("o"))
		t := lo.Must(cmd.Flags().GetString("t"))

		switch t {
		case "json":
			data, err := g.AdjacencyMap()
			if err != nil {
				fmt.Println(err)
				return
			}

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			lo.Must0(enc.Encode(data))
		case "dot":
			file := lo.Must(os.Create(o))
			lo.Must0(draw.DOT(g, file))
		default:
			fmt.Println("output type not supported")
		}
	},
}

func init() {
	callgraphCmd.PersistentFlags().String("o", "my-graph.dot", "output file name")
	callgraphCmd.PersistentFlags().String("t", "dot", "output type, support json/dot")

	rootCmd.AddCommand(callgraphCmd)
}
