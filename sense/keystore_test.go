package sense

import (
	"testing"
)

func TestKeyStore(t *testing.T) {
	tableName := "test"
	ks := NewDynamoDBKeyStore(tableName, nil)

	key, err := ks.Get("")
	if err == nil {
		t.Error("Err should not be nil")
		t.FailNow()
	}

	if len(key) != 0 {
		t.Errorf("Key should be 0 if error, key was: %d\n", len(key))
		t.FailNow()
	}
}
