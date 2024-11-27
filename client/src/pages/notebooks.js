import { useState, useEffect } from 'react';
import { jwtDecode } from 'jwt-decode';

export default function FrontPage() {
  const [notebooks, setNotebooks] = useState([]);
  const [newNotebookTitle, setNewNotebookTitle] = useState('');
  const [message, setMessage] = useState(null);
  const [ownerId, setOwnerId] = useState(null);

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
      setMessage({ type: 'error', text: 'notebook title cannot be empty.' });
      return;
    }
  
    try {
      const token = localStorage.getItem('jwt');
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/notebooks`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          title: newNotebookTitle,
          owner_id: ownerId,
        }),
      });
  
      if (response.ok) {
        const newNotebook = await response.json();
  
        // Safely update the notebooks state
        setNotebooks((prevNotebooks) =>
          Array.isArray(prevNotebooks) ? [newNotebook, ...prevNotebooks] : [newNotebook]
        );
  
        setNewNotebookTitle('');
        setMessage({ type: 'success', text: 'notebook created successfully.' });
      } else {
        const errorData = await response.json();
        setMessage({ type: 'error', text: `error: ${errorData.message}` });
      }
    } catch (error) {
      setMessage({ type: 'error', text: 'an unexpected error occurred.' });
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

  return (
    <div className="min-h-screen bg-white flex flex-col items-center justify-start px-4 py-8">
      <h1 className="text-3xl font-bold text-gray-800 tracking-wide mb-6">notebooks</h1>

      {/* Create Notebook Form */}
      <form onSubmit={handleCreateNotebook} className="w-full max-w-md mb-8">
        <div className="flex items-center space-x-4">
          <input
            type="text"
            placeholder="new notebook title"
            value={newNotebookTitle}
            onChange={(e) => setNewNotebookTitle(e.target.value)}
            className="w-full border border-gray-300 px-4 py-2 text-gray-700 font-bold shadow-sm focus:ring-gray-800 focus:border-gray-800 sm:text-sm"
          />
          <button
            type="submit"
            className="py-2 px-4 bg-gray-800 text-white font-bold hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-800"
          >
            create
          </button>
        </div>
      </form>

        {/* Message Display */}
        {message && (
        <div
          className={`mb-4 text-center ${
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
            <div
              key={notebook.id}
              className="border-b border-gray-300 py-4 flex flex-col space-y-1"
            >
              <h2 className="text-lg font-bold text-gray-800">{notebook.title}</h2>
              <p className="text-sm text-gray-500">
                last updated: {notebook.last_updated_at ? formatDate(notebook.last_updated_at) : 'not updated yet'}
              </p>
              <p className="text-sm text-gray-500">
                created: {notebook.created_at ? formatDate(notebook.created_at) : 'unknown'}
              </p>
              {notebook.latest_content && (
                <p className="text-sm text-gray-700">content: {notebook.latest_content}</p>
              )}
            </div>
          ))
        ) : (
          <p className="text-center text-gray-500">no notebooks found.</p>
        )}
      </div>
    </div>
  );
}