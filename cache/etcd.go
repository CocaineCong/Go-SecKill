package cache

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"gopkg.in/ini.v1"
	"time"
)

var ETCDClient *clientv3.Client

func ETCD() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{ETCDAddr},
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		fmt.Println("etcd", err)
	}
	ETCDClient = cli

}

func LoadETCDData(file *ini.File) {
	ETCDAddr = file.Section("etcd").Key("ETCDAddr").String()
}