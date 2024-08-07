# Example Kafka
Open three terminal follow step 
- make run-docker
- make run-producer
- make run-consumer
# JSON event
curl -X POST http://localhost:3000/order \
	-H "Content-Type: application/json" \
	-d '{
		"customer": "David",
		"coffee": "Americano"
		}'
