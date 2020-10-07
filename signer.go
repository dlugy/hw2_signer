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
	innerMd5 := DataSignerMd5(data)
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
	}
	sort.Slice(innerData, func(i, j int) bool { return innerData[i].num < innerData[j].num })
	crc32 = innerData[0].data + "~" + innerData[1].data
	out <- crc32
}

func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	quotaMd5Chan := make(chan interface{}, 1)
	for val := range in {
		data := fmt.Sprintf("%v", val)
		wg.Add(1)
		go doSingleHash(out, data, quotaMd5Chan, wg)
	}
	wg.Wait()
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
		data := fmt.Sprintf("%v", val)
		wg.Add(1)
		go innerMultiHash(out, data, wg)
	}
	wg.Wait()
}

func CombineResults(in, out chan interface{}) {
/*	data := ""
	for val := range in {
		data = data + "_" + fmt.Sprintf("%v", val)
	}
	fmt.Println(data)*/


	data := ""
	innerData := make([]string, 8)
	j := 0
	for val := range in {
		innerData[j] = fmt.Sprintf("%v", val)
		j++
	}
	sort.Slice(innerData, func(i, j int) bool { return innerData[i] < innerData[j] })
	j = 0
	for _, s1 := range innerData {
		data = data + "_" + s1
	}
	data = data[1:]
	out <- data
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
	chanPipeLine[0] = make(chan interface{}, 100)
	for i := 0; i < len(hashSignJobs); i++ {
		chanPipeLine[i+1] = make(chan interface{}, 100)
		go func(i int) {
			hashSignJobs[i](chanPipeLine[i], chanPipeLine[i+1])
			close(chanPipeLine[i+1])
		}(i)
	}

	for range chanPipeLine[len(hashSignJobs)] {
	}
}

func ExecutePipeline(hashSignJobs []job) {

	chanPipeLine := make([]chan interface{}, len(hashSignJobs)+1)
	chanPipeLine[0] = make(chan interface{}, 100)
	for i := 0; i < len(hashSignJobs); i++ {
		chanPipeLine[i+1] = make(chan interface{}, 100)
		go func(i int) {
			hashSignJobs[i](chanPipeLine[i], chanPipeLine[i+1])
			close(chanPipeLine[i+1])
		}(i)
	}

	for range chanPipeLine[len(hashSignJobs)] {
	}
}