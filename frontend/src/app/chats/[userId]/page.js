'use client';

import { useRouter, useSearchParams } from 'next/navigation';
import withAuth from '../../../lib/withAuth';
import Header from '../../../components/layout/Header';
import { chatAPI } from '../../../lib/api';
import { ArrowLeftIcon } from 'lucide-react';
import ChatInterface from '../../../components/chat/ChatInterface';

const ChatPage = ({ user, params }) => {
  const router = useRouter();
  const searchParams = useSearchParams();
  const type = searchParams.get('type') || 'private';
  const name = searchParams.get('name') || searchParams.get('nickname');
  const id = params.userId;

  const handleBackToChats = () => {
    router.push('/chats');
  };

  const recipient = {
    id,
    name,
    type,
    nickname: type === 'private' ? name : null,
  };

  return (
    <div className="min-h-screen flex flex-col" style={{ backgroundColor: 'var(--primary-background)' }}>
      <Header user={user} />
      <div className="flex-1 flex flex-col max-w-4xl w-full mx-auto p-6">
        <div className="flex items-center mb-6">
          <button
            onClick={handleBackToChats}
            className="flex items-center gap-2 px-4 py-2 rounded-lg hover:bg-gray-100 transition-colors"
            style={{ color: 'var(--primary-text)' }}
          >
            <ArrowLeftIcon className="w-5 h-5" />
            Back to Chats
          </button>
        </div>
        <h1 className="text-3xl font-bold mb-4" style={{ color: 'var(--primary-text)' }}>
          {name}
        </h1>
        <ChatInterface user={user} recipient={recipient} showSidebar={false} />
      </div>
    </div>
  );
};

export default withAuth(ChatPage);
