package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

//Model for course-file
type Course struct {
	CourseID    string  `json:"course_id"`
	CourseName  string  `json:"course_name"`
	CoursePrice int     `json:"course_price"`
	Author      *Author `json:"author"` //importing struct Author
}

type Author struct {
	FullName string `json:"full_name"`
	Website  string `json:"website"`
}

//fake DB
var courses []Course

//middleware
func isEmpty(c *Course) bool {
	if c.CourseName == "" || c.Author.FullName == "" || c.Author.Website == "" {
		return true
	}
	return false
}

func main() {
	fmt.Println("Go-Backend")
	r := mux.NewRouter()

	//Seeding the fake DB
	courses = append(courses, Course{CourseID: "1", CourseName: "Go", CoursePrice: 100, Author: &Author{FullName: "Golang", Website: "golang.com"}})
	courses = append(courses, Course{CourseID: "2", CourseName: "Python", CoursePrice: 200, Author: &Author{FullName: "Python", Website: "python.com"}})
	courses = append(courses, Course{CourseID: "3", CourseName: "Java", CoursePrice: 300, Author: &Author{FullName: "Java", Website: "java.com"}})

	//routing
	r.HandleFunc("/", home).Methods("GET")
	r.HandleFunc("/courses", getAllCourses).Methods("GET")
	r.HandleFunc("/course/{course_id}", getOneCourse).Methods("GET")
	r.HandleFunc("/course", createCourse).Methods("POST")
	r.HandleFunc("/course/{course_id}", updateOneCourse).Methods("PUT")
	r.HandleFunc("/course/{course_id}", deleteOneCourse).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}

//controllers -file

//serve home route

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Welcome to Go-Backend</h1>"))
}

func getAllCourses(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get all courses")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses) //encode the courses to bytes and send it to the client by writing it to the response writer (w) in JSON format (json.Encoder) and then close the connection (w.Close()) to avoid memory leaks and to avoid the client from sending more data.
}

func getOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get one course")
	w.Header().Set("Content-Type", "application/json") //setting the header to json format to send the response to the client in json format (w.Header().Set("Content-Type", "application/json")).
	// params := r.URL.Query()
	// courseID := params.Get("course_id")
	params := mux.Vars(r)
	courseID := params["course_id"]

	for _, course := range courses {
		if course.CourseID == courseID {
			json.NewEncoder(w).Encode(course)
			return
		}
	}
	json.NewEncoder(w).Encode("No Course found")
}

func createCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create course")
	w.Header().Set("Content-Type", "application/json")
	var course Course
	_ = json.NewDecoder(r.Body).Decode(&course) //decode the request body to the course struct (course) and then close the connection (r.Body.Close()) to avoid memory leaks and to avoid the client from sending more data.
	if isEmpty(&course) {
		json.NewEncoder(w).Encode("Please fill all the fields")
		return
	}

	//generate id
	rand.Seed(time.Now().UnixNano())
	course.CourseID = strconv.Itoa(rand.Intn(1000))

	courses = append(courses, course)
	json.NewEncoder(w).Encode(course)
}

func updateOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update one course")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	courseID := params["course_id"]

	for index, course := range courses {
		if course.CourseID == courseID {
			courses = append(courses[:index], courses[index+1:]...)
			_ = json.NewDecoder(r.Body).Decode(&course)
			course.CourseName = "Updated CourseName"
			courses = append(courses, course)
			json.NewEncoder(w).Encode(course)
			return
		}
	}

	json.NewEncoder(w).Encode("No Course found")
}

func deleteOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete one course")
	params := mux.Vars(r)
	for index, course := range courses {
		if course.CourseID == params["course_id"] {
			courses = append(courses[:index], courses[index+1:]...)
			json.NewEncoder(w).Encode(courses)
			return
		}
	}
	json.NewEncoder(w).Encode("Course ID wrong")
}
