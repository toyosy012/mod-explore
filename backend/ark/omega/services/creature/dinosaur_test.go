package creature

import (
	"testing"
)

func Test_InitHealth(t *testing.T) {
	t.Run("体力の生成テスト", func(t *testing.T) {
		health, err := NewHealth(1)
		if err != nil {
			t.Errorf("正常な体力の生成ができませんでした %d", health)
		}
	})

	t.Run("異常な体力の生成テスト", func(t *testing.T) {
		_, err := NewHealth(0)
		if err != nil {
			t.Logf("異常な体力の生成されエラーが発生しました %s", err.Error())
		}
	})
}
