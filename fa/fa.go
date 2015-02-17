package fa

import (
  // "html/template"
  "net/http"
  "encoding/json"
  "time"

  "appengine"
  "appengine/datastore"
  "appengine/user"
  // "appengine/log"
  // "github.com/melvinmt/firebase"
  "fmt"
  "strings"
  "errors"
  "github.com/zabawaba99/fireauth"
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
type TokenRes struct {
	Token string
}
// Setup routes
func init() {
  http.HandleFunc("/", root)
  http.HandleFunc("/setup", saveSecret)
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
//Take url and find if app exists
func GetApp(ad App, c appengine.Context) (a App, err error) {

  //Query for already existing app with matching FbUrl
  q := datastore.NewQuery("App").Ancestor(appKey(c)).Filter("FbUrl =", ad.FbUrl).Limit(1)
  var apps []App
  //Handle query error
  keys, err := q.GetAll(c, &apps)
  var n App
  if err != nil {
  	return n, err
  }
  c.Infof("Keys: %v", keys)
  c.Infof("Apps: %v", apps)

  if len(keys) < 1 {
		return n, errors.New("App Not Found") 
	}
	var app App
  app = apps[0]

	//Handle app existing error
	return app, nil
}
func generateAuth(w http.ResponseWriter, r *http.Request) {
  d := App{
    FbUrl: r.FormValue("fbUrl"),
  }
  	//App Engine Context
	c := appengine.NewContext(r)
	//Load App including secret from database
	a, err := GetApp(d, c)
  if err != nil {
  	http.Error(w, err.Error(), http.StatusInternalServerError)
  	return
  }
	//[TODO] Load shape of auth object from database (set by user)
	// Run required authentication check (password)
	// TokenGenerator
	gen := fireauth.New(a.Secret)
	// Build auth object
	data := fireauth.Data{"provider":"Fireadmin", "uid": "1"}
	token, err := gen.CreateToken(data, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//Token response
	ts := TokenRes{
		Token: token,
	}
	res, err := json.Marshal(ts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
  //Write response
  w.Header().Set("Content-Type", "application/json")
  w.Write(res)

    // var err error


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
