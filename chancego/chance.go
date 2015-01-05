package chancego

import (
	"math"
	"math/rand"
	"time"
	"github.com/seehuhn/mt19937"
	"errors"
	"strings"
	//"fmt"
	//"strconv"
)

const MAX_INT = math.MaxInt64
const MIN_INT = -MAX_INT
const NUMBERS = "0123456789"
const CHARS_LOWER = "abcdefghijklmnopqrstuvwxyz"
const CHARS_UPPER = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const HEX_POOL = NUMBERS + "abcdef"

type Chance struct {
	rng *rand.Rand
}

var (
	LikelihoodOutOfRangeError = errors.New("Chance: Likelihood accepts values from 0 to 100.")
	MinGreaterThanMaxError = errors.New("Chance: Min cannot be greater than Max.")
	NegativeLengthError = errors.New("Chance: length has to be bigger than 0")
	EmptyPoolError = errors.New("Chance: pool cannot be empty")
	UnknownSelectionError = errors.New("Chance: Cannot specify both alpha and symbols.")
	EmptyArrayError = errors.New("Chance: Cannot work on an empty array")
	UnequalArrayLengthsError = errors.New("Chance: Length of the arrays need to be equal")
	UnknownError = errors.New("Chance: Oh dear. This is unhelpful.")
)

// ** Private functions **

func containsInt(s []int, e int) bool {
	for _, a := range s { if a == e { return true } }
	return false
}

type arrayIntPredicate func(int, int) bool //key, val
func checkArrayIntPredicate(s []int, predicate arrayIntPredicate) bool {
	for k, a := range s { if predicate(k, a) { return true } }
	return false
}

type arrayFloatPredicate func(int, float64) bool //key, val
func checkArrayFloatPredicate(s []float64, predicate arrayFloatPredicate) bool {
	for k, a := range s { if predicate(k, a) { return true } }
	return false
}

func (chance *Chance) random() float64 {
	val := chance.rng.Float64()
	//fmt.Println("rng: " + strconv.FormatFloat(val, 'g', -1, 64))
	return val
}

// ** Public functions **

func NewChance() *Chance {
	chance := new(Chance)
	chance.rng = rand.New(mt19937.New())
	chance.rng.Seed(time.Now().UnixNano())
	return chance
}

func (chance *Chance) Bool(likelihood int, clamp bool) (bool, error) {
	if !clamp && (likelihood < 0 || likelihood > 100) {
		return false, LikelihoodOutOfRangeError
	} else {
		if likelihood < 0 {
			likelihood = 0
		}
		if likelihood > 100 {
			likelihood = 100
		}
	}
	return int(chance.random() * 100) < likelihood, nil;
}

func (chance *Chance) Integer(min int, max int) (int, error) {
	if min > max {
		return 0, MinGreaterThanMaxError
	}
	//math.Floor necessary for negative numbers.
	return int(math.Floor(chance.random() * float64(max - min + 1) + float64(min))), nil;
}

func (chance *Chance) Float(min float64, max float64) (float64, error) {
	if min > max {
		return 0, MinGreaterThanMaxError
	}

	return chance.random() * (max - min) + min, nil
}

func (chance *Chance) String(length int, pool string) (string, error) {
	if length <= 0 {
		return "", NegativeLengthError
	}

	if len(pool) == 0 {
		return "", EmptyPoolError
	}

	var ret string

	for i := 0; i < length; i ++ {
		var char, err = chance.Character("", pool, false, false)
		if err != nil {
			return "", err
		}
		ret += string(char)
	}

	return ret, nil
}

func (chance *Chance) Character(casing string, pool string, symbols bool, alpha bool) (uint8, error) {
	var temppool, letters string
	var tempsymbols = "!@#$%^&*()[]"

	if alpha && symbols {
		return 0, UnknownSelectionError
	}

	if casing == "lower" {
		letters = CHARS_LOWER
	} else if casing == "upper" {
		letters = CHARS_UPPER
	} else {
		letters = CHARS_LOWER + CHARS_UPPER
	}

	if len(pool) != 0 {
		temppool = pool
	} else if alpha {
		temppool = letters
	} else if symbols {
		temppool = tempsymbols
	} else {
		temppool = letters + NUMBERS + tempsymbols
	}

	var charAt, err = chance.Integer(0, len(temppool) - 1)
	if err != nil {
		return 0, err
	}
	return temppool[charAt], nil
}

func (chance *Chance) Capitalize(word string) string {
	return strings.ToUpper(string(word[0])) + word[1:]
}

//TODO unique
//TODO n
//TODO pad, maybe not necessary

func (chance *Chance) ShuffleInt(arr []int) []int {
	newArr := make([]int, len(arr), (cap(arr)+1)*2)
	copy(newArr, arr)
	for i := len(arr) - 1; i > 0; i-- {
		j := chance.rng.Intn(i)
		newArr[i], newArr[j] = newArr[j], newArr[i]
	}
	return newArr
}

// Returns [count] ints picked at random from an array, up to the length of said array.
func (chance *Chance) PickInt(arr []int, count int) ([]int, error) {
	if len(arr) == 0 {
		return []int{}, EmptyArrayError
	}

	if count <= 1 {
		newArr := make([]int, 1, 4)
		var i, err = chance.Integer(0, len(arr)-1)
		if err == nil {
			newArr[0] = arr[i]
		}
		return newArr, err
	}

	if count >= len(arr) {
		count = len(arr) - 1
	}

	return chance.ShuffleInt(arr)[0:count], nil
}

// Returns a single item from an array with relative weighting of odds
func (chance *Chance) WeightedInt(arr []int, weights []int) (int, error) {
	if len(arr) != len(weights) {
		return 0, UnequalArrayLengthsError
	}

	// Handle weights that are less or equal to zero.
	for weightIndex := len(weights) - 1; weightIndex >= 0; weightIndex-- {
		if weights[weightIndex] <= 0 {
			weights = append(weights[:weightIndex], weights[weightIndex+1:]...)
			arr = append(arr[:weightIndex], arr[weightIndex+1:]...)
		}
	}

	if len(arr) == 0 {
		return 0, EmptyArrayError
	}

	var sum int
	for _, a := range weights { sum += a }

	selected, err := chance.Integer(1, sum);
	if err != nil {
		return 0, err
	}

	var total int
	var chosen int
	selectPredicate := func(k int, a int) bool {
		if selected <= total + a {
			chosen = arr[k]
			return true
		}

		total += a
		return false
	}

	if !checkArrayIntPredicate(weights, selectPredicate) {
		return 0, UnknownError
	}

	return chosen, nil
}

// Returns a single item from an array with relative weighting of odds
func (chance *Chance) WeightedFloat(arr []float64, weights []float64) (float64, error) {
	if len(arr) != len(weights) {
		return 0, UnequalArrayLengthsError
	}

	// Handle weights that are less or equal to zero.
	for weightIndex := len(weights) - 1; weightIndex >= 0; weightIndex-- {
		if weights[weightIndex] <= 0 {
			weights = append(weights[:weightIndex], weights[weightIndex+1:]...)
			arr = append(arr[:weightIndex], arr[weightIndex+1:]...)
		}
	}

	if len(arr) == 0 {
		return 0, EmptyArrayError
	}

	negativePredicate := func(_ int, a float64) bool {
		if a < 1 {
			return true
		}
		return false
	}

	//If any int in weights is < 1, apply scaling
	// LOL FLOATS
	if checkArrayFloatPredicate(weights, negativePredicate) {
		var min float64
		for _, a := range weights { if a < min { min = a } }
		scalingFactor := 1 / min
		for k, a := range weights { weights[k] = a * scalingFactor }
	}

	var sum float64
	for _, a := range weights { sum += a }

	selected, err := chance.Float(1, sum);
	if err != nil {
		return 0, err
	}

	var total float64
	var chosen float64
	selectPredicate := func(k int, a float64) bool {
		if selected <= total + a {
			chosen = arr[k]
			return true
		}

		total += a
		return false
	}

	if !checkArrayFloatPredicate(weights, selectPredicate) {
		return 0, UnknownError
	}

	return chosen, nil
}
