import { useState, useEffect } from 'react';
import { jwtDecode } from 'jwt-decode';
const WebSocketReader = () => {
  const [webSocket, setWebSocket] = useState(null);
  const [messages, setMessages] = useState([]);
  const notebookId = 27; // Fixed notebook ID

  // Establish WebSocket Connection
  useEffect(() => {

    const wsURL = `${process.env.NEXT_PUBLIC_API_BASE_URL}/ws/read/${notebookId}?owner_id=7`;
    console.log(`Connecting to WebSocket at ${wsURL}...`);

    const ws = new WebSocket(wsURL);
    setWebSocket(ws);

    ws.onopen = () => {
      console.log('WebSocket connection opened.');
    };

    ws.onmessage = (event) => {
      console.log('WebSocket message received:', event.data);
      setMessages((prevMessages) => [...prevMessages, event.data]);
    };

    ws.onerror = (event) => {
      console.error('WebSocket error:', event);
    };

    ws.onclose = (event) => {
      console.log('WebSocket connection closed:', event.reason || 'No reason provided.');
      setWebSocket(null); // Reset WebSocket instance
    };

    // Cleanup WebSocket on unmount
    return () => {
      if (ws) ws.close();
    };
  });

  return (
    <div className="min-h-screen bg-white flex flex-col items-center justify-start px-4 py-5">
      <h1 className="text-3xl font-bold text-black mb-4">Notebook Reader</h1>
      <h2 className="text-xl font-bold text-black mb-4">Messages from Notebook {notebookId}</h2>
      <div className="w-full max-w-2xl border border-gray-300 rounded p-4 bg-gray-50">
        {messages.length > 0 ? (
          messages.map((message, index) => (
            <p key={index} className="text-sm text-black mb-2">
              {message}
            </p>
          ))
        ) : (
          <p className="text-sm text-gray-500">No messages received yet.</p>
        )}
      </div>
    </div>
  );
};

export default WebSocketReader;