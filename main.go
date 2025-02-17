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
	"sync"
	"time"
)

func main() {

	start := time.Now()		
	defer func() {		
		log.Println("Прошло времени", time.Since(start))	
	}()
	path_output := flag.String("src", "", "Путь к файлу с URL-адресами. Укажите полный путь к файлу, содержащему список URL, которые нужно обработать.")
	path_input := flag.String("dst", "", "Путь к директории для сохранения HTML-страниц. Укажите полный путь к директории, где будут сохранены загруженные HTML-страницы.")

	flag.Parse()	

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
		go func () {
			defer wg.Done()

			url := "http://" + domen
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
				createHtml(*path_input+"/"+domen, b)
			} else {
				log.Println("Ошибка правильность url адреса")
			}
		}()
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
