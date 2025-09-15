Payment System with Stripe, PostgreSQL, and Gin
A complete payment processing system built with Go, Stripe, PostgreSQL, and the Gin web framework. This system handles customer management, payment intents, webhook processing, and payment tracking.

🚀 Features
Customer Management - Create and manage Stripe customers

Payment Processing - Create and confirm payment intents

Webhook Handling - Process Stripe webhooks for payment events

Database Persistence - Store customers, payments, and payment intents in PostgreSQL

RESTful API - Clean API endpoints for all operations

Security - Webhook signature verification and authentication middleware

📋 Prerequisites
Go 1.21+

PostgreSQL 12+

Stripe account

Stripe CLI (for local development)

🛠️ Installation
Clone the repository

bash
git clone <your-repo-url>
cd paymentSystem
Install dependencies

bash
go mod tidy
Set up PostgreSQL database

bash
createdb payment_system
Install Stripe CLI

bash
# Windows - download from https://github.com/stripe/stripe-cli/releases
# Or use scoop:
scoop install stripe
⚙️ Configuration
Copy environment file

bash
cp .env.example .env
Update .env with your values

env
STRIPE_SECRET_KEY=sk_test_your_stripe_secret_key
STRIPE_WEBHOOK_SECRET=whsec_your_webhook_secret
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=payment_system
PORT=8080
Get your Stripe keys

Visit: https://dashboard.stripe.com/apikeys

Copy Publishable key and Secret key

🚦 Running the Application
1. Start the server
bash
go run main.go
2. Set up webhooks with Stripe CLI
bash
stripe login
stripe listen --forward-to localhost:8080/api/v1/webhooks
3. Copy the webhook secret from Stripe CLI output and update your .env file
4. Test the system
bash
stripe trigger payment_intent.succeeded
📊 API Endpoints
Customers
POST /api/v1/customers - Create a new customer

GET /api/v1/customers - List all customers

GET /api/v1/customers/:id - Get customer by ID

DELETE /api/v1/customers/:id - Delete customer

Payments
POST /api/v1/payment-intents - Create payment intent

GET /api/v1/payment-intents/:id/status - Get payment status

GET /api/v1/payments - List all payments

Webhooks
POST /api/v1/webhooks - Handle Stripe webhooks

Health
GET /health - Health check endpoint

🎯 Usage Examples
Create a Customer
bash
curl -X POST http://localhost:8080/api/v1/customers \
  -H "Content-Type: application/json" \
  -d '{
    "email": "customer@example.com",
    "name": "John Doe"
  }'
Create a Payment Intent
bash
curl -X POST http://localhost:8080/api/v1/payment-intents \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 1000,
    "currency": "usd",
    "customer_id": 1,
    "description": "Test payment"
  }'
Check Payment Status
bash
curl http://localhost:8080/api/v1/payment-intents/pi_123456789/status
🔧 Webhook Setup
Using Stripe CLI (Recommended)
bash
stripe listen --forward-to localhost:8080/api/v1/webhooks
Using ngrok
bash
ngrok http 8080
# Then use the ngrok URL in Stripe dashboard
Manual Setup
Visit: https://dashboard.stripe.com/webhooks

Add endpoint: https://yourdomain.com/api/v1/webhooks

Select events to listen for

Copy the webhook secret to your .env file

🧪 Testing
Test Cards (Stripe Test Mode)
Success: 4242 4242 4242 4242

Failure: 4000 0000 0000 0002

Authentication Required: 4000 0025 0000 3155

Test Webhooks
bash
stripe trigger payment_intent.succeeded
stripe trigger payment_intent.payment_failed
🗄️ Database Schema
Customers
id, stripe_id, email, name, created_at, updated_at

Payment Intents
id, stripe_id, customer_id, amount, currency, status, description, metadata, created_at, updated_at

Payments
id, stripe_id, payment_intent_id, amount, currency, status, created_at

🛡️ Security Features
Webhook signature verification

Environment variable configuration

Input validation

Error handling

CORS middleware

Authentication middleware (ready for JWT integration)

🚀 Deployment
Environment Variables for Production
env
STRIPE_SECRET_KEY=sk_live_your_production_key
STRIPE_WEBHOOK_SECRET=whsec_your_production_secret
DB_HOST=your_production_db_host
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=payment_system_prod
PORT=8080
Build for Production
bash
go build -o payment-system main.go
Run in Production
bash
./payment-system