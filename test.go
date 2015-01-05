package main

import (
	"./chancego"
	"fmt"
	_ "strconv"
)

func main() {
	chance := chancego.NewChance()
	for i := 1; i <= 10; i++ {
		fmt.Println(chance.Bool(50, false))
	}
	fmt.Println(chance.Bool(101, true))
	fmt.Println(chance.Bool(-1, true))
	fmt.Println(chance.Bool(101, false))
	fmt.Println(chance.Bool(-1, false))
	fmt.Println("")

	for i := 1; i <= 10; i++ {
		fmt.Println(chance.Integer(-5, 5))
	}
	fmt.Println(chance.Integer(15, 5))
	fmt.Println("")

	fmt.Println(" -- FLOAT --")
	for i := 1; i <= 10; i++ {
		fmt.Println(chance.Float(-5, 5))
	}
	fmt.Println(chance.Float(15, 5))
	fmt.Println("")

	for i := 1; i <= 10; i++ {
		var char, _ = chance.Character("", "", false, false)
		fmt.Println(string(char))
	}
	fmt.Println(chance.Character("", "", true, true))
	fmt.Println("")

	fmt.Println(chance.String(10, "abcdefghi1234567890"))
	fmt.Println("")

	fmt.Println(chance.Capitalize("word"))
	fmt.Println("")

	orderedArr := []int {1, 2, 3, 4, 5, 6}
	fmt.Println(chance.ShuffleInt(orderedArr))
	fmt.Println("")

	for i := 1; i <= 10; i++ {
		fmt.Println(chance.PickInt(orderedArr, 1))
	}
	fmt.Println(chance.PickInt(orderedArr, 20))
	fmt.Println(chance.PickInt(orderedArr, -20))
	fmt.Println(chance.PickInt([]int{}, 20))
	fmt.Println("")

	for i := 1; i <= 20; i++ {
		fmt.Println(chance.WeightedInt(orderedArr, orderedArr))
	}
	fmt.Println(chance.WeightedInt(orderedArr, []int{-1, -1, -1, -1, -1, -1}))
	fmt.Println(chance.WeightedInt(orderedArr, []int{1, 2}))
	fmt.Println("")
}
