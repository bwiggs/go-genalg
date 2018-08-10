package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"sort"
)

const (
	populationSize = 25
)

var input string
var inputLen float64
var mutationRate float64

type genotype string
type pool []genotype

// Init handles parsing the command line flags
func init() {
	flag.StringVar(&input, "i", "", "the input string to generate")
	flag.Float64Var(&mutationRate, "m", 0.01, "the mutation rate")
	flag.Parse()
	if len(input) == 0 {
		log.Fatal("expected input string")
	}
	populationSize := len(input)
	inputLen = float64(populationSize)
}

func main() {

	// Create a population of N elements, each with randomly generated DNA. ex:
	// (main.population) (len=25 cap=25) {
	// 	(main.genotype) (len=15) "^vR.z@b:\\/9eGw6",
	// 	(main.genotype) (len=15) "j'&'VM\"[-D[D^qL",
	// 	(main.genotype) (len=15) "JTsc8sCfHc^befs",
	// 	(main.genotype) (len=15) "KYj^clJ:C!U\"?:B",
	// 	(main.genotype) (len=15) "ytg1,G<K,0RX=E[",
	// 	(main.genotype) (len=15) "+mXl]C6<w,p2eYr",
	// 	(main.genotype) (len=15) "&F;3T2e&4L;Frwk",
	// }
	population := populate(populationSize, len(input))

	// init generation
	generation := 1
	fittest := population[0]

	for {
		// Selection. Evaluate the fitness of each element of the population and
		// build a mating pool. Strings are sorted to how well they match the input.
		sort.Sort(population)

		best := population[len(population)-1]
		fittest = best
		found := string(best) == input

		if generation%10 == 0 || found {
			// time.Sleep(10 * time.Millisecond)
			clearScreen()
			fmt.Printf("gen: % 6d %s\n", generation, input)
			fmt.Printf("gen: % 6d %s\n", generation, fittest.matches())
			fmt.Printf("gen: % 6d %s %3f\n\n", generation, fittest, fittest.score())
		}

		if found {
			break
		}

		// Reproduction. Repeat N times:
		// a) Pick two parents with probability according to relative fitness.
		// b) Crossover—create a “child” by combining the DNA of these two parents.
		// c) Mutation—mutate the child’s DNA based on a given probability.
		// d) Add the new child to a new population.
		// Step 4. Replace the old population with the new population and return to Step 2.

		children := make(pool, populationSize)
		for i := 0; i < populationSize; i++ {

			// breed the best item with a random member of the population
			p1 := best
			p2 := population[rand.Intn(populationSize)]

			// randomly pull a char from each parent
			child := make([]byte, len(p1))
			for i := range p1 {
				// crossover - 80% change of pulling from best parent
				if rand.Float64() > .2 {
					child[i] = p1[i]
				} else {
					child[i] = p2[i]
				}

				// mutation
				if rand.Float64() < mutationRate {
					child[i] = randASCII()
				}
			}

			children[i] = genotype(child)
		}
		population = children
		generation++
	}
}

func clearScreen() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

// populate returns a randomized population string of size strSize
func populate(populationSize, strSize int) (p pool) {
	p = make(pool, populationSize)
	for i := range p {
		var str []byte
		for i := 0; i < strSize; i++ {
			str = append(str, randASCII())
		}
		p[i] = genotype(str)
	}
	return
}

// randAscii returns a random ascii character between space and lowercase z. this
// includes all special characters except brackets in the ascii space
func randASCII() byte {
	return byte(rand.Intn(123-32) + 32)
}

// score the string based on how many characters match the target string.
func (g genotype) score() float64 {
	f := 0.0
	for i := range g {
		if g[i] == input[i] {
			f++
		}
	}
	return f / inputLen
}

func (g genotype) matches() (matches string) {
	b := make([]byte, len(g))
	for i := range g {
		if g[i] == input[i] {
			b[i] = '|'
		} else {
			b[i] = ' '
		}
	}
	return string(b)
}

func (p pool) Len() int           { return len(p) }
func (p pool) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p pool) Less(i, j int) bool { return p[i].score() < p[j].score() }
