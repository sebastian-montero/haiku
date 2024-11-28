import { useState, useEffect } from 'react';
import { jwtDecode } from 'jwt-decode';
import NavBar from '../components/nav'; // Adjust the path based on your file structure

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
        {message && (
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
      </div>
  );
}