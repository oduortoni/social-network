import Head from 'next/head';
import { FaFacebookF, FaGooglePlusG, FaLinkedinIn, FaEye, FaEyeSlash } from 'react-icons/fa';
import { useState } from 'react';
import { handleLoginFormSubmit } from '../../lib/auth';

export function LoginForm() {
    const [showPassword, setShowPassword] = useState(false);
   const [formError, setFormError] = useState("");
   const [formData, setFormData] = useState({
     email: '',
     password: '',
   });

   const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData({
      ...formData,
      [name]: value
    });
   }
    return (
        <>
           <Head>
            <title>Social Network - Login</title>
          </Head>
          <form onSubmit={e => handleLoginFormSubmit(e, formData, setFormError)} className="bg-white h-full px-[50px] flex flex-col justify-center items-center text-center">
              <h1 className="font-bold text-2xl text-[var(--tertiary-text)]">Sign in</h1>
              <div className="my-5 flex space-x-3">
                <a className="border bg-[var(--primary-accent)] border-[var(--primary-accent)] rounded-full h-10 w-10 flex justify-center items-center text-[var(--primary-background)]">
                  <FaFacebookF />
                </a>
                <a className="border bg-[var(--primary-accent)] border-[var(--primary-accent)] rounded-full h-10 w-10 flex justify-center items-center text-[var(--primary-background)]">
                  <FaGooglePlusG />
                </a>
                <a className="border bg-[var(--primary-accent)] border-[var(--primary-accent)] rounded-full h-10 w-10 flex justify-center items-center text-[var(--primary-background)]">
                  <FaLinkedinIn />
                </a>
              </div>
              <span className="text-xs mb-2 text-[var(--tertiary-text)]">or use your account</span>
              {formError && (
                <div className="font-bold text-base text-[var(--warning-color)]">
                  {formError}
                </div>
              )}
              <input
                type="email"
                name='email'
                onChange={handleChange}
                required
                placeholder="Email"
                className="bg-[var(--tertiary-background)] text-[var(--quaternary-text)] p-3 my-2 w-full outline-none"
              />
              {/* Password and Visibility Toggle Button Container */}
              <div className="relative w-full">
                <input
                  type={showPassword ? "text" : "password"}
                  name='password'
                  onChange={handleChange}
                  required
                  placeholder="Password"
                  className="bg-[var(--tertiary-background)] text-[var(--quaternary-text)] p-3 my-2 w-full outline-none pr-10"
                />
                <button
                  type="button"
                  onClick={() => setShowPassword((prev) => !prev)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-[var(--quaternary-text)]"
                  tabIndex={-1}
                >
                  {showPassword ? <FaEyeSlash /> : <FaEye />}
                </button>
              </div>
              <a href="#" className="text-sm text-[var(--tertiary-text)] mt-2 mb-4">
                Forgot your password?
              </a>
              <button className="rounded-full border border-[var(--tertiary-text)] bg-[var(--tertiary-text)] text-white text-xs font-bold py-3 px-11 uppercase hover:scale-95 transition-transform">
                Sign In
              </button>
          </form>
        </>
    )
}