package main

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
)

// сюда писать код

func DoSingleHash(out chan interface{}, data string, wg *sync.WaitGroup) {
	defer wg.Done()

	crc1 := DataSignerCrc32(data)
	crc2 := DataSignerCrc32(DataSignerMd5("0"))
	crc32 := crc1 + "~" + crc2

	out <- crc32
}

func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	for val := range in {
		data, _ := val.(string)
		wg.Add(1)
		go DoSingleHash(out, data, wg)
	}
	wg.Wait()
	close(out)
}

func MultiHash(in, out chan interface{}) {
	crcM := ""
	crcL := ""
	dataTmp := ""

	innerChan := make(chan interface{})

	for val := range in {
		for i := 0; i <= 5; i++ {
			data, _ := val.(string)
			dataTmp = strconv.Itoa(i) + data
			crcL = DataSignerCrc32(dataTmp)
			crcM += crcL
		}
		out <- crcM
	}
}

func CombineResults(in, out chan interface{}) {

}

func main() {
	var result []string
	inputData := []int{0,1}

	hashSignJobs := []job{
		job(func(in, out chan interface{}) {
			for _, fibNum := range inputData {
				out <- fibNum
			}
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
		job(func(in, out chan interface{}) {
			dataRaw := <-in
			data, ok := dataRaw.(string)
			fmt.Println(data, ok)
		}),
	}

	var chanIn []chan interface{}
	var chanOut []chan interface{}

	for i := 0; i < len(hashSignJobs); i++ {
		chanIn[i] = make(chan interface{})
		chanOut[i] = make(chan interface{})
		go hashSignJobs[i](chanIn[i], chanOut[i])
	}

	result = append(result, "end")
	sort.Strings(result)
	fmt.Scanln()
}