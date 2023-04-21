# go-subscriptions

purchase and manage recurring monthly subscriptions

## Process

- Sign up.
- Select package to subscribe to.
- Receive email to complete purchase process.
- Email has a link to pay for the package.
- On payment, stripe webhook creates a record of payment on the subscriptions service.
- Documentation sent on mail with set up / installation instructions.
- Subscription charged monthly on the account provided.

## The idea

- Run the application as a distributed system.
- Use SAGAs (orchestrator pattern) for communication between the services.
- Emit events as protobufs (gRPC) as opposed to JSON.
- Use concurrency for asynchronous requests.
- Host solution to AWS EKS.

## Technologies

- Docker
- K8s
- Go
- Gin
- Redis
- gRPC
- Stripe
