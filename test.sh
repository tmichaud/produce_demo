echo "Fetching"
echo "curl http://127.0.0.1:8080/produce"
curl http://127.0.0.1:8080/produce
echo ""

echo "Adding"
echo "curl -d '[ {"Produce Code": "AAAA-1111-2222-3333", "Name": "Pizza Pie", "Unit Price": "200.6" } ]' -X POST http://127.0.0.1:8080/produce"
curl -d '[ {"Produce Code": "AAAA-1111-2222-3333", "Name": "Pizza Pie", "Unit Price": "200.6" } ]' -X POST http://127.0.0.1:8080/produce 
echo ""

echo "Adding - ERROR"
echo " curl -d '[ {\"Produce Code\": \"-AAAA-1111-2222-3333\", \"Name\": \" Pizza Pie\", \"Unit Price\": \"2a0.6\" }, {\"Produce Code\": \"AAAA-1111-2222-4444\", \"Name\": \"Pizza Pie\", \"Unit Price\": \"200.6\" } ]' -X POST http://127.0.0.1:8080/produce "
curl -d '[ {"Produce Code": "-AAAA-1111-2222-3333", "Name": " Pizza Pie", "Unit Price": "2a0.6" }, {"Produce Code": "AAAA-1111-2222-4444", "Name": "Pizza Pie", "Unit Price": "200.6" } ]' -X POST http://127.0.0.1:8080/produce 
echo ""

echo "Fetching"
echo "curl http://127.0.0.1:8080/produce/AAAA-1111-2222-3333"
curl http://127.0.0.1:8080/produce/AAAA-1111-2222-3333
echo ""

echo "Deleting"
echo "curl -X "DELETE" http://127.0.0.1:8080/produce/AAAA-1111-2222-3333"
curl -X "DELETE" http://127.0.0.1:8080/produce/AAAA-1111-2222-3333
echo ""

echo "Deleting - ERROR"
echo "curl -X "DELETE" http://127.0.0.1:8080/produce/-AAAA-1111-2222-3333"
curl -X "DELETE" http://127.0.0.1:8080/produce/-AAAA-1111-2222-3333
echo ""

echo "Deleting - ERROR"
echo "curl -X "DELETE" http://127.0.0.1:8080/produce/AAAA-1111-2222-9999"
curl -X "DELETE" http://127.0.0.1:8080/produce/AAAA-1111-2222-9999
echo ""

echo "Fetching"
echo "curl http://127.0.0.1:8080/produce"
curl http://127.0.0.1:8080/produce
echo ""

