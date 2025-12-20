import axios from "axios";

const API_BASE = "http://localhost:8000";

export const fetchTopPlayers = async () => {
  try {
    const res = await axios.get(`${API_BASE}/api/leaderboard/top`);
    // Make sure to handle both response formats
    return Array.isArray(res.data) ? res.data : (res.data.players || []);
  } catch (error) {
    console.error('Error fetching top players:', error);
    return []; // Return empty array on error
  }
};

// SSE remains the same, just make sure server sends same shape
// In frontend/src/api/leaderboard.js
export const leaderboardSSE = () => {
  const eventSource = new EventSource(`${API_BASE}/api/leaderboard/stream`);
  return eventSource;
};
