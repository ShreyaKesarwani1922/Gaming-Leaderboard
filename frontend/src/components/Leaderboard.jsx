import { useEffect, useState } from "react";
import { fetchTopPlayers, leaderboardSSE } from "../api/leaderboard";

export default function Leaderboard() {
  const [players, setPlayers] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Initial fetch
    fetchTopPlayers()
      .then(data => setPlayers(data))
      .finally(() => setLoading(false));

    // Live updates via SSE
    const eventSource = leaderboardSSE();

    eventSource.onmessage = (event) => {
      const updatedLeaderboard = JSON.parse(event.data);
      setPlayers(updatedLeaderboard);
    };

    eventSource.onerror = () => {
      console.error("SSE connection lost");
      eventSource.close();
    };

    return () => eventSource.close();
  }, []);

  if (loading) return <p>Loading leaderboard...</p>;

  return (
    <table style={{ width: "100%", borderCollapse: "collapse" }}>
      <thead>
        <tr>
          <th>Rank</th>
          <th>User</th>
          <th>Total Score</th>
        </tr>
      </thead>
      <tbody>
        {players.map((p, index) => (
          <tr key={p.user_id}>
            <td>{index + 1}</td>
            <td>{p.user_name}</td>
            <td>{p.total_score}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}
