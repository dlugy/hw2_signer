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

	crcL := ""
	dataTmp := ""
	for i := 0; i <= 5; i++ {
		dataTmp = strconv.Itoa(i) + data

	}

}

func main() {
	var result []string

	ch1 := make(chan interface{})
	ch2 := make(chan interface{})

	go SingleHash(ch1, ch2)

	ch1 <- "0"
	res := fmt.Sprintf("%v", <- ch2)

/*
	crc1 := DataSignerCrc32("0")
	crc2 := DataSignerCrc32(DataSignerMd5("0"))
	crc32 := crc1 + "~" + crc2
	fmt.Println(crc32)
	fmt.Println("-------------")

	crcM := ""
	crcL := ""
	dataTmp := ""
	for i := 0; i <=5; i++  {
		dataTmp = strconv.Itoa(i) + crc32
		crcL = DataSignerCrc32(dataTmp)
		crcM += crcL
		fmt.Println(dataTmp, crcL)
	}

	result = append(result, crcM)
	result = append(result, "a")
	result = append(result, "c")
	result = append(result, "b")
*/

	result = append(result, res)
	result = append(result, "end")
	sort.Strings(result)
	fmt.Println(result)

	fmt.Scanln()
}