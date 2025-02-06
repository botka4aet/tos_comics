package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type FileId struct {
	gorm.Model
	//Делаем ссылку уникальной, иначе gorm clause работать не будет
	Link string `gorm:"unique"`
	Code string
}

func server_init() {
	fill_base_solved()
	fill_base_empty()
	queue(Task{},3)
}


func get_link_sql() (link_result string) {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{
		//Меняем уровень логера, тк при создании 200+ строк выскакивают предупреждения о скорости
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		return
	}
	db.AutoMigrate(&FileId{})
	var result []FileId
	db.Where("code = ?", "").Order("updated_at").First(&result)
	if len(result) == 0 {
		return
	}
	link_result = result[0].Link
	db.Model(&FileId{}).Where("Link = ?", link_result).Update("code", nil)
	//db.Order("updated_at").First(&result)
	//Закрываем соединения, тк позже база переоткрывается
	sqlDB, _ := db.DB()
	sqlDB.Close()
	return
}


func fill_base_solved() {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{
		//Меняем уровень логера, тк при создании 200+ строк выскакивают предупреждения о скорости
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		return
	}
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	db.AutoMigrate(&FileId{})

	file, err := os.Open("txtfiles\\links.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var files = []FileId{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		ntext := strings.TrimSuffix(text, "@2x.webp")
		ntext = strings.TrimSuffix(ntext, ".webp")
		ntext = strings.TrimSuffix(ntext, ".jpg")
		if text != ntext {
			files = append(files, FileId{Link: ntext[:len(ntext)-6], Code: ntext[len(ntext)-5:]})
		}
	}
	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "link"}},
		DoUpdates: clause.AssignmentColumns([]string{"code"}),
	}).CreateInBatches(&files, 1000)
}

func fill_base_empty() {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{
		//Меняем уровень логера, тк при создании 200+ строк выскакивают предупреждения о скорости
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		return
	}
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	db.AutoMigrate(&FileId{})

	file, err := os.Open("txtfiles\\check.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var files = []FileId{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		files = append(files, FileId{Link: text})
	}
	db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(&files, 1000)
}

func check_link(info Task) (result bool) {
	conn := &http.Client{
		//	Timeout:   4 * time.Second,
	}
	link, _ := get_suffix_rune(&info.Link)
	if link == "" {
		return
	}
	link = "https://cdn.townofsins.com/media/assets/images/" + info.Link + "_" + info.Code + link
	for {
		res, err := conn.Head(link)
		if err == nil && res != nil && res.StatusCode == 200 {
			return true
		} else if err == nil && res != nil && res.StatusCode == 500 {
			break
		}
	}
	return
}

// mode:
// 1 - создание; 2 - удаление, 3 - заполнение карты
func queue(info Task, mode uint8) {
	db, err := gorm.Open(sqlite.Open("queue.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		log.Println(err,info,mode)
		return
	}
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	db.AutoMigrate(&FileId{})
	if  mode == 1 {
		db.Clauses(clause.OnConflict{DoNothing: true}).Create(&FileId{Link: info.Link, Code: info.Secret})
	} else if mode == 2 {
		var result []FileId
		db.Where("Link = ?", info.Link).Delete(&result)
	} else if mode == 3 {
		rows, _ := db.Model(&FileId{}).Rows() 
		for rows.Next() {
			var result FileId
			db.ScanRows(rows, &result)
			TaskList[result.Link] = result.Code	
		}
		// for i,j  := range TaskList {
		// 	println(i,j)
		// }	
	}
	return
}


