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
  http.HandleFunc("/upload", generateAuth)
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
  	n = nameFromUrl(r.FormValue("fbUrl")) // Get name from url if it does not exist
  }
  //Build app object
  d := App{
    Secret: r.FormValue("secret"),
    FbUrl: r.FormValue("fbUrl"),
    Name: n,
    Date: time.Now(),
  }
  //Check if user is logged in
  if u := user.Current(c); u != nil {
    d.Author = u.String() //Add user as Author param
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
  //Set form val as FbUrl of App
  d := App{
    FbUrl: r.FormValue("fbUrl"),
  }
	c := appengine.NewContext(r) //App Engine Context

	la, err := GetApp(d, c) //Load App including secret from database
  if err != nil {
  	http.Error(w, err.Error(), http.StatusInternalServerError)
  	return
  }

	//[TODO] Load shape of auth object from database (set by user)
	//[TODO] Run required authentication check (password)
  //[TODO] Load auth data from db
	gen := fireauth.New(la.Secret) // TokenGenerator
  //[TODO] Fill auth object with actual data
	data := fireauth.Data{"uid": "1"} // Build auth object
	token, err := gen.CreateToken(data, nil) //Create token
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//Build Token response
	ts := TokenRes{
		Token: token,
	}
	res, err := json.Marshal(ts) //Marshal to Json for response
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
  w.Header().Set("Content-Type", "application/json") 
  w.Write(res) //Write response
}
//----------------Util Funcitons-----------------\\
//Get app name from Firebase url
func nameFromUrl(u string) string {
	u1 := strings.Replace(u, "https://", "", -1)
	return strings.Replace(u1, ".firebaseio.com", "", -1)
}
//Create a new datastore key for an app
func appKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "App", "default_app", 0, nil)
}