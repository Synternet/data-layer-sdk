# RPC over Request/Reply in Data Layer

The **Request/Reply mechanism** within the Data Layer provides a powerful and efficient way for publishers to communicate with subscribers—whether for computational tasks or fetching additional data.

By leveraging **Service Protobuf specifications**, developers can **auto-generate** client and server stubs across multiple programming languages, eliminating the need for manually coding messages and communication logic. This approach significantly **reduces boilerplate code** while ensuring **strongly typed, well-structured** interactions.

---

## Defining the Service

Service names and method names are **seamlessly translated** into Data Layer subjects.

```proto
service UserService {
  rpc AddUser(AddUserRequest) returns (AddUserResponse) {}
  rpc GetUser(GetUserRequest) returns (GetUserResponse) {}
}
```

### 🔹 Subject Mapping

Using the schema above, `UserService` is mapped as follows:

- If the **service proto type** is `org.pkg.UserService`, the resulting **subject prefix** is: `service.pkg.user`
  - Although, if the proto folder strucure results in `org.pkg.version.UserService`, then the **subject prefix** is: `service.version.user`
- If the service proto type is `pkg.UserService`, the subject prefix remains:  `service.pkg.user`
- If the service proto type is `UserService`, the subject prefix is simplified to: `service.user`

When a publisher configures a **custom prefix** and **service name**, the final subject structure follows this pattern:

```
{prefix}.{name}.service.pkg.user
```

### Method Subject Formatting

Each **service method** is tokenized based on **CamelCase**. For example:

- `AddUser` → `add.user`
- `GetUser` → `get.user`

Thus, a full subject for an RPC call would be:

```
{prefix}.{name}.service.pkg.user.add.user
```

While this may seem verbose, **customization options** allow overriding these subjects.

---

## Customization via Options

Protobuf **options** provide fine-grained control over subject naming and behavior.

### Available Options

- **`subject_prefix`** → Specifies a **custom service subject prefix**.
- **`subject_suffix`** → Defines a **custom method subject suffix**.
- **`disable_inputs`** → Marks an RPC method as a **pure stream** (no input, only continuous output).

### Example: Custom Subject Naming

```proto
import "synternet/rpc/options.proto";
import "google/protobuf/empty.proto";

service UserService {
  option (subject_prefix) = "users";

  rpc AddUser(AddUserRequest) returns (AddUserResponse) {
    option (subject_suffix) = "add";
  }
  rpc GetUser(GetUserRequest) returns (GetUserResponse) {
    option (subject_suffix) = "get";
  }

  rpc Registrations(google.protobuf.Empty) returns (stream Registrations) {
    option (subject_suffix) = "registrations";
    option (disable_inputs) = true;
  }
}
```

This schema results in the following subjects:

#### **Service Calls**

- `{prefix}.{name}.users.add`
- `{prefix}.{name}.users.get`

#### **Streaming Events**

- `{prefix}.{name}.users.registrations`

The `Registrations` method is a **free-running stream**—it does not accept input but **continuously pushes events** to subscribers.

---

## Using the Go-Generated Code

You can create testnet publisher and subscriber in order to test this example, or you can use the local NATS server.

If you want to run it locally, then install NATS tools, NATS server, and execute `setup-nats.sh` in the `examples/rpc` folder. After the NATS server is configured,
you can run it using this command:

```
nats-server -c .nats/server.conf -DV
```

After this you can run the publisher:

```
NATS_URL=127.0.0.1 PUBLISHER_CREDS=.nats/creds/test/publisher/publisher.creds go run ./publisher
```

And then the subscriber:

```
NATS_URL=127.0.0.1 SUBSCRIBER_CREDS=.nats/creds/test/subscriber/subscriber.creds go run ./subscriber
```

You should be able to see the result of interaction between the publisher and the subscriber.

### Publisher-Side Implementation

To use the service, first **generate** the protobuf messages, server, and client code. Then, integrate them into the **Data Layer**.

```go
type MyPublisher struct {
  *service.Service
  *rpc.ServiceRegistrar
  userService UserService
}

func (p *MyPublisher) Start() context.Context {
  // Registers UserService with Data Layer, defining its subject dynamically
  servicetypes.RegisterUserServiceServer(p, &p.userService)

  // Other initialization logic...
}

// Ensure UserService correctly implements the interface
var _ servicetypes.UserServiceServer = (*UserService)(nil)

type BackendService struct {
  // Implement service logic here
}
```

Don't forget to set Protobuf Codec with `service.WithCodec(codec.NewProtoJsonCodec())` option!

---

### Client-Side Integration

To make RPC calls via the Data Layer, initialize a **client connection** and invoke the service methods.

```go
func setupClient(ctx context.Context, p *service.Service) servicetypes.UserServiceClient {
  userConn := rpc.NewClientConn(ctx, p, "your-organization.publisher")
  return servicetypes.NewUserServiceClient(userConn)
}

func main() {
  ctx, cancel := context.WithCancel(context.Background())
  defer cancel()

  // Configure the service connection
  var service service.Service
  opts := []options.Option{
    service.WithContext(ctx),
    service.WithName(name),
    service.WithPrefix(prefix),
    service.WithNats(conn),
    service.WithNKeySeed(nkey),
    service.WithUserCreds(path_to_user_creds),
    service.WithCodec(codec.NewProtoJsonCodec()),
  }
  service.Configure(opts...)

  user := setupClient(ctx, &service)

  ctx1, cancel1 := context.WithTimeout(ctx, time.Second)
  result, err := user.AddUser(ctx1, &servicetypes.AddUserRequest{...})
  cancel1()

  if err != nil {
    log.Fatalf("Error calling AddUser: %v", err)
  }

  fmt.Printf("User added: %+v\n", result)

  // Other client logic...
}
```

---

## Further Exploration

For detailed usage, refer to the **individual examples** and explore more customization options.

This setup ensures a **robust, scalable, and efficient** RPC over Request/Reply mechanism within the Data Layer. 🚀
