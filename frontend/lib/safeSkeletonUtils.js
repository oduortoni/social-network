import React from 'react';
import SkeletonErrorBoundary from '../components/common/SkeletonErrorBoundary';
import { isSkeletonLoadingAvailable } from './skeletonHealthCheck';
import { 
  generateInitialSkeletons as _generateInitialSkeletons,
  generatePaginationSkeletons as _generatePaginationSkeletons,
  generateSkeletonArray as _generateSkeletonArray
} from './skeletonUtils';

/**
 * Safe wrapper for generateInitialSkeletons with error handling
 * @param {Object} options - Configuration options
 * @returns {React.ReactNode} Skeleton components or fallback
 */
export const safeGenerateInitialSkeletons = (options = {}) => {
  const fallback = (
    <div className="flex justify-center items-center py-8">
      <div style={{ color: 'var(--secondary-text)' }}>Loading posts...</div>
    </div>
  );

  try {
    const skeletons = _generateInitialSkeletons(options);
    return (
      <SkeletonErrorBoundary 
        fallback={fallback}
        loadingText="Loading posts..."
      >
        <div className="space-y-4">
          {skeletons}
        </div>
      </SkeletonErrorBoundary>
    );
  } catch (error) {
    console.warn('Failed to generate initial skeletons, using fallback:', error);
    return fallback;
  }
};

/**
 * Safe wrapper for generatePaginationSkeletons with error handling
 * @param {Object} options - Configuration options
 * @returns {React.ReactNode} Skeleton components or fallback
 */
export const safeGeneratePaginationSkeletons = (options = {}) => {
  const fallback = (
    <div className="flex justify-center py-4">
      <div style={{ color: 'var(--secondary-text)' }}>Loading more posts...</div>
    </div>
  );

  try {
    const skeletons = _generatePaginationSkeletons(options);
    return (
      <SkeletonErrorBoundary 
        fallback={fallback}
        loadingText="Loading more posts..."
      >
        <div className="space-y-4">
          {skeletons}
        </div>
      </SkeletonErrorBoundary>
    );
  } catch (error) {
    console.warn('Failed to generate pagination skeletons, using fallback:', error);
    return fallback;
  }
};

/**
 * Safe wrapper for generateSkeletonArray with error handling
 * @param {number} count - Number of skeleton components to generate
 * @param {Object} options - Configuration options
 * @returns {React.ReactNode} Skeleton components or fallback
 */
export const safeGenerateSkeletonArray = (count, options = {}) => {
  try {
    const skeletons = _generateSkeletonArray(count, options);
    return (
      <SkeletonErrorBoundary 
        fallback={
          <div className="flex justify-center items-center py-8">
            <div style={{ color: 'var(--secondary-text)' }}>Loading...</div>
          </div>
        }
        loadingText="Loading..."
      >
        {skeletons}
      </SkeletonErrorBoundary>
    );
  } catch (error) {
    console.warn('Failed to generate skeleton array, using fallback:', error);
    return (
      <div className="flex justify-center items-center py-8">
        <div style={{ color: 'var(--secondary-text)' }}>Loading...</div>
      </div>
    );
  }
};

// Re-export original functions for direct use when error handling is not needed
export {
  generateInitialSkeletons,
  generatePaginationSkeletons,
  generateSkeletonArray,
  SKELETON_CONFIG
} from './skeletonUtils';