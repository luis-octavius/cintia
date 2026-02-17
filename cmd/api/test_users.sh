printf "Beginning Tests...\n"
printf "\nRegistering user...\n"

curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "cintia",
    "email": "cintia@gmail.com",
    "password": "amoluis123"
  }' | jq 

printf "\nLogin with registed user and saving generated token...\n"

TOKEN=$(curl -X POST http://localhost:8080/api/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "cintia@gmail.com",
    "password": "amoluis123"
  }' | jq -r '.token')

printf "\nUpdating name to test authorization...\n"

curl -X PUT http://localhost:8080/api/users/me \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name":"cintia_deusa"}' | jq

echo "$TOKEN" > /tmp/cintia_token 
echo "Token saved: $TOKEN"
