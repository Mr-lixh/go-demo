package design_patterns

import (
	"fmt"
	"testing"
)

/*
工厂模式：
（1）简单工厂：用来生产同一等级结构的任意产品，对于增加新的产品无能为力；
（2）工厂方法：用来生产同一等级结构中的固定产品，支持增加任意产品；
（3）抽象工厂：用来生产不同产品族的全部产品，对于增加新的产品无能为力，支持增加产品族。
*/

// 工厂方法
type iGun interface { // 设计一个通用的接口类型
	setName(name string)
	setPower(power int)
	getName() string
	getPower() int
}

type gun struct { // 设计枪对象，然后分别实现通用接口
	name  string
	power int
}

func (g *gun) setName(name string) {
	g.name = name
}

func (g *gun) getName() string {
	return g.name
}

func (g *gun) setPower(power int) {
	g.power = power
}

func (g *gun) getPower() int {
	return g.power
}

type ak47 struct { // 定义指定产品ak47枪
	gun
}

func newAk47() iGun { // 创建ak47产品
	return &ak47{
		gun: gun{
			name:  "AK47 gun",
			power: 4,
		},
	}
}

func getGun(gunType string) (iGun, error) { // 获取指定类型的产品，为工厂方法封装体
	switch gunType {
	case "ak47":
		return newAk47(), nil
	default:
		return nil, fmt.Errorf("no support type")
	}
}

func TestFactory(t *testing.T) {
	ak47, _ := getGun("ak47")
	fmt.Println(ak47.getPower())
}
