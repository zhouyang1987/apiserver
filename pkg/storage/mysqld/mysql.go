package mysqld

import (
	"io"

	"apiserver/pkg/configz"
	"apiserver/pkg/util/log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

// --------------
// Engine

type Engine struct {
	*xorm.Engine
}

var (
	engine *Engine
)

func init() {
	log.Info(configz.GetString("mysql", "dsn", ""))
	eng, err := xorm.NewEngine(configz.GetString("mysql", "dirver", "mysql"), configz.GetString("mysql", "dsn", ""))
	if err != nil {
		log.Fatalf("init mysql connection err: %v", err)
	}
	if err = eng.Ping(); err != nil {
		log.Fatalf("access the mysql db fail ,the reason is %s", err.Error())
	}
	engine = &Engine{Engine: eng}
	// cache
	// cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 1000)
	// engine.SetDefaultCacher(cacher)
}

func GetEngine() *Engine {
	return engine
}

func (engine *Engine) Debug() {
	engine.ShowSQL(true)
}

func (engine *Engine) Close() error {
	return engine.Close()
}

type Closer interface {
	io.Closer
}

func Close(db Closer) {
	if db != nil {
		if err := db.Close(); err != nil {
			log.Warning(err)
		}
	}
}
