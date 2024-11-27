The app calculates points based on the details of a receipt and lets you check the points later. 

Step 1: Build the Docker Image
In your terminal, run the following command:

docker build -t receipt-processor .

Step 2: Run the Docker Container
Now that we have the image, we can run it in a Docker container. 
By default, the application inside the container is set to listen on port 5000. However, you can map this container port to a different port on your local machine. This allows you to access the app through a browser or any other tool on your machine.
To run the application in the Docker container and map it to a specific port on your local machine, use the following command:

docker run -p 4100:4100 -e PORT=4100 receipt-processor

After running this command, the application will be accessible at http://localhost:4100 on your local machine. You can change 4100 to any other port you prefer, as long as it’s available on your machine.

Step 3: Send a Receipt for Processing
With the app running, it's time to send a receipt for processing.We have two options depending on what environment we're using.

Option 1: Use curl in WSL (Windows Subsystem for Linux)
If you’re on Windows and have WSL (Windows Subsystem for Linux) installed, you can use curl to send the receipt. Open a WSL terminal and run this:

curl -X POST http://localhost:4100/receipts/process -H "Content-Type: application/json" -d @receipt.json

-d @receipt.json sends the contents of the receipt.json file as the data in the POST request.

Option 2: Use PowerShell’s Invoke-RestMethod (Windows Native)
If you’re using PowerShell, you can do this instead:

$headers = @{ "Content-Type" = "application/json" }
$body = Get-Content -Raw -Path "receipt.json"
Invoke-RestMethod -Method POST -Uri "http://localhost:4100/receipts/process" -Headers $headers -Body $body

After you run this, the server will process the receipt and respond with an ID for that receipt.

Step 4: Get Points Using the Receipt ID
Once the receipt is processed, you'll get an ID in the response. You can use that ID to check how many points were assigned to that receipt.

In PowerShell:
$receiptId = "e1019d99-8214-4b52-a7c5-3e1f58346e2e"  # Replace with the actual ID
Invoke-RestMethod -Method GET -Uri "http://localhost:4100/receipts/$receiptId/points"

In curl:

curl -X GET http://localhost:4100/receipts/$receiptId/points

This will give you the points for that receipt based on the rules defined in the app. Just make sure to replace $receiptId with the ID you received earlier.

Step 5: List All Running Containers
To see a list of all containers on your system (running or stopped), use the following command:
docker ps -a
This will show you the containers, their statuses, and other useful information.

Step 6: Stop a Running Container
If you need to stop a container, use:
docker stop <container-id>
Replace <container-id> with the actual ID of the container you want to stop. You can get the container ID from the output of docker ps -a.

Step 7: Remove a Stopped Container
If you no longer need a stopped container, you can remove it with:

docker rm <container-id>
This will remove the container from your system.
