package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

func main() {

	start := time.Now()		
	defer func() {		
		log.Println("Прошло времени", time.Since(start))	
	}()
	path_output := flag.String("src", "", "Путь к файлу с URL-адресами. Создайте файл и укажите полный путь к файлу, содержащему список URL, которые нужно обработать.")
	path_input := flag.String("dst", "", "Путь к директории для сохранения HTML-страниц. Укажите полный путь к директории, где будут сохранены загруженные HTML-страницы.")

	// Кастомная функция для вывода справки
	flag.Usage = func() {
		fmt.Println("Использование программы:")
		fmt.Println("  go run main.go -src <путь к файлу urls.txt> -dst <путь к директории для сохранения HTML>")
		fmt.Println("Параметры:")
		flag.PrintDefaults()
	}

	flag.Parse()

	// Проверка обязательных флагов
	if *path_output == "" || *path_input == "" {
		fmt.Println("Ошибка: оба флага обязательны: -src и -dst.")
		flag.Usage()
		os.Exit(1)
	}

	file, err := os.Open(*path_output)	
	if err != nil {																										
		log.Println("Неверный путь к файлу или файл не существует", err)	
		return
	}
	defer file.Close()

	if _, err := os.Stat(*path_input); os.IsNotExist(err) {			

		if err := os.MkdirAll(*path_input, os.ModePerm); err != nil {
			log.Println("Ошибка при создании директории", err)
			return
		}
		log.Println("Директория создана")
	}

	readFile(file, path_input)
	
}

// readFile считываем данные из файла, делаем get запросы
func readFile(file *os.File, path_input *string)  {

	scanner := bufio.NewScanner(file)

	wg := sync.WaitGroup{}

	for scanner.Scan() {
		wg.Add(1)
		domen := scanner.Text()
		go func (d string) {
			defer wg.Done()
			url := d
			if !strings.Contains(d, "http://") &&  !strings.Contains(d, "https://"){
				url = "http://" + d
			}
			parts := strings.Split(d, "//")
			d = parts[len(parts)-1]
			fmt.Println(url)
			if checkUrl(url) {
				req, err := http.Get(url)
				if err != nil {
					log.Println("Ошибка get запроса, проверьте корректность сайта", err)
					return
				}
				defer req.Body.Close()
				b, err := io.ReadAll(req.Body)
				if err != nil{
					log.Println("Ошибка чтения запроса", err)
				}
				createHtml(*path_input+"/"+d, b)
			} else {
				log.Println("Ошибка правильность url адреса", url)
			}
		}(domen)
	}

	wg.Wait()

}

// checkUrl функция валидации юрл
func checkUrl(s string) bool {		

	// Регулярное выражение для проверки корректности URL
	re := regexp.MustCompile(`^(http|https)://[a-zA-Z0-9.-]+(:[0-9]+)?(/.*)?$`)
	return re.MatchString(s)
}

// createHtml функция для создания html страниц
func createHtml(path string, data []byte) error{		

	file, err := os.Create(path + ".html")
	if err != nil {
		log.Println("Ошибка создания файла", err)
		 return err
	}
	defer file.Close()

	_, err = file.Write(data)

	if err != nil {
		log.Println("Ошибка при записи в файл", err)
		return err
	}
	log.Println("Данные успешно записаны")
	return nil
}
