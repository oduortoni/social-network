import { useState } from 'react';
import Head from 'next/head';
import { FaFacebookF, FaGooglePlusG, FaLinkedinIn, FaArrowLeft, FaArrowRight, FaEye, FaEyeSlash } from 'react-icons/fa';
import { handleRegistrationFormSubmit } from '../../lib/auth';

export function RegisterForm() {
    // State to manage registrationform data
      const [step, setStep] = useState(1);
      const [showPassword, setShowPassword] = useState(false);
      const [showConfirmPassword, setShowConfirmPassword] = useState(false);
      const [formError, setFormError] = useState("");
      const [avatarPreview, setAvatarPreview] = useState(null);
      const [formData, setFormData] = useState({
        email: '',
        password: '',
        confirmPassword: '',
        firstName: '',
        lastName: '',
        dob: '',
        profileVisibility: '',
        avatar: null,
        nickname: '',
        aboutMe: ''
      });
    
    // Function to handle form input changes
    const handleChange = (e) => {
      const { name, value, files } = e.target;
      if (name === "avatar" && files && files[0]) {
        setFormData({
          ...formData,
          [name]: files[0]
        });
        setAvatarPreview(URL.createObjectURL(files[0]));
      } else {
        setFormData({
          ...formData,
          [name]: value
        });
      }
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
                {/* Email Input */}
              <input 
                type="email"
                name="email"
                placeholder="Email"
                className="bg-[var(--tertiary-background)] text-[var(--quaternary-text)] p-3 my-2 w-full outline-none"
                onChange={handleChange} required
                />

                {/* Password Input */}
              <div className="relative w-full">
                <input
                  type={showPassword ? "text" : "password"}
                  name="password"
                  placeholder="Password"
                  className="bg-[var(--tertiary-background)] text-[var(--quaternary-text)] p-3 my-2 w-full outline-none pr-10"
                  onChange={handleChange}
                  required
                  />
                  {/* Visibility Toggle Button */}
                <button
                  type="button"
                  onClick={() => setShowPassword((prev) => !prev)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-[var(--quaternary-text)]"
                  tabIndex={-1}
                >
                  {showPassword ? <FaEyeSlash /> : <FaEye />}
                </button>
              </div>

                {/* Confirm Password Input */}
              <div className="relative w-full">
                <input
                  type={showConfirmPassword ? "text" : "password"}
                  name="confirmPassword"
                  placeholder="Confirm Password"
                  className="bg-[var(--tertiary-background)] text-[var(--quaternary-text)] p-3 my-2 w-full outline-none"
                  onChange={handleChange}
                  required
                  />
                  {/* Visibility Toggle Button */}
                <button
                  type="button"
                  onClick={() => setShowConfirmPassword((prev) => !prev)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-[var(--quaternary-text)]"
                  tabIndex={-1}
                >
                  {showPassword ? <FaEyeSlash /> : <FaEye />}
                </button>
              </div>
            </>
          )}

          {/* First Name, Last Name, and Date of Birth Inputs */}
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

          {/* Profile Visibility, Nickname, and About Me Inputs */}
          {step === 3 && (
            <>
              <select
                name="profileVisibility"
                className="bg-[var(--tertiary-background)] text-[var(--quaternary-text)] p-3 my-2 w-full outline-none"
                onChange={handleChange}
                value={formData.profileVisibility}
                required
              >
                <option value="" disabled>
                  Profile Visibility
                </option>
                <option value="private">Private</option>
                <option value="public">Public</option>
              </select>
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
                maxLength={300}
                className="bg-[var(--tertiary-background)] text-[var(--quaternary-text)] p-3 my-2 w-full outline-none h-28 resize-none"
                onChange={handleChange}
                value={formData.aboutMe}
                />
                {/* Textarea Character Limit Counter */}
              <div className="w-full text-right text-xs text-[var(--tertiary-text)]">
                {300 - (formData.aboutMe?.length || 0)} characters left
              </div>
            </>
          )}

          {/* Avatar Upload */}
          {step === 4 && (
            <>
            <input
                type="file"
                name="avatar"
                accept="image/png, image/jpeg, image/gif"
                className="bg-[var(--tertiary-background)] text-[var(--quaternary-text)] p-3 my-2 w-full outline-none"
                onChange={handleChange} />
              {avatarPreview && (
                <div className="w-full flex justify-center my-2">
                  <img
                    src={avatarPreview}
                    alt="Avatar Preview"
                    className="h-28 w-28 object-cover rounded-full border border-[var(--tertiary-text)]"
                  />
                </div>
              )}
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
            {step < 4 && (
              <button
                type="button"
                onClick={nextStep}
                className='text-[var(--tertiary-text)] hover:scale-95 transition-transform'
                >
                  {'Next '}
                  <FaArrowRight className='inline mr-1'/>
                  </button>)}
            {step === 4 && <button type="submit" className='text-[var(--tertiary-text)] hover:scale-95 transition-transform'>Register</button>}
          </div>
            
        </form>
        </>
    )
}
