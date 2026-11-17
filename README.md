# Graph solution (Go)

    chmod +x graph_solution   # if needed
    ./graph_solution test_cases/graph1.csv

First run builds the binary automatically with `go build` (requires Go, https://go.dev/dl/); later runs
reuse the cached `graph_solution_bin`.

- Parsing: CSV `source,target`; node set derived from all edges; IDs remapped to dense indices.
- `is_dag`: Kahn's algorithm (iterative; self-edges count as cycles).
- Degrees: counted from adjacency lists.
- PageRank: d=0.85, 20 iterations, uniform init. The row-stochastic transition matrix is applied
  sparsely; dangling nodes spread their rank uniformly over all N nodes (including themselves,
  which is what the spec's `1,2` example requires).
- Run `./test.sh` to compare against all `test_cases/*_output.txt`.

## Note on test_cases/graph3
`graph1` and `graph2` match the expected output exactly. For `graph3`, `is_dag`, `max_in_degree`
and `max_out_degree` match, but `pr_max`/`pr_min` differ slightly (mine: 0.003383 / 0.000007,
expected: 0.003388 / 0.000001). The expected values correspond to PageRank after **1** iteration
rather than the **20** the specification requires. This implementation follows the spec
(d = 0.85, 20 iterations, uniform start, dangling nodes spread uniformly over all nodes). The
spec's own examples (3-cycle, and `1,2` -> 0.649123 / 0.350877) are reproduced exactly.

### Iteration count and graph3
For `graph3`, running **1 PageRank iteration** reproduces the expected output exactly
(`pr_max: 0.003388`, `pr_min: 0.000001`), whereas the specified **20 iterations** give
`0.003383` / `0.000007`. I checked every iteration count from 1 to 30: only 1 matches graph3.
(`graph1` and `graph2` match for any iteration count.) However, 1 iteration contradicts the
spec's own sink example (`1,2` -> 0.649123 / 0.350877, which needs 20 iterations; 1 iteration
gives 0.712500 / 0.287500). I therefore kept the specified 20 iterations. To switch to 1,
change `iterations = 20` to `iterations = 1` in `main.go`.
