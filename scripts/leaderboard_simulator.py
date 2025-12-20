import requests
import random
import time

API_BASE_URL = "http://localhost:8000/api/leaderboard"

# Simulate score submission
def submit_score(user_id):
    score = random.randint(100, 10000)
    response = requests.post(
        f"{API_BASE_URL}/submit",
        json={"user_id": user_id, "score": score}
    )
    return response.json()

# Fetch top players
def get_top_players():
    response = requests.get(f"{API_BASE_URL}/top")
    return response.json()

# Fetch user rank
def get_user_rank(user_id):
    response = requests.get(f"{API_BASE_URL}/rank/{user_id}")
    return response.json()

if __name__ == "__main__":
    while True:
        user_id = random.randint(1, 1000000)
        print("Submitting score...")
        print(submit_score(user_id))
        print("Top players:")
        print(get_top_players())
        print(f"User {user_id} rank:")
        print(get_user_rank(user_id))
        time.sleep(random.uniform(0.5, 2))  # Simulate real user interaction
