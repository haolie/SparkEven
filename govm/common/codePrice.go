package common

import "fmt"

type CodePrice struct {
	FaceId    int
	Time      int
	Price     int
	TradeType int
	Volume    int
}

func (this *CodePrice) Print() {
	fmt.Printf("FaceId:%d\n", this.FaceId)
	fmt.Printf("Time:%d\n", this.Time)
	fmt.Printf("Price:%d\n", this.Price)
	fmt.Printf("TradeType:%d\n", this.TradeType)
	fmt.Printf("Volume:%d\n", this.Volume)
}
