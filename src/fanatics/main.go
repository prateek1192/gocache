package main

import (
    "bytes"
    "encoding/json"
    "github.com/gorilla/mux"
    "log"
    "net/http"
    "os"
    "sync"
    "time"
    "fcache"
)

type Certificate struct {
    ID		string
    StartTime   time.Time
    Expired	bool
}
var certs []Certificate
var cache fcache.Cache
var cache_exp_time time.Duration
func main() {
	log.Printf("Server starting at localhost port 8000")
	router := mux.NewRouter()
	// Cache expiration time is hardcoded but have written a method
	// to demonstrate how it can be changed
	set_expiration_time(30)
	cache = fcache.New(cache_exp_time)
	router.HandleFunc("/certs", GetAllCerts).Methods("GET")
	router.HandleFunc("/cert/{domain}", GetCert).Methods("GET")
	go requestSelfCert()
	log.Fatal(http.ListenAndServe(":8000", router))
}

func requestSelfCert() {
	time.Sleep(20 * time.Second)
	_, err := http.Get("http://localhost:8000/cert/start-cert123")
	// Since the certificate is always automatically renewed 
	// the service runs
	// If there was no certificate the service would have exited.
	if err != nil {
                log.Printf(err.Error())
                os.Exit(1)
        }

}

func set_expiration_time(time_in_seconds int){
	cache_exp_time = time.Duration(time_in_seconds) * time.Second
}

func GetAllCerts(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request to get all certificates received")
	json.NewEncoder(w).Encode(certs)
}

func GetCert(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request for certificate received")
	var m sync.Mutex
	m.Lock()
	time.Sleep(10*time.Second)
	var buffer bytes.Buffer
	params := mux.Vars(r)
	dom := params["domain"]
	present := false
	if cache.Get(dom) != nil {
		log.Printf("String found in Cache")
		buffer.WriteString("foo-")
		buffer.WriteString(dom)
		json.NewEncoder(w).Encode(buffer.String())
	} else {
		for _, cert := range certs {
			if cert.ID == dom {
				log.Printf("Certificate found.")
				// Certificate already present, will return
				// Refresh the StartTime of the certificate
				// Write method to make certificates expired
				present = true
				buffer.WriteString("foo-")
				buffer.WriteString(cert.ID)
				cache.Set(cert.ID, []byte(cert.ID))
				json.NewEncoder(w).Encode(buffer.String())
			}
		}
		if present == false {
			log.Printf("Certificate not found. Will create new.")
			certs = append(certs, Certificate{ID: dom, StartTime: time.Now(), Expired: false})
			buffer.WriteString("foo-")
			id := (certs[len(certs)-1]).ID
			buffer.WriteString(id)
			cache.Set(id, []byte(id))
			json.NewEncoder(w).Encode(buffer.String())
		}
	}
	m.Unlock()
}
