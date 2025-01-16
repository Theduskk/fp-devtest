Requires Go 1.23.4 (Requires 1.20+); Docker Enginer 27.4.0; (Tested on these, may work on older versions)

go build
docker build --tag devtask0.1 .
docker run --publish 3000:3000 devtask0.1

curl --request GET \           
-H 'Accept: application/json' \
--url 'http://localhost:3000/slots?duration=30&continuous=false'

curl --request GET \           
-H 'Accept: application/json' \
--url 'http://localhost:3000/slots?duration=30&continuous=true'