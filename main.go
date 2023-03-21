package main

import (
	h "nextlaundry_apis/helper"
	s "nextlaundry_apis/models/setup"
	r "nextlaundry_apis/routes"
)

//TODO: CREATE DASHBOARD FUNCTION UNTUK RETURN JUMLAH ANGKA TRANSAKSI DAN KEDATANGAN MEMBER BARU
//TODO: TESTING DAN INTEGRASI DI NEXTLAUNDRY REACT

func main() {
	//Init Database Setups
	s.ConnectDatabase()
	s.Migrations()

	//Init Clean Heap
	go h.CleanTokenHeap()

	//Init Router
	router := r.InitRouter()
	router.Run(":8000")
}
