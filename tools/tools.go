package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/jonseymour/12coins/lib"
	"os"
)

func main() {
	relabel := false
	reverse := false
	normalize := false
	permute := false

	flag.BoolVar(&reverse, "reverse", false, "Reverse the solution.")
	flag.BoolVar(&relabel, "relabel", false, "Relabel solution.")
	flag.BoolVar(&normalize, "normalize", false, "Normalize the solution.")
	flag.BoolVar(&permute, "permute", false, "Generate permutations of the solution.")
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)
	for {
		var line []byte
		var err error

		if line, _, err = reader.ReadLine(); err != nil {
			break
		}
		solver := &lib.Solver{}
		if err := json.Unmarshal(line, &solver); err != nil {
			fmt.Fprintf(os.Stderr, "parsing error: %v", err)
			continue
		}
		if reverse {
			if solver, err = solver.Reverse(); err != nil {
				fmt.Fprintf(os.Stderr, "bad solution: %v", err)
				continue
			}
		}

		if normalize {
			solver = solver.Normalize()
		}

		if relabel {
			solver = solver.Relabel()
		}

		if permute {
			for _, m := range []bool{false, true} {
				for _, p := range lib.Permute([]int{0, 1, 2}) {
					clone := solver.Clone()
					clone.Mirror = m
					clone.Permutation = p
					clone.Weighings = [3][2][]int{clone.Weighings[p[0]], clone.Weighings[p[1]], clone.Weighings[p[2]]}
					clone, _ = clone.Reverse()
					clone = clone.Relabel()
					if errors := lib.TestAll(clone.Decide); len(errors) != 0 {
						panic(fmt.Errorf("errors: %v", errors))
					}
					encoder.Encode(clone)
				}
			}
			continue
		}

		encoder.Encode(solver)
	}
}
