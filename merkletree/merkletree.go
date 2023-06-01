package merkletree

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"goblockchain/transaction"
	"goblockchain/utils"
)

// MT 梅克尔树 哈希树
// 通过树形哈希，组成区块头哈希，总哈希来验证区块的完整性

// SPV 快速交易验证功能
// 只存储整个区块链头部信息的节点 是轻节点
// 参与挖矿并共识的节点称为  全节点

// 交易过程
// 1、找到自己节点上的 区块ID的头部信息 -  包含MT树根哈希值
// 2、根据接收的交易信息按照  MT验证路径重新计算MT 树根哈希值
// 对比两者的哈希值
// 要点：1、验证路径不大于1KB，通信消耗小； 2、不担心伪造全节点的路径，因为对方不知道节点其他的哈希值
// 但是还是有 双花的风险 - 消耗方：多次消费； 出售方未及时同步至区块链网络
// 技术难点：
// 1、买家提供的MT验证路径的交易信息有没有被使用过，X币是否存在重复消费的状况
// 2、如何保证卖家 上传交易信息至区块链过程中，再用交易信息取做其他买卖

type MerkleTree struct {
	RootNode *MerkleNode
}

type MerkleNode struct {
	LeftNode  *MerkleNode
	RightNode *MerkleNode
	Data      []byte
}

// 先判断是否是叶子节点，叶子节点存储具体数据，非叶子节点存储hash值
func CreateMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	tempNode := MerkleNode{}

	if left == nil && right == nil {
		tempNode.Data = data
	} else {
		catenateHash := append(left.Data, right.Data...)
		hash := sha256.Sum256(catenateHash)
		tempNode.Data = hash[:]
	}

	tempNode.LeftNode = left
	tempNode.RightNode = right

	return &tempNode
}

// 创建MT树，奇数个叶子节点，将最后一个叶子节点补到最前面，补齐偶数个节点
func CrateMerkleTree(txs []*transaction.Transaction) *MerkleTree {
	txslen := len(txs)
	if txslen%2 != 0 {
		txs = append(txs, txs[txslen-1])
	}

	var nodePool []*MerkleNode
	for _, tx := range txs {
		nodePool = append(nodePool, CreateMerkleNode(nil, nil, tx.ID))
	}

	for len(nodePool) > 1 {
		var tempNodePool []*MerkleNode
		poolLen := len(nodePool)
		if poolLen%2 != 0 {
			tempNodePool = append(tempNodePool, nodePool[poolLen-1])
		}
		for i := 0; i < poolLen/2; i++ {
			tempNodePool = append(tempNodePool, CreateMerkleNode(nodePool[2*i], nodePool[2*i+1], nil))
		}
		nodePool = tempNodePool
	}

	merkleTree := MerkleTree{nodePool[0]}
	return &merkleTree
}

// route 方向 0为左 1为右
// hashroute 保存哈希值
// 深度优先搜索算法
func (mn *MerkleNode) Find(data []byte, route []int, hashroute [][]byte) (bool, []int, [][]byte) {
	findFlag := false
	if bytes.Equal(mn.Data, data) {
		findFlag = true
		return findFlag, route, hashroute
	} else {
		if mn.LeftNode != nil {
			route_t := append(route, 0)
			hashroute_t := append(hashroute, mn.RightNode.Data)

			findFlag, route_t, hashroute_t = mn.LeftNode.Find(data, route_t, hashroute_t)
			if findFlag {
				return findFlag, route_t, hashroute_t
			} else {
				if mn.RightNode != nil {
					route_t = append(route_t, 1)
					hashroute_t = append(hashroute_t, mn.LeftNode.Data)
					findFlag, route_t, hashroute_t = mn.RightNode.Find(data, route_t, hashroute_t)
					if findFlag {
						return findFlag, route_t, hashroute_t
					} else {
						return findFlag, route, hashroute
					}
				}
			}
		} else {
			return findFlag, route, hashroute
		}
	}
	return findFlag, route, hashroute
}

// 通过交易信息的ID - 交易信息的哈希值
// 返回验证路径与一个是否找到该交易信息的信号
func (mt *MerkleTree) BackValidationRoute(txid []byte) ([]int, [][]byte, bool) {
	ok, route, hashroute := mt.RootNode.Find(txid, []int{}, [][]byte{})
	return route, hashroute, ok
}

// spv 函数
// 按照MT验证路径验证交易信息是否有效
func SimplePaymentValidation(txid, mtroothash []byte, rount []int, hashroute [][]byte) bool {
	routeLen := len(rount)
	var tempHash []byte = txid

	for i := routeLen - 1; i >= 0; i-- {
		if rount[i] == 0 {
			catenateHash := append(tempHash, hashroute[i]...)
			hash := sha256.Sum256(catenateHash)
			tempHash = hash[:]
		} else if rount[i] == 1 {
			catenateHash := append(hashroute[i], tempHash...)
			hash := sha256.Sum256(catenateHash)
			tempHash = hash[:]
		} else {
			utils.Handle(errors.New("error in validation route"))
		}
	}
	return bytes.Equal(tempHash, mtroothash)
}
