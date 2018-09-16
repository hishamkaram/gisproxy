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
	URL        string `sql:"type:text;not null"`
	ServerType string `gorm:"type:varchar(100);unique;not null"`
	Active     bool   `gorm:"default:true;not null"`
	AuthInfo   ServerAuth
}

//User gisproxy user model
type User struct {
	gorm.Model
	FirstName   string
	LastName    string
	UserName    string `gorm:"size:255;unique;not null"`
	Password    string `sql:"type:text;not null"`
	Email       string `gorm:"type:varchar(100);unique_index;not null"`
	IsSuperUser bool   `gorm:"default:false"`
	IsStaff     bool   `gorm:"default:false"`
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
	db.AutoMigrate(&User{}, &ServerAuth{}, &Server{})
}

//LoadData load default data
func (server *GISProxy) LoadData() {
	db, err := gorm.Open("postgres", server.DB.buildConnectionStr())
	if err != nil {
		server.logger.Error(err)
	}
	defer db.Close()
	var user User
	db.FirstOrCreate(&user, User{UserName: "admin", FirstName: "Hisham", Password: "admin", LastName: "Karam", Email: "admin@admin.com"})
}
