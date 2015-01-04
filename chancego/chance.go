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

type ChanceArguments struct {
	max int
	min int
	len int
}

func NewChance() *Chance {
	chance := new(Chance)
	chance.rng = rand.New(mt19937.New())
	chance.rng.Seed(time.Now().UnixNano())
	return chance
}

func (chance *Chance) Bool(likelihood int, clamp bool) (bool, error) {
	if !clamp && (likelihood < 0 || likelihood > 100) {
		return false, errors.New("Chance: Likelihood accepts values from 0 to 100.")
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

func (chance *Chance) random() float64 {
	val := chance.rng.Float64()
	//fmt.Println("rng: " + strconv.FormatFloat(val, 'g', -1, 64))
	return val
}

func (chance *Chance) Integer(min int, max int) (int, error) {
	if min > max {
		return 0, errors.New("Chance: Min cannot be greater than Max.")
	}
	//math.Floor necessary for negative numbers.
	return int(math.Floor(chance.random() * float64(max - min + 1) + float64(min))), nil;
}

func (chance *Chance) String(length int, pool string) (string, error) {
	if length <= 0 {
		return "", errors.New("Chance: length has to be bigger than 0")
	}

	if len(pool) == 0 {
		return "", errors.New("Chance: pool cannot be empty")
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
		return 0, errors.New("Chance: Cannot specify both alpha and symbols.")
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

func (chance *Chance) PickInt(arr []int, count int) ([]int, error) {
	if len(arr) == 0 {
		return []int{}, errors.New("Chance: Cannot pick() from an empty array")
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
