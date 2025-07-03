import { useState } from 'react';
import Head from 'next/head';
import { FaFacebookF, FaGooglePlusG, FaLinkedinIn, FaArrowLeft, FaArrowRight } from 'react-icons/fa';
import { handleRegistrationFormSubmit } from '../../lib/auth';

export function RegisterForm() {
    // State to manage steps and form data
      const [step, setStep] = useState(1);
      const [formError, setFormError] = useState("");
      const [formData, setFormData] = useState({
        email: '',
        password: '',
        confirmPassword: '',
        firstName: '',
        lastName: '',
        dob: '',
        avatar: null,
        nickname: '',
        aboutMe: ''
      });
    
    // Function to handle form input changes
    const handleChange = (e) => {
    const { name, value, files } = e.target;
      setFormData({
        ...formData,
        [name]: files ? files[0] : value
      });
    };
    // Functions to toggle form steps
    const nextStep = () => setStep((prev) => prev + 1);
    const prevStep = () => setStep((prev) => prev - 1);
    
    return (
        <>
         <Head>
        <title>Social Network - Register</title>
        </Head>
        
        <form onSubmit={e => handleRegistrationFormSubmit(e, formData, setFormError)} className="bg-white h-full px-[50px] flex flex-col justify-center items-center text-center">
            {/* Third-Party Authentication */}
          <h1 className="font-bold text-2xl text-[var(--tertiary-text)]">Create Account</h1>
          <div className="my-5 flex space-x-3">
            <a className="border bg-[var(--primary-accent)] border-[var(--primary-accent)] rounded-full h-10 w-10 flex justify-center items-center text-[var(--quinternary-text)]">
              <FaFacebookF />
            </a>
            <a className="border bg-[var(--primary-accent)] border-[var(--primary-accent)] rounded-full h-10 w-10 flex justify-center items-center text-[var(--quinternary-text)]">
              <FaGooglePlusG />
            </a>
            <a className="border bg-[var(--primary-accent)] border-[var(--primary-accent)] rounded-full h-10 w-10 flex justify-center items-center text-[var(--quinternary-text)]">
              <FaLinkedinIn />
            </a>
          </div>

          <span className="text-xs mb-2 text-[var(--tertiary-text)]">or use your email for registration</span>

            {/* Error Message Display*/}
          {formError && (
             <div className="font-bold text-base text-[var(--warning-color)]">
                {formError}
              </div>
            )}

            {/* Step Forms */}
          {step === 1 && (
            <>
              <input 
                type="email"
                name="email"
                placeholder="Email"
                className="bg-[var(--tertiary-background)] text-[var(--quaternary-text)] p-3 my-2 w-full outline-none"
                onChange={handleChange} required
                />
              <input
                type="password"
                name="password"
                placeholder="Password"
                className="bg-[var(--tertiary-background)] text-[var(--quaternary-text)] p-3 my-2 w-full outline-none"
                onChange={handleChange}
                required
                 />
              <input
                type="password"
                name="confirmPassword"
                placeholder="Confirm Password"
                className="bg-[var(--tertiary-background)] text-[var(--quaternary-text)] p-3 my-2 w-full outline-none"
                onChange={handleChange}
                required
                />
            </>
          )}

          {step === 2 && (
            <>
              <input
                type="text"
                name="firstName"
                placeholder="First Name"
                className="bg-[var(--tertiary-background)] text-[var(--quaternary-text)] p-3 my-2 w-full outline-none"
                onChange={handleChange}
                required
                />
              <input
                type="text"
                name="lastName"
                placeholder="Last Name"
                className="bg-[var(--tertiary-background)] text-[var(--quaternary-text)] p-3 my-2 w-full outline-none"
                onChange={handleChange}
                required
                />
              <input
                type="date"
                name="dob"
                className="bg-[var(--tertiary-background)] text-[var(--quaternary-text)] p-3 my-2 w-full outline-none"
                onChange={handleChange}
                required
                />
            </>
          )}

          {step === 3 && (
            <>
              <input
                type="file"
                name="avatar"
                accept="image/png, image/jpeg, image/gif"
                className="bg-[var(--tertiary-background)] text-[var(--quaternary-text)] p-3 my-2 w-full outline-none"
                onChange={handleChange} />
              <input
                type="text"
                name="nickname"
                placeholder="Nickname (Optional)"
                className="bg-[var(--tertiary-background)] text-[var(--quaternary-text)] p-3 my-2 w-full outline-none"
                onChange={handleChange}
                />
              <textarea
                name="aboutMe"
                placeholder="About Me (Optional)"
                className="bg-[var(--tertiary-background)] text-[var(--quaternary-text)] p-3 my-2 w-full outline-none"
                onChange={handleChange}
                ></textarea>
            </>
          )}

            {/* Step Controls */}
          <div className="flex gap-4 mt-4">
            {step > 1 && (
              <button
                type="button"
                onClick={prevStep}
                className='text-[var(--tertiary-text)] hover:scale-95 transition-transform'
                >
                  <FaArrowLeft className='inline mr-1'/>
                  {' Back'}
                  </button>)}
            {step < 3 && (
              <button
                type="button"
                onClick={nextStep}
                className='text-[var(--tertiary-text)] hover:scale-95 transition-transform'
                >
                  {'Next '}
                  <FaArrowRight className='inline mr-1'/>
                  </button>)}
            {step === 3 && <button type="submit" className='text-[var(--tertiary-text)] hover:scale-95 transition-transform'>Register</button>}
          </div>
            
        </form>
        </>
    )
}
