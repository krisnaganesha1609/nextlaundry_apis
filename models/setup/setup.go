package setup

import (
	"fmt"
	"log"
	m "nextlaundry_apis/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dbName := "laundry_ukk" //rename the database to match in local mysql
	database, err := gorm.Open(mysql.Open("root:@tcp(127.0.0.1:3306)/"+dbName+"?charset=utf8mb4&parseTime=true&loc=Local"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		panic("Can't Connect To Database!")
	}

	log.Println("Connected To Database!")

	DB = database
}

func Migrations() {

	if err1 := DB.AutoMigrate(&m.Users{}); err1 != nil {
		log.Println(err1.Error())
	}
	if err2 := DB.AutoMigrate(&m.Outlets{}); err2 != nil {
		log.Println(err2.Error())
	}
	if err3 := DB.AutoMigrate(&m.Products{}); err3 != nil {
		log.Println(err3.Error())
	}
	if err4 := DB.AutoMigrate(&m.Members{}); err4 != nil {
		log.Println(err4.Error())
	}
	if err5 := DB.AutoMigrate(&m.Transactions{}); err5 != nil {
		log.Println(err5.Error())
	}
	if err6 := DB.AutoMigrate(&m.TransactionDetails{}); err6 != nil {
		log.Println(err6.Error())
	}
	if err7 := DB.AutoMigrate(&m.LogHistory{}); err7 != nil {
		log.Println(err7.Error())
	}
	UsersTriggers()
	OutletTriggers()
	MemberTriggers()
	PackageTriggers()

	// DB.Exec("CREATE OR REPLACE TRIGGER `trigger_auction` BEFORE INSERT ON `auction_histories` FOR EACH ROW UPDATE `auctions` SET last_price = new.price, bidder_id = new.user_id, bidders_count = bidders_count + 1 WHERE id = new.auction_id AND new.price > last_price AND status = 'open';")
	// DB.Exec("CREATE OR REPLACE PROCEDURE update_product(IN p_product_id VARCHAR(255), IN p_user_id INT, IN p_name VARCHAR(255), IN p_description TEXT, IN p_price DECIMAL(10,2), IN p_image VARCHAR(255)) BEGIN UPDATE products SET name = p_name, description = p_description, price = p_price, image = p_image WHERE products.id = (SELECT auctions.product_id FROM auctions WHERE auctions.product_id = p_product_id AND auctions.user_id = p_user_id); END;")
	// DB.Exec("CREATE OR REPLACE PROCEDURE get_product(IN user_id INT) BEGIN SELECT products.* FROM products INNER JOIN auctions ON products.id = auctions.product_id WHERE auctions.user_id = user_id; END")

	fmt.Println("Database Migrated!")
}
