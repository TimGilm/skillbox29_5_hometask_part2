/*
Задание 2. Graceful shutdown
Научиться правильно останавливать приложения.
В работе часто возникает потребность правильно останавливать приложения. Например, когда наш сервер
обслуживает соединения, а нам хочется, чтобы все текущие соединения были обработаны и лишь потом
произошло выключение сервиса. Для этого существует паттерн graceful shutdown.
Напишите приложение, которое выводит квадраты натуральных чисел на экран, а после получения сигнала
^С обрабатывает этот сигнал, пишет «выхожу из программы» и выходит.
Советы и рекомендации
Для реализации данного паттерна воспользуйтесь каналами и оператором select с default-кейсом.
*/
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

func main() {
	var wg sync.WaitGroup

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	go func() {
		for {
			select {
			case <-sigs:
				fmt.Println("выхожу из программы.")
				os.Exit(0)
			default:
				en := enterNumber(&wg)
				squareNumber(&wg, en)
				wg.Wait()
			}
		}
	}()
	select {}
}

func enterNumber(wg *sync.WaitGroup) chan int {
	outChan := make(chan int)
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			close(outChan)
		}()
		fmt.Println("Введите цифру, либо CTRL+с для выхода: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		if value, err := strconv.Atoi(scanner.Text()); err == nil {
			fmt.Println("Ввод: ", value)
			outChan <- value
		} else {
			fmt.Println("Некорректный ввод, повторите попытку")
		}
	}()
	return outChan
}

func squareNumber(wg *sync.WaitGroup, inChan chan int) {
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
		}()
		for val := range inChan {
			fmt.Println("Квадрат: ", val*val)
		}
	}()
}
