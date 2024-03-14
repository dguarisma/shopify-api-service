package utils

import "testing"

func (ht *HandleTest) DeleteInfo(t *testing.T, tables []interface{}) {
	for _, table := range tables {
		if err := ht.Db.Unscoped().Where("1=1").Delete(table).Error; err != nil {
			t.Fatalf(err.Error())
		}
	}
}
