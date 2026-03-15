import requests
import time

BASE = "http://localhost:8080"

session = requests.Session()

print("1️⃣ Register user")

r = session.post(
    f"{BASE}/api/user/register",
    json={"login": "testuser", "password": "12345"},
)

print("status:", r.status_code)

print("2️⃣ Login")

r = session.post(
    f"{BASE}/api/user/login",
    json={"login": "testuser", "password": "12345"},
)

print("status:", r.status_code)

token = r.headers.get("Authorization")

if token:
    session.headers["Authorization"] = token

print("3️⃣ Upload order")

order = "79927398713"

r = session.post(
    f"{BASE}/api/user/orders",
    data=order,
)

print("status:", r.status_code)

print("4️⃣ Upload same order again")

r = session.post(
    f"{BASE}/api/user/orders",
    data=order,
)

print("status:", r.status_code)

print("5️⃣ Get orders list")

r = session.get(f"{BASE}/api/user/orders")

print("status:", r.status_code)
print("orders:", r.text)

print("6️⃣ Get balance")

r = session.get(f"{BASE}/api/user/balance")

print("status:", r.status_code)
print("balance:", r.text)

print("7️⃣ Withdraw")

withdraw = {
    "order": "12345678903",
    "sum": 10
}

r = session.post(
    f"{BASE}/api/user/balance/withdraw",
    json=withdraw,
)

print("status:", r.status_code)

print("8️⃣ Get withdrawals")

r = session.get(f"{BASE}/api/user/withdrawals")

print("status:", r.status_code)
print("withdrawals:", r.text)

print("✅ TEST FINISHED")
