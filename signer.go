package main

import (
	"fmt"
	"sort"
	"strconv"
)


// сюда писать код

func SingleHash(in, out chan interface{}) {
	d1 := <- in
	data := fmt.Sprintf("%v", d1 )

	crc1 := DataSignerCrc32(data)
	crc2 := DataSignerCrc32(DataSignerMd5("0"))
	crc32 := crc1 + "~" + crc2

	out <- crc32
}

func MultiHash(in, out chan interface{}) {
	d1 := <-in
	data := fmt.Sprintf("%v", d1 )

	crcM := ""
	crcL := ""
	dataTmp := ""
	for i := 0; i <= 5; i++ {
		dataTmp = strconv.Itoa(i) + data
		crcL = DataSignerCrc32(dataTmp)
		crcM += crcL
	}
	out <- crcM
}

func CombineResults(in, out chan interface{}) {

}



func main() {
	var result []string

	ch1 := make(chan interface{})
	ch2 := make(chan interface{})

	go SingleHash(ch1, ch2)
	go MultiHash(ch2, ch2)
	go CombineResults(ch1, ch2)

	ch1 <- "0"

	result = append(result, "end")
	sort.Strings(result)

	fmt.Scanln()
}