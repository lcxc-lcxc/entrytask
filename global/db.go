package global

import (
	"github.com/bwmarrin/snowflake"
	"gorm.io/gorm"
)

var DBEngine *gorm.DB

var SnowFlakeNode1 *snowflake.Node
