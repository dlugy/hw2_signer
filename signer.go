package main

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
)

// сюда писать код

type crc32Channel struct {
	num int
	data string
}

func innerSingleHash(out chan interface{}, num int, data string, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println(num)

	var innerCrc32 crc32Channel
	crc := DataSignerCrc32(data)
	innerCrc32.num = num
	innerCrc32.data = crc

	out <- innerCrc32
}

func doSingleHash(out chan interface{}, data string, quotaChan chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	innerChan := make(chan interface{}, 2)
	innerWg := &sync.WaitGroup{}

	innerWg.Add(1)
	go innerSingleHash(innerChan, 0, data, innerWg)

	quotaChan <- struct{}{}
	innerMd5 := DataSignerMd5("0")
	<- quotaChan
	innerWg.Add(1)
	go innerSingleHash(innerChan, 1, innerMd5, innerWg)

	innerWg.Wait()
	close(innerChan)  // Закрываем канал чтобы след. цикл завершился

	crc32 := ""
	innerData := make([]crc32Channel, 2)
	j := 0
	for val := range innerChan {
		innerData[j] = val.(crc32Channel)
		j++
//		innerData = append(innerData, val.(crc32Channel))
	}
	sort.Slice(innerData, func(i, j int) bool { return innerData[i].num < innerData[j].num })
	crc32 = innerData[0].data + "`" + innerData[1].data
	out <- crc32
}

func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	quotaMd5Chan := make(chan interface{}, 1)
	for val := range in {
		data, _ := val.(string)
		wg.Add(1)
		go doSingleHash(out, data, quotaMd5Chan, wg)
	}
	wg.Wait()
	close(out)
}

func innerMultiHash(out chan interface{}, data string, wg *sync.WaitGroup) {
	defer wg.Done()

	crcM := ""
	dataTmp := ""
	innerChan := make(chan interface{}, 6)

	innerWg := &sync.WaitGroup{}
	for i := 0; i <= 5; i++ {
		dataTmp = strconv.Itoa(i) + data
		innerWg.Add(1)
		go innerSingleHash(innerChan, i, dataTmp, innerWg)
	}
	innerWg.Wait()

	close(innerChan)
	innerData := make([]crc32Channel, 6)
	j := 0
	for val := range innerChan {
		innerData[j] = val.(crc32Channel)
		j++
//		innerData = append(innerData, val.(crc32Channel))
	}
	sort.Slice(innerData, func(i, j int) bool { return innerData[i].num < innerData[j].num })
	for _, val := range innerData {
		crcM += val.data
	}
	out <- crcM

}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	for val := range in {
		data, _ := val.(string)
		wg.Add(1)
		go innerMultiHash(out, data, wg)
	}
	wg.Wait()
}

func CombineResults(in, out chan interface{}) {
/*	var innerData []string

	for val := range in {
		data, _ := val.(string)
		innerData = append(innerData, data)
	}
	sort.Slice(innerData, func(i, j int) bool { return innerData[i] < innerData[j] })

	res := innerData[0]
	for i := 1; i < len(innerData); i++ {
		res = res + "_" + innerData[i]
	}

 */
	res := "111"
//	fmt.Println("2222")
	out <- res
}

func main() {
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

	chanPipeLine := make([]chan interface{}, len(hashSignJobs)+1)

	for i := 0; i < len(hashSignJobs); i++ {
		if i == 0 {
			chanPipeLine[i] = make(chan interface{}, 100)
		}
		chanPipeLine[i+1] = make(chan interface{}, 100)
		go hashSignJobs[i](chanPipeLine[i], chanPipeLine[i+1])
	}

	fmt.Scanln()

}