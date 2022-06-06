package service

import "sync"

type singleTon chan [2]int

var instance *singleTon
var once sync.Once

func GetInstance() *singleTon {
	once.Do(func() {
		ret := make(singleTon, 100)
		instance = &ret
	})
	return instance
}
