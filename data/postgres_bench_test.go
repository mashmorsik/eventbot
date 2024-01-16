package data_test

import (
	"eventbot/data"
	"testing"
)

func BenchmarkData_AddUser(b *testing.B) {
	db := data.NewData(MustConnectTest())
	for i := 0; i < b.N; i++ {
		err := db.AddUser(int64(46234629742390))
		if err != nil {
			return
		}
	}
}
