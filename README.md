# go-subscriptions

purchase and manage recurring monthly subscriptions

## Process

- Sign up, select package to pay for.
- Receive email to complete purchase process.
- Email has a link to pay for the package.
- On payment, account is created.
- Documentation sent on mail with set up / installation instructions.
- Subscription charged monthly on the account provided.

## The idea

- Run the application as a distributed system.
- Use SAGAs (orchestrator pattern) for communication between the services.
- Use concurrency for asynchronous requests.
- Host solution to AWS EKS.

## Technologies

- Docker
- K8s
- Go
- Gin
- Redis
- Stripe
