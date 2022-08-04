package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"libsysfo-server/utility"
	bookserver "libsysfo-server/utility/book-server"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/lib/pq"
	"gorm.io/datatypes"
)

func SeedProfileAccount() {
	data := []ProfileAccount{}

	for c := 0; c < 30; c++ {
		uname := strings.ToLower(gofakeit.Gamertag())
		singelData := ProfileAccount{
			ID:       c + 1,
			Username: &uname,
			Email:    strings.ToLower(gofakeit.Email()),
			Password: "f5bb0c8de146c67b44babbf4e6584cc0",
		}
		if c == 0 {
			singelData.AccountType = 1
		} else if c < 7 {
			singelData.AccountType = 2
		} else {
			singelData.AccountType = 3
		}
		data = append(data, singelData)
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
		job.Title = strings.ToLower(job.Title)
		job.Company = strings.ToLower(job.Company)
		phoneNumber := gofakeit.Phone()
		phoneCode := "+62"
		singleData := ProfileData{
			UserID:       c + 1,
			Name:         gofakeit.Name(),
			Gender:       &(gender[rand.Intn(2)]),
			PlaceOfBirth: &(gofakeit.Address().City),
			DateOfBirth:  &(birthDate),
			Address1:     &(gofakeit.Address().Address),
			Profession:   &(job.Title),
			Institution:  &(job.Company),
			PhoneCode:    &(phoneCode),
			PhoneNo:      &(phoneNumber),
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

func SeedBookLocal() {
	content, err := ioutil.ReadFile("./sampleData/bookData.json")
	if err != nil {
		fmt.Println("Error when opening file: ", err)
		return
	}

	var payload bookDataset

	err = json.Unmarshal(content, &payload)
	if err != nil {
		fmt.Println("error Unmarshal")
		return
	}

	for i := range payload.Books {
		payload.Books[i].Slug = payload.Books[i].SlugGenerator()
		payload.Books[i].BookDetail.Description = gofakeit.LoremIpsumSentence(20)
		DB.Create(&payload.Books[i])
	}
}

func SeedBook() {
	data := []Book{}
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	for i := 0; i < 3; i++ {
		page := 3 + i
		link := fmt.Sprintf("%s/api/books?page=%d", os.Getenv("BOOK_SERVER_URL"), page)
		req, err := http.NewRequest("GET", link, nil)
		if err != nil {
			log.Fatal("error nih")
			return
		}
		req.Header.Set("user-agent", "golang application")
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Authorization", "Bearer 10|6UHnWG0z8pBYl60Dm0ioMBjwPGuRoGodYcr0X80o")
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
		if response.StatusCode == 200 {

			var template bookserver.BookResponse

			err = json.Unmarshal(b, &template)
			if err != nil {
				log.Fatal("error nih")
				return
			}

			responseBody := template.Books

			for _, content := range responseBody {
				data = append(data, Book{
					Image:  *content.Image,
					Title:  *content.Title,
					Author: *content.Author,
					Source: "gramedia",
					Slug:   *content.Slug,
				})
			}
		}
	}
	DB.Create(&data)
}

func SeedBookDetail() {
	var dataBooks []*Book
	DB.Find(&dataBooks)
	for i, dataBook := range dataBooks {
		if dataBooks[i].Source == "gramedia" {
			client := &http.Client{
				Timeout: time.Second * 10,
			}
			link := fmt.Sprintf("%s/api/books/%s/detail", os.Getenv("BOOK_SERVER_URL"), dataBook.Slug)
			req, err := http.NewRequest("GET", link, nil)
			if err != nil {
				log.Fatal("error nih")
				return
			}
			req.Header.Set("user-agent", "golang application")
			req.Header.Add("Accept", "application/json")
			req.Header.Add("Authorization", "Bearer 10|6UHnWG0z8pBYl60Dm0ioMBjwPGuRoGodYcr0X80o")
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
			if response.StatusCode == 200 {
				var template bookserver.BookResponse

				err = json.Unmarshal(b, &template)
				if err != nil {
					log.Fatal("error nih")
					return
				}

				responseBody := *template.Book.Detail

				dataDetails := BookDetail{
					ID:          dataBooks[i].ID,
					ReleaseDate: *responseBody.ReleaseDate,
					Description: *responseBody.Description,
					Language:    *responseBody.Language,
					Country:     *responseBody.Country,
					Publisher:   *responseBody.Publisher,
					PageCount:   int(*responseBody.PageCount),
					Category:    *responseBody.Category,
				}
				DB.Create(&dataDetails)
			}
		}
	}
}

func SeedLibraryData() {
	content := []string{
		"https://i.pinimg.com/474x/24/ef/00/24ef0042f7c07e1aa47280106461b853.jpg",
		"https://media.istockphoto.com/photos/library-bookshelves-with-books-and-textbooks-learning-and-education-picture-id1200326335?k=20&m=1200326335&s=612x612&w=0&h=TXy8Z48ULgGdJNWaNSXlGR5oQHCYD9rbBysf7U9w0HA=",
		"https://images.adsttc.com/media/images/5ddf/ad94/3312/fdb8/d300/011d/slideshow/_A7R9702-HDR_HUNDVEN-CLEMENTS_PHOTOGRAPHY.jpg?1574940024",
		"https://images.adsttc.com/media/images/5fd1/63ba/63c0/17e8/4500/0025/slideshow/LBCC_Library60643.jpg?1607557956",
		"https://archello.s3.eu-central-1.amazonaws.com/images/2021/03/03/g.o.-architecture-the-small--green-library-community-centres-archello.1614730542.5493.jpg",
		"https://archello.s3.eu-central-1.amazonaws.com/images/2021/03/03/g.o.-architecture-the-small--green-library-community-centres-archello.1614730556.2039.jpg",
	}

	data := []LibraryData{{
		UserID:        2,
		Name:          "dinas perpustakaan dan kearsipan kota padang",
		Address:       "Jl. Batang Anai, Rimbo Kaluang, Kec. Padang Bar., Kota Padang, Sumatera Barat",
		Coordinate:    pq.Float64Array{100.35884639177868, -0.9266827607129856},
		Description:   gofakeit.LoremIpsumParagraph(2, 5, 5, " "),
		ImagesMain:    "https://i.pinimg.com/originals/ff/96/ee/ff96eecd5f94fc82b561ef2812c541de.jpg",
		ImagesContent: pq.StringArray(content),
		Webpage:       "unand.ac.id",
	}, {
		UserID:        3,
		Name:          "perpustakaan universitas andalas",
		Address:       "3FP6+M4V Kampus Universitas Andalas, Limau Manis, Kec. Pauh, Kota Padang, Sumatera Barat 25175",
		Coordinate:    pq.Float64Array{100.46029408540059, -0.9132586264043556},
		Description:   gofakeit.LoremIpsumParagraph(2, 5, 5, " "),
		ImagesMain:    "https://pustaka.unand.ac.id/images/perpustakaan.jpg",
		ImagesContent: pq.StringArray(content),
		Webpage:       "unand.ac.id",
	}, {
		UserID:        4,
		Name:          "perpustakaan pusat uin imam bonjol padang",
		Address:       "399P+PVR, Kampus UIN Imam Bonjol Jl. Prof. Mahmud Yunus, Lubuk Lintah, Kec. Kuranji, Kota Padang, Sumatera Barat 25176",
		Coordinate:    pq.Float64Array{100.38715866243624, -0.930441718578502},
		Description:   gofakeit.LoremIpsumParagraph(2, 5, 5, " "),
		ImagesMain:    "https://asset.kompas.com/crops/gj4bxVEM-ombeC7YhdMPWTQqMwA=/0x67:800x600/750x500/data/photo/2018/01/06/3283493641.jpg",
		ImagesContent: pq.StringArray(content),
		Webpage:       "unand.ac.id",
	}, {
		UserID:        5,
		Name:          "perpustakaan universitas bung hatta",
		Address:       "38VV+HQ5, North Ulak Karang, Padang Utara, Padang City, West Sumatra",
		Coordinate:    pq.Float64Array{100.34450263160929, -0.9062321120689367},
		Description:   gofakeit.LoremIpsumParagraph(2, 5, 5, " "),
		ImagesMain:    "https://i.pinimg.com/originals/ff/96/ee/ff96eecd5f94fc82b561ef2812c541de.jpg",
		ImagesContent: pq.StringArray(content),
		Webpage:       "unand.ac.id",
	}, {
		UserID:        6,
		Name:          "UNP central library",
		Address:       "483W+HR9, West Air Tawar, Padang Utara, Padang City, West Sumatra",
		Coordinate:    pq.Float64Array{100.34692795143607, -0.8960782388043234},
		Description:   gofakeit.LoremIpsumParagraph(2, 5, 5, " "),
		ImagesMain:    "https://www.agati.com/wp-content/uploads/2017/06/Diane-Lam-Blog-header.jpg",
		ImagesContent: pq.StringArray(content),
		Webpage:       "unand.ac.id",
	}, {
		UserID:        7,
		Name:          "Perpustakaan Amanah",
		Address:       "Bundo Kanduong No.1, Belakang Tangsi, Kec. Padang Bar., Kota Padang, Sumatera Barat",
		Coordinate:    pq.Float64Array{100.359963845534, -0.9519855092298098},
		Description:   gofakeit.LoremIpsumParagraph(2, 5, 5, " "),
		ImagesMain:    "https://s26162.pcdn.co/wp-content/uploads/2021/01/bookshelf1.jpg",
		ImagesContent: pq.StringArray(content),
		Webpage:       "unand.ac.id",
	}}

	DB.Create(&data)
}

func SeedLibraryCollection() {
	var data []LibraryCollection
	currentTime := time.Now()
	rand.Seed(currentTime.UnixNano())
	sn := 8801

	for i := 0; i < 6; i++ {
		for j := 0; j < 100; j++ {
			sn++
			data = append(data, LibraryCollection{
				SerialNumber: fmt.Sprintf("1234.23.12.%d", sn),
				LibraryID:    i + 1,
				BookID:       rand.Intn(85) + 1,
				Availability: 1,
				Status:       rand.Intn(4) + 1,
			})
		}
	}
	DB.Create(&data)
}

func SeedLibraryCollectionBorrow() {
	currentTime := time.Now()
	rand.Seed(currentTime.UnixNano())

	for i := 0; i < 1000; i++ {
		collection := LibraryCollection{}

		statint := rand.Intn(5)
		randDate := utility.DateRandom("2022-01-01", "2022-07-31")
		acceptedDate := randDate.Add(time.Duration(rand.Intn(12)+1) * time.Hour)
		takedDate := acceptedDate.Add(24 * time.Hour)
		returnedDate := randDate.Add(time.Duration(24+rand.Intn(200)+1) * time.Hour)
		canceledDate := randDate.Add(time.Duration(rand.Intn(12)+13) * time.Hour)

		singleData := LibraryCollectionBorrow{
			CreatedAt:    randDate,
			CollectionID: rand.Intn(600) + 1,
			UserID:       rand.Intn(23) + 8,
		}

		DB.Where("id = ?", singleData.CollectionID).Find(&collection)

		if statint == 1 || statint == 2 || statint == 3 {
			singleData.AcceptedAt = &acceptedDate
		}

		if statint == 2 || statint == 3 {
			singleData.TakedAt = &takedDate
		}

		if statint == 3 || (statint == 2 && time.Since(takedDate).Hours() >= 168) {
			singleData.ReturnedAt = &returnedDate
		} else if statint == 2 && time.Since(takedDate).Hours() <= 168 {
			if collection.Availability == 1 {
				collection.Availability = 3
				DB.Save(&collection)
			} else {
				singleData.ReturnedAt = &returnedDate

			}

		}

		if (statint == 1 && rand.Intn(5)%2 == 0) || statint == 4 {
			singleData.CanceledAt = &canceledDate
		}
		if collection.Availability != 2 {
			DB.Create(&singleData)
		}
	}
}

func SeedLibraryPaper() {
	var data []LibraryPaper

	content, err := ioutil.ReadFile("./sampleData/paperData.json")
	if err != nil {
		fmt.Println("Error when opening file: ", err)
		return
	}

	var payload paperDataset

	err = json.Unmarshal(content, &payload)
	if err != nil {
		fmt.Println("error Unmarshal")
		return
	}

	for _, k := range payload.Papers {
		k.Abstract = gofakeit.LoremIpsumParagraph(3, 3, 10, " ")
		k.Access = (rand.Intn(2) == 0)
	}

	data = append(data, payload.Papers...)

	currentTime := time.Now()
	rand.Seed(currentTime.UnixNano())
	docType := []string{"journal", "thesis", "other document"}
	for i := 0; i < 6; i++ {
		for j := 0; j < 10; j++ {
			data = append(data, LibraryPaper{
				LibraryID:   i + 1,
				Title:       gofakeit.LoremIpsumSentence(5),
				Subject:     pq.StringArray{gofakeit.LoremIpsumWord(), gofakeit.LoremIpsumWord()},
				Abstract:    gofakeit.LoremIpsumParagraph(3, 3, 10, " "),
				Type:        docType[rand.Intn(3)],
				Description: datatypes.JSON(`{"foo": "abstract", "bar": "nice"}`),
				Access:      (rand.Intn(2) == 0),
				PaperUrl:    "https://ik.imagekit.io/libsysfo/test.pdf",
			})
		}
	}
	DB.Create(&data)
}

func SeedLibraryPaperPermission() {
	var data []LibraryPaperPermission
	currentTime := time.Now()
	rand.Seed(currentTime.UnixNano())
	for i := 0; i < 200; i++ {
		paperId := rand.Intn(80) + 1
		singleData := LibraryPaperPermission{
			PaperID: paperId,
			UserID:  rand.Intn(23) + 8,
			Purpose: gofakeit.LoremIpsumSentence(10),
		}
		if rand.Intn(2) == 0 {
			singleData.AcceptedAt = &currentTime
		} else {
			singleData.CanceledAt = &currentTime
		}
		data = append(data, singleData)

	}
	DB.Create(&data)
}

func SeedLibraryPaperAccess() {
	var data []LibraryPaperAccess
	currentTime := time.Now()
	rand.Seed(currentTime.UnixNano())
	for i := 0; i < 1000; i++ {
		var permission LibraryPaperPermission
		permissionId := rand.Intn(200) + 1
		DB.Where("Id = ?", permissionId).Find(&permission)
		if permission.AcceptedAt != nil {
			data = append(data, LibraryPaperAccess{
				CreatedAt:    utility.DateRandom("2022-01-01", "2022-07-31"),
				PermissionID: permissionId,
			})
		}
	}
	DB.Create(&data)
}
