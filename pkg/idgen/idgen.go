// Package idgen 提供分布式 ID 生成（Snowflake）。
package idgen

import (
	"fmt"

	"github.com/bwmarrin/snowflake"
)

// Generator ID 生成接口，便于测试替换。
type Generator interface {
	NextID() int64
}

type snowflakeGenerator struct {
	node *snowflake.Node
}

// NewSnowflake 创建雪花 ID 生成器，nodeID 取值 0-1023，多实例部署必须唯一。
func NewSnowflake(nodeID int64) (Generator, error) {
	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		return nil, fmt.Errorf("new snowflake node: %w", err)
	}
	return &snowflakeGenerator{node: node}, nil
}

func (g *snowflakeGenerator) NextID() int64 { return g.node.Generate().Int64() }
