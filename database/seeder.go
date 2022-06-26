package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"libsysfo-server/utility"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/lib/pq"
	"gorm.io/datatypes"
)

func SeedProfileAccount() {
	data := []ProfileAccount{}

	for c := 0; c < 30; c++ {
		data = append(data, ProfileAccount{
			Id:       c + 1,
			Username: gofakeit.Gamertag(),
			Email:    gofakeit.Email(),
			Password: "f5bb0c8de146c67b44babbf4e6584cc0",
		})
	}
	DB.Create(&data)
}

func SeedProfileData() {
	currentTime := time.Now()
	rand.Seed(currentTime.UnixNano())
	data := []ProfileData{}
	gender := []string{"M", "F"}

	for c := 7; c < 30; c++ {
		birthDate := utility.DateRandom("1900-01-01", "2016-01-01").Format(utility.Dmy)
		job := gofakeit.Job()
		singleData := ProfileData{
			UserId:       c + 1,
			Name:         gofakeit.Name(),
			Gender:       gender[rand.Intn(2)],
			PlaceOfBirth: gofakeit.Address().City,
			DateOfBirth:  birthDate,
			Address1:     gofakeit.Address().Address,
			Profession:   job.Title,
			Institution:  job.Company,
			PhoneNo:      gofakeit.Phone(),
			IsWhatsapp:   (rand.Intn(2) == 0),
			Images:       "https://i0.wp.com/global.ac.id/wp-content/uploads/2015/04/speaker-3-v2.jpg?fit=768%2C768&ssl=1",
		}
		if rand.Intn(2) == 0 {
			singleData.VerifiedAt = &currentTime
		}
		data = append(data, singleData)
	}

	DB.Create(&data)
}

func SeedBook() {
	data := []Book{}
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("GET", "http://localhost:8000/api/books?page=1", nil)
	if err != nil {
		log.Fatal("error nih")
		return
	}
	req.Header.Set("user-agent", "golang application")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer 9|IBwB98fxJElxhZPSj4IjPFyQ6lgMYbDvxYD5g9E9")
	response, err := client.Do(req)
	if err != nil {
		log.Fatal("error nih")
		return
	}
	defer response.Body.Close()
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal("error nih")
		return
	}

	var template interface{}

	err = json.Unmarshal(b, &template)
	if err != nil {
		log.Fatal("error nih")
		return
	}

	responseBody := template.(map[string]interface{})["books"].([]interface{})

	for _, singleData := range responseBody {
		content := singleData.(map[string]interface{})
		data = append(data, Book{
			Image:  content["image"].(string),
			Title:  content["title"].(string),
			Author: content["author"].(string),
			Slug:   content["slug"].(string),
		})
	}

	DB.Create(&data)
}

func SeedBookDetail() {
	var dataDetails []BookDetail
	var dataBooks []*Book
	DB.Find(&dataBooks)
	for _, dataBook := range dataBooks {
		client := &http.Client{
			Timeout: time.Second * 10,
		}
		link := fmt.Sprintf("http://localhost:8000/api/books/%s/detail", dataBook.Slug)
		req, err := http.NewRequest("GET", link, nil)
		if err != nil {
			log.Fatal("error nih")
			return
		}
		req.Header.Set("user-agent", "golang application")
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Authorization", "Bearer 9|IBwB98fxJElxhZPSj4IjPFyQ6lgMYbDvxYD5g9E9")
		response, err := client.Do(req)
		if err != nil {
			log.Fatal("error nih")
			return
		}
		defer response.Body.Close()
		b, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal("error nih")
			return
		}

		var template interface{}

		err = json.Unmarshal(b, &template)
		if err != nil {
			log.Fatal("error nih")
			return
		}
		responseBody := template.(map[string]interface{})["book"].(map[string]interface{})["detail"].(map[string]interface{})
		dataDetails = append(dataDetails, BookDetail{
			Id:          dataBook.Id,
			ReleaseDate: responseBody["release_date"].(string),
			Description: responseBody["description"].(string),
			Language:    responseBody["language"].(string),
			Country:     responseBody["country"].(string),
			Publisher:   responseBody["publisher"].(string),
			PageCount:   responseBody["page_count"].(float64),
			Category:    responseBody["category"].(string),
		})
	}
	DB.Create(&dataDetails)
}

func SeedLibraryData() {
	var data []LibraryData
	content := []string{
		"https://i.pinimg.com/474x/24/ef/00/24ef0042f7c07e1aa47280106461b853.jpg",
		"https://api.designcitylab.com/public/images/article-images/Beijing-Sub-Centre-Library-02_HK_N349156.jpg",
		"https://images.adsttc.com/media/images/5ddf/ad94/3312/fdb8/d300/011d/slideshow/_A7R9702-HDR_HUNDVEN-CLEMENTS_PHOTOGRAPHY.jpg?1574940024",
		"https://images.adsttc.com/media/images/5fd1/63ba/63c0/17e8/4500/0025/slideshow/LBCC_Library60643.jpg?1607557956",
		"https://archello.s3.eu-central-1.amazonaws.com/images/2021/03/03/g.o.-architecture-the-small--green-library-community-centres-archello.1614730542.5493.jpg",
		"https://archello.s3.eu-central-1.amazonaws.com/images/2021/03/03/g.o.-architecture-the-small--green-library-community-centres-archello.1614730556.2039.jpg",
	}

	data = append(data, LibraryData{
		UserId:        2,
		Name:          "dinas perpustakaan dan kearsipan kota padang",
		Address:       "Jl. Batang Anai, Rimbo Kaluang, Kec. Padang Bar., Kota Padang, Sumatera Barat",
		Coordinate:    pq.Float64Array{100.35884639177868, -0.9266827607129856},
		Description:   gofakeit.LoremIpsumParagraph(2, 5, 5, " "),
		ImagesMain:    "https://i.pinimg.com/originals/ff/96/ee/ff96eecd5f94fc82b561ef2812c541de.jpg",
		ImagesContent: pq.StringArray(content),
		Webpage:       "unand.ac.id",
	}, LibraryData{
		UserId:        3,
		Name:          "perpustakaan universitas andalas",
		Address:       "3FP6+M4V Kampus Universitas Andalas, Limau Manis, Kec. Pauh, Kota Padang, Sumatera Barat 25175",
		Coordinate:    pq.Float64Array{100.46029408540059, -0.9132586264043556},
		Description:   gofakeit.LoremIpsumParagraph(2, 5, 5, " "),
		ImagesMain:    "https://pustaka.unand.ac.id/images/perpustakaan.jpg",
		ImagesContent: pq.StringArray(content),
		Webpage:       "unand.ac.id",
	}, LibraryData{
		UserId:        4,
		Name:          "perpustakaan pusat uin imam bonjol padang",
		Address:       "399P+PVR, Kampus UIN Imam Bonjol Jl. Prof. Mahmud Yunus, Lubuk Lintah, Kec. Kuranji, Kota Padang, Sumatera Barat 25176",
		Coordinate:    pq.Float64Array{-0.930441718578502, 100.38715866243624},
		Description:   gofakeit.LoremIpsumParagraph(2, 5, 5, " "),
		ImagesMain:    "https://asset.kompas.com/crops/gj4bxVEM-ombeC7YhdMPWTQqMwA=/0x67:800x600/750x500/data/photo/2018/01/06/3283493641.jpg",
		ImagesContent: pq.StringArray(content),
		Webpage:       "unand.ac.id",
	}, LibraryData{
		UserId:        5,
		Name:          "perpustakaan universitas bung hatta",
		Address:       "38VV+HQ5, North Ulak Karang, Padang Utara, Padang City, West Sumatra",
		Coordinate:    pq.Float64Array{-0.9062321120689367, 100.34450263160929},
		Description:   gofakeit.LoremIpsumParagraph(2, 5, 5, " "),
		ImagesMain:    "https://i.pinimg.com/originals/ff/96/ee/ff96eecd5f94fc82b561ef2812c541de.jpg",
		ImagesContent: pq.StringArray(content),
		Webpage:       "unand.ac.id",
	}, LibraryData{
		UserId:        6,
		Name:          "UNP central library",
		Address:       "483W+HR9, West Air Tawar, Padang Utara, Padang City, West Sumatra",
		Coordinate:    pq.Float64Array{-0.8960782388043234, 100.34692795143607},
		Description:   gofakeit.LoremIpsumParagraph(2, 5, 5, " "),
		ImagesMain:    "https://www.agati.com/wp-content/uploads/2017/06/Diane-Lam-Blog-header.jpg",
		ImagesContent: pq.StringArray(content),
		Webpage:       "unand.ac.id",
	}, LibraryData{
		UserId:        7,
		Name:          "Perpustakaan Amanah",
		Address:       "Bundo Kanduong No.1, Belakang Tangsi, Kec. Padang Bar., Kota Padang, Sumatera Barat",
		Coordinate:    pq.Float64Array{-0.9519855092298098, 100.359963845534},
		Description:   gofakeit.LoremIpsumParagraph(2, 5, 5, " "),
		ImagesMain:    "https://s26162.pcdn.co/wp-content/uploads/2021/01/bookshelf1.jpg",
		ImagesContent: pq.StringArray(content),
		Webpage:       "unand.ac.id",
	})

	DB.Create(&data)
}

func SeedLibraryCollection() {
	var data []LibraryCollection
	currentTime := time.Now()
	rand.Seed(currentTime.UnixNano())

	for i := 0; i < 6; i++ {
		for j := 0; j < 10; j++ {
			data = append(data, LibraryCollection{
				LibraryId:    i + 1,
				BookId:       rand.Intn(23) + 1,
				Availability: true,
				Status:       rand.Intn(4) + 1,
			})
		}
	}
	DB.Create(&data)
}

func SeedLibraryPaper() {
	var data []LibraryPaper
	currentTime := time.Now()
	rand.Seed(currentTime.UnixNano())
	for i := 0; i < 6; i++ {
		for j := 0; j < 10; j++ {
			data = append(data, LibraryPaper{
				LibraryId:   i + 1,
				Title:       gofakeit.LoremIpsumSentence(5),
				Subject:     pq.StringArray{gofakeit.LoremIpsumWord(), gofakeit.LoremIpsumWord()},
				Abstract:    gofakeit.LoremIpsumParagraph(3, 3, 10, " "),
				Issn:        "8888888888888888",
				Description: datatypes.JSON(`{"foo": "abstract", "bar": "nice"}`),
				Access:      (rand.Intn(2) == 0),
				PaperUrl:    "https://drive.google.com/file/d/0B4eE3EAAsV6jaXg5SXQweDUyc28/view?resourcekey=0-QeSPnTIRa2FWntQ-9ev6wQ",
			})
		}
	}
	DB.Create(&data)
}

func SeedLibraryPaperPermission() {
	var data []LibraryPaperPermission
	currentTime := time.Now()
	rand.Seed(currentTime.UnixNano())
	for i := 0; i < 60; i++ {
		data = append(data, LibraryPaperPermission{
			PaperId:     rand.Intn(60) + 1,
			UserId:      rand.Intn(23) + 8,
			RedirectUrl: "https://drive.google.com/file/d/0B4eE3EAAsV6jaXg5SXQweDUyc28/view?resourcekey=0-QeSPnTIRa2FWntQ-9ev6wQ",
		})
	}
	DB.Create((&data))
}
