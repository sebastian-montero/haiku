import { useState } from 'react';
import { useRouter } from 'next/router';

export default function Login() {
  const [formData, setFormData] = useState({
    username: '',
    password: '',
  });

  const [message, setMessage] = useState(null);
  const router = useRouter();

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((prevData) => ({
      ...prevData,
      [name]: value,
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          username: formData.username,
          password: formData.password,
        }),
      });

      if (response.ok) {
        const data = await response.json();

        if (data.token) {
          localStorage.setItem('jwt', data.token);
        }
        setMessage({ type: 'success', text: 'login successful.' });
        router.push('/notebooks');
      } else {
        const errorData = await response.json();
        setMessage({ type: 'error', text: `error: ${errorData.message}` });
      }
    } catch (error) {
      setMessage({ type: 'error', text: 'an unexpected error occurred.' });
    }
  };

  return (
    <div className="min-h-screen bg-white flex items-center justify-center px-4">
      <div className="max-w-md w-full space-y-6">
        <h1 className="text-3xl font-bold text-gray-800 tracking-wide text-center">haiku⿻</h1>
        <h2 className="text-xl font-bold text-gray-800 tracking-wide text-center">log in</h2>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label
              htmlFor="username"
              className="block text-sm font-medium text-gray-700"
            >
              username
            </label>
            <input
              type="text"
              id="username"
              name="username"
              value={formData.username}
              onChange={handleChange}
              required
              className="mt-1 block w-full border-gray-300 shadow-sm focus:ring-gray-800 focus:border-gray-800 sm:text-sm text-gray-700 font-bold"
            />
          </div>

          <div>
            <label
              htmlFor="password"
              className="block text-sm font-medium text-gray-700"
            >
              password
            </label>
            <input
              type="password"
              id="password"
              name="password"
              value={formData.password}
              onChange={handleChange}
              required
              className="mt-1 block w-full border-gray-300 shadow-sm focus:ring-gray-800 focus:border-gray-800 sm:text-sm text-gray-700 font-bold"
            />
          </div>

          {message && (
            <div
              className={`p-4 text-sm ${
                message.type === 'success'
                  ? 'text-green-700 text-center'
                  : 'text-red-700 text-center'
              }`}
            >
              {message.text}
            </div>
          )}

          <button
            type="submit"
            className="w-full py-2 px-4 text-black font-medium hover:bg-black hover:text-white focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-800"
          >
            log in
          </button>
          
        </form>
        <p className="text-xs text-center text-black">
  don't have an account?
  <a href="#" className="hover:underline ml-1">sign up</a>
          </p>
      </div>
    </div>
  );
}