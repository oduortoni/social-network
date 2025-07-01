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
export const handleRegistrationFormSubmit = (e) => {
    e.preventDefault();
    // Send formData to backend
    console.log(formData);
  };