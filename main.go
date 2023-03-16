package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"personal-web/connection"
	"personal-web/middleware"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type Template struct {
	templates *template.Template
}

type User struct {
	ID                    int
	Name, Email, Password string
}

type Project struct {
	Title, Content, React, Python, Node, Golang, Duration, Waktu, Author, Image string
	StartDate, EndDate                                                          time.Time
	Id                                                                          int
	Tech                                                                        []string
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {

	connection.DatabaseConnect()
	e := echo.New()

	// route statis untuk mengakses folder public
	e.Static("/public", "public") // /public
	e.Static("/upload", "upload")

	e.Use(session.Middleware(sessions.NewCookieStore([]byte("session"))))

	// renderer
	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

	e.Renderer = t

	// Routing
	e.GET("/", home)
	e.GET("/register", formRegister)
	e.POST("/addRegister", addRegister)
	e.GET("/login", formLogin)
	e.POST("/runLogin", runLogin)
	e.GET("/contact", contact)
	e.GET("/myProject", myProject)
	e.POST("/addProject", middleware.UploadFile(addProject))
	e.GET("/deleteProject/:id", deleteProject)
	e.GET("/detailProject/:id", detailProject)
	e.GET("/editProject/:id", editProject)
	e.POST("/updateProject/:id", middleware.UploadFile(updateProject))
	e.GET("/logout", logout)

	fmt.Println("Server berjalan di port 5000")
	e.Logger.Fatal(e.Start("localhost:5000"))
}

func formLogin(c echo.Context) error {
	sess, _ := session.Get("session", c)
	flash := map[string]interface{}{
		"FlashStatus":  sess.Values["isLogin"],
		"FlashMessage": sess.Values["message"],
		"FlashName":    sess.Values["name"],
	}

	sess.Save(c.Request(), c.Response())
	delete(sess.Values, "message")
	delete(sess.Values, "status")
	tmpl, err := template.ParseFiles("views/login.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}
	return tmpl.Execute(c.Response(), flash)
}

func formRegister(c echo.Context) error {

	tmpl, err := template.ParseFiles("views/register.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}
	return tmpl.Execute(c.Response(), nil)
}

func addRegister(c echo.Context) error {
	err := c.Request().ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_user (name, email, password) VALUES ($1, $2, $3)", name, email, passwordHash)
	if err != nil {
		redirectWithMessage(c, "Register failed, please try again", false, "/register")
	}
	return redirectWithMessage(c, "Register Success, please try again", true, "/login")
}

func home(c echo.Context) error {
	sess, _ := session.Get("session", c)
	flash := map[string]interface{}{
		"FlashStatus":  sess.Values["isLogin"],
		"FlashMessage": sess.Values["message"],
		"FlashName":    sess.Values["name"],
	}
	delete(sess.Values, "message")

	sess.Save(c.Request(), c.Response())
	data, _ := connection.Conn.Query(context.Background(), "SELECT tb_project.id, tb_project.title, tb_project.content, tb_project.tech, tb_project.start_date, tb_project.end_date, tb_project.duration, tb_user.name, tb_project.image FROM tb_project inner join tb_user ON tb_project.author_id = tb_user.id")

	var result []Project
	for data.Next() {
		var each = Project{}
		err := data.Scan(&each.Id, &each.Title, &each.Content, &each.Tech, &each.StartDate, &each.EndDate, &each.Duration, &each.Author, &each.Image)
		if err != nil {
			fmt.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
		}
		result = append(result, each)
	}
	projects := map[string]interface{}{
		"Project": result,
		"Flash":   flash,
	}
	return c.Render(http.StatusOK, "index.html", projects)
}

func contact(c echo.Context) error {
	sess, _ := session.Get("session", c)
	flash := map[string]interface{}{
		"FlashStatus":  sess.Values["isLogin"],
		"FlashMessage": sess.Values["message"],
		"FlashName":    sess.Values["name"],
	}
	delete(sess.Values, "message")

	sess.Save(c.Request(), c.Response())
	return c.Render(http.StatusOK, "contact.html", flash)
}

func myProject(c echo.Context) error {
	sess, _ := session.Get("session", c)
	flash := map[string]interface{}{
		"FlashStatus":  sess.Values["status"],
		"FlashMessage": sess.Values["message"],
		"FlashName":    sess.Values["name"],
	}
	delete(sess.Values, "message")
	// delete(sess.values, "status")
	sess.Save(c.Request(), c.Response())
	return c.Render(http.StatusOK, "myProject.html", flash)
}

func addProject(c echo.Context) error {
	sess, _ := session.Get("session", c)
	author := sess.Values["id"]

	delete(sess.Values, "message")
	sess.Save(c.Request(), c.Response())
	title := c.FormValue("name")
	startDate := c.FormValue("startDate")
	endDate := c.FormValue("endDate")
	content := c.FormValue("textArea")
	react := c.FormValue("react")
	python := c.FormValue("python")
	node := c.FormValue("node")
	golang := c.FormValue("golang")

	var techs [4]string
	techs[0] = react
	techs[1] = python
	techs[2] = node
	techs[3] = golang

	image := c.Get("dataFile").(string)

	layout := "2006-01-02"
	t1, _ := time.Parse(layout, endDate)
	t2, _ := time.Parse(layout, startDate)

	diff := t1.Sub(t2)

	days := int(diff.Hours() / 24)
	months := int(diff.Hours() / 24 / 30)
	weeks := int(diff.Hours() / 24 / 7)
	years := int(diff.Hours() / 24 / 365)

	var Distance string
	if years > 0 {

		Distance = strconv.Itoa(years) + " Years Ago"
		fmt.Printf("ini tahun : %s --", Distance)
	} else if months > 0 {

		Distance = strconv.Itoa(months) + " Month Ago"
		fmt.Printf("ini bulan : %s --", Distance)
	} else if weeks > 0 {

		Distance = strconv.Itoa(weeks) + " Weeks Ago"
		fmt.Printf("ini minggu : %s --", Distance)
	} else if days > 0 {

		Distance = strconv.Itoa(days) + " Days Ago"
		fmt.Printf("ini hari : %s --", Distance)
	}

	_, err := connection.Conn.Exec(context.Background(), "INSERT INTO tb_project (title, content, start_date, end_date, tech, duration, author_id, image) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", title, content, t2, t1, techs, Distance, author, image)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func deleteProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_project WHERE id = $1", id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func detailProject(c echo.Context) error {

	id, _ := strconv.Atoi(c.Param("id"))

	sess, _ := session.Get("session", c)
	flash := map[string]interface{}{
		"FlashStatus":  sess.Values["isLogin"],
		"FlashMessage": sess.Values["message"],
		"FlashName":    sess.Values["name"],
	}
	delete(sess.Values, "message")

	sess.Save(c.Request(), c.Response())

	var ProjectDetail = Project{}
	err := connection.Conn.QueryRow(context.Background(), "SELECT tb_project.id, tb_project.title, tb_project.content, tb_project.tech, tb_project.start_date, tb_project.end_date, tb_project.duration, tb_user.name, tb_project.image FROM tb_project inner join tb_user ON tb_project.author_id = tb_user.id where tb_project.id = $1", id).Scan(&ProjectDetail.Id, &ProjectDetail.Title, &ProjectDetail.Content, &ProjectDetail.Tech, &ProjectDetail.StartDate, &ProjectDetail.EndDate, &ProjectDetail.Duration, &ProjectDetail.Author, &ProjectDetail.Image)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
	}

	ProjectDetail.Waktu = ProjectDetail.StartDate.Format("2006-01-02")

	detailProject := map[string]interface{}{
		"Project": ProjectDetail,
		"Flash":   flash,
	}
	return c.Render(http.StatusOK, "detailProject.html", detailProject)
}

func editProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var ProjectDetail = Project{}
	err := connection.Conn.QueryRow(context.Background(), "SELECT id, title, content, start_date, end_date from tb_project WHERE id = $1", id).Scan(&ProjectDetail.Id, &ProjectDetail.Title, &ProjectDetail.Content, &ProjectDetail.StartDate, &ProjectDetail.EndDate)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
	}

	detailProject := map[string]interface{}{
		"Project": ProjectDetail,
	}
	return c.Render(http.StatusOK, "editProject.html", detailProject)
}

func updateProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	title := c.FormValue("name")
	startDate := c.FormValue("startDate")
	endDate := c.FormValue("endDate")
	content := c.FormValue("textArea")
	react := c.FormValue("react")
	python := c.FormValue("python")
	node := c.FormValue("node")
	golang := c.FormValue("golang")

	image := c.Get("dataFile").(string)
	var techs [4]string
	techs[0] = react
	techs[1] = python
	techs[2] = node
	techs[3] = golang

	fmt.Print(techs)

	layout := "2006-01-02"
	t1, _ := time.Parse(layout, endDate)
	t2, _ := time.Parse(layout, startDate)

	diff := t1.Sub(t2)

	days := int(diff.Hours() / 24)
	months := int(diff.Hours() / 24 / 30)
	weeks := int(diff.Hours() / 24 / 7)
	years := int(diff.Hours() / 24 / 365)

	var Distance string
	if years > 0 {

		Distance = strconv.Itoa(years) + " Years Ago"
		fmt.Printf("ini tahun : %s --", Distance)
	} else if months > 0 {

		Distance = strconv.Itoa(months) + " Month Ago"
		fmt.Printf("ini bulan : %s --", Distance)
	} else if weeks > 0 {

		Distance = strconv.Itoa(weeks) + " Weeks Ago"
		fmt.Printf("ini minggu : %s --", Distance)
	} else if days > 0 {

		Distance = strconv.Itoa(days) + " Days Ago"
		fmt.Printf("ini hari : %s --", Distance)
	}

	_, err := connection.Conn.Exec(context.Background(), "UPDATE tb_project SET title = $1, content = $2, start_date = $3, end_date = $4, tech = $5, duration = $6, image = $7 WHERE id = $8 ", title, content, t2, t1, techs, Distance, image, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func runLogin(c echo.Context) error {
	err := c.Request().ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	email := c.FormValue("email")
	password := c.FormValue("password")

	user := User{}
	err = connection.Conn.QueryRow(context.Background(), "SELECT * FROM tb_user WHERE email = $1", email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err != nil {
		return redirectWithMessage(c, "Email salah", false, "/login")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return redirectWithMessage(c, "Password Salah!", false, "/login")
	}
	sess, _ := session.Get("session", c)
	sess.Options.MaxAge = 10800
	sess.Values["mesasge"] = "Login Success"
	sess.Values["status"] = true
	sess.Values["name"] = user.Name
	sess.Values["id"] = user.ID
	sess.Values["isLogin"] = true
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func redirectWithMessage(c echo.Context, message string, status bool, path string) error {
	sess, _ := session.Get("session", c)
	sess.Values["message"] = message
	sess.Values["status"] = status
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusMovedPermanently, path)
}

func logout(c echo.Context) error {
	sess, _ := session.Get("session", c)
	sess.Options.MaxAge = -1
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusTemporaryRedirect, "/")
}
