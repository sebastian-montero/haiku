import { useState, useEffect } from "react";
import { useRouter } from "next/router";
import NavBar from '../components/nav'; // Adjust the path based on your file structure

export default function Profile() {
  const [formData, setFormData] = useState({
    username: "",
    email: "",
  });
  const [message, setMessage] = useState(null);
  const router = useRouter();

  // Fetch current user data
  useEffect(() => {
    const fetchUser = async () => {
      try {
        const token = localStorage.getItem("jwt");
        const userId = parseJwt(token).user_id; // Decode the JWT to extract user_id

        const response = await fetch(
          `${process.env.NEXT_PUBLIC_API_BASE_URL}/users/${userId}`,
          {
            method: "GET",
            headers: {
              Authorization: `Bearer ${token}`,
            },
          }
        );

        if (response.ok) {
          const data = await response.json();
          setFormData({
            username: data.username,
            email: data.email,
          });
        } else {
          setMessage({ type: "error", text: "Failed to fetch user data." });
        }
      } catch (error) {
        setMessage({ type: "error", text: "An unexpected error occurred." });
      }
    };

    fetchUser();
  }, []);

  // Update user profile
  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const token = localStorage.getItem("jwt");
      const userId = parseJwt(token).user_id;

      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/users`,
        {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            id: userId,
            username: formData.username,
            email: formData.email,
          }),
        }
      );

      if (response.ok) {
        setMessage({ type: "success", text: "Profile updated successfully." });
      } else {
        const errorData = await response.json();
        setMessage({ type: "error", text: `Error: ${errorData.message}` });
      }
    } catch (error) {
      setMessage({ type: "error", text: "An unexpected error occurred." });
    }
  };

  // Handle form input changes
  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((prevData) => ({
      ...prevData,
      [name]: value,
    }));
  };

  // Sign-out logic
  const handleSignOut = () => {
    localStorage.removeItem("jwt");
    router.push("/login"); // Redirect to login page
  };

  // Utility to decode JWT
  const parseJwt = (token) => {
    try {
      return JSON.parse(atob(token.split(".")[1]));
    } catch (e) {
      return null;
    }
  };
    
  return (
    <div className="min-h-screen bg-white flex flex-col items-center justify-start px-4 py-5">
  <NavBar />

  <div className="w-full max-w-lg items-start py-2">
    <h2 className="text-xl font-bold text-black tracking-wide text-left py-4">profile</h2>

    {/* Update Profile Form */}
    <form onSubmit={handleSubmit} className="w-full space-y-4">
      <div>
        <label
          htmlFor="username"
          className="block text-sm font-bold text-black"
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
          className="mt-1 block w-full border-gray-300 shadow-sm focus:ring-black focus:border-black sm:text-sm text-black"
        />
      </div>

      <div>
        <label htmlFor="email" className="block text-sm font-bold text-black">
          email
        </label>
        <input
          type="email"
          id="email"
          name="email"
          value={formData.email}
          onChange={handleChange}
          required
          className="mt-1 block w-full border-gray-300 shadow-sm focus:ring-black focus:border-black sm:text-sm text-black"
        />
      </div>

      {/* Message Display */}
      {message && (
        <div
          className={`mb-4 text-xs text-center ${
            message.type === "success" ? "text-green-700" : "text-red-700"
          }`}
        >
          {message.text}
        </div>
      )}

      <div className="flex justify-center">
        <button
          type="submit"
          className="mt-2 w-1/2 max-w-sm py-2 px-4 text-sm bg-white text-black font-medium hover:bg-black hover:text-white focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-black"
        >
          update profile
        </button>
      </div>
    </form>

    {/* Sign-Out Button */}
    <div className="flex justify-center">
      <button
        onClick={handleSignOut}
        className="mt-2 w-1/2 max-w-sm py-2 px-4 text-sm bg-white text-black font-medium hover:bg-red-600 hover:text-white focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-600"
      >
        sign out
      </button>
    </div>
  </div>
</div>
  );
}