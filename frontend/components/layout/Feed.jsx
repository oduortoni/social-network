import React from 'react';
import { ImageIcon, VideoIcon, BarChart2Icon, MoreHorizontalIcon, ThumbsUpIcon, ThumbsDownIcon, MessageCircleIcon, SendIcon } from 'lucide-react';
import UserCircle from '../homepage/UserCircle';
import VerifiedBadge from '../homepage/VerifiedBadge';

const Feed = ({user = null, connectedUsers = []}) => {
  const users = connectedUsers.filter(u => u.user_id != user.id);

  return <div className="flex-1 flex flex-col gap-4">
      <div className="flex overflow-x-auto gap-3 pb-2 cursor-pointer">
        {
          users.map((u, index) => {
              return <UserCircle avatar={u.avatar? u.avatar : ''} name={u.nickname} active={false} highlight="#3f3fd3" key={index} />;
          })
        }
      </div>
      <div className="rounded-xl p-4" style={{ backgroundColor: 'var(--primary-background)' }}>
        <div className="flex items-center gap-3 rounded-xl p-3" style={{ backgroundColor: 'var(--secondary-background)' }}>
          <img src="https://randomuser.me/api/portraits/men/30.jpg" alt="Profile" className="w-10 h-10 rounded-full" />
          <input type="text" placeholder="Tell your friends about your thoughts..." className="bg-transparent flex-1 focus:outline-none text-sm" />
        </div>
        <div className="flex justify-between mt-4 border-t border-[#3f3fd3]/30 pt-3">
          <button className="flex items-center gap-2 text-sm py-1.5 px-3 rounded-lg cursor-pointer hover:bg-[#3f3fd3]/20"
            style={{ color: 'var(--secondary-text)'}}>
            <ImageIcon className="w-4 h-4" />
            <span>Photo</span>
          </button>
          <button className="flex items-center gap-2 text-sm py-1.5 px-3 cursor-pointer rounded-lg hover:bg-[#3f3fd3]/20"
            style={{ color: 'var(--secondary-text)' }}>
            <VideoIcon className="w-4 h-4" />
            <span>Video</span>
          </button>
          <button className="flex items-center gap-2 text-sm py-1.5 px-3 cursor-pointer rounded-lg hover:bg-[#3f3fd3]/20"
            style={{ color: 'var(--secondary-text)' }}>
            <BarChart2Icon className="w-4 h-4" />
            <span>Poll</span>
          </button>
          <button className="flex items-center gap-2 text-sm py-1.5 px-3 cursor-pointer rounded-lg hover:bg-[#3f3fd3]/20"
            style={{ color: 'var(--secondary-text)' }}>
            <SendIcon className="w-4 h-4" />
            <span>Post</span>
          </button>
        </div>
      </div>
      <div className="rounded-xl p-4" style={{ backgroundColor: 'var(--primary-background)' }}>
        <div className="flex justify-between mb-3">
          <div className="flex items-center gap-2">
            <div className="relative">
              <img src="https://randomuser.me/api/portraits/men/42.jpg" alt="Mudreh" className="w-10 h-10 rounded-full" />
              <div className="absolute -bottom-1 -right-1">
                <VerifiedBadge />
              </div>
            </div>
            <div>
              <div className="flex items-center gap-1 cursor-pointer">
                <p className="font-medium">@Muhadrehh</p>
                <VerifiedBadge />
              </div>
              <div className="flex items-center gap-2 cursor-pointer">
                <p className="text-sm font-medium">Mudreh Kumbirai</p>
                <span className="text-xs" style={{ color: 'var(--secondary-text)' }}>• 1 hr ago</span>
              </div>
            </div>
          </div>
          <button>
            <MoreHorizontalIcon className="w-5 h-5 cursor-pointer" />
          </button>
        </div>
        <p className="text-sm mb-3">
          In some cases you may see a third-party client name, which indicates
          the Tweet came from a non-Twitter application.
        </p>
        <div className="rounded-xl overflow-hidden mb-3">
          <img src="/cat.jpg" alt="Post image" className="w-full h-64 object-cover" />
        </div>
        <div className="flex justify-between">
          <div className="flex gap-6 items-center">
            <button className="flex items-center gap-1 cursor-pointer">
              <ThumbsUpIcon className="w-5 h-5 text-blue-400" />
              <span className="text-xs text-white">24 Likes</span>
            </button>
            <button className="flex items-center cursor-pointer gap-1">
              <ThumbsDownIcon className="w-5 h-5 text-red-400" />
              <span className="text-xs text-white">3 Dislikes</span>
            </button>
            <button className="flex items-center gap-1 cursor-pointer">
              <MessageCircleIcon className="w-5 h-5" />
            </button>
          </div>
          <button className="text-black px-4 py-1 rounded-full cursor-pointer text-sm font-medium"
            style={{ backgroundColor: 'var(--primary-accent)', color: 'black' }}>
            #SampleTag
          </button>
        </div>
      </div>
      <div className="rounded-xl p-4" style={{ backgroundColor: 'var(--primary-background)' }}>
        <div className="flex items-center gap-3 rounded-xl p-3" style={{ backgroundColor: 'var(--secondary-background)' }}>
          <img src="https://randomuser.me/api/portraits/men/30.jpg" alt="Profile" className="w-10 h-10 rounded-full" />
          <input type="text" placeholder="Write your comment..." className="bg-transparent flex-1 focus:outline-none text-sm" />
          <button className="cursor-pointer" style={{ color: 'var(--secondary-text)' }}>
            <ImageIcon className="w-5 h-5" />
          </button>
          <button className= "cursor-pointer" style={{ color: 'var(--secondary-text)' }}>
            <SendIcon className="w-5 h-5" />
          </button>
        </div>
      </div>
    </div>;
};

export default Feed;