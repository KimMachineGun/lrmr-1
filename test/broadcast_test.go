package test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/therne/lrmr/test/integration"
	"github.com/therne/lrmr/test/testutils"
)

func TestBroadcast(t *testing.T) {
	Convey("Given running nodes", t, integration.WithLocalCluster(2, func(cluster *integration.LocalCluster) {
		Convey("When running Map with broadcasts", func() {
			ds := BroadcastTester(cluster.Session)

			Convey("It should run without preserving broadcast values from master", func() {
				rows, err := ds.Collect()
				So(err, ShouldBeNil)
				So(rows, ShouldHaveLength, 1)
				So(testutils.StringValue(rows[0]), ShouldEqual, "throughStruct=foo, throughContext=bar")
			})
		})
	}))
}
