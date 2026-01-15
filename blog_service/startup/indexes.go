package startup

import (
	blog "github.com/afteracademy/gomicro/blog-service/api/blog/model"
	"github.com/afteracademy/goserve/arch/mongo"
)

func EnsureDbIndexes(db mongo.Database) {
	go mongo.Document[blog.Blog](&blog.Blog{}).EnsureIndexes(db)
}
