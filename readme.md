curl -X GET "http://localhost:50055/alive" -w "\n"
curl -X POST "http://localhost:50055/add" -H "Content-Type: text/plain" -d "5" -w "\n"
curl -X POST "http://localhost:50055/work" -H "Content-Type: text/plain" -d "job5,job6,job7,job8" -w "\n"
curl -X DELETE "http://localhost:50055/delete" -H "Content-Type: text/plain" -d "1,2,3" -w "\n"
