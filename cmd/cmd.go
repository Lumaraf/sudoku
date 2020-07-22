package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/lumaraf/sudoku-solver/definition"
	"github.com/lumaraf/sudoku-solver/generator"
	"github.com/lumaraf/sudoku-solver/grid"
	"gopkg.in/yaml.v2"
	"os"
	"time"
	"unsafe"
)

type args struct {
	DefinitionFile string `arg:"positional" help:"write memory profile to 'file''"`
}

func main() {
	// Parse args
	args := args{
		"",
	}
	arg.MustParse(&args)

	var def definition.Definition
	var err error
	if args.DefinitionFile == "" {
		decoder := yaml.NewDecoder(os.Stdin)
		if err = decoder.Decode(&def); err != nil {
			fmt.Println(err)
			return
		}
	} else if def, err = definition.FromFile(args.DefinitionFile); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(def.Name, "by", def.Author)
	fmt.Println()

	solutions := []grid.Grid{}
	start := time.Now()
	generator.Generate(def.Rules, func(g grid.Grid) bool {
		solutions = append(solutions, g)
		return len(solutions) < 2
	})
	fmt.Println("runtime:", time.Now().Sub(start))

	if len(solutions) >= 1 {
		if len(solutions) > 1 {
			fmt.Println("WARNING: multiple valid solutions found")
			fmt.Println()
		}
		solutions[0].Print()
	}

	fmt.Println(unsafe.Sizeof(generator.GeneratorState{}))

	//memStats := runtime.MemStats{}
	//runtime.ReadMemStats(&memStats)
	//fmt.Printf("%+v\n", memStats)
}
