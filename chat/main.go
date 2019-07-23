package main

import(
	"github.com/stretchr/gomniauth/providers/facebook"
        "github.com/stretchr/gomniauth/providers/github"
        "github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
	"log"
	"net/http"
	"text/template"
	"path/filepath"
	"sync"
	"fmt"
	"flag"
	"os"
	"trace"
)

type templateHandler struct{
	once sync.Once
	filename string
	templ *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter , r *http.Request){
	t.once.Do(func(){
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates" , t.filename)))
	})

	data := map[string]interface{}{
		"Host" : r.Host,
	}

	if authCookie , err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w,data)

}

func main(){

	var addr = flag.String("addr",":8080","the address of the application.")

	flag.Parse()

	gomniauth.SetSecurityKey("PUT YOUR AUTH KEY HERE")
        gomniauth.WithProviders(
                facebook.New("key", "secret","http://localhost:8080/auth/callback/facebook"),
                github.New("key", "secret","http://localhost:8080/auth/callback/github"),
                google.New("594120381762-de2cijq8o6kko28naj27eqjj177clpq4.apps.googleusercontent.com", "iaCWHmLH4CSogYtjJ_38NH4V","http://localhost:8080/auth/callback/google"),
	)

	r := NewRoom()

	r.tracer = trace.New(os.Stdout)

	http.Handle("/chat" ,MustAuth(&templateHandler{filename : "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room" , r)

	go r.run()

	fmt.Println("starting web server on :" , *addr)
	log.Println("starting web server and printing through log : " , *addr)

	if err := http.ListenAndServe(*addr , nil) ; err != nil {
		log.Fatal("ListenAndServe:" , err)
	} 
}
