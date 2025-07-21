import React from 'react';
import { SearchIcon, PlusIcon } from 'lucide-react';
import CommunityItem from '../homepage/CommunityItem';

// TODO: Plug in backend API calls when available
function fetchProfileImage() {
  // TODO: Fetch profile image from backend
  return null;
}

function fetchVerifiedBadge() {
  // TODO: Fetch verified badge status from backend
  return null;
}

function fetchFollowers() {
  // TODO: Fetch followers count from backend
  return null;
}

function fetchFollowing() {
  // TODO: Fetch following count from backend
  return null;
}

function fetchProfileStatus() {
  // TODO: Fetch profile status from backend
  return null;
}

function fetchCommunities() {
  // TODO: Fetch communities list from backend
  return null;
}

const ProfileSidebar = () => {
  // TODO: Use useEffect to call fetchProfileImage, fetchVerifiedBadge, fetchFollowers, fetchFollowing, fetchProfileStatus, and fetchCommunities when backend is available

  // This component is a sidebar for user profile 
  return <div className="w-72 flex flex-col gap-6">
    {/* TODO: Replace hardcoded profile image, verified badge, followers, following, and profile status with backend data */}
    <div className="bg-[#101b70] rounded-xl p-4 flex flex-col items-center">
      <div className="relative">
        <div className="w-24 h-24 rounded-full bg-[#ffd700] flex items-center justify-center">
          <img src="https://randomuser.me/api/portraits/men/30.jpg" alt="Profile" className="w-20 h-20 rounded-full" />
          {/* TODO: Show verified badge if user is verified */}
        </div>
        <button className="absolute bottom-0 right-0 bg-[#3f3fd3] p-1 rounded-full">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" stroke="white" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
          </svg>
        </button>
      </div>
      
      <div className="flex justify-between w-full mt-4">
        <div className="text-center">
          <div className="text-xl font-bold">1984</div>
          <div className="text-xs text-[#afafed]">Followers</div>
        </div>
        <div className="text-center">
          <div className="text-xl font-bold">1002</div>
          <div className="text-xs text-[#afafed]">Following</div>
        </div>
      </div>
      
      <div className="mt-4 text-center">
        <h2 className="text-lg font-bold">Zeus</h2>
        <p className="text-sm text-[#afafed]">@zeus</p>
      </div>
      
      <div className="mt-4 text-center text-sm">
        <p>✨ Hello, I'm Zeus. I know everything ✨</p>
        {/* TODO: Replace with profile status from backend */}
      </div>
      
      <button className="mt-4 w-full py-2 bg-[#101b70] border border-[#3f3fd3] rounded-lg text-sm">
        My Profile
      </button>
    </div>

    {/* TODO: Replace hardcoded communities list with backend data */}
    <div className="bg-[#101b70] rounded-xl p-4">
      <div className="flex justify-between items-center mb-3">
        <h3 className="font-bold">Communities</h3>
        <div className="flex gap-2">
          <button className="p-1.5 rounded-full hover:bg-[#3f3fd3]/20">
            <SearchIcon className="w-4 h-4" />
          </button>
          <button className="p-1.5 rounded-full hover:bg-[#3f3fd3]/20">
            <PlusIcon className="w-4 h-4" />
          </button>
        </div>
      </div>
      
      <div className="flex flex-col gap-3">
        {/* TODO: Map communities from backend data */}
        <CommunityItem icon="https://randomuser.me/api/portraits/men/22.jpg" name="UX designers community" memberCount={32} />
        <CommunityItem icon="https://randomuser.me/api/portraits/men/7.jpg" name="Frontend developers" memberCount={12} />
        <CommunityItem icon="https://randomuser.me/api/portraits/men/1.jpg"  name="Frontend developers" memberCount={3} />
        <CommunityItem icon="https://randomuser.me/api/portraits/men/23.jpg"  name="Frontend developers" memberCount={3} />
        <CommunityItem icon="https://randomuser.me/api/portraits/men/4.jpg"  name="Frontend developers" memberCount={3} />
        <CommunityItem icon="https://randomuser.me/api/portraits/men/25.jpg"  name="Frontend developers" memberCount={3} />
      </div>
    </div>
  </div>;
};

export default ProfileSidebar;