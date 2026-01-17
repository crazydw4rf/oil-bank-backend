POST {{BASE_URL}}/transactions
{
  "email": "john.seller@example.com",
  "oil_volume": 100,
  "price": 23000,
  "transaction_type": "SELL"
}
HTTP 201

POST {{BASE_URL}}/transactions
{
  "email": "contact@techcorp.com",
  "oil_volume": 20,
  "price": 32000,
  "transaction_type": "BUY"
}
HTTP 201
