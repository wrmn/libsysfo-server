package server

import (
	"libsysfo-server/database"
	"libsysfo-server/database/book"
	"libsysfo-server/database/library"
	"libsysfo-server/utility"
)

var booksData = database.BooksData{
	Data: []book.BookData{
		{
			Id:        1,
			Title:     "something",
			Author:    "someone",
			Publisher: "somewhat",
		},
	},
}

var librariesData = database.LibrariesData{
	Data: []library.LibraryData{

		{
			Id: 1,

			Name:    "dinas perpustakaan dan kearsipan kota padang",
			Address: "Jl. Batang Anai, Rimbo Kaluang, Kec. Padang Bar., Kota Padang, Sumatera Barat",
			Coordinate: utility.Location{
				Latitude:  -0.9266827607129856,
				Longitude: 100.35884639177868,
			},
			Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
			Images: utility.ImagesData{
				Main: "https://i.pinimg.com/originals/ff/96/ee/ff96eecd5f94fc82b561ef2812c541de.jpg",
				Content: []string{
					"https://i.pinimg.com/474x/24/ef/00/24ef0042f7c07e1aa47280106461b853.jpg",
					"https://api.designcitylab.com/public/images/article-images/Beijing-Sub-Centre-Library-02_HK_N349156.jpg",
					"https://images.adsttc.com/media/images/5ddf/ad94/3312/fdb8/d300/011d/slideshow/_A7R9702-HDR_HUNDVEN-CLEMENTS_PHOTOGRAPHY.jpg?1574940024",
					"https://images.adsttc.com/media/images/5fd1/63ba/63c0/17e8/4500/0025/slideshow/LBCC_Library60643.jpg?1607557956",
					"https://archello.s3.eu-central-1.amazonaws.com/images/2021/03/03/g.o.-architecture-the-small--green-library-community-centres-archello.1614730542.5493.jpg",
					"https://archello.s3.eu-central-1.amazonaws.com/images/2021/03/03/g.o.-architecture-the-small--green-library-community-centres-archello.1614730556.2039.jpg",
				},
			},
		}, {
			Id:      2,
			Name:    "taman bacaan p-Mee",
			Address: "Jl. Gajah Mada No.46, RW.04, Gn. Pangilun, Kec. Padang Utara, Kota Padang, Sumatera Barat 25137",
			Coordinate: utility.Location{
				Latitude:  -0.9189879475080617,
				Longitude: 100.36549420528756,
			},
			Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
			Images: utility.ImagesData{
				Main: "https://www.gannett-cdn.com/presto/2022/03/11/PCNJ/1fb8bf4f-e83a-4c3a-aef5-72924721aff9-New_Brunswick_Library_3.10.22-1.jpg?crop=5745,3232,x1,y177&width=3200&height=1801&format=pjpg&auto=webp",
				Content: []string{
					"https://i.pinimg.com/474x/24/ef/00/24ef0042f7c07e1aa47280106461b853.jpg",
					"https://api.designcitylab.com/public/images/article-images/Beijing-Sub-Centre-Library-02_HK_N349156.jpg",
					"https://images.adsttc.com/media/images/5ddf/ad94/3312/fdb8/d300/011d/slideshow/_A7R9702-HDR_HUNDVEN-CLEMENTS_PHOTOGRAPHY.jpg?1574940024",
					"https://images.adsttc.com/media/images/5fd1/63ba/63c0/17e8/4500/0025/slideshow/LBCC_Library60643.jpg?1607557956",
					"https://archello.s3.eu-central-1.amazonaws.com/images/2021/03/03/g.o.-architecture-the-small--green-library-community-centres-archello.1614730542.5493.jpg",
					"https://archello.s3.eu-central-1.amazonaws.com/images/2021/03/03/g.o.-architecture-the-small--green-library-community-centres-archello.1614730556.2039.jpg",
				},
			},
		}, {
			Id:      3,
			Name:    "UNP central library",
			Address: "483W+HR9, West Air Tawar, Padang Utara, Padang City, West Sumatra",
			Coordinate: utility.Location{
				Latitude:  -0.8960782388043234,
				Longitude: 100.34692795143607,
			},
			Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
			Images: utility.ImagesData{
				Main: "https://www.agati.com/wp-content/uploads/2017/06/Diane-Lam-Blog-header.jpg",
				Content: []string{
					"https://i.pinimg.com/474x/24/ef/00/24ef0042f7c07e1aa47280106461b853.jpg",
					"https://api.designcitylab.com/public/images/article-images/Beijing-Sub-Centre-Library-02_HK_N349156.jpg",
					"https://images.adsttc.com/media/images/5ddf/ad94/3312/fdb8/d300/011d/slideshow/_A7R9702-HDR_HUNDVEN-CLEMENTS_PHOTOGRAPHY.jpg?1574940024",
					"https://images.adsttc.com/media/images/5fd1/63ba/63c0/17e8/4500/0025/slideshow/LBCC_Library60643.jpg?1607557956",
					"https://archello.s3.eu-central-1.amazonaws.com/images/2021/03/03/g.o.-architecture-the-small--green-library-community-centres-archello.1614730542.5493.jpg",
					"https://archello.s3.eu-central-1.amazonaws.com/images/2021/03/03/g.o.-architecture-the-small--green-library-community-centres-archello.1614730556.2039.jpg",
				},
			},
		}, {
			Id:      4,
			Name:    "rumah baca darman moenir",
			Address: "Jl. Pasaman Jl. Pagang Raya-Siteba No.170, Surau Gadang, Kec. Nanggalo, Kota Padang, Sumatera Barat 25146",
			Coordinate: utility.Location{
				Latitude:  -0.8952247516786328,
				Longitude: 100.3677440695125,
			},
			Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
			Images: utility.ImagesData{
				Main: "https://media.gettyimages.com/photos/young-boy-climbing-steps-to-a-library-building-picture-id157314421?s=612x612",
				Content: []string{
					"https://i.pinimg.com/474x/24/ef/00/24ef0042f7c07e1aa47280106461b853.jpg",
					"https://api.designcitylab.com/public/images/article-images/Beijing-Sub-Centre-Library-02_HK_N349156.jpg",
					"https://images.adsttc.com/media/images/5ddf/ad94/3312/fdb8/d300/011d/slideshow/_A7R9702-HDR_HUNDVEN-CLEMENTS_PHOTOGRAPHY.jpg?1574940024",
					"https://images.adsttc.com/media/images/5fd1/63ba/63c0/17e8/4500/0025/slideshow/LBCC_Library60643.jpg?1607557956",
					"https://archello.s3.eu-central-1.amazonaws.com/images/2021/03/03/g.o.-architecture-the-small--green-library-community-centres-archello.1614730542.5493.jpg",
					"https://archello.s3.eu-central-1.amazonaws.com/images/2021/03/03/g.o.-architecture-the-small--green-library-community-centres-archello.1614730556.2039.jpg",
				},
			},
		}, {
			Id:      5,
			Name:    "Perpustakaan Amanah",
			Address: "Bundo Kanduong No.1, Belakang Tangsi, Kec. Padang Bar., Kota Padang, Sumatera Barat",
			Coordinate: utility.Location{
				Latitude:  -0.9519855092298098,
				Longitude: 100.359963845534,
			},
			Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
			Images: utility.ImagesData{
				Main: "https://s26162.pcdn.co/wp-content/uploads/2021/01/bookshelf1.jpg",
				Content: []string{
					"https://i.pinimg.com/474x/24/ef/00/24ef0042f7c07e1aa47280106461b853.jpg",
					"https://api.designcitylab.com/public/images/article-images/Beijing-Sub-Centre-Library-02_HK_N349156.jpg",
					"https://images.adsttc.com/media/images/5ddf/ad94/3312/fdb8/d300/011d/slideshow/_A7R9702-HDR_HUNDVEN-CLEMENTS_PHOTOGRAPHY.jpg?1574940024",
					"https://images.adsttc.com/media/images/5fd1/63ba/63c0/17e8/4500/0025/slideshow/LBCC_Library60643.jpg?1607557956",
					"https://archello.s3.eu-central-1.amazonaws.com/images/2021/03/03/g.o.-architecture-the-small--green-library-community-centres-archello.1614730542.5493.jpg",
					"https://archello.s3.eu-central-1.amazonaws.com/images/2021/03/03/g.o.-architecture-the-small--green-library-community-centres-archello.1614730556.2039.jpg",
				},
			},
		}, {
			Id:      6,
			Name:    "perpustakaan daerah provinsi sumatera barat",
			Address: "Diponegoro nmr 4, Belakang Tangsi, Kec. Padang Bar., Kota Padang, Sumatera Barat 25118",
			Coordinate: utility.Location{
				Latitude:  -0.9535401233912927,
				Longitude: 100.35610911181004,
			},
			Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
			Images: utility.ImagesData{
				Main: "https://4.bp.blogspot.com/-w9FsGMYaBiY/UkyLVb3RnbI/AAAAAAAAAKo/eT1FxTxUKsM/s1600/Library+Building.JPG",
				Content: []string{
					"https://i.pinimg.com/474x/24/ef/00/24ef0042f7c07e1aa47280106461b853.jpg",
					"https://api.designcitylab.com/public/images/article-images/Beijing-Sub-Centre-Library-02_HK_N349156.jpg",
					"https://images.adsttc.com/media/images/5ddf/ad94/3312/fdb8/d300/011d/slideshow/_A7R9702-HDR_HUNDVEN-CLEMENTS_PHOTOGRAPHY.jpg?1574940024",
					"https://images.adsttc.com/media/images/5fd1/63ba/63c0/17e8/4500/0025/slideshow/LBCC_Library60643.jpg?1607557956",
					"https://archello.s3.eu-central-1.amazonaws.com/images/2021/03/03/g.o.-architecture-the-small--green-library-community-centres-archello.1614730542.5493.jpg",
					"https://archello.s3.eu-central-1.amazonaws.com/images/2021/03/03/g.o.-architecture-the-small--green-library-community-centres-archello.1614730556.2039.jpg",
				},
			},
		}, {
			Id:      7,
			Name:    "perpustakaan universitas andalas",
			Address: "3FP6+M4V Kampus Universitas Andalas, Limau Manis, Kec. Pauh, Kota Padang, Sumatera Barat 25175",
			Coordinate: utility.Location{
				Latitude:  -0.9132586264043556,
				Longitude: 100.46029408540059,
			},
			Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
			Images: utility.ImagesData{
				Main: "https://pustaka.unand.ac.id/images/perpustakaan.jpg",
				Content: []string{
					"https://i.pinimg.com/474x/24/ef/00/24ef0042f7c07e1aa47280106461b853.jpg",
					"https://api.designcitylab.com/public/images/article-images/Beijing-Sub-Centre-Library-02_HK_N349156.jpg",
					"https://images.adsttc.com/media/images/5ddf/ad94/3312/fdb8/d300/011d/slideshow/_A7R9702-HDR_HUNDVEN-CLEMENTS_PHOTOGRAPHY.jpg?1574940024",
					"https://images.adsttc.com/media/images/5fd1/63ba/63c0/17e8/4500/0025/slideshow/LBCC_Library60643.jpg?1607557956",
					"https://archello.s3.eu-central-1.amazonaws.com/images/2021/03/03/g.o.-architecture-the-small--green-library-community-centres-archello.1614730542.5493.jpg",
					"https://archello.s3.eu-central-1.amazonaws.com/images/2021/03/03/g.o.-architecture-the-small--green-library-community-centres-archello.1614730556.2039.jpg",
				},
			},
		}, {
			Id:      8,
			Name:    "perpustakaan pusat uin imam bonjol padang",
			Address: "399P+PVR, Kampus UIN Imam Bonjol Jl. Prof. Mahmud Yunus, Lubuk Lintah, Kec. Kuranji, Kota Padang, Sumatera Barat 25176",
			Coordinate: utility.Location{
				Latitude:  -0.930441718578502,
				Longitude: 100.38715866243624,
			},
			Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
			Images: utility.ImagesData{
				Main: "https://asset.kompas.com/crops/gj4bxVEM-ombeC7YhdMPWTQqMwA=/0x67:800x600/750x500/data/photo/2018/01/06/3283493641.jpg",
				Content: []string{
					"https://i.pinimg.com/474x/24/ef/00/24ef0042f7c07e1aa47280106461b853.jpg",
					"https://api.designcitylab.com/public/images/article-images/Beijing-Sub-Centre-Library-02_HK_N349156.jpg",
					"https://images.adsttc.com/media/images/5ddf/ad94/3312/fdb8/d300/011d/slideshow/_A7R9702-HDR_HUNDVEN-CLEMENTS_PHOTOGRAPHY.jpg?1574940024",
					"https://images.adsttc.com/media/images/5fd1/63ba/63c0/17e8/4500/0025/slideshow/LBCC_Library60643.jpg?1607557956",
					"https://archello.s3.eu-central-1.amazonaws.com/images/2021/03/03/g.o.-architecture-the-small--green-library-community-centres-archello.1614730542.5493.jpg",
					"https://archello.s3.eu-central-1.amazonaws.com/images/2021/03/03/g.o.-architecture-the-small--green-library-community-centres-archello.1614730556.2039.jpg",
				},
			},
		}, {
			Id:      9,
			Name:    "perpustakaan masjid mujahiddin",
			Address: "Jl. Sutan Syahrir No.135, Mata Air, Kec. Padang Sel., Kota Padang, Sumatera Barat",
			Coordinate: utility.Location{
				Latitude:  -0.9643840029104725,
				Longitude: 100.37922164020581,
			},
			Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
			Images: utility.ImagesData{
				Main: "https://img.beritasatu.com/cache/beritasatu/620x350-2/1601367160.jpg",
				Content: []string{
					"https://i.pinimg.com/474x/24/ef/00/24ef0042f7c07e1aa47280106461b853.jpg",
					"https://api.designcitylab.com/public/images/article-images/Beijing-Sub-Centre-Library-02_HK_N349156.jpg",
					"https://images.adsttc.com/media/images/5ddf/ad94/3312/fdb8/d300/011d/slideshow/_A7R9702-HDR_HUNDVEN-CLEMENTS_PHOTOGRAPHY.jpg?1574940024",
					"https://images.adsttc.com/media/images/5fd1/63ba/63c0/17e8/4500/0025/slideshow/LBCC_Library60643.jpg?1607557956",
					"https://archello.s3.eu-central-1.amazonaws.com/images/2021/03/03/g.o.-architecture-the-small--green-library-community-centres-archello.1614730542.5493.jpg",
					"https://archello.s3.eu-central-1.amazonaws.com/images/2021/03/03/g.o.-architecture-the-small--green-library-community-centres-archello.1614730556.2039.jpg",
				},
			},
		}, {
			Id:      10,
			Name:    "tbm lentera kota tua padang",
			Address: "Jl. Batang Arau No.52, Berok Nipah, Kec. Padang Bar., Kota Padang, Sumatera Barat",
			Coordinate: utility.Location{
				Latitude:  -0.9642306758290158,
				Longitude: 100.3600132688974,
			},
			Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
			Images: utility.ImagesData{
				Main: "http://www.brantas-abipraya.co.id/frontend/uploads/defaults/4i1mwN20170116151536.jpg",
				Content: []string{
					"https://i.pinimg.com/474x/24/ef/00/24ef0042f7c07e1aa47280106461b853.jpg",
					"https://api.designcitylab.com/public/images/article-images/Beijing-Sub-Centre-Library-02_HK_N349156.jpg",
					"https://images.adsttc.com/media/images/5ddf/ad94/3312/fdb8/d300/011d/slideshow/_A7R9702-HDR_HUNDVEN-CLEMENTS_PHOTOGRAPHY.jpg?1574940024",
					"https://images.adsttc.com/media/images/5fd1/63ba/63c0/17e8/4500/0025/slideshow/LBCC_Library60643.jpg?1607557956",
					"https://archello.s3.eu-central-1.amazonaws.com/images/2021/03/03/g.o.-architecture-the-small--green-library-community-centres-archello.1614730542.5493.jpg",
					"https://archello.s3.eu-central-1.amazonaws.com/images/2021/03/03/g.o.-architecture-the-small--green-library-community-centres-archello.1614730556.2039.jpg",
				},
			},
		}, {
			Id:      11,
			Name:    "perpustakaan universitas bung hatta",
			Address: "38VV+HQ5, North Ulak Karang, Padang Utara, Padang City, West Sumatra",
			Coordinate: utility.Location{
				Latitude:  -0.9062321120689367,
				Longitude: 100.34450263160929,
			},
			Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
			Images: utility.ImagesData{
				Main: "https://i.pinimg.com/originals/ff/96/ee/ff96eecd5f94fc82b561ef2812c541de.jpg",
				Content: []string{
					"https://i.pinimg.com/474x/24/ef/00/24ef0042f7c07e1aa47280106461b853.jpg",
					"https://api.designcitylab.com/public/images/article-images/Beijing-Sub-Centre-Library-02_HK_N349156.jpg",
					"https://images.adsttc.com/media/images/5ddf/ad94/3312/fdb8/d300/011d/slideshow/_A7R9702-HDR_HUNDVEN-CLEMENTS_PHOTOGRAPHY.jpg?1574940024",
					"https://images.adsttc.com/media/images/5fd1/63ba/63c0/17e8/4500/0025/slideshow/LBCC_Library60643.jpg?1607557956",
					"https://archello.s3.eu-central-1.amazonaws.com/images/2021/03/03/g.o.-architecture-the-small--green-library-community-centres-archello.1614730542.5493.jpg",
					"https://archello.s3.eu-central-1.amazonaws.com/images/2021/03/03/g.o.-architecture-the-small--green-library-community-centres-archello.1614730556.2039.jpg",
				},
			},
		},
	},
}
