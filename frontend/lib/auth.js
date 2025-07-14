import { useState } from 'react';

// Function to handle form input changes
export const handleRegistrationFormChange = (e) => {
    const { name, value, files } = e.target;
    setFormData({
      ...formData,
      [name]: files ? files[0] : value
    });
  };

  // Functions to toggle form steps
export const registrationFormNextStep = () => setStep((prev) => prev + 1);
export const registrationFormPrevStep = () => setStep((prev) => prev - 1);

  // Function to handle form submission
export const handleRegistrationFormSubmit = async (e, formData, setFormError) => {
    e.preventDefault();
  
    // Validate required fields
    const {
      email,
      password,
      confirmPassword,
      firstName,
      lastName,
      dob,
      avatar,
      nickname,
      aboutMe
    } = formData;
  
    const requiredFields = [email, password, confirmPassword, firstName, lastName, dob];
    if (requiredFields.some((field) => !field || (typeof field === "string" && !field.trim()))) {
      setFormError("Please fill in all required fields.");
      return;
    }
  
    // Validate email format
    const isValidEmail = /\S+@\S+\.\S+/.test(email);
    if (!isValidEmail) {
      setFormError("Please enter a valid email address.");
      return;
    }
  
    // Validate password match
    if (password !== confirmPassword) {
      setFormError("Passwords do not match.");
      return;
    }
  
    // Validate password length
    if (password.length < 6) {
      setFormError("Password must be at least 6 characters.");
      return;
    }
  
    // Validate DOB is not in the future
    const dateOfBirth = new Date(dob);
    if (dateOfBirth > new Date()) {
      setFormError("Date of birth cannot be in the future.");
      return;
    }
  
    // Validate avatar type if uploaded
    if (avatar) {
      const allowedTypes = ["image/png", "image/jpeg", "image/gif"];
      if (!allowedTypes.includes(avatar.type)) {
        setFormError("Only PNG, JPEG, or GIF images are allowed for the avatar.");
        return;
      }
    }
  
    // Prepare FormData for submission
    const form = new FormData();
    form.append("email", email);
    form.append("password", password);
    form.append("firstName", firstName.trim());
    form.append("lastName", lastName.trim());
    form.append("dob", dob); // In 'yyyy-mm-dd' format
    if (avatar) form.append("avatar", avatar);
    if (nickname) form.append("nickname", nickname.trim());
    if (aboutMe) form.append("aboutMe", aboutMe.trim());
  
    try {
      // Send to your Go backend
      const response = await fetch("http://localhost:9000/register", {
        method: "POST",
        body: form,
      });
  
      if (!response.ok) {
        const error = await response.text();
        setFormError("Registration failed: " + error);
        return;
      }
  
      const result = await response.json();
      setFormError("Registration successful!");
      // Redirect or reset form
    } catch (err) {
      console.error("Registration error:", err);
      setFormError("Something went wrong. Please try again later.");
    }
  };

export const handleLoginFormSubmit = async (e, formData, setFormError) => {
  e.preventDefault();
  const { email, password } = formData;
  
  // Add proper email validation
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  const isValidEmail = emailRegex.test(email);
  
  if (!isValidEmail) {
    setFormError("Please enter a valid email address.");
    return { success: false, error: "Invalid email" };
  }

  const user = {
    email: email,
    password: password
  };

  try {
    // Send to your Go backend - include credentials for cookies
    const response = await fetch("http://localhost:9000/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include", // Important for cookie-based auth
      body: JSON.stringify(user),
    });

    if (!response.ok) {
      const error = await response.json();
      console.error("Login error:", error);
      setFormError(error.message || "Login failed");
      return { success: false, error: error.message || "Login failed" };
    }

    const result = await response.json();
    console.log("Login successful:", result);
    
    // Clear any previous errors
    setFormError("");
    
    return { success: true, data: result };
    
  } catch (err) {
    console.error("Login error:", err);
    setFormError("Something went wrong. Please try again later.");
    return { success: false, error: "Something went wrong. Please try again later." };
  }
};
