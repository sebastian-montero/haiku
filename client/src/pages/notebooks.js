import { useState, useEffect } from 'react';
import { jwtDecode } from 'jwt-decode';
import NavBar from '../components/nav'; // Adjust the path based on your file structure

let ws = null;

export default function FrontPage() {
  const [notebooks, setNotebooks] = useState([]);
  const [newNotebookTitle, setNewNotebookTitle] = useState("");
  const [message, setMessage] = useState(null);
  const [ownerId, setOwnerId] = useState(null);
  const [selectedNotebook, setSelectedNotebook] = useState(null);
  const [sessionId, setSessionId] = useState(null);
  const [modalOpen, setModalOpen] = useState(false);
  const [webSocket, setWebSocket] = useState(null);
  const [textBoxContent, setTextBoxContent] = useState("");

  // Decode JWT and set owner_id
  useEffect(() => {
    const token = localStorage.getItem('jwt');
    if (token) {
      try {
        const decoded = jwtDecode(token); // Decode the JWT
        setOwnerId(decoded.user_id); // Assuming `id` is the field in the JWT payload
      } catch (error) {
        console.error('Failed to decode JWT:', error);
      }
    }
  }, []);

  // Fetch Notebooks by Owner
  useEffect(() => {
    if (!ownerId) return;

    console.log(ownerId);

    async function fetchNotebooks() {
      try {
        const token = localStorage.getItem('jwt');
        const response = await fetch(
          `${process.env.NEXT_PUBLIC_API_BASE_URL}/notebooks/by_owner/${ownerId}`,
          {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          }
        );
        if (response.ok) {
          const data = await response.json();
          setNotebooks(data);
        } else {
          setMessage({ type: 'error', text: 'failed to fetch notebooks.' });
        }
      } catch (error) {
        setMessage({ type: 'error', text: 'an unexpected error occurred.' });
      }
    }

    fetchNotebooks();
  }, [ownerId]);

  // Create Notebook
  const handleCreateNotebook = async (e) => {
    e.preventDefault();
  
    if (!newNotebookTitle.trim()) {
      setMessage({ type: "error", text: "Notebook title cannot be empty." });
      return;
    }
  
    try {
      const token = localStorage.getItem("jwt");
  
      // Create the notebook
      const notebookResponse = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/notebooks`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            title: newNotebookTitle,
            owner_id: ownerId,
          }),
        }
      );
  
      if (notebookResponse.ok) {
        const newNotebook = await notebookResponse.json();
  
        // Safely update the notebooks state
        setNotebooks((prevNotebooks) =>
          Array.isArray(prevNotebooks) ? [newNotebook, ...prevNotebooks] : [newNotebook]
        );
  
        // Create an inactive session for the new notebook
        const sessionResponse = await fetch(
          `${process.env.NEXT_PUBLIC_API_BASE_URL}/sessions`,
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
              Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify({
              owner_id: ownerId,
              notebook_id: newNotebook.id,
              is_active: false
            }),
          }
        );
  
        if (sessionResponse.ok) {
          setMessage({ type: "success", text: "Notebook created successfully." });
        } else {
          setMessage({ type: "error", text: "Notebook created, but failed to create an inactive session." });
        }
      } else {
        const errorData = await notebookResponse.json();
        setMessage({ type: "error", text: `Error: ${errorData.message}` });
      }
  
      setNewNotebookTitle("");
    } catch (error) {
      setMessage({ type: "error", text: "An unexpected error occurred." });
    }
  };

  function formatDate(dateString) {
    const date = new Date(dateString);
    return new Intl.DateTimeFormat('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    }).format(date);
  }

  const handleNotebookClick = async (notebookId) => {
    try {
      const token = localStorage.getItem("jwt");
  
      // First, try to fetch the session by notebook ID
      const getSessionResponse = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/sessions/by_notebook/${notebookId}`,
        {
          method: "GET",
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
  
      let session;
      if (getSessionResponse.ok) {
        session = await getSessionResponse.json();
      } else if (getSessionResponse.status === 404) {
        console.log("No existing session found. Creating a new session...");
        const createSessionResponse = await fetch(
          `${process.env.NEXT_PUBLIC_API_BASE_URL}/sessions`,
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
              Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify({
              owner_id: ownerId,
              notebook_id: notebookId,
            }),
          }
        );
  
        if (!createSessionResponse.ok) {
          setMessage({ type: "error", text: "Failed to create a new session." });
          return;
        }
  
        session = await createSessionResponse.json();
        console.log("New session created:", session);
        setMessage("");
      } else {
        setMessage({ type: "error", text: "Failed to fetch or create session." });
        return;
      }
  
      // Update session to set is_active to true
      const updateSessionResponse = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/sessions/${session.id}`,
        {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            is_active: true,
          }),
        }
      );
  
      if (!updateSessionResponse.ok) {
        setMessage({ type: "error", text: "Failed to activate the session." });
        return;
      }
    
      // Fetch the latest content for the notebook
      const getContentResponse = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/content/by_session/${session.id}`,
        {
          method: "GET",
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
  
      let latestContent = "";
      if (getContentResponse.ok) {
        const content = await getContentResponse.json();
        latestContent = content.content || "";
      } else {
        console.log("No content found for the session.");
      }

      try {
        const data = { message: "Notebook opened", timestamp: new Date() };   
        await openWebSocketConnection(notebookId, ownerId, data);
      } catch (error) {
        console.error("Error handling notebook click:", error);
        setMessage({ type: "error", text: "Failed to process notebook click." });
      }
  
      // Set session state and open the modal with content
      setSessionId(session.id);
      setSelectedNotebook(notebookId);
      setTextBoxContent(latestContent); // Display latest content in the modal's text box
      setModalOpen(true);

    } catch (error) {
      setMessage({ type: "error", text: "An unexpected error occurred." });
      console.error(error);
    }


  };

  const openWebSocketConnection = async (notebookId, ownerId, data) => {
    try {
      // Construct WebSocket URL
      const wsURL = `${process.env.NEXT_PUBLIC_API_BASE_URL}/ws/write/${notebookId}?owner_id=${ownerId}`;
      console.log(`Connecting to WebSocket at ${wsURL}...`);
  
      // If WebSocket is already open, just send the data
      if (ws && ws.readyState === WebSocket.OPEN) {
        console.log("WebSocket already open. Sending data...");
        ws.send(JSON.stringify(data));
        return;
      }
  
      ws = new WebSocket(wsURL);
  
      ws.onopen = () => {
        console.log("WebSocket connection opened.");
        ws.send(JSON.stringify(data));
        console.log("Data sent:", data);
      };
  
      ws.onerror = (event) => {
        console.error("WebSocket error:", event);
      };
  
      ws.onclose = (event) => {
        console.log("WebSocket connection closed:", event.reason || "No reason provided.");
        ws = null;
      };
    } catch (error) {
      console.error("Error during WebSocket operation:", error);
    }
  };

  const sendMessageToWebSocket = (data) => {
    console.log(ws);  
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify(data));
    } else {
      console.error("WebSocket is not open. Cannot send message.");
    }
  };

  const handleSaveContent = async () => {
    let content = textBoxContent;
    if (!content.trim()) {
      content = "";
    }

    try {
      const token = localStorage.getItem("jwt");
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/content`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          session_id: sessionId,
          content: content,
        }),
      });

      if (response.ok) {
        setMessage({ type: "success", text: "Content saved successfully." });
      } else {
        setMessage({ type: "error", text: "Failed to save content." });
      }
    } catch (error) {
      console.error("An error occurred while saving content:", error);
      setMessage({ type: "error", text: "An unexpected error occurred." });
    }
  };

  const closeModal = async () => {
    try {
      const token = localStorage.getItem("jwt");
      if (!sessionId) {
        console.error("Session ID is not set.");
        return;
      }

      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/sessions/${sessionId}/end`,
        {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            owner_id: ownerId,
          }),
        }
      );
  
      if (!response.ok) {
        console.error("Failed to end session.");
      } else {
        console.log("Session ended successfully.");
        setModalOpen(false);
      }
    } catch (error) {
      console.error("An error occurred while ending the session:", error);
    } finally {
      setSelectedNotebook(null);
      if (ws) {
        sendMessageToWebSocket({'type':'end'});
      }
      setSessionId(null);
    }
  };

  return (
      <div className="min-h-screen bg-white flex flex-col items-center justify-start px-4 py-5">
      <NavBar />


      <div className="items-center justify-start py-2">
      <h2 className="text-xl font-bold text-gray-800 tracking-wide text-left py-4 px-2">notebooks</h2>

      {/* Create Notebook Form */}
      <form onSubmit={handleCreateNotebook} className="w-full max-w-md mb-1 px-2">
        <div className="flex items-center space-x-4">
          <input
            type="text"
            placeholder="new notebook title"
            value={newNotebookTitle}
            onChange={(e) => setNewNotebookTitle(e.target.value)}
            className="w-full border border-gray-300 px-4 py-2 text-black font-bold shadow-sm focus:ring-gray-800 focus:border-gray-800 sm:text-sm"
          />
          <button
            type="submit"
            className="py-2 px-4 text-black font-bold hover:bg-black hover:text-white focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-800"
          >
            create
          </button>
        </div>
      </form>

        {/* Message Display */}
        {message && !modalOpen && ( 
        <div
          className={`text-xs my-3 text-center ${
            message.type === 'success' ? 'text-green-700' : 'text-red-700'
          }`}
        >
          {message.text}
        </div>
      )}

      {/* List of Notebooks */}
      <div className="w-full max-w-md">
        {notebooks && notebooks.length > 0 ? (
          notebooks.map((notebook) => (
            <button
              key={notebook.id}
              className="w-full text-left border-b border-gray-300 text-gray-700 py-4 flex flex-col space-y-1 bg-white hover:bg-black hover:text-white focus:outline-none focus:ring-2 focus:ring-gray-300 p-2"
              onClick={() => handleNotebookClick(notebook.id)}
            >
              <h2 className="font-bold">{notebook.title}</h2>
              {notebook.last_updated_at && (
                <p className="text-xs">
                  last updated: {notebook.last_updated_at ? formatDate(notebook.last_updated_at) : 'not updated yet'}
                </p>
              )}
              {notebook.latest_content && (
                <div className="text-xs">
                  <p>
                    {notebook.latest_content.length > 200
                      ? `${notebook.latest_content.slice(0, 200)}...`
                      : notebook.latest_content}
                  </p>
                </div>
              )}
            </button>
          ))
        ) : (
          <p className="my-10 text-center text-gray-500 text-xs">no notebooks found.</p>
        )}
      </div>
      </div>
      {/* Modal */}
      {modalOpen && (
        <div className="fixed inset-0 bg-gray-800 bg-opacity-50 flex items-center justify-center">
          <div className="bg-white p-6 rounded shadow-md w-full max-w-md">
            <h2 className="text-xl font-bold mb-4 text-black">streaming</h2>
            <textarea
              value={textBoxContent}
              onChange={(e) => {
                const content = e.target.value;
                setTextBoxContent(content);   
                sendMessageToWebSocket({'type':'send', 'content': content});
              }}
              placeholder="start writing..."
              className="w-full h-32 border border-gray-300 p-2 mb-4 text-black"
            />
            <div className="flex justify-between">
              <button
                onClick={closeModal}
                className="px-4 py-2 bg-gray-300 text-black bg-transparent hover:text-white hover:bg-black font-medium"
              >
                close
              </button>
              <button
                onClick={handleSaveContent}
                className="px-4 py-2 bg-gray-300 text-black bg-transparent hover:text-white hover:bg-blue-900 font-medium"
              >
                save
              </button>
            </div>
            {/* Message Display */}
                {message && (
                <div
                  className={`text-xs my-3 text-center ${
                    message.type === 'success' ? 'text-green-700' : 'text-red-700'
                  }`}
                >
                  {message.text}
                </div>
              )}
          </div>
        </div>
      )}
      </div>
  );
}