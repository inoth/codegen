package codegen

import (
	"codegenfornet/components/cache"
	"codegenfornet/register"
	"testing"
	"time"

	_ "codegenfornet/src/plugins/struct_crawl/all"
)

func TestCodegen(t *testing.T) {
	t1 := time.Now()

	reg := register.Register(
		&cache.CacheComponents{},
	)

	err := reg.Init().Run(&CodegenForNet{
		DbType:      "",
		ProjectName: "projectName",
		TempPath:    "temp/temp.zip",
		SourceDb:    false,
	})
	if err != nil {
		t.Logf(err.Error())
	}

	t.Logf("ok; time: %v", time.Since(t1))
}
