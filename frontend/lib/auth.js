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

// Post creation function
export const createPost = async (formData) => {
  try {
    const response = await fetch('http://localhost:9000/posts', {
      method: 'POST',
      credentials: 'include',
      body: formData, // FormData object
    });

    const data = await response.json();

    if (response.ok) {
      return { success: true, data };
    } else {
      return { success: false, error: data.message || 'Failed to create post' };
    }
  } catch (error) {
    console.error('Error creating post:', error);
    return { success: false, error: 'Network error occurred' };
  }
};

// Fetch posts function
export const fetchPosts = async () => {
  try {
    const response = await fetch('http://localhost:9000/posts', {
      method: 'GET',
      credentials: 'include',
    });

    const data = await response.json();

    if (response.ok) {
      return { success: true, data };
    } else {
      return { success: false, error: data.message || 'Failed to fetch posts' };
    }
  } catch (error) {
    console.error('Error fetching posts:', error);
    return { success: false, error: 'Network error occurred' };
  }
};

// Create comment function
export const createComment = async (postId, formData) => {
  try {
    const response = await fetch(`http://localhost:9000/posts/${postId}/comments`, {
      method: 'POST',
      credentials: 'include',
      body: formData, // FormData object
    });

    const data = await response.json();

    if (response.ok) {
      return { success: true, data };
    } else {
      return { success: false, error: data.message || 'Failed to create comment' };
    }
  } catch (error) {
    console.error('Error creating comment:', error);
    return { success: false, error: 'Network error occurred' };
  }
};

// Fetch comments for a post
export const fetchComments = async (postId) => {
  try {
    const response = await fetch(`http://localhost:9000/posts/${postId}/comments`, {
      method: 'GET',
      credentials: 'include',
    });

    const data = await response.json();

    if (response.ok) {
      return { success: true, data };
    } else {
      return { success: false, error: data.message || 'Failed to fetch comments' };
    }
  } catch (error) {
    console.error('Error fetching comments:', error);
    return { success: false, error: 'Network error occurred' };
  }
};

// Delete post function
export const deletePost = async (postId) => {
  try {
    const response = await fetch(`http://localhost:9000/posts/${postId}`, {
      method: 'DELETE',
      credentials: 'include',
    });

    if (response.ok) {
      return { success: true };
    } else {
      const errorData = await response.json();
      return { success: false, error: errorData.message || 'Failed to delete post' };
    }
  } catch (error) {
    console.error('Error deleting post:', error);
    return { success: false, error: 'Network error occurred' };
  }
};

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
        const error = await response.json();
        setFormError("Registration failed: " + error.message);
        return;
      }
  
      const result = await response.json();
      setFormError("Registration successful!");
      // Redirect or reset form
      window.location.href="/"
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
     // console.error("Login error:", error);
      setFormError(error.message || "Login failed");
      return { success: false, error: error.message || "Login failed" };
    }

    const result = await response.json();
    
    // Clear any previous errors
    setFormError("");
    
    return { success: true, data: result };
    
  } catch (err) {
    console.error("Login error:", err);
    setFormError("Something went wrong. Please try again later.");
    return { success: false, error: "Something went wrong. Please try again later." };
  }
};

export const handleLogout = async () => {
  try {
    const response = await fetch("http://localhost:9000/logout", {
      method: "POST",
      credentials: "include", 
    });

    if (!response.ok) {
      throw new Error("Logout failed");
    }

    window.location.href = "/";
  } catch (error) {
    console.error("Logout error:", error);
  }
};


export const validateStepOne = async (userEmail, userPassword, confirmPassword) => {
  try {
    const response = await fetch("http://localhost:9000/validate/step1", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        email: userEmail,
        password: userPassword,
        confirm_password: confirmPassword,
      }),
    });

    if (!response.ok) {
      const errorData = await response.json();
      return {
        status: "failed",
        errorMessage: errorData.message,
      };
    }

    return {
      status: "success",
    };
  } catch (error) {
    console.error(error);
    return {
      status: "error",
      errorMessage: "Something went wrong. Please try again.",
    };
  }
};


export const validateStepTwo=(userFirstName,UserLastName,UserDateOfBirth)=>{  

  //check if user filled firstname, last name and date_of_birth
  if (!userFirstName || !UserLastName || !UserDateOfBirth){
       return {
        status:"failed",
        errorMessage:"Please fill in all required fields."
       }
  }

  // Validate DOB is not in the future
    const dateOfBirth = new Date(UserDateOfBirth);
    if (dateOfBirth > new Date()) {
      return {
        status:"failed",
        errorMessage:"Date of birth cannot be in the future."
      }
    }
   return {
    status:"success"
   }
}
