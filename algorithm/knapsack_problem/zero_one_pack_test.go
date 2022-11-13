package knapsack_problem

import (
	"fmt"
	"github.com/Mr-lixh/go-demo/utils"
	"github.com/pkg/errors"
	"math/rand"
	"testing"
	"time"
)

/*
01背包问题
描述：有一个最多能装 m 公斤的背包，现在有 n 件物品，第 i 件物品的重量是 c[i]，价值是 w[i]。求解将哪些物品装入背包可使价值总和达到最大？
*/

var (
	m = 10
	n = 4
	c = []int{2, 3, 4, 7}
	w = []int{1, 3, 5, 9}
)

func TestBaseResolve(t *testing.T) {
	// 初始化数组
	dp := [5][11]int{}
	gp := [5][11]int{}
	for i := 0; i <= n; i++ {
		for j := 0; j <= m; j++ {
			dp[i][j] = 0
			gp[i][j] = 0
		}
	}

	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if j < c[i-1] {
				dp[i][j] = dp[i-1][j]
			} else {
				dp[i][j] = dp[i-1][j]
				if dp[i-1][j-c[i-1]]+w[i-1] > dp[i][j] {
					dp[i][j] = dp[i-1][j-c[i-1]] + w[i-1]
					gp[i][j] = 1
				}
			}
		}
	}

	for i := 0; i <= n; i++ {
		for j := 0; j <= m; j++ {
			fmt.Printf("%d\t", dp[i][j])
		}
		fmt.Println()
	}

OUT:
	for i := n; i >= 1; i-- {
		for j := m; j >= 1; j-- {
			if j >= c[i-1] && dp[i][j] == dp[i-1][j-c[i-1]]+w[i-1] {
				fmt.Printf("%d\t", i)
				m = m - c[i-1]
				continue OUT
			}
		}
	}

}

/*
动态规划法求解一维0-1背包问题
*/

type Item struct {
	Index         int
	Weight        int
	Value         int
	ValueByWeight float64
}

func printItems(items []Item) {
	for i := 0; i < len(items); i++ {
		fmt.Printf("ITEM: %d VALOR: %d PESO: %d VALOR/PESO: %f\n", items[i].Index, items[i].Value, items[i].Weight, items[i].ValueByWeight)
	}
}

func TestDynamicKnapsackSolve(t *testing.T) {
	items := generateItems(50000, 50000, 50000)

	startTime := time.Now()
	knapsack, totalValue, totalWeight := dynamicKnapsackSolve(items, 100000)
	elapsedTime := time.Since(startTime).Seconds()

	fmt.Printf("Total time: %.2f\n", elapsedTime)
	utils.PrintMemUsage()

	printItems(knapsack)
	fmt.Printf("TotalValue: %d, TotalWeight: %d\n", totalValue, totalWeight)
}

func generateItems(numberOfItems, maxWeight, maxValue int) []Item {
	var items []Item
	//rand.Seed(time.Now().UnixNano())
	rand.Seed(1)
	for i := 0; i < numberOfItems; i++ {
		items = append(items, Item{
			Index:         i,
			Weight:        rand.Intn(maxWeight) + 1,
			Value:         rand.Intn(maxValue),
			ValueByWeight: 0,
		})
	}
	return items
}

// 动态规划法求解一维0-1背包问题
func dynamicKnapsackSolve(items []Item, knapSackCapacity int) ([]Item, int, int) {
	var knapsack []Item
	table := make([][]int, len(items)+1)
	for i := 0; i < len(items)+1; i++ {
		line := make([]int, knapSackCapacity+1)
		table[i] = line

		for j := 0; j < knapSackCapacity+1; j++ {
			if i == 0 || j == 0 {
				continue
			}
			if items[i-1].Weight <= j {
				aux := items[i-1].Value + table[i-1][j-items[i-1].Weight]
				if aux > table[i-1][j] {
					table[i][j] = aux
				} else {
					table[i][j] = table[i-1][j]
				}
			} else {
				table[i][j] = table[i-1][j]
			}
		}

	}
	currentLine := len(items)
	currentColumn := knapSackCapacity
	totalWeight := 0

	for currentLine > 0 && currentColumn > 0 {
		if table[currentLine][currentColumn] != table[currentLine-1][currentColumn] {
			knapsack = append(knapsack, items[currentLine-1])
			totalWeight += items[currentLine-1].Weight
			currentColumn -= items[currentLine-1].Weight
		}
		currentLine--
	}

	return knapsack, table[len(items)][knapSackCapacity], totalWeight
}

/*
带有自定义资源的一维0-1背包问题
*/
type Gpu struct {
	Model           string // GPU型号
	CardMemory      int    // 每张卡显存，单位GiB
	TotalCount      int    // 总卡数
	TotalMemory     int
	RemainingMemory int
	Cards           []GpuCard
}

type GpuCard struct {
	Index           int
	Uuid            string
	TotalMemory     int
	RemainingMemory int
}

func (gpu *Gpu) IsAllocatable(require int) bool {
	if require > gpu.CardMemory {
		return false
	}

	allocatable := false
	for _, c := range gpu.Cards {
		if c.RemainingMemory >= require {
			allocatable = true
			break
		}
	}

	return allocatable
}

func (gpu *Gpu) Allocate(require int) error {
	if !gpu.IsAllocatable(require) {
		return errors.New("can not allocate require gpu")
	}

	minRemain := gpu.CardMemory
	minIndex := -1

	for i, c := range gpu.Cards {
		if c.RemainingMemory < require {
			continue
		}

		if (c.RemainingMemory - require) < minRemain {
			minRemain = c.RemainingMemory - require
			minIndex = i
		}
	}
	if minIndex == -1 {
		return errors.New("can not allocate require gpu")
	}

	gpu.RemainingMemory -= require
	gpu.Cards[minIndex].RemainingMemory -= require

	return nil
}

func (gpu *Gpu) SetTotalMemory(totalMemory int) {
	if gpu.CardMemory == 0 {
		return
	}
	gpu.TotalMemory = totalMemory
	gpu.RemainingMemory = totalMemory

	if totalMemory <= gpu.CardMemory {
		gpu.Cards = []GpuCard{
			{
				Index:           0,
				TotalMemory:     totalMemory,
				RemainingMemory: totalMemory,
			},
		}
		gpu.TotalCount = 1
		return
	}

	count := totalMemory / gpu.CardMemory
	var cards []GpuCard
	for c := 0; c < count; c++ {
		cards = append(cards, GpuCard{
			Index:           c,
			TotalMemory:     gpu.CardMemory,
			RemainingMemory: gpu.CardMemory,
		})
	}
	if totalMemory%gpu.CardMemory > 0 {
		count += 1
		cards = append(cards, GpuCard{
			Index:           count - 1,
			TotalMemory:     totalMemory % gpu.CardMemory,
			RemainingMemory: totalMemory % gpu.CardMemory,
		})
	}

	gpu.TotalCount = count
	gpu.Cards = cards
	return
}

type CRItem struct {
	Index int
	Gpu   Gpu
	Value int
}

type CRKnapsack struct {
	Index      int
	Gpu        Gpu
	Items      []CRItem
	TotalValue int
}

func TestCRKSolve(t *testing.T) {
	items := []CRItem{
		{
			Index: 1,
			Gpu: Gpu{
				TotalMemory: 3,
				TotalCount:  1,
			},
			Value: 100,
		},
		{
			Index: 2,
			Gpu: Gpu{
				TotalMemory: 4,
				TotalCount:  1,
			},
			Value: 200,
		},
		{
			Index: 3,
			Gpu: Gpu{
				TotalMemory: 3,
				TotalCount:  1,
			},
			Value: 300,
		},
	}

	knapsack := CRKnapsack{
		Index: 1,
		Gpu: Gpu{
			Model:           "T4",
			CardMemory:      6,
			TotalCount:      2,
			TotalMemory:     12,
			RemainingMemory: 12,
			Cards: []GpuCard{
				{
					Index:           0,
					TotalMemory:     6,
					RemainingMemory: 6,
				},
				{
					Index:           1,
					TotalMemory:     6,
					RemainingMemory: 6,
				},
			},
		},
	}

	k := dynamicCRKnapsackSolve(items, knapsack)
	fmt.Printf("%v\n", k)
}

func dynamicCRKnapsackSolve(items []CRItem, knapsack CRKnapsack) CRKnapsack {
	ableToPack := func(item CRItem, knapsack CRKnapsack) bool {
		return knapsack.Gpu.IsAllocatable(item.Gpu.TotalMemory)
	}

	getLastLineGpuMem := func(gpu *Gpu, mem int) (*Gpu, int) {
		_ = gpu.Allocate(mem)
		maxRemain := 0
		for _, c := range gpu.Cards {
			if c.RemainingMemory > maxRemain {
				maxRemain = c.RemainingMemory
			}
		}

		if maxRemain == gpu.CardMemory {
			maxRemain = gpu.RemainingMemory
		}

		return gpu, maxRemain
	}

	table := make([][]int, len(items)+1)

	for i := 0; i < len(items)+1; i++ {
		table[i] = make([]int, knapsack.Gpu.CardMemory*knapsack.Gpu.TotalCount+1)

		for g := 0; g < knapsack.Gpu.CardMemory*knapsack.Gpu.TotalCount+1; g++ {
			if i == 0 {
				continue
			}

			k := knapsack
			k.Gpu.SetTotalMemory(g)

			if ableToPack(items[i-1], k) {
				_, lastGpuMemory := getLastLineGpuMem(&k.Gpu, items[i-1].Gpu.TotalMemory)
				aux := table[i-1][lastGpuMemory] + items[i-1].Value
				if aux > table[i-1][g] {
					table[i][g] = aux
				} else {
					table[i][g] = table[i-1][g]
				}
			} else {
				table[i][g] = table[i-1][g]
			}

		}
	}

	var packedItems []CRItem
	currentLine := len(items)
	currentGpu := &knapsack.Gpu
	currentColumn := knapsack.Gpu.TotalMemory

	for currentLine > 0 && currentColumn > 0 {
		if table[currentLine][currentColumn] != table[currentLine-1][currentColumn] {
			packedItems = append(packedItems, items[currentLine-1])
			currentGpu, currentColumn = getLastLineGpuMem(currentGpu, items[currentLine-1].Gpu.TotalMemory)
		}
		currentLine--
	}

	knapsack.Items = packedItems

	return knapsack
}

func generateCRItems(numberOfItems, maxGpuMemory, maxValue int) []CRItem {
	var items []CRItem

	rand.Seed(1)
	for i := 0; i < numberOfItems; i++ {
		items = append(items, CRItem{
			Index: i,
			Gpu: Gpu{
				Model:       "T4",
				TotalMemory: rand.Intn(maxGpuMemory) + 1,
				TotalCount:  1,
			},
			Value: 100,
		})
	}
	return items
}

func generateCRKnapsack(cardMemory, gpuCount int) CRKnapsack {
	gpu := Gpu{
		Model:       "T4",
		TotalCount:  gpuCount,
		CardMemory:  cardMemory,
		TotalMemory: cardMemory * gpuCount,
	}

	var cards []GpuCard
	for i := 0; i < gpuCount; i++ {
		cards = append(cards, GpuCard{
			Index:           i,
			TotalMemory:     cardMemory,
			RemainingMemory: cardMemory,
		})
	}
	gpu.Cards = cards

	return CRKnapsack{
		Index: 0,
		Gpu:   gpu,
	}
}

/*
动态规划法求解二维0-1背包问题
*/

// 二维物体，其中资源维度包含Length和Width
type TwoDimItem struct {
	Index  int
	Length int
	Width  int
	Value  int
}

type Knapsack struct {
	Index           int
	Length          int
	Width           int
	Items           []TwoDimItem
	RemainingLength int
	RemainingWidth  int
	TotalVale       int
}

type Knapsacks []Knapsack

func TestDynamicTwoDimsKnapsackSolve(t *testing.T) {
	items := generateTwoDimItems(100, 20, 100, 10)
	knapsack := Knapsack{
		Length: 72,
		Width:  512,
	}
	knapsack.RemainingLength = knapsack.Length
	knapsack.RemainingWidth = knapsack.Width

	startTime := time.Now()
	packedKnapsack := dynamicTwoDimsKnapsackSolve(items, knapsack)
	elapsedTime := time.Since(startTime).Seconds()

	fmt.Printf("Total time: %.2f\n", elapsedTime)
	utils.PrintMemUsage()

	fmt.Printf("Packed items: %v\n", packedKnapsack.Items)
	fmt.Printf("Total value: %v\n", packedKnapsack.TotalVale)
	fmt.Printf("Remaining length: %v, remaining width: %v\n", packedKnapsack.RemainingLength, packedKnapsack.RemainingWidth)
}

func TestMultiKnapsackSolve(t *testing.T) {
	items := generateTwoDimItems(100, 20, 100, 10)
	knapsacks := generateMultyKnapsack(100, 72, 512)

	var packedKnapsacks Knapsacks
	for i := 0; i < len(knapsacks); i++ {
		packedKnapsack := dynamicTwoDimsKnapsackSolve(items, knapsacks[i])
		packedKnapsacks = append(packedKnapsacks, *packedKnapsack)
	}

	// 挑选本轮资源利用率最高的机器，如果所有维度资源都剩0，则直接返回
	bestKnapsack, err := packedKnapsacks.GetBestKnapsack()
	if err != nil {
		t.Fatal(err)
	}
	if bestKnapsack == nil {
		t.Fatal("can not get best packed knapsack")
	}

	// TODO: 去除本轮已经挑选的items和knapsack，然后重新计算下一轮最优解

	fmt.Printf("Best packed knapsack: %v, packed items: %v\n", bestKnapsack.Index, bestKnapsack.Items)
	fmt.Printf("Remaining: %v, %v\n", bestKnapsack.RemainingLength, bestKnapsack.RemainingWidth)
}

func (ks Knapsacks) GetBestKnapsack() (*Knapsack, error) {
	if len(ks) == 0 {
		return nil, errors.New("empty knapsacks")
	}
	if len(ks) == 1 {
		return &ks[0], nil
	}

	bestKnapsack := ks[0]
	for i := 0; i < len(ks); i++ {
		if ks[i].RemainingLength == 0 && ks[i].RemainingWidth == 0 {
			return &ks[i], nil
		}

		if ks[i].RemainingLength < bestKnapsack.RemainingLength && ks[i].RemainingWidth < bestKnapsack.RemainingWidth {
			bestKnapsack = ks[i]
		}
	}

	return &bestKnapsack, nil
}

func generateTwoDimItems(numberOfItems, maxLength, maxWidth, maxValue int) []TwoDimItem {
	var items []TwoDimItem
	//rand.Seed(time.Now().UnixNano())
	rand.Seed(1)
	for i := 0; i < numberOfItems; i++ {
		items = append(items, TwoDimItem{
			Index:  i,
			Width:  rand.Intn(maxWidth) + 1,
			Length: rand.Intn(maxLength) + 1,
			Value:  rand.Intn(maxValue),
		})
	}

	return items
}

func generateMultyKnapsack(numberOfKnapsack, length, wideth int) []Knapsack {
	var knapsacks []Knapsack

	rand.Seed(1)
	for i := 0; i < numberOfKnapsack; i++ {
		knapsack := Knapsack{
			Index:  i,
			Length: rand.Intn(length) + 32,
			Width:  rand.Intn(wideth) + 128,
		}
		knapsack.RemainingLength = knapsack.Length
		knapsack.RemainingWidth = knapsack.Width

		knapsacks = append(knapsacks, knapsack)
	}

	return knapsacks
}

func dynamicTwoDimsKnapsackSolve(items []TwoDimItem, knapsack Knapsack) *Knapsack {
	var packedItems []TwoDimItem

	state := make([][][]int, len(items)+1) // 存储中间状态

	ableToPack := func(item TwoDimItem, length, width int) bool {
		if item.Width > width || item.Length > length {
			return false
		}
		return true
	}

	for i := 0; i < len(items)+1; i++ {
		stateLenDim := make([][]int, knapsack.Length+1)
		state[i] = stateLenDim

		for l := 0; l < knapsack.Length+1; l++ {
			stateWidDim := make([]int, knapsack.Width+1)
			state[i][l] = stateWidDim

			for w := 0; w < knapsack.Width+1; w++ {
				if i == 0 || l == 0 || w == 0 {
					continue
				}

				if ableToPack(items[i-1], l, w) {
					aux := items[i-1].Value + state[i-1][l-items[i-1].Length][w-items[i-1].Width]
					if aux > state[i-1][l][w] {
						state[i][l][w] = aux
					} else {
						state[i][l][w] = state[i-1][l][w]
					}
				} else {
					state[i][l][w] = state[i-1][l][w]
				}
			}
		}
	}

	// 获取最优结果项
	currentLine := len(items)

	for currentLine > 0 && knapsack.RemainingLength > 0 && knapsack.RemainingWidth > 0 {
		if state[currentLine][knapsack.RemainingLength][knapsack.RemainingWidth] != state[currentLine-1][knapsack.RemainingLength][knapsack.RemainingWidth] {
			packedItems = append(packedItems, items[currentLine-1])
			knapsack.RemainingLength = knapsack.RemainingLength - items[currentLine-1].Length
			knapsack.RemainingWidth = knapsack.RemainingWidth - items[currentLine-1].Width
			knapsack.TotalVale += items[currentLine-1].Value
		}
		currentLine--
	}
	knapsack.Items = packedItems

	return &knapsack
}
