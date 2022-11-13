package logger

import (
	"context"
	"flag"
	"fmt"
	"go.uber.org/zap/zapcore"
	"k8s.io/apimachinery/pkg/util/uuid"
	"testing"
)

func TestLogger(t *testing.T) {
	o := &Options{
		Encoder: "json",
		Caller:  true,
	}

	Init(o, func(config *zapcore.EncoderConfig) {
		config.TimeKey = "customer_time_key"
	})

	G().AddGlobalFields("idc", "dx", "name", "dcos_api")

	test1()

	test2()
}

func test1() {
	G().AddTempFields("request_id", uuid.NewUUID())
	defer G().ResetTempFields()

	G().Info("test1")
}
func test3() {
	G().Info("test3")
}

func test2() {
	G().AddTempFields("request_id", uuid.NewUUID())
	defer G().ResetTempFields()

	G().Errorw("test2", "parameters", struct {
		Name string
		Age  int
	}{Name: "xhli18", Age: 18})
}

func TestLog(t *testing.T) {
	o := &Options{
		Encoder: "json",
		Caller:  false,
	}

	o.BindFlags(flag.NewFlagSet("log", flag.ExitOnError))

	flag.Parse()
	Init(o, func(config *zapcore.EncoderConfig) {
		config.TimeKey = "customer_time_key"
	})

	ctx := context.WithValue(context.Background(), "params", []interface{}{"name", "xhli18", "age", 18})
	test(ctx)
}

func test(ctx context.Context) {
	params, _ := ctx.Value("params").([]interface{})
	G().Infow(fmt.Sprintf("test1=%s", "aaa"), params...)
	G().Info("test222")
}
