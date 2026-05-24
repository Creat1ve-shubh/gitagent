package sdk

import (
	"github.com/open-gitagent/gitagent/internal/agent"
	"github.com/open-gitagent/gitagent/internal/agent/bench"
	"github.com/open-gitagent/gitagent/internal/agent/diff"
)

type RunOptions = agent.RunOptions

type BenchOptions = bench.Options

func Run(opts RunOptions) (string, error) {
	return agent.Run(opts)
}

func SemanticDiff(dir string) (*diff.Report, error) {
	return diff.SemanticDiff(dir)
}

func RunBench(opts BenchOptions) (*bench.Result, error) {
	return bench.RunBench(opts)
}
