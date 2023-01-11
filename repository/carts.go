package repository

import (
	"a21hc3NpZ25tZW50/db"
	"a21hc3NpZ25tZW50/model"
	"encoding/json"
	"fmt"
)

type CartRepository struct {
	db db.DB
}

func NewCartRepository(db db.DB) CartRepository {
	return CartRepository{db}
}

func (u *CartRepository) ReadCart() ([]model.Cart, error) {
	records, err := u.db.Load("carts")
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("Cart not found!")
	}

	var cart []model.Cart
	err = json.Unmarshal([]byte(records), &cart)
	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (u *CartRepository) UpdateCart(cart model.Cart) error {
	listCart, _ := u.ReadCart()
	fmt.Println("ini list cart Sebelum : ", cart)
	listCartUpdate := []model.Cart{}
	for _, element := range listCart {
		if element.Name == cart.Name {
			listCartUpdate = append(listCartUpdate, cart)
			continue
		} else {
			listCartUpdate = append(listCartUpdate, element)
		}
	}

	jsonData, _ := json.Marshal(listCartUpdate)

	fmt.Println("ini list cart Update : ", jsonData)

	err := u.db.Save("carts", jsonData)
	if err != nil {
		return err
	}

	return nil
}

func (u *CartRepository) AddCart(cart model.Cart) error {

	card, _ := u.ReadCart()
	card = append(card, cart)
	jcard, _ := json.Marshal(card)
	err := u.db.Save("carts", jcard)
	if err != nil {
		return fmt.Errorf("Request Product Not Found")
	}
	fmt.Println(cart)
	return nil // TODO: replace this
}

func (u *CartRepository) ResetCarts() error {
	err := u.db.Reset("carts", []byte("[]"))
	if err != nil {
		return err
	}

	return nil
}

func (u *CartRepository) CartUserExist(name string) (model.Cart, error) {
	listcCart, err := u.ReadCart()
	if err != nil {
		return model.Cart{}, err
	}
	for _, element := range listcCart {
		if element.Name == name {
			return element, nil
		}
	}
	return model.Cart{}, fmt.Errorf("Cart Empty!")
}
