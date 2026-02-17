TOKEN=$(cat /tmp/cintia_token 2>/dev/null)

printf "\nMake sure to run test_users.sh before this test\n"

if [ -z "$TOKEN" ]; then 
  echo "Error: No token found. Run test_users.sh first"
  exit 1 
fi 


printf "\nCreating a Job and getting his ID...\n"

JOB_ID_GO=$(curl -X POST http://localhost:8080/api/jobs/ \
  -H "Content-Type: application/json" \
  -d '{
    "title":"Backend Developer Go",
    "company":"Google",
    "location":"remote",
    "description":"Very good job",
    "salary_range":"1200 USD",
    "requirements":"Java, Go, Kubernetes, Docker",
    "link":"http://linkedin.com/google-developer-go",
    "source":"linkedin"
  }' \
  -H "Authorization: Bearer $TOKEN" | jq -r '.job.id')


JOB_ID_JAVA=$(curl -X POST http://localhost:8080/api/jobs/ \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Senior Java Backend Developer", 
    "company":"Mercado Livre",
    "location":"on site",
    "description":"Other very good job",
    "salary_range": "2500 USD",
    "requirements": "Java, DDD, Docker, Kubernetes, Spring, Kafka, AWS",
    "link": "http://linkedin.com/senior-java-developer",
    "source": "linkedin"
  }' \
  -H "Authorization: Bearer $TOKEN" | jq -r '.job.id')

echo "Java: $JOB_ID_JAVA"

printf "\nGetting job from ID...\n"

curl http://localhost:8080/api/jobs/$JOB_ID_GO | jq 
curl http://localhost:8080/api/jobs/$JOB_ID_JAVA | jq 

printf "\nTesting search queries...\n"

printf "\nTesting with company = 'mercado'"
curl 'http://localhost:8080/api/jobs/?company=mercado' | jq 

printf "\nTesting with requirements = 'kubernetes'"
curl 'http://localhost:8080/api/jobs/?requirements=kubernetes' | jq

printf "\nChanging job to inactive...\n"

curl -X PATCH http://localhost:8080/api/jobs/$JOB_ID_GO \
  -H "Authorization: Bearer $TOKEN" | jq 

curl -X PATCH http://localhost:8080/api/jobs/$JOB_ID_JAVA \
  -H "Authorization: Bearer $TOKEN" | jq 



