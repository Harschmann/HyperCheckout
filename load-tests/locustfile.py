import random
from locust import HttpUser, task, between


class WebsiteUser(HttpUser):
    # Simulating a very aggressive bot (fast spamming)
    wait_time = between(1, 3)

    @task
    def purchase_item(self):
        user_id = random.randint(1, 10000)
        payload = {"user_id": user_id, "product_id": 1, "quantity": 1}

        with self.client.post(
            "/purchase", json=payload, catch_response=True
        ) as response:
            if response.status_code == 200:
                response.success()
            elif response.status_code == 409:
                response.success()  # Out of stock is NOT a bug
            elif response.status_code == 429:
                # We mark 429 as "Success" because the system DID what it was supposed to do!
                # Or you can log it as a custom failure message if you prefer.
                response.failure("Rate Limited (Expected)")
            else:
                response.failure(f"Unexpected status: {response.status_code}")
