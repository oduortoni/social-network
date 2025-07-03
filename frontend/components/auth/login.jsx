import Head from 'next/head';
import { FaFacebookF, FaGooglePlusG, FaLinkedinIn, FaEye, FaEyeSlash } from 'react-icons/fa';
import { useState } from 'react';

export function LoginForm() {
    const [showPassword, setShowPassword] = useState(false);

    return (
        <>
           <Head>
            <title>Social Network - Login</title>
          </Head>
          <form className="bg-white h-full px-[50px] flex flex-col justify-center items-center text-center">
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
              <input
                type="email"
                placeholder="Email"
                className="bg-[var(--tertiary-background)] text-[var(--quaternary-text)] p-3 my-2 w-full outline-none"
              />
              <input
                type={showPassword ? "text" : "password"}
                placeholder="Password"
                className="bg-[var(--tertiary-background)] text-[var(--quaternary-text)] p-3 my-2 w-full outline-none"
              />
              <button
                type="button"
                onClick={() => setShowPassword((prev) => !prev)}
                className="absolute right-3 top-1/2 transform -translate-y-1/2 text-[var(--quaternary-text)]"
                tabIndex={-1}
              >
                {showPassword ? <FaEyeSlash /> : <FaEye />}
              </button>
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