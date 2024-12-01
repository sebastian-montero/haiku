import { useState } from "react";
import { useRouter } from "next/router";

export default function Signup() {
  const [formData, setFormData] = useState({
    username: "",
    email: "",
    password: "",
    confirmPassword: "",
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

    // Check if passwords match
    if (formData.password !== formData.confirmPassword) {
      setMessage({ type: "error", text: "passwords do not match!" });
      return;
    }

    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_BASE_URL}/signup`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            username: formData.username,
            email: formData.email,
            password: formData.password,
          }),
        },
      );

      if (response.ok) {
        const data = await response.json();

        // Store JWT in localStorage
        if (data.token) {
          localStorage.setItem("jwt", data.token);
        }

        setMessage({ type: "success", text: `welcome!` });
        router.push("/notebooks");
      } else {
        const errorData = await response.json();
        setMessage({ type: "error", text: `error: ${errorData.message}` });
      }
    } catch (error) {
      setMessage({ type: "error", text: "an unexpected error occurred." });
    }
  };

  return (
    <div className="min-h-screen bg-white flex items-center justify-center px-4">
      <div className="max-w-md w-full space-y-6">
        <h1 className="text-3xl font-bold text-black tracking-wide text-center">
          haiku⿻
        </h1>
        <h2 className="text-xl font-bold text-black tracking-wide text-center">
          sign up
        </h2>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label
              htmlFor="username"
              className="block text-sm font-medium text-black"
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
            <label
              htmlFor="email"
              className="block text-sm font-medium text-black"
            >
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

          <div>
            <label
              htmlFor="password"
              className="block text-sm font-medium text-black"
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
              className="mt-1 block w-full border-gray-300 shadow-sm focus:ring-black focus:border-black sm:text-sm text-black"
            />
          </div>

          <div>
            <label
              htmlFor="confirmPassword"
              className="block text-sm font-medium text-black"
            >
              re-enter password
            </label>
            <input
              type="password"
              id="confirmPassword"
              name="confirmPassword"
              value={formData.confirmPassword}
              onChange={handleChange}
              required
              className="mt-1 block w-full border-gray-300 shadow-sm focus:ring-black focus:border-black sm:text-sm text-black"
            />
          </div>

          {message && (
            <div
              className={`p-4 text-sm ${
                message.type === "success"
                  ? "text-green-700 text-center"
                  : "text-red-700 text-center"
              }`}
            >
              {message.text}
            </div>
          )}

          <button
            type="submit"
            className="w-full py-2 px-4 text-black font-medium hover:bg-black hover:text-white focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-black"
          >
            sign up
          </button>
        </form>
        <p className="text-xs text-center text-black">
          have an account?
          <a href="/login" className="hover:underline ml-1">
            log in
          </a>
        </p>
      </div>
    </div>
  );
}
