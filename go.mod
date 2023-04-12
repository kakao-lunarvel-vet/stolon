module github.com/sorintlab/stolon

require (
	// github.com/coreos/bbolt v1.3.3 // indirect
	// github.com/coreos/etcd v3.3.18+incompatible // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/docker/leadership v0.1.0
	github.com/docker/libkv v0.2.1
	github.com/evanphx/json-patch v4.5.0+incompatible
	github.com/gofrs/uuid v4.2.0+incompatible
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/mock v1.4.0
	github.com/google/btree v1.0.1 // indirect
	github.com/google/go-cmp v0.5.6
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/hashicorp/consul/api v1.4.0
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/lib/pq v1.3.0
	github.com/mattn/go-colorable v0.1.11 // indirect
	github.com/mattn/go-isatty v0.0.14
	github.com/mitchellh/copystructure v1.0.0
	github.com/onsi/ginkgo v1.13.0 // indirect
	github.com/prometheus/client_golang v1.11.1
	github.com/rogpeppe/go-internal v1.8.1 // indirect
	github.com/sgotti/gexpect v0.0.0-20210315095146-1ec64e69809b
	github.com/sirupsen/logrus v1.7.0 // indirect
	github.com/soheilhy/cmux v0.1.5 // indirect
	github.com/sorintlab/pollon v0.0.0-20181009091703-248c68238c16
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/tmc/grpc-websocket-proxy v0.0.0-20201229170055-e5319fda7802 // indirect
	go.etcd.io/bbolt v1.3.7 // indirect
	go.etcd.io/etcd v3.3.27+incompatible
	go.etcd.io/etcd/client/v3 v3.5.7
	go.uber.org/zap v1.17.0
	golang.org/x/crypto v0.0.0-20220411220226-7b82a4e95df4 // indirect
	golang.org/x/net v0.4.0 // indirect
	golang.org/x/time v0.0.0-20210220033141-f8bda1e9f3ba // indirect
	google.golang.org/genproto v0.0.0-20210624195500-8bfb893ecb84 // indirect
	google.golang.org/grpc v1.54.0 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	k8s.io/api v0.17.3
	k8s.io/apimachinery v0.17.3
	k8s.io/client-go v0.17.3
)

go 1.12

replace (
	github.com/coreos/bbolt v1.3.3 => github.com/etcd-io/bbolt v1.3.3
	google.golang.org/grpc v1.54.0 => google.golang.org/grpc v1.29.1
)
