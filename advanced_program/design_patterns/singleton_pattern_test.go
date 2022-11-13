package design_patterns

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

/*
单例模式：
	单例是一种创建型设计模式，能够保证一个类只有一个实例，并提供一个访问该实例的全局节点。
*/

var once sync.Once // 一种安全的初始化机制，能够保证并发的情况下只执行一次

type single struct { // 实体对象结构
}

var singleInstance *single // 声明一个全局的实体对象

func getInstance() *single {
	fmt.Println("----------", singleInstance)
	if singleInstance == nil {
		once.Do(func() {
			fmt.Println("Creating single instance now")
			singleInstance = &single{}
		})
	} else {
		fmt.Println("single instance already exists, reuse it")
	}
	return singleInstance
}

func TestSingle(t *testing.T) {
	for i := 0; i < 20; i++ {
		go getInstance() // 模拟并发请求去获取单例的实例对象
	}
	time.Sleep(5 * time.Second)
}
