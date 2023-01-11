package api

import (
	"a21hc3NpZ25tZW50/model"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func (api *API) AddCart(w http.ResponseWriter, r *http.Request) {

	username := fmt.Sprintf("%s", r.Context().Value("username"))
	// w.Write([]byte(fmt.Sprintf("Welcome %s!", username)))

	fmt.Println(username)
	// if err := r.ParseForm(); err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.Write([]byte(`{"error":"Request Product Not Found"}`))
	// 	return
	// } else {
	// 	w.WriteHeader(200)
	// }
	r.ParseForm()
	// if r.FormValue("id") == "" {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.Write([]byte(`{"error":"Request Product Not Found"}`))
	// 	return
	// } else {
	// 	w.WriteHeader(200)
	// }
	fmt.Println(r.Form)
	if len(r.Form) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error":"Request Product Not Found"}`))
		return
	}
	// Check r.Form with key product, if not found then return response code 400 and message "Request Product Not Found".
	// TODO: answer here

	var list []model.Product
	for _, formList := range r.Form {
		for _, v := range formList {
			item := strings.Split(v, ",")
			p, _ := strconv.ParseFloat(item[2], 64)
			q, _ := strconv.ParseFloat(item[3], 64)
			total := p * q
			list = append(list, model.Product{
				Id:       item[0],
				Name:     item[1],
				Price:    item[2],
				Quantity: item[3],
				Total:    total,
			})
		}
	}
	totals := 0.0
	for _, list := range list {
		totals += list.Total
	}

	// Add data field Name, Cart and TotalPrice with struct model.Cart.
	cart := model.Cart{Name: username, Cart: list, TotalPrice: totals} // TODO: replace this

	_, err := api.cartsRepo.CartUserExist(cart.Name)
	if err != nil {
		api.cartsRepo.AddCart(cart)
	} else {
		api.cartsRepo.UpdateCart(cart)
	}
	api.dashboardView(w, r)

}
