// components/AuthPanel.tsx

'use client';

import { useState } from 'react';
import { FaFacebookF, FaGooglePlusG, FaLinkedinIn } from 'react-icons/fa';

export function AuthPanel() {
  const [rightPanelActive, setRightPanelActive] = useState(false);

  return (
    <div
      className={`relative w-full max-w-[768px] min-h-[480px] bg-white rounded-[10px] shadow-[0_14px_28px_rgba(0,0,0,0.25),_0_10px_10px_rgba(0,0,0,0.22)] overflow-hidden transition-all duration-700 ${
        rightPanelActive ? 'right-panel-active' : ''
      }`}
    >
      {/* Sign Up Form */}
   <div
      className={`absolute top-0 left-1/2 h-full w-1/2 transition-all duration-700
        ${rightPanelActive
          ? 'opacity-100 z-10 animate-fadeIn'
          : 'opacity-0 z-0'
        }`}
    >
        <form className="bg-white h-full px-[50px] flex flex-col justify-center items-center text-center">
          <h1 className="font-bold text-2xl">Create Account</h1>
          <div className="my-5 flex space-x-3">
            <a className="border border-[var(--primary-accent)] rounded-full h-10 w-10 flex justify-center items-center text-[var(--primary-accent)]">
              <FaFacebookF />
            </a>
            <a className="border border-[var(--primary-accent)] rounded-full h-10 w-10 flex justify-center items-center text-[var(--primary-accent)]">
              <FaGooglePlusG />
            </a>
            <a className="border border-[var(--primary-accent)] rounded-full h-10 w-10 flex justify-center items-center text-[var(--primary-accent)]">
              <FaLinkedinIn />
            </a>
          </div>
          <span className="text-xs mb-2 text-[var(--tertiary-text)]">or use your email for registration</span>
          <input
            type="text"
            placeholder="Name"
            className="bg-[var(--tertiary-background)] p-3 my-2 w-full outline-none"
          />
          <input
            type="email"
            placeholder="Email"
            className="bg-[var(--tertiary-background)] p-3 my-2 w-full outline-none"
          />
          <input
            type="password"
            placeholder="Password"
            className="bg-[var(--tertiary-background)] p-3 my-2 w-full outline-none"
          />
          <button className="rounded-full border border-[var(--tertiary-text)] bg-[var(--tertiary-text)] text-white text-xs font-bold py-3 px-11 mt-4 uppercase hover:scale-95 transition-transform">
            Sign Up
          </button>
        </form>
      </div>

      {/* Sign In Form */}
      <div
        className={`absolute top-0 left-0 h-full w-1/2 transition-all duration-700
          ${rightPanelActive
            ? 'opacity-0 z-0'
            : 'opacity-100 z-10 animate-fadeIn'
          }`}
      >
        <form className="bg-white h-full px-[50px] flex flex-col justify-center items-center text-center">
          <h1 className="font-bold text-2xl text-[var(--tertiary-text)]">Sign in</h1>
          <div className="my-5 flex space-x-3">
            <a className="border bg-[var(--primary-accent)] border-[var(--primary-accent)] rounded-full h-10 w-10 flex justify-center items-center text-[var(--primary-background)]">
              <FaFacebookF />
            </a>
            <a className="border border-[var(--primary-accent)] rounded-full h-10 w-10 flex justify-center items-center text-[var(--primary-background)]">
              <FaGooglePlusG />
            </a>
            <a className="border border-[var(--primary-accent)] rounded-full h-10 w-10 flex justify-center items-center text-[var(--primary-background)]">
              <FaLinkedinIn />
            </a>
          </div>
          <span className="text-xs mb-2 text-[var(--tertiary-text)]">or use your account</span>
          <input
            type="email"
            placeholder="Email"
            className="bg-[var(--tertiary-background)] p-3 my-2 w-full outline-none"
          />
          <input
            type="password"
            placeholder="Password"
            className="bg-[var(--tertiary-background)] p-3 my-2 w-full outline-none"
          />
          <a href="#" className="text-sm text-[var(--tertiary-text)] mt-2 mb-4">
            Forgot your password?
          </a>
          <button className="rounded-full border border-[var(--tertiary-text)] bg-[var(--tertiary-text)] text-white text-xs font-bold py-3 px-11 uppercase hover:scale-95 transition-transform">
            Sign In
          </button>
        </form>
      </div>

      {/* Overlay Panel */}
      <div
        className={`absolute top-0 left-1/2 w-1/2 h-full overflow-hidden transition-transform duration-700 z-30 ${
          rightPanelActive ? '-translate-x-full' : ''
        }`}
      >
        <div
          className={`absolute left-[-100%] w-[200%] h-full bg-gradient-to-r from-[#363a80] to-[var(--secondary-background)] text-white flex transition-transform duration-700 ${
            rightPanelActive ? 'translate-x-1/2' : ''
          }`}
        >
          {/* Overlay Left */}
          <div
            className={`w-1/2 h-full flex flex-col justify-center items-center text-center px-10 transition-transform duration-700 ${
              rightPanelActive ? '' : '-translate-x-5'
            }`}
          >
            <h1 className="text-2xl font-bold">Welcome Back!</h1>
            <p className="text-sm mt-5 mb-8">
              To keep connected with us please login with your personal info
            </p>
            <button
              className="ghost bg-transparent border border-white text-white text-xs font-bold py-3 px-11 rounded-full uppercase hover:scale-95 transition-transform"
              onClick={() => setRightPanelActive(false)}
            >
              Sign In
            </button>
          </div>
        
          {/* Overlay Right */}
          <div
            className={`w-1/2 h-full flex flex-col justify-center items-center text-center px-10 transition-transform duration-700 ${
              rightPanelActive ? 'translate-x-5' : ''
            }`}
          >
            <h1 className="text-2xl font-bold">Hello, Friend!</h1>
            <p className="text-sm mt-5 mb-8">
              Enter your personal details and start journey with us
            </p>
            <button
              className="ghost bg-transparent border border-white text-white text-xs font-bold py-3 px-11 rounded-full uppercase hover:scale-95 transition-transform"
              onClick={() => setRightPanelActive(true)}
            >
              Sign Up
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
