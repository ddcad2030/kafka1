package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/IBM/sarama"
)

type Order struct {
	Customer string `json:"customer"`
	Coffee   string `json:"coffee"`
}

func main() {
	http.HandleFunc("/order", placeOrder)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func connectProducer(broker []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	return sarama.NewSyncProducer(broker, config)
}

func pushOrderToQueue(topic string, event []byte) error {
	broker := []string{"localhost:9092"}
	producer, err := connectProducer(broker)
	if err != nil {
		return err
	}

	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(event),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return err
	}
	log.Printf("Topic %s partition %d offset %d\n",
		topic, partition, offset,
	)
	return nil
}

func placeOrder(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	order := new(Order)
	if err := json.NewDecoder(r.Body).Decode(order); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	orderInBytes, err := json.Marshal(order)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = pushOrderToQueue("coffee_order", orderInBytes)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"msg":     order.Customer + " order sucessfully",
	}

	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println(err)
		http.Error(w, "Error placing order", http.StatusInternalServerError)
		return
	}

}
