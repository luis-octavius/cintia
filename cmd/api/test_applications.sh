TOKEN=$(cat /tmp/cintia_token 2>/dev/null)
JAVA_JOB=$(cat /tmp/job_id_java 2>/dev/null)
GO_JOB=$(cat /tmp/job_id_go 2>/dev/null)

echo Java: $JAVA_JOB
echo GO: $GO_JOB

echo "Creating applications..."

JAVA_APP_ID=$(curl -X POST http://localhost:8080/api/applications/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
     \"job_id\": \"$JAVA_JOB\",
     \"notes\": \"very good job, applicable\"
   }" | jq -r '.application.id'
 )

 echo Java APP ID: $JAVA_APP_ID

GO_APP_ID=$(curl -X POST http://localhost:8080/api/applications/ \
 -H "Content-Type: application/json" \
 -H "Authorization: Bearer $TOKEN" \
 -d "{
   \"job_id\": \"$GO_JOB\",
   \"notes\": \"better than java job, go is much better\"
 }" | jq '.application.id'
)
 
echo Go APP ID: $GO_APP_ID

echo "Listing applications..."

curl http://localhost:8080/api/applications/ \
  -H "Authorization: Bearer $TOKEN" | jq
   

echo "Testing getting details of an application"

curl http://localhost:8080/api/applications/$GO_APP_ID \
  -H "Authorization: Bearer $TOKEN" | jq 
curl http://localhost:8080/api/applications/$JAVA_APP_ID \
  -H "Authorization: Bearer $TOKEN" | jq

echo "Testing update application..."

curl -X PUT http://localhost:8080/api/applications/$GO_APP_ID \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "salary_offer": "1600 USD"
  }' | jq 

echo "Testing status update..."

curl -X PATCH http://localhost:8080/api/applications/$GO_APP_ID/status \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "status": "interviewing"
  }' | jq
