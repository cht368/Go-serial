package main

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/tarm/serial"

	"fmt"
)

func Counter(t *testing.T) {
	port0 := "COM2"
	port1 := "COM1"
	if port0 == "" || port1 == "" {
		t.Skip("Skipping test because PORT0 or PORT1 environment variable is not set")
	}
	c0 := &serial.Config{Name: port0, Baud: 2400}
	c1 := &serial.Config{Name: port1, Baud: 2400}

	s1, err := serial.OpenPort(c0)
	if err != nil {
		fmt.Println(err)
	}

	s2, err := serial.OpenPort(c1)
	if err != nil {
		fmt.Println(err)
	}

	ch := make(chan int, 1)
	go func() {
		buf := make([]byte, 128)
		var readCount int
		var temp string

		maxCounter := 1
		for {
			n, err := s2.Read(buf)
			if err != nil {
				fmt.Print(err)
			}
			readCount++
			// fmt.Printf("Read %v %v bytes: % 02x %s", readCount, n, buf[:n], buf[:n])
			weight := fmt.Sprintf("%s ", buf[:n])
			if strings.ContainsAny(weight, "+") {
				weight = strings.Trim(weight, "+")
				// fmt.Print(weight + " ")
				EndTime = time.Now().Format(TimeFormat)
				temp += EndTime + ",		" + weight + "\n"
				WriteTxt(temp)

				intWeight, err := strconv.Atoi(strings.Replace(weight, " ", "", -1)) //convert weight to int and remove white space
				if err != nil {
					fmt.Println("hell no, ", err)
				}
				maxCounter = len(AllTempMax)
				fmt.Println(weight, ": weight in int: ", intWeight, "counter: ", maxCounter)
				if intWeight >= MAX {
					tempMax := &ExcelTable{
						No:    strconv.Itoa(maxCounter),
						Jam:   EndTime,
						Max:   strconv.Itoa(intWeight),
						Lama:  "test",
						Awal:  "test",
						Akhir: "test",
					}
					TempMaxs = append(TempMaxs, *tempMax)
				} else {
					if TempMaxs != nil {
						var max int
						var awal time.Time
						var akhir time.Time
						var timeDef string

						for i, tempMax := range TempMaxs {
							time, _ := time.Parse(TimeFormat, tempMax.Jam)
							tMax, _ := strconv.Atoi(tempMax.Max)
							if tMax > max {
								max = tMax
							}
							if i == 0 {
								awal = time
							}
							if i == (len(TempMaxs) - 1) {
								akhir = time
							}
							ttd := awal.Sub(akhir)
							timeDef = fmt.Sprintf("%.0f:%.0f:%.0f", ttd.Hours()*-1, ttd.Minutes()*-1, ttd.Seconds()*-1)
						}
						passTemp := &ExcelTable{
							No:    strconv.Itoa(maxCounter),
							Jam:   TempMaxs[0].Jam,
							Max:   strconv.Itoa(max),
							Lama:  timeDef,
							Awal:  awal.Format("15:04:05"),
							Akhir: akhir.Format("15:04:05"),
						}
						AllTempMax = append(AllTempMax, *passTemp)
						TempMaxs = nil
					}
				}
			}

			select {
			case <-ch:
				ch <- readCount
				close(ch)
			default:
			}
		}
	}()

	if _, err = s1.Write([]byte(" ")); err != nil {
		fmt.Println(err)
	}
	if _, err = s1.Write([]byte(" ")); err != nil {
		fmt.Println(err)
	}
	time.Sleep(time.Millisecond / 1)
	if _, err = s1.Write([]byte(" ")); err != nil {
		fmt.Println(err)
	}
	time.Sleep(time.Millisecond / 1)

	ch <- 0
	s1.Write([]byte(" ")) // We could be blocked in the read without this
	c := <-ch
	exp := 5
	if c >= exp {
		fmt.Println("Expected less than %v read, got %v", exp, c)
	}
}

func byteSlice(arr []byte) byte {
	sum := byte(0)
	for _, b := range arr {
		sum += b
	}
	return sum
}
