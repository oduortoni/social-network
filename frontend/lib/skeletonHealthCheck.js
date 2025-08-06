/**
 * Health check utilities for skeleton loading functionality
 */

/**
 * Tests if skeleton loading is available and working
 * @returns {boolean} True if skeleton loading is functional
 */
export const isSkeletonLoadingAvailable = () => {
  try {
    // Basic environment checks
    if (typeof window === 'undefined') {
      return false; // Server-side rendering
    }

    // Test if CSS variables are supported (basic check)
    if (window.CSS && window.CSS.supports) {
      if (!window.CSS.supports('color', 'var(--test)')) {
        console.warn('CSS variables not supported, falling back to text loading');
        return false;
      }
    }

    // If we get here, assume skeleton loading is available
    // The animation check was too strict and causing false negatives
    return true;
  } catch (error) {
    console.warn('Skeleton loading health check failed:', error);
    return false;
  }
};

/**
 * Gets the appropriate loading component based on skeleton availability
 * @param {React.ReactNode} skeletonComponent - The skeleton component to use
 * @param {React.ReactNode} fallbackComponent - The fallback component to use
 * @returns {React.ReactNode} The appropriate loading component
 */
export const getLoadingComponent = (skeletonComponent, fallbackComponent) => {
  if (isSkeletonLoadingAvailable()) {
    return skeletonComponent;
  }
  return fallbackComponent;
};