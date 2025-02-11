package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {

	

	start := time.Now()		// время запуска программы
	defer func() {		// сколько времени ушло на программу
		fmt.Println("Прошло времени", time.Since(start))	
	}()
	path_output := flag.String("src", "urls.txt", "Путь к файлу c url")		// путь к файлы с юрл
	path_input := flag.String("dst", "./result", "Путь куда сохранить html страницы")	// путь к директории где хранить

	flag.Parse()	

	file, err := os.Open(*path_output)		// проверка на сущетсвования файла с юрл

	if err != nil {																										
		fmt.Println("Неверный путь к файлу или файл не существует")	
		return
	}
	defer file.Close()

	if _, err := os.Stat(*path_input); os.IsNotExist(err) {			// проверка на существаование директории

		if err := os.MkdirAll(*path_input, os.ModePerm); err != nil {
			fmt.Println("Ошибка при создании директории", err)
			return
		}
		fmt.Println("Директория создана")
	}

	scanner := bufio.NewScanner(file)		// 

	wg := sync.WaitGroup{}

	for scanner.Scan() {
		wg.Add(1)
		domen := scanner.Text()
		go func () {
			defer wg.Done()
			url := "https://" + domen
			fmt.Println(url)
			if checkUrl(url) {
				req, err := http.Get(url)
				if err != nil {
					fmt.Println("Ошибка get запроса, проверьте правильность url", err)
				}
				defer req.Body.Close()
				b, err := io.ReadAll(req.Body)
				if err != nil{
					fmt.Println("Ошибка чтения запроса")
				}
				createHtml(*path_input+"/"+domen, b)
			} else {
				fmt.Println("Ошибка правильность url адреса")
			}
		}()
	}

	wg.Wait()

}

func checkUrl(s string) bool {		// функция валидации юрл

	array := strings.Split(s, ".")
	return len(array) == 2

}

func createHtml(path string, data []byte) {		// функция для создания html страниц

	file, err := os.Create(path + ".html")
	if err != nil {
		fmt.Println("Ошибка создания файла")
		return
	}
	defer file.Close()

	_, err = file.Write(data)

	if err != nil {
		fmt.Println("Ошибка при записи в файл")
		return
	}
	fmt.Println("Данные успешно записаны")

}
