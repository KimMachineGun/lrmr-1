package integration

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/therne/lrmr/coordinator"
	"github.com/thoas/go-funk"
)

const (
	etcdEndpointEnvKey  = "LRMR_TEST_ETCD_ENDPOINT"
	defaultEtcdEndpoint = "127.0.0.1:2379"
)

// ProvideEtcd provides coordinator.Etcd on integration tests.
// Otherwise, coordinator.LocalMemory is provided.
func ProvideEtcd() coordinator.Coordinator {
	if !IsIntegrationTest {
		return coordinator.NewLocalMemory()
	}
	rand.Seed(time.Now().Unix())
	testNs := fmt.Sprintf("lrmr_test_%s/", funk.RandomString(10))

	etcdEndpoint, ok := os.LookupEnv(etcdEndpointEnvKey)
	if !ok {
		etcdEndpoint = defaultEtcdEndpoint
	}
	etcd, err := coordinator.NewEtcd([]string{etcdEndpoint}, testNs)
	if err != nil {
		So(err, ShouldBeNil)
	}

	// clean all items under test namespace
	Reset(func() {
		time.Sleep(400 * time.Millisecond)

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		log.Verbose("Closing etcd")
		if _, err := etcd.Delete(ctx, ""); err != nil {
			So(err, ShouldBeNil)
		}
		if err := etcd.Close(); err != nil {
			So(err, ShouldBeNil)
		}
	})
	return etcd
}
