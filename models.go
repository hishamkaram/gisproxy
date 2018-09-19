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
	FirstName       string
	LastName        string
	UserName        string  `gorm:"size:255;unique;not null"`
	Password        string  `sql:"type:text;not null"`
	Email           string  `gorm:"type:varchar(100);unique_index;not null"`
	IsSuperUser     bool    `gorm:"default:false"`
	IsStaff         bool    `gorm:"default:false"`
	Layers          []Layer `gorm:"ForeignKey:UserID"`
	LayerPermission []*LayerPermission
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
	Layer       *Layer
	User        *User
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
	Name             string `gorm:"size:255;not null"`
	ServerID         uint
	UserID           uint
	URLS             []LayerURL `gorm:"ForeignKey:LayerID"`
	LayerPermissions []*LayerPermission
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
func (proxyServer *GISProxy) MigrateDatabase() {
	db, err := gorm.Open("postgres", proxyServer.DB.buildConnectionStr())
	if err != nil {
		proxyServer.logger.Error(err)
	}
	defer db.Close()
	db.AutoMigrate(&User{}, &ServerAuth{}, &LayerURL{}, &Layer{}, &Server{}, &LayerPermission{})
}

//GetDB return a Database Connection
func (proxyServer *GISProxy) GetDB() (db *gorm.DB, err error) {
	db, err = gorm.Open("postgres", proxyServer.DB.buildConnectionStr())
	if err != nil {
		proxyServer.logger.Error(err)
	}
	return
}

//LoadData load default data
func (proxyServer *GISProxy) LoadData() {
	db, err := proxyServer.GetDB()
	if err != nil {
		proxyServer.logger.Error(err)
	}
	defer db.Close()
	var user User
	db.FirstOrCreate(&user, User{UserName: "admin", FirstName: "Hisham", Password: "admin", LastName: "Karam", Email: "admin@admin.com"})
	layerServer := Server{
		URL:        "http://localhost:8080/geoserver",
		ServerType: "geoserver",
		Name:       "geoserver2",
		AuthInfo: ServerAuth{
			UserName: "admin",
			Password: "geoserver",
			Type:     "Basic",
		},
		Active: true,
		Layers: []Layer{
			Layer{
				Name: "geonode:other_healthcare_60cfefd3",
				LayerPermissions: []*LayerPermission{
					{User: &user},
				},
			},
		},
	}
	db.Create(&layerServer)
	// layerServer.Layers = append(layerServer.Layers, Layer{Name: "geonode:other_healthcare_60cfefd3"})
	// db.Save(&layerServer)
	var users []User
	var servers []Server
	db.Find(&users)
	db.Preload("Layers").Preload("Layers.LayerPermissions").Preload("Layers.LayerPermissions.User").Find(&servers)
	for _, user := range users {
		fmt.Println(user.UserName)
		// fmt.Println(user.Layers)
	}
	for _, server := range servers {
		fmt.Println(server.ServerType)
		fmt.Println(server.AuthInfo)
		fmt.Printf("^^^^^^^^^^^^^^%v\n", server.Layers[0].LayerPermissions[0])
	}
}
