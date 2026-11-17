package main

// imports
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// constants
const (
	damping    = 0.85
	iterations = 20
)

type graph struct {
	n      int
	outAdj [][]int
	inDeg  []int
}

func parse(path string) (*graph, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	ids := map[uint64]int{}
	intern := func(id uint64) int {
		if i, ok := ids[id]; ok {
			return i
		}
		i := len(ids)
		ids[id] = i
		return i
	}
	var edges [][2]int
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1<<20), 1<<20)
	for ln := 1; sc.Scan(); ln++ {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			return nil, fmt.Errorf("line %d: expected 'source,target'", ln)
		}
		a, err1 := strconv.ParseUint(strings.TrimSpace(parts[0]), 10, 64)
		b, err2 := strconv.ParseUint(strings.TrimSpace(parts[1]), 10, 64)
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("line %d: invalid node id", ln)
		}
		edges = append(edges, [2]int{intern(a), intern(b)})
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	g := &graph{n: len(ids)}
	g.outAdj = make([][]int, g.n)
	g.inDeg = make([]int, g.n)
	for _, e := range edges {
		g.outAdj[e[0]] = append(g.outAdj[e[0]], e[1])
		g.inDeg[e[1]]++
	}
	return g, nil
}

func (g *graph) isDAG() bool {
	indeg := append([]int(nil), g.inDeg...)
	stack := []int{}
	for i, d := range indeg {
		if d == 0 {
			stack = append(stack, i)
		}
	}
	removed := 0
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		removed++
		for _, v := range g.outAdj[u] {
			indeg[v]--
			if indeg[v] == 0 {
				stack = append(stack, v)
			}
		}
	}
	return removed == g.n
}

func (g *graph) pageRank() []float64 {
	n := float64(g.n)
	rank := make([]float64, g.n)
	next := make([]float64, g.n)
	for i := range rank {
		rank[i] = 1 / n
	}
	for it := 0; it < iterations; it++ {
		for i := range next {
			next[i] = 0
		}
		dangling := 0.0
		for u, outs := range g.outAdj {
			if len(outs) == 0 {
				dangling += rank[u]
				continue
			}
			share := rank[u] / float64(len(outs))
			for _, v := range outs {
				next[v] += share
			}
		}
		base := (1-damping)/n + damping*dangling/n
		for v := range next {
			next[v] = base + damping*next[v]
		}
		rank, next = next, rank
	}
	return rank
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: graph_solution path/to/graph.csv")
		os.Exit(2)
	}
	g, err := parse(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	if g.n == 0 {
		fmt.Fprintln(os.Stderr, "error: empty graph")
		os.Exit(1)
	}
	maxIn, maxOut := 0, 0
	for i := 0; i < g.n; i++ {
		if g.inDeg[i] > maxIn {
			maxIn = g.inDeg[i]
		}
		if len(g.outAdj[i]) > maxOut {
			maxOut = len(g.outAdj[i])
		}
	}
	pr := g.pageRank()
	hi, lo := pr[0], pr[0]
	for _, v := range pr {
		if v > hi {
			hi = v
		}
		if v < lo {
			lo = v
		}
	}
	fmt.Printf("is_dag: %t\n", g.isDAG())
	fmt.Printf("max_in_degree: %d\n", maxIn)
	fmt.Printf("max_out_degree: %d\n", maxOut)
	fmt.Printf("pr_max: %.6f\n", hi)
	fmt.Printf("pr_min: %.6f\n", lo)
}
