import React from 'react';
import PostSkeleton from '../components/posts/PostSkeleton';

/**
 * Configuration object for skeleton generation
 */
export const SKELETON_CONFIG = {
  initialCount: 5,        // Number of skeletons for initial load
  paginationCount: 2,     // Number of skeletons for pagination
  animationDelay: 100     // Staggered animation delay in milliseconds
};

/**
 * Generates an array of PostSkeleton components with unique keys
 * @param {number} count - Number of skeleton components to generate
 * @param {Object} options - Configuration options
 * @param {number} options.animationDelay - Delay between skeleton animations (ms)
 * @param {string} options.keyPrefix - Prefix for component keys
 * @returns {Array} Array of PostSkeleton React components
 */
export const generateSkeletonArray = (count = SKELETON_CONFIG.initialCount, options = {}) => {
  const {
    animationDelay = SKELETON_CONFIG.animationDelay,
    keyPrefix = 'skeleton'
  } = options;

  return Array.from({ length: count }, (_, index) => {
    const key = `${keyPrefix}-${index}`;
    const delay = animationDelay > 0 ? index * animationDelay : 0;
    
    return React.createElement(PostSkeleton, {
      key,
      style: delay > 0 ? { animationDelay: `${delay}ms` } : undefined
    });
  });
};

/**
 * Generates skeleton array for initial loading state
 * @param {Object} options - Configuration options
 * @returns {Array} Array of PostSkeleton components for initial load
 */
export const generateInitialSkeletons = (options = {}) => {
  return generateSkeletonArray(SKELETON_CONFIG.initialCount, {
    keyPrefix: 'initial-skeleton',
    ...options
  });
};

/**
 * Generates skeleton array for pagination loading state
 * @param {Object} options - Configuration options
 * @returns {Array} Array of PostSkeleton components for pagination
 */
export const generatePaginationSkeletons = (options = {}) => {
  return generateSkeletonArray(SKELETON_CONFIG.paginationCount, {
    keyPrefix: 'pagination-skeleton',
    ...options
  });
};