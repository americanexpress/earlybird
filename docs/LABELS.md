# Labels

By modifying `.go-earlybird/labels.json`, labels can be added to the hits based on context.  This can be done either on a same line search or on a scan of context throughout the whole file.  For example, the following line will add the `oracle` label to the returned hit:

```
password = 'abc123def456` //this is my oracle password
```

The following file will identify the password as a hit and label it as `mysql` because the word `mysql` is in the file context:

```
$servername = "localhost";
$username = "username";
$password = "password";

// Create connection
$conn = new mysqli($servername, $username, $password);

// Check connection
if ($conn->connect_error) {
    die("Connection failed: " . $conn->connect_error);
} 
echo "Connected successfully";
```
Each label configuration is applied to a hit based on its `Code`.