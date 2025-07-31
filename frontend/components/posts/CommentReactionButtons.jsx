import React, { useState } from 'react';
import { ThumbsUpIcon, ThumbsDownIcon } from 'lucide-react';
import { reactToComment, unreactToComment } from '../../lib/auth';

const CommentReactionButtons = ({ comment, user }) => {
  const [likes, setLikes] = useState(comment.likes_count || 0);
  const [dislikes, setDislikes] = useState(comment.dislikes_count || 0);
  const [userReaction, setUserReaction] = useState(comment.user_reaction || null);

  const handleReaction = async (reactionType) => {
    const newReaction = userReaction === reactionType ? null : reactionType;

    try {
      if (newReaction) {
        await reactToComment(comment.id, newReaction);
      } else {
        await unreactToComment(comment.id);
      }

      setUserReaction(newReaction);

      if (newReaction === 'like') {
        setLikes(likes + 1);
        if (userReaction === 'dislike') {
          setDislikes(dislikes - 1);
        }
      } else if (newReaction === 'dislike') {
        setDislikes(dislikes + 1);
        if (userReaction === 'like') {
          setLikes(likes - 1);
        }
      } else {
        if (userReaction === 'like') {
          setLikes(likes - 1);
        } else if (userReaction === 'dislike') {
          setDislikes(dislikes - 1);
        }
      }

    } catch (error) {
      console.error(`Failed to ${newReaction ? '' : 'un'}react to comment:`, error);
    }
  };

  return (
    <div className="flex items-center gap-4">
      <button
        onClick={() => handleReaction('like')}
        className="flex items-center gap-2 text-sm py-1.5 px-3 rounded-lg transition-colors"
        style={{ color: userReaction === 'like' ? 'var(--primary-accent)' : 'var(--secondary-text)' }}
        onMouseOver={(e) => e.currentTarget.style.backgroundColor = 'var(--hover-background)'}
        onMouseOut={(e) => e.currentTarget.style.backgroundColor = 'transparent'}
      >
        <ThumbsUpIcon className="w-4 h-4" />
        <span>{likes}</span>
      </button>
      <button
        onClick={() => handleReaction('dislike')}
        className="flex items-center gap-2 text-sm py-1.5 px-3 rounded-lg transition-colors"
        style={{ color: userReaction === 'dislike' ? 'var(--warning-color)' : 'var(--secondary-text)' }}
        onMouseOver={(e) => e.currentTarget.style.backgroundColor = 'var(--hover-background)'}
        onMouseOut={(e) => e.currentTarget.style.backgroundColor = 'transparent'}
      >
        <ThumbsDownIcon className="w-4 h-4" />
        <span>{dislikes}</span>
      </button>
    </div>
  );
};

export default CommentReactionButtons;
