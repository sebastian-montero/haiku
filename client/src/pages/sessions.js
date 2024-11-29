import { useState, useEffect } from "react";
import NavBar from "../components/nav"; // Adjust the path based on your file structure

let ws = null; // WebSocket instance

export default function ActiveSessionsGrid() {
  const [sessions, setSessions] = useState([]);
  const [notebooks, setNotebooks] = useState({});
  const [message, setMessage] = useState(null);
  const [webSocketData, setWebSocketData] = useState(""); // Latest WebSocket message content
  const [isModalOpen, setIsModalOpen] = useState(false); // Modal state
  const [activeNotebook, setActiveNotebook] = useState(null); // Active notebook title for modal
  const [isConnecting, setIsConnecting] = useState(false); // Connection state

  // Fetch active sessions
  useEffect(() => {
    async function fetchActiveSessions() {
      try {
        const token = localStorage.getItem("jwt");
        const response = await fetch(
          `${process.env.NEXT_PUBLIC_API_BASE_URL}/sessions`,
          {
            method: "GET",
            headers: {
              Authorization: `Bearer ${token}`,
            },
          }
        );

        if (response.ok) {
          const data = await response.json();
          setSessions(data);
        } else {
          setMessage({ type: "error", text: "Failed to fetch active sessions." });
        }
      } catch (error) {
        setMessage({ type: "error", text: "An unexpected error occurred." });
      }
    }

    fetchActiveSessions();
  }, []);

  // Fetch notebook details for each session
  useEffect(() => {
    async function fetchNotebooksForSessions() {
      const token = localStorage.getItem("jwt");
      const notebookPromises = sessions.map(async (session) => {
        const response = await fetch(
          `${process.env.NEXT_PUBLIC_API_BASE_URL}/notebooks/${session.notebook_id}`,
          {
            method: "GET",
            headers: {
              Authorization: `Bearer ${token}`,
            },
          }
        );

        if (response.ok) {
          const notebook = await response.json();
          return { notebookId: session.notebook_id, notebook };
        } else {
          console.error(
            `Failed to fetch notebook for notebook_id ${session.notebook_id}`
          );
          return { notebookId: session.notebook_id, notebook: null };
        }
      });

      const notebooksData = await Promise.all(notebookPromises);
      const notebooksMap = {};
      notebooksData.forEach(({ notebookId, notebook }) => {
        if (notebook) {
          notebooksMap[notebookId] = notebook;
        }
      });

      setNotebooks(notebooksMap);
    }

    if (sessions.length > 0) {
      fetchNotebooksForSessions();
    }
  }, [sessions]);

  // Handle notebook button click and open WebSocket
  const handleNotebookClick = (notebookId, notebookTitle) => {
    try {
      const wsURL = `${process.env.NEXT_PUBLIC_API_BASE_URL}/ws/read/${notebookId}`;
      console.log(`Connecting to WebSocket: ${wsURL}`);
      setActiveNotebook(notebookTitle); // Set active notebook title
      setIsConnecting(true); // Show loading spinner

      // Close existing WebSocket if open
      if (ws) {
        ws.close();
      }

      // Create a new WebSocket connection
      ws = new WebSocket(wsURL);

      ws.onopen = () => {
        console.log("WebSocket connection opened.");
        setIsConnecting(false); // Stop loading spinner
        setIsModalOpen(true); // Open modal
      };

      ws.onmessage = (event) => {
        try {
          const parsedMessage = JSON.parse(event.data);
          if (parsedMessage.type === "send" && parsedMessage.content) {
            setWebSocketData(parsedMessage.content); // Update latest content
          } else {
            console.error("Unexpected WebSocket message format:", event.data);
          }
        } catch (error) {
          console.error("Failed to parse WebSocket message:", event.data);
        }
      };

      ws.onerror = (error) => {
        console.error("WebSocket error:", error);
        setMessage({ type: "error", text: "WebSocket connection error." });
        setIsConnecting(false); // Stop loading spinner
      };

      ws.onclose = () => {
        console.log("WebSocket connection closed.");
        ws = null; // Reset WebSocket
      };
    } catch (error) {
      console.error("Error during WebSocket operation:", error);
      setMessage({ type: "error", text: "Failed to open WebSocket connection." });
      setIsConnecting(false); // Stop loading spinner
    }
  };

  const closeModal = () => {
    if (ws) {
      ws.close();
    }
    setIsModalOpen(false);
    setWebSocketData(""); // Clear latest message
    setActiveNotebook(null); // Reset active notebook
  };

  return (
    <div className="min-h-screen bg-white flex flex-col items-center justify-start py-5 p-5">
      <NavBar />

      <div className="items-center justify-start py-2">
        <h2 className="text-xl font-bold text-black tracking-wide text-left py-4">
          live sessions
        </h2>
        {message && (
          <div
            className={`text-sm my-3 text-center ${
              message.type === "success" ? "text-green-700" : "text-red-700"
            }`}
          >
            {message.text}
          </div>
        )}
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 w-full max-w-7xl">
          {sessions.map((session) => {
            const notebook = notebooks[session.notebook_id];
            return (
              <button
  key={session.id}
  onClick={() =>
    handleNotebookClick(session.notebook_id, notebook?.title)
  }
  className="p-4 shadow rounded bg-white hover:bg-gray-100 focus:outline-none text-left"
>
  {notebook ? (
    <>
      <h2 className="text-lg font-bold text-black mb-2">
        {notebook.title}
      </h2>
      <p className="text-xs text-gray-500">
        {session.started_at
          ? new Date(session.started_at).toLocaleString()
          : "Unknown"}
      </p>
      <p
        className="text-xs text-gray-500 mt-2 line-clamp-2"
        style={{
          display: "-webkit-box",
          WebkitLineClamp: 2,
          WebkitBoxOrient: "vertical",
          overflow: "hidden",
          textOverflow: "ellipsis",
        }}
      >
        {notebook.latest_content || "No content yet."}
      </p>
    </>
  ) : (
    <p className="text-gray-500 text-sm">Loading notebook...</p>
  )}
</button>
            );
          })}
        </div>
      </div>

      {/* Loading Spinner */}
      {isConnecting && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center">
          <div className="bg-white p-6 rounded shadow-md w-full max-w-md text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-t-4 border-gray-500 mx-auto mb-4"></div>
            <p className="text-lg font-bold text-black">
              Connecting to {activeNotebook || "Notebook"}...
            </p>
          </div>
        </div>
      )}

      {/* Modal */}
      {isModalOpen && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center">
          <div className="bg-white p-6 rounded shadow-md w-full max-w-md">
          <div className="flex items-center justify-between mb-5">
              <h2 className="text-xl font-bold text-black">{activeNotebook}</h2>
        <div className="flex items-center">
          <div className="w-2 h-2 mb-4 bg-red-600 rounded-full animate-pulse mr-1"></div>
          <span className="text-sm mb-4 font-bold text-red-600">LIVE</span>
        </div>
      </div>
            <div className="text-black text-sm whitespace-pre-wrap break-words">
              {webSocketData || "waiting..."}
            </div>
            <div className="flex justify-end mt-4">
              <button
                onClick={closeModal}
                className="px-4 py-2 bg-gray-300 text-black bg-transparent hover:text-white hover:bg-black font-medium"
              >
                close
              </button>
            </div>
          </div>
        </div>
      )}


      
    </div>
  );
}