module github.com/weedge/wedis

go 1.19

require (
	github.com/cloudwego/hertz v0.6.2
	github.com/cloudwego/kitex v0.5.2
	github.com/google/wire v0.5.0
	github.com/hertz-contrib/obs-opentelemetry/provider v0.2.1
	github.com/hertz-contrib/obs-opentelemetry/tracing v0.2.1
	github.com/spf13/cobra v1.7.0
	github.com/tidwall/redcon v1.6.2
	github.com/weedge/openkv-goleveldb v0.0.0-20230527091727-961892f63755
	github.com/weedge/pkg v0.0.0-20230616035038-f6c50b5a2707
	github.com/weedge/xdis-storager v0.0.0-20230616035224-eddcc703b8bb
)

require (
	github.com/andeya/goutil v1.0.0 // indirect
	github.com/apache/thrift v0.13.0 // indirect
	github.com/bytedance/go-tagexpr/v2 v2.9.6 // indirect
	github.com/bytedance/gopkg v0.0.0-20221122125632-68358b8ecec6 // indirect
	github.com/bytedance/sonic v1.8.1 // indirect
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
	github.com/choleraehyq/pid v0.0.16 // indirect
	github.com/cloudwego/netpoll v0.3.2 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gofrs/flock v0.8.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.15.2 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/henrylee2cn/ameda v1.5.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/kitex-contrib/obs-opentelemetry/logging/zap v0.0.0-20230512024524-5f5a227105f7 // indirect
	github.com/klauspost/cpuid/v2 v2.2.4 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/nyaruka/phonenumbers v1.1.6 // indirect
	github.com/pelletier/go-toml/v2 v2.0.6 // indirect
	github.com/spf13/afero v1.9.3 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.15.0 // indirect
	github.com/subosito/gotenv v1.4.2 // indirect
	github.com/syndtr/goleveldb v1.0.0 // indirect
	github.com/tidwall/btree v1.6.0 // indirect
	github.com/tidwall/gjson v1.14.4 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	go.opentelemetry.io/contrib/instrumentation/runtime v0.41.0 // indirect
	go.opentelemetry.io/contrib/propagators/b3 v1.16.0 // indirect
	go.opentelemetry.io/contrib/propagators/ot v1.16.0 // indirect
	go.opentelemetry.io/otel v1.15.1 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.15.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric v0.38.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v0.38.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.15.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.15.0 // indirect
	go.opentelemetry.io/otel/metric v0.38.0 // indirect
	go.opentelemetry.io/otel/sdk v1.15.1 // indirect
	go.opentelemetry.io/otel/sdk/metric v0.38.0 // indirect
	go.opentelemetry.io/otel/trace v1.15.1 // indirect
	go.opentelemetry.io/proto/otlp v0.19.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.24.0 // indirect
	golang.org/x/arch v0.2.0 // indirect
	golang.org/x/net v0.9.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1 // indirect
	google.golang.org/grpc v1.54.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/tidwall/redcon => github.com/weedge/redcon v0.0.0-20230613050749-77c37f99e042

//replace github.com/tidwall/redcon => ../redcon

//replace github.com/weedge/pkg => ../pkg

//replace github.com/weedge/xdis-storager => ../xdis-storager

//replace github.com/weedge/openkv-goleveldb => ../openkv-goleveldb
