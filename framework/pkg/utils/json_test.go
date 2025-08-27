package utils

import "testing"

func Test_ConvertToJSONString(t *testing.T) {
	jsonString, err := ConvertToJSONString("address:TGQm8VYkHkQnAPuNRwCiBUaJ8MLSkGfm3X,amount:15978100,blockHigh:66683747,coinType:TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t,decimals:6,fee:13844850,mainCoinType:195,memo:,status:3,tradeId:1302998592649158656,tradeType:1,txId:56f937a1446b6c7fea62433834d1ef05a9860141a478df63c2def51fe418c3fb}")
	if err != nil {
		t.Error(err)
	}
	t.Log(jsonString)
}
