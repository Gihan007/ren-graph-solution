# Graph solution (Go)

## Run
- `chmod +x graph_solution` (if needed)
- `./graph_solution test_cases/graph1.csv`
- First run builds the binary automatically (needs Go: https://go.dev/dl/); later runs reuse the cached `graph_solution_bin`
- `./test.sh` compares against all `test_cases/*_output.txt`

## Approach
- Parsing: CSV `source,target`; nodes taken from all edges, remapped to dense indices
- `is_dag`: Kahn's algorithm (iterative; self-edges count as cycles)
- Degrees: counted from adjacency lists
- PageRank: d=0.85, 20 iterations, uniform start, row-stochastic matrix applied sparsely
- Dangling nodes spread rank over all N nodes, including themselves (needed to match the spec's `1,2` example)

## Test results
- graph1, graph2: match exactly
- graph3: `is_dag`, `max_in_degree`, `max_out_degree` match
- graph3 PageRank: mine `0.003383 / 0.000007`, expected `0.003388 / 0.000001`
- Expected graph3 values match **1** iteration (checked 1–30; only 1 matches), not the 20 in the spec
- 1 iteration breaks the spec's `1,2` example (gives `0.712500 / 0.287500`, spec says `0.649123 / 0.350877`)
- Kept 20 iterations per the spec; to use 1, set `iterations = 1` in `main.go`
