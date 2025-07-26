import React, { useState, useRef } from 'react';
import { ImageIcon, SendIcon, X, Globe, Users, Lock } from 'lucide-react';
import { createPost } from '../../lib/auth';

const PostCreation = ({ user, onPostCreated }) => {
  const [content, setContent] = useState('');
  const [privacy, setPrivacy] = useState('public');
  const [selectedImage, setSelectedImage] = useState(null);
  const [imagePreview, setImagePreview] = useState(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState('');
  const [showPrivacyDropdown, setShowPrivacyDropdown] = useState(false);
  const fileInputRef = useRef(null);

  const privacyOptions = [
    { value: 'public', label: 'Public', icon: Globe, description: 'Anyone can see this post' },
    { value: 'almost_private', label: 'Followers', icon: Users, description: 'Only your followers can see this' },
    { value: 'private', label: 'Private', icon: Lock, description: 'Only specific people can see this' }
  ];

  const handleImageSelect = (e) => {
    const file = e.target.files[0];
    if (file) {
      // Validate file type
      const allowedTypes = ['image/jpeg', 'image/png', 'image/gif'];
      if (!allowedTypes.includes(file.type)) {
        setError('Only JPEG, PNG, and GIF images are allowed');
        return;
      }

      // Validate file size (20MB limit)
      if (file.size > 20 * 1024 * 1024) {
        setError('Image size must be less than 20MB');
        return;
      }

      setSelectedImage(file);
      setImagePreview(URL.createObjectURL(file));
      setError('');
    }
  };

  const removeImage = () => {
    setSelectedImage(null);
    setImagePreview(null);
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!content.trim()) {
      setError('Post content is required');
      return;
    }

    setIsSubmitting(true);
    setError('');

    try {
      const formData = new FormData();
      formData.append('content', content.trim());
      formData.append('privacy', privacy);
      
      if (selectedImage) {
        formData.append('image', selectedImage);
      }

      const result = await createPost(formData);
      
      if (result.success) {
        // Reset form
        setContent('');
        setPrivacy('public');
        setSelectedImage(null);
        setImagePreview(null);
        if (fileInputRef.current) {
          fileInputRef.current.value = '';
        }
        
        // Notify parent component
        if (onPostCreated) {
          onPostCreated(result.data);
        }
      } else {
        setError(result.error);
      }
    } catch (error) {
      setError('Failed to create post. Please try again.');
    } finally {
      setIsSubmitting(false);
    }
  };

  const getCurrentPrivacyOption = () => {
    return privacyOptions.find(option => option.value === privacy);
  };

  return (
    <div className="rounded-xl p-4 mb-4" style={{ backgroundColor: 'var(--primary-background)' }}>
      <form onSubmit={handleSubmit}>
        {/* Post Input Area */}
        <div className="flex items-start gap-3 rounded-xl p-3 mb-4" style={{ backgroundColor: 'var(--secondary-background)' }}>
          <img
            src={user?.avatar && user.avatar !== "no profile photo" ? `http://localhost:9000/avatar?avatar=${user.avatar}` : "http://localhost:9000/avatar?avatar=user-profile-circle-svgrepo-com.svg"}
            alt="Profile"
            className="w-10 h-10 rounded-full flex-shrink-0"
          />
          <div className="flex-1">
            <textarea
              value={content}
              onChange={(e) => setContent(e.target.value)}
              placeholder="Tell your friends about your thoughts..."
              className="bg-transparent w-full focus:outline-none text-sm resize-none min-h-[60px] text-white placeholder-gray-400"
              rows="3"
              disabled={isSubmitting}
            />
            
            {/* Image Preview */}
            {imagePreview && (
              <div className="relative mt-3 inline-block">
                <img 
                  src={imagePreview} 
                  alt="Preview" 
                  className="max-w-full max-h-48 rounded-lg"
                />
                <button
                  type="button"
                  onClick={removeImage}
                  className="absolute top-2 right-2 bg-black bg-opacity-50 text-white rounded-full p-1 hover:bg-opacity-70"
                  disabled={isSubmitting}
                >
                  <X className="w-4 h-4" />
                </button>
              </div>
            )}
          </div>
        </div>

        {/* Error Message */}
        {error && (
          <div className="mb-4 p-3 bg-red-500 bg-opacity-20 border border-red-500 rounded-lg text-red-300 text-sm">
            {error}
          </div>
        )}

        {/* Privacy Selector and Action Buttons */}
        <div className="flex justify-between items-center border-t border-[#3f3fd3]/30 pt-3">
          <div className="flex items-center gap-2">
            {/* Image Upload Button */}
            <button 
              type="button"
              onClick={() => fileInputRef.current?.click()}
              className="flex items-center gap-2 text-sm py-1.5 px-3 rounded-lg cursor-pointer hover:bg-[#3f3fd3]/20 transition-colors"
              style={{ color: 'var(--secondary-text)' }}
              disabled={isSubmitting}
            >
              <ImageIcon className="w-4 h-4" />
              <span>Photo</span>
            </button>
            
            {/* Hidden File Input */}
            <input
              ref={fileInputRef}
              type="file"
              accept="image/jpeg,image/png,image/gif"
              onChange={handleImageSelect}
              className="hidden"
              disabled={isSubmitting}
            />

            {/* Privacy Selector */}
            <div className="relative">
              <button
                type="button"
                onClick={() => setShowPrivacyDropdown(!showPrivacyDropdown)}
                className="flex items-center gap-2 text-sm py-1.5 px-3 rounded-lg cursor-pointer hover:bg-[#3f3fd3]/20 transition-colors"
                style={{ color: 'var(--secondary-text)' }}
                disabled={isSubmitting}
              >
                {React.createElement(getCurrentPrivacyOption().icon, { className: "w-4 h-4" })}
                <span>{getCurrentPrivacyOption().label}</span>
              </button>

              {/* Privacy Dropdown */}
              {showPrivacyDropdown && (
                <div className="absolute bottom-full left-0 mb-2 bg-gray-800 rounded-lg shadow-lg border border-gray-600 min-w-[200px] z-10">
                  {privacyOptions.map((option) => (
                    <button
                      key={option.value}
                      type="button"
                      onClick={() => {
                        setPrivacy(option.value);
                        setShowPrivacyDropdown(false);
                      }}
                      className={`w-full text-left px-3 py-2 hover:bg-gray-700 first:rounded-t-lg last:rounded-b-lg transition-colors ${
                        privacy === option.value ? 'bg-[#3f3fd3]/20' : ''
                      }`}
                      disabled={isSubmitting}
                    >
                      <div className="flex items-center gap-2 mb-1">
                        {React.createElement(option.icon, { className: "w-4 h-4" })}
                        <span className="text-sm font-medium text-white">{option.label}</span>
                      </div>
                      <div className="text-xs text-gray-400 ml-6">{option.description}</div>
                    </button>
                  ))}
                </div>
              )}
            </div>
          </div>

          {/* Post Button */}
          <button
            type="submit"
            disabled={!content.trim() || isSubmitting}
            className={`flex items-center gap-2 text-sm py-2 px-4 rounded-lg transition-all ${
              !content.trim() || isSubmitting
                ? 'bg-gray-600 text-gray-400 cursor-not-allowed'
                : 'bg-[#3f3fd3] text-white hover:bg-[#3f3fd3]/80 cursor-pointer'
            }`}
          >
            <SendIcon className="w-4 h-4" />
            <span>{isSubmitting ? 'Posting...' : 'Post'}</span>
          </button>
        </div>
      </form>
    </div>
  );
};

export default PostCreation;
