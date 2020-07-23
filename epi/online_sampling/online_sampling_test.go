package online_sampling_test

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"testing"

	"github.com/stefantds/csvdecoder"

	. "github.com/stefantds/go-epi-judge/epi/online_sampling"
	"github.com/stefantds/go-epi-judge/iterator"
	"github.com/stefantds/go-epi-judge/random"
	"github.com/stefantds/go-epi-judge/utils"
)

func TestOnlineRandomSample(t *testing.T) {
	testFileName := testConfig.TestDataFolder + "/" + "online_sampling.tsv"
	file, err := os.Open(testFileName)
	if err != nil {
		t.Fatalf("could not open file %s: %v", testFileName, err)
	}
	defer file.Close()

	type TestCase struct {
		Stream  []int
		K       int
		Details string
	}

	parser, err := csvdecoder.NewWithConfig(file, csvdecoder.Config{Comma: '\t', IgnoreHeaders: true})
	if err != nil {
		t.Fatalf("could not parse file %s: %s", testFileName, err)
	}

	for i := 0; parser.Next(); i++ {
		tc := TestCase{}
		if err := parser.Scan(
			&tc.Stream,
			&tc.K,
			&tc.Details,
		); err != nil {
			t.Fatal(err)
		}

		t.Run(fmt.Sprintf("Test Case %d", i), func(t *testing.T) {
			if err := onlineRandomSampleWrapper(tc.Stream, tc.K); err != nil {
				t.Error(err)
			}
		})
	}
	if err = parser.Err(); err != nil {
		t.Fatalf("parsing error: %s", err)
	}
}

func onlineRandomSampleWrapper(stream []int, k int) error {
	return random.RunFuncWithRetries(
		func() bool {
			return onlineRandomSampleRunner(stream, k)
		},
		errors.New("the results don't match the expected distribution"),
	)
}

func onlineRandomSampleRunner(stream []int, k int) bool {
	const nbRuns = 1000000

	results := make([][]int, nbRuns)
	for i := 0; i < nbRuns; i++ {
		iter := iterator.New(iterator.Ints(stream))
		results[i] = OnlineRandomSample(iter, k)
	}

	totalPossibleOutcomes := random.BinomialCoefficient(len(stream), k)

	combinations := make([][]int, totalPossibleOutcomes)
	for i := 0; i < totalPossibleOutcomes; i++ {
		combinations[i] = random.ComputeCombinationIdx(stream, k, i)
	}

	sort.Slice(combinations, func(i, j int) bool {
		return utils.LexIntsCompare(combinations[i], combinations[j])
	})

	sequence := make([]int, nbRuns)
	for i, r := range results {
		sort.Ints(r)
		sequence[i] = sort.Search(
			len(combinations),
			func(i int) bool { return !utils.LexIntsCompare(r, combinations[i]) },
		)
	}
	return random.CheckSequenceIsUniformlyRandom(sequence, totalPossibleOutcomes, 0.01)
}