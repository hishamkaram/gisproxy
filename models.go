package gisproxy

import (
	"fmt"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//ServerAuth ServerAuth
type ServerAuth struct {
	gorm.Model
	UserName    string `gorm:"size:255;"`
	Password    string `sql:"type:text;"`
	AccessToken string `sql:"type:text;"`
	Type        string `gorm:"size:100;"`
}

//Server gisproxy server model
type Server struct {
	gorm.Model
	Name       string     `gorm:"size:255;unique;not null"`
	URL        string     `sql:"type:text;not null"`
	ServerType string     `gorm:"type:varchar(100);not null"`
	Active     bool       `gorm:"default:true;not null"`
	AuthInfo   ServerAuth `gorm:"foreignkey:AuthInfoID"`
	AuthInfoID uint
	Layers     []Layer `gorm:"ForeignKey:ServerID"`
}

//User gisproxy user model
type User struct {
	gorm.Model
	FirstName    string
	LastName     string
	UserName     string   `gorm:"size:255;unique;not null"`
	Password     string   `sql:"type:text;not null"`
	Email        string   `gorm:"type:varchar(100);unique_index;not null"`
	IsSuperUser  bool     `gorm:"default:false"`
	IsStaff      bool     `gorm:"default:false"`
	Layers       []Layer  `gorm:"ForeignKey:UserID"`
	SharedLayers []*Layer `gorm:"many2many:layer_permissions;;foreignkey:LayerID"`
}

//LayerPermission Layer Permissions
type LayerPermission struct {
	gorm.Model
	CanEdit     bool `gorm:"default:true"`
	CanView     bool `gorm:"default:true"`
	CanDelete   bool `gorm:"default:true"`
	CanDownload bool `gorm:"default:true"`
	LayerID     uint
	UserID      uint
}

//LayerURL Layer Base URL
type LayerURL struct {
	gorm.Model
	URL     string `gorm:"not null"`
	Type    string `gorm:"not null"`
	LayerID uint
}

//Layer Layer Def
type Layer struct {
	gorm.Model
	Name     string `gorm:"size:255;not null"`
	ServerID uint
	UserID   uint
	URLS     []LayerURL `gorm:"ForeignKey:LayerID"`
	Users    []*User    `gorm:"many2many:layer_permissions;foreignkey:UserID"`
}

//BeforeCreate hashing password before insert
func (user *User) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("Password", HashPassword(user.Password))
	return nil
}

//ChangePassword hashing password before insert
func (user *User) ChangePassword(server *GISProxy, password string) error {
	db, err := gorm.Open("postgres", server.DB.buildConnectionStr())
	db.Model(&user).Update("password", HashPassword(password))
	return err
}
func (databaseConn *DBConnection) buildConnectionStr() string {
	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", databaseConn.Host, strconv.Itoa(databaseConn.Port), databaseConn.Username, databaseConn.Name, databaseConn.Password)
	return connStr
}

//MigrateDatabase MigrateDatabase
func (server *GISProxy) MigrateDatabase() {
	db, err := gorm.Open("postgres", server.DB.buildConnectionStr())
	if err != nil {
		server.logger.Error(err)
	}
	defer db.Close()
	db.AutoMigrate(&User{}, &ServerAuth{}, &LayerURL{}, &Layer{}, &Server{}, &LayerPermission{})
}

//LoadData load default data
func (server *GISProxy) LoadData() {
	db, err := gorm.Open("postgres", server.DB.buildConnectionStr())
	if err != nil {
		server.logger.Error(err)
	}
	defer db.Close()
	db.LogMode(true)
	var user User
	db.FirstOrCreate(&user, User{UserName: "admin", FirstName: "Hisham", Password: "admin", LastName: "Karam", Email: "admin@admin.com"})
	// layerServer := Server{
	// 	URL:        "http://localhost:8080/geoserver",
	// 	ServerType: "geoserver",
	// 	Name:       "geoserver2",
	// 	AuthInfo: ServerAuth{
	// 		UserName: "admin",
	// 		Password: "geoserver",
	// 		Type:     "Basic",
	// 	},
	// 	Active: true,
	// 	Layers: []Layer{
	// 		Layer{
	// 			Name: "geonode:other_healthcare_60cfefd3",
	// 		},
	// 	},
	// }
	// db.Create(&layerServer)
	// layerServer.Layers = append(layerServer.Layers, Layer{Name: "hisham:other_healthcare_60cfefd3"})
	// db.Save(&layerServer)
	// var users []User
	// var servers []Server
	// db.Find(&users)
	// db.Preload("Layers").Find(&servers)
	// for _, user := range users {
	// 	fmt.Println(user.UserName)
	// 	// fmt.Println(user.Layers)
	// }
	// for _, server := range servers {
	// 	fmt.Println(server.ServerType)
	// 	fmt.Println(server.AuthInfo)
	// 	fmt.Println(server.Layers)
	// 	fmt.Println(server.AuthInfo)
	// }
}
