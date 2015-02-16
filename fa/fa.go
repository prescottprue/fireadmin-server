package fa

import (
  // "html/template"
  "net/http"
  "encoding/json"
  "time"

  "appengine"
  "appengine/datastore"
  "appengine/user"
  // "github.com/melvinmt/firebase"
  "fmt"
  "strings"
)

type App struct {
  Secret string
  FbUrl string  
  Author string 
  Name string 
  Date time.Time 
}

type PersonName struct {
  First string
  Last  string
}

type Person struct {
  Name PersonName
}
// Setup routes
func init() {
  http.HandleFunc("/", root)
  http.HandleFunc("/secret", saveSecret)
  http.HandleFunc("/auth", generateAuth)

}
//Display Home Page
func root(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-type", "text/html; charset=utf-8")
  c := appengine.NewContext(r)
  u := user.Current(c)
  if u == nil {
    url, _ := user.LoginURL(c, "/")
    fmt.Fprintf(w, `<a href="%s">Sign in or register</a>`, url)
    return
  }
  url, _ := user.LogoutURL(c, "/")
  fmt.Fprintf(w, `Welcome, %s! (<a href="%s">sign out</a>)`, u, url)
}

/** Save App information to database
 */
func saveSecret(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  n := r.FormValue("name")
  if n == "" {
  	n = nameFromUrl(r.FormValue("fbUrl"))
  } 
  d := App{
    Secret: r.FormValue("secret"),
    FbUrl: r.FormValue("fbUrl"),
    Name: n,
    Date:    time.Now(),
  }
  //Add user if a user is currently logged in
  if u := user.Current(c); u != nil {
    d.Author = u.String()
  }
  //Query for already existing app with matching FbUrl
  q := datastore.NewQuery("App").Ancestor(appKey(c)).Filter("FbUrl =", d.FbUrl)
  var apps []App
  //Handle query error
  if _, err := q.GetAll(c, &apps); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
	}
	//Handle app existing error
	if l := len(apps); l >= 1 {
		http.Error(w, "App already exists", http.StatusInternalServerError)
    return
	}
	//Create new key for app
  key := datastore.NewIncompleteKey(c, "App", appKey(c))
  //Put new app in database
  _, err := datastore.Put(c, key, &d)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
}
func generateAuth(w http.ResponseWriter, r *http.Request) {
  d := App{
    Secret: r.FormValue("secret"),
    FbUrl: r.FormValue("fbUrl"),
    Date:    time.Now(),
  }
	c := appengine.NewContext(r)
  //Query for already existing app with matching FbUrl
  q := datastore.NewQuery("App").Ancestor(appKey(c)).Filter("FbUrl =", d.FbUrl)
  var apps []App
  //Handle query error
  if _, err := q.GetAll(c, &apps); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
	}
	//Handle app existing error
	if l := len(apps); l >= 1 {
		http.Error(w, "App already exists", http.StatusInternalServerError)
    return
	}
  if _, err := q.GetAll(c, &apps); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
	}
  js, _ := json.Marshal(apps)
   w.Header().Set("Content-Type", "application/json")
  w.Write(js)

    // var err error

    // // url := "https://SampleChat.firebaseio.com/users/fred/name"

    // // Can also be your Firebase secret:
    // // authToken := "MqL0c8tKCtheLSYcygYNtGhU8Z2hULOFs9OKPdEp"

    // // Auth is optional:
    // ref := firebase.NewReference(g.FbUrl).Auth(g.Secret)

    // // Create the value.
    // personName := PersonName{
    //     First: "Fred",
    //     Last:  "Swanson",
    // }

    // // Write the value to Firebase.
    // if err = ref.Write(personName); err != nil {
    //     panic(err)
    // }

    // // Now, we're going to retrieve the person.
    // personUrl := "https://SampleChat.firebaseIO-demo.com/users/fred"

    // personRef := firebase.NewReference(personUrl).Export(false)

    // fred := Person{}

    // if err = personRef.Value(fred); err != nil {
    //     panic(err)
    // }

    // fmt.Println(fred.Name.First, fred.Name.Last) // prints: Fred Swanson
}
func nameFromUrl(u string) string {
	u1 := strings.Replace(u, "https://", "", -1)
	return strings.Replace(u1, ".firebaseio.com", "", -1)
}
//Create a new datastore key for an app
func appKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "App", "default_app", 0, nil)
}
// func page(w http.ResponseWriter, r *http.Request) {
//   if err := homeTemplate.Execute(w, "Welcome to the Fireadmin Server"); err != nil {
//     http.Error(w, err.Error(), http.StatusInternalServerError)
//   }
// }
// var homeTemplate = template.Must(template.New("book").Parse(`
// <html>
//   <head>
//     <title>Fireadmin Server</title>
//   </head>
//   <body>
//   <div style="text-align:center; margin-top:20%;">{{.}}</div>
    
//   </body>
// </html>
// `))
