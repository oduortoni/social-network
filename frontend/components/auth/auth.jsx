// components/AuthPanel.jsx

'use client';

import { useState } from 'react';
import { RegisterForm } from './register';
import { LoginForm } from './login';

export function AuthPanel() {
  // State to manage the active panel and form data
  const [rightPanelActive, setRightPanelActive] = useState(false);

  return (
    <div
      className={`relative w-full max-w-[768px] min-h-[480px] bg-white rounded-[10px] shadow-[0_14px_28px_rgba(0,0,0,0.25),_0_10px_10px_rgba(0,0,0,0.22)] overflow-hidden transition-all duration-700 ${
        rightPanelActive ? 'right-panel-active' : ''
      }`}
    >
      {/* Sign Up Form */}
      <div
      className={`
        absolute top-0 left-1/2 h-full w-1/2 transition-all duration-700
        ${rightPanelActive
          ? 'translate-x-0 opacity-100 z-10'
          : 'translate-x-[100%] opacity-0 z-0'
        }
      `}
      >
        {/* Render Sign Up Form */}
        <RegisterForm />
      </div>

      {/* Sign In Form */}
       <div
      className={`
        absolute top-0 left-0 h-full w-1/2 transition-all duration-700
        ${rightPanelActive
          ? 'translate-x-[-100%] opacity-0 z-0'
          : 'translate-x-0 opacity-100 z-10'
        }
      `}
        >
        {/* Render Sign In Form */}
        <LoginForm />
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
