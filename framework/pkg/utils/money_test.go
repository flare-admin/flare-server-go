package utils

import (
	"fmt"
	"testing"
)

func TestFlPointsToYuan(t *testing.T) {
	t.Log(FlPointsToYuan(121.5617))
}
func Test_CalculateActualAmount(t *testing.T) {
	amount := "500298299999999950000" // 示例金额
	decimals := 18
	actualAmount, err := CalculateActualAmount(amount, decimals)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("实际金额: %s\n", actualAmount.FloatString(4))
}

func Test_SplitAmount(t *testing.T) {
	amount := "500298299999999950000" // 示例金额
	decimals := 18
	splitAmount, i, err := SplitAmount(amount, decimals, 4)
	if err != nil {
		fmt.Println("Error:", err)
	}
	t.Logf("Split amount: %d, decimals: %d", splitAmount, i)
}

func Test_PointsToIntegerYuan(t *testing.T) {
	yuan := PointsToIntegerYuan(123453)
	t.Log(IntegerToDBMoney(yuan))
}

func Test_ConvertAmountToCents(t *testing.T) {
	cents, err := ConvertAmountToCents("400413900000000000000", 18)
	if err != nil {
		t.Error(err)
	}
	t.Log(cents)
}
