"use client";

import { useEffect, useState } from "react";



const ContactForm = () => {
  const [formData, setFormData] = useState({
    email: '',
    password: '',
    first_name: '',
    last_name: '',
    date_of_birth: '',
    nickname: '',
    about_me: '',
    is_profile_public: false,
    avatar: null // <-- for profile photo
  });

  const handleChange = (e) => {
    const { name, value, type, checked, files } = e.target;
    setFormData((prevData) => ({
      ...prevData,
      [name]: type === 'checkbox' ? checked : type === 'file' ? files[0] : value,
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    const fd = new FormData();
    for (let key in formData) {
      if (formData[key] !== null) {
        fd.append(key, formData[key]);
      }
    }

    try {
      const response = await fetch("http://localhost:9000/register", {
        method: "POST",
        body: fd,
      });

      const result = await response.json();
      console.log(result);
    } catch (err) {
      console.error("Form submission failed:", err);
    }
  };

  return (
    <main className="flex flex-col gap-[32px] row-start-2 items-center sm:items-start">
      <h1 className="text-2xl font-bold">User Registration</h1>

      <form onSubmit={handleSubmit} className="flex flex-col gap-4 w-full max-w-sm">
        {/* Other fields */}
        <label className="flex flex-col text-sm">
          Email
          <input type="email" name="email" value={formData.email}
            onChange={handleChange} className="border rounded px-3 py-2 mt-1" required />
        </label>

        <label className="flex flex-col text-sm">
          Password
          <input type="password" name="password" value={formData.password}
            onChange={handleChange} className="border rounded px-3 py-2 mt-1" required />
        </label>

        <label className="flex flex-col text-sm">
          First Name
          <input type="text" name="first_name" value={formData.first_name}
            onChange={handleChange} className="border rounded px-3 py-2 mt-1" required />
        </label>

        <label className="flex flex-col text-sm">
          Last Name
          <input type="text" name="last_name" value={formData.last_name}
            onChange={handleChange} className="border rounded px-3 py-2 mt-1" required />
        </label>

        <label className="flex flex-col text-sm">
          Date of Birth
          <input type="date" name="date_of_birth" value={formData.date_of_birth}
            onChange={handleChange} className="border rounded px-3 py-2 mt-1" required />
        </label>

        <label className="flex flex-col text-sm">
          Nickname
          <input type="text" name="nickname" value={formData.nickname}
            onChange={handleChange} className="border rounded px-3 py-2 mt-1" />
        </label>

        <label className="flex flex-col text-sm">
          About Me
          <textarea name="about_me" rows="3" value={formData.about_me}
            onChange={handleChange} className="border rounded px-3 py-2 mt-1" />
        </label>

        <label className="flex items-center text-sm gap-2">
          <input type="checkbox" name="is_profile_public"
            checked={formData.is_profile_public}
            onChange={handleChange} />
          Make Profile Public
        </label>

        {/* 🔽 Profile Photo Upload */}
        <label className="flex flex-col text-sm">
          Profile Photo
          <input
            type="file"
            name="avatar"
            accept="image/*"
            onChange={handleChange}
            className="border rounded px-3 py-2 mt-1"
          />
        </label>

        <button type="submit"
          className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700">
          Submit
        </button>
      </form>
    </main>
  );
};


export default function Home() {
  const [data, setData] = useState(null);
  const [error, setError] = useState(null);

  useEffect(() => {
    const apiUrl = process.env.NEXT_PUBLIC_API_URL;

    fetch(`${apiUrl}/`)
      .then((res) => {
        if (!res.ok) throw new Error(`API returned ${res.status}`);
        return res.json();
      })
      .then(setData)
      .catch((err) => {
        console.error(err);
        setError("Failed to load API status");
      });
  }, []);

  return (
    <div className="grid grid-rows-[20px_1fr_20px] items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20 font-[family-name:var(--font-geist-sans)]">
          <ContactForm/>

      <footer className="row-start-3 flex gap-[24px] flex-wrap items-center justify-center">
        &copy; 2025 tajjjjr
      </footer>
    </div>
  );
}
