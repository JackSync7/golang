package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"personal-web/connection"
	"strconv"
	"text/template"
	"time"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

type Blog struct {
	ID                                      int
	Title, Content, Author, PostDate, Image string
}

type Project struct {
	Title, StartDate, EndDate, Content, React, Python, Node, Golang, Selisih string
	Id                                                                       int
	Tech                                                                     []string
}

var dataBlog = []Blog{
	{
		Title:    "Dumbways Web App",
		Content:  "A web application (web app) is an application program that is stored on a remote server and delivered over the internet through a browser interface. Web services are web apps by definition and many, although not all, websites contain web apps. Developers design web applications for a wide variety of uses and users, from an organization to an individual for numerous reasons. Commonly used web applications can include webmail, online calculators or e-commerce shops. While users can only access some web apps by a specific browser, most are available no matter the browser.",
		Author:   "Jeri Utama",
		PostDate: "09 Maret 2023",
	},
	{
		Title:    "Dumbways Mobile Developer",
		Content:  "A mobile application, most commonly referred to as an app, is a type of application software designed to run on a mobile device, such as a smartphone or tablet computer. Mobile applications frequently serve to provide users with similar services to those accessed on PCs. Apps are generally small, individual software units with limited function. This use of app software was originally popularized by Apple Inc. and its App Store, which offers thousands of applications for the iPhone, iPad and iPod Touch.",
		Author:   "Jeri Utama",
		PostDate: "09 Maret 2023",
	},
}

var dataProject = []Project{
	{
		Title:     "Dumbways Web App",
		StartDate: "22-11-2023",
		EndDate:   "12-12-2023",
		Content:   "A web application (web app) is an application program that is stored on a remote server and delivered over the internet through a browser interface. Web services are web apps by definition and many, although not all, websites contain web apps. Developers design web applications for a wide variety of uses and users, from an organization to an individual for numerous reasons. Commonly used web applications can include webmail, online calculators or e-commerce shops. While users can only access some web apps by a specific browser, most are available no matter the browser.",
		React:     "<i class='fa-brands fa-react fa-xl me-3'></i>",
		Python:    "<i class='fa-brands fa-python fa-xl me-3'></i>",
		Node:      "",
		Golang:    "<i class='fa-brands fa-golang fa-xl me-3'></i>",
		Selisih:   "3 Month Ago",
		Tech:      []string{"<i class='fa-brands fa-react fa-xl me-3'></i>", "<i class='fa-brands fa-python fa-xl me-3'></i>", "Node"},
	},
	{
		Title:     "Dumbways Mobile Developer",
		StartDate: "12-11-2023",
		EndDate:   "05-22-2024",
		Content:   "A mobile application, most commonly referred to as an app, is a type of application software designed to run on a mobile device, such as a smartphone or tablet computer. Mobile applications frequently serve to provide users with similar services to those accessed on PCs. Apps are generally small, individual software units with limited function. This use of app software was originally popularized by Apple Inc. and its App Store, which offers thousands of applications for the iPhone, iPad and iPod Touch.",
		React:     "<i class='fa-brands fa-react fa-xl me-3'></i>",
		Python:    "<i class='fa-brands fa-python fa-xl me-3'></i>",
		Node:      "<i class='fa-brands fa-node-js fa-xl me-3'></i>",
		Golang:    "<i class='fa-brands fa-golang fa-xl me-3'></i>",
		Selisih:   "3 Weeks Ago",
		Tech:      []string{"<i class='fa-brands fa-python fa-xl me-3'></i>", "<i class='fa-brands fa-node-js fa-xl me-3'></i>", "<i class='fa-brands fa-golang fa-xl me-3'></i>"},
	},
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	connection.DatabaseConnect()
	e := echo.New()

	// route statis untuk mengakses folder public
	e.Static("/public", "public") // /public

	// renderer
	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

	e.Renderer = t

	// Routing
	e.GET("/hello", helloWorld)           //localhost:5000/hello
	e.GET("/", home)                      //localhost:5000
	e.GET("/contact", contact)            //localhost:5000/contact
	e.GET("/blog", blog)                  //localhost:5000/blog
	e.GET("/blog-detail/:id", blogDetail) //localhost:5000/blog-detail/0 | :id = url params
	e.GET("/form-blog", formAddBlog)      //localhost:5000/form-blog
	e.POST("/add-blog", addBlog)          //localhost:5000/add-blog
	e.GET("/delete-blog/:id", deleteBlog)
	e.GET("/myProject", myProject)
	e.POST("/addProject", addProject)
	e.GET("/deleteProject/:id", deleteProject)
	e.GET("/detailProject/:id", detailProject)
	e.GET("/editProject/:id", editProject)
	e.POST("/updateProject/:id", updateProject)

	fmt.Println("Server berjalan di port 5000")
	e.Logger.Fatal(e.Start("localhost:5000"))
}

func helloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World!")
}

func home(c echo.Context) error {
	// data, _ := connection.Conn.Query(context.Background(), "SELECT id, title, content, tech, start_date, end_date FROM tb_project")

	// var result []Project
	// for data.Next() {
	// 	var each = Project{}

	// 	err := data.Scan(&each.Id, &each.Title, &each.Content, &each.Tech, &each.StartDate, &each.EndDate)
	// 	if err != nil {
	// 		fmt.Println(err.Error())
	// 		return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
	// 	}

	// 	result = append(result, each)
	// }
	projects := map[string]interface{}{
		"Project": dataProject,
	}
	return c.Render(http.StatusOK, "index.html", projects)
}

func contact(c echo.Context) error {
	return c.Render(http.StatusOK, "contact.html", nil)
}

func blog(c echo.Context) error {
	data, _ := connection.Conn.Query(context.Background(), "SELECT id, title, content, image, post_date FROM tb_blog")

	var result []Blog
	for data.Next() {
		var each = Blog{}

		err := data.Scan(&each.ID, &each.Title, &each.Content, &each.Image, &each.PostDate)
		if err != nil {
			fmt.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
		}

		result = append(result, each)

	}

	blogs := map[string]interface{}{
		"Blogs": result,
	}
	return c.Render(http.StatusOK, "blog.html", blogs)
}

func blogDetail(c echo.Context) error {
	// http://localhost:5000/blog-detail/1
	// "1" => 1
	id, _ := strconv.Atoi(c.Param("id")) // url params | dikonversikan dari string menjadi int/integer

	var BlogDetail = Blog{}

	for i, data := range dataBlog {
		if id == i {
			BlogDetail = Blog{
				Title:    data.Title,
				Content:  data.Content,
				PostDate: data.PostDate,
				Author:   data.Author,
			}
		}
	}

	detailBlog := map[string]interface{}{
		"Blog": BlogDetail,
	}

	return c.Render(http.StatusOK, "blog-detail.html", detailBlog)
}

func formAddBlog(c echo.Context) error {
	return c.Render(http.StatusOK, "add-blog.html", nil)
}

func Pro(c echo.Context) error {
	return c.Render(http.StatusOK, "add-blog.html", nil)
}

func addBlog(c echo.Context) error {
	title := c.FormValue("inputTitle")
	content := c.FormValue("inputContent")

	println("Title: " + title)
	println("Content: " + content)

	var newBlog = Blog{
		Title:    title,
		Content:  content,
		Author:   "Dandi Saputra",
		PostDate: time.Now().String(),
	}

	dataBlog = append(dataBlog, newBlog)

	return c.Redirect(http.StatusMovedPermanently, "/blog")
}

func deleteBlog(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	dataBlog = append(dataBlog[:id], dataBlog[id+1:]...)

	return c.Redirect(http.StatusMovedPermanently, "/blog")
}

func myProject(c echo.Context) error {

	projects := map[string]interface{}{
		"Project": dataProject,
	}
	return c.Render(http.StatusOK, "myProject.html", projects)
}

func addProject(c echo.Context) error {
	title := c.FormValue("name")
	startDate := c.FormValue("startDate")
	endDate := c.FormValue("endDate")
	content := c.FormValue("textArea")
	react := c.FormValue("react")
	python := c.FormValue("python")
	node := c.FormValue("node")
	golang := c.FormValue("golang")

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
	println("awal: " + startDate)
	println("end: " + endDate)
	println("selisih: " + Distance)

	var newProject = Project{
		Title:     title,
		StartDate: startDate,
		EndDate:   endDate,
		Content:   content,
		React:     react,
		Python:    python,
		Node:      node,
		Golang:    golang,
		Selisih:   Distance,
	}
	dataProject = append(dataProject, newProject)
	return c.Redirect(http.StatusMovedPermanently, "/")
}

func deleteProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	dataProject = append(dataProject[:id], dataProject[id+1:]...)

	return c.Redirect(http.StatusMovedPermanently, "/")
}
func detailProject(c echo.Context) error {

	id, _ := strconv.Atoi(c.Param("id"))

	var ProjectDetail = Project{}

	for i, data := range dataProject {
		if id == i {
			ProjectDetail = Project{
				Title:     data.Title,
				StartDate: data.StartDate,
				EndDate:   data.EndDate,
				Content:   data.Content,
				React:     data.React,
				Python:    data.Python,
				Node:      data.Node,
				Golang:    data.Golang,
				Id:        id,
			}
		}

	}
	detailProject := map[string]interface{}{
		"Project": ProjectDetail,
	}
	return c.Render(http.StatusOK, "detailProject.html", detailProject)
}

func editProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var ProjectDetail = Project{}

	for i, data := range dataProject {
		if id == i {
			ProjectDetail = Project{
				Title:     data.Title,
				StartDate: data.StartDate,
				EndDate:   data.EndDate,
				Content:   data.Content,
				React:     data.React,
				Python:    data.Python,
				Node:      data.Node,
				Golang:    data.Golang,
				Id:        id,
			}
		}
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

	// layout := "02-01-2006"
	// t1, _ := time.Parse(layout, "12-04-2022")
	// t2, _ := time.Parse(layout, "12-04-2021")

	// diff := t1.Sub(t2)
	// selisih := 0
	// if int(diff.Hours() / 24 / 365) > 0 {
	// 	selisih = int(diff.Hours() / 24 / 365)
	// }else if int(diff.Hours() / 24) > 0

	var newProject = Project{
		Title:     title,
		StartDate: startDate,
		EndDate:   endDate,
		Content:   content,
		React:     react,
		Python:    python,
		Node:      node,
		Golang:    golang,
	}
	dataProject[id] = newProject
	return c.Redirect(http.StatusMovedPermanently, "/")
}
