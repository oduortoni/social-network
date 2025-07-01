import { FaFacebookF, FaGooglePlusG, FaLinkedinIn } from 'react-icons/fa';

export function LoginForm() {
    return (
        <>
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
        </>
    )
}