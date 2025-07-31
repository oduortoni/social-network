import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import CommentList from '../components/posts/CommentList';
import * as auth from '../lib/auth';

// Mock the auth module
jest.mock('../lib/auth', () => ({
  fetchComments: jest.fn(),
  updateComment: jest.fn(),
  deleteComment: jest.fn(),
  followUser: jest.fn(),
}));

const mockUser = {
  id: 1,
  first_name: 'John',
  last_name: 'Doe',
  nickname: 'johndoe',
  avatar: 'avatar.jpg'
};

const mockComments = [
  {
    id: 1,
    post_id: 1,
    user_id: 1,
    content: 'Test comment 1',
    image: null,
    created_at: '2024-01-01T10:00:00Z',
    updated_at: null,
    is_edited: false,
    author: mockUser
  },
  {
    id: 2,
    post_id: 1,
    user_id: 2,
    content: 'Test comment 2',
    image: null,
    created_at: '2024-01-01T11:00:00Z',
    updated_at: '2024-01-01T12:00:00Z',
    is_edited: true,
    author: {
      id: 2,
      first_name: 'Jane',
      last_name: 'Smith',
      nickname: 'janesmith',
      avatar: 'avatar2.jpg'
    }
  }
];

describe('CommentList', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('Comment Loading', () => {
    test('displays loading state initially', () => {
      auth.fetchComments.mockImplementation(() => new Promise(() => {}));
      
      render(<CommentList postId={1} user={mockUser} />);
      
      expect(screen.getByText('Loading comments...')).toBeInTheDocument();
    });

    test('displays comments after successful fetch', async () => {
      auth.fetchComments.mockResolvedValue({
        success: true,
        data: mockComments
      });

      render(<CommentList postId={1} user={mockUser} />);

      await waitFor(() => {
        expect(screen.getByText('Test comment 1')).toBeInTheDocument();
        expect(screen.getByText('Test comment 2')).toBeInTheDocument();
      });
    });

    test('displays error message on fetch failure', async () => {
      auth.fetchComments.mockResolvedValue({
        success: false,
        error: 'Failed to load comments'
      });

      render(<CommentList postId={1} user={mockUser} />);

      await waitFor(() => {
        expect(screen.getByText('Failed to load comments')).toBeInTheDocument();
      });
    });

    test('displays no comments message when empty', async () => {
      auth.fetchComments.mockResolvedValue({
        success: true,
        data: []
      });

      render(<CommentList postId={1} user={mockUser} />);

      await waitFor(() => {
        expect(screen.getByText('No comments yet. Be the first to comment!')).toBeInTheDocument();
      });
    });
  });

  describe('Dropdown Menu Interactions', () => {
    beforeEach(async () => {
      auth.fetchComments.mockResolvedValue({
        success: true,
        data: mockComments
      });
    });

    test('shows dropdown menu when three-dot button is clicked', async () => {
      render(<CommentList postId={1} user={mockUser} />);

      await waitFor(() => {
        expect(screen.getByText('Test comment 1')).toBeInTheDocument();
      });

      const dropdownButtons = screen.getAllByRole('button');
      const threeDotsButton = dropdownButtons.find(btn => 
        btn.querySelector('svg')?.classList.contains('lucide-more-horizontal')
      );
      
      fireEvent.click(threeDotsButton);

      expect(screen.getByText('Edit')).toBeInTheDocument();
      expect(screen.getByText('Delete')).toBeInTheDocument();
    });

    test('shows only Follow option for other users comments', async () => {
      render(<CommentList postId={1} user={mockUser} />);

      await waitFor(() => {
        expect(screen.getByText('Test comment 2')).toBeInTheDocument();
      });

      const dropdownButtons = screen.getAllByRole('button');
      const secondCommentButton = dropdownButtons[1];
      
      fireEvent.click(secondCommentButton);

      expect(screen.getByText('Follow')).toBeInTheDocument();
      expect(screen.queryByText('Edit')).not.toBeInTheDocument();
      expect(screen.queryByText('Delete')).not.toBeInTheDocument();
    });

    test('closes dropdown when clicking outside', async () => {
      render(<CommentList postId={1} user={mockUser} />);

      await waitFor(() => {
        expect(screen.getByText('Test comment 1')).toBeInTheDocument();
      });

      const dropdownButtons = screen.getAllByRole('button');
      const threeDotsButton = dropdownButtons[0];
      
      fireEvent.click(threeDotsButton);
      expect(screen.getByText('Edit')).toBeInTheDocument();

      fireEvent.mouseDown(document.body);
      expect(screen.queryByText('Edit')).not.toBeInTheDocument();
    });
  });

  describe('Edit Modal Functionality', () => {
    beforeEach(async () => {
      auth.fetchComments.mockResolvedValue({
        success: true,
        data: mockComments
      });
    });

    test('opens edit modal when Edit is clicked', async () => {
      render(<CommentList postId={1} user={mockUser} />);

      await waitFor(() => {
        expect(screen.getByText('Test comment 1')).toBeInTheDocument();
      });

      const dropdownButtons = screen.getAllByRole('button');
      const threeDotsButton = dropdownButtons[0];
      
      fireEvent.click(threeDotsButton);
      fireEvent.click(screen.getByText('Edit'));

      expect(screen.getByText('Edit Comment')).toBeInTheDocument();
      expect(screen.getByDisplayValue('Test comment 1')).toBeInTheDocument();
    });

    test('validates empty content in edit modal', async () => {
      render(<CommentList postId={1} user={mockUser} />);

      await waitFor(() => {
        expect(screen.getByText('Test comment 1')).toBeInTheDocument();
      });

      const dropdownButtons = screen.getAllByRole('button');
      const threeDotsButton = dropdownButtons[0];
      
      fireEvent.click(threeDotsButton);
      fireEvent.click(screen.getByText('Edit'));

      const textarea = screen.getByDisplayValue('Test comment 1');
      fireEvent.change(textarea, { target: { value: '' } });

      const saveButton = screen.getByText('Save Changes');
      expect(saveButton).toBeDisabled();
    });

    test('submits edit form with valid content', async () => {
      auth.updateComment.mockResolvedValue({
        success: true,
        data: { ...mockComments[0], content: 'Updated comment', is_edited: true }
      });

      render(<CommentList postId={1} user={mockUser} />);

      await waitFor(() => {
        expect(screen.getByText('Test comment 1')).toBeInTheDocument();
      });

      const dropdownButtons = screen.getAllByRole('button');
      const threeDotsButton = dropdownButtons[0];
      
      fireEvent.click(threeDotsButton);
      fireEvent.click(screen.getByText('Edit'));

      const textarea = screen.getByDisplayValue('Test comment 1');
      fireEvent.change(textarea, { target: { value: 'Updated comment' } });
      fireEvent.click(screen.getByText('Save Changes'));

      await waitFor(() => {
        expect(auth.updateComment).toHaveBeenCalledWith(1, 1, 'Updated comment', null);
      });
    });

    test('handles edit form submission error', async () => {
      auth.updateComment.mockResolvedValue({
        success: false,
        error: 'Update failed'
      });

      const consoleSpy = jest.spyOn(console, 'error').mockImplementation();

      render(<CommentList postId={1} user={mockUser} />);

      await waitFor(() => {
        expect(screen.getByText('Test comment 1')).toBeInTheDocument();
      });

      const dropdownButtons = screen.getAllByRole('button');
      const threeDotsButton = dropdownButtons[0];
      
      fireEvent.click(threeDotsButton);
      fireEvent.click(screen.getByText('Edit'));

      const textarea = screen.getByDisplayValue('Test comment 1');
      fireEvent.change(textarea, { target: { value: 'Updated comment' } });
      fireEvent.click(screen.getByText('Save Changes'));

      await waitFor(() => {
        expect(consoleSpy).toHaveBeenCalledWith('Failed to update comment:', 'Update failed');
      });

      consoleSpy.mockRestore();
    });

    test('closes edit modal when Cancel is clicked', async () => {
      render(<CommentList postId={1} user={mockUser} />);

      await waitFor(() => {
        expect(screen.getByText('Test comment 1')).toBeInTheDocument();
      });

      const dropdownButtons = screen.getAllByRole('button');
      const threeDotsButton = dropdownButtons[0];
      
      fireEvent.click(threeDotsButton);
      fireEvent.click(screen.getByText('Edit'));

      expect(screen.getByText('Edit Comment')).toBeInTheDocument();
      
      fireEvent.click(screen.getByText('Cancel'));
      expect(screen.queryByText('Edit Comment')).not.toBeInTheDocument();
    });
  });

  describe('Delete Confirmation Flow', () => {
    beforeEach(async () => {
      auth.fetchComments.mockResolvedValue({
        success: true,
        data: mockComments
      });
    });

    test('opens delete confirmation modal when Delete is clicked', async () => {
      render(<CommentList postId={1} user={mockUser} />);

      await waitFor(() => {
        expect(screen.getByText('Test comment 1')).toBeInTheDocument();
      });

      const dropdownButtons = screen.getAllByRole('button');
      const threeDotsButton = dropdownButtons[0];
      
      fireEvent.click(threeDotsButton);
      fireEvent.click(screen.getByText('Delete'));

      expect(screen.getByText('Delete Comment')).toBeInTheDocument();
      expect(screen.getByText('Are you sure you want to delete this comment? This action cannot be undone.')).toBeInTheDocument();
    });

    test('closes delete modal when Cancel is clicked', async () => {
      render(<CommentList postId={1} user={mockUser} />);

      await waitFor(() => {
        expect(screen.getByText('Test comment 1')).toBeInTheDocument();
      });

      const dropdownButtons = screen.getAllByRole('button');
      const threeDotsButton = dropdownButtons[0];
      
      fireEvent.click(threeDotsButton);
      fireEvent.click(screen.getByText('Delete'));

      expect(screen.getByText('Delete Comment')).toBeInTheDocument();
      
      fireEvent.click(screen.getByText('Cancel'));
      expect(screen.queryByText('Delete Comment')).not.toBeInTheDocument();
    });

    test('deletes comment when confirmed', async () => {
      auth.deleteComment.mockResolvedValue({ success: true });

      render(<CommentList postId={1} user={mockUser} />);

      await waitFor(() => {
        expect(screen.getByText('Test comment 1')).toBeInTheDocument();
      });

      const dropdownButtons = screen.getAllByRole('button');
      const threeDotsButton = dropdownButtons[0];
      
      fireEvent.click(threeDotsButton);
      fireEvent.click(screen.getByText('Delete'));
      fireEvent.click(screen.getByText('Delete'));

      await waitFor(() => {
        expect(auth.deleteComment).toHaveBeenCalledWith(1, 1);
        expect(screen.queryByText('Test comment 1')).not.toBeInTheDocument();
      });
    });

    test('handles delete error', async () => {
      auth.deleteComment.mockResolvedValue({
        success: false,
        error: 'Delete failed'
      });

      const consoleSpy = jest.spyOn(console, 'error').mockImplementation();

      render(<CommentList postId={1} user={mockUser} />);

      await waitFor(() => {
        expect(screen.getByText('Test comment 1')).toBeInTheDocument();
      });

      const dropdownButtons = screen.getAllByRole('button');
      const threeDotsButton = dropdownButtons[0];
      
      fireEvent.click(threeDotsButton);
      fireEvent.click(screen.getByText('Delete'));
      fireEvent.click(screen.getByText('Delete'));

      await waitFor(() => {
        expect(consoleSpy).toHaveBeenCalledWith('Failed to delete comment:', 'Delete failed');
      });

      consoleSpy.mockRestore();
    });
  });

  describe('Follow User Functionality', () => {
    beforeEach(async () => {
      auth.fetchComments.mockResolvedValue({
        success: true,
        data: mockComments
      });
    });

    test('calls followUser when Follow is clicked', async () => {
      auth.followUser.mockResolvedValue({ success: true });

      render(<CommentList postId={1} user={mockUser} />);

      await waitFor(() => {
        expect(screen.getByText('Test comment 2')).toBeInTheDocument();
      });

      const dropdownButtons = screen.getAllByRole('button');
      const secondCommentButton = dropdownButtons[1];
      
      fireEvent.click(secondCommentButton);
      fireEvent.click(screen.getByText('Follow'));

      await waitFor(() => {
        expect(auth.followUser).toHaveBeenCalledWith(2);
      });
    });

    test('handles follow error', async () => {
      auth.followUser.mockResolvedValue({
        success: false,
        error: 'Follow failed'
      });

      const consoleSpy = jest.spyOn(console, 'error').mockImplementation();

      render(<CommentList postId={1} user={mockUser} />);

      await waitFor(() => {
        expect(screen.getByText('Test comment 2')).toBeInTheDocument();
      });

      const dropdownButtons = screen.getAllByRole('button');
      const secondCommentButton = dropdownButtons[1];
      
      fireEvent.click(secondCommentButton);
      fireEvent.click(screen.getByText('Follow'));

      await waitFor(() => {
        expect(consoleSpy).toHaveBeenCalledWith('Failed to follow user:', 'Follow failed');
      });

      consoleSpy.mockRestore();
    });
  });

  describe('Edited Indicator', () => {
    test('displays edited indicator for edited comments', async () => {
      auth.fetchComments.mockResolvedValue({
        success: true,
        data: mockComments
      });

      render(<CommentList postId={1} user={mockUser} />);

      await waitFor(() => {
        expect(screen.getByText('• edited')).toBeInTheDocument();
      });
    });

    test('does not display edited indicator for non-edited comments', async () => {
      auth.fetchComments.mockResolvedValue({
        success: true,
        data: [mockComments[0]] // Only the non-edited comment
      });

      render(<CommentList postId={1} user={mockUser} />);

      await waitFor(() => {
        expect(screen.getByText('Test comment 1')).toBeInTheDocument();
        expect(screen.queryByText('• edited')).not.toBeInTheDocument();
      });
    });
  });

  describe('Loading States', () => {
    test('shows loading state during edit submission', async () => {
      auth.updateComment.mockImplementation(() => new Promise(() => {}));

      render(<CommentList postId={1} user={mockUser} />);

      await waitFor(() => {
        expect(screen.getByText('Test comment 1')).toBeInTheDocument();
      });

      const dropdownButtons = screen.getAllByRole('button');
      const threeDotsButton = dropdownButtons[0];
      
      fireEvent.click(threeDotsButton);
      fireEvent.click(screen.getByText('Edit'));

      const textarea = screen.getByDisplayValue('Test comment 1');
      fireEvent.change(textarea, { target: { value: 'Updated comment' } });
      fireEvent.click(screen.getByText('Save Changes'));

      expect(screen.getByText('Saving...')).toBeInTheDocument();
    });

    test('disables buttons during edit submission', async () => {
      auth.updateComment.mockImplementation(() => new Promise(() => {}));

      render(<CommentList postId={1} user={mockUser} />);

      await waitFor(() => {
        expect(screen.getByText('Test comment 1')).toBeInTheDocument();
      });

      const dropdownButtons = screen.getAllByRole('button');
      const threeDotsButton = dropdownButtons[0];
      
      fireEvent.click(threeDotsButton);
      fireEvent.click(screen.getByText('Edit'));

      const textarea = screen.getByDisplayValue('Test comment 1');
      fireEvent.change(textarea, { target: { value: 'Updated comment' } });
      fireEvent.click(screen.getByText('Save Changes'));

      expect(screen.getByText('Cancel')).toBeDisabled();
    });
  });
});