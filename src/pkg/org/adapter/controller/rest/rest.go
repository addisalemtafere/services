package rest

import (
	"auth/src/pkg/org/usecase"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Error struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (err Error) Error() string {
	return err.Message
}

type Controller struct {
	log        *log.Logger
	interactor usecase.Interactor
	sm         *http.ServeMux
}

// Organization
type Organization struct {
	Id       string    `json:"id"`
	Name     string    `json:"name"`
	Capital  float64   `json:"capital"`
	RegDate  time.Time `json:"reg_date"`
	Country  string    `json:"country"`
	Logo     string    `json:"logo"`
	Category *struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"category"`
	LegalCondtion *struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"legal_condition"`
	Taxes       []struct{}   `json:"taxes"`
	Departments []Department `json:"departments"`
	Details     interface{}  `json:"details"`
	CreatedAt   time.Time    `json:"created_at"`
}

// Department
type Department struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Categories  []struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"categories"`
	CountryWhitelist []string    `json:"country_whitelist"`
	CountryBlacklist []string    `json:"country_blacklist"`
	Details          interface{} `json:"details"`
	CreatedAt        time.Time   `json:"created_at"`
}

// Eth Bus Org
type EthBusOrg struct {
	TIN     string `json:"tin"`
	TinFile string `json:"tin_file"`
	Status  struct {
		Verified bool   `json:"verified"`
		Status   string `json:"status"`
		Message  string `json:"message"`
	} `json:"status"`
	RegNo   string `json:"reg_no"`
	RegFile string `json:"reg_file"`
}

// Response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   error       `json:"error,omitempty"`
}

func New(log *log.Logger, interactor usecase.Interactor, sm *http.ServeMux) Controller {
	var controller = Controller{log: log, interactor: interactor}

	// [TODO] Routing
	sm.HandleFunc("/org/check-tin", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			{
				controller.CheckTIN(w, r)
			}
		}
	})

	sm.HandleFunc("/orgs/categories", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			{
				controller.GetCategories(w, r)
			}
		case http.MethodPost:
			{
				controller.GetAddCategory(w, r)
			}
		}
	})

	sm.HandleFunc("/orgs/legal-conditions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			{
				controller.GetLegalConditions(w, r)
			}
		case http.MethodPost:
			{
				controller.GetAddLegalCondition(w, r)
			}
		}
	})

	sm.HandleFunc("/orgs/init", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			{
				controller.GetInitOrgRegistration(w, r)
			}
		}
	})

	// Taxes
	sm.HandleFunc("/orgs/taxes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				controller.GetAddTax(w, r)
			}
		case http.MethodGet:
			{
				controller.GetTaxes(w, r)
			}
		}
	})

	controller.sm = sm

	return controller
}

func SendJSONResponse(w http.ResponseWriter, data interface{}, status int) {
	serData, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(serData)
}
