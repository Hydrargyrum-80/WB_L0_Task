package app

import (
	"WB_L0_Task/models"
	"WB_L0_Task/storage"
	"encoding/json"
	"fmt"
	"github.com/nats-io/stan.go"
	"io"
	"log"
	"net/http"
)

var st storage.IStorage

func getDataHandler(w http.ResponseWriter, r *http.Request) {
	ids, ok := r.URL.Query()["id"]
	if !ok || len(ids[0]) < 1 {
		http.Error(w, "Missing 'id' parameter", http.StatusBadRequest)
		return
	}
	id := ids[0]
	order := st.GetOrder(id)
	if order == nil {
		http.Error(w, "Order not found", http.StatusNotFound)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func logCloser(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Printf("close error: %s", err)
	}
}

func Start() {
	st = &storage.OrdersStorage{}
	st.Init()
	conn, err := stan.Connect("test-cluster", "test-client")
	if err != nil {
		log.Print(err)
		return
	}
	defer logCloser(conn)

	handle := func(msg *stan.Msg) {
		var order models.Order
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			fmt.Println(err)
		}
		st.AddOrder(order, order.OrderUid)
		log.Print(order)
	}

	sub, err := conn.Subscribe(
		"stream-name",
		handle,
		stan.DeliverAllAvailable())
	if err != nil {
		log.Print(err)
		return
	}
	defer logCloser(sub)
	http.HandleFunc("/getData", getDataHandler)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	http.ListenAndServe(":8080", nil)
}
