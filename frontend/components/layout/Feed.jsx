import React from 'react';
import { ImageIcon, VideoIcon, BarChart2Icon, CalendarIcon, MoreHorizontalIcon, ThumbsUpIcon, ThumbsDownIcon, MessageCircleIcon, ShareIcon, SendIcon } from 'lucide-react';
import UserCircle from '../homepage/UserCircle';
import VerifiedBadge from '../homepage/VerifiedBadge';

const Feed = () => {
  return <div className="flex-1 flex flex-col gap-4">
      <div className="flex overflow-x-auto gap-3 pb-2">
        <UserCircle image="https://randomuser.me/api/portraits/women/22.jpg" name="Amanda" active={false} highlight="#3f3fd3" />
        <UserCircle image="https://randomuser.me/api/portraits/men/22.jpg" name="John" active={false} highlight="#ff4444" />
        <UserCircle image="https://randomuser.me/api/portraits/men/32.jpg" name="Andrew" active={false} highlight="#3f3fd3" />
        <UserCircle image="https://randomuser.me/api/portraits/women/32.jpg" name="Rosaline" active={false} highlight="#00ff00" />
        <UserCircle image="https://randomuser.me/api/portraits/men/42.jpg" name="Mudreh" active={false} highlight="#ff4444" />
        <UserCircle image="https://randomuser.me/api/portraits/women/42.jpg" name="Juliet" active={false} highlight="#3f3fd3" />
        <UserCircle image="https://randomuser.me/api/portraits/men/52.jpg" name="Bob" active={false} highlight="#3f3fd3" />
        <UserCircle image="https://randomuser.me/api/portraits/men/2.jpg" name="Mudreh" active={false} highlight="#ff4444" />
        <UserCircle image="https://randomuser.me/api/portraits/women/4.jpg" name="Juliet" active={false} highlight="#3f3fd3" />
      </div>
      <div className="bg-[#101b70] rounded-xl p-4">
        <div className="flex items-center gap-3">
          <img src="https://randomuser.me/api/portraits/men/32.jpg" alt="Profile" className="w-10 h-10 rounded-full" />
          <input type="text" placeholder="Tell your friends about your thoughts..." className="bg-transparent flex-1 focus:outline-none text-sm" />
        </div>
        <div className="flex justify-between mt-4 border-t border-[#3f3fd3]/30 pt-3">
          <button className="flex items-center gap-2 text-sm text-[#afafed] py-1.5 px-3 rounded-lg hover:bg-[#3f3fd3]/20">
            <ImageIcon className="w-4 h-4" />
            <span>Photo</span>
          </button>
          <button className="flex items-center gap-2 text-sm text-[#afafed] py-1.5 px-3 rounded-lg hover:bg-[#3f3fd3]/20">
            <VideoIcon className="w-4 h-4" />
            <span>Video</span>
          </button>
          <button className="flex items-center gap-2 text-sm text-[#afafed] py-1.5 px-3 rounded-lg hover:bg-[#3f3fd3]/20">
            <BarChart2Icon className="w-4 h-4" />
            <span>Poll</span>
          </button>
          <button className="flex items-center gap-2 text-sm text-[#afafed] py-1.5 px-3 rounded-lg hover:bg-[#3f3fd3]/20">
            <SendIcon className="w-4 h-4" />
            <span>Post</span>
          </button>
        </div>
      </div>
      <div className="bg-[#101b70] rounded-xl p-4">
        <div className="flex justify-between mb-3">
          <div className="flex items-center gap-2">
            <div className="relative">
              <img src="https://randomuser.me/api/portraits/men/42.jpg" alt="Mudreh" className="w-10 h-10 rounded-full" />
              <div className="absolute -bottom-1 -right-1">
                <VerifiedBadge />
              </div>
            </div>
            <div>
              <div className="flex items-center gap-1">
                <p className="font-medium">@Muhadrehh</p>
                <VerifiedBadge />
              </div>
              <div className="flex items-center gap-2">
                <p className="text-sm font-medium">Mudreh Kumbirai</p>
                <span className="text-xs text-[#afafed]">• 1 hr ago</span>
              </div>
            </div>
          </div>
          <button>
            <MoreHorizontalIcon className="w-5 h-5" />
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
            <button className="flex items-center gap-1">
              <ThumbsUpIcon className="w-5 h-5 text-blue-400" />
              <span className="text-xs text-white">24 Likes</span>
            </button>
            <button className="flex items-center gap-1">
              <ThumbsDownIcon className="w-5 h-5 text-red-400" />
              <span className="text-xs text-white">3 Dislikes</span>
            </button>
            <button className="flex items-center gap-1">
              <MessageCircleIcon className="w-5 h-5" />
            </button>
          </div>
          <button className="bg-[#ffd700] text-black px-4 py-1 rounded-full text-sm font-medium">
            #SampleTag
          </button>
        </div>
      </div>
      <div className="bg-[#101b70] rounded-xl p-4">
        <div className="flex items-center gap-3">
          <img src="https://randomuser.me/api/portraits/men/32.jpg" alt="Profile" className="w-10 h-10 rounded-full" />
          <input type="text" placeholder="Write your comment..." className="bg-transparent flex-1 focus:outline-none text-sm" />
          <button className="text-[#afafed]">
            <ImageIcon className="w-5 h-5" />
          </button>
          <button className="text-[#afafed]">
            <SendIcon className="w-5 h-5" />
          </button>
        </div>
      </div>
    </div>;
};

export default Feed;