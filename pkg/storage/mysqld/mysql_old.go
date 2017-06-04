package mysqld

/*// Copyright © 2017 huang jia <449264675@qq.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mysqld

import (
	"io"

	// "apiserver/pkg/configz"
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
	// eng, err := xorm.NewEngine(configz.GetString("mysql", "dirver", "mysql"), configz.GetString("mysql", "dsn", ""))
	// if err != nil {
	// 	log.Fatalf("init mysql connection err: %v", err)
	// }
	// if err = eng.Ping(); err != nil {
	// 	log.Fatalf("access the mysql db fail ,the reason is %s", err.Error())
	// }
	// eng.ShowSQL(configz.MustBool("system", "debug", false))
	// engine = &Engine{Engine: eng}
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
*/
