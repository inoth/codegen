package handle

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/inoth/codegen/db"
	"github.com/inoth/codegen/templet"
	"github.com/inoth/codegen/util"
)

func CreateFolder(projectName string) {
	for k, _ := range templet.FileTemplate {
		path := fmt.Sprintf("./%v/%v", projectName, k)

		if util.PathExists(path) {
			continue
		}
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Fatal(err.Error())
			return
		}
	}
}

func HandlerTable(ctx context.Context, ch_progress chan string, ch_limit chan struct{}, table db.TableInfo) {
	wg := sync.WaitGroup{}
	for _, temp := range templet.FileTemplate {
		wg.Add(1)
		go temp.Process(&wg, table.TableName)
	}
	wg.Wait()
	ch_progress <- fmt.Sprintf("%v 处理完成", table.TableName)
	<-ch_limit
	return
}
