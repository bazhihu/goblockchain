package wallet

import (
	"bytes"
	"encoding/gob"
	"errors"
	"goblockchain/constcoe"
	"goblockchain/utils"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// 用户和钱包 是一对多的关系
// 钱包管理模块 管理一台机器上保存的所有钱包

// key值 为钱包地址
// value为钱包的别名
type RefList map[string]string

func (r *RefList) Save() {
	filename := constcoe.WalletsRefList + "ref_list.data"
	var content bytes.Buffer
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(r)
	utils.Handle(err)
	err = ioutil.WriteFile(filename, content.Bytes(), 0644)
	utils.Handle(err)
}

func (r *RefList) Update() {
	// 遍历目录 判断文件是否是钱包
	err := filepath.Walk(constcoe.Wallets, func(path string, info fs.FileInfo, err error) error {
		if info == nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		fileName := info.Name()
		// 比较文件后缀 是否是.wlt
		if strings.Compare(fileName[len(fileName)-4:], ".wlt") == 0 {
			if _, ok := (*r)[fileName[:len(fileName)-4]]; !ok {
				(*r)[fileName[:len(fileName)-4]] = ""
			}
		}
		return nil
	})
	utils.Handle(err)
}

func LoadRefList() *RefList {
	filename := constcoe.WalletsRefList + "ref_list.data"
	var reflist RefList
	if utils.FileExists(filename) {
		fileContent, err := ioutil.ReadFile(filename)
		utils.Handle(err)
		decoder := gob.NewDecoder(bytes.NewBuffer(fileContent))
		err = decoder.Decode(&reflist)
		utils.Handle(err)
	} else {
		reflist = make(RefList)
		reflist.Update()
	}
	return &reflist
}

// 绑定别名
func (r *RefList) BindRef(address, refname string) {
	(*r)[address] = refname
}

// 通过别名调取钱包地址
func (r *RefList) FindRef(refname string) (string, error) {
	temp := ""
	for key, val := range *r {
		if val == refname {
			temp = key
			break
		}
	}
	if temp == "" {
		err := errors.New("the refname is not found")
		return temp, err
	}
	return temp, nil
}
